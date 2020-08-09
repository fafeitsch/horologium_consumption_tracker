package horologium

import (
	"fmt"
	"io"
	"math"
	"strings"
)

// MonthlyStatistics is a slice of Statistics which contain one Statistics per month
// over a certain timespan.
type MonthlyStatistics []Statistics

// Total returns the sum of consumptions and costs (in this order)
// of the statistics.
func (m MonthlyStatistics) Total() (float64, float64) {
	consumption := 0.0
	costs := 0.0
	for _, part := range m {
		consumption = consumption + part.Consumption
		costs = costs + part.Costs
	}
	return consumption, costs
}

// RenderTable converts the MonthlyStatistics to a nice-looking table (see example).
// This method assumes that the MonthlyStatistics are sorted (earliest month first).
func (s MonthlyStatistics) RenderTable(writer io.Writer) {
	rows := make([][]string, 0, len(s))
	consumptionFormat := "%.2f"
	currencyFormat := "%.2f"
	for _, stat := range s {
		month := stat.ValidFrom.Month().String()
		year := fmt.Sprintf("%d", stat.ValidFrom.Year())
		cons := stat.FormatConsumption()
		costs := stat.FormatCosts()
		if stat.ConsumptionFormat != "" {
			consumptionFormat = stat.ConsumptionFormat
		}
		if stat.CurrencyFormat != "" {
			currencyFormat = stat.CurrencyFormat
		}
		row := []string{month, year, cons, costs}
		rows = append(rows, row)
	}
	totalConsumption, totalCosts := s.Total()
	rows = append(rows, []string{"TOTAL", "", fmt.Sprintf(consumptionFormat, totalConsumption), fmt.Sprintf(currencyFormat, totalCosts)})
	yearColLength := 6
	monthColLength := 11
	consColLength := int(math.Max(float64(longestEntry(2, rows)), 13))
	costsColLength := int(math.Max(float64(longestEntry(3, rows)), 7))
	_, _ = fmt.Fprintf(writer, "|%s|%s|%s|%s|\n", padCenter(monthColLength, "MONTH"), padCenter(yearColLength, "YEAR"), padCenter(consColLength, "CONSUMPTION"), padCenter(costsColLength, "COSTS"))
	_, _ = fmt.Fprintf(writer, "|%s|%s|%s|%s|\n", repeat(monthColLength, "-"), repeat(yearColLength, "-"), repeat(consColLength, "-"), repeat(costsColLength, "-"))
	for index, row := range rows {
		month := fmt.Sprintf(" %-10v", row[0])
		year := fmt.Sprintf(" %-5v", row[1])
		if index > 0 && strings.Contains(year, rows[index-1][1]) {
			year = repeat(yearColLength, " ")
		}
		consumptionFormat := " %" + fmt.Sprintf("%d", consColLength-2) + "v "
		consValue := fmt.Sprintf(consumptionFormat, row[2])
		costsFormat := " %" + fmt.Sprintf("%d", costsColLength-2) + "v "
		costs := fmt.Sprintf(costsFormat, row[3])
		if index == len(rows)-1 {
			_, _ = fmt.Fprintf(writer, "|%s|%s|%s|%s|\n", repeat(monthColLength, "-"), repeat(yearColLength, "-"), repeat(consColLength, "-"), repeat(costsColLength, "-"))
		}
		_, _ = fmt.Fprintf(writer, "|%s|%s|%s|%s|\n", month, year, consValue, costs)
	}
	_, _ = fmt.Fprintf(writer, "|%s|%s|%s|%s|\n", repeat(monthColLength, "-"), repeat(yearColLength, "-"), repeat(consColLength, "-"), repeat(costsColLength, "-"))
}

func longestEntry(col int, rows [][]string) int {
	max := 0
	for _, row := range rows {
		max = int(math.Max(float64(len(row[col])), float64(max)))
	}
	return max + 2
}

func repeat(number int, text string) string {
	result := ""
	for i := 0; i < number; i++ {
		result = result + text
	}
	return result
}

func padCenter(totalSize int, text string) string {
	padding := totalSize - len(text)
	padLeftFormat := "%" + fmt.Sprintf("%d", padding/2+len(text)) + "v"
	result := fmt.Sprintf(padLeftFormat, text)
	padRightFormat := "%-" + fmt.Sprintf("%d", (padding+1)/2+len(result)) + "v"
	return fmt.Sprintf(padRightFormat, result)
}
