# 360-Degree Feedback System

This is a simple Go server that implements a basic 360-degree feedback system. It uses the Go `html/template` package for templating and SQLite for storing feedback answers.

## Requirements

- Go (>= 1.13)
- SQLite (go-sqlite3 package)

## Installation

1. Clone the repository to your local machine.

   ```bash
   git clone https://github.com/yourusername/360-degree-feedback-system.git
   ```

2. Navigate to the project directory.

   ```bash
   cd 360-degree-feedback-system
   ```

3. Install the SQLite Go package.

   ```bash
   go get -u github.com/mattn/go-sqlite3
   ```

## Usage

1. Run the Go server:

   ```bash
   go run main.go
   ```

2. Open your web browser and navigate to `http://localhost:8080/` to fill out the feedback form.

3. Submit the form, and you'll be redirected to the results page at `http://localhost:8080/results`.

## Structure

- `main.go`: The main Go file that sets up the server, routes, and SQLite database.
- `form.html`: HTML template for the feedback form.
- `results.html`: HTML template for displaying the stored feedback.

## Features

- Basic 360-degree feedback form.
- SQLite database for storing feedback answers.
- Results page for viewing all stored feedback.

## Future Improvements

- Adding data validation
- Improved error handling
- More comprehensive questions
- Authentication for admin route to view results
