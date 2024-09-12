package predicate

import "github.com/PakkaSys/fluidapi/database/util"

var APIToDatabasePredicates = map[Predicate]util.Predicate{
	GREATER:                util.GREATER,
	GREATER_SHORT:          util.GREATER,
	GREATER_OR_EQUAL:       util.GREATER_OR_EQUAL,
	GREATER_OR_EQUAL_SHORT: util.GREATER_OR_EQUAL,
	EQUAL:                  util.EQUAL,
	EQUAL_SHORT:            util.EQUAL,
	NOT_EQUAL:              util.NOT_EQUAL,
	NOT_EQUAL_SHORT:        util.NOT_EQUAL,
	LESS:                   util.LESS,
	LESS_SHORT:             util.LESS,
	LESS_OR_EQUAL:          util.LESS_OR_EQUAL,
	LESS_OR_EQUAL_SHORT:    util.LESS_OR_EQUAL,
	IN:                     util.IN,
	NOT_IN:                 util.NOT_IN,
}
