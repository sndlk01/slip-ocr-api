package ocr

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type ExtractedData struct {
	Amount    float64
	Date      string
	Time      string
	Reference string
	Bank      string
	Sender    string
	Receiver  string
}

type BankPattern struct {
	Name             string
	Identifiers      []string 
	AmountPatterns   []string
	DatePatterns     []string
	TimePatterns     []string
	RefPatterns      []string
	SenderPatterns   []string
	ReceiverPatterns []string
}

var bankPatterns = []BankPattern{
	{
		Name: "SCB",
		Identifiers: []string{
			"(?i)siam\\s*commercial\\s*bank",
			"(?i)scb",
			"(?i)ธนาคารไทยพาณิชย์",
		},
		AmountPatterns: []string{
			`(?i)(?:amount|จำนวนเงิน|ยอดเงิน)[:\s]*([0-9,]+\.?\d{0,2})`,
			`(?i)(?:THB|บาท)[:\s]*([0-9,]+\.?\d{0,2})`,
			`([0-9,]+\.\d{2})\s*(?:THB|บาท|BAHT)`,
		},
		DatePatterns: []string{
			`(\d{1,2}[/-]\d{1,2}[/-]\d{2,4})`,
			`(?i)(?:date|วันที่)[:\s]*(\d{1,2}[/-]\d{1,2}[/-]\d{2,4})`,
		},
		TimePatterns: []string{
			`(\d{1,2}:\d{2}(?::\d{2})?)`,
			`(?i)(?:time|เวลา)[:\s]*(\d{1,2}:\d{2}(?::\d{2})?)`,
		},
		RefPatterns: []string{
			`(?i)(?:ref(?:erence)?|อ้างอิง|เลขที่อ้างอิง)[:\s#]*([A-Z0-9]+)`,
			`(?i)transaction\s*(?:ref|id)[:\s]*([A-Z0-9]+)`,
		},
		SenderPatterns: []string{
			`(?i)(?:from|จาก)[:\s]*([^\n]+)`,
			`(?i)sender[:\s]*([^\n]+)`,
		},
		ReceiverPatterns: []string{
			`(?i)(?:to|ถึง|ไปยัง)[:\s]*([^\n]+)`,
			`(?i)(?:receiver|ผู้รับ)[:\s]*([^\n]+)`,
		},
	},
	{
		Name: "KBank",
		Identifiers: []string{
			"(?i)kasikorn\\s*bank",
			"(?i)kbank",
			"(?i)k-bank",
			"(?i)ธนาคารกสิกรไทย",
		},
		AmountPatterns: []string{
			`(?i)(?:amount|จำนวนเงิน|ยอดเงิน)[:\s]*([0-9,]+\.?\d{0,2})`,
			`(?i)(?:THB|บาท)[:\s]*([0-9,]+\.?\d{0,2})`,
			`([0-9,]+\.\d{2})\s*(?:THB|บาท|BAHT)`,
		},
		DatePatterns: []string{
			`(\d{1,2}[/-]\d{1,2}[/-]\d{2,4})`,
			`(?i)(?:date|วันที่)[:\s]*(\d{1,2}[/-]\d{1,2}[/-]\d{2,4})`,
		},
		TimePatterns: []string{
			`(\d{1,2}:\d{2}(?::\d{2})?)`,
			`(?i)(?:time|เวลา)[:\s]*(\d{1,2}:\d{2}(?::\d{2})?)`,
		},
		RefPatterns: []string{
			`(?i)(?:ref(?:erence)?|อ้างอิง|เลขที่อ้างอิง)[:\s#]*([A-Z0-9]+)`,
			`(?i)transaction\s*(?:ref|no)[:\s]*([A-Z0-9]+)`,
		},
		SenderPatterns: []string{
			`(?i)(?:from|จาก)[:\s]*([^\n]+)`,
		},
		ReceiverPatterns: []string{
			`(?i)(?:to|ถึง)[:\s]*([^\n]+)`,
		},
	},
	{
		Name: "BBL",
		Identifiers: []string{
			"(?i)bangkok\\s*bank",
			"(?i)bbl",
			"(?i)ธนาคารกรุงเทพ",
		},
		AmountPatterns: []string{
			`(?i)(?:amount|จำนวนเงิน|ยอดเงิน)[:\s]*([0-9,]+\.?\d{0,2})`,
			`(?i)(?:THB|บาท)[:\s]*([0-9,]+\.?\d{0,2})`,
			`([0-9,]+\.\d{2})\s*(?:THB|บาท|BAHT)`,
		},
		DatePatterns: []string{
			`(\d{1,2}[/-]\d{1,2}[/-]\d{2,4})`,
			`(?i)(?:date|วันที่)[:\s]*(\d{1,2}[/-]\d{1,2}[/-]\d{2,4})`,
		},
		TimePatterns: []string{
			`(\d{1,2}:\d{2}(?::\d{2})?)`,
		},
		RefPatterns: []string{
			`(?i)(?:ref(?:erence)?|อ้างอิง)[:\s#]*([A-Z0-9]+)`,
		},
		SenderPatterns: []string{
			`(?i)(?:from|จาก)[:\s]*([^\n]+)`,
		},
		ReceiverPatterns: []string{
			`(?i)(?:to|ถึง)[:\s]*([^\n]+)`,
		},
	},
	{
		Name: "KTB",
		Identifiers: []string{
			"(?i)krung\\s*thai\\s*bank",
			"(?i)ktb",
			"(?i)ธนาคารกรุงไทย",
		},
		AmountPatterns: []string{
			`(?i)(?:amount|จำนวนเงิน|ยอดเงิน)[:\s]*([0-9,]+\.?\d{0,2})`,
			`(?i)(?:THB|บาท)[:\s]*([0-9,]+\.?\d{0,2})`,
			`([0-9,]+\.\d{2})\s*(?:THB|บาท|BAHT)`,
		},
		DatePatterns: []string{
			`(\d{1,2}[/-]\d{1,2}[/-]\d{2,4})`,
		},
		TimePatterns: []string{
			`(\d{1,2}:\d{2}(?::\d{2})?)`,
		},
		RefPatterns: []string{
			`(?i)(?:ref(?:erence)?|อ้างอิง)[:\s#]*([A-Z0-9]+)`,
		},
		SenderPatterns: []string{
			`(?i)(?:from|จาก)[:\s]*([^\n]+)`,
		},
		ReceiverPatterns: []string{
			`(?i)(?:to|ถึง)[:\s]*([^\n]+)`,
		},
	},
}

func ExtractData(ocrText string) (*ExtractedData, error) {
	if ocrText == "" {
		return nil, fmt.Errorf("OCR text is empty")
	}

	data := &ExtractedData{}

	data.Bank = detectBank(ocrText)
	log.Printf("Detected bank: %s", data.Bank)

	var patterns BankPattern
	for _, bp := range bankPatterns {
		if bp.Name == data.Bank {
			patterns = bp
			break
		}
	}
	if patterns.Name == "" {
		patterns = bankPatterns[0] 
	}

	data.Amount = extractAmount(ocrText, patterns.AmountPatterns)

	data.Date = extractField(ocrText, patterns.DatePatterns)

	data.Time = extractField(ocrText, patterns.TimePatterns)

	data.Reference = extractField(ocrText, patterns.RefPatterns)

	data.Sender = extractField(ocrText, patterns.SenderPatterns)

	data.Receiver = extractField(ocrText, patterns.ReceiverPatterns)

	if data.Amount == 0 {
		return nil, fmt.Errorf("failed to extract amount from OCR text")
	}

	return data, nil
}

func detectBank(text string) string {
	for _, bp := range bankPatterns {
		for _, identifier := range bp.Identifiers {
			re := regexp.MustCompile(identifier)
			if re.MatchString(text) {
				return bp.Name
			}
		}
	}
	return "Unknown"
}

func extractAmount(text string, patterns []string) float64 {
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(text)
		if len(matches) > 1 {
			// Remove commas and parse
			amountStr := strings.ReplaceAll(matches[1], ",", "")
			amount, err := strconv.ParseFloat(amountStr, 64)
			if err == nil && amount > 0 {
				log.Printf("Extracted amount: %.2f using pattern: %s", amount, pattern)
				return amount
			}
		}
	}
	return 0
}

func extractField(text string, patterns []string) string {
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(text)
		if len(matches) > 1 {
			result := strings.TrimSpace(matches[1])
			if result != "" {
				log.Printf("Extracted field: %s using pattern: %s", result, pattern)
				return result
			}
		}
	}
	return ""
}

func NormalizeDate(dateStr string) string {
	if dateStr == "" {
		return ""
	}

	dateStr = strings.TrimSpace(dateStr)

	re := regexp.MustCompile(`(\d{1,2})[/-](\d{1,2})[/-](\d{2,4})`)
	matches := re.FindStringSubmatch(dateStr)

	if len(matches) == 4 {
		day := matches[1]
		month := matches[2]
		year := matches[3]

		if len(day) == 1 {
			day = "0" + day
		}
		if len(month) == 1 {
			month = "0" + month
		}

		if len(year) == 2 {
			yearInt, _ := strconv.Atoi(year)
			if yearInt > 50 {
				year = "19" + year
			} else {
				year = "20" + year
			}
		}

		return fmt.Sprintf("%s/%s/%s", day, month, year)
	}

	return dateStr
}
