package queueing

import (
	"math"
	"testing"
)

func TestUse(t *testing.T) {
	cases := []struct {
		m  string
		ok bool
	}{
		{MD1, true},
		{MM1, true},
		{"M/M/1", true},
		{MMC, true},
		{"XXX", false},
	}
	for _, c := range cases {
		_, err := Use(c.m)
		if c.ok && err != nil {
			t.Error()
		}
	}
}

func TestMD1(t *testing.T) {
	//
	// See https://personalpages.bradley.edu/~rf/wait-md1.htm
	//
	m, _ := Use(MD1)
	p := &Problem{
		Serving: 1,
		Arrival: 20,
		Service: 25,
	}
	result, _ := m.Analyse(p)
	if math.Round(1000*result.Queued)/1000 != 1.6 {
		t.Error()
	}
	if math.Round(1000*result.Sojourn)/1000 != 0.12 {
		t.Error()
	}
}

func TestMM1(t *testing.T) {
	//
	// See Queueing Theory page 35
	//
	cases := []struct {
		arrival float64
		result  float64 // Note, these are waiting times only.
	}{
		{0.5, 1},
		{0.8, 4},
		{0.9, 9},
		{0.95, 19},
	}
	m, _ := Use(MM1)
	p := &Problem{Serving: 1, Service: 1.0}
	for _, c := range cases {
		p.Arrival = c.arrival
		result, _ := m.Analyse(p)
		//
		// Substract service time as the expected results are waiting times only.
		//
		result.Sojourn -= 1.0
		if math.Round(1000*result.Sojourn)/1000 != c.result {
			t.Error()
		}
	}
}

func TestErlangB(t *testing.T) {
	//
	// See Queueing Theory page 114
	//
	cases := []struct {
		servers int
		result  float64
	}{
		{150, 0.062},
		{155, 0.044},
		{160, 0.028},
		{165, 0.017},
		{170, 0.009},
	}
	for _, c := range cases {
		b := ErlangB(c.servers, 60.0, 1/2.5)
		if math.Round(1000*b)/1000 != c.result {
			t.Error()
		}
	}
}

func TestErlangC(t *testing.T) {
	//
	// See Queueing Theory page 45
	//
	cases := []struct {
		servers int
		result  float64
	}{
		{1, 0.90},
		{2, 0.85},
		{5, 0.76},
		{10, 0.67},
		{20, 0.55},
	}
	for _, c := range cases {
		ec := ErlangC(c.servers, 0.9, 1.0)
		if math.Round(100*ec)/100 != c.result {
			t.Error()
		}
	}
}

func TestMMC(t *testing.T) {
	//
	// See Queueing Theory page 45
	//
	cases := []struct {
		servers int
		arrival float64
		queued  int
		sojourn float64
	}{
		{1, 0.9, 8, 9},         //   9 -  1
		{2, 0.95, 17, 9.26},    //  19 -  2
		{5, 0.98, 46, 9.5},     //  51 -  5
		{10, 0.99, 95, 9.64},   // 105 - 10
		{20, 0.995, 193, 9.74}, // 214 - 20 ... but rounding
	}
	p := &Problem{Service: 1.0}
	m, _ := Use(MMC)
	for _, c := range cases {
		p.Serving = c.servers
		p.Arrival = c.arrival
		result, _ := m.Analyse(p)
		if int(result.Queued) != c.queued {
			t.Error()
		}
		if math.Round(100*result.Sojourn)/100 != c.sojourn {
			t.Error()
		}
	}
}
