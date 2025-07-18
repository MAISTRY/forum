// Content management functionality for edit/delete posts and comments

// Edit post functionality
async function editPost(postId) {
    try {
        // Get post data for editing
        const response = await fetch(`/Data-GetPostForEdit?postId=${postId}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'X-Requested-With': 'XMLHttpRequest'
            }
        });

        if (!response.ok) {
            throw new Error(`${response.status}: ${response.statusText}`);
        }

        const postData = await response.json();
        showEditPostModal(postData);
        
    } catch (error) {
        console.error('Error loading post for edit:', error);
        alert('Failed to load post for editing');
    }
}

// Show edit post modal
function showEditPostModal(postData) {
    const modal = document.createElement('div');
    modal.className = 'edit-modal-overlay';
    modal.innerHTML = `
        <div class="edit-modal">
            <div class="edit-modal-header">
                <h3>Edit Post</h3>
                <button class="close-modal" onclick="closeEditModal()">&times;</button>
            </div>
            <div class="edit-modal-body">
                <form id="edit-post-form">
                    <input type="hidden" id="edit-post-id" value="${postData.post_id}">
                    <div class="form-group">
                        <label for="edit-post-title">Title:</label>
                        <input type="text" id="edit-post-title" value="${escapeHtml(postData.title)}" required>
                    </div>
                    <div class="form-group">
                        <label for="edit-post-content">Content:</label>
                        <textarea id="edit-post-content" rows="6" required>${escapeHtml(postData.content)}</textarea>
                    </div>
                    <div class="form-actions">
                        <button type="button" onclick="closeEditModal()">Cancel</button>
                        <button type="submit">Save Changes</button>
                    </div>
                </form>
            </div>
        </div>
    `;

    document.body.appendChild(modal);

    // Handle form submission
    document.getElementById('edit-post-form').addEventListener('submit', async (e) => {
        e.preventDefault();
        await savePostEdit();
    });
}

// Save post edit
async function savePostEdit() {
    const postId = document.getElementById('edit-post-id').value;
    const title = document.getElementById('edit-post-title').value.trim();
    const content = document.getElementById('edit-post-content').value.trim();

    if (!title || !content) {
        alert('Title and content cannot be empty');
        return;
    }

    try {
        const response = await fetch('/Data-EditPost', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'X-Requested-With': 'XMLHttpRequest'
            },
            body: JSON.stringify({
                post_id: postId,
                title: title,
                content: content
            })
        });

        const result = await response.json();

        if (result.success) {
            alert('Post updated successfully');
            closeEditModal();
            // Refresh the current page to show updated content
            location.reload();
        } else {
            alert(result.message || 'Failed to update post');
        }
        
    } catch (error) {
        console.error('Error saving post edit:', error);
        alert('Failed to save changes');
    }
}

// Edit comment functionality
async function editComment(commentId) {
    try {
        // Get comment data for editing
        const response = await fetch(`/Data-GetCommentForEdit?commentId=${commentId}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'X-Requested-With': 'XMLHttpRequest'
            }
        });

        if (!response.ok) {
            throw new Error(`${response.status}: ${response.statusText}`);
        }

        const commentData = await response.json();
        showEditCommentModal(commentData);
        
    } catch (error) {
        console.error('Error loading comment for edit:', error);
        alert('Failed to load comment for editing');
    }
}

// Show edit comment modal
function showEditCommentModal(commentData) {
    const modal = document.createElement('div');
    modal.className = 'edit-modal-overlay';
    modal.innerHTML = `
        <div class="edit-modal">
            <div class="edit-modal-header">
                <h3>Edit Comment</h3>
                <button class="close-modal" onclick="closeEditModal()">&times;</button>
            </div>
            <div class="edit-modal-body">
                <form id="edit-comment-form">
                    <input type="hidden" id="edit-comment-id" value="${commentData.comment_id}">
                    <div class="form-group">
                        <label for="edit-comment-content">Comment:</label>
                        <textarea id="edit-comment-content" rows="4" required>${escapeHtml(commentData.content)}</textarea>
                    </div>
                    <div class="form-actions">
                        <button type="button" onclick="closeEditModal()">Cancel</button>
                        <button type="submit">Save Changes</button>
                    </div>
                </form>
            </div>
        </div>
    `;

    document.body.appendChild(modal);

    // Handle form submission
    document.getElementById('edit-comment-form').addEventListener('submit', async (e) => {
        e.preventDefault();
        await saveCommentEdit();
    });
}

// Save comment edit
async function saveCommentEdit() {
    const commentId = document.getElementById('edit-comment-id').value;
    const content = document.getElementById('edit-comment-content').value.trim();

    if (!content) {
        alert('Comment content cannot be empty');
        return;
    }

    try {
        const response = await fetch('/Data-EditComment', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'X-Requested-With': 'XMLHttpRequest'
            },
            body: JSON.stringify({
                comment_id: commentId,
                content: content
            })
        });

        const result = await response.json();

        if (result.success) {
            alert('Comment updated successfully');
            closeEditModal();
            // Refresh the current page to show updated content
            location.reload();
        } else {
            alert(result.message || 'Failed to update comment');
        }
        
    } catch (error) {
        console.error('Error saving comment edit:', error);
        alert('Failed to save changes');
    }
}

// Delete post functionality
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
            // Refresh the current page
            location.reload();
        } else {
            alert(result.message || 'Failed to delete post');
        }
        
    } catch (error) {
        console.error('Error deleting post:', error);
        alert('Failed to delete post');
    }
}

// Delete comment functionality
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
            // Refresh the current page
            location.reload();
        } else {
            alert(result.message || 'Failed to delete comment');
        }
        
    } catch (error) {
        console.error('Error deleting comment:', error);
        alert('Failed to delete comment');
    }
}

// Close edit modal
function closeEditModal() {
    const modal = document.querySelector('.edit-modal-overlay');
    if (modal) {
        modal.remove();
    }
}

// Escape HTML to prevent XSS
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}
