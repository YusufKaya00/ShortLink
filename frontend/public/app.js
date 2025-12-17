// API Configuration
const API_BASE = {
    user: 'https://shortlink-b3tq.onrender.com',
    url: 'https://url-service-p46t.onrender.com',
    stats: 'https://stats-service-jlys.onrender.com'
};

// State
let token = localStorage.getItem('token');
let user = JSON.parse(localStorage.getItem('user') || 'null');

// Initialize
document.addEventListener('DOMContentLoaded', () => {
    updateAuthUI();
    loadOverallStats();
    loadRecentLinks();

    // Form submission
    document.getElementById('shorten-form').addEventListener('submit', handleShorten);
});

// ===== Auth Functions =====

function updateAuthUI() {
    const authButtons = document.getElementById('auth-buttons');
    const userMenu = document.getElementById('user-menu');
    const userName = document.getElementById('user-name');

    if (token && user) {
        authButtons.classList.add('hidden');
        userMenu.classList.remove('hidden');
        userName.textContent = user.name || user.email;
    } else {
        authButtons.classList.remove('hidden');
        userMenu.classList.add('hidden');
    }
}

async function handleLogin(e) {
    e.preventDefault();

    const email = document.getElementById('login-email').value;
    const password = document.getElementById('login-password').value;

    try {
        const response = await fetch(`${API_BASE.user}/api/users/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password })
        });

        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.error || 'Login failed');
        }

        token = data.token;
        user = data.user;
        localStorage.setItem('token', token);
        localStorage.setItem('user', JSON.stringify(user));

        closeModals();
        updateAuthUI();
        showToast('Welcome back! 👋');
        loadRecentLinks();
    } catch (error) {
        showToast(error.message, 'error');
    }
}

async function handleRegister(e) {
    e.preventDefault();

    const name = document.getElementById('register-name').value;
    const email = document.getElementById('register-email').value;
    const password = document.getElementById('register-password').value;

    try {
        const response = await fetch(`${API_BASE.user}/api/users/register`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name, email, password })
        });

        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.error || 'Registration failed');
        }

        token = data.token;
        user = data.user;
        localStorage.setItem('token', token);
        localStorage.setItem('user', JSON.stringify(user));

        closeModals();
        updateAuthUI();
        showToast('Account created successfully! 🎉');
    } catch (error) {
        showToast(error.message, 'error');
    }
}

function logout() {
    token = null;
    user = null;
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    updateAuthUI();
    showToast('Logged out successfully');
    loadRecentLinks();
}

// ===== URL Functions =====

async function handleShorten(e) {
    e.preventDefault();

    const urlInput = document.getElementById('url-input');
    const customCodeInput = document.getElementById('custom-code-input');
    const customCodeToggle = document.getElementById('custom-code-toggle');

    const body = {
        original_url: urlInput.value
    };

    if (customCodeToggle.checked && customCodeInput.value) {
        body.custom_code = customCodeInput.value;
    }

    try {
        const headers = { 'Content-Type': 'application/json' };
        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }

        const response = await fetch(`${API_BASE.url}/api/urls`, {
            method: 'POST',
            headers,
            body: JSON.stringify(body)
        });

        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.error || 'Failed to shorten URL');
        }

        // Show result
        const resultDiv = document.getElementById('result');
        const shortUrlLink = document.getElementById('short-url');

        shortUrlLink.href = data.short_url;
        shortUrlLink.textContent = data.short_url;
        resultDiv.classList.remove('hidden');

        // Reset form
        urlInput.value = '';
        customCodeInput.value = '';
        customCodeToggle.checked = false;
        customCodeInput.classList.add('hidden');

        // Reload links
        loadRecentLinks();
        loadOverallStats();

        showToast('Link shortened successfully! 🔗');
    } catch (error) {
        showToast(error.message, 'error');
    }
}

function toggleCustomCode() {
    const toggle = document.getElementById('custom-code-toggle');
    const input = document.getElementById('custom-code-input');

    if (toggle.checked) {
        input.classList.remove('hidden');
        input.focus();
    } else {
        input.classList.add('hidden');
        input.value = '';
    }
}

async function copyToClipboard() {
    const shortUrl = document.getElementById('short-url').textContent;

    try {
        await navigator.clipboard.writeText(shortUrl);
        showToast('Copied to clipboard! 📋');
    } catch (error) {
        showToast('Failed to copy', 'error');
    }
}

async function copyLink(shortUrl) {
    try {
        await navigator.clipboard.writeText(shortUrl);
        showToast('Copied to clipboard! 📋');
    } catch (error) {
        showToast('Failed to copy', 'error');
    }
}

// ===== Stats Functions =====

async function loadOverallStats() {
    try {
        const response = await fetch(`${API_BASE.stats}/api/stats/overall`);
        const data = await response.json();

        document.getElementById('total-urls').textContent = formatNumber(data.total_urls || 0);
        document.getElementById('total-clicks').textContent = formatNumber(data.total_clicks || 0);
        document.getElementById('today-clicks').textContent = formatNumber(data.today_clicks || 0);
    } catch (error) {
        console.error('Failed to load stats:', error);
    }
}

async function loadRecentLinks() {
    try {
        const response = await fetch(`${API_BASE.url}/api/urls/all?limit=6`);
        const data = await response.json();

        const container = document.getElementById('recent-links');

        if (!data.urls || data.urls.length === 0) {
            container.innerHTML = `
                <div class="link-card" style="text-align: center; grid-column: 1 / -1;">
                    <p style="color: var(--text-muted);">No links yet. Create your first short link above! 🚀</p>
                </div>
            `;
            return;
        }

        container.innerHTML = data.urls.map(url => `
            <div class="link-card">
                <div class="link-header">
                    <a href="${url.short_url}" target="_blank" class="link-short">${url.short_code}</a>
                    <span class="link-clicks">
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path>
                            <circle cx="12" cy="12" r="3"></circle>
                        </svg>
                        ${formatNumber(url.click_count)}
                    </span>
                </div>
                <p class="link-original" title="${url.original_url}">${url.original_url}</p>
                <div class="link-actions">
                    <button class="btn btn-ghost" onclick="copyLink('${url.short_url}')">
                        Copy
                    </button>
                    <button class="btn btn-ghost" onclick="showStats('${url.short_code}')">
                        Stats
                    </button>
                </div>
            </div>
        `).join('');
    } catch (error) {
        console.error('Failed to load links:', error);
    }
}

async function showStats(shortCode) {
    try {
        const response = await fetch(`${API_BASE.stats}/api/stats/${shortCode}`);
        const data = await response.json();

        const content = document.getElementById('stats-content');

        content.innerHTML = `
            <div class="stats-detail">
                <div class="stats-row">
                    <div class="stats-box">
                        <h4>Total Clicks</h4>
                        <div class="stat-value" style="font-size: 2rem;">${formatNumber(data.total_clicks)}</div>
                    </div>
                    <div class="stats-box">
                        <h4>Short Code</h4>
                        <div style="font-size: 1.25rem; color: var(--primary-light);">${data.short_code}</div>
                    </div>
                </div>
                <div class="stats-row">
                    <div class="stats-box">
                        <h4>By Device</h4>
                        <ul class="stats-list">
                            ${(data.by_device || []).map(d => `
                                <li><span>${d.device}</span><span>${d.count}</span></li>
                            `).join('') || '<li>No data yet</li>'}
                        </ul>
                    </div>
                    <div class="stats-box">
                        <h4>By Browser</h4>
                        <ul class="stats-list">
                            ${(data.by_browser || []).map(b => `
                                <li><span>${b.browser}</span><span>${b.count}</span></li>
                            `).join('') || '<li>No data yet</li>'}
                        </ul>
                    </div>
                </div>
                <div class="stats-box">
                    <h4>Top Referrers</h4>
                    <ul class="stats-list">
                        ${(data.by_referer || []).map(r => `
                            <li><span>${r.referer || 'Direct'}</span><span>${r.count}</span></li>
                        `).join('') || '<li>No data yet</li>'}
                    </ul>
                </div>
            </div>
        `;

        showModal('stats');
    } catch (error) {
        showToast('Failed to load stats', 'error');
    }
}

// ===== Modal Functions =====

function showModal(type) {
    closeModals();
    document.getElementById(`${type}-modal`).classList.remove('hidden');
    document.body.style.overflow = 'hidden';
}

function closeModals() {
    document.querySelectorAll('.modal').forEach(m => m.classList.add('hidden'));
    document.body.style.overflow = '';
}

// Close modal on Escape key
document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape') closeModals();
});

// ===== Toast Function =====

function showToast(message, type = 'success') {
    const toast = document.getElementById('toast');
    toast.textContent = message;
    toast.style.background = type === 'error'
        ? 'linear-gradient(135deg, #ef4444, #dc2626)'
        : 'linear-gradient(135deg, var(--primary), var(--primary-dark))';
    toast.classList.remove('hidden');

    setTimeout(() => {
        toast.classList.add('hidden');
    }, 3000);
}

// ===== Utility Functions =====

function formatNumber(num) {
    if (num >= 1000000) {
        return (num / 1000000).toFixed(1) + 'M';
    }
    if (num >= 1000) {
        return (num / 1000).toFixed(1) + 'K';
    }
    return num.toString();
}
