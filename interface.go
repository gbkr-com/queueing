package queueing

// Interface defines the queueing problem to be analysed.
//
type Interface interface {

	// Servers returns the number of servers.
	//
	Servers() int

	// ArrivalRate returns the mean number of arrivals per unit time.
	//
	ArrivalRate() float64

	// ServiceRate returns the mean number serviced per unit time.
	//
	ServiceRate() float64
}

// Problem is a simple implementation of Interface.
//
type Problem struct {
	Serving int
	Arrival float64
	Service float64
}

// Servers implements Interface.
//
func (p Problem) Servers() int { return p.Serving }

// ArrivalRate implements Interface.
//
func (p Problem) ArrivalRate() float64 { return p.Arrival }

// ServiceRate implements Interface.
//
func (p Problem) ServiceRate() float64 { return p.Service }
