//go:build integration

package experiment

import (
	"context"
	"database/sql"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/G-Research/fasttrackml/pkg/api/aim/response"
	"github.com/G-Research/fasttrackml/pkg/api/mlflow/api"
	"github.com/G-Research/fasttrackml/pkg/api/mlflow/common"
	"github.com/G-Research/fasttrackml/pkg/api/mlflow/dao/models"
	"github.com/G-Research/fasttrackml/tests/integration/golang/helpers"
)

type DeleteExperimentTestSuite struct {
	helpers.BaseTestSuite
}

func TestDeleteExperimentTestSuite(t *testing.T) {
	suite.Run(t, new(DeleteExperimentTestSuite))
}

func (s *DeleteExperimentTestSuite) Test_Ok() {
	defer func() {
		require.Nil(s.T(), s.NamespaceFixtures.UnloadFixtures())
	}()

	namespace, err := s.NamespaceFixtures.CreateNamespace(context.Background(), &models.Namespace{
		ID:                  1,
		Code:                "default",
		DefaultExperimentID: common.GetPointer(int32(0)),
	})
	require.Nil(s.T(), err)

	experiment, err := s.ExperimentFixtures.CreateExperiment(context.Background(), &models.Experiment{
		Name: "Test Experiment",
		Tags: []models.ExperimentTag{
			{
				Key:   "key1",
				Value: "value1",
			},
		},
		CreationTime: sql.NullInt64{
			Int64: time.Now().UTC().UnixMilli(),
			Valid: true,
		},
		NamespaceID: namespace.ID,
		LastUpdateTime: sql.NullInt64{
			Int64: time.Now().UTC().UnixMilli(),
			Valid: true,
		},
		LifecycleStage:   models.LifecycleStageActive,
		ArtifactLocation: "/artifact/location",
	})
	require.Nil(s.T(), err)

	experiments, err := s.ExperimentFixtures.GetTestExperiments(context.Background())
	require.Nil(s.T(), err)
	length := len(experiments)

	var resp response.DeleteExperiment
	require.Nil(
		s.T(),
		s.AIMClient.WithMethod(
			http.MethodDelete,
		).WithResponse(
			&resp,
		).DoRequest(
			"/experiments/%d", *experiment.ID,
		),
	)

	remainingExperiments, err := s.ExperimentFixtures.GetTestExperiments(context.Background())
	require.Nil(s.T(), err)
	assert.Equal(s.T(), length-1, len(remainingExperiments))
}

func (s *DeleteExperimentTestSuite) Test_Error() {
	defer func() {
		require.Nil(s.T(), s.NamespaceFixtures.UnloadFixtures())
	}()

	_, err := s.NamespaceFixtures.CreateNamespace(context.Background(), &models.Namespace{
		ID:                  1,
		Code:                "default",
		DefaultExperimentID: common.GetPointer(int32(0)),
	})
	require.Nil(s.T(), err)

	tests := []struct {
		name  string
		ID    string
		error string
	}{
		{
			ID:    "123",
			name:  "DeleteWithUnknownIDFails",
			error: "Not Found",
		},
		{
			name: "DeleteIncorrectExperimentID",
			error: `: unable to parse experiment id "incorrect_experiment_id": strconv.ParseInt:` +
				` parsing "incorrect_experiment_id": invalid syntax`,
			ID: "incorrect_experiment_id",
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(T *testing.T) {
			var resp api.ErrorResponse
			require.Nil(
				s.T(),
				s.AIMClient.WithMethod(
					http.MethodDelete,
				).WithResponse(
					&resp,
				).DoRequest(
					"/experiments/%s", tt.ID,
				),
			)
			assert.Contains(s.T(), resp.Error(), tt.error)
			assert.NoError(s.T(), err)
		})
	}
}
