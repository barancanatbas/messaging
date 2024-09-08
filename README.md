# Go Messaging Service

This project is a Go application that manages messaging services. The application can be run using Docker Compose and integrates with services like Redis and MySQL, while automating various tasks using Taskfile.

## Table of Contents
- [Installation](#installation)
- [Running](#running)
- [API Documentation](#api-documentation)

---

## Installation

Follow the steps below to run the project:

### 1. Required Tools:
- **Go** (v1.22+)
- **Docker** and **Docker Compose**
- **Task CLI**: We use the `task` command-line tool for task management.

  **Task CLI Installation (MacOS and Linux)**: You can download and install Task CLI from this link: https://taskfile.dev/installation/

---

## Running

To run the services, you can use either of the following commands:

```bash
task compose-up
```
Or

```bash
docker-compose up -d --build
```

## API Documentation
The API documentation in the project is automatically generated using Swagger. To access the API documentation:

You can update the Swagger documentation with the following task command:

```bash
task swagger
```

```bash
http://localhost:8080/swagger/index.html
```

This address allows you to explore and test your APIs via Swagger UI.

---# messaging
