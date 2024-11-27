package middleware

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

type Message struct {
	Status string `json:"status"`
	Body   string `json:"body"`
}

// RateLimiter is a middleware function that implements a rate limiter for HTTP requests.
// It limits the number of requests from a single client within a specified time frame.
// If a client exceeds the rate limit, they will be temporarily blocked from making further requests.
//
// The rate limiter uses a token bucket algorithm, where tokens are generated at a constant rate.
// Clients are assigned a burst limit of tokens, and they can only make requests when they have tokens available.
// If a client exceeds the burst limit, they will be temporarily blocked.
//
// The rate limiter maintains a map of clients, where each client's IP address is used as the key.
// The map is protected by a mutex to ensure thread safety.
//
// The rate limiter has the following configurable constants:
// - rateLimit: The maximum number of requests allowed per second.
// - burstLimit: The maximum number of requests allowed within a single time frame.
// - cleanupInterval: The interval at which the rate limiter cleans up expired clients from the map.
// - clientTimeout: The maximum duration of time a client can be inactive before being removed from the map.
// - blockDuration: The duration of time a client will be temporarily blocked if they exceed the rate limit.
//
// The middleware function returns an http.Handler that can be used as a middleware in an HTTP server.
func RateLimiter(next http.Handler) http.Handler {
	type Client struct {
		tokens        float64
		lastTimestamp time.Time
		blockUntil    time.Time
	}

	var (
		mtx     sync.Mutex
		clients = make(map[string]*Client)
	)

	const (
		rateLimit       = 100
		burstLimit      = 100
		cleanupInterval = time.Minute
		clientTimeout   = 2 * time.Minute
		blockDuration   = 1 * time.Minute
	)

	go func() {
		for {
			time.Sleep(cleanupInterval)
			mtx.Lock()
			for ip, client := range clients {
				if time.Since(client.lastTimestamp) > clientTimeout {
					delete(clients, ip)
				}
			}
			mtx.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ip string

		if ip = r.Header.Get("X-Real-IP"); ip == "" {
			if ip = r.Header.Get("X-Forwarded-For"); ip == "" {
				ip, _, _ = net.SplitHostPort(r.RemoteAddr)
			}
		}

		if net.ParseIP(ip) == nil {
			http.Error(w, "Invalid IP address", http.StatusBadRequest)
			return
		}

		now := time.Now()

		mtx.Lock()
		client, exists := clients[ip]
		if !exists {
			client = &Client{
				tokens:        burstLimit,
				lastTimestamp: now,
			}
			clients[ip] = client
		} else {
			if client.blockUntil.After(now) {
				mtx.Unlock()
				message := Message{
					Status: "error",
					Body:   "You are temporarily blocked due to excessive requests.",
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				if err := json.NewEncoder(w).Encode(&message); err != nil {
					log.Printf("Failed to encode JSON response: %v", err)
				}
				return
			}

			elapsed := now.Sub(client.lastTimestamp).Seconds()
			client.tokens += elapsed * rateLimit
			if client.tokens > burstLimit {
				client.tokens = burstLimit
			}
			client.lastTimestamp = now
		}

		if client.tokens >= 1.0 {
			client.tokens -= 1.0
			mtx.Unlock()
			next.ServeHTTP(w, r)
			return
		} else {
			client.blockUntil = now.Add(blockDuration)
			mtx.Unlock()

			message := Message{
				Status: "error",
				Body:   "You have exceeded the rate limit and are temporarily blocked.",
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			if err := json.NewEncoder(w).Encode(&message); err != nil {
				log.Printf("Failed to encode JSON response: %v", err)
			}
			return
		}
	})
}

// package middleware

// import (
// 	"database/sql"
// 	"encoding/json"
// 	"errors"
// 	"log"
// 	"net"
// 	"net/http"
// 	"sync"
// 	"time"
// )

// const (
// 	getUserIDQuery = `SELECT user_id FROM Session WHERE session_id = ?`
// )

// // Message represents a JSON response message.
// type Message struct {
// 	Status string `json:"status"`
// 	Body   string `json:"body"`
// }

// // RateLimiter returns a middleware function that rate-limits requests based on a key (user ID or IP address).
// func RateLimiter(getKey func(*http.Request) (string, error)) func(http.Handler) http.Handler {
// 	type Client struct {
// 		tokens        float64
// 		lastTimestamp time.Time
// 		blockUntil    time.Time
// 	}

// 	var (
// 		mtx     sync.Mutex
// 		clients = make(map[string]*Client)
// 	)

// 	const (
// 		rateLimit       = 5.0
// 		burstLimit      = 100.0
// 		cleanupInterval = time.Minute
// 		clientTimeout   = 2 * time.Minute
// 		blockDuration   = 1 * time.Minute
// 	)

// 	go func() {
// 		for {
// 			time.Sleep(cleanupInterval)
// 			mtx.Lock()
// 			for key, client := range clients {
// 				if time.Since(client.lastTimestamp) > clientTimeout {
// 					delete(clients, key)
// 				}
// 			}
// 			mtx.Unlock()
// 		}
// 	}()

// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			key, err := getKey(r)
// 			if err != nil {
// 				http.Error(w, "Unable to determine client identity", http.StatusBadRequest)
// 				return
// 			}

// 			now := time.Now()

// 			mtx.Lock()
// 			client, exists := clients[key]
// 			if !exists {
// 				client = &Client{
// 					tokens:        burstLimit,
// 					lastTimestamp: now,
// 				}
// 				clients[key] = client
// 			} else {
// 				if client.blockUntil.After(now) {
// 					mtx.Unlock()
// 					message := Message{
// 						Status: "error",
// 						Body:   "You are temporarily blocked due to excessive requests.",
// 					}

// 					w.Header().Set("Content-Type", "application/json")
// 					w.WriteHeader(http.StatusTooManyRequests)
// 					if err := json.NewEncoder(w).Encode(&message); err != nil {
// 						log.Printf("Failed to encode JSON response: %v", err)
// 					}
// 					return
// 				}

// 				elapsed := now.Sub(client.lastTimestamp).Seconds()
// 				client.tokens += elapsed * rateLimit
// 				if client.tokens > burstLimit {
// 					client.tokens = burstLimit
// 				}
// 				client.lastTimestamp = now
// 			}

// 			if client.tokens >= 1.0 {
// 				client.tokens -= 1.0
// 				mtx.Unlock()
// 				next.ServeHTTP(w, r)
// 				return
// 			} else {
// 				client.blockUntil = now.Add(blockDuration)
// 				mtx.Unlock()

// 				message := Message{
// 					Status: "error",
// 					Body:   "You have exceeded the rate limit and are temporarily blocked.",
// 				}

// 				w.Header().Set("Content-Type", "application/json")
// 				w.WriteHeader(http.StatusTooManyRequests)
// 				if err := json.NewEncoder(w).Encode(&message); err != nil {
// 					log.Printf("Failed to encode JSON response: %v", err)
// 				}
// 				return
// 			}
// 		})
// 	}
// }

// // getUserIDByCookie extracts the user ID from the session cookie.
// func getUserIDByCookie(r *http.Request, db *sql.DB) (string, error) {
// 	yummyCookie, err := r.Cookie("sessionID")
// 	if err != nil {
// 		return "", errors.New("error getting cookie")
// 	}

// 	stmt, err := db.Prepare(getUserIDQuery)
// 	if err != nil {
// 		return "", errors.New("error preparing statement")
// 	}
// 	defer stmt.Close()

// 	var userID string

// 	err = stmt.QueryRow(yummyCookie.Value).Scan(&userID)
// 	if err != nil {
// 		return "", errors.New("error getting user ID")
// 	}

// 	return userID, nil
// }

// // getKeyFunction returns a function that extracts the key for rate limiting.
// func getKeyFunction(db *sql.DB) func(*http.Request) (string, error) {
// 	return func(r *http.Request) (string, error) {
// 		// Attempt to get user ID from cookie
// 		userID, err := getUserIDByCookie(r, db)
// 		if err == nil && userID != "" {
// 			return userID, nil
// 		}

// 		// Fallback to IP address
// 		var ip string
// 		if ip = r.Header.Get("X-Real-IP"); ip == "" {
// 			if ip = r.Header.Get("X-Forwarded-For"); ip == "" {
// 				ip, _, _ = net.SplitHostPort(r.RemoteAddr)
// 			}
// 		}

// 		if net.ParseIP(ip) == nil {
// 			return "", errors.New("invalid IP address")
// 		}
// 		return ip, nil
// 	}
// }
