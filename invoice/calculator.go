package invoice

type PerformanceCalculator interface {
	Amount() float64
	VolumeCredits() int
}

func NewPerformanceCalculator(aPerformance Performance, aPlay Play) PerformanceCalculator {
	switch aPlay.PlayType {
	case "tragedy":
		return &tragedyCalculator{aPerformance, aPlay}
	case "comedy":
		return &comedyCalculator{aPerformance, aPlay}
	default:
		panic("unknown type: " + aPlay.PlayType)
	}
}

type comedyCalculator struct {
	performance Performance
	play        Play
}

func (c comedyCalculator) Amount() float64 {
	result := 30000.0
	if c.performance.Audience > 20 {
		result += 10000 + 500*float64(c.performance.Audience-20)
	}
	result += 300 * float64(c.performance.Audience)
	return result
}

func (c comedyCalculator) VolumeCredits() int {
	return Max(c.performance.Audience-30, 0) + +c.performance.Audience/5
}

type tragedyCalculator struct {
	performance Performance
	play        Play
}

func (t tragedyCalculator) Amount() float64 {
	result := 40000.0
	if t.performance.Audience > 30 {
		result += 1000 * float64(t.performance.Audience-30)
	}
	return result
}

func (t tragedyCalculator) VolumeCredits() int {
	return Max(t.performance.Audience-30, 0)
}
