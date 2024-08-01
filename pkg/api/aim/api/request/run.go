package request

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// BaseSearchRequest defines some shared fields for search requestes.
type BaseSearchRequest struct {
	ReportProgress bool `query:"report_progress"`
}

// GetRunInfoRequest is a request object for `GET /runs/:id/info` endpoint.
type GetRunInfoRequest struct {
	ID         string   `params:"id"`
	SkipSystem bool     `query:"skip_system"`
	Sequences  []string `query:"sequence"`
}

// GetRunMetricsRequest is a request object for `POST /runs/:id/metric/get-batch` endpoint.
type GetRunMetricsRequest []struct {
	Name    string            `json:"name"`
	Context map[string]string `json:"context"`
}

// GetRunImagesRequest is a request object for `POST /runs/:id/images/get-batch` endpoint.
type GetRunImagesRequest []struct {
	Name    string            `json:"name"`
	Context map[string]string `json:"context"`
}

// GetRunImagesBatchRequest is a request object for `POST /runs/images/get-batch` endpoint.
type GetRunImagesBatchRequest []string

// GetRunsActiveRequest is a request object for `GET /runs/active` endpoint.
type GetRunsActiveRequest struct {
	BaseSearchRequest
}

// UpdateRunRequest is a request struct for `PUT /runs/:id` endpoint.
type UpdateRunRequest struct {
	ID          string  `params:"id"`
	RunID       *string `json:"run_id"`
	RunUUID     *string `json:"run_uuid"`
	Name        *string `json:"run_name"`
	Description *string `json:"description"`
	Status      *string `json:"status"`
	EndTime     *int64  `json:"end_time"`
	Archived    *bool   `json:"archived"`
}

// GetRunLogsRequest is a request struct for `GET /runs/:id/logs` endpoint.
type GetRunLogsRequest struct {
	ID string `params:"id"`
}

// SearchRunsRequest is a request object for `GET /runs/search/run` endpoint.
type SearchRunsRequest struct {
	BaseSearchRequest
	Query           string   `query:"q"`
	Limit           int      `query:"limit"`
	Offset          string   `query:"offset"`
	Action          string   `query:"action"`
	SkipSystem      bool     `query:"skip_system"`
	ExcludeParams   bool     `query:"exclude_params"`
	ExcludeTraces   bool     `query:"exclude_traces"`
	ExperimentNames []string `query:"experiment_names"`
}

// MetricTuple represents a metric with key and context.
type MetricTuple struct {
	Key     string    `json:"key"`
	Context fiber.Map `json:"context"`
}

// SearchMetricsRequest is a request struct for `GET /runs/search/metric` endpoint.
type SearchMetricsRequest struct {
	BaseSearchRequest
	Metrics    []MetricTuple `json:"metrics"`
	Query      string        `json:"query"`
	Steps      int           `json:"steps"`
	XAxis      string        `json:"x_axis"`
	SkipSystem bool          `json:"skip_system"`
}

// SearchAlignedMetricsRequest is a request struct for `GET /runs/search/metric/align` endpoint.
type SearchAlignedMetricsRequest struct {
	Runs []struct {
		ID     string `json:"run_id"`
		Traces []struct {
			Name    string    `json:"name"`
			Slice   [3]int    `json:"slice"`
			Context fiber.Map `json:"context"`
		} `json:"traces"`
	} `json:"runs"`
	AlignBy string `json:"align_by"`
}

// SearchArtifactsRequest is a request struct for `POST /runs/search/image` endpoint.
type SearchArtifactsRequest struct {
	BaseSearchRequest
	Query         string `json:"q"`
	SkipSystem    bool   `json:"skip_system"`
	RecordDensity any    `json:"record_density"`
	IndexDensity  any    `json:"index_density"`
	RecordRange   string `json:"record_range"`
	IndexRange    string `json:"index_range"`
	CalcRanges    bool   `json:"calc_ranges"`
}

// DeleteRunRequest is a request struct for `DELETE /runs/:id` endpoint.
type DeleteRunRequest struct {
	ID string `params:"id"`
}

// ArchiveBatchRequest is a request struct for `DELETE /runs/archive-batch` endpoint.
type ArchiveBatchRequest []string

// DeleteBatchRequest is a request struct for `DELETE /runs/delete-batch` endpoint.
type DeleteBatchRequest []string

// AddRunTagRequest is a request for `POST /runs/:id/tags/new` endpoint.
type AddRunTagRequest struct {
	RunID   string `params:"id"`
	TagName string `json:"tag_name"`
}

// DeleteRunTagRequest is a request for `DELETE /runs/:id/tags/:tagID` endpoint.
type DeleteRunTagRequest struct {
	RunID string `params:"id"`
	TagID string `params:"tagID"`
}

// RecordRangeMin returns the low end of the record range.
func (req SearchArtifactsRequest) RecordRangeMin() int {
	return rangeMin(req.RecordRange)
}

// RecordRangeMax returns the high end of the record range.
func (req SearchArtifactsRequest) RecordRangeMax(dflt int) int {
	return rangeMax(req.RecordRange, dflt)
}

// IndexRangeMin returns the low end of the index range.
func (req SearchArtifactsRequest) IndexRangeMin() int {
	return rangeMin(req.IndexRange)
}

// IndexRangeMax returns the high end of the index range.
func (req SearchArtifactsRequest) IndexRangeMax(dflt int) int {
	return rangeMax(req.IndexRange, dflt)
}

// rangeMin will extract the lower end of a range string in the request.
func rangeMin(r string) int {
	rangeVals := strings.Split(r, ":")
	if len(rangeVals) != 2 {
		return 0
	}
	num, err := strconv.Atoi(rangeVals[0])
	if err == nil {
		return num
	}
	return 0
}

// rangeMax will extract the lower end of a range string in the request.
func rangeMax(r string, dflt int) int {
	rangeVals := strings.Split(r, ":")
	if len(rangeVals) != 2 {
		return dflt
	}
	num, err := strconv.Atoi(rangeVals[1])
	if err == nil {
		return num
	}
	return dflt
}
