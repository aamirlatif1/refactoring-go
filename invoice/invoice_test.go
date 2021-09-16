package invoice

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestStatementGeneratorSuite(t *testing.T) {
	suite.Run(t, new(StatementGeneratorSuite))
}

type StatementGeneratorSuite struct {
	suite.Suite
}

func (suite *StatementGeneratorSuite) SetupTest() {
}

func (suite *StatementGeneratorSuite) TestUnknownTypeInInvoice() {
	invoice := Invoice{
		Customer: "BigCo",
		Performances: []Performance{
			{"hamlet", 55},
			{"as-like", 35},
			{"othello", 40},
		},
	}
	plays := map[string]Play{
		"hamlet":  {"Hamlet", "sci-fi"},
		"as-like": {"As You Like It", "comedy"},
		"othello": {"Othello", "tragedy"},
	}
	_, err := statement(invoice, plays)

	assert.EqualError(suite.T(), err, "unknown type: sci-fi")
}

func (suite *StatementGeneratorSuite) TestPlayIdNotFoundInPerformances() {
	invoice := Invoice{
		Customer: "BigCo",
		Performances: []Performance{
			{"hamlet2", 55},
			{"as-like", 35},
			{"othello", 40},
		},
	}
	plays := map[string]Play{
		"hamlet":  {"Hamlet", "sci-fi"},
		"as-like": {"As You Like It", "comedy"},
		"othello": {"Othello", "tragedy"},
	}
	_, err := statement(invoice, plays)
	assert.EqualError(suite.T(), err, "unknown type: hamlet2")
}

func (suite *StatementGeneratorSuite) TestGenerateStatementSuccess() {
	invoice := Invoice{
		Customer: "BigCo",
		Performances: []Performance{
			{"hamlet", 55},
			{"as-like", 35},
			{"othello", 40},
		},
	}
	plays := map[string]Play{
		"hamlet":  {"Hamlet", "tragedy"},
		"as-like": {"As You Like It", "comedy"},
		"othello": {"Othello", "tragedy"},
	}

	actual, err := statement(invoice, plays)

	expected := "Statement for BigCo\n" +
		" Hamlet: $650.00 (55 seats)\n" +
		" As You Like It: $580.00 (35 seats)\n" +
		" Othello: $500.00 (40 seats)\n" +
		"Amount owed is $1,730.00\n" +
		"You earned 47 credits\n"

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), expected, *actual)
}
