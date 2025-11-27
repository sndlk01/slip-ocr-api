# OCR Payment Slip API

A RESTful API built with Go for processing Thai bank payment slips using OCR (Optical Character Recognition). Extract transaction details from slip images automatically.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)](https://www.docker.com/)

---

## 🎯 Overview

A **complete financial management API** with OCR capabilities, budget tracking, and analytics:

1. **Upload slip(s)** → Auto-extract data via OCR → Detect duplicates → Save transactions
2. **Budget management** → Set monthly limits → Track spending → Get warnings
3. **Subscription tracking** → Auto-detect services → Monitor recurring payments
4. **Analytics & insights** → Monthly trends → Yearly comparison → Category breakdown

**Core Features:**
- 🔐 User authentication with JWT tokens (Login/Register)
- 📸 Multi-file upload with duplicate detection
- 🤖 Advanced OCR for Thai bank slips (SCB, KBANK, BBL, KTB)
- 💰 Budget system with automatic warnings
- 🔄 Subscription auto-detection (Netflix, Spotify, etc.)
- 📊 Dashboard APIs for data visualization
- 📝 Category-based expense tracking
- ✏️ Manual transaction creation and editing

**Tech Stack:** Go, Gin, GORM, SQLite, Tesseract OCR

---

## 🐳 Quick Start with Docker (Recommended)

```bash
# Start the service
docker-compose up -d

# Test it
curl http://localhost:8077/health
```

That's it! API is running at `http://localhost:8077`

---

## 📦 Manual Installation

### Prerequisites

1. **Go 1.21+**
   ```bash
   # Download from https://golang.org/dl/
   go version
   ```

2. **Tesseract OCR with Thai language**

   **Windows:**
   ```powershell
   choco install tesseract
   # Or download from: https://github.com/UB-Mannheim/tesseract/wiki
   # Select Thai language during installation
   ```

   **macOS:**
   ```bash
   brew install tesseract tesseract-lang
   ```

   **Linux:**
   ```bash
   sudo apt install tesseract-ocr tesseract-ocr-tha tesseract-ocr-eng
   ```

   **Verify:**
   ```bash
   tesseract --version
   tesseract --list-langs  # Should show: eng, tha
   ```

3. **C Compiler (for CGO)**
   - **Windows:** Install [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) or [MinGW-w64](https://www.mingw-w64.org/)
   - **macOS:** `xcode-select --install`
   - **Linux:** `sudo apt install build-essential`

### Installation

```bash
# Navigate to project
cd D:\Projects\ocr-api

# Install dependencies
go mod download

# Run the application
go run main.go
```

Server starts on `http://localhost:8077`

---

## 📡 API Documentation

**Base URL:** `http://localhost:8077`

### Endpoints

#### Authentication
| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/v1/auth/register` | Register new user |
| `POST` | `/api/v1/auth/login` | Login and get JWT token |
| `GET` | `/api/v1/auth/profile` | Get user profile (requires auth) |

#### Core Transactions
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/health` | Health check |
| `POST` | `/api/v1/upload` | Upload slip(s) - **multi-file, duplicate detection, auto-detect subscriptions** |
| `POST` | `/api/v1/transactions` | Create manual transaction |
| `GET` | `/api/v1/transactions` | List transactions (filter by type/bank/category) |
| `GET` | `/api/v1/transactions/:id` | Get transaction details |
| `PUT/PATCH` | `/api/v1/transactions/:id` | Update transaction |
| `DELETE` | `/api/v1/transactions/:id` | Delete transaction |

#### Budget Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/v1/budgets` | Create monthly budget |
| `GET` | `/api/v1/budgets` | List all budgets |
| `GET` | `/api/v1/budgets/status` | Get budget status (spent/remaining/warnings) |
| `DELETE` | `/api/v1/budgets/:id` | Delete budget |

#### Subscription Tracking
| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/v1/subscriptions` | Add subscription |
| `GET` | `/api/v1/subscriptions` | List all subscriptions |
| `DELETE` | `/api/v1/subscriptions/:id` | Delete subscription |

#### Analytics & Dashboard
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/v1/dashboard/monthly` | Monthly trend (12 months) |
| `GET` | `/api/v1/dashboard/yearly` | Yearly comparison |
| `GET` | `/api/v1/dashboard/categories` | Category breakdown (pie chart data) |
| `GET` | `/api/v1/summary/monthly` | Monthly summary with categories |

---

### 1. Health Check

```bash
GET /health
```

**Example:**
```bash
curl http://localhost:8077/health
```

**Response:**
```json
{
  "status": "ok",
  "message": "OCR API is running"
}
```

---

### 2. User Registration

```bash
POST /api/v1/auth/register
Content-Type: application/json
```

**Parameters:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `username` | String | Yes | Unique username |
| `email` | String | Yes | Valid email address |
| `password` | String | Yes | Password (min 6 characters) |
| `full_name` | String | No | Full name |

**Example:**
```bash
curl -X POST http://localhost:8077/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "somchai",
    "email": "somchai@example.com",
    "password": "mypassword123",
    "full_name": "นายสมชาย ใจดี"
  }'
```

**Response (201):**
```json
{
  "message": "User registered successfully",
  "user": {
    "id": 1,
    "username": "somchai",
    "email": "somchai@example.com",
    "full_name": "นายสมชาย ใจดี",
    "created_at": "2025-11-27T10:00:00Z"
  }
}
```

**Error Responses:**
```json
// 400 - Validation error
{ "error": "Invalid request body" }

// 400 - Duplicate user
{ "error": "username or email already exists" }
```

---

### 3. User Login

```bash
POST /api/v1/auth/login
Content-Type: application/json
```

**Parameters:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `username` | String | Yes | Username or email |
| `password` | String | Yes | Password |

**Example:**
```bash
curl -X POST http://localhost:8077/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "somchai",
    "password": "mypassword123"
  }'
```

**Response (200):**
```json
{
  "message": "Login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "somchai",
    "email": "somchai@example.com",
    "full_name": "นายสมชาย ใจดี",
    "created_at": "2025-11-27T10:00:00Z"
  }
}
```

**Error Response (401):**
```json
{ "error": "invalid username or password" }
```

**Note:** Save the `token` value and include it in subsequent API requests using the `Authorization` header:
```bash
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

---

### 4. Get User Profile

```bash
GET /api/v1/auth/profile
Authorization: Bearer <token>
```

**Example:**
```bash
curl http://localhost:8077/api/v1/auth/profile \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Response (200):**
```json
{
  "user": {
    "id": 1,
    "username": "somchai",
    "email": "somchai@example.com",
    "full_name": "นายสมชาย ใจดี",
    "created_at": "2025-11-27T10:00:00Z"
  }
}
```

**Error Response (401):**
```json
{ "error": "Unauthorized" }
```

---

### 5. Upload Payment Slip(s)

```bash
POST /api/v1/upload
Content-Type: multipart/form-data
```

**Parameters:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `slip` or `slips` | File(s) | Yes | Image file(s) (JPG, JPEG, PNG, max 10MB each) |
| `type` | String | Yes | `income` or `expense` |

**Example - Single File:**
```bash
curl -X POST http://localhost:8077/api/v1/upload \
  -F "slip=@/path/to/slip.jpg" \
  -F "type=expense"
```

**Example - Multiple Files:**
```bash
curl -X POST http://localhost:8077/api/v1/upload \
  -F "slips=@slip1.jpg" \
  -F "slips=@slip2.jpg" \
  -F "slips=@slip3.jpg" \
  -F "type=expense"
```

**Example - PowerShell (Multiple Files):**
```powershell
$uri = "http://localhost:8077/api/v1/upload"
$form = @{
    slips = Get-Item -Path "C:\path\to\slip1.jpg","C:\path\to\slip2.jpg"
    type = "expense"
}
Invoke-RestMethod -Uri $uri -Method Post -Form $form
```

**Example - JavaScript:**
```javascript
const formData = new FormData();
// Single file
formData.append('slip', fileInput.files[0]);
// Or multiple files
for (let file of fileInput.files) {
  formData.append('slips', file);
}
formData.append('type', 'expense');

fetch('http://localhost:8077/api/v1/upload', {
  method: 'POST',
  body: formData
})
  .then(res => res.json())
  .then(data => console.log(data));
```

**Success Response (201) - Single File:**
```json
{
  "message": "Processed 1 out of 1 slips successfully",
  "success_count": 1,
  "total_count": 1,
  "transactions": [
    {
      "id": 1,
      "type": "expense",
      "amount": 1500.00,
      "date": "25/11/2025",
      "time": "14:30:25",
      "reference": "T123456789012",
      "bank": "SCB",
      "sender": "นายสมชาย ใจดี",
      "receiver": "นางสมหญิง รักสนุก",
      "detail": "",
      "raw_ocr_text": "ชําระเงินสําเร็จ\n21 ต.ค. 68 14:00 น.\n...",
      "created_at": "2025-11-25T08:30:00Z"
    }
  ]
}
```

**Success Response (201) - Multiple Files with Errors:**
```json
{
  "message": "Processed 2 out of 3 slips successfully",
  "success_count": 2,
  "total_count": 3,
  "transactions": [...],
  "errors": [
    "Failed to process 'slip3.jpg': image too dark"
  ]
}
```

**Error Responses:**
```json
// 400 - Missing/invalid type
{ "error": "Missing or invalid 'type' field. Must be 'income' or 'expense'" }

// 400 - No file
{ "error": "No files uploaded. Use 'slip' or 'slips' as the form field name" }

// 400 - File too large
{ "error": "File 'slip1.jpg' is too large (max 10MB)" }

// 400 - Invalid file type
{ "error": "invalid file type: .pdf. Allowed types: jpg, jpeg, png" }
```

---

### 6. Create Transaction Manually

```bash
POST /api/v1/transactions
Content-Type: application/json
```

**Use Case:** Create income/expense transactions without a physical slip (e.g., cash income, online transfers)

**Parameters:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `type` | String | Yes | `income` or `expense` |
| `amount` | Number | Yes | Transaction amount |
| `date` | String | No | Date (e.g., "26/11/2025") |
| `time` | String | No | Time (e.g., "10:30") |
| `reference` | String | No | Reference number |
| `bank` | String | No | Bank name (SCB, KBANK, BBL, KTB) |
| `sender` | String | No | Sender name |
| `receiver` | String | No | Receiver name |
| `detail` | String | No | Custom description/notes |

**Example:**
```bash
curl -X POST http://localhost:8077/api/v1/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "type": "income",
    "amount": 5000,
    "date": "26/11/2025",
    "time": "10:00",
    "bank": "SCB",
    "sender": "นายสมชาย",
    "detail": "รายได้จากการขาย"
  }'
```

**Response (201):**
```json
{
  "message": "Transaction created successfully",
  "transaction": {
    "id": 10,
    "type": "income",
    "amount": 5000,
    "date": "26/11/2025",
    "time": "10:00",
    "bank": "SCB",
    "sender": "นายสมชาย",
    "detail": "รายได้จากการขาย",
    "created_at": "2025-11-26T10:00:00Z",
    "updated_at": "2025-11-26T10:00:00Z"
  }
}
```

---

### 7. Get All Transactions

```bash
GET /api/v1/transactions
```

**Query Parameters (Optional):**
- `type`: Filter by `income` or `expense`
- `bank`: Filter by `SCB`, `KBANK`, `BBL`, or `KTB`

**Examples:**
```bash
# Get all
curl http://localhost:8077/api/v1/transactions

# Filter by type
curl "http://localhost:8077/api/v1/transactions?type=income"

# Filter by bank
curl "http://localhost:8077/api/v1/transactions?bank=SCB"

# Combine filters
curl "http://localhost:8077/api/v1/transactions?type=income&bank=SCB"
```

**Response (200):**
```json
{
  "transactions": [
    {
      "id": 1,
      "type": "income",
      "amount": 1500.00,
      "date": "25/11/2025",
      "time": "14:30:25",
      "reference": "T123456789012",
      "bank": "SCB",
      "sender": "นายสมชาย ใจดี",
      "receiver": "นางสมหญิง รักสนุก",
      "detail": "",
      "created_at": "2025-11-25T08:30:00Z"
    }
  ]
}
```

---

### 8. Get Transaction by ID

```bash
GET /api/v1/transactions/:id
```

**Example:**
```bash
curl http://localhost:8077/api/v1/transactions/1
```

**Response (200):**
```json
{
  "transaction": {
    "id": 1,
    "type": "income",
    "amount": 1500.00,
    "date": "25/11/2025",
    "time": "14:30:25",
    "reference": "T123456789012",
    "bank": "SCB",
    "sender": "นายสมชาย ใจดี",
    "receiver": "นางสมหญิง รักสนุก",
    "detail": "",
    "raw_ocr_text": "ธนาคารไทยพาณิชย์\nจำนวนเงิน 1,500.00\n...",
    "created_at": "2025-11-25T08:30:00Z"
  }
}
```

**Error (404):**
```json
{ "error": "Transaction not found" }
```

---

### 9. Update Transaction

```bash
PUT /api/v1/transactions/:id
PATCH /api/v1/transactions/:id
Content-Type: application/json
```

**Use Case:** Update transaction details (e.g., add description, correct amount)

**Parameters (all optional, send only what you want to update):**

| Field | Type | Description |
|-------|------|-------------|
| `amount` | Number | Transaction amount |
| `date` | String | Date |
| `time` | String | Time |
| `reference` | String | Reference number |
| `bank` | String | Bank name |
| `sender` | String | Sender name |
| `receiver` | String | Receiver name |
| `detail` | String | Custom description/notes |

**Example - Update detail only:**
```bash
curl -X PATCH http://localhost:8077/api/v1/transactions/1 \
  -H "Content-Type: application/json" \
  -d '{"detail": "ค่าน้ำมัน รถยนต์"}'
```

**Example - Update multiple fields:**
```bash
curl -X PUT http://localhost:8077/api/v1/transactions/1 \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 1600,
    "detail": "ค่าอาหาร (แก้ไขจำนวนเงิน)",
    "date": "26/11/2025"
  }'
```

**Response (200):**
```json
{
  "message": "Transaction updated successfully",
  "transaction": {
    "id": 1,
    "type": "expense",
    "amount": 1600,
    "date": "26/11/2025",
    "time": "14:30:25",
    "bank": "SCB",
    "detail": "ค่าอาหาร (แก้ไขจำนวนเงิน)",
    "updated_at": "2025-11-26T11:00:00Z"
  }
}
```

**Error Responses:**
```json
// 404 - Not found
{ "error": "transaction not found: record not found" }

// 400 - No fields to update
{ "error": "No fields to update" }

// 400 - Invalid request
{ "error": "Invalid request body" }
```

---

### 10. Delete Transaction

```bash
DELETE /api/v1/transactions/:id
```

**Example:**
```bash
curl -X DELETE http://localhost:8077/api/v1/transactions/1
```

**Response (200):**
```json
{
  "message": "Transaction deleted successfully"
}
```

---

## 🏗️ Project Structure

```
ocr-api/
├── main.go                          # Entry point
├── docker-compose.yml               # Docker setup
├── config/
│   ├── config.go                   # Configuration
│   └── database.go                 # Database setup (auto-migration)
├── models/
│   ├── user.go                     # User model (Authentication)
│   ├── transaction.go              # Transaction model
│   ├── budget.go                   # Budget model
│   └── subscription.go             # Subscription model
├── controllers/
│   ├── auth_controller.go          # Authentication (Login/Register)
│   ├── upload_controller.go        # Upload handler (duplicate detection)
│   ├── transaction_controller.go   # Transaction CRUD
│   ├── budget_controller.go        # Budget management
│   ├── subscription_controller.go  # Subscription tracking
│   └── dashboard_controller.go     # Analytics endpoints
├── services/
│   ├── auth_service.go             # Authentication service (JWT)
│   ├── ocr_service.go              # OCR workflow + subscription detection
│   ├── transaction_service.go      # Transaction service + duplicate check
│   ├── budget_service.go           # Budget calculations
│   ├── subscription_service.go     # Subscription auto-detection
│   └── dashboard_service.go        # Analytics & reporting
├── ocr/
│   ├── tesseract.go                # Tesseract wrapper
│   ├── preprocessor.go             # Image preprocessing
│   └── extractor.go                # Data extraction (Thai date support)
├── routes/
│   └── routes.go                   # API routes (26 endpoints)
└── utils/
    ├── helpers.go                  # Utilities (OCR text cleaning)
    └── jwt.go                      # JWT token generation & validation
```

---

## ⚙️ Configuration

Environment variables (optional):

```bash
SERVER_PORT=8077                    # Server port
DATABASE_PATH=./data/db.sqlite      # Database file
UPLOAD_DIR=./uploads                # Temp upload directory
TESSERACT_LANG=tha+eng             # OCR languages
MAX_UPLOAD_SIZE=10485760           # Max file size (10MB)
```

---

## 🧪 Testing

### Basic Tests
```bash
# Health check
curl http://localhost:8077/health

# Upload slip (with duplicate detection & subscription auto-detect)
curl -X POST http://localhost:8077/api/v1/upload \
  -F "slip=@netflix_slip.jpg" \
  -F "type=expense"

# Upload multiple slips
curl -X POST http://localhost:8077/api/v1/upload \
  -F "slips=@slip1.jpg" \
  -F "slips=@slip2.jpg" \
  -F "type=expense"
```

### Transaction Management
```bash
# Create manual transaction with category
curl -X POST http://localhost:8077/api/v1/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "type": "income",
    "amount": 5000,
    "date": "26/11/2025",
    "category": "รายได้จากการขาย",
    "detail": "ขายสินค้าออนไลน์"
  }'

# List all transactions
curl http://localhost:8077/api/v1/transactions

# Filter by category
curl "http://localhost:8077/api/v1/transactions?type=expense&category=ค่าอาหาร"

# Update transaction
curl -X PATCH http://localhost:8077/api/v1/transactions/1 \
  -H "Content-Type: application/json" \
  -d '{"category": "ค่าเดินทาง", "detail": "ค่าน้ำมัน"}'
```

### Budget Management
```bash
# Set budget
curl -X POST http://localhost:8077/api/v1/budgets \
  -H "Content-Type: application/json" \
  -d '{"category": "ค่าอาหาร", "monthly_limit": 5000, "month": 11, "year": 2025}'

# Check budget status
curl "http://localhost:8077/api/v1/budgets/status?year=2025&month=11"

# List all budgets
curl http://localhost:8077/api/v1/budgets
```

### Subscription Tracking
```bash
# List auto-detected subscriptions
curl http://localhost:8077/api/v1/subscriptions

# Add manual subscription
curl -X POST http://localhost:8077/api/v1/subscriptions \
  -H "Content-Type: application/json" \
  -d '{"name": "Spotify", "amount": 129, "billing_cycle": "monthly"}'
```

### Analytics & Dashboard
```bash
# Monthly trend
curl "http://localhost:8077/api/v1/dashboard/monthly?year=2025"

# Yearly comparison
curl "http://localhost:8077/api/v1/dashboard/yearly?years=2024,2025"

# Category breakdown
curl "http://localhost:8077/api/v1/dashboard/categories?year=2025&month=11&type=expense"

# Monthly summary
curl "http://localhost:8077/api/v1/summary/monthly?year=2025&month=11"
```

---

## 🐛 Common Issues

### Tesseract Not Found
```bash
# Install Tesseract and add to PATH
tesseract --version  # Verify installation
```

### Thai Language Missing
```bash
# Check available languages
tesseract --list-langs  # Should show 'tha'

# Install if missing:
# Windows: Reinstall Tesseract with Thai option
# macOS: brew install tesseract-lang
# Linux: sudo apt install tesseract-ocr-tha
```

### CGO Build Errors
```bash
# Use Docker (recommended) or install C compiler
# Windows: TDM-GCC or MinGW-w64
# macOS: xcode-select --install
# Linux: sudo apt install build-essential
```

### 11. Budget Management

```bash
POST /api/v1/budgets
Content-Type: application/json
```

**Use Case:** Set monthly spending limits per category

**Example:**
```bash
curl -X POST http://localhost:8077/api/v1/budgets \
  -H "Content-Type: application/json" \
  -d '{
    "category": "ค่าอาหาร",
    "monthly_limit": 5000,
    "month": 11,
    "year": 2025
  }'
```

**Check Budget Status:**
```bash
curl "http://localhost:8077/api/v1/budgets/status?year=2025&month=11"
```

**Response:**
```json
{
  "budget_status": [
    {
      "category": "ค่าอาหาร",
      "monthly_limit": 5000,
      "spent": 3200,
      "remaining": 1800,
      "percent_used": 64,
      "status": "ok"
    }
  ]
}
```

---

### 12. Subscription Tracking

**Auto-detected subscriptions** are created automatically when uploading slips with recognized services (Netflix, Spotify, etc.)

**Manual subscription:**
```bash
curl -X POST http://localhost:8077/api/v1/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Netflix Premium",
    "amount": 419,
    "category": "บันเทิง",
    "billing_cycle": "monthly",
    "next_billing_date": "01/12/2025"
  }'
```

**List subscriptions:**
```bash
curl http://localhost:8077/api/v1/subscriptions
```

---

### 13. Dashboard Analytics

**Monthly Trend (for line charts):**
```bash
curl "http://localhost:8077/api/v1/dashboard/monthly?year=2025"
```

**Response:**
```json
{
  "monthly_trend": [
    {"month": "01/2025", "income": 50000, "expense": 30000},
    {"month": "02/2025", "income": 55000, "expense": 32000}
  ]
}
```

**Yearly Comparison:**
```bash
curl "http://localhost:8077/api/v1/dashboard/yearly?years=2023,2024,2025"
```

**Category Breakdown (for pie charts):**
```bash
curl "http://localhost:8077/api/v1/dashboard/categories?year=2025&month=11&type=expense"
```

**Response:**
```json
{
  "category_breakdown": [
    {"category": "ค่าอาหาร", "amount": 5000, "count": 25},
    {"category": "ค่าเดินทาง", "amount": 3000, "count": 15}
  ]
}
```

---

## 🆕 What's New in v3.1

### Latest Features (v3.1.0)
- 🔐 **User Authentication** - Register and login with JWT tokens
- 👤 **User Profiles** - Get authenticated user information
- 🔒 **Secure API** - Password hashing with bcrypt
- 🎫 **JWT Tokens** - 24-hour token expiration

### Previous Features (v3.0)
- 💰 **Budget Management** - Set monthly limits, track spending, get automatic warnings
- 🔄 **Subscription Tracker** - Auto-detect recurring payments from OCR
- 📊 **Dashboard Analytics** - Monthly trends, yearly comparison, category breakdowns
- 🚫 **Duplicate Detection** - Prevent re-uploading the same slip
- 🏷️ **Category System** - Organize expenses and income by category
- 📅 **Thai Date Support** - Recognize "23 พ.ย. 68" format automatically

### API Changes
- **Budget System:** `POST /budgets`, `GET /budgets/status`
- **Subscriptions:** `POST /subscriptions`, auto-detection from OCR
- **Analytics:** `GET /dashboard/monthly`, `/dashboard/yearly`, `/dashboard/categories`
- **Duplicate Check:** Automatic before saving transactions
- **Category Field:** Added to all transaction endpoints

### Auto-Detection
Upload a Netflix/Spotify slip → System automatically creates subscription entry!

Supported services: Netflix, Spotify, YouTube Premium, LINE MAN, Grab, True ID, Disney+, iCloud, Google One, Adobe

### Migration from v2.0
1. Rebuild Docker: `docker-compose up -d --build`
2. Database auto-migrates new tables: `budgets`, `subscriptions`
3. New field `category` added to transactions
4. All existing API endpoints still work

See [CHANGELOG.md](CHANGELOG.md) for full details.

---

## 📝 License

This project is provided as-is for educational and commercial use.

---

**Built with:** Go • Gin • GORM • SQLite • Tesseract OCR

**Version:** 3.1.0
*Last updated: November 27, 2025*

---

## 📈 API Statistics

- **26 Total Endpoints**
- **5 Main Features:** Authentication, Transactions, Budgets, Subscriptions, Analytics
- **10+ Auto-detected Services:** Netflix, Spotify, YouTube, etc.
- **5 Database Tables:** users, transactions, budgets, subscriptions, analytics-ready structure
