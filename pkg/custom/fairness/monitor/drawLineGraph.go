package monitor

import (
	"fmt"
	"github.com/apache/yunikorn-core/pkg/log"
	"github.com/xuri/excelize/v2"
)

// GenerateChartSeries generates the chart series format for tenants
func (tenantsMonitor *TenantsMonitor) GenerateChartSeries() []excelize.ChartSeries {
	var series []excelize.ChartSeries
	// Generate series for each tenant
	for _, colIdx := range tenantsMonitor.tenantsList {
		series = append(series, excelize.ChartSeries{
			Name:       fmt.Sprintf("log!$%s$1", colToLetter(colIdx)),                                           // Name of the tenant
			Categories: fmt.Sprintf("log!$A$2:$A$%d", tenantsMonitor.currRow-1),                                             // Time range (X-axis)
			Values:     fmt.Sprintf("log!$%s$2:$%s$%d", colToLetter(colIdx), colToLetter(colIdx), tenantsMonitor.currRow-1), // Tenant data range (Y-axis)
		})
	}
	return series
}

// colToLetter converts column index to letter, e.g., 1 -> A, 2 -> B, 3 -> C
func colToLetter(colIdx int) string {
	result := ""
	for colIdx > 0 {
		colIdx--
		result = string('A'+(colIdx%26)) + result
		colIdx /= 26
	}
	return result
}

func (tenantsMonitor *TenantsMonitor) createGraph() {
	if err := tenantsMonitor.excelFile.AddChart("Sheet1", "A1", &excelize.Chart{
		Type: excelize.Line,
		Series: tenantsMonitor.GenerateChartSeries(),
		Title: []excelize.RichTextRun{
			{
				Text: "User Resource Usage Over Time",
			},
		},
		XAxis: excelize.ChartAxis{
			MajorGridLines: true,
			Title: []excelize.RichTextRun{
				{
					Text: "Time",
				},
			},
		},
		YAxis: excelize.ChartAxis{
			MajorGridLines: true,
			Title: []excelize.RichTextRun{
				{
					Text: "Donimant Resource Usage",
				},
			},
		},
	}); err != nil {
		log.Log(log.Custom).Info("create graph error occur")
		return
	}
}
