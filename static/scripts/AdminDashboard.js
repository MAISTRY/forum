// Admin Dashboard JavaScript

// Load admin dashboard data when the page is shown
async function loadAdminDashboard() {
    try {
        await Promise.all([
            loadStatistics(),
            loadUsers(),
            loadCategories(),
            loadModerationRequests(),
            loadPostReports()
        ]);
    } catch (error) {
        console.error('Error loading admin dashboard:', error);
    }
}

// Load statistics
async function loadStatistics() {
    try {
        const response = await fetch('/Data-AdminStats', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'X-Requested-With': 'XMLHttpRequest'
            }
        });

        if (!response.ok) {
            throw new Error('Failed to load statistics');
        }

        const stats = await response.json();
        
        document.getElementById('admin-count').textContent = stats.AdminCount || 0;
        document.getElementById('mod-count').textContent = stats.ModeratorCount || 0;
        document.getElementById('post-count').textContent = stats.PostCount || 0;
        document.getElementById('comment-count').textContent = stats.CommentCount || 0;
    } catch (error) {
        console.error('Error loading statistics:', error);
        // Set default values on error
        document.getElementById('admin-count').textContent = '0';
        document.getElementById('mod-count').textContent = '0';
        document.getElementById('post-count').textContent = '0';
        document.getElementById('comment-count').textContent = '0';
    }
}

// Load users for management
async function loadUsers(searchTerm = '') {
    try {
        const url = searchTerm ? `/Data-AdminUsers?search=${encodeURIComponent(searchTerm)}` : '/Data-AdminUsers';
        const response = await fetch(url, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'X-Requested-With': 'XMLHttpRequest'
            }
        });

        if (!response.ok) {
            throw new Error('Failed to load users');
        }

        const users = await response.json();
        displayUsers(users);
    } catch (error) {
        console.error('Error loading users:', error);
        document.getElementById('users-list').innerHTML = '<div class="empty-message">Error loading users</div>';
    }
}

// Display users in the management section
function displayUsers(users) {
    const container = document.getElementById('users-list');

    // Handle case when users is null, undefined, or empty
    if (!users || !Array.isArray(users) || users.length === 0) {
        container.innerHTML = '<div class="empty-message">No users found</div>';
        return;
    }

    const usersHTML = users.map(user => {
        const roleClass = user.Privilege === 3 ? 'admin' : user.Privilege === 2 ? 'moderator' : 'user';
        const roleName = user.Privilege === 3 ? 'Admin' : user.Privilege === 2 ? 'Moderator' : 'User';

        let actionsHTML = '';
        if (user.Privilege === 1) {
            // Normal user - can promote to moderator
            actionsHTML = `<button class="btn-promote" onclick="promoteUser(${user.UserID}, 2)">Promote to Mod</button>`;
        } else if (user.Privilege === 2) {
            // Moderator - can only demote to user
            actionsHTML = `<button class="btn-demote" onclick="demoteUser(${user.UserID}, 1)">Demote to User</button>`;
        } else if (user.Privilege === 3) {
            // Admin - no actions allowed (admins cannot demote other admins)
            actionsHTML = '<span style="color: #666; font-size: 12px;">Admin</span>';
        }

        return `
            <div class="user-item">
                <div class="user-info">
                    <div class="username">@${user.Username}</div>
                    <div class="user-role ${roleClass}">${roleName}</div>
                </div>
                <div class="user-actions">
                    ${actionsHTML}
                </div>
            </div>
        `;
    }).join('');

    container.innerHTML = usersHTML;
}

// Search users
function searchUsers() {
    const searchTerm = document.getElementById('user-search').value.trim();
    loadUsers(searchTerm);
}

// Promote user
async function promoteUser(userId, newPrivilege) {
    try {
        const response = await fetch('/Data-AdminPromoteUser', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'X-Requested-With': 'XMLHttpRequest'
            },
            body: JSON.stringify({
                userId: userId,
                privilege: newPrivilege
            })
        });

        if (!response.ok) {
            throw new Error('Failed to promote user');
        }

        const result = await response.json();
        if (result.success) {
            // Reload users and statistics
            await loadUsers();
            await loadStatistics();
        } else {
            alert('Failed to promote user: ' + (result.message || 'Unknown error'));
        }
    } catch (error) {
        console.error('Error promoting user:', error);
        alert('Error promoting user. Please try again.');
    }
}

// Demote user
async function demoteUser(userId, newPrivilege) {
    try {
        const response = await fetch('/Data-AdminDemoteUser', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'X-Requested-With': 'XMLHttpRequest'
            },
            body: JSON.stringify({
                userId: userId,
                privilege: newPrivilege
            })
        });

        if (!response.ok) {
            throw new Error('Failed to demote user');
        }

        const result = await response.json();
        if (result.success) {
            // Reload users and statistics
            await loadUsers();
            await loadStatistics();
        } else {
            alert('Failed to demote user: ' + (result.message || 'Unknown error'));
        }
    } catch (error) {
        console.error('Error demoting user:', error);
        alert('Error demoting user. Please try again.');
    }
}

// Load moderation requests
async function loadModerationRequests() {
    try {
        const response = await fetch('/Data-AdminModerationRequests', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'X-Requested-With': 'XMLHttpRequest'
            }
        });

        if (!response.ok) {
            throw new Error('Failed to load moderation requests');
        }

        const requests = await response.json();
        displayModerationRequests(requests);
    } catch (error) {
        console.error('Error loading moderation requests:', error);
        document.getElementById('moderation-requests').innerHTML = '<div class="empty-message">Error loading requests</div>';
    }
}

// Display moderation requests
function displayModerationRequests(requests) {
    const container = document.getElementById('moderation-requests');

    // Handle case when requests is null, undefined, or empty
    if (!requests || !Array.isArray(requests) || requests.length === 0) {
        container.innerHTML = '<div class="empty-message">No moderation requests</div>';
        return;
    }

    const requestsHTML = requests.map(request => {
        let actionsHTML = '';
        if (request.Status === 'pending') {
            actionsHTML = `
                <div class="request-actions">
                    <button class="btn-approve" onclick="respondToRequest(${request.RequestID}, 'approved')">Approve</button>
                    <button class="btn-reject" onclick="respondToRequest(${request.RequestID}, 'rejected')">Reject</button>
                </div>
            `;
        }

        return `
            <div class="request-item">
                <div class="request-header">
                    <div class="request-user">@${request.Username}</div>
                    <div class="request-date">${formatDate(request.RequestDate)}</div>
                </div>
                <div class="request-status ${request.Status}">${request.Status}</div>
                ${actionsHTML}
            </div>
        `;
    }).join('');

    container.innerHTML = requestsHTML;
}

// Respond to moderation request
async function respondToRequest(requestId, status) {
    try {
        const response = await fetch('/Data-AdminRespondRequest', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'X-Requested-With': 'XMLHttpRequest'
            },
            body: JSON.stringify({
                requestId: requestId,
                status: status
            })
        });

        if (!response.ok) {
            throw new Error('Failed to respond to request');
        }

        const result = await response.json();
        if (result.success) {
            // Reload requests, users, and statistics
            await loadModerationRequests();
            await loadUsers();
            await loadStatistics();
        } else {
            alert('Failed to respond to request: ' + (result.message || 'Unknown error'));
        }
    } catch (error) {
        console.error('Error responding to request:', error);
        alert('Error responding to request. Please try again.');
    }
}

// Add event listener for Enter key in search input
document.addEventListener('DOMContentLoaded', function() {
    const searchInput = document.getElementById('user-search');
    if (searchInput) {
        searchInput.addEventListener('keypress', function(e) {
            if (e.key === 'Enter') {
                searchUsers();
            }
        });
    }
});

// Load categories for management
async function loadCategories() {
    try {
        const response = await fetch('/Data-AdminCategories', {
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
        displayCategories(categories);
    } catch (error) {
        console.error('Error loading categories:', error);
        document.getElementById('categories-list').innerHTML = '<div class="empty-message">Error loading categories</div>';
    }
}

// Display categories in the management section
function displayCategories(categories) {
    const container = document.getElementById('categories-list');

    // Handle case when categories is null, undefined, or empty
    if (!categories || !Array.isArray(categories) || categories.length === 0) {
        container.innerHTML = '<div class="empty-message">No categories found</div>';
        return;
    }

    const categoriesHTML = categories.map(category => {
        return `
            <div class="category-item">
                <div class="category-info">
                    <div class="category-title">${category.title}</div>
                    <div class="category-description">${category.description}</div>
                </div>
                <div class="category-actions">
                    <button class="btn-delete-category" onclick="deleteCategory(${category.CategoryID})">Delete</button>
                </div>
            </div>
        `;
    }).join('');

    container.innerHTML = categoriesHTML;
}

// Add new category
async function addCategory() {
    const title = document.getElementById('new-category-title').value.trim();
    const description = document.getElementById('new-category-description').value.trim();

    if (!title) {
        alert('Category title is required');
        return;
    }

    if (!description) {
        alert('Category description is required');
        return;
    }

    try {
        const response = await fetch('/Data-AdminAddCategory', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'X-Requested-With': 'XMLHttpRequest'
            },
            body: JSON.stringify({
                title: title,
                description: description
            })
        });

        if (!response.ok) {
            throw new Error('Failed to add category');
        }

        const result = await response.json();
        if (result.success) {
            // Clear form
            document.getElementById('new-category-title').value = '';
            document.getElementById('new-category-description').value = '';

            // Reload categories
            await loadCategories();
        } else {
            alert('Failed to add category: ' + (result.message || 'Unknown error'));
        }
    } catch (error) {
        console.error('Error adding category:', error);
        alert('Error adding category. Please try again.');
    }
}

// Delete category
async function deleteCategory(categoryId) {
    if (!confirm('Are you sure you want to delete this category? This action cannot be undone and will affect all posts in this category.')) {
        return;
    }

    try {
        const response = await fetch('/Data-AdminDeleteCategory', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'X-Requested-With': 'XMLHttpRequest'
            },
            body: JSON.stringify({
                categoryId: categoryId
            })
        });

        if (!response.ok) {
            throw new Error('Failed to delete category');
        }

        const result = await response.json();
        if (result.success) {
            // Reload categories
            await loadCategories();
        } else {
            alert('Failed to delete category: ' + (result.message || 'Unknown error'));
        }
    } catch (error) {
        console.error('Error deleting category:', error);
        alert('Error deleting category. Please try again.');
    }
}

// Make functions globally available
window.loadAdminDashboard = loadAdminDashboard;
window.searchUsers = searchUsers;
window.promoteUser = promoteUser;
window.demoteUser = demoteUser;
window.respondToRequest = respondToRequest;
window.addCategory = addCategory;
window.deleteCategory = deleteCategory;
window.respondToReport = respondToReport;

// Load post reports
async function loadPostReports() {
    try {
        const response = await fetch('/Data-AdminReports', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'X-Requested-With': 'XMLHttpRequest'
            }
        });

        if (!response.ok) {
            throw new Error('Failed to load post reports');
        }

        const reports = await response.json();
        displayPostReports(reports);
    } catch (error) {
        console.error('Error loading post reports:', error);
        document.getElementById('post-reports').innerHTML = '<div class="empty-message">Error loading reports</div>';
    }
}

// Display post reports
function displayPostReports(reports) {
    const container = document.getElementById('post-reports');

    if (!reports || !Array.isArray(reports) || reports.length === 0) {
        container.innerHTML = '<div class="empty-message">No post reports found</div>';
        return;
    }

    const reportsHTML = reports.map(report => {
        let actionsHTML = '';
        if (report.Status === 'pending') {
            actionsHTML = `
                <div class="report-actions">
                    <textarea id="response-${report.ReportID}" placeholder="Admin response (optional)..." rows="2"></textarea>
                    <div class="action-buttons">
                        <button class="btn-approve" onclick="respondToReport(${report.ReportID}, 'approved')">Approve & Delete Post</button>
                        <button class="btn-reject" onclick="respondToReport(${report.ReportID}, 'rejected')">Reject</button>
                    </div>
                </div>
            `;
        }

        // Truncate content if too long
        const truncatedContent = report.PostContent.length > 100
            ? report.PostContent.substring(0, 100) + '...'
            : report.PostContent;

        return `
            <div class="report-item">
                <div class="report-header">
                    <div class="report-info">
                        <div class="report-post-title"><strong>${report.PostTitle}</strong></div>
                        <div class="report-post-author">by @${report.PostAuthor}</div>
                    </div>
                    <div class="report-date">${formatDate(report.ReportDate)}</div>
                </div>
                <div class="report-content">${truncatedContent}</div>
                <div class="report-details">
                    <div class="report-moderator">Reported by: @${report.ModeratorName}</div>
                    <div class="report-reason"><strong>Reason:</strong> ${report.Reason}</div>
                </div>
                <div class="report-status ${report.Status}">${report.Status}</div>
                ${report.AdminResponse ? `<div class="admin-response"><strong>Admin Response:</strong> ${report.AdminResponse}</div>` : ''}
                ${actionsHTML}
            </div>
        `;
    }).join('');

    container.innerHTML = reportsHTML;
}

// Respond to post report
async function respondToReport(reportId, status) {
    const responseTextarea = document.getElementById(`response-${reportId}`);
    const response = responseTextarea ? responseTextarea.value.trim() : '';

    let confirmMessage = status === 'approved'
        ? 'Are you sure you want to approve this report? This will DELETE the post permanently.'
        : 'Are you sure you want to reject this report?';

    if (!confirm(confirmMessage)) {
        return;
    }

    try {
        const apiResponse = await fetch('/Data-AdminRespondReport', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'X-Requested-With': 'XMLHttpRequest'
            },
            body: JSON.stringify({
                reportId: reportId,
                status: status,
                response: response
            })
        });

        if (!apiResponse.ok) {
            throw new Error('Failed to respond to report');
        }

        const result = await apiResponse.json();
        if (result.success) {
            // Reload reports and statistics
            await loadPostReports();
            await loadStatistics();

            if (status === 'approved') {
                alert('Report approved and post deleted successfully.');
            } else {
                alert('Report rejected successfully.');
            }
        } else {
            alert('Failed to respond to report: ' + (result.message || 'Unknown error'));
        }
    } catch (error) {
        console.error('Error responding to report:', error);
        alert('Error responding to report. Please try again.');
    }
}
