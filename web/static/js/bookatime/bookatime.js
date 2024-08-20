


        // Get references to the elements
        const durationSlider = document.getElementById('duration');
        const durationValue = document.getElementById('duration-value');
        const bookButton = document.getElementById('bookButton');

        // Handle the "Book" button click

  document.getElementById('bookButton').addEventListener('click', () => {
        const date = document.getElementById('date').value;
        const time = document.getElementById('time').value;
        const file = document.getElementById('file').files[0];

        const formData = new FormData();
        formData.append('date', date);
        formData.append('time', time);
        formData.append('file', file);

        bookButton.addEventListener('click', () => {
            // Get the values from the form
            const date = document.getElementById('date').value;
            const time = document.getElementById('time').value;
            const file = document.getElementById('file').files[0];

            // ... (Your logic to process the booking data) ...
        });
        fetch('/book', {
          method: 'POST',
  headers: {
    'Content-Type': 'application/json',
              //TODO: wtf i should write here if my token contains in the cookie
    'Authorization': `Bearer `
  },
          body: formData
        })
        .then(response => {
          if (!response.ok) {
            throw new Error('Network response was not ok');
          }
          return response.json();  
 // Or handle response as needed
        })
        .then(data => {
          console.log('Booking successful:', data);
          // Handle successful booking, e.g., display confirmation message
        })
        .catch(error => {
          console.error('Error booking:', error);
          // Handle booking error, e.g., display error message
        });
      });
