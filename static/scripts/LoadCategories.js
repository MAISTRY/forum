// Load Categories Dynamically for Create Post Form

// Load categories when the create post page is shown
async function loadCategoriesForForm() {
    try {
        const response = await fetch('/Data-PublicCategories', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'X-Requested-With': 'XMLHttpRequest'
            }
        });

        if (!response.ok) {
            throw new Error('Failed to load categories');
        }

        const categories = await response.json();
        console.log('Loaded categories:', categories);
        displayCategoriesInForm(categories);
    } catch (error) {
        console.error('Error loading categories for form:', error);
        // Fallback to default categories if loading fails
        displayDefaultCategories();
    }
}

// Display categories in the create post form
function displayCategoriesInForm(categoriesData) {
    const categoriesContainer = document.getElementById('categories');
    if (!categoriesContainer) {
        console.error('Categories container not found');
        return;
    }

    console.log('Displaying categories:', categoriesData);

    // Handle case when categoriesData is null, undefined, or empty
    if (!categoriesData || !Array.isArray(categoriesData) || categoriesData.length === 0) {
        console.log('No categories data, using defaults');
        displayDefaultCategories();
        return;
    }

    // Generate HTML for categories from admin categories data
    const categoriesHTML = categoriesData.map(category => `
        <label class="multi-option">
            <input type="checkbox" name="categories" value="${category.title}">
            <span>${category.title}</span>
        </label>
    `).join('');

    console.log('Generated categories HTML:', categoriesHTML);
    categoriesContainer.innerHTML = categoriesHTML;
}

// Display default categories as fallback
function displayDefaultCategories() {
    const categoriesContainer = document.getElementById('categories');
    if (!categoriesContainer) return;

    const defaultCategories = [
        'Technology', 'Education', 'Entertainment', 'Travel', 
        'Cars', 'Sports', 'Lifestyle', 'Science', 'Business'
    ];

    const categoriesHTML = defaultCategories.map(categoryName => `
        <label class="multi-option">
            <input type="checkbox" name="categories" value="${categoryName}">
            <span>${categoryName}</span>
        </label>
    `).join('');

    categoriesContainer.innerHTML = categoriesHTML;
}

// Load categories when DOM is ready
document.addEventListener('DOMContentLoaded', function() {
    // Load categories initially
    loadCategoriesForForm();
});

// Also load categories when navigating to create post page
// This will be called by the navigation system
function onCreatePostPageShow() {
    loadCategoriesForForm();
}

// Make functions globally available
window.loadCategoriesForForm = loadCategoriesForForm;
window.onCreatePostPageShow = onCreatePostPageShow;
