package plotting

import (
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func PlotHistogram(values plotter.Values, name string) {
	p := plot.New()
	p.Title.Text = fmt.Sprintf("Frequency of %s", name)

	p.X.Label.Text = "Hour"
	p.X.Max = 24
	p.Y.Label.Text = "% of occurence"

	barChart, err := plotter.NewBarChart(values, 5)
	if err != nil {
		panic(err)
	}

	barChart.Color = plotutil.Color(3)
	barChart.XMin = 0
	barChart.LineStyle.Width = vg.Length(0)

	p.Add(barChart)

	if err := p.Save(3*vg.Inch, 3*vg.Inch, "barChart.png"); err != nil {
		panic(err)
	}
}

