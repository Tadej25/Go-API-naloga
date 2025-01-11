
# Go API naloga

A basic API for displaying temperature data for cities.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

- Go (version 1.16 or higher)
- Git

### Installation

1. Clone the repository:

```bash
git clone https://github.com/your-username/Go-API-naloga.git
```

2. Change into the project directory:

```bash
cd Go-API-naloga
```

3. Install the required dependencies:

```bash
go get -u github.com/gin-gonic/gin
go get -u github.com/ahmetb/go-linq/v3
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
```

4. Run the application:

```bash
go run main.go
```

The API will start running on `http://localhost:8080`.

### API Documentation

The API documentation is generated using Swagger. To access the Swagger UI, navigate to `http://localhost:8080/swagger/index.html`.

## Built With

- [Go](https://golang.org/) - The programming language used
- [Gin](https://github.com/gin-gonic/gin) - A web framework for Go
- [go-linq](https://github.com/ahmetb/go-linq) - A LINQ-like library for Go
- [gin-swagger](https://github.com/swaggo/gin-swagger) - A Swagger 2.0 documentation generator for Gin
- [swaggo/files](https://github.com/swaggo/files) - A Swagger 2.0 file serving package for Gin

## Authors

- **Tadej Lipar** - *Initial work* - [Tadej25](https://github.com/Tadej25)
