package transfernsq

import (
	"time"
)

func NewHandler(uc TransferUsecase, baseRequeueDelay, maxRequeueDelay time.Duration) Handler {
	if baseRequeueDelay <= 0 {
		baseRequeueDelay = 2 * time.Second
	}
	if maxRequeueDelay <= 0 {
		maxRequeueDelay = 10 * time.Second
	}
	return &handler{
		usecase:          uc,
		baseRequeueDelay: baseRequeueDelay,
		maxRequeueDelay:  maxRequeueDelay,
	}
}
