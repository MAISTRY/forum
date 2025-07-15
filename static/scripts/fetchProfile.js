function loadProfileData() {
    const createdPostsContainer = document.getElementById('Created');
    const likedPostsContainer = document.getElementById('Liked');
    const dislikedPostsContainer = document.getElementById('Disliked');
    const commentsContainer = document.getElementById('Profile');

    // Load user reports for moderators
    loadUserReports();

    // Validate containers
    if (!createdPostsContainer || !likedPostsContainer || !dislikedPostsContainer || !commentsContainer) {
        console.error("One or more containers are missing from the DOM.");
        return;
    }

    // Load moderation request section for normal users
    loadModerationRequestSection();

    // Show loading messages
    [createdPostsContainer, likedPostsContainer, dislikedPostsContainer].forEach(container => {
        container.innerHTML = '<p style="text-align: center">Loading posts...</p>';
    });

    // First get user authentication data
    fetch("/auth/status", {
        method: "GET",
        credentials: "same-origin"
    })
    .then(response => response.json())
    .then(authData => {
        // Then fetch profile data
        return fetch('/Data-Profile', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'X-Requested-With': 'XMLHttpRequest'
            }
        }).then(response => {
            if (!response.ok) {
                throw new Error(`${response.status}: ${response.statusText || 'Unknown Error'}`);
            }
            return response.json();
        }).then(profileData => {
            return { profileData, currentUserId: authData.user_id || 0 };
        });
    })
    .then(data => {
        const { profileData, currentUserId } = data;
        const createdPostsContainer = document.getElementById('Created');
        const likedPostsContainer = document.getElementById('Liked');
        const dislikedPostsContainer = document.getElementById('Disliked');
        
        // Clear the containers first
        createdPostsContainer.innerHTML = '';
        likedPostsContainer.innerHTML = '';
        dislikedPostsContainer.innerHTML = '';
        
        // Loop through the profile data for Created, Liked, and Disliked posts
        ['CreatedPosts', 'LikedPosts', 'DislikedPosts'].forEach(type => {
            const posts = profileData[type];
            const containerMap = {
                'CreatedPosts': createdPostsContainer,
                'LikedPosts': likedPostsContainer,
                'DislikedPosts': dislikedPostsContainer
            };
            const container = containerMap[type];

            // Handle case when posts is null, undefined, or empty
            if (!posts || !Array.isArray(posts) || posts.length === 0) {
                container.innerHTML = '<div class="empty-message" style="text-align: center; padding: 40px; color: #666; font-style: italic;">No posts found</div>';
                return;
            }
            posts.forEach(post => {
                const postElement = document.createElement('div');
                postElement.id = 'Profile-' + post.PostID;
                postElement.className = 'post-card';

                // Categories Section
                const postCategoryContainer = document.createElement('div');
                postCategoryContainer.className = 'post-category';
    
                const categorySpan = document.createElement('span');
                categorySpan.className = 'text-category';
                categorySpan.textContent = 'Categories:';
                postCategoryContainer.appendChild(categorySpan);
    
                if (post.Categories && post.Categories.length > 0) {
                    post.Categories.forEach(CTG => {
                        const categoryElement = document.createElement('a');
                        categoryElement.textContent = CTG;
                        postCategoryContainer.appendChild(categoryElement);
                    });
                } else {
                    const noCategoryMessage = document.createElement('span');
                    noCategoryMessage.textContent = ' None';
                    postCategoryContainer.appendChild(noCategoryMessage);
                }
    
                postElement.insertBefore(postCategoryContainer, postElement.firstChild);
                
    
                // Post Header
                const postHeader = document.createElement('div');
                postHeader.classList.add('post-header');
    
                const postTitle = document.createElement('h3');
                postTitle.classList.add('post-title');
                postTitle.textContent = post.title;
    
                const postMeta = document.createElement('div');
                postMeta.classList.add('post-meta');
                postMeta.textContent = formatDate(post.PostDate);
    
                postHeader.appendChild(postTitle);
                postHeader.appendChild(postMeta);
    
                // Post Content
                const postContent = document.createElement('div');
                postContent.classList.add('post-content');
                postContent.textContent = post.content;
    
                // Post Footer
                const postFooter = document.createElement('div');
                postFooter.classList.add('post-footer');
    
                const buttonsContainer = document.createElement('div');
                buttonsContainer.id = `Profile-post-${post.PostID}`;
                buttonsContainer.classList.add('buttons-contant');
    
                // Like Button
                const likeForm = document.createElement('form');
                likeForm.classList.add('like-form');
                likeForm.onsubmit = (event) => {
                    likeForm.addEventListener('submit', handlePostInteraction(event,"Profile"));
                };
    
                const likeInput = document.createElement('input');
                likeInput.type = 'hidden';
                likeInput.name = 'postId';
                likeInput.value = post.PostID;
    
                const likeButton = document.createElement('button');
                likeButton.type = 'submit';
                likeButton.classList.add('footer-buttons', 'post-button', 'like-buttons');
                likeButton.title = 'Like';
    
                const likeIcon = document.createElement('i');
                likeIcon.classList.add('material-icons');
                likeIcon.textContent = 'thumb_up';
    
                const likeCount = document.createElement('span');
                likeCount.classList.add('likes');
                likeCount.textContent = post.Likes;
    
                likeButton.appendChild(likeIcon);
                likeButton.appendChild(likeCount);
                likeForm.appendChild(likeInput);
                likeForm.appendChild(likeButton);
    
                // Dislike Button
                const dislikeForm = document.createElement('form');
                dislikeForm.classList.add('dislike-form');
                dislikeForm.onclick = (event) => {
                    dislikeForm.addEventListener('submit', handlePostInteraction(event,"Profile"));
                }
    
                const dislikeInput = document.createElement('input');
                dislikeInput.type = 'hidden';
                dislikeInput.name = 'postId';
                dislikeInput.value = post.PostID;
    
                const dislikeButton = document.createElement('button');
                dislikeButton.type = 'submit';
                dislikeButton.classList.add('footer-buttons', 'post-button', 'dislike-buttons');
                dislikeButton.title = 'Dislike';
    
                const dislikeIcon = document.createElement('i');
                dislikeIcon.classList.add('material-icons');
                dislikeIcon.textContent = 'thumb_down';
    
                const dislikeCount = document.createElement('span');
                dislikeCount.classList.add('dislikes');
                dislikeCount.textContent = post.Dislikes;
    
                dislikeButton.appendChild(dislikeIcon);
                dislikeButton.appendChild(dislikeCount);
                dislikeForm.appendChild(dislikeInput);
                dislikeForm.appendChild(dislikeButton);
    
                // Comment Button
                const commentButton = document.createElement('a');
                commentButton.classList.add('footer-buttons', 'post-button', 'comment-buttons');
                commentButton.title = 'comment';
                commentButton.onclick = () => {
                    commentButton.addEventListener('click', commentShow("Profile",post.PostID));
                }
    
                const commentIcon = document.createElement('i');
                commentIcon.classList.add('material-icons');
                commentIcon.textContent = 'comment';
    
                const commentCount = document.createElement('span');
                commentCount.classList.add('comments');
                commentCount.textContent = post.CmtCount;
    
                commentButton.appendChild(commentIcon);
                commentButton.appendChild(commentCount);
    
                // User Info
                const postUser = document.createElement('div');
                postUser.classList.add('footer-buttons', 'post-user');
                postUser.textContent = `@${post.username}`;
    
                buttonsContainer.appendChild(likeForm);
                buttonsContainer.appendChild(dislikeForm);
                buttonsContainer.appendChild(commentButton);

                // Add edit/delete buttons for created posts only
                if (type === 'CreatedPosts' && currentUserId === post.UserID) {
                    // Edit button
                    const editButton = document.createElement('button');
                    editButton.classList.add('footer-buttons', 'post-button', 'edit-btn');
                    editButton.title = 'Edit Post';
                    editButton.onclick = () => editPost(post.PostID);

                    const editIcon = document.createElement('i');
                    editIcon.classList.add('material-icons');
                    editIcon.textContent = 'edit';
                    editButton.appendChild(editIcon);
                    buttonsContainer.appendChild(editButton);

                    // Delete button
                    const deleteButton = document.createElement('button');
                    deleteButton.classList.add('footer-buttons', 'post-button', 'delete-btn');
                    deleteButton.title = 'Delete Post';
                    deleteButton.onclick = () => deletePost(post.PostID);

                    const deleteIcon = document.createElement('i');
                    deleteIcon.classList.add('material-icons');
                    deleteIcon.textContent = 'delete';
                    deleteButton.appendChild(deleteIcon);
                    buttonsContainer.appendChild(deleteButton);
                }
    
                postFooter.appendChild(buttonsContainer);
                postFooter.appendChild(postUser);
    
                // Add Post Header, Content, and Footer
                postElement.appendChild(postHeader);
                postElement.appendChild(postContent);
    
                // Post Image (if exists)
                if (post.imagePath) {
                    const postImageContainer = document.createElement('div');
                    postImageContainer.classList.add('post-image');
                    const postImage = document.createElement('img');
                    postImage.src = post.imagePath;
                    postImage.alt = 'Post Image';
                    postImage.classList.add('image-size');
                    postImageContainer.appendChild(postImage);
                    postElement.appendChild(postImageContainer);
                }
    
                postElement.appendChild(postFooter);
    
                // Comments Section
                const commentsContainer = document.createElement('div');
                commentsContainer.className = 'comments-section';
                commentsContainer.style.display = 'none';
                commentsContainer.id = `Profile-comments-${post.PostID}`;
                postElement.appendChild(commentsContainer);
    
                // Append post to the correct container (Created, Liked, or Disliked)
                if (type === 'CreatedPosts') {
                    createdPostsContainer.appendChild(postElement);
                } else if (type === 'LikedPosts') {
                    likedPostsContainer.appendChild(postElement);
                } else if (type === 'DislikedPosts') {
                    dislikedPostsContainer.appendChild(postElement);
                }
            });
        });
    
    })
    .catch(error => {
        console.error(error);
    });
}

// Load moderation request section
async function loadModerationRequestSection() {
    try {
        // Get user privilege level
        const authResponse = await fetch("/auth/status", {
            method: "GET",
            credentials: "same-origin"
        });
        const authData = await authResponse.json();

        const moderationSection = document.getElementById('moderation-request-section');
        const statusDiv = document.getElementById('moderation-status');

        if (!moderationSection || !statusDiv) return;

        // Only show for normal users (privilege 1)
        if (authData.authenticated && authData.privilege === 1) {
            moderationSection.style.display = 'block';

            // Check if user already has a pending request
            checkModerationRequestStatus();
        } else {
            moderationSection.style.display = 'none';
        }
    } catch (error) {
        console.error('Error loading moderation request section:', error);
    }
}

// Check moderation request status
async function checkModerationRequestStatus() {
    try {
        const response = await fetch('/Data-AdminModerationRequests', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'X-Requested-With': 'XMLHttpRequest'
            }
        });

        if (response.ok) {
            const requests = await response.json();
            const statusDiv = document.getElementById('moderation-status');
            const requestBtn = document.getElementById('request-moderation-btn');

            // Find current user's request
            const userRequest = requests.find(req => req.Status === 'pending');

            if (userRequest) {
                statusDiv.innerHTML = '<p style="color: orange;">Your moderation request is pending review.</p>';
                requestBtn.disabled = true;
                requestBtn.textContent = 'Request Pending';
            } else {
                statusDiv.innerHTML = '';
                requestBtn.disabled = false;
                requestBtn.textContent = 'Request Moderator Role';
            }
        }
    } catch (error) {
        console.error('Error checking moderation request status:', error);
    }
}

// Request moderation function
async function requestModeration() {
    try {
        const response = await fetch('/Data-CreateModerationRequest', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'X-Requested-With': 'XMLHttpRequest'
            }
        });

        if (response.ok) {
            const result = await response.json();
            if (result.success) {
                const statusDiv = document.getElementById('moderation-status');
                const requestBtn = document.getElementById('request-moderation-btn');

                statusDiv.innerHTML = '<p style="color: green;">Moderation request submitted successfully!</p>';
                requestBtn.disabled = true;
                requestBtn.textContent = 'Request Pending';
            } else {
                alert('Failed to submit request: ' + (result.message || 'Unknown error'));
            }
        } else {
            const errorText = await response.text();
            alert('Failed to submit request: ' + errorText);
        }
    } catch (error) {
        console.error('Error requesting moderation:', error);
        alert('Error submitting request. Please try again.');
    }
}

// Load user reports (for moderators)
async function loadUserReports() {
    try {
        const response = await fetch('/Data-UserReports', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'X-Requested-With': 'XMLHttpRequest'
            }
        });

        if (!response.ok) {
            // If unauthorized, user is not a moderator, so hide the section
            if (response.status === 401) {
                return;
            }
            throw new Error('Failed to load user reports');
        }

        const reports = await response.json();
        displayUserReports(reports);
    } catch (error) {
        console.error('Error loading user reports:', error);
        const container = document.getElementById('user-reports');
        if (container) {
            container.innerHTML = '<div class="empty-message">Error loading reports</div>';
        }
    }
}

// Display user reports
function displayUserReports(reports) {
    const container = document.getElementById('user-reports');
    if (!container) return;

    if (!reports || !Array.isArray(reports) || reports.length === 0) {
        container.innerHTML = '<div class="empty-message">No reports found</div>';
        return;
    }

    const reportsHTML = reports.map(report => {
        // Truncate content if too long
        const truncatedContent = report.PostContent.length > 100
            ? report.PostContent.substring(0, 100) + '...'
            : report.PostContent;

        let statusClass = report.Status;
        let statusText = report.Status.charAt(0).toUpperCase() + report.Status.slice(1);

        let adminResponseHTML = '';
        if (report.AdminResponse && report.AdminResponse.trim() !== '') {
            adminResponseHTML = `
                <div class="admin-response">
                    <strong>Admin Response:</strong> ${report.AdminResponse}
                    ${report.ResponseDate ? `<div class="response-date">Responded on: ${formatDate(report.ResponseDate)}</div>` : ''}
                </div>
            `;
        }

        return `
            <div class="user-report-item">
                <div class="report-header">
                    <div class="report-info">
                        <div class="report-post-title"><strong>${report.PostTitle}</strong></div>
                        <div class="report-post-author">by @${report.PostAuthor}</div>
                    </div>
                    <div class="report-date">${formatDate(report.ReportDate)}</div>
                </div>
                <div class="report-content">${truncatedContent}</div>
                <div class="report-reason"><strong>Reason:</strong> ${report.Reason}</div>
                <div class="report-status ${statusClass}">${statusText}</div>
                ${adminResponseHTML}
            </div>
        `;
    }).join('');

    container.innerHTML = reportsHTML;
}

// Helper function to format dates (if not already available)
function formatDate(dateString) {
    if (!dateString) return '';
    const date = new Date(dateString);
    return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'});
}

// Make function globally available
window.requestModeration = requestModeration;