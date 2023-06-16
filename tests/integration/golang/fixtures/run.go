package fixtures

import (
	"context"
	"time"

	"github.com/rotisserie/eris"

	"github.com/G-Research/fasttrackml/pkg/api/mlflow/dao/models"
	"github.com/G-Research/fasttrackml/pkg/api/mlflow/dao/repositories"
	"github.com/G-Research/fasttrackml/pkg/database"
)

// RunFixtures represents data fixtures object.
type RunFixtures struct {
	baseFixtures
	runRepository        repositories.RunRepositoryProvider
}

// NewRunFixtures creates new instance of RunFixtures.
func NewRunFixtures(databaseDSN string) (*RunFixtures, error) {
	db, err := database.ConnectDB(
		databaseDSN,
		1*time.Second,
		20,
		false,
		false,
		"",
	)
	if err != nil {
		return nil, eris.Wrap(err, "error connection to database")
	}
	return &RunFixtures{
		baseFixtures:  baseFixtures{db: db},
		runRepository: repositories.NewRunRepository(db),
	}, nil
}

// CreateTestRun creates a new test Run.
func (f RunFixtures) CreateTestRun(
	ctx context.Context, run *models.Run,
) (*models.Run, error) {
	if err := f.runRepository.Create(ctx, run); err != nil {
		return nil, eris.Wrap(err, "error creating test run")
	}
	return run, nil
}

// GetTestRuns fetches all runs for an experiment
func (f RunFixtures) GetTestRuns(
	ctx context.Context, experimentID int32) ([]models.Run, error) {
	return f.runRepository.List(ctx, experimentID)
}

// FindMinMaxRowNums finds min and max rownum for an experiment's runs
func (f RunFixtures) FindMinMaxRowNums(
	ctx context.Context, experimentID int32) (int64, int64, error) {
	runs, err := f.runRepository.List(ctx, experimentID)
	if err != nil {
		return 0, 0, eris.Wrap(err, "error fetching test runs")
	}
	var min, max models.RowNum
	for _, run := range runs {
		if run.RowNum < min {
			min = run.RowNum
		}
		if run.RowNum > max {
			max = run.RowNum
		}
	}
	return int64(min), int64(max), nil
}
