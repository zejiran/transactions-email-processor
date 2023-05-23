package main

import (
	"testing"
	"time"
)

func TestCalculateTotalBalance(t *testing.T) {
	transactions := []Transaction{
		{ID: 0, Date: time.Now(), Transaction: 100},
		{ID: 1, Date: time.Now(), Transaction: -50},
		{ID: 2, Date: time.Now(), Transaction: 75},
	}

	expectedBalance := 125.0
	totalBalance := calculateTotalBalance(transactions)

	if totalBalance != expectedBalance {
		t.Errorf("Expected total balance %.2f, got %.2f", expectedBalance, totalBalance)
	}
}

func TestGroupTransactionsByMonth(t *testing.T) {
	transactions := []Transaction{
		{ID: 0, Date: time.Date(2023, time.January, 10, 0, 0, 0, 0, time.UTC), Transaction: 100},
		{ID: 1, Date: time.Date(2023, time.February, 15, 0, 0, 0, 0, time.UTC), Transaction: -50},
		{ID: 2, Date: time.Date(2023, time.January, 20, 0, 0, 0, 0, time.UTC), Transaction: 75},
	}

	expectedCounts := map[string]int{
		"January":  2,
		"February": 1,
	}
	transactionCounts := groupTransactionsByMonth(transactions)

	if len(transactionCounts) != len(expectedCounts) {
		t.Errorf("Expected %d months, got %d", len(expectedCounts), len(transactionCounts))
	}

	for month, expectedCount := range expectedCounts {
		if count, ok := transactionCounts[month]; !ok || count != expectedCount {
			t.Errorf("Expected transaction count for %s: %d, got %d", month, expectedCount, count)
		}
	}
}

func TestCalculateAverageAmounts(t *testing.T) {
	transactions := []Transaction{
		{ID: 0, Date: time.Now(), Transaction: 100},
		{ID: 1, Date: time.Now(), Transaction: -50},
		{ID: 2, Date: time.Now(), Transaction: -25},
		{ID: 3, Date: time.Now(), Transaction: 75},
	}

	expectedDebit := -37.5
	expectedCredit := 87.5
	averageDebit, averageCredit := calculateAverageAmounts(transactions)

	if averageDebit != expectedDebit {
		t.Errorf("Expected average debit %.2f, got %.2f", expectedDebit, averageDebit)
	}
	if averageCredit != expectedCredit {
		t.Errorf("Expected average credit %.2f, got %.2f", expectedCredit, averageCredit)
	}
}
