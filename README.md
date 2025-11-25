# OCR Payment Slip API

A RESTful API built with Go for processing Thai bank payment slips using OCR (Optical Character Recognition). Extract transaction details from slip images automatically.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)](https://www.docker.com/)

---

## 🎯 Overview

Upload a payment slip image → Extract transaction data (amount, date, bank, reference, etc.) → Store in database

**Supported Banks:** SCB (ไทยพาณิชย์), KBANK (กสิกรไทย), BBL (กรุงเทพ), KTB (กรุงไทย)

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
| `POST` | `/api/v1/upload` | Upload and process slip |
| `GET` | `/api/v1/transactions` | Get all transactions |
| `GET` | `/api/v1/transactions/:id` | Get transaction by ID |
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

### 2. Upload Payment Slip

```bash
POST /api/v1/upload
Content-Type: multipart/form-data
```

**Parameters:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `slip` | File | Yes | Image file (JPG, JPEG, PNG, max 10MB) |
| `type` | String | Yes | `income` or `expense` |

**Example - cURL:**
```bash
curl -X POST http://localhost:8077/api/v1/upload \
  -F "slip=@/path/to/slip.jpg" \
  -F "type=income"
```

**Example - PowerShell:**
```powershell
$form = @{
    slip = Get-Item -Path "C:\path\to\slip.jpg"
    type = "income"
}
Invoke-RestMethod -Uri "http://localhost:8077/api/v1/upload" -Method Post -Form $form
```

**Example - JavaScript:**
```javascript
const formData = new FormData();
formData.append('slip', fileInput.files[0]);
formData.append('type', 'income');

fetch('http://localhost:8077/api/v1/upload', {
  method: 'POST',
  body: formData
})
  .then(res => res.json())
  .then(data => console.log(data));
```

**Success Response (201):**
```json
{
  "message": "Slip processed successfully",
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
    "raw_ocr_text": "...",
    "created_at": "2025-11-25T08:30:00Z"
  }
}
```

**Error Responses:**
```json
// 400 - Missing/invalid type
{ "error": "Missing or invalid 'type' field. Must be 'income' or 'expense'" }

// 400 - No file
{ "error": "No file uploaded. Use 'slip' as the form field name" }

// 400 - File too large
{ "error": "File too large. Maximum size is 10MB" }

// 400 - Invalid file type
{ "error": "invalid file type: .pdf. Allowed types: jpg, jpeg, png" }
```

---

### 3. Get All Transactions

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
      "created_at": "2025-11-25T08:30:00Z"
    }
  ]
}
```

---

### 4. Get Transaction by ID

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

### 5. Delete Transaction

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

# Upload slip
curl -X POST http://localhost:8077/api/v1/upload \
  -F "slip=@test.jpg" \
  -F "type=income"

# View transactions
curl http://localhost:8077/api/v1/transactions

# View specific transaction
curl http://localhost:8077/api/v1/transactions/1

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

## 📝 License

SNDLK01
---

**Built with:** Go • Gin • GORM • Tesseract OCR

*Last updated: November 25, 2025*
