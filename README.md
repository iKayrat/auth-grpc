# Authentication gRPC Service

This project is a gRPC service written in Go that provides authentication functionality. It includes features for user registration, login, and token generation.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
  - [Running the Service](#running-the-service)
  - [gRPC Endpoints](#grpc-endpoints)
- [Configuration](#configuration)
- [Authentication](#authentication)
- [Examples](#examples)


## Installation

Make sure you have Go installed on your machine. Clone the repository and install dependencies:

```bash
git clone https://github.com/iKayrat/auth-grpc.git
cd cmd/app/main.go
```

## Usage
### Running the Service

To run the authentication gRPC service, use the following command:

```bash
go run main.go
```

By default, the service will be available at localhost:8081.
### gRPC Endpoints

The service exposes the following gRPC endpoints:

    SignUp: Allows users to create a new account.
    Login: Authenticates a user and generates an authentication token.
    CreateUser: Allows admins to create a new user.
    GetUsers: Allows users to get a list of users.
    GetUserById: Allows admins to get a user info.
    UpdateUser: Allows admins to update a user.
    DeleteUser: Allows admins to delete a user.

## Configuration
You can configure the service using environment variables or a configuration file. 
Create a .env file with the following content:
```bash
# .env
PORT=8081

SERVER_ADDRESS=0.0.0.0:8081

DB_CONNECTION_STRING=mongodb://localhost:27017/authdb
JWT_SECRET=your-secret-key
```

## Authentication
The service uses JSON Web Tokens (JWT) for authentication. 
Include the JWT token in the Authorization header of your gRPC requests.

Example:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
```
