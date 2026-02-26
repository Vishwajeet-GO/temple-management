// Current Language
let currentLang = localStorage.getItem('site_lang') || 'en';

document.addEventListener('DOMContentLoaded', () => {
    applyLanguage(currentLang);
    renderLangButton();
});

function toggleLanguage() {
    currentLang = currentLang === 'en' ? 'hi' : 'en';
    localStorage.setItem('site_lang', currentLang);
    applyLanguage(currentLang);
    renderLangButton();
}

function applyLanguage(lang) {
    // 1. Update simple text elements
    document.querySelectorAll('[data-lang]').forEach(el => {
        const key = el.getAttribute('data-lang');
        if (translations[lang][key]) {
            // Handle inputs with placeholders vs regular text
            if (el.tagName === 'INPUT' || el.tagName === 'TEXTAREA') {
                el.placeholder = translations[lang][key];
            } else {
                // If element has icon, keep icon and update text
                const icon = el.querySelector('i');
                if (icon) {
                    el.innerHTML = ''; // Clear
                    el.appendChild(icon); // Add icon back
                    el.append(' ' + translations[lang][key]); // Add text
                } else {
                    el.innerText = translations[lang][key];
                }
            }
        }
    });

    // 2. Update specific font for Hindi (optional, makes it look better)
    if (lang === 'hi') {
        document.body.classList.add('lang-hi');
    } else {
        document.body.classList.remove('lang-hi');
    }
}

function renderLangButton() {
    const nav = document.querySelector('.top-navbar .nav-buttons') || document.querySelector('.top-nav .nav-btns');
    if (!nav) return;

    let btn = document.getElementById('langToggleBtn');
    if (!btn) {
        btn = document.createElement('button');
        btn.id = 'langToggleBtn';
        btn.className = 'nav-btn';
        btn.onclick = toggleLanguage;
        // Insert before the last button (usually logout or login)
        nav.insertBefore(btn, nav.firstChild);
    }

    // Update button look
    if (currentLang === 'en') {
        btn.innerHTML = '<i class="fas fa-language"></i> HI';
        btn.style.background = 'rgba(255,255,255,0.1)';
        btn.style.border = '1px solid rgba(255,255,255,0.3)';
    } else {
        btn.innerHTML = '<i class="fas fa-language"></i> EN';
        btn.style.background = 'linear-gradient(135deg, #ff9800, #f57c00)';
        btn.style.border = 'none';
    }
}