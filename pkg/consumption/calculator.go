package consumption

import "time"

type Calculator interface {
	forMonths(time.Time, time.Time) []float64
}

type calculatorImpl struct {
}
