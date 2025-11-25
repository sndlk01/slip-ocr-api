# OCR Payment Slip API

A production-ready Golang backend API for processing Thai bank payment slips using OCR (Optical Character Recognition). The API extracts transaction details from slip images and stores them in a database.

## 🚀 Features

- **Upload payment slip images** (JPG, PNG)
- **OCR processing** using Tesseract (Thai + English)
- **Image preprocessing** for better OCR accuracy
- **Automatic data extraction** using regex patterns:
  - Amount
  - Date (dd/mm/yyyy format)
  - Time
  - Reference number
  - Sender/Receiver (when available)
- **Bank detection** (SCB, KBank, BBL, KTB)
- **Transaction management** (Create, Read, Delete)
- **RESTful API** with JSON responses
- **SQLite database** with GORM ORM
- **Automatic file cleanup** (images are not stored)

## 🐳 Quick Start with Docker (Recommended)

The easiest way to run this application is using Docker:

```bash
# Clone or navigate to the project
cd d:\Projects\ocr-api

# Start the application
docker-compose up -d

# Check if it's running
curl http://localhost:8080/health

# View logs
docker-compose logs -f
```

That's it! The API is now running at `http://localhost:8080`

**For detailed Docker instructions, see [DOCKER.md](DOCKER.md)**

---

## 📋 Manual Installation Prerequisites

If you prefer to run without Docker, follow these steps:

### 1. Install Go

Download and install Go 1.21 or higher from [https://golang.org/dl/](https://golang.org/dl/)

Verify installation:
```bash
go version
```

### 2. Install Tesseract OCR

#### **Windows**

1. Download the installer from [UB Mannheim Tesseract](https://github.com/UB-Mannheim/tesseract/wiki)
2. Run the installer (e.g., `tesseract-ocr-w64-setup-5.3.3.20231005.exe`)
3. During installation, make sure to select **Thai language data**
4. Add Tesseract to PATH:
   - Default installation path: `C:\Program Files\Tesseract-OCR`
   - Add to System Environment Variables → Path

Verify installation:
```powershell
tesseract --version
```

#### **macOS**

Using Homebrew:
```bash
brew install tesseract
brew install tesseract-lang  # For Thai language support
```

Verify installation:
```bash
tesseract --version
tesseract --list-langs  # Should show 'tha' and 'eng'
```

#### **Linux (Ubuntu/Debian)**

```bash
sudo apt update
sudo apt install tesseract-ocr
sudo apt install tesseract-ocr-tha  # Thai language data
sudo apt install tesseract-ocr-eng  # English language data
```

Verify installation:
```bash
tesseract --version
tesseract --list-langs
```

### 3. Install GCC (for CGO)

Tesseract Go bindings require CGO, which needs a C compiler.

#### **Windows**
- Install [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) or [MinGW-w64](https://www.mingw-w64.org/)
- Add to PATH

#### **macOS**
```bash
xcode-select --install
```

#### **Linux**
```bash
sudo apt install build-essential
```

## 🛠️ Installation

### 1. Clone or navigate to the project directory

```bash
cd d:\Projects\ocr-api
```

### 2. Install Go dependencies

```bash
go mod download
```

### 3. Build the application

```bash
go build -o ocr-api.exe
```

## 🏃 Running the Application

### Development Mode

```bash
go run main.go
```

### Production Mode

```bash
# Build first
go build -o ocr-api.exe

# Run the executable
./ocr-api.exe
```

The server will start on `http://localhost:8080`

## 📡 API Endpoints

### 1. Health Check

```bash
GET /health
```

**Response:**
```json
{
  "status": "ok",
  "message": "OCR API is running"
}
```

### 2. Upload Payment Slip

```bash
POST /api/v1/upload
```

**Parameters:**
- `slip` (file, required): Image file (JPG, PNG)
- `type` (form field, required): Transaction type (`income` or `expense`)

**Example using curl:**

```bash
# Upload as income
curl -X POST http://localhost:8080/api/v1/upload \
  -F "slip=@payment_slip.jpg" \
  -F "type=income"

# Upload as expense
curl -X POST http://localhost:8080/api/v1/upload \
  -F "slip=@payment_slip.png" \
  -F "type=expense"
```

**Example using PowerShell:**

```powershell
$form = @{
    slip = Get-Item -Path "payment_slip.jpg"
    type = "income"
}
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/upload" -Method Post -Form $form
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
    "time": "14:30:00",
    "reference": "ABC123456",
    "bank": "SCB",
    "sender": "John Doe",
    "receiver": "Jane Smith",
    "created_at": "2025-11-25T15:30:00Z",
    "updated_at": "2025-11-25T15:30:00Z"
  }
}
```

### 3. Get All Transactions

```bash
GET /api/v1/transactions
```

**Optional Query Parameters:**
- `type`: Filter by type (`income` or `expense`)
- `bank`: Filter by bank (`SCB`, `KBank`, `BBL`, `KTB`)

**Examples:**

```bash
# Get all transactions
curl http://localhost:8080/api/v1/transactions

# Get only income transactions
curl http://localhost:8080/api/v1/transactions?type=income

# Get only SCB transactions
curl http://localhost:8080/api/v1/transactions?bank=SCB
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
      "time": "14:30:00",
      "reference": "ABC123456",
      "bank": "SCB",
      "created_at": "2025-11-25T15:30:00Z",
      "updated_at": "2025-11-25T15:30:00Z"
    }
  ]
}
```

### 4. Get Transaction by ID

```bash
GET /api/v1/transactions/:id
```

**Example:**

```bash
curl http://localhost:8080/api/v1/transactions/1
```

**Response (200):**
```json
{
  "transaction": {
    "id": 1,
    "type": "income",
    "amount": 1500.00,
    "date": "25/11/2025",
    "reference": "ABC123456",
    "bank": "SCB"
  }
}
```

### 5. Delete Transaction

```bash
DELETE /api/v1/transactions/:id
```

**Example:**

```bash
curl -X DELETE http://localhost:8080/api/v1/transactions/1
```

**Response (200):**
```json
{
  "message": "Transaction deleted successfully"
}
```

## 🏗️ Project Structure

```
ocr-api/
├── main.go                 # Application entry point
├── go.mod                  # Go module dependencies
├── config/
│   ├── config.go          # Configuration management
│   └── database.go        # Database initialization
├── models/
│   └── transaction.go     # Transaction model
├── controllers/
│   ├── upload_controller.go      # Upload endpoint handler
│   └── transaction_controller.go # Transaction CRUD handlers
├── services/
│   ├── ocr_service.go            # OCR processing service
│   └── transaction_service.go    # Transaction business logic
├── ocr/
│   ├── tesseract.go       # Tesseract OCR wrapper
│   ├── preprocessor.go    # Image preprocessing
│   └── extractor.go       # Data extraction with regex
├── routes/
│   └── routes.go          # API route definitions
├── utils/
│   └── helpers.go         # Utility functions
├── uploads/               # Temporary upload directory (auto-created)
└── db.sqlite             # SQLite database (auto-created)
```

## ⚙️ Configuration

Environment variables (optional):

```bash
# Server configuration
SERVER_PORT=8080

# Database
DATABASE_PATH=./db.sqlite

# Upload settings
UPLOAD_DIR=./uploads

# Tesseract language
TESSERACT_LANG=tha+eng
```

## 🔧 Extending OCR Mapping

### Adding a New Bank

Edit `ocr/extractor.go` and add a new `BankPattern`:

```go
{
    Name: "NewBank",
    Identifiers: []string{
        "(?i)new\\s*bank",
        "(?i)ธนาคารใหม่",
    },
    AmountPatterns: []string{
        `(?i)(?:amount|จำนวนเงิน)[:\s]*([0-9,]+\.?\d{0,2})`,
    },
    DatePatterns: []string{
        `(\d{1,2}[/-]\d{1,2}[/-]\d{2,4})`,
    },
    // Add more patterns...
}
```

### Improving Regex Patterns

1. **Test with real slips**: Collect sample OCR output
2. **Identify patterns**: Look for common text structures
3. **Add patterns**: Add new regex to the appropriate pattern array
4. **Test**: Upload slips and verify extraction

### Debugging OCR Output

The raw OCR text is stored in the `raw_ocr_text` field of each transaction. Use this to:

1. View what Tesseract extracted
2. Identify missing patterns
3. Refine regex patterns

## 🐛 Troubleshooting

### Tesseract not found

**Error:** `Failed to perform OCR: exec: "tesseract": executable file not found`

**Solution:** 
- Ensure Tesseract is installed
- Add Tesseract to system PATH
- Restart terminal/IDE

### Thai language not available

**Error:** `Failed to set language: Error opening data file`

**Solution:**
- Install Thai language data for Tesseract
- Verify with: `tesseract --list-langs`

### CGO errors

**Error:** `gcc: command not found`

**Solution:**
- Install GCC/MinGW (see Prerequisites)
- Ensure it's in PATH

### Image preprocessing fails

**Error:** `failed to preprocess image`

**Solution:**
- Check image file is not corrupted
- Ensure image format is supported (JPG, PNG)
- Check file permissions

## 📝 Testing

### Manual Testing

1. Prepare test images of Thai bank payment slips
2. Use curl or Postman to upload
3. Check response and database

### Example Test Commands

```bash
# Test upload
curl -X POST http://localhost:8080/api/v1/upload \
  -F "slip=@test_slip.jpg" \
  -F "type=income"

# Test get all
curl http://localhost:8080/api/v1/transactions

# Test get by ID
curl http://localhost:8080/api/v1/transactions/1

# Test delete
curl -X DELETE http://localhost:8080/api/v1/transactions/1
```

## 🔒 Production Considerations

1. **Add authentication**: Implement JWT or API key authentication
2. **Rate limiting**: Add middleware to prevent abuse
3. **Input validation**: Enhance file validation (magic bytes, size limits)
4. **Error logging**: Implement structured logging (e.g., logrus, zap)
5. **Database**: Consider PostgreSQL or MySQL for production
6. **File storage**: If needed, use cloud storage (S3, GCS)
7. **Monitoring**: Add health checks and metrics
8. **HTTPS**: Use TLS certificates in production
9. **Environment config**: Use proper environment variable management

## 📄 License

This project is provided as-is for educational and commercial use.

## 🤝 Contributing

To improve OCR accuracy:

1. Collect more sample slips
2. Refine regex patterns
3. Add preprocessing techniques
4. Test with different image qualities

---

**Built with ❤️ using Go, Gin, GORM, and Tesseract**
