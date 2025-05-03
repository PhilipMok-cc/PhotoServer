document.addEventListener('DOMContentLoaded', function() {
    // Selectors and state variables
    let images = document.querySelectorAll('img');
    let chapters = document.querySelectorAll('.chapter');
    let indexDisplay = document.getElementById('indexDisplay');
    let totalImages = images.length;
    let currentIndex = 0;

    // Constants for gesture detection
    const minSwipeDistance = 30; // Minimum horizontal swipe distance
    const minVerticalSwipeDistance = 50; // Minimum vertical swipe distance
    const doubleTapThreshold = 300; // Maximum time between taps for a double tap (ms)

    // Variables for touch gestures
    let touchStartX = 0, touchEndX = 0, touchStartY = 0, touchEndY = 0;
    let lastTapTime = 0;

    /**
     * Get the current chapter index based on the current image index.
     * @param {number} index - The index of the current image.
     * @returns {number} - The index of the current chapter.
     */
    function getCurrentChapterIndex(index) {
        for (let i = 0; i < chapters.length; i++) {
            let chapterImages = chapters[i].querySelectorAll('img');
            if (Array.from(chapterImages).includes(images[index])) {
                return i;
            }
        }
        return -1;
    }

    /**
     * Scroll to a specific image by index.
     * @param {number} index - The index of the image to scroll to.
     */
    function scrollToImage(index) {
        if (index >= 0 && index < images.length) {
            images[index].scrollIntoView({ behavior: 'auto', block: 'start' });
            currentIndex = index;
            updateIndexDisplay(currentIndex);
        }
    }

    /**
     * Scroll to the first image of the next chapter.
     */
    function scrollToNextChapter() {
        let currentChapterIndex = getCurrentChapterIndex(currentIndex);
        if (currentChapterIndex >= 0 && currentChapterIndex < chapters.length - 1) {
            let nextChapterImages = chapters[currentChapterIndex + 1].querySelectorAll('img');
            if (nextChapterImages.length > 0) {
                scrollToImage(Array.from(images).indexOf(nextChapterImages[0]));
            }
        }
    }

    /**
     * Scroll to the first image of the previous chapter or the current chapter.
     */
    function scrollToPreviousChapter() {
        let currentChapterIndex = getCurrentChapterIndex(currentIndex);
        if (currentChapterIndex > 0) {
            let currentChapterImages = chapters[currentChapterIndex].querySelectorAll('img');
            if (currentIndex !== Array.from(images).indexOf(currentChapterImages[0])) {
                scrollToImage(Array.from(images).indexOf(currentChapterImages[0])); // Scroll to the first image of the current chapter
            } else {
                let previousChapterImages = chapters[currentChapterIndex - 1].querySelectorAll('img');
                if (previousChapterImages.length > 0) {
                    scrollToImage(Array.from(images).indexOf(previousChapterImages[0])); // Scroll to the first image of the previous chapter
                }
            }
        } else {
            console.log("Already at the first chapter. No further scrolling.");
        }
    }

    /**
     * Update the index display with the current chapter and image information.
     * @param {number} index - The index of the current image.
     */
    function updateIndexDisplay(index) {
        if (indexDisplay) {
            let currentChapter = getCurrentChapterIndex(index);
            let chapterName = chapters[currentChapter]?.getAttribute('data-chapter-name') || 'Unknown';
            indexDisplay.textContent = `Chapter: ${chapterName}, Image ${index + 1} of ${totalImages}`;
        }
    }

    /**
     * Handle scroll events to update the current index based on the viewport position.
     */
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

    /**
     * Handle touch gestures for swiping and double-tap actions.
     */
    function handleGesture() {
        const horizontalSwipeDistance = Math.abs(touchEndX - touchStartX);
        const verticalSwipeDistance = Math.abs(touchEndY - touchStartY);

        if (horizontalSwipeDistance > minSwipeDistance && verticalSwipeDistance < minSwipeDistance) {
            // Horizontal swipe
            if (touchEndX < touchStartX) {
                scrollToImage(currentIndex + 1); // Swipe left to go to the next image
            } else if (touchEndX > touchStartX) {
                scrollToImage(currentIndex - 1); // Swipe right to go to the previous image
            }
        } else if (verticalSwipeDistance > minVerticalSwipeDistance && horizontalSwipeDistance < minSwipeDistance) {
            // Vertical swipe
            let currentChapterIndex = getCurrentChapterIndex(currentIndex);
            if (touchEndY < touchStartY && currentChapterIndex < chapters.length - 1) {
                scrollToNextChapter(); // Swipe up to go to the next chapter
            } else if (touchEndY > touchStartY && currentChapterIndex > 0) {
                scrollToPreviousChapter(); // Swipe down to go to the previous chapter
            }
        }

        // Reset touch coordinates
        touchStartX = touchEndX = touchStartY = touchEndY = 0;
    }

    /**
     * Handle double-tap gestures to go to the next image.
     */
    function handleDoubleTap() {
        scrollToImage(currentIndex + 1);
    }

    // Event listeners
    document.addEventListener('scroll', onScroll);

    document.addEventListener('keydown', function(e) {
        if (e.key === 'ArrowDown' || e.key === ' ') {
            e.preventDefault();
            scrollToImage(currentIndex + 1);
        } else if (e.key === 'ArrowUp') {
            e.preventDefault();
            scrollToImage(currentIndex - 1);
        } else if (e.key === 'PageDown' || e.key === 'ArrowRight') {
            e.preventDefault();
            scrollToNextChapter();
        } else if (e.key === 'PageUp' || e.key === 'ArrowLeft') {
            e.preventDefault();
            scrollToPreviousChapter();
        }
    });

    document.addEventListener('touchstart', function(e) {
        touchStartX = e.changedTouches[0].screenX;
        touchStartY = e.changedTouches[0].screenY;

        const currentTime = new Date().getTime();
        if (currentTime - lastTapTime < doubleTapThreshold) {
            handleDoubleTap();
        }
        lastTapTime = currentTime;
    });

    document.addEventListener('touchend', function(e) {
        touchEndX = e.changedTouches[0].screenX;
        touchEndY = e.changedTouches[0].screenY;
        handleGesture();
    });

    // Initialize the display
    updateIndexDisplay(currentIndex);
});
