# Bin/IIN Lookup service

![Example image](https://raw.githubusercontent.com/aaronburt/repo-image-host/main/chrome_f60APRpqEe.png "Example image")


The BIN/IIN Lookup Service is a simple web application that allows users to search for information about Bank Identification Numbers (BINs). A BIN is the first six digits of a payment card number and is used to identify the issuing institution. This service reads from a [CSV file](https://github.com/venelinkochev/bin-list-data) containing BIN data and returns detailed information in JSON format.


* Server.go is the Go file that handles HTTP requests, CSV file reading, and JSON response generation.

* A .env file with the following environment variable 
    - PORT: The port on which the server should run, must be prefixed with a colon (:)

```env
PORT=:80
```

### Environment
Built on GO Version 1.23.0 (Windows)

### Acknowledgments

The data used in this service is provided by the [bin-list-data](https://github.com/venelinkochev/bin-list-data) repository.