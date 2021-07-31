package plotting

import (
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"io"
)

func PlotHistogram(values plotter.Values, name string) io.WriterTo {
	p := plot.New()
	p.Title.Text = fmt.Sprintf("Frequency of %s", name)

	p.X.Label.Text = "Hour"
	p.X.Max = 24
	p.X.Tick.Marker = plot.DefaultTicks{NumberOfTicks: 24}

	p.Y.Label.Text = "% of occurence"
	p.Y.Tick.Marker = plot.DefaultTicks{NumberOfTicks: 10}

	barChart, err := plotter.NewBarChart(values, 5)
	if err != nil {
		panic(err)
	}

	barChart.Color = plotutil.Color(3)
	barChart.XMin = 0
	barChart.LineStyle.Width = vg.Length(0)

	p.Add(barChart)

	writer, err := p.WriterTo(6*vg.Inch, 3*vg.Inch, "png")
	if err != nil {
		panic(err)
	}

	return writer
}

