package queueing

// An Analysis of a queueing problem.
//
type Analysis struct {

	// Utilisation is the expected utilisation of each server.
	//
	Utilisation float64

	// Queued is the number of arrivals waiting to be serviced at any moment.
	//
	Queued float64

	// Sojourn is the mean time an arrival is in the system, service and waiting
	// time.
	//
	Sojourn float64

	// Queue is the probabilty of an arrival being queued.
	//
	Queue float64

	// Loss is the probability of an arrival being lost if queueing
	// is not allowed.
	//
	Loss float64
}

// Model is a queueing model.
type Model interface {
	Name() string
	Analyse(Interface) (*Analysis, error)
}

type md1 struct{}

func (*md1) Name() string { return MD1 }

func (*md1) Analyse(problem Interface) (analysis *Analysis, err error) {
	rho := problem.ArrivalRate() / problem.ServiceRate()
	if rho >= 1.0 {
		err = ErrUtilisation
		return
	}
	analysis = new(Analysis)
	analysis.Utilisation = rho
	analysis.Queue = rho
	analysis.Loss = rho
	//
	// Extract a common factor.
	//
	f := 0.5 * odds(rho)
	analysis.Queued = rho * f
	analysis.Sojourn = (1.0 + f) / problem.ServiceRate()
	return
}

type mm1 struct{}

func (*mm1) Name() string { return MM1 }

func (*mm1) Analyse(problem Interface) (analysis *Analysis, err error) {
	rho := problem.ArrivalRate() / problem.ServiceRate()
	if rho >= 1.0 {
		err = ErrUtilisation
		return
	}
	analysis = new(Analysis)
	analysis.Utilisation = rho
	analysis.Queue = rho
	analysis.Loss = rho
	//
	// Common factor.
	//
	f := odds(rho)
	analysis.Queued = rho * f
	analysis.Sojourn = (1 / problem.ArrivalRate()) * f
	return
}

type mmc struct{}

func (*mmc) Name() string { return MMC }

func (*mmc) Analyse(problem Interface) (analysis *Analysis, err error) {
	rho := problem.ArrivalRate() / (float64(problem.Servers()) * problem.ServiceRate())
	if rho >= 1.0 {
		err = ErrUtilisation
		return
	}
	analysis = new(Analysis)
	analysis.Utilisation = rho
	analysis.Queue = ErlangC(problem.Servers(), problem.ArrivalRate(), problem.ServiceRate())
	analysis.Loss = ErlangB(problem.Servers(), problem.ArrivalRate(), problem.ServiceRate())
	//
	// Rho now redefined as the total intensity.
	//
	rho = problem.ArrivalRate() / problem.ServiceRate()
	analysis.Queued = analysis.Queue * odds(rho)
	capacity := float64(problem.Servers()) * problem.ServiceRate()
	analysis.Sojourn = (analysis.Queue / capacity) + (analysis.Queued / capacity)
	return
}

var (
	md1Model = new(md1)
	mm1Model = new(mm1)
	mmcModel = new(mmc)
)

// Use a named queueing model.
//
func Use(m string) (Model, error) {
	switch m {
	case MD1:
		return md1Model, nil
	case MM1:
		return mm1Model, nil
	case MMC:
		return mmcModel, nil
	default:
		return nil, ErrModel
	}
}

func odds(x float64) float64 {
	return x / (1 - x)
}

// ErlangB returns the blocking probability, or loss, of an arrival to a system.
//
func ErlangB(servers int, arrivalRate, serviceRate float64) float64 {
	return erlangB(servers, arrivalRate/serviceRate)
}

func erlangB(c int, rho float64) float64 {
	switch c {
	case 0:
		return 1.0
	default:
		x := rho * erlangB(c-1, rho)
		return x / (x + float64(c))
	}
}

// ErlangC returns the probability of an arrival being queued.
//
func ErlangC(servers int, arrivalRate, serviceRate float64) float64 {
	rho := arrivalRate / serviceRate
	x := rho * erlangB(servers-1, float64(servers)*rho)
	return x / ((1 - rho) + x)
}
