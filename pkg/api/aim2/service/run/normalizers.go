package run

import (
	"github.com/G-Research/fasttrackml/pkg/api/aim2/api/request"
)

// NormaliseGetRunInfoRequest normalizes
func NormaliseGetRunInfoRequest(req *request.GetRunInfoRequest) *request.GetRunInfoRequest {
	if len(req.Sequences) == 0 {
		req.Sequences = SupportedSequences
	}
	return req
}
