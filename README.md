# OCR Payment Slip API

A RESTful API built with Go for processing Thai bank payment slips using OCR (Optical Character Recognition). Extract transaction details from slip images automatically.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)](https://www.docker.com/)

---

## 🎯 Overview

A complete transaction management API with OCR capabilities:

1. **Upload slip(s)** → Extract transaction data via OCR → Store in database
2. **Create manual transactions** → For income/transfers without slips → Store in database
3. **Update transactions** → Add details, correct amounts, categorize expenses
4. **Query & manage** → Filter, view, and delete transactions

**Features:**
- 📸 Multi-file upload support (process multiple slips at once)
- 🤖 OCR for Thai bank slips (SCB, KBANK, BBL, KTB)
- ✏️ Manual transaction creation (income/expense)
- 🔄 Update transaction details anytime
- 🔍 Filter by type (income/expense) and bank
- 📝 Add custom descriptions/notes

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

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/health` | Health check |
| `POST` | `/api/v1/upload` | Upload and process slip(s) - **supports multiple files** |
| `POST` | `/api/v1/transactions` | **Create transaction manually** (income/expense) |
| `GET` | `/api/v1/transactions` | Get all transactions (filter by type/bank) |
| `GET` | `/api/v1/transactions/:id` | Get transaction by ID |
| `PUT/PATCH` | `/api/v1/transactions/:id` | **Update transaction details** |
| `DELETE` | `/api/v1/transactions/:id` | Delete transaction |

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

### 2. Upload Payment Slip(s)

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

### 3. Create Transaction Manually

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

### 4. Get All Transactions

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

### 5. Get Transaction by ID

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

### 6. Update Transaction

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

### 7. Delete Transaction

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
├── main.go                      # Entry point
├── docker-compose.yml           # Docker setup
├── config/
│   ├── config.go               # Configuration
│   └── database.go             # Database setup
├── models/
│   └── transaction.go          # Transaction model
├── controllers/
│   ├── upload_controller.go    # Upload handler
│   └── transaction_controller.go
├── services/
│   ├── ocr_service.go          # OCR workflow
│   └── transaction_service.go
├── ocr/
│   ├── tesseract.go            # Tesseract wrapper
│   ├── preprocessor.go         # Image preprocessing
│   └── extractor.go            # Data extraction (regex)
├── routes/
│   └── routes.go               # API routes
└── utils/
    └── helpers.go              # Utilities
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

```bash
# Health check
curl http://localhost:8077/health

# Upload single slip
curl -X POST http://localhost:8077/api/v1/upload \
  -F "slip=@test.jpg" \
  -F "type=expense"

# Upload multiple slips
curl -X POST http://localhost:8077/api/v1/upload \
  -F "slips=@slip1.jpg" \
  -F "slips=@slip2.jpg" \
  -F "slips=@slip3.jpg" \
  -F "type=expense"

# Create manual transaction (income without slip)
curl -X POST http://localhost:8077/api/v1/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "type": "income",
    "amount": 5000,
    "date": "26/11/2025",
    "detail": "รายได้จากการขาย"
  }'

# View all transactions
curl http://localhost:8077/api/v1/transactions

# Filter transactions
curl "http://localhost:8077/api/v1/transactions?type=income"
curl "http://localhost:8077/api/v1/transactions?bank=SCB"

# View specific transaction
curl http://localhost:8077/api/v1/transactions/1

# Update transaction details
curl -X PATCH http://localhost:8077/api/v1/transactions/1 \
  -H "Content-Type: application/json" \
  -d '{"detail": "ค่าน้ำมัน"}'

# Delete transaction
curl -X DELETE http://localhost:8077/api/v1/transactions/1
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

---

## 🆕 What's New in v2.0

### Major Features
- ✅ **Multi-file upload** - Upload multiple slips in one request
- ✅ **Manual transaction creation** - Add income/expense without slips
- ✅ **Transaction updates** - Edit details, amounts, dates anytime
- ✅ **Detail field** - Add custom descriptions/notes to transactions
- ✅ **Improved OCR text** - Cleaner, more readable `raw_ocr_text`

### API Changes
- `POST /api/v1/upload` now supports `slips` (multiple) and `slip` (single)
- `POST /api/v1/transactions` - Create manual transactions
- `PUT/PATCH /api/v1/transactions/:id` - Update transactions
- All transactions now include `detail` field in responses

### Migration from v1.0
1. Rebuild Docker: `docker-compose up -d --build`
2. Database auto-migrates the new `detail` field
3. Old API calls still work (backward compatible)

See [CHANGELOG.md](CHANGELOG.md) for full details.

---

## 📝 License

This project is provided as-is for educational and commercial use.

---

**Built with:** Go • Gin • GORM • Tesseract OCR

*Last updated: November 26, 2025 (v2.0)*
