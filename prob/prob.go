package prob

type Prob struct {
	Val    int    // prob value
	DisVal string // prob display value
}

func NewProb(v int, dis string) Prob {
	return Prob{Val: v, DisVal: dis}
}
