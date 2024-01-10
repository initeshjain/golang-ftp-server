# Golang FTP Service

This project provides a simple FTP-like service built with Golang, using SQLite as the metadata storage. The service allows you to upload, get, and delete files from a local directory.

## Setup

### Install Golang:
   Make sure you have Golang installed on your machine. You can download it from [https://golang.org/dl/](https://golang.org/dl/).

### Clone the Repository:
```bash
git clone https://github.com/initeshjain/golang-ftp-service.git
cd golang-ftp-service
```

### Install Dependencies:
Run the following commands to install the required dependencies:
```bash
go get -u github.com/gin-gonic/gin
go get -u github.com/mattn/go-sqlite3
```

### Build the Application:
```bash
go build -o ftp-service main.go
```

### Run the Application:
```bash
./ftp-service
```

The server will start on http://localhost:8080.

## Configuration:
You can configure the service by modifying the config.json file. Update the UploadDir to specify the local directory path where uploaded files will be stored.

```json
{
"UploadDir": "/path/to/upload/directory",
"DbPath": "ftp_metadata.db"
}
```

## Usage
### Upload File
To upload a file, you can use curl from the command line. Replace `<file_path>` and `<your_server_url>` with the actual file path and server URL.
```bash
curl -X POST -F "file=@<file_path>/example.txt" http://localhost:8080/upload
```

### Get File
To retrieve a file, use the below endpoint in your browser run below command. Replace `<filename>` with the name of the file you want to retrieve.
```bash
curl -X GET http://localhost:8080/get/<filename> -o <filename>
```

### Delete File
To delete a file, use the following endpoint in your browser or any HTTP client. Replace `<filename>` with the name of the file you want to delete.
```bash
curl -X DELETE http://localhost:8080/delete/<filename>
```

## Notes
No authentication is required as it's intended for local development purposes.
This is a basic implementation, and you may need to enhance security and error handling based on your requirements.


# Connect with me
Connect with me on social media:

[![LinkedIn](https://img.shields.io/badge/LinkedIn-initeshjain-blue)](https://www.linkedin.com/in/initeshjain/)
[![Instagram](https://img.shields.io/badge/Instagram-initeshjain-orange)](https://www.instagram.com/initeshjain/)
[![Twitter](https://img.shields.io/badge/Twitter-initeshjain-black)](https://twitter.com/initeshjain)
