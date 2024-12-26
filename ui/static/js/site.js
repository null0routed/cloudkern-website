window.onload = (event) => {
    console.log("Page Loaded");
    
    const video = document.getElementById('bgvid');
    if (!video) {
        console.error('Video element not found');
        return;
    }

    video.addEventListener('loadedmetadata', () => {
        video.playbackRate = 0.5;
        console.log('Video playback rate set to 0.5');
    });

    video.load();
};