function initCvMenu() {
    const toggle = document.getElementById('cv-toggle');
    const menu = document.getElementById('cv-menu');

    if (!toggle || !menu) {
        console.error('CV menu elements not found');
        return;
    }

    const accordionToggles = document.querySelectorAll('.accordion-toggle');

    // Collapsed accordion panels keep the `hidden` class (display: none),
    // so they stay out of both the viewport and the accessibility tree.
    const panelFor = (button) => {
        const id = button.getAttribute('aria-controls');
        return id ? document.getElementById(id) : null;
    };

    accordionToggles.forEach(button => {
        button.addEventListener('click', () => {
            const panel = panelFor(button);
            if (!panel) return;

            // Toggle visibility and ARIA state (aria-expanded drives the CSS expand rule)
            const isExpanded = button.getAttribute('aria-expanded') === 'true';
            button.setAttribute('aria-expanded', String(!isExpanded));

            panel.classList.toggle('hidden');
        });
    });

    const closeMenu = () => {
        menu.classList.remove('open');
        toggle.setAttribute('aria-expanded', 'false');
    };

    const openMenu = () => {
        menu.classList.add('open');
        toggle.setAttribute('aria-expanded', 'true');
    };

    toggle.addEventListener('click', (event) => {
        event.stopPropagation();
        if (menu.classList.contains('open')) {
            closeMenu();
        } else {
            // When opening the main menu, ensure all accordions start collapsed
            accordionToggles.forEach(btn => {
                btn.setAttribute('aria-expanded', 'false');
                const panel = panelFor(btn);
                if (panel) {
                    panel.classList.add('hidden');
                }
            });
            openMenu();
        }
    });

    document.addEventListener('click', (event) => {
        if (!menu.contains(event.target) && !toggle.contains(event.target)) {
            closeMenu();
        }
    });

    document.addEventListener('keydown', (event) => {
        if (event.key === 'Escape') {
            closeMenu();
        }
    });
}

// Rocket Loader (and similar async script rewriters) can delay execution
// until after DOMContentLoaded has already fired, which means a plain
// 'DOMContentLoaded' listener registered too late never runs. Checking
// readyState covers both cases: run now if the DOM is already parsed,
// otherwise wait for the event like normal.
if (document.readyState === 'loading') {
    window.addEventListener('DOMContentLoaded', initCvMenu);
} else {
    initCvMenu();
}
