package controller

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/G-Research/fasttrackml/pkg/api/aim2/api/request"

	"github.com/gofiber/fiber/v2"
	"github.com/rotisserie/eris"
	log "github.com/sirupsen/logrus"

	"github.com/G-Research/fasttrackml/pkg/common/api"
	"github.com/G-Research/fasttrackml/pkg/common/middleware/namespace"
	"github.com/G-Research/fasttrackml/pkg/database"
)

func (c Controller) GetProject(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"name":              "FastTrackML",
		"path":              database.DB.Dialector.Name(),
		"description":       "",
		"telemetry_enabled": 0,
	})
}

func (c Controller) GetProjectActivity(ctx *fiber.Ctx) error {
	ns, err := namespace.GetNamespaceFromContext(ctx.Context())
	if err != nil {
		return api.NewInternalError("error getting namespace from context")
	}
	log.Debugf("getProjectActivity namespace: %s", ns.Code)

	tzOffset, err := strconv.Atoi(ctx.Get("x-timezone-offset", "0"))
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "x-timezone-offset header is not a valid integer")
	}

	var numExperiments int64
	if tx := database.DB.Model(
		&database.Experiment{},
	).Where(
		"lifecycle_stage = ?", database.LifecycleStageActive,
	).Where(
		"namespace_id = ?", ns.ID,
	).Count(&numExperiments); tx.Error != nil {
		return fmt.Errorf("error counting experiments: %w", tx.Error)
	}

	var runs []database.Run
	if tx := database.DB.Select(
		"runs.status",
		"runs.start_time",
		"runs.lifecycle_stage",
	).Joins(
		"INNER JOIN experiments ON experiments.experiment_id = runs.experiment_id AND experiments.namespace_id = ?",
		ns.ID,
	).Find(
		&runs,
	); tx.Error != nil {
		return fmt.Errorf("error retrieving runs: %w", tx.Error)
	}

	numArchivedRuns := 0
	numActiveRuns := 0
	activity := map[string]int{}
	for _, r := range runs {
		key := time.UnixMilli(r.StartTime.Int64).Add(time.Duration(-tzOffset) * time.Minute).Format("2006-01-02T15:00:00")
		activity[key] += 1
		switch {
		case r.LifecycleStage == database.LifecycleStageDeleted:
			numArchivedRuns += 1
		case r.Status == database.StatusRunning:
			numActiveRuns += 1
		}
	}

	return ctx.JSON(fiber.Map{
		"num_runs":          len(runs),
		"activity_map":      activity,
		"num_active_runs":   numActiveRuns,
		"num_experiments":   numExperiments,
		"num_archived_runs": numArchivedRuns,
	})
}

// TODO
func (c Controller) GetProjectPinnedSequences(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"sequences": []string{},
	})
}

// TODO
func (c Controller) UpdateProjectPinnedSequences(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"sequences": []string{},
	})
}

func (c Controller) GetProjectParams(ctx *fiber.Ctx) error {
	ns, err := namespace.GetNamespaceFromContext(ctx.Context())
	if err != nil {
		return api.NewInternalError("error getting namespace from context")
	}
	log.Debugf("getProjectParams namespace: %s", ns.Code)

	req := request.GetProjectParamsRequest{}
	if err := ctx.QueryParser(&req); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	resp := fiber.Map{}

	if !req.ExcludeParams {
		// fetch and process params.
		query := database.DB.Distinct().Model(
			&database.Param{},
		).Joins(
			"JOIN runs USING(run_uuid)",
		).Joins(
			"INNER JOIN experiments ON experiments.experiment_id = runs.experiment_id AND experiments.namespace_id = ?",
			ns.ID,
		).Where(
			"runs.lifecycle_stage = ?", database.LifecycleStageActive,
		)
		if len(req.Experiments) != 0 {
			query.Where("experiments.experiment_id IN ?", req.Experiments)
		}
		var paramKeys []string
		if err = query.Pluck("Key", &paramKeys).Error; err != nil {
			return fmt.Errorf("error retrieving param keys: %w", err)
		}

		params := make(map[string]any, len(paramKeys)+1)
		for _, p := range paramKeys {
			params[p] = map[string]string{
				"__example_type__": "<class 'str'>",
			}
		}

		// fetch and process tags.
		query = database.DB.Distinct().Model(
			&database.Tag{},
		).Joins(
			"JOIN runs USING(run_uuid)",
		).Joins(
			"INNER JOIN experiments ON experiments.experiment_id = runs.experiment_id AND experiments.namespace_id = ?",
			ns.ID,
		).Where(
			"runs.lifecycle_stage = ?", database.LifecycleStageActive,
		)
		if len(req.Experiments) != 0 {
			query.Where("experiments.experiment_id IN ?", req.Experiments)
		}
		var tagKeys []string
		if err = query.Pluck("Key", &tagKeys).Error; err != nil {
			return fmt.Errorf("error retrieving tag keys: %w", err)
		}

		tags := make(map[string]map[string]string, len(tagKeys))
		for _, t := range tagKeys {
			tags[t] = map[string]string{
				"__example_type__": "<class 'str'>",
			}
		}

		params["tags"] = tags
		resp["params"] = params
	}

	if len(req.Sequences) == 0 {
		req.Sequences = []string{
			"metric",
			"images",
			"texts",
			"figures",
			"distributions",
			"audios",
		}
	}

	for _, s := range req.Sequences {
		switch s {
		case "images", "texts", "figures", "distributions", "audios":
			resp[s] = fiber.Map{}
		case "metric":
			query := database.DB.Distinct().Model(
				&database.LatestMetric{},
			).Joins(
				"JOIN runs USING(run_uuid)",
			).Joins(
				"INNER JOIN experiments ON experiments.experiment_id = runs.experiment_id AND experiments.namespace_id = ?",
				ns.ID,
			).Joins(
				"Context",
			).Where(
				"runs.lifecycle_stage = ?", database.LifecycleStageActive,
			)
			if len(req.Experiments) != 0 {
				query.Where("experiments.experiment_id IN ?", req.Experiments)
			}
			var metrics []database.LatestMetric
			if err = query.Find(&metrics).Error; err != nil {
				return fmt.Errorf("error retrieving metric keys: %w", err)
			}

			data, mapped := make(map[string][]fiber.Map, len(metrics)), make(map[string]map[string]fiber.Map, len(metrics))
			for _, metric := range metrics {
				if mapped[metric.Key] == nil {
					mapped[metric.Key] = map[string]fiber.Map{}
				}
				if _, ok := mapped[metric.Key][metric.Context.GetJsonHash()]; !ok {
					// to be properly decoded by AIM UI, json should be represented as a key:value object.
					context := fiber.Map{}
					if err := json.Unmarshal(metric.Context.Json, &context); err != nil {
						return eris.Wrap(err, "error unmarshalling `context` json to `fiber.Map` object")
					}
					mapped[metric.Key][metric.Context.GetJsonHash()] = context
					data[metric.Key] = append(data[metric.Key], context)
				}
			}
			resp[s] = data
		default:
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("%q is not a valid Sequence", s))
		}
	}

	return ctx.JSON(resp)
}

func (c Controller) GetProjectStatus(ctx *fiber.Ctx) error {
	return ctx.JSON("up-to-date")
}
