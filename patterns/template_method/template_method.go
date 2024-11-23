package templatemethod

type SumTerm func() int

type Summer struct {
	first  SumTerm
	second SumTerm
}

func NewSummer(f, s SumTerm) *Summer {
	return &Summer{f, s}
}

func (s *Summer) Sum() int {
	return s.first() + s.second()
}
