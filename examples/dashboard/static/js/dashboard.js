/**
 * Dashboard JavaScript
 * Handles all interactive functionality for the dashboard
 */

class Dashboard {
    constructor() {
        this.init();
    }

    init() {
        this.bindEvents();
        this.initializeComponents();
    }

    bindEvents() {
        // Sidebar toggle
        document.addEventListener('click', (e) => {
            if (e.target.matches('[onclick*="toggleSidebar"]')) {
                this.toggleSidebar();
            }
            if (e.target.matches('[onclick*="toggleSidebarCollapse"]')) {
                this.toggleSidebarCollapse();
            }
        });

        // Navigation
        document.addEventListener('click', (e) => {
            if (e.target.matches('[onclick*="toggleNavGroup"]')) {
                const groupId = e.target.getAttribute('onclick').match(/toggleNavGroup\('([^']+)'\)/)[1];
                this.toggleNavGroup(groupId);
            }
        });

        // Dropdowns
        document.addEventListener('click', (e) => {
            if (e.target.matches('[onclick*="toggleNotifications"]')) {
                this.toggleDropdown('.notification-dropdown');
            }
            if (e.target.matches('[onclick*="toggleUserMenu"]')) {
                this.toggleDropdown('.user-dropdown');
            }
        });

        // Close dropdowns when clicking outside
        document.addEventListener('click', (e) => {
            if (!e.target.closest('.notification-dropdown') && !e.target.matches('[onclick*="toggleNotifications"]')) {
                this.closeDropdown('.notification-dropdown');
            }
            if (!e.target.closest('.user-dropdown') && !e.target.matches('[onclick*="toggleUserMenu"]')) {
                this.closeDropdown('.user-dropdown');
            }
        });

        // Modal handling
        document.addEventListener('click', (e) => {
            if (e.target.matches('[onclick*="closeModal"]')) {
                this.closeModal();
            }
        });

        // Table actions
        document.addEventListener('click', (e) => {
            if (e.target.matches('[onclick*="editItem"]')) {
                const id = e.target.getAttribute('onclick').match(/editItem\('([^']+)'\)/)[1];
                this.editItem(id);
            }
            if (e.target.matches('[onclick*="deleteItem"]')) {
                const id = e.target.getAttribute('onclick').match(/deleteItem\('([^']+)'\)/)[1];
                this.deleteItem(id);
            }
        });

        // Pagination
        document.addEventListener('click', (e) => {
            if (e.target.matches('[onclick*="goToPage"]')) {
                const page = parseInt(e.target.getAttribute('onclick').match(/goToPage\((\d+)\)/)[1]);
                this.goToPage(page);
            }
        });
    }

    initializeComponents() {
        // Initialize tooltips
        this.initTooltips();
        
        // Initialize search
        this.initSearch();
        
        // Initialize theme
        this.initTheme();
    }

    // Sidebar functionality
    toggleSidebar() {
        const sidebar = document.querySelector('.sidebar');
        if (sidebar) {
            sidebar.classList.toggle('hidden');
            sidebar.classList.toggle('lg:block');
        }
    }

    toggleSidebarCollapse() {
        const sidebar = document.querySelector('.sidebar');
        if (sidebar) {
            sidebar.classList.toggle('sidebar-collapsed');
            localStorage.setItem('sidebarCollapsed', sidebar.classList.contains('sidebar-collapsed'));
        }
    }

    // Navigation functionality
    toggleNavGroup(groupId) {
        const submenu = document.getElementById(`nav-group-${groupId}`);
        const button = document.querySelector(`[onclick*="toggleNavGroup('${groupId}')"]`);
        const icon = button?.querySelector('i[class*="chevron"]');
        
        if (submenu) {
            submenu.classList.toggle('hidden');
            if (icon) {
                icon.classList.toggle('rotate-90');
            }
        }
    }

    // Dropdown functionality
    toggleDropdown(selector) {
        const dropdown = document.querySelector(selector);
        if (dropdown) {
            dropdown.classList.toggle('hidden');
        }
    }

    closeDropdown(selector) {
        const dropdown = document.querySelector(selector);
        if (dropdown) {
            dropdown.classList.add('hidden');
        }
    }

    // Modal functionality
    openModal(modalSelector) {
        const modal = document.querySelector(modalSelector);
        if (modal) {
            modal.classList.add('modal-open');
            document.body.classList.add('overflow-hidden');
        }
    }

    closeModal() {
        const modals = document.querySelectorAll('.modal');
        modals.forEach(modal => {
            modal.classList.remove('modal-open');
        });
        document.body.classList.remove('overflow-hidden');
    }

    // Table functionality
    editItem(id) {
        console.log('Edit item:', id);
        // Implement edit functionality
        this.showNotification('Edit functionality not implemented yet', 'info');
    }

    deleteItem(id) {
        if (confirm('Are you sure you want to delete this item?')) {
            console.log('Delete item:', id);
            // Implement delete functionality
            this.showNotification('Delete functionality not implemented yet', 'info');
        }
    }

    // Pagination functionality
    goToPage(page) {
        console.log('Go to page:', page);
        // Implement pagination
        this.showNotification(`Going to page ${page}`, 'info');
    }

    // Search functionality
    initSearch() {
        const searchInput = document.querySelector('input[placeholder*="Search"]');
        if (searchInput) {
            let debounceTimer;
            searchInput.addEventListener('input', (e) => {
                clearTimeout(debounceTimer);
                debounceTimer = setTimeout(() => {
                    this.performSearch(e.target.value);
                }, 300);
            });
        }
    }

    performSearch(query) {
        console.log('Search query:', query);
        // Implement search functionality
    }

    // Notification functionality
    showNotification(message, type = 'info') {
        const notification = document.createElement('div');
        notification.className = `fixed top-4 right-4 z-50 p-4 rounded-lg shadow-lg max-w-sm transform transition-transform duration-300 translate-x-full`;
        
        const colors = {
            success: 'bg-green-500 text-white',
            error: 'bg-red-500 text-white',
            warning: 'bg-yellow-500 text-white',
            info: 'bg-blue-500 text-white'
        };
        
        notification.className += ` ${colors[type] || colors.info}`;
        notification.innerHTML = `
            <div class="flex items-center justify-between">
                <span>${message}</span>
                <button onclick="this.parentElement.parentElement.remove()" class="ml-4 text-white hover:text-gray-200">
                    <i class="fas fa-times"></i>
                </button>
            </div>
        `;
        
        document.body.appendChild(notification);
        
        // Animate in
        setTimeout(() => {
            notification.classList.remove('translate-x-full');
        }, 100);
        
        // Auto remove after 5 seconds
        setTimeout(() => {
            notification.classList.add('translate-x-full');
            setTimeout(() => {
                notification.remove();
            }, 300);
        }, 5000);
    }

    // Theme functionality
    initTheme() {
        const savedTheme = localStorage.getItem('theme');
        if (savedTheme) {
            document.documentElement.classList.toggle('dark', savedTheme === 'dark');
        }
    }

    toggleTheme() {
        document.documentElement.classList.toggle('dark');
        const isDark = document.documentElement.classList.contains('dark');
        localStorage.setItem('theme', isDark ? 'dark' : 'light');
    }

    // Tooltip functionality
    initTooltips() {
        const tooltipElements = document.querySelectorAll('[data-tooltip]');
        tooltipElements.forEach(element => {
            element.addEventListener('mouseenter', this.showTooltip);
            element.addEventListener('mouseleave', this.hideTooltip);
        });
    }

    showTooltip(e) {
        const text = e.target.getAttribute('data-tooltip');
        const tooltip = document.createElement('div');
        tooltip.className = 'absolute z-50 px-2 py-1 text-sm text-white bg-gray-900 rounded shadow-lg tooltip';
        tooltip.textContent = text;
        
        document.body.appendChild(tooltip);
        
        const rect = e.target.getBoundingClientRect();
        tooltip.style.left = rect.left + (rect.width / 2) - (tooltip.offsetWidth / 2) + 'px';
        tooltip.style.top = rect.top - tooltip.offsetHeight - 8 + 'px';
    }

    hideTooltip() {
        const tooltip = document.querySelector('.tooltip');
        if (tooltip) {
            tooltip.remove();
        }
    }

    // Utility functions
    markAllAsRead() {
        console.log('Mark all notifications as read');
        this.showNotification('All notifications marked as read', 'success');
    }

    openAddUserModal() {
        console.log('Open add user modal');
        this.showNotification('Add user modal would open here', 'info');
    }

    generateReport() {
        console.log('Generate report');
        this.showNotification('Generating report...', 'info');
        setTimeout(() => {
            this.showNotification('Report generated successfully!', 'success');
        }, 2000);
    }
}

// Chart initialization
function initializeCharts(data) {
    // Revenue Chart
    const revenueCtx = document.getElementById('revenueChart');
    if (revenueCtx) {
        new Chart(revenueCtx, {
            type: 'line',
            data: {
                labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun'],
                datasets: [{
                    label: 'Revenue',
                    data: [12000, 19000, 15000, 25000, 22000, 30000],
                    borderColor: 'rgb(79, 70, 229)',
                    backgroundColor: 'rgba(79, 70, 229, 0.1)',
                    tension: 0.4,
                    fill: true
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: {
                        display: false
                    }
                },
                scales: {
                    y: {
                        beginAtZero: true,
                        ticks: {
                            callback: function(value) {
                                return '$' + value.toLocaleString();
                            }
                        }
                    }
                }
            }
        });
    }

    // User Chart
    const userCtx = document.getElementById('userChart');
    if (userCtx) {
        new Chart(userCtx, {
            type: 'doughnut',
            data: {
                labels: ['Active Users', 'Inactive Users'],
                datasets: [{
                    data: [data.userStats.active, data.userStats.inactive],
                    backgroundColor: [
                        'rgb(34, 197, 94)',
                        'rgb(156, 163, 175)'
                    ],
                    borderWidth: 0
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: {
                        position: 'bottom'
                    }
                }
            }
        });
    }
}

// Global functions for onclick handlers
window.toggleSidebar = () => window.dashboard.toggleSidebar();
window.toggleSidebarCollapse = () => window.dashboard.toggleSidebarCollapse();
window.toggleNavGroup = (id) => window.dashboard.toggleNavGroup(id);
window.toggleNotifications = () => window.dashboard.toggleDropdown('.notification-dropdown');
window.toggleUserMenu = () => window.dashboard.toggleDropdown('.user-dropdown');
window.closeModal = () => window.dashboard.closeModal();
window.editItem = (id) => window.dashboard.editItem(id);
window.deleteItem = (id) => window.dashboard.deleteItem(id);
window.goToPage = (page) => window.dashboard.goToPage(page);
window.markAllAsRead = () => window.dashboard.markAllAsRead();
window.openAddUserModal = () => window.dashboard.openAddUserModal();
window.generateReport = () => window.dashboard.generateReport();
