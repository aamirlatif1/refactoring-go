package invoice

import (
	"errors"
	"fmt"
	"github.com/leekchan/accounting"
)

func statement(invoice Invoice, plays map[string]Play) (*string, error) {
	g := generator{invoice, plays}
	return g.generate()
}

type generator struct {
	invoice Invoice
	plays   map[string]Play
}

func (g *generator) generate() (*string, error) {
	totalAmount := 0.0
	volumeCredits := 0
	result := fmt.Sprintf("Statement for %s\n", g.invoice.Customer)

	for _, perf := range g.invoice.Performances {
		if _, ok := g.plays[perf.PlayID]; !ok {
			return nil, errors.New("unknown type: " + perf.PlayID)
		}
		volumeCredits += g.volumeCreditFor(perf)

		result += fmt.Sprintf(" %s: %s (%d seats)\n", g.playFor(perf).Name, g.usd(g.amountFor(perf)), perf.Audience)
		totalAmount += g.amountFor(perf)
	}
	result += fmt.Sprintf("Amount owed is %s\n", g.usd(totalAmount))
	result += fmt.Sprintf("You earned %d credits\n", volumeCredits)
	return &result, nil
}

func (g *generator) usd(totalAmount float64) string {
	lc := accounting.LocaleInfo["USD"]
	ac := accounting.Accounting{Symbol: lc.ComSymbol, Precision: 2, Thousand: lc.ThouSep, Decimal: lc.DecSep}
	return ac.FormatMoney(totalAmount / 100.0)
}

func (g *generator) volumeCreditFor(perf Performance) int {
	result := Max(perf.Audience-30, 0)
	if "comedy" == g.playFor(perf).PlayType {
		result += perf.Audience / 5
	}
	return result
}

func (g *generator) playFor(perf Performance) Play {
	play, _ := g.plays[perf.PlayID]
	return play
}

func (g *generator) amountFor(perf Performance) float64 {
	amount := 0.0
	switch g.playFor(perf).PlayType {
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
		panic("unknown type: " + g.playFor(perf).PlayType)
	}
	return amount
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
