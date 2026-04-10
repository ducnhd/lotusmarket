package portfolio

type ReturnSegment struct {
	StartValue float64
	EndValue   float64
}

func CompoundReturns(segments []ReturnSegment) float64 {
	if len(segments) == 0 {
		return 0
	}
	product := 1.0
	for _, seg := range segments {
		if seg.StartValue == 0 {
			return 0
		}
		r := (seg.EndValue - seg.StartValue) / seg.StartValue
		product *= (1 + r)
	}
	return product - 1
}
