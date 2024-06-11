package repositories

import (
	"context"

	"github.com/rotisserie/eris"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/G-Research/fasttrackml/pkg/api/aim/dao/models"
	"github.com/G-Research/fasttrackml/pkg/common/dao/repositories"
)

// TagRepositoryProvider provides an interface to work with models.Tag entity.
type TagRepositoryProvider interface {
	repositories.BaseRepositoryProvider
	// GetTagsByNamespace returns the list of tags.
	GetTagsByNamespace(ctx context.Context, namespaceID uint) ([]models.TagData, error)
	// CreateExperimentTag creates new models.ExperimentTag entity connected to models.Experiment.
	CreateExperimentTag(ctx context.Context, experimentTag *models.ExperimentTag) error
	// CreateRunTag creates new models.Tag entity connected to models.Run.
	CreateRunTag(ctx context.Context, runTag *models.Tag) error
	// GetTagKeysByParameters returns list of tag keys by requested parameters.
	GetTagKeysByParameters(ctx context.Context, namespaceID uint, experiments []int) ([]string, error)
}

// TagRepository repository to work with models.Tag entity.
type TagRepository struct {
	repositories.BaseRepositoryProvider
}

// NewTagRepository creates repository to work with models.Tag entity.
func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{
		repositories.NewBaseRepository(db),
	}
}

// CreateExperimentTag creates new models.ExperimentTag entity connected to models.Experiment.
func (r TagRepository) CreateExperimentTag(ctx context.Context, experimentTag *models.ExperimentTag) error {
	if err := r.GetDB().WithContext(ctx).Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(experimentTag).Error; err != nil {
		return eris.Wrapf(err, "error creating tag for experiment with id: %d", experimentTag.ExperimentID)
	}
	return nil
}

// GetTagsByNamespace returns the list of TagData, with virtual rows populated from the Tag table.
func (r TagRepository) GetTagsByNamespace(ctx context.Context, namespaceID uint) ([]models.TagData, error) {
	var tagDatas []models.TagData
	if err := r.GetDB().WithContext(ctx).
		Raw(`
                   SELECT *
                   FROM tag_data
                   WHERE is_archived = false
                   AND namespace_id = ?
                   UNION
                   SELECT "00000000-0000-0000-0000-000000000000" as id,
                          'false' as is_archived,
                          tags.key as key,
                          '#cccccc' as color,
                          '' as description,
                          ? as namespace_id
                   FROM tags JOIN runs USING(run_uuid)
                   INNER JOIN experiments ON experiments.experiment_id = runs.experiment_id
                   AND experiments.namespace_id = ?
                   AND key NOT IN (select key from tag_data where is_archived = false and namespace_id = ?)
                   ORDER by key
                 `, namespaceID, namespaceID, namespaceID, namespaceID).
		Preload("Runs").
		Find(&tagDatas).Error; err != nil {
		return nil, eris.Wrap(err, "unable to fetch tag_data")
	}
	return tagDatas, nil
}

// CreateRunTagn creates new models.Tag entity connected to models.Run.
func (r TagRepository) CreateRunTag(ctx context.Context, runTag *models.Tag) error {
	if err := r.GetDB().WithContext(ctx).Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create([]models.Tag{{
		Key:   runTag.Key,
		Value: runTag.Value,
		RunID: runTag.RunID,
	}}).Error; err != nil {
		return eris.Wrapf(err, "error creating tag for run with id: %s", runTag.RunID)
	}
	return nil
	}
}

// GetTagKeysByParameters returns list of tag keys by requested parameters.
func (r TagRepository) GetTagKeysByParameters(
	ctx context.Context, namespaceID uint, experiments []int,
) ([]string, error) {
	// fetch and process tags.
	query := r.GetDB().WithContext(ctx).Model(
		&models.Tag{},
	).Joins(
		"JOIN runs USING(run_uuid)",
	).Joins(
		"INNER JOIN experiments ON experiments.experiment_id = runs.experiment_id AND experiments.namespace_id = ?",
		namespaceID,
	).Where(
		"runs.lifecycle_stage = ?", models.LifecycleStageActive,
	)
	if len(experiments) != 0 {
		query = query.Where("experiments.experiment_id IN ?", experiments)
	}

	var keys []string
	if err := query.Pluck("Key", &keys).Error; err != nil {
		return nil, eris.Wrap(err, "error getting tag keys by parameters")
	}
	return keys, nil
}
