# Filesharing

A filesharing application built with [gin](https://github.com/gin-gonic/gin).

## Table of contents

- [About](#about)
- [Usage](#usage)
- [License](#license)

## About

### Endpoints

- `GET /api/download/{id}` - Download file
- `GET /api/file/{id}` - File info <br>
Example response:
```json
{
  "Status": "OK",
  "FileName": "Name",
  "FileSize": "Size in Mb",
  "ExpirationDate": "Date",
  "DownloadLink": "/api/download/{id}"
}
```
- `POST /api/upload` - Upload file <br>
Example request:
```json
{
  "File": "file",
  "Expiration": "Options in .env"
}
```

## Usage

### Configuration

1. Copy `config/.env.example` contents to `config/.env`, change some parameters if necessary.
2. In `/config/nginx.conf` if necessary, change:
   1. `server_name` - according [this instructions](https://nginx.org/en/docs/http/server_names.html),
   2. `client_max_body_size` - to maximum desired file size,
   3. Add HTTPS configuration following [this guide](https://nginx.org/en/docs/http/configuring_https_servers.html).

### Deploying

Run `docker compose up -d` on your local machine or your server.

## License

Licensed under the MIT License. See [LICENSE.txt](LICENSE.txt) for details.
