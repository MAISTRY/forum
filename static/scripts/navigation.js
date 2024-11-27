document.addEventListener('DOMContentLoaded', () => {

    const validPages = ['Home', 'Error', 'Categories'];

    fetch("/auth/status", {
        method: "GET",
        credentials: "same-origin"  // Important to include cookies in the request
    })
    .then(response => response.json())
    .then(data => {
        const loginButtons = document.querySelectorAll(".login-buttons");
        const logoutButtons = document.querySelectorAll(".logout-buttons");
        const privilage = document.querySelectorAll(".privilage");
        
            if (data.authenticated && data.privilege === 1) {
                loginButtons.forEach(section => section.classList.add('deactive'));
                logoutButtons.forEach(section => section.classList.remove('deactive'));

                privilage.forEach(section => section.classList.remove('disabled-link'));
                if (!validPages.includes('Createpost')) {
                    validPages.push('Createpost');
                }
                
            } else if (data.authenticated && data.privilege === 2) {
                loginButtons.forEach(section => section.classList.add('deactive'));
                logoutButtons.forEach(section => section.classList.remove('deactive'));

            } else if (data.authenticated && data.privilege === 3) {
                loginButtons.forEach(section => section.classList.add('deactive'));
                logoutButtons.forEach(section => section.classList.remove('deactive'));

            } else {
                loginButtons.forEach(section => section.classList.remove('deactive'));
                logoutButtons.forEach(section => section.classList.add('deactive'));
                privilage.forEach(section => section.classList.add('disabled-link'));
                if (!validPages.includes('Login','Register')) {
                    validPages.push('Login','Register');
                }
            }
            setupNavigationListeners();
    })
    .catch(error => {
        console.error("Error checking auth status:", error)
        navigateToPage('Error');
    });
    
    
    function setupNavigationListeners(){
        
        // Event listeners for navigation buttons
        document.querySelectorAll('.sidebar [data-page], .main-content [data-page]').forEach(element => {
            element.addEventListener('click', () => {
                const pageId = element.dataset.page;
                navigateToPage(pageId); 
            });
        });
    
        // Handle the browser's back/forward navigation
        window.addEventListener('popstate', () => {
            const currentPath = window.location.pathname.slice(1);
            let pageId = currentPath.charAt(0).toUpperCase() + currentPath.slice(1);

            if (!validPages.includes(pageId)) {
                pageId = 'Error'; 
            }
            showPage(pageId);
        });

        // On page load, check the URL for the current path
        const currentPath = window.location.pathname.slice(1);
        const pageIdToShow = currentPath ? currentPath.charAt(0).toUpperCase() + currentPath.slice(1) : 'Home';
        showPage(validPages.includes(pageIdToShow) ? pageIdToShow : 'Error');
    }

    function navigateToPage(pageId) {
        if (validPages.includes(pageId)) {
            const path = `/${pageId.toLowerCase()}`;
            history.pushState({}, '', path);
            showPage(pageId);
            localStorage.setItem('currentPage', pageId);
        } else {
            console.error(`Invalid page ID: ${pageId}`);
            navigateToPage('Error');
        }
    }
    window.navigateToPage = navigateToPage;
    // Show the requested page
    function showPage(page) {
        const contents = document.querySelectorAll('.deactive');
        contents.forEach(content => {
            content.classList.remove('active');
        });
        
        const activeContent = document.getElementById(page);
        if (activeContent) {
            activeContent.classList.add('active');
        }
        if (page === 'Home') {
            loadPosts();
            console.log("Posts loaded, applying handlers...");
        } else if (page === 'Categories') {
            loadCategories();
            console.log("Categories loaded, applying handlers...");
        }
    }
});
