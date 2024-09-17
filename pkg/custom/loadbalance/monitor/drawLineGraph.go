package monitor

import (
	"fmt"
	"github.com/apache/yunikorn-core/pkg/log"
	"github.com/xuri/excelize/v2"
)

//  generates the chart series format for tenants
func (nodesMonitor *NodesMonitor) GenerateChartSeries(sheetName string) []excelize.ChartSeries {
	var series []excelize.ChartSeries
	// Generate series for each tenant
	for _, colIdx := range nodesMonitor.nodesList {
		series = append(series, excelize.ChartSeries{
			Name:       fmt.Sprintf(sheetName + "!$%s$1", colToLetter(colIdx)),                                           // Name of the tenant
			Categories: fmt.Sprintf(sheetName + "!$A$2:$A$%d", nodesMonitor.currRow-1),                                             // Time range (X-axis)
			Values:     fmt.Sprintf(sheetName + "!$%s$2:$%s$%d", colToLetter(colIdx), colToLetter(colIdx), nodesMonitor.currRow-1), // Tenant data range (Y-axis)
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

func (nodesMonitor *NodesMonitor) createGraph(sheetName string) {
	if err := nodesMonitor.excelFile.AddChart(sheetName, "I3", &excelize.Chart{
		Type: excelize.Line,
		Series: nodesMonitor.GenerateChartSeries(sheetName),
		Title: []excelize.RichTextRun{
			{
				Text: "Node "+ sheetName + " Usage Over Time",
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
					Text: sheetName + " Usage",
				},
			},
		},
		Dimension: excelize.ChartDimension{Width: 800, Height: 600},  
	}); err != nil {
		log.Log(log.Custom).Info("create graph error occur")
		return
	}
}
