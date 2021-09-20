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
	playID       string
	audience     int
	play         Play
	amount       float64
	volumeCredit int
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
		enrichPerf = append(enrichPerf, s.enrichPerformances(perf))
	}
	data.performances = enrichPerf
	data.customer = s.invoice.Customer
	data.totalAmount = s.totalAmount(data)
	data.totalVolumeCredits = s.totalVolumeCredits(data)
	return data
}

func (s *statement) enrichPerformances(aPerformance Performance) performanceExt {
	calculator := NewPerformanceCalculator(aPerformance, s.playFor(aPerformance))
	ext := performanceExt{
		playID:       aPerformance.PlayID,
		audience:     aPerformance.Audience,
		play:         s.playFor(aPerformance),
		amount:       calculator.Amount(),
		volumeCredit: calculator.VolumeCredits(),
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

func (s *statement) totalAmount(data statementData) float64 {
	result := 0.0
	for _, perf := range data.performances {
		result += perf.amount
	}
	return result
}

func (s *statement) totalVolumeCredits(data statementData) int {
	result := 0
	for _, perf := range data.performances {
		result += perf.volumeCredit
	}
	return result
}

func (s *statement) usd(totalAmount float64) string {
	lc := accounting.LocaleInfo["USD"]
	ac := accounting.Accounting{Symbol: lc.ComSymbol, Precision: 2, Thousand: lc.ThouSep, Decimal: lc.DecSep}
	return ac.FormatMoney(totalAmount / 100.0)
}

func (s *statement) playFor(perf Performance) Play {
	play, ok := s.plays[perf.PlayID]
	if !ok {
		panic("unknown play: " + perf.PlayID)
	}
	return play
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
