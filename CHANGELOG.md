# Changelog

All notable changes to this project will be documented in this file.

## [3.0.0] - 2025-11-26

### Added - Major Features

#### 1. Budget Management System
- **NEW: Budget tracking and monitoring**
  - Set monthly budgets per category
  - Real-time budget status tracking
  - Automatic warning system (ok/warning/exceeded)
  - Track spending vs. budget limits
- **Endpoints:**
  - `POST /api/v1/budgets` - Create budget
  - `GET /api/v1/budgets` - List all budgets
  - `GET /api/v1/budgets/status` - Check budget status with spending details
  - `DELETE /api/v1/budgets/:id` - Remove budget

#### 2. Subscription Tracker
- **NEW: Recurring payment management**
  - Track monthly/yearly subscriptions
  - Auto-detect popular services from OCR (Netflix, Spotify, YouTube Premium, etc.)
  - Monitor active subscriptions and billing dates
  - Calculate total monthly subscription costs
- **Auto-Detection Support:** Netflix, Spotify, YouTube Premium, LINE MAN, Grab Unlimited, True ID, Disney+, iCloud, Google One, Adobe
- **Endpoints:**
  - `POST /api/v1/subscriptions` - Add subscription
  - `GET /api/v1/subscriptions` - List subscriptions
  - `DELETE /api/v1/subscriptions/:id` - Remove subscription

#### 3. Dashboard & Analytics API
- **NEW: Data visualization endpoints**
  - Monthly trends (12-month view)
  - Yearly comparison (multi-year analysis)
  - Category breakdown (pie chart data)
- **Endpoints:**
  - `GET /api/v1/dashboard/monthly?year=2025` - Get 12-month income/expense trend
  - `GET /api/v1/dashboard/yearly?years=2023,2024,2025` - Compare multiple years
  - `GET /api/v1/dashboard/categories?year=2025&month=11&type=expense` - Category breakdown for charts

#### 4. Duplicate Slip Detection
- **NEW: Automatic duplicate prevention**
  - Checks reference numbers before saving
  - Validates amount + date + time + bank combination
  - Prevents accidental re-uploads
  - Returns helpful error messages with existing transaction ID

#### 5. Category Support
- **Added `category` field to Transaction model**
  - Categorize income and expenses
  - Filter and group by category
  - Essential for budget tracking and analytics

### Improved

#### Enhanced Thai Date Recognition
- **Better OCR date extraction**
  - Now supports Thai month abbreviations (พ.ย., ธ.ค., etc.)
  - Automatic Buddhist-to-Christian year conversion (68 → 2025)
  - Handles both "23 พ.ย. 68" and "23/11/2025" formats
  - Improved pattern matching for "จํานวน" and "เลขที่รายการ"

#### Monthly Summary Enhancement
- **Expanded summary endpoint**
  - Now includes category breakdown
  - Shows transaction counts per category
  - Separates income/expense categories
  - Better data structure for frontend display

---

## [2.0.0] - 2025-11-26

### Added

#### Multiple File Upload Support
- **Upload API now supports multiple slips in a single request**
  - Changed form field from `slip` (single) to `slips` (multiple)
  - Backward compatible: still supports `slip` for single file upload
  - Returns detailed response with `success_count`, `total_count`, and `errors` array
  - Each file is processed independently and creates separate transactions

#### Manual Transaction Creation
- **New POST `/api/v1/transactions` endpoint** for creating transactions without slips
  - Useful for income transactions that don't have physical slips
  - Required fields: `type` (income/expense), `amount`
  - Optional fields: `date`, `time`, `reference`, `bank`, `sender`, `receiver`, `detail`, `category`

#### Transaction Update Support
- **New PUT/PATCH `/api/v1/transactions/:id` endpoint** for updating transactions
  - Supports partial updates (send only fields you want to change)
  - Can update: `amount`, `date`, `time`, `reference`, `bank`, `sender`, `receiver`, `category`, `detail`
  - Cannot update: `type` (income/expense must remain unchanged)

#### Detail Field
- **Added `detail` field to Transaction model**
  - Allows users to add custom descriptions/notes to transactions
  - Always visible in JSON response (not omitted when empty)
  - Can be set during creation or updated later

### Improved

#### OCR Text Readability
- **Enhanced `raw_ocr_text` formatting**
  - Automatically removes excessive whitespace
  - Trims leading/trailing spaces from each line
  - Removes empty lines
  - Replaces multiple spaces with single space
  - Much more readable than before

### API Changes Summary

| Change Type | Old Behavior | New Behavior |
|------------|--------------|--------------|
| Upload field name | `slip` only | `slip` or `slips` (both supported) |
| Upload count | Single file only | Multiple files supported |
| Manual transaction | Not possible | POST `/api/v1/transactions` |
| Update transaction | Not possible | PUT/PATCH `/api/v1/transactions/:id` |
| Transaction fields | No `detail` field | Added `detail` field |
| OCR text format | Raw with whitespace | Cleaned and formatted |

### Migration Notes

1. **Database**: Run the application to auto-migrate the new `detail` field
2. **API Clients**: Update to use `slips` field name for multiple files (or keep using `slip` for backward compatibility)
3. **Docker**: Rebuild image with `docker-compose up -d --build`

### Files Modified

- `models/transaction.go` - Added `detail` field
- `controllers/upload_controller.go` - Multi-file upload support
- `controllers/transaction_controller.go` - Added Create and Update methods
- `services/transaction_service.go` - Added Update method
- `services/ocr_service.go` - Clean OCR text before saving
- `utils/helpers.go` - Added CleanOCRText function
- `routes/routes.go` - Added new endpoints

---

## [1.0.0] - 2025-11-25

### Initial Release

- OCR processing for Thai bank payment slips
- Support for SCB, KBANK, BBL, KTB banks
- Basic CRUD operations for transactions
- SQLite database storage
- Docker support
- Tesseract OCR integration
