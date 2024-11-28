
function loadCategories() {
    const postsContainer = document.getElementById('category-container');
    const categoryNav = document.getElementById('category-nav');
    postsContainer.innerHTML = '<p style="text-align: center">Loading posts...</p>';
    sessionStorage.removeItem('reloaded');

    fetch('/Data-Categories', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        }
    })
    .then(response => {
        if (!response.ok) {

            const fallbackMessages = {
                400: 'Bad Request - Please check your input.',
                401: 'Unauthorized - Please log in.',
                403: 'Forbidden - You do not have permission to access this resource.',
                404: 'Not Found - The requested resource was not found.',
                405: 'Method Not Allowed - The action is not supported.',
                429: 'Too Many Requests - Please try again later.',
                500: 'Internal Server Error - Please try again later.',
                502: 'Bad Gateway - The server received an invalid response.',
                503: 'Service Unavailable - The server is temporarily unavailable.',
                504: 'Gateway Timeout - The server took too long to respond.'
            };
            const statusText = response.statusText || fallbackMessages[response.status] || 'Unknown Error';

            throw new Error(`${response.status}: ${statusText}`);
        }
        return response.json();
    })
    .then(categories => {
        const fragment = document.createDocumentFragment();
        categoryNav.innerHTML = '';

        const showAll = document.createElement('button');
        showAll.className = 'category-button';
        showAll.type = 'button';
        showAll.textContent = 'Show All';
        showAll.addEventListener('click', () => {
            const allsections = document.querySelectorAll('.category-section');
            allsections.forEach(section => section.classList.remove('deactive'));
        });
        categoryNav.appendChild(showAll);
        categories.forEach(category => {
            
            const categoryButton = document.createElement('button');
            categoryButton.className = 'category-button';
            categoryButton.type = 'button';
            categoryButton.textContent = category.CategoryName;
            categoryButton.addEventListener('click', () => {

                const allsections = document.querySelectorAll('.category-section');
                allsections.forEach(section => section.classList.add('deactive'));

                const currentSection = document.getElementById(`category-${category.CategoryName}`);
                if (currentSection) currentSection.classList.remove('deactive');

            });
            categoryNav.appendChild(categoryButton);

            const categoryElement = document.createElement('div');
            categoryElement.className = 'category-section deactive';
            categoryElement.id = `category-${category.CategoryName}`;
            categoryElement.innerHTML = `<h2 class="category-title">${category.CategoryName}</h2>`;

            const postsList = document.createElement('div');
            postsList.className = 'posts-list';

            category.Posts.forEach(post => {

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
                buttonsContainer.id = `Category-post-${post.PostID}`;
                buttonsContainer.classList.add('buttons-contant');

                // Like Button
                const likeForm = document.createElement('form');
                likeForm.classList.add('like-form');
                likeForm.onclick = (event) => {
                    likeForm.addEventListener('submit', handlePostInteraction(event,"Category"));
                }

                const likeInput = document.createElement('input');
                likeInput.type = 'hidden';
                likeInput.name = 'postId';
                likeInput.value = post.PostID;

                const likeButton = document.createElement('button');
                likeButton.type = 'submit';
                likeButton.classList.add('footer-buttons', 'post-button');
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
                    dislikeForm.addEventListener('submit', handlePostInteraction(event,"Category"));
                }

                const dislikeInput = document.createElement('input');
                dislikeInput.type = 'hidden';
                dislikeInput.name = 'postId';
                dislikeInput.value = post.PostID;

                const dislikeButton = document.createElement('button');
                dislikeButton.type = 'submit';
                dislikeButton.classList.add('footer-buttons', 'post-button');
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
                commentButton.classList.add('footer-buttons', 'post-button');
                commentButton.title = 'comment';
                commentButton.onclick = () => {
                    commentButton.addEventListener('click', commentShow("Category",post.PostID));
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
                commentsContainer.id = `Category-comments-${post.PostID}`;
                postElement.appendChild(commentsContainer);
                postsList.appendChild(postElement);

            });
            
            categoryElement.appendChild(postsList);
            fragment.appendChild(categoryElement);
        });

        postsContainer.innerHTML = '';
        postsContainer.appendChild(fragment);
        
        if (categories.length > 0) {
            const allsections = document.querySelectorAll('.category-section');
            allsections.forEach(section => section.classList.remove('deactive'));
        }

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
}
