package repository

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRenderHTML(t *testing.T) {
	testCounter := int64(100)
	testGauge := float64(110)
	metrics := []Metrics{
		{
			ID:    "testCounter",
			MType: Counter,
			Delta: &testCounter,
		},
		{
			ID:    "testGauge",
			MType: Gauge,
			Value: &testGauge,
		},
	}

	want := `
	<html>
	<body>
	<h1>All metrics</h1>
	<h3>gauges</h3>
	<ul>
	  <li>testGauge : 110.000</li>

	</ul>
	<h3>counters</h3>
	<ul>
	  <li>testCounter : 100</li>

	</ul>
	</body>
	</html>`

	assert.Equal(t, RenderHTML(metrics), want)
}
