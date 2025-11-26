# Changelog

All notable changes to this project will be documented in this file.

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
  - Optional fields: `date`, `time`, `reference`, `bank`, `sender`, `receiver`, `detail`

#### Transaction Update Support
- **New PUT/PATCH `/api/v1/transactions/:id` endpoint** for updating transactions
  - Supports partial updates (send only fields you want to change)
  - Can update: `amount`, `date`, `time`, `reference`, `bank`, `sender`, `receiver`, `detail`
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
