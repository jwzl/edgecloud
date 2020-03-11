package devicetwin

import (
	"time"
	"sync"
	"errors"
	"strings"
	"k8s.io/klog"
	"github.com/jwzl/edgeOn/common"
	"github.com/jwzl/wssocket/model"
	"github.com/jwzl/beehive/pkg/core/context"
)


type DTContext struct {
	Context			*context.Context
	/*
	* This is all edge in this cluster.
	*/
	EdgeMap		*sync.Map
	EdgeMutex	*sync.Map
}

func NewDTContext(c *context.Context) *DTContext {
	if c == nil {
		return nil
	}

	var edges sync.Map
	var edgesMutex sync.Map

	return &DTContext{
		Context:	c,
		EdgeMap:	&edges,
		EdgeMutex:  &edgesMutex,
	}
}


//GetEdgeMutex get the edge mutex
func (dtc *DTContext) GetEdgeMutex (edgeID string) (*sync.Mutex, bool) {
	v, exist := dtc.EdgeMutex.Load(edgeID)
	if !exist {
		return nil, false
	}

	mutex, isMutex := v.(*sync.Mutex)
	if !isMutex {
		return nil, false
	}

	return mutex, true
}
