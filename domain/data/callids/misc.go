package callids

import	"github.com/josephbudd/cwt/domain/types"

var nextid types.CallID

func nextCallID() types.CallID {
	id := nextid
	nextid++
	return id
}

