window.addEventListener('DOMContentLoaded', () => {
    const toggle = document.getElementById('cv-toggle');
    const menu = document.getElementById('cv-menu');

    if (!toggle || !menu) {
        console.error('CV menu elements not found');
        return;
    }

    const accordionToggles = document.querySelectorAll('.accordion-toggle');

    accordionToggles.forEach(button => {
        button.addEventListener('click', (event) => {
            const parentLi = button.closest('li[role="none"]'); // Get the parent <li>
            const content = parentLi ? parentLi.querySelector('.accordion-content') : null;

            if (!content) return;

            // Toggle visibility and ARIA state
            const isExpanded = button.getAttribute('aria-expanded') === 'true' || false;
            button.setAttribute('aria-expanded', !isExpanded);
            
            content.classList.toggle('hidden');
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
            // When opening the main menu, ensure all accordions are closed initially
            accordionToggles.forEach(btn => btn.setAttribute('aria-expanded', 'false'));
            menu.classList.add('open');
            toggle.setAttribute('aria-expanded', 'true');
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
});
