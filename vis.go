package main

import (
	"fmt"
	giles "github.com/gtfierro/giles/archiver"
	"github.com/gtfierro/go-gnuplot"
	"github.com/joliv/spark"
)

func visSpark(data []giles.SmapMessage) {
	for _, msg := range data {
		readings := extractDataNumeric(msg)
		sparkline := spark.Line(readings)
		fmt.Printf("%s %v %v %v\n", msg.UUID, readings[0], sparkline, readings[len(readings)-1])
	}
}

func visPlot(data []giles.SmapMessage) {
	fname := ""
	persist := false
	debug := true
	for _, msg := range data {
		readings := extractDataNumeric(msg)
		fmt.Println(len(readings))
		p, err := gnuplot.NewPlotter(fname, persist, debug)
		if err != nil {
			err_string := fmt.Sprintf("** err: %v\n", err)
			panic(err_string)
		}
		defer p.Close()

		p.CheckedCmd("set terminal dumb")
		p.PlotX(readings, msg.UUID)
		p.CheckedCmd("q")
	}
}

func visPlotTime(data []giles.SmapMessage) {
	fname := ""
	persist := false
	debug := false
	for _, msg := range data {
		readings := extractDataNumeric(msg)
		times := extractTime(msg)
		p, err := gnuplot.NewPlotter(fname, persist, debug)
		if err != nil {
			err_string := fmt.Sprintf("** err: %v\n", err)
			panic(err_string)
		}
		defer p.Close()

		p.CheckedCmd("set terminal dumb")
		p.PlotXY(times, readings, msg.UUID)
		p.CheckedCmd("q")
	}
}
