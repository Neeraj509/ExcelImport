# Golang Excel Import and CRUD System

## Installation

1. Clone the repository
2. Run `go mod tidy` to install the necessary packages.
3. Set up MySQL and Redis locally.
4. Update the MySQL credentials in `.env` file.
5. Run `go run main.go`.
6. Check the logs in `app.log` file.

## Endpoints

1. **POST /import**: Upload Excel file for processing.
2. **GET /records**: View all records (from Redis or MySQL).
3. **PUT /record/:id**: Edit a record.
4. **/DELETE /records/:id**: Delete a single record.
5.  **/DELETE /records/delete-all**: Delete all the record.
