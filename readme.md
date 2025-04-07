# ITV Test Project - Movie Management API

This is a RESTful API built in Go for managing movies, featuring user authentication, rate limiting, and Redis caching. It uses [Gin](https://github.com/gin-gonic/gin) for routing, [GORM](https://gorm.io) for database interactions, [Redis](https://redis.io) for caching, and [Fx](https://github.com/uber-go/fx) for dependency injection.

## Features

- **Movie Management**: Create, read, update, and delete (CRUD) movies.
- **Authentication**: User registration, login, and token refresh with JWT-based access and refresh tokens.
- **Caching**: Redis caching for movie data with TTL expiration.
- **Rate Limiting**: Token bucket-based rate limiting using Redis.
- **Swagger Documentation**: Interactive API docs available at `/swagger/index.html`.
- **Health Check**: Simple health endpoint at `/health`.

## Tech Stack

- **Language**: Go 1.x
- **Web Framework**: Gin
- **Database**: PostgreSQL (via GORM)
- **Caching**: Redis
- **Dependency Injection**: Uber Fx
- **Authentication**: JWT
- **Rate Limiting**: Custom token bucket implementation
- **API Documentation**: Swagger (via gin-swagger)

## Prerequisites

- **Go**: 1.18 or later
- **PostgreSQL**: Running instance
- **Redis**: Running instance
- **Docker** (optional): For running services in containers

## Setup Instructions

### 1. Clone the Repository

```bash
git clone https://github.com/ruziba3vich/itv_test_project.git
cd itv_test_project
```

## Run the Application
```bash
docker compose up -d --build
```


## API Endpoints
Authentication Routes (/api/v1)
Method	Endpoint	Description	Request Body	Response Body	Authentication
POST	/register	Register a new user	CreateUserRequest	Success: 200, Error: 409/500	None
POST	/login	Login and get tokens	LoginUserRequest	LoginUserResponse	None
POST	/refresh	Refresh access token	RefreshTokenReq	RefreshTokenResponse	None
Movie Routes (/api/v1)
Method	Endpoint	Description	Request Body/Params	Response Body	Authentication
POST	/movies	Create a new movie	CreateMovieRequest	CreateMovieResponse	Required
GET	/movies	Get all movies (paginated)	Query: limit, offset	GetAllResponse	None
GET	/movies/:id	Get a movie by ID	URI: id	GetByIDResponse or null	None
PUT	/movies/:id	Update a movie by ID	URI: id, UpdateMovieRequest	UpdateMovieResponse	Required
DELETE	/movies/:id	Delete a movie by ID	URI: id	DeleteMovieResponse	Required
Utility Routes
Method	Endpoint	Description	Response Body
GET	/health	Check server health	200 OK
GET	/swagger/*any	Swagger UI	Swagger HTML
Authentication

    Access Token: Include in the Authorization header as Bearer <access_token> for protected routes.
    Refresh Token: Send in the refresh_token field of the /refresh request body.

Example:
```bash
curl -X POST http://localhost:7777/api/v1/movies \
  -H "Authorization: Bearer <access_token>" \
  -d '{"title":"Inception","director":"Nolan","year":2010,"plot":"A dream heist"}'
Request/Response Types
Authentication

    CreateUserRequest: { "full_name": string, "username": string, "password": string }
    LoginUserRequest: { "username": string, "password": string }```
