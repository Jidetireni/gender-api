# Gender API - HNG Stage 0 Assessment

This is a lightweight RESTful API built in Go that acts as an intelligent wrapper around the [Genderize API](https://genderize.io/). It accepts a name, fetches its gender prediction, and applies additional data processing and confidence evaluation before returning the results.

## Features

- **API Integration**: Fetches real-time gender data from the Genderize API.
- **Confidence Scoring**: Calculates a confidence boolean (`is_confident`) based on strict probability (>= 0.7) and sample size (>= 100) thresholds.
- **Data Transformation**: Formats the external API response, renaming `count` to `sample_size` and injecting an ISO 8601 UTC `processed_at` timestamp.
- **Validation & Error Handling**: 
  - Validates missing or empty names (400 Bad Request).
  - Rejects numbers in names (422 Unprocessable Entity).
  - Handles edge cases where no prediction is available or gender is null (404 Not Found).
- **CORS Support**: Fully accessible across domains (`Access-Control-Allow-Origin: *`).
- **Containerized**: Includes a multi-stage Dockerfile for lightweight and efficient deployment.

## Prerequisites

- [Go](https://golang.org/dl/) 1.22 or higher (if running locally)
- [Docker](https://docs.docker.com/get-docker/) (optional, for containerized deployment)

## Environment Variables

Create a `.env` file in the root directory (or export these variables in your environment) with the following structure:

```env
HOST=0.0.0.0
PORT=8000
ENV=development
GENDERIZED_API_BASE_URL=https://api.genderize.io
```

## Running the Application

### Option 1: Running Locally with Go

1. Clone the repository and navigate into the project directory:
   ```bash
   git clone <repository-url>
   cd gender-api
   ```
2. Download dependencies:
   ```bash
   go mod download
   ```
3. Run the server:
   ```bash
   go run ./cmd
   ```

### Option 2: Running with Docker

1. Build the Docker image:
   ```bash
   docker build -t gender-api .
   ```
2. Run the container:
   ```bash
   docker run -p 8000:8000 --env-file .env gender-api
   ```

The API will be available at `http://localhost:8000`.

## API Documentation

### Classify Name

**Endpoint:** `GET /api/classify`

**Query Parameters:**
- `name` (string, required): The name to classify. Must only contain alphabetical characters.

#### Example Request

```http
GET /api/classify?name=peter
```

#### Example Success Response (200 OK)

```json
{
  "status": "success",
  "data": {
    "name": "peter",
    "gender": "male",
    "probability": 0.99,
    "sample_size": 165452,
    "is_confident": true,
    "processed_at": "2024-04-01T12:00:00Z"
  }
}
```

#### Example Error Responses

**Missing Name (400 Bad Request)**
```json
{
  "status": "error",
  "message": "name query parameter is required"
}
```

**Invalid Name containing numbers (422 Unprocessable Entity)**
```json
{
  "status": "error",
  "message": "name must not contain numbers"
}
```

**No Prediction Available (404 Not Found)**
```json
{
  "status": "error",
  "message": "No prediction available for the provided name"
}
```
