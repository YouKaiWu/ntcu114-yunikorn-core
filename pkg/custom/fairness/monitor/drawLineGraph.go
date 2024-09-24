package monitor

import (
	"fmt"
	"github.com/apache/yunikorn-core/pkg/log"
	"github.com/xuri/excelize/v2"
)

// GenerateChartSeries generates the chart series format for tenants
func (tenantsMonitor *TenantsMonitor) GenerateLineChartSeries(sheetName string) []excelize.ChartSeries {
	var series []excelize.ChartSeries
	// Generate series for each tenant
	for _, colIdx := range tenantsMonitor.tenantsList {
		series = append(series, excelize.ChartSeries{
			Name:       fmt.Sprintf(sheetName+"!$%s$1", colToLetter(colIdx)),                                                       // Name of the tenant
			Categories: fmt.Sprintf(sheetName+"!$A$2:$A$%d", tenantsMonitor.currRow-1),                                             // Time range (X-axis)
			Values:     fmt.Sprintf(sheetName+"!$%s$2:$%s$%d", colToLetter(colIdx), colToLetter(colIdx), tenantsMonitor.currRow-1), // Tenant data range (Y-axis)
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

func (tenantsMonitor *TenantsMonitor) createLineChart(sheetName string) {
	if err := tenantsMonitor.excelFile.AddChart(sheetName, "I3", &excelize.Chart{
		Type:   excelize.Line,
		Series: tenantsMonitor.GenerateLineChartSeries(sheetName),
		Title: []excelize.RichTextRun{
			{
				Text: "User " + sheetName + " Usage Over Time",
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

func (tenantsMonitor *TenantsMonitor) GenerateTotalWaitTimeBarChartSeries() []excelize.ChartSeries {
	var series []excelize.ChartSeries

	series = append(series, excelize.ChartSeries{
		Name:       "Total Wait Time",                                            
		Categories: "WaitTime!$B$1:$" + colToLetter(len(tenantsMonitor.tenantsList)+1) + "$1", 
		Values:     "WaitTime!$B$2:$" + colToLetter(len(tenantsMonitor.tenantsList)+1) + "$2", 
	})

	return series
}

func (tenantsMonitor *TenantsMonitor) GenerateAvgWaitTimeBarChartSeries() []excelize.ChartSeries {
	var series []excelize.ChartSeries

	series = append(series, excelize.ChartSeries{
		Name:       "Average Wait Time",                                         
		Categories: "WaitTime!$B$1:$" + colToLetter(len(tenantsMonitor.tenantsList)+1) + "$1", 
		Values:     "WaitTime!$B$3:$" + colToLetter(len(tenantsMonitor.tenantsList)+1) + "$3", 
	})

	return series
}

func (tenantsMonitor *TenantsMonitor) createTotalWaitTimeBarChart(sheetName string) {
	if err := tenantsMonitor.excelFile.AddChart(sheetName, "A5", &excelize.Chart{
		Type:   excelize.Col, 
		Series: tenantsMonitor.GenerateTotalWaitTimeBarChartSeries(),
		Title: []excelize.RichTextRun{
			{
				Text: "User Total Wait Time Overview",
			},
		},
		XAxis: excelize.ChartAxis{
			MajorGridLines: true,
			Title: []excelize.RichTextRun{
				{
					Text: "Users",
				},
			},
		},
		YAxis: excelize.ChartAxis{
			MajorGridLines: true,
			Title: []excelize.RichTextRun{
				{
					Text: "Total Wait Time (Seconds)",
				},
			},
		},
		Dimension: excelize.ChartDimension{Width: 400, Height: 200},
	}); err != nil {
		log.Log(log.Custom).Info("create horizontal bar chart error occur")
		return
	}
}

func (tenantsMonitor *TenantsMonitor) createAvgWaitTimeBarChart(sheetName string) {
	if err := tenantsMonitor.excelFile.AddChart(sheetName, "I5", &excelize.Chart{
		Type:   excelize.Col,
		Series: tenantsMonitor.GenerateAvgWaitTimeBarChartSeries(),
		Title: []excelize.RichTextRun{
			{
				Text: "User Average Wait Time Overview",
			},
		},
		XAxis: excelize.ChartAxis{
			MajorGridLines: true,
			Title: []excelize.RichTextRun{
				{
					Text: "Users",
				},
			},
		},
		YAxis: excelize.ChartAxis{
			MajorGridLines: true,
			Title: []excelize.RichTextRun{
				{
					Text: "Average Wait Time (Seconds)",
				},
			},
		},
		Dimension: excelize.ChartDimension{Width: 400, Height: 200},
	}); err != nil {
		log.Log(log.Custom).Info("create average wait time bar chart error occur")
		return
	}
}
