# Auth API Service

This project is a simple authentication system implemented in Go. It allows users to sign up, sign in, access protected endpoints, refresh tokens, and revoke tokens. The API uses JSON Web Tokens (JWT) for authentication and authorization. The backend is built using the Go programming language and net/http.

## Features
- **Sign Up**: Register a new user with an email and password.
- **Sign In**: Authenticate a user and generate an access token (JWT).
- **Protected Endpoints**: Only accessible by authenticated users with a valid JWT.
- **Token Refresh**: Renew an access token using a valid refresh token.
- **Token Revocation**: Revoke a refresh token to invalidate it.
- **JWT Authentication**: Using JSON Web Tokens (JWT) for user authentication.



### Detailed Breakdown of Folders and Files:
- **/config**: Contains the configuration file for initializing the database connection.
- **/controllers**: Contains the logic for handling HTTP requests. It has functions like `SignUp`, `SignIn`, and `ProtectedEndpoint`.
- **/db**: Contains database-related logic, including migration functions and models for interacting with the database.
- **/middleware**: Middleware functions like `AuthMiddleware` to check the validity of JWT tokens in incoming requests.
- **/models**: Contains the data model for user authentication, specifically the `User` model and functions to interact with the database.
- **/services**: Contains the business logic for signing up and signing in users.
- **/utils**: Contains helper functions like `GenerateToken` to create JWT tokens and validate them.
- **/migrations**: Database migration files that define the structure of the database (e.g., user table creation).

## Dependencies

- Go 1.18+ (Ensure Go is installed on your machine)
- JWT-Go (for generating and validating JWT tokens)
- SQLite  `modernc.org/sqlite` for in-memory database (or use file-based storage)

You can install the required dependencies by running:

```bash
go mod tidy


## How to Run the Application

- 1. Clone the Repository
    git clone https://github.com/yourusername/auth-api.git
    cd auth-api

- 2. Build and Run the server
    go run main.go

-This will start the server at http://localhost:8080.

```



## API Endpoints

- Sign Up (Create a new user)
```bash
    curl -X POST http://localhost:8080/signup -H "Content-Type: application/json" -d '{"email": "testuser@example.com", "password": "password123"}'
```

- Sign In (Get authentication token)
```bash
    curl -X POST http://localhost:8080/signin -H "Content-Type: application/json" -d '{"email": "testuser@example.com", "password": "password123"}'
```
    Response:
    {
        "token": "token-here",
        "refresh_token": "refresh-token-here",
        "expires_in": 900
    }

- Access Protected Route (Requires authentication)
```bash
    curl -X GET http://localhost:8080/protected -H "Authorization: Bearer your-access-token-here"
```
    Response:
    {
        "message": "Welcome to the protected endpoint",
        "user": "testuser@example.com"
    }

- Refresh Token (Renew the access token)
```bash
    curl -X POST http://localhost:8080/refresh -H "Content-Type: application/json" -d '{"refresh_token": "token-here"}'
```


## Note
- The application currently uses an in-memory SQLite database or file storage (modernc.org/sqlite).
- The JWT tokens have a default expiration of 15 minutes (configurable).
- Token refresh and revocation mechanisms are implemented for session management.