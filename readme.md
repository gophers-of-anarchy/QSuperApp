# Q Super App

## Component #1 Tasks

### Database Design and Setup (Narges)
- **Design the database schema for airplanes, users, orders, customizations, and payments:** Map out a detailed schema that captures all necessary data and relationships for the application's core functionalities.
- **Document the database design:** Create thorough documentation that covers all aspects of the database schema, which can be used for future reference and by new team members.
- **Set up the database using an ORM:** Configure the selected ORM to connect to the database and set up the initial models based on the designed schema.
- **Implement database migrations:** Write migration scripts to create and update the database schema without downtime or data loss.

### Authentication and User Management (Parsa)
- **Implement JWT-based authentication:** Code a secure authentication mechanism that uses JWT for verifying user identities and handling sessions.
- **Create user roles and permissions:** Establish a system for role-based access control that defines what each user type can see and do within the application.
- **Develop functionality to create and manage user accounts:** Build robust APIs to allow users to register, update their profiles, and manage their accounts securely.

### Airplane Management (Saman)
- **API endpoints for admin to manage airplanes:** Create a set of API endpoints to handle the creation, modification, retrieval, and deletion of airplane data by the admin.
- **Classify airplanes by type:** Write the business logic that enables airplanes to be categorized and filtered by their respective classes - military, passenger, or training.
- **Functionality for admins to view and manage orders:** Develop backend functionalities that allow admins to monitor and manage all orders, including updating their statuses.

### Order System and Customization (MHosein)
- **API endpoints for order placement and viewing:** Develop endpoints for users to request the construction of airplanes and view both active and past orders.
- **Customization options for airplanes:** Implement the backend logic that handles various customization options chosen by the users for different types of airplanes.
- **Save incomplete orders:** Code a system that can persistently store the progress of incomplete orders so that users can return and complete their personalization later.

### Pricing and Payments (Mohsen)
- **Pricing algorithm based on selections:** Construct an algorithm that dynamically calculates the cost of the airplane based on the base price and added customizations.
- **Mock payment service integration:** Connect the system with a mock banking service to simulate financial transactions for advance and final payments.
- **Processing payments:** Create APIs that manage the payment process, including capturing, processing, and verifying payment transactions.

### Order Status Management (Saman)
- **Approve, reject, or update order status:** Program the logic that allows the admin to manage the lifecycle of orders, from approval to construction and final delivery.
- **Notification system:** Implement backend functionality to notify users when there's a change in the status of their orders, enhancing the user experience.

### Testing (Narges & MHosein)
- **Write unit tests:** Develop a suite of unit tests to verify the functionality of individual components or methods.
- **Integration tests for API and database:** Write comprehensive integration tests that ensure the application components work together as expected, particularly for API endpoints and database interactions.

### API Documentation (Parsa)
- **Comprehensive API documentation:** Utilize a tool like Swagger to produce user-friendly and interactive API documentation.
- **Document all endpoints:** Make sure all API endpoints are well-documented, including clear descriptions, required parameters, and example responses.

### Performance and Security Considerations (All Developers)
- **Optimize queries for performance:** Ensure that database queries are efficient and optimized to handle the expected load.
- **Secure API endpoints:** Implement best practices to make API endpoints secure against common security threats and vulnerabilities.

### Maintenance and Clean Code (All Developers)
- **Refactor code for clean code principles:** Regularly review and refactor code to align with clean coding principles, aiming for readability and maintainability.
- **Regular code reviews:** Engage in peer code reviews to maintain high code quality and to share knowledge within the team.

### Cross-cutting Concerns (All Developers)
- **Write appropriate commit messages:** Practice writing clear, concise commit messages that accurately reflect the changes made.
- **Comment and document code:** Ensure code is accompanied by relevant comments and documentation to provide context and clarity.
- **Ensure well-tested work:** Vigilantly write and maintain tests for new and existing features to reduce bugs and improve code quality.
- **Follow security best practices:** Consistently apply security best practices throughout all coding activities to mitigate potential risks.


## "Recommended" Models

### User Model
- **Attributes:**
    - User ID
    - Name
    - Email
    - Password (hashed)
    - Role (e.g., admin, military, passenger, educational)
    - Account status
    - Created at
    - Updated at

### Airplane Model
- **Attributes:**
    - Airplane ID
    - Type (military, passenger, training)
    - Base price
    - Specifications (e.g., range, capacity, engine type)
    - Created at
    - Updated at

### Order Model
- **Attributes:**
    - Order ID
    - User ID (Foreign Key from User Model)
    - Airplane ID (Foreign Key from Airplane Model)
    - Customization details (e.g., VIP seats, painting, interior facilities)
    - Status (pending, approved, under construction, built, delivered)
    - Price
    - Advance payment status
    - Final payment status
    - Created at
    - Updated at

### Customization Model (could be part of Order Model or separate based on complexity)
- **Attributes:**
    - Customization ID
    - Order ID (Foreign Key from Order Model)
    - VIP seats count
    - Exterior painting design and color
    - Seat configuration
    - Additional facilities (TV, sockets, etc.)
    - Cockpit facilities level
    - Created at
    - Updated at

### Payment Model
- **Attributes:**
    - Payment ID
    - Order ID (Foreign Key from Order Model)
    - User ID (Foreign Key from User Model)
    - Amount
    - Payment type (advance, final)
    - Payment status (completed, pending, failed)
    - Transaction ID (from Banktest or other payment gateway)
    - Created at
    - Updated at

### Authentication Model (usually not persisted but can be for tokens or sessions)
- **Attributes:**
    - Token ID
    - User ID (Foreign Key from User Model)
    - JWT Token
    - Expiry
    - Created at


## "Recommended" API Endpoints

### User-related Endpoints

- **POST /users/register**
    - Register a new user account.

- **POST /users/login**
    - Authenticate a user and return a JWT token.

- **GET /users/profile**
    - Retrieve the logged-in user's profile information.

- **PUT /users/profile**
    - Update the logged-in user's profile information.

- **GET /users/:id**
    - Retrieve a specific user's profile (admin only).

### Airplane-related Endpoints

- **POST /airplanes**
    - Admin can add a new airplane with specifications.

- **GET /airplanes**
    - List all airplanes, possibly with query parameters for filtering by class (military, passenger, training).

- **GET /airplanes/:id**
    - Get details of a specific airplane.

- **PUT /airplanes/:id**
    - Admin can update specifications of an existing airplane.

- **DELETE /airplanes/:id**
    - Admin can remove an airplane from the listing.

### Order-related Endpoints

- **POST /orders**
    - Place a new order for an airplane with customizations.

- **GET /orders**
    - Retrieve a list of orders placed by the logged-in user.

- **GET /orders/:id**
    - Retrieve details of a specific order.

- **PUT /orders/:id/customize**
    - Update the customization details of an existing order.

- **PUT /orders/:id/status**
    - Admin can update the status of an order (e.g., approved, under construction, built).

- **GET /admin/orders**
    - Admin can view all orders placed by all users.

### Payment-related Endpoints

- **POST /payments/advance**
    - Process an advance payment for an order.

- **POST /payments/finalize**
    - Process the final payment upon order completion.

- **GET /payments/orders/:orderId**
    - Retrieve payment details for a specific order.

### Admin-specific Endpoints (for managing the application)

- **GET /admin/users**
    - Retrieve a list of all users and their roles.

- **PUT /admin/users/:id**
    - Update user roles or statuses (activate, deactivate accounts).

- **DELETE /admin/users/:id**
    - Delete a user account.

- **GET /admin/orders/pending**
    - Retrieve a list of all pending orders for review.

- **PUT /admin/orders/:id/approve**
    - Approve a specific order.

- **PUT /admin/orders/:id/reject**
    - Reject a specific order.

### Reporting and Analytics Endpoints (optional for future development)

- **GET /admin/reports/sales**
    - Get sales reports for a given time period.

- **GET /admin/reports/user-activity**
    - Get user activity reports.
