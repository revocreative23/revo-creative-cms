// Navbar — scroll state (transparent dark → solid white once scrolled past hero)
const navbarEl = document.querySelector(".navbar");
if (navbarEl) {
  const SCROLL_THRESHOLD = 80;
  const updateNavbarState = () => {
    if (window.scrollY > SCROLL_THRESHOLD) {
      navbarEl.classList.add("scrolled");
    } else {
      navbarEl.classList.remove("scrolled");
    }
  };
  window.addEventListener("scroll", updateNavbarState, { passive: true });
  updateNavbarState();
}

// Navbar — mobile toggle
const navToggle = document.querySelector(".nav-toggle");
const navLinks = document.querySelector(".nav-links");

if (navToggle && navLinks) {
  navToggle.addEventListener("click", () => {
    navToggle.classList.toggle("active");
    navLinks.classList.toggle("open");
  });

  navLinks.querySelectorAll("a").forEach((link) => {
    link.addEventListener("click", () => {
      navToggle.classList.remove("active");
      navLinks.classList.remove("open");
    });
  });
}

// Highlight active nav link based on page
const navAnchors = document.querySelectorAll(".nav-links a");
const path = window.location.pathname.split("/").pop() || "index.html";
navAnchors.forEach((a) => {
  const href = a.getAttribute("href");
  if (!href) return;
  if (href === path || (path === "" && href === "index.html")) {
    a.classList.add("active");
  }
});

// Reveal on scroll
const revealItems = document.querySelectorAll(".reveal");
if ("IntersectionObserver" in window && revealItems.length) {
  const io = new IntersectionObserver(
    (entries) => {
      entries.forEach((e) => {
        if (e.isIntersecting) {
          e.target.classList.add("is-visible");
          io.unobserve(e.target);
        }
      });
    },
    { threshold: 0.12 }
  );
  revealItems.forEach((el) => io.observe(el));
}

// Portfolio filter
const filters = document.querySelectorAll(".filter-pill");
const portfolioItems = document.querySelectorAll(".portfolio-card");

filters.forEach((pill) => {
  pill.addEventListener("click", () => {
    filters.forEach((p) => p.classList.remove("active"));
    pill.classList.add("active");
    const cat = pill.dataset.filter;
    portfolioItems.forEach((card) => {
      const cardCat = card.dataset.category;
      if (cat === "all" || cardCat === cat) {
        card.style.display = "";
      } else {
        card.style.display = "none";
      }
    });
  });
});

// Contact form (front-end only stub)
const contactForm = document.querySelector(".contact-form");
if (contactForm) {
  contactForm.addEventListener("submit", (e) => {
    e.preventDefault();
    const success = contactForm.querySelector(".form-success");
    if (success) {
      success.classList.add("show");
      contactForm.reset();
      setTimeout(() => success.classList.remove("show"), 5000);
    }
  });
}

// Products carousel — auto-rotating
const productsCarousel = document.querySelector(".products-carousel");
if (productsCarousel) {
  const track = productsCarousel.querySelector(".products-track");
  const slides = productsCarousel.querySelectorAll(".product-slide");
  const dots = productsCarousel.querySelectorAll(".carousel-dot");
  const prevBtn = productsCarousel.querySelector(".carousel-arrow.prev");
  const nextBtn = productsCarousel.querySelector(".carousel-arrow.next");
  const total = slides.length;
  const ROTATE_MS = 5000;
  let current = 0;
  let timer = null;

  const goTo = (index) => {
    current = (index + total) % total;
    track.style.transform = `translateX(-${current * 100}%)`;
    dots.forEach((d, i) => d.classList.toggle("active", i === current));
  };

  const next = () => goTo(current + 1);
  const prev = () => goTo(current - 1);

  const start = () => {
    stop();
    timer = setInterval(next, ROTATE_MS);
  };
  const stop = () => {
    if (timer) {
      clearInterval(timer);
      timer = null;
    }
  };

  nextBtn && nextBtn.addEventListener("click", () => { next(); start(); });
  prevBtn && prevBtn.addEventListener("click", () => { prev(); start(); });
  dots.forEach((dot) => {
    dot.addEventListener("click", () => {
      goTo(parseInt(dot.dataset.slide, 10));
      start();
    });
  });

  // Pause on hover (desktop)
  productsCarousel.addEventListener("mouseenter", stop);
  productsCarousel.addEventListener("mouseleave", start);

  // Pause when tab not visible — avoids drift
  document.addEventListener("visibilitychange", () => {
    document.hidden ? stop() : start();
  });

  // Touch swipe (mobile)
  let touchStartX = 0;
  track.addEventListener("touchstart", (e) => {
    touchStartX = e.changedTouches[0].screenX;
    stop();
  }, { passive: true });
  track.addEventListener("touchend", (e) => {
    const diff = e.changedTouches[0].screenX - touchStartX;
    if (Math.abs(diff) > 40) diff > 0 ? prev() : next();
    start();
  }, { passive: true });

  start();
}

// Year in footer
const yearEl = document.getElementById("year");
if (yearEl) yearEl.textContent = new Date().getFullYear();
