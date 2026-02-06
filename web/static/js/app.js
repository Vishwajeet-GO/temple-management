const API_BASE = '/api/v1';

function formatCurrency(amount) {
    return '₹' + Number(amount).toLocaleString('en-IN');
}

function formatDate(dateString) {
    if (!dateString) return '-';
    const date = new Date(dateString);
    return date.toLocaleDateString('en-IN');
}

async function apiCall(endpoint, method = 'GET', data = null) {
    const options = {
        method,
        headers: {
            'Content-Type': 'application/json'
        }
    };
    if (data) {
        options.body = JSON.stringify(data);
    }
    const response = await fetch(API_BASE + endpoint, options);
    return await response.json();
}

async function loadDashboard() {
    try {
        const summary = await apiCall('/dashboard/summary');
        if (summary.success) {
            document.getElementById('totalDonations').textContent = formatCurrency(summary.data.total_donations);
            document.getElementById('totalExpenses').textContent = formatCurrency(summary.data.total_expenses);
            document.getElementById('balance').textContent = formatCurrency(summary.data.balance);
            document.getElementById('upcomingFestivals').textContent = summary.data.upcoming_festivals;
        }

        const donations = await apiCall('/dashboard/recent-donations');
        if (donations.success) {
            const container = document.getElementById('recentDonations');
            if (donations.data && donations.data.length > 0) {
                container.innerHTML = donations.data.slice(0, 5).map(d => `
                    <div class="list-group-item d-flex justify-content-between align-items-center">
                        <div>
                            <strong>${d.donor_name || 'Anonymous'}</strong>
                            <small class="text-muted d-block">${formatDate(d.date)}</small>
                        </div>
                        <span class="badge bg-success">${formatCurrency(d.amount)}</span>
                    </div>
                `).join('');
            } else {
                container.innerHTML = '<div class="list-group-item text-center text-muted">No donations yet</div>';
            }
        }

        const expenses = await apiCall('/dashboard/recent-expenses');
        if (expenses.success) {
            const container = document.getElementById('recentExpenses');
            if (expenses.data && expenses.data.length > 0) {
                container.innerHTML = expenses.data.slice(0, 5).map(e => `
                    <div class="list-group-item d-flex justify-content-between align-items-center">
                        <div>
                            <strong>${e.title}</strong>
                            <small class="text-muted d-block">${e.category}</small>
                        </div>
                        <span class="badge bg-danger">${formatCurrency(e.amount)}</span>
                    </div>
                `).join('');
            } else {
                container.innerHTML = '<div class="list-group-item text-center text-muted">No expenses yet</div>';
            }
        }
    } catch (error) {
        console.error('Error loading dashboard:', error);
    }
}

async function loadTempleInfo() {
    try {
        const response = await apiCall('/temple');
        if (response.success) {
            document.getElementById('templeName').textContent = '🛕 ' + response.data.name;
            document.getElementById('templeAddress').textContent =
                `${response.data.address}, ${response.data.city}, ${response.data.state}`;
        }
    } catch (error) {
        console.error('Error loading temple info:', error);
    }
}