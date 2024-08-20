# Web and Software Architecture
Welcome to the repository for the Web and Software Architecture course. This README outlines the project work across four homework assignments, focusing on API design, backend development, frontend creation, and deployment using Docker.
Table of Contents
* Homework 1: API Definition
* Homework 2: Backend Development in Go
* Homework 3: Frontend Development in Vue.js
* Homework 4: Deployment Using Docker

## API Definition
This assignment involved creating a comprehensive API specification using the OpenAPI standard. The YAML document detailing the API is thoroughly structured to cover all endpoints, responses, and security definitions required for the WASAphoto service.
* Document Path: doc/api.yaml
  ## Key Features:
* Defined RESTful endpoints for user authentication, profile management, and social interactions.
* Detailed request and response models with appropriate status codes to enhance client-server communication.
* Incorporation of security schemes to ensure that API interactions are secure.

## Backend Development in Go
Using the Fantastic Coffee Template, the backend was developed in Go, focusing on robustness and secure data handling.
Highlights:
* Functions and Endpoints: Implementation of all defined API endpoints with full support for CRUD operations.
* Data Structures: Designed to efficiently handle and manipulate data throughout the application lifecycle.
* Database Interactions: Secure interactions with an SQLite database, ensuring that data integrity and access control are maintained.

## Frontend Development in Vue.js
The frontend of our application harnesses the power of Vue.js to create a dynamic and responsive user interface.
Features:
* Component Design: Emphasis on reusable and modular components.
* Data Management: Efficient handling of state and props to ensure smooth data flow across components.
* User Experience: Focus on creating an intuitive and user-friendly interface.

## Deployment Using Docker
Finalizing the project, the application was containerized using Docker, focusing on lightweight containers and efficient deployment.
Process:
* Dockerfile Creation: Crafted to build a minimal and optimized image.
* Container Management: Configured Docker to handle the application in a lightweight and scalable manner.
* Performance Optimization: Ensured that the container runs efficiently with minimal resource usage.
