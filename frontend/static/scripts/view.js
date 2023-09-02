const fileId = getFileId()

// Fetch file information from the API
fetch(`/api/file/${fileId}`)
    .then(response => response.json())
    .then(data => {
         if (data["Status"] === "OK") {
            document.getElementById("fileName").textContent = data["FileName"];
             document.getElementById("fileSize").textContent = data["FileSize"];
             document.getElementById("expirationDate").textContent = data["ExpirationDate"];
             document.getElementById("downloadLink").href = data["DownloadLink"];
             document.getElementById("downloadLink").download = data["FileName"];
         } else if (data["Status"] === "Not found") {
            handleNotFound()
            console.error("File not found");
        } else {
            // Handle other errors
            console.error("Error:", data["Status"]);
        }
    })
    .catch(error => {
        console.error("Error fetching file information:", error);
    });

function getFileId() {
  const urlPath = window.location.pathname;
  const pathParts = urlPath.split("/");

  // The last part of the path should be the file ID
  return pathParts[pathParts.length - 1];
}

function handleNotFound() {
    // Replace the current page content with the content of the 404.html template
    fetch("/static/html/404.html")
        .then(response => response.text())
        .then(html => {
            document.body.innerHTML = html;
        })
        .catch(error => {
            // Handle any errors that occur while fetching the 404.html template
            console.error('Error fetching 404.html:', error);
        });
}