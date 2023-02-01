package agtwo

import "fmt"

type QueryFilter struct {
	Query string
	Args  []any
}

func (qf *QueryFilter) And(query string, args ...any) *QueryFilter {
	if qf.Query == "" {
		qf.Query = query
	} else {
		if query != "" {
			qf.Query = qf.Query + " AND " + query
		}
	}

	qf.addArgs(args...)
	return qf
}

func (qf *QueryFilter) Or(query string, args ...any) *QueryFilter {
	if qf.Query == "" {
		qf.Query = query
	} else {
		if query != "" {
			qf.Query = fmt.Sprintf("(%s) OR (%s)", qf.Query, query)
		}
	}
	qf.addArgs(args...)
	return qf
}
func (qf *QueryFilter) addArgs(args ...any) {
	qf.Args = append(qf.Args, args...)
}
