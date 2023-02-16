# usrsvc
[![.github/workflows/ci.yaml](https://github.com/alesr/usrsvc/actions/workflows/ci.yaml/badge.svg?branch=master)](https://github.com/alesr/usrsvc/actions/workflows/ci.yaml)

The User Service is a microservice designed to provide basic CRUD operations for users, using gRPC for communication, PostgreSQL for data storage, Docker for containerization, and Go for the implementation. It has been implemented following Layered Architecture and Clean Architecture principles, with each layer having its own responsibilities to ensure maintainability and scalability.

### Note for the reviewer
A good place to start looking at the code is the end-to-end tests in the tests directory. The tests are commented to explain some of the design decisions and the rationale behind the implementation. I also expanded on some comments in the code itself to compensate for the lack of discussions during the implementation.

## How to Run

To run the application, Docker must be installed on your machine. If you don't have it installed, you can find instructions for installing Docker [here](https://docs.docker.com/get-docker/). Once Docker is installed, you can run the following command:

```bash
make run
```

This command will spin up a PostgreSQL container and a container for the application. The application will be available on `http://localhost:50051`. 


## How to Test

For code formatting, static analysis, unit and integration tests you can run the following command:
```bash
make test
```

For end-to-end tests you can run the following command:
```bash
make test-e2e
```

For more information about the available commands, you can run the following command:
```bash
make help
```
