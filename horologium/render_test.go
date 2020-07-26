package horologium

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

const wantTable = `|   MONTH   | YEAR | CONSUMPTION |   COSTS    |
|-----------|------|-------------|------------|
| December  | 2019 |       42.23 |    1341.12 |
| January   | 2020 |       53.76 |    1343.28 |
| February  |      |       75.34 | 1123744.74 |
| March     |      |       12.53 |     633.28 |
`

func TestMonthlyStatistics(t *testing.T) {
	stats := []Statistics{
		{
			ValidFrom:   CreateDate(2019, 12, 1),
			ValidTo:     CreateDate(2020, 1, 1),
			Costs:       1341.12,
			Consumption: 42.23,
		},
		{
			ValidFrom:   CreateDate(2020, 1, 1),
			ValidTo:     CreateDate(2020, 2, 1),
			Costs:       1343.28,
			Consumption: 53.76,
		},
		{
			ValidFrom:   CreateDate(2020, 2, 1),
			ValidTo:     CreateDate(2020, 3, 1),
			Costs:       1123744.74,
			Consumption: 75.34,
		},
		{
			ValidFrom:   CreateDate(2020, 3, 1),
			ValidTo:     CreateDate(2020, 4, 1),
			Costs:       633.28,
			Consumption: 12.53,
		},
	}
	buffer := bytes.Buffer{}
	MonthlyStatistics(&buffer, stats)
	assert.Equal(t, wantTable, buffer.String(), "rendered table is wrong")
}
