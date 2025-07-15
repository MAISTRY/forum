function loadPosts() {
    const postsContainers = document.querySelectorAll('#posts-container');

    postsContainers.forEach(container => {
        container.innerHTML = '<p style="text-align: center">Loading posts...</p>';

        // First get user privilege level
        fetch("/auth/status", {
            method: "GET",
            credentials: "same-origin"
        })
        .then(response => response.json())
        .then(authData => {
            // Then fetch posts
            return fetch('/Data-Post', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'X-Requested-With': 'XMLHttpRequest'
                }
            }).then(response => {
                if (!response.ok) {
                    const fallbackMessages = {
                        400: 'Bad Request - Please check your input.',
                        401: 'Unauthorized - Please log in.',
                        403: 'Forbidden - You do not have permission to access this resource.',
                        404: 'Not Found - The requested resource was not found.',
                        405: 'Method Not Allowed - The action is not supported.',
                        500: 'Internal Server Error - Please try again later.',
                        502: 'Bad Gateway - The server received an invalid response.',
                        503: 'Service Unavailable - The server is temporarily unavailable.',
                        504: 'Gateway Timeout - The server took too long to respond.'
                    };
                    const statusText = response.statusText || fallbackMessages[response.status] || 'Unknown Error';

                    throw new Error(`${response.status}: ${statusText}`);
                }
                return response.json();
            }).then(posts => {
                return { posts, userPrivilege: authData.privilege || 0, currentUserId: authData.user_id || 0 };
            });
        })
        .then(data => {
            const { posts, userPrivilege, currentUserId } = data;
            const fragment = document.createDocumentFragment();

            // Handle case when posts is null, undefined, or empty
            if (!posts || !Array.isArray(posts) || posts.length === 0) {
                container.innerHTML = '<div class="empty-message" style="text-align: center; padding: 40px; color: #666; font-style: italic;">No posts available</div>';
                return;
            }

            posts.forEach(post => {

                const postElement = document.createElement('div');
                postElement.id = 'post-' + post.PostID;
                postElement.className = 'post-card';
                
                const postCategoryContainer = document.createElement('div');
                postCategoryContainer.className = 'post-category';

                const categorySpan = document.createElement('span');
                categorySpan.className = 'text-category';
                categorySpan.textContent = 'Categories:';
                postCategoryContainer.appendChild(categorySpan);
            
                post.Categories.forEach(CTG => {
                    const categoryElement = document.createElement('a');
                    categoryElement.textContent = CTG;
                    postCategoryContainer.appendChild(categoryElement);
                });
                
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
                buttonsContainer.id = `Home-post-${post.PostID}`;
                buttonsContainer.classList.add('buttons-contant');

                // Like Button
                const likeForm = document.createElement('form');
                likeForm.classList.add('like-form');
                likeForm.onclick = (event) => {
                    likeForm.addEventListener('submit', handlePostInteraction(event,"Home"));
                }

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
                    dislikeForm.addEventListener('submit', handlePostInteraction(event,"Home"));
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
                    commentButton.addEventListener('click', commentShow("Home",post.PostID));
                }

                const commentIcon = document.createElement('i');
                commentIcon.classList.add('material-icons');
                commentIcon.textContent = 'comment';

                const commentCount = document.createElement('span');
                commentCount.id = `comment-count-${post.PostID}`;
                commentCount.classList.add('comments');
                commentCount.textContent = post.CmtCount;

                commentButton.appendChild(commentIcon);
                commentButton.appendChild(commentCount);

                // Edit/Delete buttons for post owner
                let editButton = null;
                let userDeleteButton = null;
                if (currentUserId === post.UserID) {
                    // Edit button
                    editButton = document.createElement('button');
                    editButton.classList.add('footer-buttons', 'post-button', 'edit-btn');
                    editButton.title = 'Edit Post';
                    editButton.onclick = () => editPost(post.PostID);

                    const editIcon = document.createElement('i');
                    editIcon.classList.add('material-icons');
                    editIcon.textContent = 'edit';
                    editButton.appendChild(editIcon);

                    // Delete button for owner
                    userDeleteButton = document.createElement('button');
                    userDeleteButton.classList.add('footer-buttons', 'post-button', 'delete-btn');
                    userDeleteButton.title = 'Delete Post';
                    userDeleteButton.onclick = () => deletePost(post.PostID);

                    const userDeleteIcon = document.createElement('i');
                    userDeleteIcon.classList.add('material-icons');
                    userDeleteIcon.textContent = 'delete';
                    userDeleteButton.appendChild(userDeleteIcon);
                }

                // Admin/Moderator Delete Button (for admins and moderators only)
                let adminDeleteButton = null;
                if (userPrivilege >= 2) { // Moderators (2) and Admins (3)
                    adminDeleteButton = document.createElement('button');
                    adminDeleteButton.classList.add('footer-buttons', 'post-button', 'delete-button');
                    adminDeleteButton.title = 'Delete Post (Admin)';
                    adminDeleteButton.onclick = () => adminDeletePost(post.PostID);

                    const deleteIcon = document.createElement('i');
                    deleteIcon.classList.add('material-icons');
                    deleteIcon.textContent = 'delete';

                    adminDeleteButton.appendChild(deleteIcon);
                }

                // Report Button (for moderators only)
                let reportButton = null;
                if (userPrivilege === 2) { // Moderators only (not admins)
                    reportButton = document.createElement('button');
                    reportButton.classList.add('footer-buttons', 'post-button', 'report-button');
                    reportButton.title = 'Report Post';
                    reportButton.onclick = () => reportPost(post.PostID);

                    const reportIcon = document.createElement('i');
                    reportIcon.classList.add('material-icons');
                    reportIcon.textContent = 'report';

                    reportButton.appendChild(reportIcon);
                }

                // User Info
                const postUser = document.createElement('div');
                postUser.classList.add('footer-buttons', 'post-user');
                postUser.textContent = `@${post.username}`;

                buttonsContainer.appendChild(likeForm);
                buttonsContainer.appendChild(dislikeForm);
                buttonsContainer.appendChild(commentButton);
                if (editButton) {
                    buttonsContainer.appendChild(editButton);
                }
                if (userDeleteButton) {
                    buttonsContainer.appendChild(userDeleteButton);
                }
                if (adminDeleteButton) {
                    buttonsContainer.appendChild(adminDeleteButton);
                }
                if (reportButton) {
                    buttonsContainer.appendChild(reportButton);
                }

                postFooter.appendChild(buttonsContainer);
                postFooter.appendChild(postUser);

                // Append all components to the main post element
                postElement.appendChild(postHeader);
                postElement.appendChild(postContent);
                // Post Image (if exists)
                if (post.imagePath) {
                    const postImage = document.createElement('img');
                    postImage.src = post.imagePath;
                    postImage.alt = "Post Image";
                    postImage.classList.add('image-size');

                    const postImageContainer = document.createElement('div');
                    postImageContainer.classList.add('post-image');
                    postImageContainer.appendChild(postImage);

                    postElement.appendChild(postImageContainer);
                }
                postElement.appendChild(postFooter);

                const commentsContainer = document.createElement('div'); 
                commentsContainer.className = 'comments-section';
                commentsContainer.style.display = 'none';
                commentsContainer.id = `Home-comments-${post.PostID}`;
                postElement.appendChild(commentsContainer);
                fragment.appendChild(postElement);

            });
            
            container.innerHTML = '';
            container.appendChild(fragment);
        })
        .catch(error => {
            console.error(error);
            navigateToPage('Error');
            const errorCode = document.getElementById('error-id');
            const errorMessage = document.getElementById('error-message');
    
            const [status, message] = error.message.split(':'); 
            errorCode.innerHTML = status.trim();
            errorMessage.innerHTML = message.trim() || 'an unexpected error occurred.';
        });
    });
}

// Delete post function (for post owners)
async function deletePost(postId) {
    if (!confirm('Are you sure you want to delete this post? This action cannot be undone.')) {
        return;
    }

    try {
        const formData = new FormData();
        formData.append('postId', postId);

        const response = await fetch('/Data-UserDeletePost', {
            method: 'POST',
            body: formData
        });

        const result = await response.json();

        if (result.success) {
            alert('Post deleted successfully');
            location.reload();
        } else {
            alert(result.message || 'Failed to delete post');
        }

    } catch (error) {
        console.error('Error deleting post:', error);
        alert('Failed to delete post');
    }
}

// Admin delete post function (for admins and moderators)
async function adminDeletePost(postId) {
    if (!confirm('Are you sure you want to delete this post? This action cannot be undone.')) {
        return;
    }

    try {
        const response = await fetch('/Data-DeletePost', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
                'X-Requested-With': 'XMLHttpRequest'
            },
            body: `postId=${postId}`
        });

        if (!response.ok) {
            throw new Error('Failed to delete post');
        }

        // Remove the post from the DOM
        const postElement = document.getElementById(`post-${postId}`);
        if (postElement) {
            postElement.remove();
        }

        // Show success message (optional)
        console.log('Post deleted successfully');

    } catch (error) {
        console.error('Error deleting post:', error);
        alert('Failed to delete post. Please try again.');
    }
}

// Report post function (for moderators only)
async function reportPost(postId) {
    const reason = prompt('Please provide a reason for reporting this post:');
    if (!reason || reason.trim() === '') {
        return;
    }

    try {
        const response = await fetch('/Data-ReportPost', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
                'X-Requested-With': 'XMLHttpRequest'
            },
            body: `postId=${postId}&reason=${encodeURIComponent(reason.trim())}`
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(errorText || 'Failed to report post');
        }

        const result = await response.json();
        if (result.success) {
            alert('Post reported successfully. An admin will review it.');
        } else {
            alert('Failed to report post: ' + (result.message || 'Unknown error'));
        }

    } catch (error) {
        console.error('Error reporting post:', error);
        if (error.message.includes('already reported')) {
            alert('You have already reported this post.');
        } else {
            alert('Failed to report post. Please try again.');
        }
    }
}
