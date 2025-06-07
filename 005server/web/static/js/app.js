document.addEventListener('DOMContentLoaded', function() {
    console.log('Go HTTP Server - Client-side loaded');
    
    // Add click handlers for API endpoint buttons
    const apiButtons = document.querySelectorAll('.endpoint .btn');
    
    apiButtons.forEach(button => {
        button.addEventListener('click', function(e) {
            // For JSON endpoints, open in new tab to show raw JSON
            if (this.href.includes('/api/')) {
                e.preventDefault();
                window.open(this.href, '_blank');
            }
        });
    });
    
    // Add a simple loading animation for buttons
    const buttons = document.querySelectorAll('.btn');
    buttons.forEach(button => {
        button.addEventListener('click', function() {
            this.style.opacity = '0.7';
            setTimeout(() => {
                this.style.opacity = '1';
            }, 200);
        });
    });
    
    // Add current time update functionality (if on home page)
    const timeElement = document.querySelector('footer p');
    if (timeElement && timeElement.textContent.includes('Server time:')) {
        updateTime();
        setInterval(updateTime, 1000);
    }
    
    function updateTime() {
        const now = new Date();
        const timeString = now.toISOString();
        const parts = timeElement.textContent.split(' | ');
        if (parts.length >= 2) {
            timeElement.textContent = `Client time: ${timeString} | ${parts[1]}`;
        }
    }
    

    console.log('%cGo HTTP Server', 'color: #667eea; font-size: 20px; font-weight: bold;');
    console.log('API Endpoints:');
    console.log('GET /api/info - Server information');
    console.log('GET /api/time - Current server time');
    console.log('GET /api/echo - Echo request data');
    console.log('GET /health - Health check');
});