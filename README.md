# WASAphoto — Full-Stack Photo Sharing Prototype

A photo-sharing web application built with **Go**, **Vue 3**, **SQLite**, and **Docker**. It implements a REST API and a browser-based client for publishing photos and managing social interactions.

## Overview

WASAphoto separates application logic, data persistence, and presentation into distinct components. The Go service exposes HTTP endpoints backed by SQLite, while a Vue 3 single-page application communicates with the API through Axios. An OpenAPI 3 specification documents the HTTP interface. Not a production-hosted service.

## Features

- **User profiles:** username-based session creation, profile retrieval, profile search, and nickname updates.
- **Photo management:** upload and delete photos; display a user's published photos.
- **Social interactions:** like and unlike photos, add and remove comments, follow and unfollow users.
- **Feed:** retrieve a photo stream based on followed users.
- **User controls:** ban and unban users.
- **Web interface:** views for login, feed, profile, search, and settings.
- **API contract:** OpenAPI 3.0.3 specification in [`doc/api.yaml`](doc/api.yaml).

## Architecture

```text
                 Vue 3 SPA (Vite + Axios)
                            |
                       HTTP / JSON
                            |
                      Go REST API
                      /         \
               SQLite DB      Photo files
```

The Go application keeps HTTP handlers in `service/api/` and database operations in `service/database/`. Photo metadata and social relationships are stored in SQLite; uploaded image files are stored on the server filesystem. The frontend is organized into views, reusable components, routing, and API services.

## Tech Stack

| Layer | Technology |
| --- | --- |
| Backend | Go, `net/http`, `httprouter` |
| Frontend | Vue 3, JavaScript, Vue Router, Axios, Vite |
| Database | SQLite |
| API documentation | OpenAPI 3.0.3 (YAML) |
| Containerization | Docker, Docker Compose, Nginx |

## Getting Started

### Prerequisites

- Docker and Docker Compose

### Run with Docker Compose

From the repository root:

```bash
# Build and start the frontend and backend containers
docker compose up --build
```

The Compose configuration publishes:

| Service | Local address |
| --- | --- |
| Frontend | http://localhost:8080 |
| Backend API | http://localhost:3000 |

Stop the services with `Ctrl+C`, or use `docker compose down` if started in detached mode.

> **Prototype data:** The supplied Compose file does not configure persistent volumes. The backend's default SQLite database and uploaded images reside in the container filesystem; data may be lost when a container is recreated.

### Run the frontend for development

With Node.js and npm installed:

```bash
cd webui
npm ci
npm run dev
```

The Vite client is configured to call the API at `http://localhost:3000`, so the backend must also be running. The frontend scripts are defined in [`webui/package.json`](webui/package.json).

### Run the backend locally

With Go installed, from the repository root:

```bash
go run ./cmd/webapi
```

By default, the service listens on port `3000` and uses a SQLite database under `/tmp`. See `cmd/webapi/load-configuration.go` for configuration options.

## Repository Structure

```text
.
├── cmd/
│   ├── webapi/             # HTTP server entry point and configuration
│   └── healthcheck/        # Health-check utility
├── service/
│   ├── api/                # Request handlers and application logic
│   └── database/           # SQLite persistence layer
├── webui/
│   └── src/                # Vue components, views, router, API client
├── doc/
│   └── api.yaml            # OpenAPI specification
├── demo/                   # Sample service configuration
├── Dockerfile.backend
├── Dockerfile.frontend
└── docker-compose.yml
```

## Scope and Limitations

This project is a full-stack application prototype focused on RESTful API design, backend development, data persistence, and containerized deployment.

The current implementation uses SQLite for data persistence, local filesystem storage for uploaded images, and Docker Compose for local execution.

Authentication is simplified for demonstration purposes and is not suitable for deployment as a public-facing service without additional security measures.

## License

This repository includes an MIT license; consult [`LICENSE`](LICENSE) for its terms and existing attribution.
