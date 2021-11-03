package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/gbkr-com/queueing"
)

const (
	stringRow = "%-24s %10s\n"
	intRow    = "%-24s %10d\n"
	floatRow  = "%-24s %10.3f\n"
)

func main() {
	//
	// Command line flags.
	//
	var (
		m string
		c *int
		a string
		s string
		j *bool
	)
	flag.StringVar(&m, "m", "MM1", "the queueing model distribution, MM1, MMC or MD1")
	c = flag.Int("c", 1, "the number of servers")
	flag.StringVar(&a, "a", "", "the arrival rate")
	flag.StringVar(&s, "s", "", "the service rate")
	j = flag.Bool("json", false, "specify for output as JSON")
	flag.Parse()
	//
	// Validation.
	//
	var (
		model queueing.Model
		err   error
	)
	p := new(queueing.Problem)
	switch m {
	case "MM1":
		model, _ = queueing.Use(queueing.MM1)
	case "MMC":
		model, _ = queueing.Use(queueing.MMC)
	case "MD1":
		model, _ = queueing.Use(queueing.MD1)
	default:
		fmt.Println("model must be MM1, MMC or MD1")
		os.Exit(1)
	}
	p.Serving = *c
	if p.Serving < 1 {
		fmt.Println("number of servers must be one or more")
		os.Exit(1)
	}
	if p.Arrival, err = strconv.ParseFloat(a, 64); err != nil {
		fmt.Println("arrival rate must be a valid floating point number")
		os.Exit(1)
	}
	if p.Service, err = strconv.ParseFloat(s, 64); err != nil {
		fmt.Println("service rate must be a valid floating point number")
		os.Exit(1)
	}
	//
	// Calculation.
	//
	var analysis *queueing.Analysis
	if analysis, err = model.Analyse(p); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if *j {
		obj, _ := json.Marshal(analysis)
		fmt.Printf("%s\n", obj)
		os.Exit(0)
	}
	fmt.Printf(stringRow, "Model", model.Name())
	fmt.Printf(intRow, "Servers", *c)
	fmt.Printf(floatRow, "Arrival rate", p.Arrival)
	fmt.Printf(floatRow, "Service rate", p.Service)
	fmt.Printf(floatRow, "Server utilisation", analysis.Utilisation)
	fmt.Printf(floatRow, "Queue length", analysis.Queued)
	fmt.Printf(floatRow, "Sojourn time", analysis.Sojourn)
	fmt.Printf(floatRow, "Probability of queueing", analysis.Queue)
	fmt.Printf(floatRow, "Probability of loss", analysis.Loss)
}
