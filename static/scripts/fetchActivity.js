// Activity page functionality
let notificationCount = 0;

// Load activity page data
function loadActivityData() {
    loadNotifications();
    loadUserActivity();
}

// Load notifications
async function loadNotifications() {
    const notificationsContainer = document.getElementById('notifications-container');
    
    if (!notificationsContainer) {
        console.error("Notifications container not found");
        return;
    }

    try {
        const response = await fetch('/Data-Notifications', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'X-Requested-With': 'XMLHttpRequest'
            }
        });

        if (!response.ok) {
            throw new Error(`${response.status}: ${response.statusText}`);
        }

        const notifications = await response.json();
        displayNotifications(notifications);

        // Get unread count for badge
        loadNotificationCount();
        
    } catch (error) {
        console.error('Error loading notifications:', error);
        notificationsContainer.innerHTML = '<div class="error-message">Failed to load notifications</div>';
    }
}

// Display notifications
function displayNotifications(notifications) {
    const container = document.getElementById('notifications-container');
    
    if (!notifications || notifications.length === 0) {
        container.innerHTML = '<div class="empty-message">No new notifications</div>';
        return;
    }

    let notificationsHTML = '';
    notifications.forEach(notification => {
        const notificationHTML = createNotificationHTML(notification);
        notificationsHTML += notificationHTML;
    });

    container.innerHTML = notificationsHTML;
}

// Create HTML for a single notification
function createNotificationHTML(notification) {
    const timeAgo = formatTimeAgo(notification.created_at);
    let message = '';
    let icon = '';

    switch (notification.notification_type) {
        case 'PostLike':
            message = `${notification.username} liked your post "${notification.post_title}"`;
            icon = 'thumb_up';
            break;
        case 'PostDislike':
            message = `${notification.username} disliked your post "${notification.post_title}"`;
            icon = 'thumb_down';
            break;
        case 'Comment':
            message = `${notification.username} commented on your post "${notification.post_title}"`;
            icon = 'comment';
            break;
        case 'CommentLike':
            message = `${notification.username} liked your comment`;
            icon = 'thumb_up';
            break;
        case 'CommentDislike':
            message = `${notification.username} disliked your comment`;
            icon = 'thumb_down';
            break;
        default:
            message = `${notification.username} interacted with your content`;
            icon = 'notifications';
    }

    const isRead = notification.is_read;
    const readClass = isRead ? 'notification-read' : 'notification-unread';

    return `
        <div class="notification-item ${readClass}" data-notification-id="${notification.notification_id}">
            <div class="notification-icon">
                <i class="material-icons">${icon}</i>
            </div>
            <div class="notification-content">
                <div class="notification-message">${message}</div>
                <div class="notification-time">${timeAgo}</div>
            </div>
            <div class="notification-actions">
                ${!isRead ? `<button class="mark-read-btn" onclick="markAsRead(${notification.notification_id})">
                    <i class="material-icons">done</i>
                </button>` : '<span class="read-indicator">Read</span>'}
            </div>
        </div>
    `;
}

// Mark notification as read
async function markAsRead(notificationId) {
    console.log('Marking notification as read:', notificationId);

    try {
        const response = await fetch('/Data-MarkAsRead', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: `notificationID=${encodeURIComponent(notificationId)}`
        });

        console.log('Mark as read response status:', response.status);

        if (response.ok) {
            console.log('Successfully marked notification as read');

            // Update the notification visual state instead of removing it
            const notificationElement = document.querySelector(`[data-notification-id="${notificationId}"]`);
            if (notificationElement) {
                // Change visual state to read
                notificationElement.classList.remove('notification-unread');
                notificationElement.classList.add('notification-read');

                // Replace the mark-read button with read indicator
                const actionsDiv = notificationElement.querySelector('.notification-actions');
                if (actionsDiv) {
                    actionsDiv.innerHTML = '<span class="read-indicator">Read</span>';
                }

                console.log('Updated notification visual state to read');
            } else {
                console.error('Could not find notification element to update');
            }

            // Update notification badge count
            loadNotificationCount();
        } else {
            const errorText = await response.text();
            console.error('Failed to mark notification as read:', response.status, errorText);
            alert('Failed to mark notification as read');
        }
    } catch (error) {
        console.error('Error marking notification as read:', error);
        alert('Error marking notification as read');
    }
}

// Update notification badge
function updateNotificationBadge(count) {
    const badge = document.getElementById('notification-badge');
    if (badge) {
        notificationCount = count;
        if (count > 0) {
            badge.textContent = count;
            badge.style.display = 'flex';
        } else {
            badge.style.display = 'none';
        }
    }
}

// Load user activity data
async function loadUserActivity() {
    try {
        const response = await fetch('/Data-Profile', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'X-Requested-With': 'XMLHttpRequest'
            }
        });

        if (!response.ok) {
            throw new Error(`${response.status}: ${response.statusText}`);
        }

        const profileData = await response.json();
        displayUserActivity(profileData);
        
    } catch (error) {
        console.error('Error loading user activity:', error);
        displayActivityError();
    }
}

// Display user activity
function displayUserActivity(profileData) {
    displayActivityPosts('activity-created-posts', profileData.CreatedPosts, 'No posts created yet');
    displayActivityPosts('activity-liked-posts', profileData.LikedPosts, 'No posts liked yet');
    displayActivityPosts('activity-disliked-posts', profileData.DislikedPosts, 'No posts disliked yet');
    
    // Load comments separately
    loadUserComments();
}

// Display posts in activity sections
function displayActivityPosts(containerId, posts, emptyMessage) {
    const container = document.getElementById(containerId);
    if (!container) return;

    if (!posts || posts.length === 0) {
        container.innerHTML = `<div class="empty-message">${emptyMessage}</div>`;
        return;
    }

    let postsHTML = '';
    posts.forEach(post => {
        postsHTML += createActivityPostHTML(post);
    });

    container.innerHTML = postsHTML;
}

// Create HTML for activity post
function createActivityPostHTML(post) {
    const timeAgo = formatTimeAgo(post.PostDate);
    
    return `
        <div class="activity-post-item">
            <div class="activity-post-header">
                <h4 class="activity-post-title">${escapeHtml(post.title)}</h4>
                <span class="activity-post-time">${timeAgo}</span>
            </div>
            <div class="activity-post-content">
                ${escapeHtml(post.content.substring(0, 150))}${post.content.length > 150 ? '...' : ''}
            </div>
            <div class="activity-post-stats">
                <span><i class="material-icons">thumb_up</i> ${post.Likes}</span>
                <span><i class="material-icons">thumb_down</i> ${post.Dislikes}</span>
                <span><i class="material-icons">comment</i> ${post.CmtCount}</span>
            </div>
        </div>
    `;
}

// Load user comments
async function loadUserComments() {
    const container = document.getElementById('activity-comments');
    if (!container) return;

    try {
        const response = await fetch('/Data-Activity', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'X-Requested-With': 'XMLHttpRequest'
            }
        });

        if (!response.ok) {
            throw new Error(`${response.status}: ${response.statusText}`);
        }

        const activityData = await response.json();
        displayUserComments(activityData.user_comments);

    } catch (error) {
        console.error('Error loading user comments:', error);
        container.innerHTML = '<div class="error-message">Failed to load comments</div>';
    }
}

// Display user comments
function displayUserComments(comments) {
    const container = document.getElementById('activity-comments');
    if (!container) return;

    if (!comments || comments.length === 0) {
        container.innerHTML = '<div class="empty-message">No comments made yet</div>';
        return;
    }

    let commentsHTML = '';
    comments.forEach(comment => {
        commentsHTML += createActivityCommentHTML(comment);
    });

    container.innerHTML = commentsHTML;
}

// Create HTML for activity comment
function createActivityCommentHTML(comment) {
    const timeAgo = formatTimeAgo(comment.comment_date);

    return `
        <div class="activity-comment-item">
            <div class="activity-comment-header">
                <div class="activity-comment-post-info">
                    <span class="activity-comment-post-title">On: "${escapeHtml(comment.post_title)}"</span>
                    <span class="activity-comment-post-author">by ${escapeHtml(comment.post_author)}</span>
                </div>
                <span class="activity-comment-time">${timeAgo}</span>
            </div>
            <div class="activity-comment-content">
                ${escapeHtml(comment.comment)}
            </div>
            <div class="activity-comment-stats">
                <span><i class="material-icons">thumb_up</i> ${comment.likes}</span>
                <span><i class="material-icons">thumb_down</i> ${comment.dislikes}</span>
            </div>
        </div>
    `;
}

// Display activity error
function displayActivityError() {
    const containers = ['activity-created-posts', 'activity-liked-posts', 'activity-disliked-posts', 'activity-comments'];
    containers.forEach(containerId => {
        const container = document.getElementById(containerId);
        if (container) {
            container.innerHTML = '<div class="error-message">Failed to load activity data</div>';
        }
    });
}

// Format time ago
function formatTimeAgo(dateString) {
    const date = new Date(dateString);
    const now = new Date();
    const diffInSeconds = Math.floor((now - date) / 1000);

    if (diffInSeconds < 60) return 'Just now';
    if (diffInSeconds < 3600) return `${Math.floor(diffInSeconds / 60)} minutes ago`;
    if (diffInSeconds < 86400) return `${Math.floor(diffInSeconds / 3600)} hours ago`;
    if (diffInSeconds < 604800) return `${Math.floor(diffInSeconds / 86400)} days ago`;
    
    return date.toLocaleDateString();
}

// Escape HTML to prevent XSS
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// Auto-refresh notifications every 30 seconds
setInterval(() => {
    if (document.getElementById('Activity') && !document.getElementById('Activity').classList.contains('deactive')) {
        loadNotifications();
    } else {
        // Just update the badge count if not on activity page
        loadNotificationCount();
    }
}, 30000);

// Load notification count on page load
document.addEventListener('DOMContentLoaded', () => {
    // Check if user is authenticated and load notification count
    fetch("/auth/status", {
        method: "GET",
        credentials: "same-origin"
    })
    .then(response => response.json())
    .then(authData => {
        if (authData.authenticated) {
            loadNotificationCount();
        }
    })
    .catch(error => {
        console.error('Error checking auth status:', error);
    });
});

// Load notification count for badge
async function loadNotificationCount() {
    try {
        const response = await fetch('/Data-NotificationCount', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'X-Requested-With': 'XMLHttpRequest'
            }
        });

        if (response.ok) {
            const data = await response.json();
            updateNotificationBadge(data.count);
        }
    } catch (error) {
        console.error('Error loading notification count:', error);
    }
}
