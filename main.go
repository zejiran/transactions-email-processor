package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"gopkg.in/gomail.v2"
)

type Transaction struct {
	ID          int
	Date        time.Time
	Transaction float64
}

func main() {
	// Read the input file
	filePath := "files/example_transactions.csv"
	transactions, err := readTransactions(filePath)
	if err != nil {
		log.Fatal(err)
	}

	// Calculate summary information
	totalBalance := calculateTotalBalance(transactions)
	transactionCounts := groupTransactionsByMonth(transactions)
	averageDebit, averageCredit := calculateAverageAmounts(transactions, transactionCounts)

	// Generate the email body
	emailBody, err := generateEmailBody(totalBalance, transactionCounts, averageDebit, averageCredit)
	if err != nil {
		fmt.Println("Error generating email body:", err)
		return
	}

	// Send the email
	if err := sendEmail(emailBody); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Email sent successfully!")
}

func readTransactions(filePath string) ([]Transaction, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 3

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var transactions []Transaction
	for i, record := range records {
		// Skip the header row
		if i == 0 {
			continue
		}

		id, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, err
		}

		date, err := time.Parse("1/2", record[1])
		if err != nil {
			return nil, err
		}

		transaction, err := strconv.ParseFloat(strings.TrimPrefix(record[2], "+"), 64)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, Transaction{
			ID:          id,
			Date:        date,
			Transaction: transaction,
		})
	}

	return transactions, nil
}

func calculateTotalBalance(transactions []Transaction) float64 {
	var totalBalance float64

	for _, transaction := range transactions {
		totalBalance += transaction.Transaction
	}

	return totalBalance
}

func groupTransactionsByMonth(transactions []Transaction) map[string]int {
	transactionCounts := make(map[string]int)

	for _, transaction := range transactions {
		month := transaction.Date.Format("January")
		transactionCounts[month]++
	}

	return transactionCounts
}

func calculateAverageAmounts(transactions []Transaction, transactionCounts map[string]int) (float64, float64) {
	var debitSum, creditSum float64
	var debitCount, creditCount int

	for _, transaction := range transactions {
		if transaction.Transaction < 0 {
			debitSum += transaction.Transaction
			debitCount++
		} else {
			creditSum += transaction.Transaction
			creditCount++
		}
	}

	averageDebit := debitSum / float64(debitCount)
	averageCredit := creditSum / float64(creditCount)

	return averageDebit, averageCredit
}

func generateEmailBody(totalBalance float64, transactionCounts map[string]int, averageDebit, averageCredit float64) (string, error) {
	// Read the HTML template from the external file
	templateFile := "email/email_template.html"
	templateContent, err := ioutil.ReadFile(templateFile)
	if err != nil {
		return "", err
	}

	// Prepare the data for the email template
	emailData := struct {
		TotalBalance      string
		TransactionCounts map[string]int
		AverageDebit      string
		AverageCredit     string
	}{
		TotalBalance:      fmt.Sprintf("$%s", humanize.FormatFloat("#,###.##", totalBalance)),
		TransactionCounts: transactionCounts,
		AverageDebit:      fmt.Sprintf("$%s", humanize.FormatFloat("#,###.##", averageDebit)),
		AverageCredit:     fmt.Sprintf("$%s", humanize.FormatFloat("#,###.##", averageCredit)),
	}

	// Create a new template and parse the email template content
	tmpl := template.New("emailTemplate")
	tmpl, err = tmpl.Parse(string(templateContent))
	if err != nil {
		return "", err
	}

	// Generate the email body by executing the template with the email data
	var emailBody bytes.Buffer
	err = tmpl.Execute(&emailBody, emailData)
	if err != nil {
		return "", err
	}

	return emailBody.String(), nil
}

func sendEmail(body string) error {
	sender := os.Getenv("SENDER_MAIL")
	password := os.Getenv("PASSWORD")
	recipient := os.Getenv("RECIPIENT_MAIL")

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", sender)
	mailer.SetHeader("To", recipient)
	mailer.SetHeader("Subject", "Transaction Summary")
	mailer.SetBody("text/html", body)

	dialer := gomail.NewDialer("smtp.gmail.com", 587, sender, password)

	if err := dialer.DialAndSend(mailer); err != nil {
		return err
	}

	return nil
}
