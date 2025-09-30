# City Autocomplete API

A high-performance Go-based API for city name autocompletion with over 32,000 cities from around the world.

## Features
- Fast autocomplete search for city names
- Prioritizes prefix matches over substring matches
- Configurable result limits
- CORS-enabled for web usage
- Over 32,000 cities worldwide

## Endpoints

### GET /autocomplete
Search for cities by name with autocomplete functionality.

#### Parameters
- `q` (required): Search query string
- `limit` (optional): Maximum number of results (default: 10, max: 100)

#### Example Request
```
GET /autocomplete?q=lon&limit=3
```

#### Example Response
```json
[
  {
    "name": "London",
    "country": "Canada", 
    "subcountry": "Ontario",
    "geonameid": "6058560"
  },
  {
    "name": "Londonderry County Borough",
    "country": "United Kingdom",
    "subcountry": "Northern Ireland", 
    "geonameid": "2643734"
  }
]
```

## Running the API

1. Make sure you have Go 1.21+ installed
2. Ensure the `world-cities.csv` file is in the root directory
3. Run the application:
   ```bash
   go run main.go
   ```
4. The API will start on `http://localhost:8080`

## Configuration
- To change the port, set the `PORT` environment variable:
  ```bash
  PORT=3000 go run main.go
  ```