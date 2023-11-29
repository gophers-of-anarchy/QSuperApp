# QSuperApp Main Application

This Go application is the main entry point for the QSuperApp. It sets up an Echo web server with various routes for handling authentication, user accounts, airplanes, orders, and payments.

## Dependencies

- [Echo](https://github.com/labstack/echo/v4): A high-performance, minimalist Go web framework.
- [Joho/godotenv](https://github.com/joho/godotenv): A Go (golang) port of the Ruby dotenv project (Loads environment variables from `.env` files).
- [PostgreSQL](https://www.postgresql.org/): A powerful, open-source relational database system.

## Application Structure

The application is structured as follows:

- **Main Package**: `main`
  - Imports necessary packages and initializes the application.
  - Connects to the PostgreSQL database using configurations from the `.env` file.

- **Template Struct**: `Template`
  - Implements a custom template renderer for rendering HTML templates using the `html/template` package.

- **Main Function**: `main()`
  - Checks the application environment and loads environment variables from the `.env` file if not set.
  - Initializes an Echo instance and sets up middleware for logging, recovery, and CORS.
  - Connects to the PostgreSQL database.
  - Configures the application to render HTML templates.
  - Defines API routes for authentication, account management, airplanes, orders, and payments.
  - Runs the Echo server on port `8080`.

## API Routes

### Auth Routes: `/api/v1/auth`

- `POST /register`: Handles user registration.

    **Request:**
    ```json
    {
        "username": "example_user",
        "password": "example_password",
        "email": "example@gmail.com",
        "cellphone": "09999999999"
    }

- `POST /login`: Handles user login.

    **Request:**
    ```json
    {
        "username": "example_user",
        "password": "example_password"
    }

### Account Management Routes: `/api/v1/users`

- `POST /register`: Creates a new user account (requires authentication)

    **Request**
    ```json
    {
        "name": "example"
    }

- `PUT /profile`: Updates user profile (requires authentication)
    
    **Request**
    ```json
    {
        "username": "example",
        "email": "updated_email@example.com",
        "cellphone": "0999999999"
    }

- `GET /profile/:id`: Retrieves user profile by ID (requires authentication and admin privileges).

    **Request**
    ```json
    {
        "username": "example",
        "email": "example@gmail.com",
        "cellphone": "099999999",
        "created_at": "example",
        "updated_at": "example",
        "is_admin": false
    }

### Airplane Routes: `/api/v1/airplane`

- `POST /add/`:  Adds a new airplane.
    
    **Request**
    ```json
    {
        "type": "example",
        "base_price": 1.5,
        "number": "example"
    }

- `PUT /update/`: Updates an airplane.

    **Request**
    ```json
    {
        "type": "example",
        "base_price": 1.5,
        "number": "example",
        "ID": 1
    }

- `GET /all/`: Get all airplanes.

    **Response**
    ```json
    [
        {
            "id": 1,
            "type": "example",
            "base_price": 1.5,
            "number": "example",
            "created_at": "2000/1/1",
            "updated_at": "2000/1/1"
        }
    ]
    
- `GET /:id/`: Get an airplane detail
    **Response**
    ```json
    
    {
        "id": 1,
        "type": "example",
        "base_price": 1.5,
        "number": "example",
        "created_at": "2000/1/1",
        "updated_at": "2000/1/1"
    }

-  `DELETE /:id/`: Delete an airplane

### Order Routes: `/api/v1/order-management`

- `POST /admin/orders/?order_id=<order_id>&status=<status>`: Change order status by admin
  
    **Response**
    ```json
    {
        "message": "Order status changed successfully",
        "OrderStatus": "Approved"
    }


- `POST /admin/orders/status/?order_id=<order_id>&status=<status>`: 


- `GET /admin/orders/list/`: Get all orders
    **Response**
    ```json
    [
        {
            "id": "example",
            "user_id": 1,
            "airplane_id": 1,
            "status": "examole",
            "number": "example",
            "created_at": "2000/1/1",
            "updated_at": "2000/1/1",
        }
    ]
    
### Advance Payment `/api/v1/payment`

- `POST /advance`: advance payment
    **Request**
    ```json
    {
        "order_id": 1
    }

- `POST /finalize`: finalize payment
    **Request**
    ```json
    {
        "order_id": 1 
    }

- `GET /orders/:order_id`: Get order payment status
    **Response**
    ```json
    [
        {
            "ID": 3,
            "CreatedAt": "2023-11-29T23:13:33.648297+03:30",
            "UpdatedAt": "2023-11-29T23:13:33.814969+03:30",
            "Amount": 7600000,
            "PaymentType": 0,
            "PaymentStatus": 2,
            "OrderID": 1,
            "PaymentTypeString": "Advance",
            "PaymentStatusString": "Failed"
        },
        {
            "ID": 6,
            "CreatedAt": "2023-11-29T23:18:39.834706+03:30",
            "UpdatedAt": "2023-11-29T23:18:39.835082+03:30",
            "Amount": 7600000,
            "PaymentType": 1,
            "PaymentStatus": 0,
            "OrderID": 1,
            "PaymentTypeString": "Final",
            "PaymentStatusString": "Completed"
        }
    ]

### Verify Payment `/api/v1/verify`

- `POST /page`: 

- `POST /payment`


## How to Run

To run the application, execute the following command:

```bash
go run main.go
