function loadProfileData() {
    const createdPostsContainer = document.getElementById('Created');
    const likedPostsContainer = document.getElementById('Liked');
    const dislikedPostsContainer = document.getElementById('Disliked');
    const commentsContainer = document.getElementById('Profile');

    // Validate containers
    if (!createdPostsContainer || !likedPostsContainer || !dislikedPostsContainer || !commentsContainer) {
        console.error("One or more containers are missing from the DOM.");
        return;
    }

    // Show loading messages
    [createdPostsContainer, likedPostsContainer, dislikedPostsContainer].forEach(container => {
        container.innerHTML = '<p style="text-align: center">Loading posts...</p>';
    });

    // Fetch profile data
    fetch('/Data-Profile', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'X-Requested-With': 'XMLHttpRequest'
        }
    })
    .then(response => {
        if (!response.ok) {
            throw new Error(`${response.status}: ${response.statusText || 'Unknown Error'}`);
        }
        return response.json();
    })
    .then(profileData => {
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
            if (!posts || posts.length === 0) {
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