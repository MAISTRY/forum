async function commentShow(distination,postId) {
    const commentsContainer = document.getElementById(`${distination}-comments-${postId}`);
    if (!commentsContainer) {
        console.error(`Comments container not found for post ID: ${postId}`);
        return;
    }

    if (commentsContainer.style.display === 'none') {

        document.querySelectorAll('.comments-section').forEach(element => {
            element.style.display = 'none';
            commentsContainer.innerHTML = '';
        });

        await loadComments(distination,postId);

    } else {
        commentsContainer.style.display = 'none';
        commentsContainer.innerHTML = '';
    }
}

async function loadComments(destination, postId) {
    const commentsContainer = document.getElementById(`${destination}-comments-${postId}`);

    try {
        // First get user privilege level
        const authResponse = await fetch("/auth/status", {
            method: "GET",
            credentials: "same-origin"
        });
        const authData = await authResponse.json();
        const userPrivilege = authData.privilege || 0;
        const currentUserId = authData.user_id || 0;

        // Then fetch comments
        const response = await fetch(`/Data-Comment?postid=${postId}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'X-Requested-With': 'XMLHttpRequest'
            }
        });

        if (!response.ok) throw new Error('Network response was not ok');

        const data = await response.json();

        const commentsSection = document.createElement('div');
        commentsSection.classList.add('comments-section');

        // Handle case when data is null, undefined, or empty
        if (!data || !Array.isArray(data) || data.length === 0) {
            const noCommentsMessage = document.createElement('div');
            noCommentsMessage.className = 'empty-message';
            noCommentsMessage.style.cssText = 'text-align: center; padding: 20px; color: #666; font-style: italic;';
            noCommentsMessage.textContent = 'No comments yet';
            commentsSection.appendChild(noCommentsMessage);
        } else {
            data.forEach(comment => {
                // Create Comment Card
                const commentCard = document.createElement('div');
                commentCard.classList.add('comment-card');

                // Create Comment Header
                const commentHeader = document.createElement('div');
                commentHeader.classList.add('comment-header');

                const usernameSpan = document.createElement('span');
                usernameSpan.classList.add('comment-username');
                usernameSpan.textContent = `@${comment.CmtUsername}`;

                const dateSpan = document.createElement('span');
                dateSpan.classList.add('comment-date');
                dateSpan.textContent = formatDate(comment.CmtDate);

                commentHeader.appendChild(usernameSpan);
                commentHeader.appendChild(dateSpan);

                // Create Comment Content
                const commentContent = document.createElement('div');
                commentContent.classList.add('comment-content');
                commentContent.textContent = comment.CmtContent;

                // Create Comment Footer
                const commentFooter = document.createElement('div');
                commentFooter.classList.add('comment-footer');
                commentFooter.id = `comment-${comment.CmtID}`;

                // Like Form
                const likeForm = document.createElement('form');
                likeForm.classList.add('like-cmt');
                likeForm.onclick = () => {
                    likeForm.addEventListener('submit', handleCommentInteraction);
                }

                const likeInput = document.createElement('input');
                likeInput.type = 'hidden';
                likeInput.name = 'commentId';
                likeInput.value = comment.CmtID;

                const likeButton = document.createElement('button');
                likeButton.classList.add('footer-buttons', 'comment-button', 'like-buttons');
                likeButton.title = 'Like';

                const likeIcon = document.createElement('i');
                likeIcon.classList.add('material-icons');
                likeIcon.textContent = 'thumb_up';

                const likeCount = document.createElement('span');
                likeCount.classList.add('cmtLikes');
                likeCount.textContent = comment.CmtLikes;

                likeButton.appendChild(likeIcon);
                likeButton.appendChild(likeCount);
                likeForm.appendChild(likeInput);
                likeForm.appendChild(likeButton);

                // Dislike Form
                const dislikeForm = document.createElement('form');
                dislikeForm.classList.add('dislike-cmt');
                dislikeForm.onclick = () => {
                    dislikeForm.addEventListener('submit', handleCommentInteraction);
                }
                const dislikeInput = document.createElement('input');
                dislikeInput.type = 'hidden';
                dislikeInput.name = 'commentId';
                dislikeInput.value = comment.CmtID;

                const dislikeButton = document.createElement('button');
                dislikeButton.classList.add('footer-buttons', 'comment-button', 'dislike-buttons');
                dislikeButton.title = 'Dislike';

                const dislikeIcon = document.createElement('i');
                dislikeIcon.classList.add('material-icons');
                dislikeIcon.textContent = 'thumb_down';

                const dislikeCount = document.createElement('span');
                dislikeCount.classList.add('cmtdisLikes');
                dislikeCount.textContent = comment.CmtDislikes;

                dislikeButton.appendChild(dislikeIcon);
                dislikeButton.appendChild(dislikeCount);
                dislikeForm.appendChild(dislikeInput);
                dislikeForm.appendChild(dislikeButton);

                commentFooter.appendChild(likeForm);
                commentFooter.appendChild(dislikeForm);

                // Edit/Delete buttons for comment owner
                if (currentUserId === comment.CmtUserID) {
                    // Edit button
                    const editButton = document.createElement('button');
                    editButton.classList.add('comment-edit-button', 'edit-btn');
                    editButton.title = 'Edit Comment';
                    editButton.onclick = () => editComment(comment.CmtID);

                    const editIcon = document.createElement('i');
                    editIcon.classList.add('material-icons');
                    editIcon.textContent = 'edit';
                    editButton.appendChild(editIcon);
                    commentFooter.appendChild(editButton);

                    // Delete button for owner
                    const userDeleteButton = document.createElement('button');
                    userDeleteButton.classList.add('comment-delete-button', 'delete-btn');
                    userDeleteButton.title = 'Delete Comment';
                    userDeleteButton.onclick = () => deleteComment(comment.CmtID);

                    const userDeleteIcon = document.createElement('i');
                    userDeleteIcon.classList.add('material-icons');
                    userDeleteIcon.textContent = 'delete';
                    userDeleteButton.appendChild(userDeleteIcon);
                    commentFooter.appendChild(userDeleteButton);
                }

                // Admin Delete Button (for admins only)
                if (userPrivilege === 3) { // Admin only
                    const adminDeleteButton = document.createElement('button');
                    adminDeleteButton.classList.add('comment-delete-button');
                    adminDeleteButton.title = 'Delete Comment (Admin)';
                    adminDeleteButton.onclick = () => adminDeleteComment(comment.CmtID);

                    const deleteIcon = document.createElement('i');
                    deleteIcon.classList.add('material-icons');
                    deleteIcon.textContent = 'delete';

                    adminDeleteButton.appendChild(deleteIcon);
                    commentFooter.appendChild(adminDeleteButton);
                }

                // Assemble Comment Card
                commentCard.appendChild(commentHeader);
                commentCard.appendChild(commentContent);
                commentCard.appendChild(commentFooter);

                commentsSection.appendChild(commentCard);
            });
        }

        // Create Comment Form
        const commentForm = document.createElement('form');
        commentForm.id = `form-${postId}`;
        commentForm.classList.add('comment-form-container');
        commentForm.onclick = () => {
            commentForm.addEventListener('submit', handleCommentSubmission);
        }

        const postIdInput = document.createElement('input');
        postIdInput.type = 'hidden';
        postIdInput.name = 'postId';
        postIdInput.value = postId;

        const commentTextarea = document.createElement('textarea');
        commentTextarea.name = 'comment';
        commentTextarea.classList.add('comment-input');
        commentTextarea.placeholder = 'Add your comment...';
        commentTextarea.required = true;

        const commentActions = document.createElement('div');
        commentActions.classList.add('comment-actions');

        const submitButton = document.createElement('button');
        submitButton.type = 'submit';
        submitButton.classList.add('comment-submit');
        submitButton.textContent = 'Post Comment';

        const loadingSpinner = document.createElement('div');
        loadingSpinner.classList.add('loading-spinner');
        submitButton.appendChild(loadingSpinner);

        const errorDiv = document.createElement('div');
        errorDiv.classList.add('comment-error');
        errorDiv.style.display = 'none';
        errorDiv.style.color = 'red';
        errorDiv.textContent = 'Failed to post comment. Please try again.';

        commentActions.appendChild(submitButton);
        commentActions.appendChild(errorDiv);

        commentForm.appendChild(postIdInput);
        commentForm.appendChild(commentTextarea);
        commentForm.appendChild(commentActions);

        commentsSection.appendChild(commentForm);

        commentsContainer.appendChild(commentsSection);
        commentsContainer.style.display = 'block';
        commentsContainer.scrollIntoView({
            behavior: 'smooth', 
            block: 'start',     
            inline: 'nearest'   
        });
        
    } catch (error) {
        console.error('Error loading comments:', error);
        commentsContainer.innerHTML = '<p>Error loading comments. Please try again later.</p>';
    }
}

async function handleCommentInteraction(event) {
    event.preventDefault(); 
    
    const form = event.currentTarget;
    const commentID = form.querySelector('input[name="commentId"]').value;
    const isLike = form.classList.contains('like-cmt');
    
    try {
        const endpoint = isLike ? '/Data-CommentLike' : '/Data-CommentDisLike';
        
        const response = await fetch(`${endpoint}?commentId=${commentID}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'X-Requested-With': 'XMLHttpRequest'
            }
        });
        
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        
        const data = await response.json();
        
        const buttonsContainer = document.getElementById(`comment-${commentID}`);
        if (buttonsContainer) {
            
            buttonsContainer.querySelector('.cmtLikes ').textContent = data.LikeCount;
            buttonsContainer.querySelector('.cmtdisLikes').textContent = data.DislikeCount;

        }
    } catch (error) {
        navigateToPage('Login');
    }
}

async function handleCommentSubmission(event) {
    event.preventDefault();

    const form = event.currentTarget;
    const parentDiv = form.parentElement;
    const postId = form.querySelector('input[name="postId"]').value;
    const comment = form.querySelector('textarea[name="comment"]').value;

    try {
        // Get user privilege level
        const authResponse = await fetch("/auth/status", {
            method: "GET",
            credentials: "same-origin"
        });
        const authData = await authResponse.json();
        const userPrivilege = authData.privilege || 0;
        const response = await fetch(`/Data-CreatComment`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'X-Requested-With': 'XMLHttpRequest'
            },
            body: JSON.stringify({
                postId: postId,
                comment: comment
            })

        });

        if (!response.ok) {
            throw new Error('Network response was not ok');
        }

        const data = await response.json();
        if (parentDiv) {
            // Create Comment Card
            const commentCard = document.createElement('div');
            commentCard.classList.add('comment-card');

            // Create Comment Header
            const commentHeader = document.createElement('div');
            commentHeader.classList.add('comment-header');

            const usernameSpan = document.createElement('span');
            usernameSpan.classList.add('comment-username');
            usernameSpan.textContent = `@${data.UserName}`;

            const dateSpan = document.createElement('span');
            dateSpan.classList.add('comment-date');
            dateSpan.textContent = data.CreateDate;

            commentHeader.appendChild(usernameSpan);
            commentHeader.appendChild(dateSpan);

            // Create Comment Content
            const commentContent = document.createElement('div');
            commentContent.classList.add('comment-content');
            commentContent.textContent = data.Comment;

            // Create Comment Footer
            const commentFooter = document.createElement('div');
            commentFooter.classList.add('comment-footer');
            commentFooter.id = `comment-${data.CommentID}`;

            // Like Form
            const likeForm = document.createElement('form');
            likeForm.classList.add('like-cmt');
            likeForm.onclick = () => {
                likeForm.addEventListener('submit', handleCommentInteraction);
            }
            
            const likeInput = document.createElement('input');
            likeInput.type = 'hidden';
            likeInput.name = 'commentId';
            likeInput.value = data.CommentID;

            const likeButton = document.createElement('button');
            likeButton.classList.add('footer-buttons', 'comment-button', 'like-buttons');
            likeButton.title = 'Like';

            const likeIcon = document.createElement('i');
            likeIcon.classList.add('material-icons');
            likeIcon.textContent = 'thumb_up';

            const likeCount = document.createElement('span');
            likeCount.classList.add('cmtLikes');
            likeCount.textContent = data.Likes;

            likeButton.appendChild(likeIcon);
            likeButton.appendChild(likeCount);
            likeForm.appendChild(likeInput);
            likeForm.appendChild(likeButton);

            // Dislike Form
            const dislikeForm = document.createElement('form');
            dislikeForm.classList.add('dislike-cmt');
            dislikeForm.onclick = () => {
                dislikeForm.addEventListener('submit', handleCommentInteraction);
            }
            const dislikeInput = document.createElement('input');
            dislikeInput.type = 'hidden';
            dislikeInput.name = 'commentId';
            dislikeInput.value = data.CommentID;

            const dislikeButton = document.createElement('button');
            dislikeButton.classList.add('footer-buttons', 'comment-button', 'dislike-buttons');
            dislikeButton.title = 'Dislike';

            const dislikeIcon = document.createElement('i');
            dislikeIcon.classList.add('material-icons');
            dislikeIcon.textContent = 'thumb_down';

            const dislikeCount = document.createElement('span');
            dislikeCount.classList.add('cmtdisLikes');
            dislikeCount.textContent = data.Dislikes;

            dislikeButton.appendChild(dislikeIcon);
            dislikeButton.appendChild(dislikeCount);
            dislikeForm.appendChild(dislikeInput);
            dislikeForm.appendChild(dislikeButton);

            commentFooter.appendChild(likeForm);
            commentFooter.appendChild(dislikeForm);

            // Edit/Delete buttons for comment owner (newly created comment is always owned by current user)
            const editButton = document.createElement('button');
            editButton.classList.add('comment-edit-button', 'edit-btn');
            editButton.title = 'Edit Comment';
            editButton.onclick = () => editComment(data.CommentID);

            const editIcon = document.createElement('i');
            editIcon.classList.add('material-icons');
            editIcon.textContent = 'edit';
            editButton.appendChild(editIcon);
            commentFooter.appendChild(editButton);

            const userDeleteButton = document.createElement('button');
            userDeleteButton.classList.add('comment-delete-button', 'delete-btn');
            userDeleteButton.title = 'Delete Comment';
            userDeleteButton.onclick = () => deleteComment(data.CommentID);

            const userDeleteIcon = document.createElement('i');
            userDeleteIcon.classList.add('material-icons');
            userDeleteIcon.textContent = 'delete';
            userDeleteButton.appendChild(userDeleteIcon);
            commentFooter.appendChild(userDeleteButton);

            // Admin Delete Button (for admins only)
            if (userPrivilege === 3) { // Admin only
                const adminDeleteButton = document.createElement('button');
                adminDeleteButton.classList.add('comment-delete-button');
                adminDeleteButton.title = 'Delete Comment (Admin)';
                adminDeleteButton.onclick = () => adminDeleteComment(data.CommentID);

                const deleteIcon = document.createElement('i');
                deleteIcon.classList.add('material-icons');
                deleteIcon.textContent = 'delete';

                adminDeleteButton.appendChild(deleteIcon);
                commentFooter.appendChild(adminDeleteButton);
            }

            // Assemble Comment Card
            commentCard.appendChild(commentHeader);
            commentCard.appendChild(commentContent);
            commentCard.appendChild(commentFooter);

            // Insert the new comment before submit form)
            parentDiv.insertBefore(commentCard, parentDiv.lastElementChild);

            // Clear the comment input field
            form.querySelector('textarea[name="comment"]').value = '';

        }
    } catch (error) {
        console.error('Error:', error);
        const errorElement = form.querySelector('.comment-error');
        if (errorElement) {
            errorElement.style.display = 'block';
            setTimeout(() => {
                errorElement.style.display = 'none';
            }, 3000);
        }
    }
}

// Delete comment function (for comment owners)
async function deleteComment(commentId) {
    if (!confirm('Are you sure you want to delete this comment? This action cannot be undone.')) {
        return;
    }

    try {
        const formData = new FormData();
        formData.append('commentId', commentId);

        const response = await fetch('/Data-UserDeleteComment', {
            method: 'POST',
            body: formData
        });

        const result = await response.json();

        if (result.success) {
            alert('Comment deleted successfully');
            location.reload();
        } else {
            alert(result.message || 'Failed to delete comment');
        }

    } catch (error) {
        console.error('Error deleting comment:', error);
        alert('Failed to delete comment');
    }
}

// Admin delete comment function (for admins only)
async function adminDeleteComment(commentId) {
    if (!confirm('Are you sure you want to delete this comment? This action cannot be undone.')) {
        return;
    }

    try {
        const response = await fetch('/Data-DeleteComment', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
                'X-Requested-With': 'XMLHttpRequest'
            },
            body: `commentId=${commentId}`
        });

        if (!response.ok) {
            throw new Error('Failed to delete comment');
        }

        // Find and remove the comment from the DOM
        const commentElement = document.getElementById(`comment-${commentId}`);
        if (commentElement) {
            const commentCard = commentElement.closest('.comment-card');
            if (commentCard) {
                commentCard.remove();
            }
        }

        // Show success message (optional)
        console.log('Comment deleted successfully');

    } catch (error) {
        console.error('Error deleting comment:', error);
        alert('Failed to delete comment. Please try again.');
    }
}
