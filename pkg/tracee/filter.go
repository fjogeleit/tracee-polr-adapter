package tracee

type Filter struct {
	minumumSeverity int
	excludeRules    []string
}

func (f *Filter) Check(event Event) bool {
	for _, rule := range f.excludeRules {
		if rule == event.SigMetadata.ID {
			return false
		}
	}

	return event.SigMetadata.Severity >= f.minumumSeverity
}

func NewFilter(minumumSeverity int, excludeRules []string) *Filter {
	return &Filter{minumumSeverity, excludeRules}
}
