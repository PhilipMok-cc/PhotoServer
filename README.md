# PhotoServer

PhotoServer is a simple web server for serving photo albums. It allows you to browse and view photos stored in a specified directory. Utilizing the features of HTTP in Go 1.22, it is the simplest way to serve your photos on NAS.

## Features

- Serve photos from a specified directory
- Display photo albums with thumbnails
- View individual photos with navigation
- The Docker version is designed to run on NAS

## Requirements

- Go 1.22 or later
- Docker (optional, for containerized deployment)

## Installation

### Using Go

1. Clone the repository:
    ```sh
    git clone https://github.com/yourusername/PhotoServer.git
    cd PhotoServer
    ```

2. Build the application:
    ```sh
    go build -o photo-server
    ```

3. Create a `config.yaml` file with the following content:
    ```yaml
    photoDir: /path/to/your/photos
    port: 8080
    ```

4. Run the application:
    ```sh
    ./photo-server
    ```

### Using Docker

1. Build the Docker image:
    ```sh
    ./build_and_save.sh
    ```

2. Run the Docker container:
    ```sh
    docker run -d -p 8080:8080 -v /path/to/your/photos:/photos photo-server
    ```

## Configuration

The application is configured using a `config.yaml` file. The following options are available:

- `photoDir`: The directory where your photos are stored.
- `port`: The port on which the server will listen.

Example `config.yaml`:
```yaml
photoDir: /path/to/your/photos
port: 8080
```

## Usage

1. Open your web browser and navigate to `http://localhost:8080`.
2. Browse the photo albums and click on an album to view the photos.

## Development

### Running Tests

To run the tests, use the following command:
```sh
go test ./...
```

### Directory Structure

- `PhotoServer.go`: The main application code.
- `PhotoServer_test.go`: The test code.
- `Dockerfile`: The Dockerfile for building the Docker image.
- `build_and_save.sh`: The script for building and saving the Docker image.
- `config.yaml`: The configuration file.
- `static/`: The directory for static assets (CSS, JS).
- `templates/`: The directory for HTML templates.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.