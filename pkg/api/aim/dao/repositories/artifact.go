package repositories

import (
	"context"
	"database/sql"

	"gorm.io/gorm"

	"github.com/G-Research/fasttrackml/pkg/api/aim/api/request"
	"github.com/G-Research/fasttrackml/pkg/api/aim/dao/models"
	"github.com/G-Research/fasttrackml/pkg/api/aim/query"
	"github.com/G-Research/fasttrackml/pkg/common/dao/repositories"
	"github.com/rotisserie/eris"
)

// ImageSearchStepInfo is a search summary for a Run Step.
type ImageSearchStepInfo struct {
	RunUUID  string `gorm:"column:run_uuid"`
	Step     int    `gorm:"column:step"`
	ImgCount int    `gorm:"column:count"`
}

// ImageSearchSummary is a search summary for whole run.
type ImageSearchSummary map[string][]ImageSearchStepInfo

// TotalSteps figures out how many steps belong to the runID.
func (r ImageSearchSummary) TotalSteps(runID string) int {
	return len(r[runID])
}

// StepImageCount figures out how many steps belong to the runID and step.
func (r ImageSearchSummary) StepImageCount(runID string, step int) int {
	runStepImages := r[runID]
	return runStepImages[step].ImgCount
}

// ArtifactRepositoryProvider provides an interface to work with `artifact` entity.
type ArtifactRepositoryProvider interface {
	repositories.BaseRepositoryProvider
	// Search will find artifacts based on the request.
	Search(
		ctx context.Context,
		namespaceID uint,
		timeZoneOffset int,
		req request.SearchArtifactsRequest,
	) (*sql.Rows, int64, ImageSearchSummary, error)
}

// ArtifactRepository repository to work with `artifact` entity.
type ArtifactRepository struct {
	repositories.BaseRepositoryProvider
}

// NewArtifactRepository creates a repository to work with `artifact` entity.
func NewArtifactRepository(db *gorm.DB) *ArtifactRepository {
	return &ArtifactRepository{
		repositories.NewBaseRepository(db),
	}
}

// Search will find artifacts based on the request.
func (r ArtifactRepository) Search(
	ctx context.Context,
	namespaceID uint,
	timeZoneOffset int,
	req request.SearchArtifactsRequest,
) (*sql.Rows, int64, ImageSearchSummary, error) {
	qp := query.QueryParser{
		Default: query.DefaultExpression{
			Contains:   "run.archived",
			Expression: "not run.archived",
		},
		Tables: map[string]string{
			"runs":        "runs",
			"experiments": "experiments",
		},
		TzOffset:  timeZoneOffset,
		Dialector: r.GetDB().Dialector.Name(),
	}
	pq, err := qp.Parse(req.Query)
	if err != nil {
		return nil, 0, nil, err
	}

	var totalRuns int64
	if err := r.GetDB().WithContext(ctx).Model(&models.Run{}).Count(&totalRuns).Error; err != nil {
		return nil, 0, nil, eris.Wrap(err, "error counting metrics")
	}

	runIDs := []string{}
	if tx := pq.Filter(r.GetDB().WithContext(ctx).
		Select("runs.run_uuid").
		Table("runs").
		Joins(
			"INNER JOIN experiments ON experiments.experiment_id = runs.experiment_id AND experiments.namespace_id = ?",
			namespaceID,
		)).
		Find(&runIDs); tx.Error != nil {
		return nil, 0, nil, eris.Wrap(err, "error finding runs for artifact search")
	}

	resultSummaries := []ImageSearchStepInfo{}
	if tx := r.GetDB().WithContext(ctx).
		Raw(`SELECT run_uuid, step, count(id)
			  FROM artifacts
			  WHERE run_uuid IN (?)
			  GROUP BY run_uuid, step;`, runIDs).
		Find(&resultSummaries); tx.Error != nil {
		return nil, 0, nil, eris.Wrap(err, "error find result summary for artifact search")
	}

	runImages := make(ImageSearchSummary, len(runIDs))
	for _, rslt := range resultSummaries {
		runImages[rslt.RunUUID] = append(runImages[rslt.RunUUID], rslt)
	}

	tx := r.GetDB().WithContext(ctx).
		Table("artifacts").
		Where("run_uuid IN ?", runIDs).
		Order("run_uuid").
		Order("step").
		Order("created_at")

	rows, err := tx.Rows()
	if err != nil {
		return nil, 0, nil, eris.Wrap(err, "error searching artifacts")
	}
	if err := rows.Err(); err != nil {
		return nil, 0, nil, eris.Wrap(err, "error getting artifacts rows cursor")
	}

	return rows, int64(len(runIDs)), runImages, nil
}
