//go:build integration

package run

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/G-Research/fasttrackml/pkg/api/mlflow"
	"github.com/G-Research/fasttrackml/pkg/api/mlflow/api"
	"github.com/G-Research/fasttrackml/pkg/api/mlflow/api/request"
	"github.com/G-Research/fasttrackml/pkg/api/mlflow/common"
	"github.com/G-Research/fasttrackml/pkg/common/dao/models"
	"github.com/G-Research/fasttrackml/tests/integration/golang/helpers"
)

type DeleteRunTestSuite struct {
	suite.Suite
	helpers.BaseTestSuite
}

func TestDeleteRunTestSuite(t *testing.T) {
	suite.Run(t, new(DeleteRunTestSuite))
}

func (s *DeleteRunTestSuite) SetupTest() {
	s.BaseTestSuite.SetupTest(s.T())
}

func (s *DeleteRunTestSuite) Test_Ok() {
	defer func() {
		assert.Nil(s.T(), s.NamespaceFixtures.UnloadFixtures())
	}()

	// create experiment
	namespace, err := s.NamespaceFixtures.CreateNamespace(context.Background(), &models.Namespace{
		ID:                  0,
		Code:                "default",
		DefaultExperimentID: common.GetPointer(int32(0)),
	})
	assert.Nil(s.T(), err)

	experiment, err := s.ExperimentFixtures.CreateExperiment(context.Background(), &models.Experiment{
		Name:           uuid.New().String(),
		NamespaceID:    namespace.ID,
		LifecycleStage: models.LifecycleStageActive,
	})
	assert.Nil(s.T(), err)

	// create run for the experiment
	run, err := s.RunFixtures.CreateRun(context.Background(), &models.Run{
		ID:             strings.ReplaceAll(uuid.New().String(), "-", ""),
		Name:           "TestRun",
		Status:         models.StatusRunning,
		SourceType:     "JOB",
		ExperimentID:   *experiment.ID,
		LifecycleStage: models.LifecycleStageActive,
	})
	assert.Nil(s.T(), err)

	tests := []struct {
		name    string
		request request.DeleteRunRequest
	}{
		{
			name:    "DeleteRunSucceedsWithExistingRunID",
			request: request.DeleteRunRequest{RunID: run.ID},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(T *testing.T) {
			resp := map[string]any{}
			err := s.MlflowClient.DoPostRequest(
				fmt.Sprintf("%s%s", mlflow.RunsRoutePrefix, mlflow.RunsDeleteRoute),
				tt.request,
				&resp,
			)
			assert.Nil(s.T(), err)
			assert.Empty(s.T(), resp)

			archivedRuns, err := s.RunFixtures.GetRuns(context.Background(), run.ExperimentID)

			assert.Nil(s.T(), err)
			assert.Equal(T, 1, len(archivedRuns))
			assert.Equal(s.T(), run.ID, archivedRuns[0].ID)
			assert.Equal(s.T(), models.LifecycleStageDeleted, archivedRuns[0].LifecycleStage)
		})
	}
}

func (s *DeleteRunTestSuite) Test_Error() {
	_, err := s.NamespaceFixtures.CreateNamespace(context.Background(), &models.Namespace{
		ID:                  0,
		Code:                "default",
		DefaultExperimentID: common.GetPointer(int32(0)),
	})
	assert.Nil(s.T(), err)

	tests := []struct {
		name    string
		request request.DeleteRunRequest
	}{
		{
			name:    "DeleteRunFailsWithNonExistingRunID",
			request: request.DeleteRunRequest{RunID: "not-an-id"},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(T *testing.T) {
			resp := api.ErrorResponse{}
			err := s.MlflowClient.DoPostRequest(
				fmt.Sprintf("%s%s", mlflow.RunsRoutePrefix, mlflow.RunsDeleteRoute),
				tt.request,
				&resp,
			)
			assert.Nil(s.T(), err)
			assert.Equal(s.T(), "RESOURCE_DOES_NOT_EXIST: unable to find run 'not-an-id': error getting 'run' entity by id: not-an-id: record not found", resp.Error())
		})
	}
}
