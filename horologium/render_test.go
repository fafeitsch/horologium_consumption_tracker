package horologium

import (
	"os"
)

func ExampleMonthlyStatistics_Render() {
	stats := MonthlyStatistics{
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
			Costs:       3252.74,
			Consumption: 75.34,
		},
		{
			ValidFrom:   CreateDate(2020, 3, 1),
			ValidTo:     CreateDate(2020, 4, 1),
			Costs:       633.28,
			Consumption: 12.53,
		},
	}
	stats.Render(os.Stdout)
	// Output:
	// |   MONTH   | YEAR | CONSUMPTION |  COSTS  |
	// |-----------|------|-------------|---------|
	// | December  | 2019 |       42.23 | 1341.12 |
	// | January   | 2020 |       53.76 | 1343.28 |
	// | February  |      |       75.34 | 3252.74 |
	// | March     |      |       12.53 |  633.28 |
}
