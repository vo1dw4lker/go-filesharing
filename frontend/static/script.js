const uploadButton = document.getElementById("upload-btn");
const inputFile = document.getElementById("file-input");
const fileLink = document.getElementById("file-link");
const fileLinkLabel = document.getElementById("file-link-label");
const expOption = document.getElementById("expiration-options");

const currentUrl = window.location.href;
const baseUrl = currentUrl.split('/').slice(0, 3).join('/');

uploadButton.addEventListener('click', () => {
    const file = inputFile.files[0];

    // Create a FormData object and append the file to it
    const formData = new FormData();
    formData.append("file", file);
    formData.append("exp", expOption.value);

    // Send the form data to the server using Fetch API
    fetch("/api/upload", {
        method: "POST",
        body: formData
    })
    .then(response => response.json())
    .then(data => {
        fileLink.value = `${baseUrl}/files/${data.link}`;
        fileLink.style.display = 'block';
        fileLinkLabel.style.display = 'block';
    })
    .catch(error => {
        console.error('Error uploading file:', error);
    });
});
