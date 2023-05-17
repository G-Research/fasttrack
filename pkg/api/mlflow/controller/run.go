package controller

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"

	"github.com/G-Research/fasttrackml/pkg/api/mlflow/api"
	"github.com/G-Research/fasttrackml/pkg/api/mlflow/api/request"
	"github.com/G-Research/fasttrackml/pkg/api/mlflow/api/response"
	"github.com/G-Research/fasttrackml/pkg/api/mlflow/dao/models"
)

// CreateRun handles `POST /runs/create` endpoint.
func (c Controller) CreateRun(ctx *fiber.Ctx) error {
	var req request.CreateRunRequest
	if err := ctx.BodyParser(&req); err != nil {
		return api.NewBadRequestError("Unable to decode request body: %s", err)
	}

	log.Debugf("Create request: %#v", &req)
	run, err := c.runService.CreateRun(ctx.Context(), &req)
	if err != nil {
		return err
	}
	resp := response.NewCreateRunResponse(run)
	log.Debugf("Create response: %#v", resp)

	return ctx.JSON(resp)
}

// UpdateRun handles `POST /runs/update` endpoint.
func (c Controller) UpdateRun(ctx *fiber.Ctx) error {
	var req request.UpdateRunRequest
	if err := ctx.BodyParser(&req); err != nil {
		return api.NewBadRequestError("Unable to decode request body: %s", err)
	}

	run, err := c.runService.UpdateRun(ctx.Context(), &req)
	if err != nil {
		return err
	}
	log.Debugf("UpdateRun request: %#v", req)
	resp := response.NewUpdateRunResponse(run)

	log.Debugf("UpdateRun response: %#v", resp)

	return ctx.JSON(resp)
}

// GetRun handles `GET /runs/get` endpoint.
func (c Controller) GetRun(ctx *fiber.Ctx) error {
	req := request.GetRunRequest{}
	if err := ctx.QueryParser(&req); err != nil {
		return api.NewBadRequestError(err.Error())
	}

	log.Debugf("GetRun request: %#v", req)
	run, err := c.runService.GetRun(ctx.Context(), &req)
	if err != nil {
		return err
	}

	resp := response.NewGetRunResponse(run)

	log.Debugf("GetRun response: %#v", resp)

	return ctx.JSON(resp)
}

// SearchRuns handles `POST /runs/search` endpoint.
func (c Controller) SearchRuns(ctx *fiber.Ctx) error {
	var req request.SearchRunsRequest
	if err := ctx.BodyParser(&req); err != nil {
		return api.NewBadRequestError("Unable to decode request body: %s", err)
	}
	log.Debugf("SearchRuns request: %#v", req)

	runs, limit, offset, err := c.runService.SearchRuns(ctx.Context(), &req)
	if err != nil {
		return err
	}

	resp, err := response.NewSearchRunsResponse(runs, limit, offset)
	if err != nil {
		return api.NewInternalError("Unable to build next_page_token: %s", err)
	}

	log.Debugf("SearchRuns response: %#v", resp)

	return ctx.JSON(resp)
}

// DeleteRun handles `POST /runs/delete` endpoint.
func (c Controller) DeleteRun(ctx *fiber.Ctx) error {
	var req request.DeleteRunRequest
	if err := ctx.BodyParser(&req); err != nil {
		return api.NewBadRequestError("Unable to decode request body: %s", err)
	}
	log.Debugf("DeleteRun request: %#v", req)

	if err := c.runService.DeleteRun(ctx.Context(), &req); err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{})
}

// RestoreRun handles `POST /runs/restore` endpoint.
func (c Controller) RestoreRun(ctx *fiber.Ctx) error {
	var req request.RestoreRunRequest
	if err := ctx.BodyParser(&req); err != nil {
		return api.NewBadRequestError("Unable to decode request body: %s", err)
	}
	log.Debugf("RestoreRun request: %#v", req)

	if err := c.runService.RestoreRun(ctx.Context(), &req); err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{})
}

// LogMetric handles `POST /runs/log-metric` endpoint.
func (c Controller) LogMetric(ctx *fiber.Ctx) error {
	var req request.LogMetricRequest
	if err := ctx.BodyParser(&req); err != nil {
		if err, ok := err.(*json.UnmarshalTypeError); ok {
			return api.NewInvalidParameterValueError("Invalid value for parameter '%s' supplied. Hint: Value was of type '%s'. See the API docs for more information about request parameters.", err.Field, err.Value)
		}
		return api.NewBadRequestError("Unable to decode request body: %s", err)
	}
	log.Debugf("LogMetric request: %#v", req)

	if err := c.runService.LogMetric(ctx.Context(), &req); err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{})
}

// LogParam handles `POST /runs/log-parameter` endpoint.
func (c Controller) LogParam(ctx *fiber.Ctx) error {
	var req request.LogParamRequest
	if err := ctx.BodyParser(&req); err != nil {
		if err, ok := err.(*json.UnmarshalTypeError); ok {
			return api.NewInvalidParameterValueError("Invalid value for parameter '%s' supplied. Hint: Value was of type '%s'. See the API docs for more information about request parameters.", err.Field, err.Value)
		}
		return api.NewBadRequestError("Unable to decode request body: %s", err)
	}
	log.Debugf("LogParam request: %#v", req)

	if err := c.runService.LogParam(ctx.Context(), &req); err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{})
}

// SetRunTag handles `POST /runs/set-tag` endpoint.
func (c Controller) SetRunTag(ctx *fiber.Ctx) error {
	var req request.SetRunTagRequest
	if err := ctx.BodyParser(&req); err != nil {
		if err, ok := err.(*json.UnmarshalTypeError); ok {
			return api.NewInvalidParameterValueError("Invalid value for parameter '%s' supplied. Hint: Value was of type '%s'. See the API docs for more information about request parameters.", err.Field, err.Value)
		}
		return api.NewBadRequestError("Unable to decode request body: %s", err)
	}
	log.Debugf("SetRunTag request: %#v", req)

	if err := c.runService.SetRunTag(ctx.Context(), &req); err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{})
}

// DeleteRunTag handles `POST /runs/delete-tag` endpoint.
func (c Controller) DeleteRunTag(ctx *fiber.Ctx) error {
	var req request.DeleteRunTagRequest
	if err := ctx.BodyParser(&req); err != nil {
		return api.NewBadRequestError("Unable to decode request body: %s", err)
	}
	log.Debugf("DeleteRunTag request: %#v", req)

	if err := c.runService.DeleteRunTag(ctx.Context(), &req); err != nil {
		return err
	}
	return ctx.JSON(fiber.Map{})
}

// LogBatch handles `POST /runs/log-batch` endpoint.
func (c Controller) LogBatch(ctx *fiber.Ctx) error {
	var req request.LogBatchRequest
	if err := ctx.BodyParser(&req); err != nil {
		if err, ok := err.(*json.UnmarshalTypeError); ok {
			return api.NewInvalidParameterValueError("Invalid value for parameter '%s' supplied. Hint: Value was of type '%s'. See the API docs for more information about request parameters.", err.Field, err.Value)
		}
		return api.NewBadRequestError("Unable to decode request body: %s", err)
	}
	log.Debugf("LogBatch request: %#v", req)

	if err := c.runService.LogBatch(ctx.Context(), &req); err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{})
}

// TODO:Dsuhinin lets keep it here for now.
func modelRunToAPI(r *models.Run) response.RunPartialResponse {
	metrics := make([]response.RunMetricPartialResponse, len(r.LatestMetrics))
	for n, m := range r.LatestMetrics {
		metrics[n] = response.RunMetricPartialResponse{
			Key:       m.Key,
			Value:     m.Value,
			Timestamp: m.Timestamp,
			Step:      m.Step,
		}
		if m.IsNan {
			metrics[n].Value = "NaN"
		}
	}

	params := make([]response.RunParamPartialResponse, len(r.Params))
	for n, p := range r.Params {
		params[n] = response.RunParamPartialResponse{
			Key:   p.Key,
			Value: p.Value,
		}
	}

	tags := make([]response.RunTagPartialResponse, len(r.Tags))
	for n, t := range r.Tags {
		tags[n] = response.RunTagPartialResponse{
			Key:   t.Key,
			Value: t.Value,
		}
		switch t.Key {
		case "mlflow.runName":
			r.Name = t.Value
		case "mlflow.user":
			r.UserID = t.Value
		}
	}

	return response.RunPartialResponse{
		Info: response.RunInfoPartialResponse{
			ID:             r.ID,
			UUID:           r.ID,
			Name:           r.Name,
			ExperimentID:   fmt.Sprint(r.ExperimentID),
			UserID:         r.UserID,
			Status:         string(r.Status),
			StartTime:      r.StartTime.Int64,
			EndTime:        r.EndTime.Int64,
			ArtifactURI:    r.ArtifactURI,
			LifecycleStage: string(r.LifecycleStage),
		},
		Data: response.RunDataPartialResponse{
			Metrics: metrics,
			Params:  params,
			Tags:    tags,
		},
	}
}
