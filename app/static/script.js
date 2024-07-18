  const watchButton = document.getElementById('watchButton' );
        const bookTimeButton = document.getElementById('bookTimeButton');

        watchButton.addEventListener('click', () => {
            window.location.href = 'watch.html';
        });

        bookTimeButton.addEventListener('click', () => {
            window.location.href = 'bookatime.html';
        });