package stats

type label struct {
	name     string
	negative bool
}

func newLabel(name string) label {
	if name == "" {
		panic("empty label name")
	}
	if name[0] == '-' {
		return label{name: name[1:], negative: true}
	}
	return label{name: name, negative: false}
}

type labelGroup struct {
	labels  []label
	labeled map[string]bool
}

func newLabelGroup(names []string) labelGroup {
	if len(names) == 0 {
		panic("empty label group")
	}
	res := labelGroup{labeled: map[string]bool{}}
	for _, name := range names {
		res.labels = append(res.labels, newLabel(name))
	}
	return res
}

func (lg *labelGroup) areWeConcernedByThisLabel(name string) bool {
	res := false
	for _, l := range lg.labels {
		if l.name == name {
			res = true
			break
		}
	}
	return res
}

func (lg *labelGroup) label(name string) bool {
	if !lg.areWeConcernedByThisLabel(name) {
		return false
	}
	lg.labeled[name] = true
	return true
}

func (lg *labelGroup) unlabel(name string) bool {
	if !lg.areWeConcernedByThisLabel(name) {
		return false
	}
	lg.labeled[name] = false
	return true
}

func (lg *labelGroup) isFlagged() bool {
	res := true
	for _, l := range lg.labels {
		if l.negative {
			// the label must be absent
			labeled := lg.labeled[l.name]
			if labeled {
				res = false
				break
			}
		} else {
			// the label must be present
			labeled := lg.labeled[l.name]
			if !labeled {
				res = false
				break
			}
		}
	}
	return res
}
