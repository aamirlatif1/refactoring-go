package invoice

import (
	"errors"
	"fmt"
	"github.com/leekchan/accounting"
	"strings"
)

func plainStatement(invoice Invoice, plays map[string]Play) (string, error) {
	g := statement{invoice, plays}
	return g.renderPainText()
}

type statement struct {
	invoice Invoice
	plays   map[string]Play
}

type statementData struct {
	customer string
}

func (s *statement) renderPainText() (string, error) {
	data := statementData{
		customer: s.invoice.Customer,
	}
	return s.renderPlainText(data)
}

func (s *statement) renderPlainText(data statementData) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Statement for %s\n", data.customer))
	for _, perf := range s.invoice.Performances {
		if _, ok := s.plays[perf.PlayID]; !ok {
			return "", errors.New("unknown type: " + perf.PlayID)
		}
		sb.WriteString(fmt.Sprintf(" %s: %s (%d seats)\n", s.playFor(perf).Name, s.usd(s.amountFor(perf)), perf.Audience))
	}
	sb.WriteString(fmt.Sprintf("Amount owed is %s\n", s.usd(s.totalAmount())))
	sb.WriteString(fmt.Sprintf("You earned %d credits\n", s.totalVolumeCredits()))
	return sb.String(), nil
}

func (s *statement) totalAmount() float64 {
	result := 0.0
	for _, perf := range s.invoice.Performances {
		result += s.amountFor(perf)
	}
	return result
}

func (s *statement) totalVolumeCredits() int {
	result := 0
	for _, perf := range s.invoice.Performances {
		result += s.volumeCreditFor(perf)
	}
	return result
}

func (s *statement) usd(totalAmount float64) string {
	lc := accounting.LocaleInfo["USD"]
	ac := accounting.Accounting{Symbol: lc.ComSymbol, Precision: 2, Thousand: lc.ThouSep, Decimal: lc.DecSep}
	return ac.FormatMoney(totalAmount / 100.0)
}

func (s *statement) volumeCreditFor(perf Performance) int {
	result := Max(perf.Audience-30, 0)
	if "comedy" == s.playFor(perf).PlayType {
		result += perf.Audience / 5
	}
	return result
}

func (s *statement) playFor(perf Performance) Play {
	play, _ := s.plays[perf.PlayID]
	return play
}

func (s *statement) amountFor(perf Performance) float64 {
	amount := 0.0
	switch s.playFor(perf).PlayType {
	case "tragedy":
		amount = 40000.0
		if perf.Audience > 30 {
			amount += 1000 * float64(perf.Audience-30)
		}
	case "comedy":
		amount = 30000.0
		if perf.Audience > 20 {
			amount += 10000 + 500*float64(perf.Audience-20)
		}
		amount += 300 * float64(perf.Audience)
	default:
		panic("unknown type: " + s.playFor(perf).PlayType)
	}
	return amount
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
