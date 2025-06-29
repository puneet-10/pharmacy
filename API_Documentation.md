# Pharmacy API Documentation

## Base URL
```
http://localhost:8080
```

## Authentication
Most endpoints require JWT authentication. Include the JWT token in the Authorization header:
```
Authorization: Bearer <jwt_token>
```

---

## User Management APIs

### 1. User Signup
**Endpoint:** `POST /signup`  
**Authentication:** Not required  
**Description:** Register a new user

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "phone": "1234567890",
  "password": "password123",
  "firm_name": "ABC Pharmacy",
  "is_admin": false
}
```

**Response (201 Created):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "phoneNumber": "1234567890",
    "name": "John Doe",
    "email": "john@example.com",
    "firmName": "ABC Pharmacy",
    "isAdmin": false
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid input or missing required fields (Name, Phone, and Password are required)
- `409 Conflict`: Phone number already in use

---

### 2. User Authentication
**Endpoint:** `POST /authenticate`  
**Authentication:** Not required  
**Description:** Authenticate user and get JWT token

**Request Body:**
```json
{
  "identifier": "1234567890",
  "password": "password123"
}
```

**Response (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "phoneNumber": "1234567890",
    "name": "John Doe",
    "email": "john@example.com",
    "firmName": "ABC Pharmacy",
    "isAdmin": false
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid input
- `401 Unauthorized`: Invalid credentials

---

### 3. Update User
**Endpoint:** `PUT /user/{id}`  
**Authentication:** Not required  
**Description:** Update user information

**Path Parameters:**
- `id` (integer): User ID

**Request Headers:**
- `X-Updated-By` (optional): Who updated the user (defaults to "system")

**Request Body:**
```json
{
  "name": "John Smith",
  "phone": "0987654321",
  "firm_name": "XYZ Pharmacy",
  "is_admin": true
}
```

**Response (200 OK):**
```json
{
  "id": 1,
  "name": "John Smith",
  "email": "john@example.com",
  "phone": "0987654321",
  "password": "hashed_password",
  "firm_name": "XYZ Pharmacy",
  "is_admin": true,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-02T00:00:00Z",
  "updated_by": "admin"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid input
- `500 Internal Server Error`: Database error

---

## Company Management APIs

### 4. Create Company
**Endpoint:** `POST /companies`  
**Authentication:** Not required  
**Description:** Create a new company

**Request Body:**
```json
{
  "company_name": "Pharma Corp",
  "description": "Leading pharmaceutical company",
  "updated_by": "admin",
  "logo_url": "https://example.com/logo.png"
}
```

**Response (201 Created):**
```json
{
  "id": 1,
  "company_name": "Pharma Corp",
  "description": "Leading pharmaceutical company",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "updated_by": "admin",
  "logo_url": "https://example.com/logo.png"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid data
- `500 Internal Server Error`: Could not create company

---

### 5. Get All Companies
**Endpoint:** `GET /companies`  
**Authentication:** Not required  
**Description:** Retrieve all companies

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "company_name": "Pharma Corp",
    "description": "Leading pharmaceutical company",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "updated_by": "admin",
    "logo_url": "https://example.com/logo.png"
  }
]
```

**Error Responses:**
- `500 Internal Server Error`: Could not retrieve companies

---

### 6. Get Company by ID
**Endpoint:** `GET /companies/{id}`  
**Authentication:** Not required  
**Description:** Retrieve a specific company

**Path Parameters:**
- `id` (integer): Company ID

**Response (200 OK):**
```json
{
  "id": 1,
  "company_name": "Pharma Corp",
  "description": "Leading pharmaceutical company",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "updated_by": "admin",
  "logo_url": "https://example.com/logo.png"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid company ID
- `404 Not Found`: Company not found

---

### 7. Update Company
**Endpoint:** `PUT /companies/{id}`  
**Authentication:** Not required  
**Description:** Update an existing company

**Path Parameters:**
- `id` (integer): Company ID

**Request Body:**
```json
{
  "company_name": "Updated Pharma Corp",
  "description": "Updated description",
  "updated_by": "admin",
  "logo_url": "https://example.com/new-logo.png"
}
```

**Response (200 OK):**
```json
{
  "id": 1,
  "company_name": "Updated Pharma Corp",
  "description": "Updated description",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-02T00:00:00Z",
  "updated_by": "admin",
  "logo_url": "https://example.com/new-logo.png"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid company ID or data
- `500 Internal Server Error`: Could not update company

---

### 8. Delete Company
**Endpoint:** `DELETE /companies/{id}`  
**Authentication:** Not required  
**Description:** Delete a company

**Path Parameters:**
- `id` (integer): Company ID

**Response (204 No Content):**
```
(Empty response body)
```

**Error Responses:**
- `400 Bad Request`: Invalid company ID
- `500 Internal Server Error`: Could not delete company

---

## Medicine Management APIs

### 9. Create Medicine
**Endpoint:** `POST /medicines`  
**Authentication:** Not required  
**Description:** Create a new medicine

**Request Body:**
```json
{
  "name": "Aspirin",
  "description": "Pain reliever",
  "company_id": 1,
  "updated_by": "admin",
  "offer": "10% off"
}
```

**Response (201 Created):**
```json
{
  "id": 1,
  "name": "Aspirin",
  "description": "Pain reliever",
  "company_id": 1,
  "company": {
    "id": 1,
    "company_name": "Pharma Corp",
    "description": "Leading pharmaceutical company",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "updated_by": "admin",
    "logo_url": "https://example.com/logo.png"
  },
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "updated_by": "admin",
  "offer": "10% off"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid data
- `500 Internal Server Error`: Could not create medicine

---

### 10. Get All Medicines
**Endpoint:** `GET /medicines`  
**Authentication:** Not required  
**Description:** Retrieve all medicines grouped by company

**Response (200 OK):**
```json
[
  {
    "companyId": 1,
    "companyName": "Pharma Corp",
    "medicines": [
      {
        "medicineId": 1,
        "name": "Aspirin",
        "offer": "10% off"
      },
      {
        "medicineId": 2,
        "name": "Ibuprofen",
        "offer": "Buy 2 Get 1 Free"
      }
    ]
  }
]
```

**Error Responses:**
- `500 Internal Server Error`: Could not retrieve medicines

---

### 11. Get Medicine by ID
**Endpoint:** `GET /medicines/{id}`  
**Authentication:** Not required  
**Description:** Retrieve a specific medicine

**Path Parameters:**
- `id` (integer): Medicine ID

**Response (200 OK):**
```json
{
  "id": 1,
  "name": "Aspirin",
  "description": "Pain reliever",
  "company_id": 1,
  "company": {
    "id": 1,
    "company_name": "Pharma Corp",
    "description": "Leading pharmaceutical company",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "updated_by": "admin",
    "logo_url": "https://example.com/logo.png"
  },
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "updated_by": "admin",
  "offer": "10% off"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid medicine ID
- `500 Internal Server Error`: Could not retrieve medicine

---

### 12. Update Medicine
**Endpoint:** `PUT /medicines/{id}`  
**Authentication:** Not required  
**Description:** Update an existing medicine

**Path Parameters:**
- `id` (integer): Medicine ID

**Request Body:**
```json
{
  "name": "Updated Aspirin",
  "description": "Updated pain reliever",
  "company_id": 1,
  "updated_by": "admin",
  "offer": "15% off"
}
```

**Response (200 OK):**
```json
{
  "id": 1,
  "name": "Updated Aspirin",
  "description": "Updated pain reliever",
  "company_id": 1,
  "company": {
    "id": 1,
    "company_name": "Pharma Corp",
    "description": "Leading pharmaceutical company",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "updated_by": "admin",
    "logo_url": "https://example.com/logo.png"
  },
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-02T00:00:00Z",
  "updated_by": "admin",
  "offer": "15% off"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid medicine ID or data
- `500 Internal Server Error`: Could not update medicine

---

### 13. Delete Medicine
**Endpoint:** `DELETE /medicines/{id}`  
**Authentication:** Not required  
**Description:** Delete a medicine

**Path Parameters:**
- `id` (integer): Medicine ID

**Response (200 OK):**
```json
{
  "message": "Medicine deleted successfully"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid medicine ID
- `500 Internal Server Error`: Could not delete medicine

---

### 14. Upload Medicines CSV
**Endpoint:** `POST /medicines/upload`  
**Authentication:** Not required  
**Description:** Upload medicines from CSV file  
**Content-Type:** `multipart/form-data`

**Request Body (Form Data):**
- `file`: CSV file containing medicine data

**CSV Format:**
```csv
Name,Description,CompanyName
Aspirin,Pain reliever,Pharma Corp
Ibuprofen,Anti-inflammatory,Pharma Corp
```

**Response (200 OK):**
```json
{
  "message": "Medicines uploaded successfully"
}
```

**Error Responses:**
- `400 Bad Request`: CSV file is required
- `500 Internal Server Error`: File processing error

---

### 15. Update Medicine Offer
**Endpoint:** `PUT /medicines/offer`  
**Authentication:** Not required  
**Description:** Update offer for specific medicine or all medicines in a company

**Request Body:**
```json
{
  "medicine_id": 1,
  "company_id": 1,
  "offer": "20% off",
  "updated_by": "admin"
}
```

**Note:** 
- If `medicine_id` is 0 or omitted, the offer will be updated for all medicines in the specified company
- If `medicine_id` is provided, the offer will be updated for that specific medicine only

**Response (200 OK):**
```json
{
  "message": "Offer updated for the specific medicine"
}
```

**Or for company-wide update:**
```json
{
  "message": "Offer updated for all medicines in the company"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid data or missing required fields
- `500 Internal Server Error`: Could not update offer

---

## Order Management APIs

### 16. Create Order
**Endpoint:** `POST /orders`  
**Authentication:** Required (JWT Token)  
**Description:** Create a new order

**Request Body:**
```json
{
  "items": [
    {
      "medicineId": 1,
      "companyId": 1,
      "quantity": 2
    },
    {
      "medicineId": 2,
      "companyId": 1,
      "quantity": 1
    }
  ]
}
```

**Response (200 OK):**
```json
{
  "orderId": 1,
  "userId": 1,
  "items": [
    {
      "medicineId": 1,
      "medicineName": "Aspirin",
      "companyId": 1,
      "companyName": "Pharma Corp",
      "quantity": 2
    },
    {
      "medicineId": 2,
      "medicineName": "Ibuprofen",
      "companyId": 1,
      "companyName": "Pharma Corp",
      "quantity": 1
    }
  ],
  "createdAt": "2024-01-01T00:00:00Z"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request
- `401 Unauthorized`: Missing or invalid JWT token
- `500 Internal Server Error`: Failed to create order

---

### 17. Get Order by ID
**Endpoint:** `GET /orders/{id}`  
**Authentication:** Not required  
**Description:** Retrieve a specific order

**Path Parameters:**
- `id` (integer): Order ID

**Response (200 OK):**
```json
{
  "orderId": 1,
  "userId": 1,
  "items": [
    {
      "medicineId": 1,
      "medicineName": "Aspirin",
      "companyId": 1,
      "companyName": "Pharma Corp",
      "quantity": 2
    }
  ],
  "status": "pending",
  "createdAt": "2024-01-01T00:00:00Z"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid ID
- `404 Not Found`: Order not found

---

### 18. Get All Orders
**Endpoint:** `GET /orders`  
**Authentication:** Required (JWT Token)  
**Description:** Retrieve all orders (filtered by user if not admin)

**Response (200 OK):**
```json
[
  {
    "orderId": 1,
    "userId": 1,
    "items": [
      {
        "medicineId": 1,
        "medicineName": "Aspirin",
        "companyId": 1,
        "companyName": "Pharma Corp",
              "quantity": 2
    }
  ],
  "status": "pending",
  "createdAt": "2024-01-01T00:00:00Z"
  }
]
```

**Note:** 
- Regular users can only see their own orders
- Admin users can see all orders

**Error Responses:**
- `401 Unauthorized`: Missing or invalid JWT token
- `500 Internal Server Error`: Failed to fetch orders

---

### 19. Update Order
**Endpoint:** `PUT /orders/{id}`  
**Authentication:** Not required  
**Description:** Update an existing order

**Path Parameters:**
- `id` (integer): Order ID

**Request Body:**
```json
{
  "userId": 1,
  "items": [
    {
      "medicineId": 1,
      "companyId": 1,
      "quantity": 3
    }
  ]
}
```

**Response (200 OK):**
```json
{
  "orderId": 1,
  "userId": 1,
  "items": [
    {
      "medicineId": 1,
      "medicineName": "Aspirin",
      "companyId": 1,
      "companyName": "Pharma Corp",
      "quantity": 3
    }
  ],
  "status": "pending",
  "createdAt": "2024-01-01T00:00:00Z"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid ID or request body
- `500 Internal Server Error`: Failed to update order

---

### 20. Update Order Status
**Endpoint:** `PUT /orders/{id}/status`  
**Authentication:** Required (JWT Token)  
**Description:** Update the status of an existing order

**Path Parameters:**
- `id` (integer): Order ID

**Request Body:**
```json
{
  "status": "processing"
}
```

**Valid Status Values:**
- `pending` - Order is pending (default when created)
- `processing` - Order is being processed
- `shipped` - Order has been shipped
- `delivered` - Order has been delivered
- `cancelled` - Order has been cancelled

**Response (200 OK):**
```json
{
  "orderId": 1,
  "userId": 1,
  "items": [
    {
      "medicineId": 1,
      "medicineName": "Aspirin",
      "companyId": 1,
      "companyName": "Pharma Corp",
      "quantity": 2
    }
  ],
  "status": "processing",
  "createdAt": "2024-01-01T00:00:00Z"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid ID, request body, or invalid status value
- `401 Unauthorized`: Missing or invalid JWT token
- `500 Internal Server Error`: Failed to update order status

---

### 21. Delete Order
**Endpoint:** `DELETE /orders/{id}`  
**Authentication:** Not required  
**Description:** Delete an order

**Path Parameters:**
- `id` (integer): Order ID

**Response (200 OK):**
```json
{
  "message": "Order deleted successfully"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid ID
- `500 Internal Server Error`: Failed to delete order

---

## Error Response Format

All error responses follow this format:

```json
{
  "message": "Error description"
}
```

Common HTTP status codes used:
- `200 OK`: Successful operation
- `201 Created`: Resource created successfully
- `204 No Content`: Successful operation with no response body
- `400 Bad Request`: Invalid request data
- `401 Unauthorized`: Authentication required or failed
- `404 Not Found`: Resource not found
- `409 Conflict`: Resource conflict (e.g., duplicate phone number)
- `500 Internal Server Error`: Server-side error

---

## JWT Token Structure

The JWT token contains the following claims:
```json
{
  "userId": 1,
  "name": "John Doe",
  "email": "john@example.com",
  "phone": "1234567890",
  "firmName": "ABC Pharmacy",
  "isAdmin": false,
  "exp": 1704067200
}
```

Token expires after 24 hours from generation.