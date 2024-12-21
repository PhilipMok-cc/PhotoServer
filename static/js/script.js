document.addEventListener('DOMContentLoaded', function() {
    let images = document.querySelectorAll('img');
    let indexDisplay = document.getElementById('indexDisplay');
    let totalImages = images.length;
    let currentIndex = 0;

    function updateIndexDisplay(index) {
        if (indexDisplay) {
            indexDisplay.textContent = `Image ${index + 1} of ${totalImages}`;
        }
    }

    function scrollToImage(index) {
        if (index >= 0 && index < images.length) {
            images[index].scrollIntoView({behavior: 'smooth'});
            currentIndex = index;
            updateIndexDisplay(currentIndex);
        }
    }

    function onScroll() {
        let scrollPosition = window.scrollY + window.innerHeight / 2;
        images.forEach((img, index) => {
            let imgTop = img.offsetTop;
            let imgBottom = imgTop + img.offsetHeight;
            if (scrollPosition >= imgTop && scrollPosition < imgBottom) {
                if (currentIndex !== index) {
                    currentIndex = index;
                    updateIndexDisplay(currentIndex);
                }
            }
        });
    }

    document.addEventListener('scroll', onScroll);

    document.addEventListener('keydown', function(e) {
        if (e.key === 'ArrowDown' || e.key === 'PageDown' || e.key === ' ') {
            e.preventDefault();
            scrollToImage(currentIndex + 1);
        } else if (e.key === 'ArrowUp' || e.key === 'PageUp') {
            e.preventDefault();
            scrollToImage(currentIndex - 1);
        }
    });

    let touchStartX = 0;
    let touchEndX = 0;
    const minSwipeDistance = 30; // Minimum distance for a swipe to be detected
    let lastTapTime = 0;
    const doubleTapThreshold = 300; // Maximum time between taps for a double tap (in milliseconds)

    function handleGesture() {
        console.log(`Touch start: ${touchStartX}, Touch end: ${touchEndX}`);
        if (Math.abs(touchEndX - touchStartX) > minSwipeDistance) {
            if (touchEndX < touchStartX) {
                scrollToImage(currentIndex + 1);
            }
            if (touchEndX > touchStartX) {
                scrollToImage(currentIndex - 1);
            }
        }
        // Reset touch coordinates
        touchStartX = 0;
        touchEndX = 0;
    }

    function handleDoubleTap() {
        scrollToImage(currentIndex + 1);
    }

    document.addEventListener('touchstart', function(e) {
        touchStartX = e.changedTouches[0].screenX;
        console.log(`Touch start detected at: ${touchStartX}`);

        const currentTime = new Date().getTime();
        const tapInterval = currentTime - lastTapTime;
        if (tapInterval < doubleTapThreshold && tapInterval > 0) {
            handleDoubleTap();
        }
        lastTapTime = currentTime;
    });

    document.addEventListener('touchend', function(e) {
        touchEndX = e.changedTouches[0].screenX;
        console.log(`Touch end detected at: ${touchEndX}`);
        handleGesture();
    });

    // Initialize the display
    updateIndexDisplay(currentIndex);
});
