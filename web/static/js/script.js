const host = "http://localhost:8080/";
const watchButton = document.getElementById('watchButton' );
        const bookTimeButton = document.getElementById('bookTimeButton');

        watchButton.addEventListener('click', () => {
            window.location.href = host + 'watch';
        });

        bookTimeButton.addEventListener('click', () => {
            window.location.href = host + 'book';
        });
