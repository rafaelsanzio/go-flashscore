package money

import (
	"fmt"
	"math"

	"github.com/leekchan/accounting"

	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
)

// Money contains an amount and a currency.
type Money struct {
	// Cents is used internally as an imprecise term for "minor currency unit"
	// because "minor currency unit" sounds awful -- apologies to folks outside
	// the US and EU who have a different name for their "minor currency unit"
	// required: true
	// example: 350
	Cents int
	// required: true
	Currency Currency // required...and needs to be valid
}

func (m Money) String() string {
	ac := accounting.Accounting{Precision: m.Currency.Scale, FormatNegative: "(%v)"}

	amtAsFloat := m.Float()

	return fmt.Sprintf("%s %s", ac.FormatMoneyFloat64(amtAsFloat), m.Currency.AlphaCode)
}

func (m Money) ToAmountString() string {
	ac := accounting.Accounting{Precision: m.Currency.Scale}
	return ac.FormatMoneyFloat64(m.Float())
}

func (m Money) Float() float64 {
	div := math.Pow10(m.Currency.Scale)
	return float64(m.Cents) / div
}

func (m Money) Validate() errs.AppError {
	// Validate currency
	_, ok := SafeCurrencyLookup(m.Currency.AlphaCode)

	if !ok {
		return errs.ErrInvalidMoney.Throwf(applog.Log, errs.ErrFmt, "Currency", m.Currency.AlphaCode, "ISO currency Code")
	}

	return nil
}

func (m Money) IsZeroMoney() bool {
	return m.Cents == 0
}
