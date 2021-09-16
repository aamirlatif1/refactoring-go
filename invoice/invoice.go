package invoice

import (
	"errors"
	"fmt"
	"github.com/leekchan/accounting"
)

func statement(invoice Invoice, plays map[string]Play) (*string, error) {
	totalAmount := 0.0
	volumeCredits := 0
	result := fmt.Sprintf("Statement for %s\n", invoice.Customer)
	lc := accounting.LocaleInfo["USD"]
	ac := accounting.Accounting{Symbol: lc.ComSymbol, Precision: 2, Thousand: lc.ThouSep, Decimal: lc.DecSep}

	for _, perf := range invoice.Performances {
		if _, ok := plays[perf.PlayID]; !ok {
			return nil, errors.New("unknown type: " + perf.PlayID)
		}
		play, _ := plays[perf.PlayID]
		thisAmount := 0.0

		switch play.PlayType {
		case "tragedy":
			thisAmount = 40000.0
			if perf.Audience > 30 {
				thisAmount += 1000 * float64(perf.Audience-30)
			}
		case "comedy":
			thisAmount = 30000.0
			if perf.Audience > 20 {
				thisAmount += 10000 + 500*float64(perf.Audience-20)
			}
			thisAmount += 300 * float64(perf.Audience)
		default:
			return nil, errors.New("unknown type: " + play.PlayType)
		}

		volumeCredits += Max(perf.Audience-30, 0)
		if "comedy" == play.PlayType {
			volumeCredits += perf.Audience / 5
		}

		result += fmt.Sprintf(" %s: %s (%d seats)\n", play.Name, ac.FormatMoney(thisAmount/100.0), perf.Audience)
		totalAmount += thisAmount
	}
	result += fmt.Sprintf("Amount owed is %s\n", ac.FormatMoney(totalAmount/100.0))
	result += fmt.Sprintf("You earned %d credits\n", volumeCredits)
	return &result, nil
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
