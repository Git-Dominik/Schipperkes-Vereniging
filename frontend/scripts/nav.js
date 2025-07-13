(function initNavbar() {
  const mobileToggle = document.getElementById('mobileToggle');
  const navLinks = document.getElementById('navLinks');
  const navItems = document.querySelectorAll('.nav-item');

  if (!mobileToggle || !navLinks || navItems.length === 0) return;

  mobileToggle.addEventListener('click', function () {
    navLinks.classList.toggle('active');
    const icon = mobileToggle.querySelector('i');
    if (navLinks.classList.contains('active')) {
      icon.classList.replace('fa-bars', 'fa-times');
    } else {
      icon.classList.replace('fa-times', 'fa-bars');
    }
  });

  function setupMobileDropdowns() {
    navItems.forEach(item => {
      const link = item.querySelector('.nav-link');
      const dropdown = item.querySelector('.dropdown');
      if (dropdown) {
        link.removeEventListener('click', handleMobileLinkClick);
        link.addEventListener('click', handleMobileLinkClick);
      }
    });
  }

  function handleMobileLinkClick(e) {
    if (window.innerWidth <= 992) {
      e.preventDefault();
      const dropdown = this.parentElement.querySelector('.dropdown');
      const icon = this.querySelector('.nav-icon');

      dropdown.classList.toggle('show');

      if (dropdown.classList.contains('show')) {
        icon.classList.replace('fa-chevron-down', 'fa-chevron-up');
      } else {
        icon.classList.replace('fa-chevron-up', 'fa-chevron-down');
      }

      navItems.forEach(otherItem => {
        if (otherItem !== this.parentElement) {
          const otherDropdown = otherItem.querySelector('.dropdown');
          if (otherDropdown && otherDropdown.classList.contains('show')) {
            otherDropdown.classList.remove('show');
            const otherIcon = otherItem.querySelector('.nav-icon');
            if (otherIcon) {
              otherIcon.classList.replace('fa-chevron-up', 'fa-chevron-down');
            }
          }
        }
      });
    }
  }

  function initDropdowns() {
    if (window.innerWidth <= 992) {
      setupMobileDropdowns();
    } else {
      navItems.forEach(item => {
        const link = item.querySelector('.nav-link');
        if (link) {
          link.removeEventListener('click', handleMobileLinkClick);
        }
      });
    }
  }

  initDropdowns();

  window.addEventListener('resize', function () {
    if (window.innerWidth > 992) {
      navLinks.classList.remove('active');
      const icon = mobileToggle.querySelector('i');
      icon.classList.replace('fa-times', 'fa-bars');

      document.querySelectorAll('.dropdown').forEach(dropdown => {
        dropdown.classList.remove('show');
      });

      document.querySelectorAll('.nav-icon').forEach(icon => {
        icon.classList.replace('fa-chevron-up', 'fa-chevron-down');
      });
    }
    initDropdowns();
  });

  // Touch ripple effect
  if ('ontouchstart' in window) {
    document.querySelectorAll('.nav-link, .dropdown-link, .contact-button').forEach(link => {
      link.addEventListener('touchstart', function (e) {
        const rect = link.getBoundingClientRect();
        const ripple = document.createElement('span');
        ripple.classList.add('ripple');
        ripple.style.left = `${e.touches[0].clientX - rect.left}px`;
        ripple.style.top = `${e.touches[0].clientY - rect.top}px`;
        link.appendChild(ripple);
        setTimeout(() => ripple.remove(), 600);
      });
    });
  }
})();
