package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"

	// hndls "forum/handlers"
	"forum/utils"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	googleAuthURL  = "https://accounts.google.com/o/oauth2/v2/auth"
	googleTokenURL = "https://oauth2.googleapis.com/token"
	githubAuthURL  = "https://github.com/login/oauth/authorize"
	githubTokenURL = "https://github.com/login/oauth/access_token"
)
const (
	InsertSessionQuery = `INSERT INTO Session (session_id, user_id, created_at, expiry_date, ip_address) VALUES (?,?,?,?,?)`
)

// HandleGoogleLogin initiates the OAuth2 flow for Google authentication.
// It constructs a URL to the Google OAuth2 authorization endpoint with the required parameters.
// After the user authorizes the application, Google will redirect the user back to the specified redirect URI.
// The function then redirects the user's browser to the constructed authorization URL.
//
// Parameters:
// - w: http.ResponseWriter to write the response.
// - r: *http.Request containing the HTTP request data.
func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	params := url.Values{}
	params.Add("client_id", GoogleClientID)
	params.Add("redirect_uri", GoogleRedirectURL)
	params.Add("response_type", "code")
	params.Add("scope", "https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile")

	authURL := googleAuthURL + "?" + params.Encode()
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// HandleGitHubLogin initiates the OAuth2 flow for GitHub authentication.
// It constructs a URL to the GitHub OAuth2 authorization endpoint with the required parameters.
// After the user authorizes the application, GitHub will redirect the user back to the specified redirect URI.
// The function then redirects the user's browser to the constructed authorization URL.
//
// Parameters:
// - w: http.ResponseWriter to write the HTTP response.
// - r: *http.Request containing the HTTP request data.
func HandleGitHubLogin(w http.ResponseWriter, r *http.Request) {
	params := url.Values{}
	params.Add("client_id", GitHubClientID)
	params.Add("redirect_uri", GitHubRedirectURL)
	params.Add("scope", "user:email")

	authURL := githubAuthURL + "?" + params.Encode()
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// HandleOAuthCallback processes the callback from OAuth providers (Google or GitHub).
// It exchanges the authorization code for an access token, retrieves user information,
// and either logs in an existing user or redirects to registration for a new user.
//
// Parameters:
//   - w: http.ResponseWriter to write the HTTP response.
//   - r: *http.Request containing the HTTP request data, including the OAuth code and provider.
//
// The function doesn't return any value, but it writes to the ResponseWriter:
//   - In case of errors, it sends appropriate HTTP error responses.
//   - For existing users, it sets a session cookie and redirects to the home page.
//   - For new users, it redirects to the registration page with pre-filled user information.
func HandleOAuthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	provider := r.URL.Query().Get("provider")

	var userInfo map[string]interface{}
	var err error

	switch provider {
	case "google":
		token, err := exchangeGoogleToken(code)
		if err != nil {
			http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
			return
		}
		userInfo, err = getGoogleUserInfo(token)
	case "github":
		token, err := exchangeGitHubToken(code)
		if err != nil {
			http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
			return
		}
		userInfo, err = getGitHubUserInfo(token)
		// Add email fetch for GitHub
		if userInfo["email"] == nil {
			// Fetch email separately for GitHub
			email, err := getGitHubEmail(token)
			if err != nil {
				http.Error(w, "Failed to get email", http.StatusInternalServerError)
				return
			}
			userInfo["email"] = email
		}
	}

	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}

	email := userInfo["email"].(string)
	exists, err := checkEmailExists(email)

	if exists {
		handleExistingUser(w, r, email)
	} else {
		redirectToRegistration(w, r, userInfo)
	}
}

// exchangeGoogleToken exchanges an authorization code for a Google OAuth2 access token.
// It sends a POST request to Google's token endpoint with the necessary parameters.
//
// Parameters:
//   - code: A string representing the authorization code received from Google's OAuth2 flow.
//
// Returns:
//   - string: The access token obtained from Google, which can be used to access Google APIs.
//   - error: An error if the token exchange fails, or nil if successful.
func exchangeGoogleToken(code string) (string, error) {
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", GoogleClientID)
	data.Set("client_secret", GoogleClientSecret)
	data.Set("redirect_uri", GoogleRedirectURL)
	data.Set("grant_type", "authorization_code")

	resp, err := http.PostForm(googleTokenURL, data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.AccessToken, nil
}

// exchangeGitHubToken exchanges an authorization code for a GitHub OAuth2 access token.
// It sends a POST request to GitHub's token endpoint with the necessary parameters.
//
// Parameters:
//   - code: A string representing the authorization code received from GitHub's OAuth2 flow.
//
// Returns:
//   - string: The access token obtained from GitHub, which can be used to access GitHub APIs.
//   - error: An error if the token exchange fails, or nil if successful.
func exchangeGitHubToken(code string) (string, error) {
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", GitHubClientID)
	data.Set("client_secret", GitHubClientSecret)
	data.Set("redirect_uri", GitHubRedirectURL)

	req, _ := http.NewRequest("POST", githubTokenURL, strings.NewReader(data.Encode()))
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.AccessToken, nil
}

// ! helper functions

// getGoogleUserInfo retrieves user information from Google's API using the provided access token.
//
// Parameters:
//   - accessToken: A string representing the OAuth2 access token obtained from Google.
//
// Returns:
//   - map[string]interface{}: A map containing the user's information retrieved from Google.
//     The keys in this map correspond to the fields in Google's user info response.
//   - error: An error if the API request fails or if there's an issue decoding the response.
//     Returns nil if the operation is successful.
func getGoogleUserInfo(accessToken string) (map[string]interface{}, error) {
	req, _ := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&userInfo); err != nil {
		return nil, err
	}
	return userInfo, nil
}

// getGitHubUserInfo retrieves user information from GitHub's API using the provided access token.
//
// Parameters:
//   - accessToken: A string representing the OAuth2 access token obtained from GitHub.
//
// Returns:
//   - map[string]interface{}: A map containing the user's information retrieved from GitHub.
//     The keys in this map correspond to the fields in GitHub's user info response.
//   - error: An error if the API request fails or if there's an issue decoding the response.
//     Returns nil if the operation is successful.
func getGitHubUserInfo(accessToken string) (map[string]interface{}, error) {
	req, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	req.Header.Set("Authorization", "token "+accessToken)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&userInfo); err != nil {
		return nil, err
	}
	return userInfo, nil
}

// checkEmailExists queries the database to determine if a user with the given email exists.
//
// Parameters:
//   - email: A string representing the email address to check in the database.
//
// Returns:
//   - bool: true if a user with the given email exists, false otherwise.
//   - error: An error if the database query fails, or nil if the operation is successful.
func checkEmailExists(email string) (bool, error) {
	db, err := sql.Open("sqlite3", "meow.db")
	if err != nil {
		return false, err
	}
	defer db.Close()

	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM User WHERE email = ?)", email).Scan(&exists)
	return exists, err
}

// redirectToRegistration redirects the user to the registration page with pre-filled user information.
// This function is typically used after OAuth authentication when a new user needs to complete registration.
//
// Parameters:
//   - w: http.ResponseWriter to write the HTTP response for the redirection.
//   - r: *http.Request containing the original HTTP request data.
//   - userInfo: map[string]interface{} containing user information obtained from the OAuth provider.
//     Expected to contain at least an "email" key, and optionally a "name" key.
//
// The function doesn't return any value, but it performs an HTTP redirection to the registration page,
// appending user information (email, firstname, lastname) as URL query parameters.
func redirectToRegistration(w http.ResponseWriter, r *http.Request, userInfo map[string]interface{}) {
	data := url.Values{}
	data.Set("email", userInfo["email"].(string))

	// Handle different field names for GitHub vs Google
	firstName := ""
	lastName := ""
	if name, ok := userInfo["name"].(string); ok {
		names := strings.Split(name, " ")
		if len(names) > 0 {
			firstName = names[0]
			if len(names) > 1 {
				lastName = strings.Join(names[1:], " ")
			}
		}
	}

	data.Set("firstname", firstName)
	data.Set("lastname", lastName)

	http.Redirect(w, r, "/register?"+data.Encode(), http.StatusTemporaryRedirect)
}

// handleExistingUser processes authentication for an existing user.
// It retrieves the user's ID from the database, generates a new session token,
// creates a session in the database, sets a session cookie, and redirects to the home page.
//
// Parameters:
//   - w: http.ResponseWriter to write the HTTP response.
//   - r: *http.Request containing the HTTP request data.
//   - email: string representing the user's email address.
//
// The function doesn't return any value, but it performs the following actions:
//   - Sets a session cookie for the authenticated user.
//   - Redirects the user to the home page upon successful authentication.
//   - Writes HTTP error responses to w in case of any errors during the process.
func handleExistingUser(w http.ResponseWriter, r *http.Request, email string) {
	db, err := sql.Open("sqlite3", "meow.db")
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var userID string
	err = db.QueryRow("SELECT UserID FROM User WHERE email = ?", email).Scan(&userID)
	if err != nil {
		http.Error(w, "Error fetching user", http.StatusInternalServerError)
		return
	}

	sessionToken, err := utils.GenerateSessionToken()
	if err != nil {
		http.Error(w, "Error generating session", http.StatusInternalServerError)
		return
	}

	expiryDate := time.Now().Add(72 * time.Hour)
	ipAddr := utils.GetIP(r)

	_, err = db.Exec(InsertSessionQuery, sessionToken, userID, time.Now(), expiryDate, ipAddr)
	if err != nil {
		http.Error(w, "Error creating session", http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:     "sessionID",
		Value:    sessionToken,
		Expires:  expiryDate,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

// getGitHubEmail retrieves the primary email address associated with the authenticated GitHub user.
// It sends a GET request to the GitHub API endpoint for user emails, including the access token
// to authenticate the request. The function then parses the JSON response, extracts the primary email,
// and returns it along with any encountered errors.
//
// Parameters:
//   - accessToken string: The OAuth2 access token obtained from GitHub during the OAuth2 flow.
//     This token is used to authenticate the API request and authorize access to the user's email.
//
// Return:
//   - string: The primary email address associated with the authenticated GitHub user.
//   - error: An error if the API request fails, the JSON response cannot be parsed, or no primary email is found.
//     Returns nil if the operation is successful.
func getGitHubEmail(accessToken string) (string, error) {
	req, _ := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", "token "+accessToken)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}

	body, _ := io.ReadAll(response.Body)
	if err := json.Unmarshal(body, &emails); err != nil {
		return "", err
	}

	for _, email := range emails {
		if email.Primary {
			return email.Email, nil
		}
	}

	return "", fmt.Errorf("no primary email found")
}
