// Package repository содержит функционал по генерации рандомного html по списку метрик
package repository

import (
	"fmt"
)

// RenderHTML генерация html с данными метрик в произвольной форме.
func RenderHTML(metrics []Metrics) string {
	htmlTmpl := `
	<html>
	<body>
	<h1>All metrics</h1>
	<h3>gauges</h3>
	<ul>
	  %s
	</ul>
	<h3>counters</h3>
	<ul>
	  %s
	</ul>
	</body>
	</html>`

	gaugesText := ""
	countersText := ""
	for _, value := range metrics {
		switch value.MType {
		case Gauge:
			gaugesText += fmt.Sprintf("<li>%s : %.3f</li>\n", value.ID, *value.Value)
		case Counter:
			countersText += fmt.Sprintf("<li>%s : %d</li>\n", value.ID, *value.Delta)
		}
	}
	return fmt.Sprintf(htmlTmpl, gaugesText, countersText)
}
