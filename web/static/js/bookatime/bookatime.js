// Get references to the elements
const dateInput = document.getElementById('date');
const timeInput = document.getElementById('time');
const fileInput = document.getElementById('file');
const bookButton = document.getElementById('bookButton');

// Handle the "Book" button click
bookButton.addEventListener('click', () => {
  // Get the values from the form
  const date = dateInput.value;
  const time = timeInput.value;
  const file = fileInput.files[0]; // Assuming only one file is allowed

  // Validate the inputs (optional)
  if (!date || !time || !file) {
    alert('Please fill in all fields.');
    return;
  }

  // Create a FormData object to send the data
  const formData = new FormData();
  formData.append('date', date);
  formData.append('time', time);
  formData.append('file', file);

  // TODO: fix no such file error
  // Send the data to the server using fetch
  fetch('/book', {
    method: 'POST',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    body: formData
  })
  .then(response => {
    if (!response.ok) {
      throw new Error('Network response was not ok');
    }
    return response.json();
  })
  .then(data => {
    console.log('Booking successful:', data);
    // Handle successful booking, e.g., display a confirmation message
  })
  .catch(error => {
    console.error('Error booking:', error);
    // Handle booking error, e.g., display an error message
  });
});
