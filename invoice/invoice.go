package invoice

import (
	"fmt"
	"github.com/leekchan/accounting"
	"strings"
)

func plainStatement(invoice Invoice, plays map[string]Play) string {
	g := statement{invoice, plays}
	return g.plainStatement()
}

type statement struct {
	invoice Invoice
	plays   map[string]Play
}

type performanceExt struct {
	playID   string
	audience int
	play     Play
	amount   float64
}

type statementData struct {
	customer           string
	totalAmount        float64
	totalVolumeCredits int
	performances       []performanceExt
}

func (s *statement) plainStatement() string {
	return s.renderPlainStatement(s.getStatementData())
}

func (s *statement) getStatementData() statementData {
	data := statementData{}
	var enrichPerf []performanceExt
	for _, perf := range s.invoice.Performances {
		enrichPerf = append(enrichPerf, s.enrich(perf))
	}
	data.performances = enrichPerf
	data.customer = s.invoice.Customer
	data.totalAmount = s.totalAmount()
	data.totalVolumeCredits = s.totalVolumeCredits()
	return data
}

func (s *statement) enrich(perf Performance) performanceExt {
	ext := performanceExt{
		playID:   perf.PlayID,
		audience: perf.Audience,
		play:     s.playFor(perf),
		amount:   s.amountFor(perf),
	}
	return ext
}

func (s *statement) renderPlainStatement(data statementData) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Statement for %s\n", data.customer))
	for _, perf := range data.performances {
		sb.WriteString(fmt.Sprintf(" %s: %s (%d seats)\n", perf.play.Name, s.usd(perf.amount), perf.audience))
	}
	sb.WriteString(fmt.Sprintf("Amount owed is %s\n", s.usd(data.totalAmount)))
	sb.WriteString(fmt.Sprintf("You earned %d credits\n", data.totalVolumeCredits))
	return sb.String()
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
	play, ok := s.plays[perf.PlayID]
	if !ok {
		panic("unknown play: " + perf.PlayID)
	}
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
