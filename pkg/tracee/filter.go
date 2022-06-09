package tracee

type Filter struct {
	minumumSeverity int
}

func (f *Filter) Check(event Event) bool {
	return event.SigMetadata.Severity >= f.minumumSeverity
}

func NewFilter(minumumSeverity int) *Filter {
	return &Filter{minumumSeverity}
}
