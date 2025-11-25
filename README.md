# OCR Payment Slip API

A production-ready RESTful API built with Go for processing Thai bank payment slips using OCR (Optical Character Recognition). Extract transaction details from slip images automatically and store them in a database.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)](https://www.docker.com/)
[![Tesseract](https://img.shields.io/badge/Tesseract-OCR-88C0D0?style=flat)](https://github.com/tesseract-ocr/tesseract)

---

## 📖 Table of Contents

- [Overview](#-overview)
- [Features](#-features)
- [Tech Stack](#-tech-stack)
- [Quick Start with Docker](#-quick-start-with-docker)
- [Manual Installation](#-manual-installation)
- [API Documentation](#-api-documentation)
- [Project Structure](#-project-structure)
- [Configuration](#-configuration)
- [Testing Guide](#-testing-guide)
- [Extending the OCR System](#-extending-the-ocr-system)
- [Troubleshooting](#-troubleshooting)
- [Production Deployment](#-production-deployment)

---

## 🎯 Overview

**OCR Payment Slip API** is a backend service designed to automate the extraction of transaction data from Thai bank payment slips. Simply upload an image of a payment slip, and the API will:

1. **Preprocess** the image for optimal OCR accuracy
2. **Extract text** using Tesseract OCR (Thai + English)
3. **Parse transaction details** using intelligent regex patterns
4. **Store** the data in SQLite database
5. **Return** structured JSON response

### How It Works

```
┌─────────────┐         ┌──────────────┐         ┌───────────────┐
│   Upload    │────────▶│  Image Pre-  │────────▶│   Tesseract   │
│   Image     │         │  processing  │         │   OCR Engine  │
└─────────────┘         └──────────────┘         └───────────────┘
                                                          │
                                                          ▼
┌─────────────┐         ┌──────────────┐         ┌───────────────┐
│   Return    │◀────────│   Save to    │◀────────│  Data Extract │
│   JSON      │         │   Database   │         │  (Regex Parse)│
└─────────────┘         └──────────────┘         └───────────────┘
```

### Supported Banks

The API can detect and process slips from major Thai banks:

| Bank | Code | Thai Name | Detection |
|------|------|-----------|-----------|
| **Siam Commercial Bank** | SCB | ไทยพาณิชย์ | ✅ |
| **Kasikorn Bank** | KBANK | กสิกรไทย | ✅ |
| **Bangkok Bank** | BBL | กรุงเทพ | ✅ |
| **Krung Thai Bank** | KTB | กรุงไทย | ✅ |

---

## ✨ Features

### Core Functionality
- ✅ **Upload payment slip images** (JPG, JPEG, PNG)
- ✅ **OCR text extraction** with Thai + English language support
- ✅ **Image preprocessing** (JPEG conversion, contrast enhancement)
- ✅ **Smart data extraction** using bank-specific regex patterns
- ✅ **Bank detection** from slip content
- ✅ **Transaction type** classification (income/expense)
- ✅ **Automatic file cleanup** (images deleted after processing)

### Extracted Data Fields
- 💰 **Amount** (THB)
- 📅 **Date** (normalized to dd/mm/yyyy)
- ⏰ **Time** (hh:mm:ss)
- 🏦 **Bank name**
- 🔢 **Reference number**
- 👤 **Sender** (when available)
- 👤 **Receiver** (when available)
- 📄 **Raw OCR text** (for debugging)

### API Features
- 🔍 **Query transactions** by type or bank
- 🗑️ **Delete transactions**
- 📊 **RESTful JSON responses**
- 🐳 **Docker support** (recommended)
- 🔄 **CORS enabled** for frontend integration
- ❤️ **Health check** endpoint

---

## 🛠️ Tech Stack

| Component | Technology |
|-----------|-----------|
| **Language** | Go 1.21+ |
| **Web Framework** | Gin |
| **ORM** | GORM |
| **Database** | SQLite |
| **OCR Engine** | Tesseract 4.x |
| **Image Processing** | disintegration/imaging |
| **Tesseract Bindings** | otiai10/gosseract/v2 |
| **Containerization** | Docker + Docker Compose |

---

## 🐳 Quick Start with Docker

**Prerequisites:** Docker and Docker Compose installed

### 1. Start the Application

```bash
# Navigate to project directory
cd D:\Projects\ocr-api

# Start the service (builds automatically)
docker-compose up -d

# Check status
docker-compose ps
```

### 2. Verify It's Running

```bash
# Test health endpoint
curl http://localhost:8077/health

# Expected response:
# {"status":"ok","message":"OCR API is running"}
```

### 3. View Logs

```bash
# Follow logs in real-time
docker-compose logs -f

# View last 100 lines
docker-compose logs --tail=100
```

### 4. Stop the Application

```bash
docker-compose down
```

### Docker Configuration

The `docker-compose.yml` configuration:
- **Port:** Host `8077` → Container `8080`
- **Database:** Persisted in `./data` directory
- **Health Check:** Automatic with 30s interval
- **Auto-restart:** Unless manually stopped

---

## 📦 Manual Installation

If you prefer to run without Docker, follow these steps.

### Prerequisites

#### 1. Install Go 1.21+

Download from [golang.org/dl](https://golang.org/dl/)

```bash
# Verify installation
go version
```

#### 2. Install Tesseract OCR

**Windows:**
```powershell
# Using Chocolatey
choco install tesseract

# Or download installer from:
# https://github.com/UB-Mannheim/tesseract/wiki
```

**macOS:**
```bash
brew install tesseract
brew install tesseract-lang  # Thai + other languages
```

**Linux (Ubuntu/Debian):**
```bash
sudo apt update
sudo apt install tesseract-ocr tesseract-ocr-tha tesseract-ocr-eng
```

**Verify installation:**
```bash
tesseract --version
tesseract --list-langs  # Should show: eng, tha
```

#### 3. Install C Compiler (for CGO)

The Go Tesseract bindings require CGO, which needs a C compiler.

**Windows:**
- Install [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) or [MinGW-w64](https://www.mingw-w64.org/)
- Add to system PATH

**macOS:**
```bash
xcode-select --install
```

**Linux:**
```bash
sudo apt install build-essential
```

### Installation Steps

```bash
# 1. Navigate to project
cd D:\Projects\ocr-api

# 2. Download dependencies
go mod download

# 3. Build the application
go build -o ocr-api.exe

# 4. Run the application
./ocr-api.exe
```

The server starts on `http://localhost:8077`

---

## 📡 API Documentation

Base URL: `http://localhost:8077`

### Endpoints Overview

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/health` | Health check |
| `POST` | `/api/v1/upload` | Upload and process slip |
| `GET` | `/api/v1/transactions` | Get all transactions |
| `GET` | `/api/v1/transactions/:id` | Get transaction by ID |
| `DELETE` | `/api/v1/transactions/:id` | Delete transaction |

---

### 1. Health Check

Check if the API is running.

**Request:**
```bash
GET /health
```

**Example:**
```bash
curl http://localhost:8077/health
```

**Response (200 OK):**
```json
{
  "status": "ok",
  "message": "OCR API is running"
}
```

---

### 2. Upload Payment Slip

Upload a payment slip image for OCR processing.

**Request:**
```bash
POST /api/v1/upload
Content-Type: multipart/form-data
```

**Parameters:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `slip` | File | ✅ Yes | Image file (JPG, JPEG, PNG) |
| `type` | String | ✅ Yes | Transaction type: `income` or `expense` |

**File Requirements:**
- Format: JPG, JPEG, or PNG
- Max size: 10 MB
- Recommended: Clear, well-lit image

**Example - cURL:**
```bash
# Upload as income
curl -X POST http://localhost:8077/api/v1/upload \
  -F "slip=@/path/to/slip-image.jpg" \
  -F "type=income"

# Upload as expense
curl -X POST http://localhost:8077/api/v1/upload \
  -F "slip=@/path/to/slip-image.png" \
  -F "type=expense"
```

**Example - PowerShell:**
```powershell
$form = @{
    slip = Get-Item -Path "C:\Users\YourName\Downloads\slip.jpg"
    type = "income"
}
Invoke-RestMethod -Uri "http://localhost:8077/api/v1/upload" -Method Post -Form $form
```

**Example - JavaScript (Fetch API):**
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

**Success Response (201 Created):**
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
    "raw_ocr_text": "ธนาคารไทยพาณิชย์\nจำนวนเงิน 1,500.00\n...",
    "created_at": "2025-11-25T08:30:00Z",
    "updated_at": "2025-11-25T08:30:00Z"
  }
}
```

**Error Responses:**

**400 Bad Request - Missing type:**
```json
{
  "error": "Missing or invalid 'type' field. Must be 'income' or 'expense'"
}
```

**400 Bad Request - Invalid type:**
```json
{
  "error": "Invalid transaction type. Must be 'income' or 'expense'"
}
```

**400 Bad Request - No file:**
```json
{
  "error": "No file uploaded. Use 'slip' as the form field name"
}
```

**400 Bad Request - File too large:**
```json
{
  "error": "File too large. Maximum size is 10MB"
}
```

**400 Bad Request - Invalid file type:**
```json
{
  "error": "invalid file type: .pdf. Allowed types: jpg, jpeg, png"
}
```

**500 Internal Server Error - OCR failed:**
```json
{
  "error": "Failed to process slip: failed to perform OCR: ..."
}
```

---

### 3. Get All Transactions

Retrieve all transactions with optional filtering.

**Request:**
```bash
GET /api/v1/transactions
```

**Query Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `type` | String | ❌ No | Filter by type: `income` or `expense` |
| `bank` | String | ❌ No | Filter by bank code: `SCB`, `KBANK`, `BBL`, `KTB` |

**Examples:**

```bash
# Get all transactions
curl http://localhost:8077/api/v1/transactions

# Get only income transactions
curl "http://localhost:8077/api/v1/transactions?type=income"

# Get only expense transactions
curl "http://localhost:8077/api/v1/transactions?type=expense"

# Get only SCB transactions
curl "http://localhost:8077/api/v1/transactions?bank=SCB"

# Get SCB income transactions
curl "http://localhost:8077/api/v1/transactions?type=income&bank=SCB"
```

**Success Response (200 OK):**
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
      "raw_ocr_text": "...",
      "created_at": "2025-11-25T08:30:00Z",
      "updated_at": "2025-11-25T08:30:00Z"
    },
    {
      "id": 2,
      "type": "expense",
      "amount": 500.00,
      "date": "24/11/2025",
      "time": "10:15:30",
      "reference": "K987654321098",
      "bank": "KBANK",
      "sender": "นางสมหญิง รักสนุก",
      "receiver": "บริษัท ABC จำกัด",
      "raw_ocr_text": "...",
      "created_at": "2025-11-24T03:15:00Z",
      "updated_at": "2025-11-24T03:15:00Z"
    }
  ]
}
```

**Empty Result:**
```json
{
  "transactions": []
}
```

---

### 4. Get Transaction by ID

Retrieve a specific transaction.

**Request:**
```bash
GET /api/v1/transactions/:id
```

**Path Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | Integer | ✅ Yes | Transaction ID |

**Example:**
```bash
curl http://localhost:8077/api/v1/transactions/1
```

**Success Response (200 OK):**
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
    "created_at": "2025-11-25T08:30:00Z",
    "updated_at": "2025-11-25T08:30:00Z"
  }
}
```

**Error Responses:**

**400 Bad Request - Invalid ID:**
```json
{
  "error": "Invalid transaction ID"
}
```

**404 Not Found:**
```json
{
  "error": "Transaction not found"
}
```

---

### 5. Delete Transaction

Delete a transaction from the database.

**Request:**
```bash
DELETE /api/v1/transactions/:id
```

**Path Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | Integer | ✅ Yes | Transaction ID |

**Example:**
```bash
curl -X DELETE http://localhost:8077/api/v1/transactions/1
```

**Success Response (200 OK):**
```json
{
  "message": "Transaction deleted successfully"
}
```

**Error Responses:**

**400 Bad Request:**
```json
{
  "error": "Invalid transaction ID"
}
```

**404 Not Found:**
```json
{
  "error": "Transaction not found"
}
```

---

## 🏗️ Project Structure

```
ocr-api/
├── main.go                          # Application entry point & CORS middleware
├── go.mod                           # Go module dependencies
├── go.sum                           # Dependency checksums
├── Dockerfile                       # Docker image configuration
├── docker-compose.yml               # Docker Compose orchestration
├── .env.example                     # Environment variables template
├── Makefile                         # Build automation (if exists)
│
├── config/
│   ├── config.go                    # App configuration & env loading
│   └── database.go                  # Database connection & migrations
│
├── models/
│   └── transaction.go               # Transaction GORM model
│
├── controllers/
│   ├── upload_controller.go         # Handles /upload endpoint
│   └── transaction_controller.go    # Handles transaction CRUD
│
├── services/
│   ├── ocr_service.go              # Orchestrates OCR workflow
│   └── transaction_service.go      # Transaction business logic
│
├── ocr/
│   ├── tesseract.go                # Tesseract OCR wrapper
│   ├── preprocessor.go             # Image preprocessing (JPEG, contrast)
│   └── extractor.go                # Data extraction with bank patterns
│
├── routes/
│   └── routes.go                   # API route definitions
│
├── utils/
│   └── helpers.go                  # Utility functions
│
├── uploads/                        # Temporary upload directory (auto-created)
├── data/
│   └── db.sqlite                   # SQLite database (auto-created)
└── README.md                       # This file
```

### Key Files Explained

| File | Purpose |
|------|---------|
| `main.go` | Initializes server, database, routes, CORS |
| `config/config.go` | Loads environment variables, sets defaults |
| `ocr/extractor.go` | Contains bank-specific regex patterns |
| `ocr/preprocessor.go` | Enhances image quality for better OCR |
| `models/transaction.go` | Database schema definition |
| `docker-compose.yml` | Production-ready container setup |

---

## ⚙️ Configuration

### Environment Variables

Create a `.env` file or set environment variables:

```bash
# Server Configuration
SERVER_PORT=8077              # Port to listen on

# Database
DATABASE_PATH=./data/db.sqlite  # SQLite database file path

# Upload Settings
UPLOAD_DIR=./uploads            # Temporary upload directory
MAX_UPLOAD_SIZE=10485760        # Max file size (10MB in bytes)

# OCR Settings
TESSERACT_LANG=tha+eng         # Tesseract languages (Thai + English)

# Gin Mode
GIN_MODE=release               # release or debug
```

### Docker Environment

The `docker-compose.yml` sets these automatically:
- Port mapping: `8077:8080` (host:container)
- Persistent database in `./data`
- Thai + English OCR support
- Auto-restart enabled

---

## 🧪 Testing Guide

### 1. Quick Test with cURL

```bash
# Health check
curl http://localhost:8077/health

# Upload test slip (replace path)
curl -X POST http://localhost:8077/api/v1/upload \
  -F "slip=@test-slip.jpg" \
  -F "type=income"

# View all transactions
curl http://localhost:8077/api/v1/transactions

# View specific transaction
curl http://localhost:8077/api/v1/transactions/1

# Delete transaction
curl -X DELETE http://localhost:8077/api/v1/transactions/1
```

### 2. Testing with Postman

**Collection Setup:**

1. **POST** `/api/v1/upload`
   - Body → `form-data`
   - Key: `slip` (File) → Select image file
   - Key: `type` (Text) → `income` or `expense`

2. **GET** `/api/v1/transactions`
   - Params: `type=income` (optional)
   - Params: `bank=SCB` (optional)

3. **GET** `/api/v1/transactions/:id`
   - Params: `id=1`

4. **DELETE** `/api/v1/transactions/:id`
   - Params: `id=1`

### 3. Prepare Test Images

**Best Practices for Test Slips:**
- ✅ Clear, well-lit photos
- ✅ High contrast
- ✅ Minimal shadows
- ✅ Straight orientation
- ❌ Avoid blurry images
- ❌ Avoid low resolution

### 4. View Raw OCR Output

To debug extraction issues, check the `raw_ocr_text` field:

```bash
curl http://localhost:8077/api/v1/transactions/1 | jq '.transaction.raw_ocr_text'
```

This shows exactly what Tesseract extracted.

---

## 🔧 Extending the OCR System

### Adding a New Bank

Edit `ocr/extractor.go` and add a new bank pattern:

```go
BankPattern{
    Name: "TMB",  // Bank code
    Identifiers: []string{
        `(?i)tmb`,
        `(?i)ทหารไทยธนชาต`,
        `(?i)ธนาคารทหารไทย`,
    },
    AmountPatterns: []string{
        `(?i)(?:amount|จำนวนเงิน)[:\s]*([0-9,]+\.?\d{0,2})`,
        `(?i)รับเงิน[:\s]*([0-9,]+\.?\d{0,2})`,
    },
    DatePatterns: []string{
        `(\d{1,2}[/-]\d{1,2}[/-]\d{2,4})`,
    },
    TimePatterns: []string{
        `(\d{1,2}:\d{2}(?::\d{2})?)`,
    },
    ReferencePatterns: []string{
        `(?i)(?:ref|reference|อ้างอิง)[:\s]*([A-Z0-9]+)`,
    },
}
```

### Improving OCR Accuracy

**1. Enhance Image Preprocessing**

Edit `ocr/preprocessor.go`:
```go
// Add more preprocessing steps
func PreprocessImage(inputPath string) (string, error) {
    img := ... // Load image

    // Increase contrast (current)
    img = imaging.AdjustContrast(img, 30)

    // Add: Increase sharpness
    img = imaging.Sharpen(img, 2.0)

    // Add: Convert to grayscale
    img = imaging.Grayscale(img)

    // Add: Adjust brightness
    img = imaging.AdjustBrightness(img, 10)

    return outputPath, nil
}
```

**2. Add More Regex Patterns**

Study the `raw_ocr_text` from failed extractions, then add patterns:

```go
// Example: Add more amount pattern variations
AmountPatterns: []string{
    `(?i)(?:amount|จำนวนเงิน)[:\s]*([0-9,]+\.?\d{0,2})`,
    `(?i)รับเงิน[:\s]*([0-9,]+\.?\d{0,2})`,
    `(?i)โอน[:\s]*([0-9,]+\.?\d{0,2})`,  // New pattern
    `บาท[:\s]*([0-9,]+\.?\d{0,2})`,      // New pattern
},
```

**3. Train Custom Tesseract Data**

For very specific slip formats, train custom Tesseract LSTM models.

### Adding New Extracted Fields

**1. Update Model** (`models/transaction.go`):
```go
type Transaction struct {
    // Existing fields...
    Channel    string `json:"channel" gorm:"type:varchar(50)"`  // New field
}
```

**2. Update Extractor** (`ocr/extractor.go`):
```go
type ExtractedData struct {
    // Existing fields...
    Channel    string  // New field
}

// Add pattern in BankPattern
ChannelPatterns: []string{
    `(?i)channel[:\s]*([A-Z0-9]+)`,
},
```

**3. Migrate Database:**
```bash
# Delete old database (for testing)
rm ./data/db.sqlite

# Restart app to create new schema
docker-compose restart
```

---

## 🐛 Troubleshooting

### Problem: Tesseract Not Found

**Error:**
```
Failed to perform OCR: exec: "tesseract": executable file not found
```

**Solution:**
1. Install Tesseract (see [Manual Installation](#-manual-installation))
2. Add to system PATH
3. Restart terminal/IDE
4. Verify: `tesseract --version`

---

### Problem: Thai Language Not Available

**Error:**
```
Failed to set language: Error opening data file
```

**Solution:**
```bash
# Check available languages
tesseract --list-langs

# If 'tha' is missing:
# Windows (Reinstall Tesseract with Thai option)
# macOS
brew install tesseract-lang

# Linux
sudo apt install tesseract-ocr-tha
```

---

### Problem: CGO Build Errors

**Error:**
```
gcc: command not found
undefined: gosseract.NewClient
```

**Solution:**
1. **Use Docker** (recommended - no CGO issues)
2. Or install C compiler:
   - Windows: Install [TDM-GCC](https://jmeubank.github.io/tdm-gcc/)
   - macOS: `xcode-select --install`
   - Linux: `sudo apt install build-essential`

---

### Problem: Port Already in Use

**Error:**
```
bind: address already in use
```

**Solution:**
```bash
# Windows
netstat -ano | findstr :8077
taskkill /PID <PID> /F

# Linux/macOS
lsof -ti:8077 | xargs kill -9

# Or change port in docker-compose.yml
ports:
  - "8078:8080"  # Use different host port
```

---

### Problem: Poor OCR Accuracy

**Symptoms:**
- Wrong amounts extracted
- Missing fields
- Incorrect dates

**Solutions:**
1. **Improve image quality:**
   - Use better lighting
   - Take photo straight-on
   - Ensure high resolution

2. **Check raw OCR output:**
   ```bash
   curl http://localhost:8077/api/v1/transactions/1 | jq '.transaction.raw_ocr_text'
   ```

3. **Enhance preprocessing** (see [Extending the OCR System](#-extending-the-ocr-system))

4. **Add bank-specific patterns** in `ocr/extractor.go`

---

### Problem: Database Locked

**Error:**
```
database is locked
```

**Solution:**
```bash
# Stop all processes
docker-compose down

# Remove database lock
rm ./data/db.sqlite-shm
rm ./data/db.sqlite-wal

# Restart
docker-compose up -d
```

---

## 🚀 Production Deployment

### Security Checklist

- [ ] **Add authentication** (JWT, API keys)
- [ ] **Enable HTTPS** (TLS certificates)
- [ ] **Add rate limiting** (prevent abuse)
- [ ] **Validate file uploads** (check magic bytes, not just extensions)
- [ ] **Set up CORS properly** (don't use `*` in production)
- [ ] **Use environment secrets** (don't commit `.env`)
- [ ] **Add request logging** (structured logs)
- [ ] **Set file size limits** (already 10MB, adjust if needed)

### Database

**For production, consider:**
- PostgreSQL (concurrent writes)
- MySQL (better performance)
- Add database backups

```bash
# Example: Backup SQLite
cp ./data/db.sqlite ./data/backup-$(date +%Y%m%d).sqlite
```

### Monitoring

Add health checks and metrics:
- Uptime monitoring (UptimeRobot, Pingdom)
- Error tracking (Sentry)
- Performance monitoring (New Relic, DataDog)

### Scaling

**Horizontal Scaling:**
```yaml
# docker-compose.yml
services:
  ocr-api:
    deploy:
      replicas: 3  # Run 3 instances
```

**Add Load Balancer:**
- Nginx
- Traefik
- Cloud load balancers (AWS ALB, GCP LB)

### CI/CD Pipeline

Example GitHub Actions:

```yaml
name: Deploy
on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Build and push Docker image
        run: |
          docker build -t ocr-api:latest .
          docker push yourregistry/ocr-api:latest
```

---

## 📊 Performance Considerations

### Processing Times

Typical processing times per slip:

| Step | Duration |
|------|----------|
| Upload & Validation | ~50ms |
| JPEG Conversion | ~200ms |
| Image Preprocessing | ~300ms |
| OCR Text Extraction | ~1-3s |
| Data Extraction | ~10ms |
| Database Save | ~20ms |
| **Total** | **~2-4s** |

### Optimization Tips

1. **Use Docker** (dependencies pre-configured)
2. **SSD storage** (faster image I/O)
3. **Increase Tesseract threads** (for multiple uploads)
4. **Cache preprocessed images** (if reprocessing)
5. **Use database connection pooling** (if switching to PostgreSQL)

---

## 🤝 Contributing

### Improving OCR Accuracy

1. Collect sample slips (with permission)
2. Process and check `raw_ocr_text`
3. Identify patterns in the text
4. Add regex patterns to `ocr/extractor.go`
5. Test with various slip formats
6. Submit improvements

### Reporting Issues

When reporting OCR issues, include:
- Slip bank name
- Screenshot (redact sensitive info)
- Raw OCR text output
- Expected vs actual extraction

---

## 📝 License

This project is provided as-is for educational and commercial use.

---

## 🙏 Acknowledgments

**Built with:**
- [Go](https://golang.org/) - Programming language
- [Gin](https://gin-gonic.com/) - Web framework
- [GORM](https://gorm.io/) - ORM library
- [Tesseract](https://github.com/tesseract-ocr/tesseract) - OCR engine
- [gosseract](https://github.com/otiai10/gosseract) - Go Tesseract bindings
- [imaging](https://github.com/disintegration/imaging) - Image processing

---

## 📞 Support

For issues and questions:
- Check [Troubleshooting](#-troubleshooting) section
- Review [API Documentation](#-api-documentation)
- Check Tesseract documentation for OCR issues

---

**Made with ❤️ for automating payment slip processing**

*Last updated: November 25, 2025*
