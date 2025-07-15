// Post Form Handler - Prevents duplicate submissions and provides better UX

let isSubmitting = false;
let lastSubmissionData = null;

// Initialize form handler when DOM is loaded
document.addEventListener('DOMContentLoaded', function() {
    const postForm = document.getElementById('postForm');
    if (postForm) {
        setupPostFormHandler(postForm);
    }
});

function setupPostFormHandler(form) {
    const submitButton = form.querySelector('button[type="submit"]');
    const originalButtonText = submitButton.textContent;

    // Add event listener to prevent duplicate submissions
    form.addEventListener('submit', function(e) {
        // Check if already submitting
        if (isSubmitting) {
            e.preventDefault();
            return false;
        }

        // Get form data for duplicate check
        const formData = new FormData(form);
        const currentSubmissionData = {
            title: formData.get('title'),
            content: formData.get('content'),
            categories: formData.getAll('categories')
        };

        // Check if this is a duplicate of the last submission (within 30 seconds)
        if (lastSubmissionData && 
            JSON.stringify(currentSubmissionData) === JSON.stringify(lastSubmissionData) &&
            Date.now() - lastSubmissionData.timestamp < 30000) {
            
            e.preventDefault();
            alert('Please wait before submitting the same post again.');
            return false;
        }

        // Validate form before submission
        if (!validatePostForm(formData)) {
            e.preventDefault();
            return false;
        }

        // Set submitting state
        isSubmitting = true;
        submitButton.disabled = true;
        submitButton.textContent = 'Creating Post...';
        submitButton.style.opacity = '0.6';

        // Store submission data
        lastSubmissionData = {
            ...currentSubmissionData,
            timestamp: Date.now()
        };

        // Reset state after timeout (in case of network issues)
        setTimeout(() => {
            resetSubmissionState(submitButton, originalButtonText);
        }, 10000); // 10 seconds timeout
    });

    // Listen for HTMX events to reset state
    form.addEventListener('htmx:afterRequest', function(e) {
        resetSubmissionState(submitButton, originalButtonText);
        
        // If successful, clear form
        if (e.detail.xhr.status === 200) {
            clearPostForm(form);
        }
    });

    // Handle HTMX errors
    form.addEventListener('htmx:responseError', function(e) {
        resetSubmissionState(submitButton, originalButtonText);
    });
}

function validatePostForm(formData) {
    const title = formData.get('title');
    const content = formData.get('content');
    const categories = formData.getAll('categories');

    // Clear previous error messages
    const errorField = document.getElementById('CreatePost_err_field');
    if (errorField) {
        errorField.innerHTML = '';
    }

    // Validate title
    if (!title || title.trim().length === 0) {
        showError('Post title is required');
        return false;
    }

    if (title.trim().length > 100) {
        showError('Post title must be less than 100 characters');
        return false;
    }

    // Validate content
    if (!content || content.trim().length === 0) {
        showError('Post content is required');
        return false;
    }

    if (content.trim().length > 1000) {
        showError('Post content must be less than 1000 characters');
        return false;
    }

    // Validate categories
    if (categories.length === 0) {
        showError('Please select at least one category');
        return false;
    }

    return true;
}

function showError(message) {
    const errorField = document.getElementById('CreatePost_err_field');
    if (errorField) {
        errorField.innerHTML = `<div style="color: #dc3545; padding: 10px; background: #f8d7da; border: 1px solid #f5c6cb; border-radius: 4px; margin: 10px 0;">${message}</div>`;
    }
}

function resetSubmissionState(submitButton, originalButtonText) {
    isSubmitting = false;
    submitButton.disabled = false;
    submitButton.textContent = originalButtonText;
    submitButton.style.opacity = '1';
}

function clearPostForm(form) {
    // Clear text inputs
    const titleInput = form.querySelector('input[name="title"]');
    const contentTextarea = form.querySelector('textarea[name="content"]');
    const fileInput = form.querySelector('input[name="image"]');
    
    if (titleInput) titleInput.value = '';
    if (contentTextarea) contentTextarea.value = '';
    if (fileInput) fileInput.value = '';

    // Clear category checkboxes
    const categoryCheckboxes = form.querySelectorAll('input[name="categories"]');
    categoryCheckboxes.forEach(checkbox => {
        checkbox.checked = false;
    });

    // Clear error messages
    const errorField = document.getElementById('CreatePost_err_field');
    if (errorField) {
        errorField.innerHTML = '';
    }
}

// Character counter for content textarea
document.addEventListener('DOMContentLoaded', function() {
    const contentTextarea = document.getElementById('content');
    if (contentTextarea) {
        // Create character counter
        const counter = document.createElement('div');
        counter.id = 'content-counter';
        counter.style.cssText = 'text-align: right; font-size: 12px; color: #666; margin-top: 5px;';
        counter.textContent = '0/1000';
        
        contentTextarea.parentNode.insertBefore(counter, contentTextarea.nextSibling);

        // Update counter on input
        contentTextarea.addEventListener('input', function() {
            const length = this.value.length;
            counter.textContent = `${length}/1000`;
            
            if (length > 900) {
                counter.style.color = '#dc3545';
            } else if (length > 800) {
                counter.style.color = '#ffc107';
            } else {
                counter.style.color = '#666';
            }
        });
    }
});

// Make functions globally available
window.setupPostFormHandler = setupPostFormHandler;
