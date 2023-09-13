//go:build integration

package experiment

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/hetiansu5/urlquery"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/G-Research/fasttrackml/pkg/api/mlflow"
	"github.com/G-Research/fasttrackml/pkg/api/mlflow/api"
	"github.com/G-Research/fasttrackml/pkg/api/mlflow/api/request"
	"github.com/G-Research/fasttrackml/pkg/api/mlflow/api/response"
	"github.com/G-Research/fasttrackml/pkg/api/mlflow/common"
	"github.com/G-Research/fasttrackml/pkg/api/mlflow/dao/models"
	"github.com/G-Research/fasttrackml/tests/integration/golang/helpers"
)

type GetExperimentByNameTestSuite struct {
	suite.Suite
	helpers.BaseTestSuite
}

func TestGetExperimentByNameTestSuite(t *testing.T) {
	suite.Run(t, new(GetExperimentByNameTestSuite))
}

func (s *GetExperimentByNameTestSuite) SetupTest() {
	s.BaseTestSuite.SetupTest(s.T())
}

func (s *GetExperimentByNameTestSuite) Test_Ok() {
	defer func() {
		assert.Nil(s.T(), s.NamespaceFixtures.UnloadFixtures())
	}()
	// 1. prepare database with test data.
	namespace, err := s.NamespaceFixtures.CreateNamespace(context.Background(), &models.Namespace{
		ID:                  1,
		Code:                "default",
		DefaultExperimentID: common.GetPointer(int32(0)),
	})
	assert.Nil(s.T(), err)

	experiment, err := s.ExperimentFixtures.CreateExperiment(context.Background(), &models.Experiment{
		Name: "Test Experiment",
		Tags: []models.ExperimentTag{
			{
				Key:   "key1",
				Value: "value1",
			},
		},
		NamespaceID: namespace.ID,
		CreationTime: sql.NullInt64{
			Int64: time.Now().UTC().UnixMilli(),
			Valid: true,
		},
		LastUpdateTime: sql.NullInt64{
			Int64: time.Now().UTC().UnixMilli(),
			Valid: true,
		},
		LifecycleStage:   models.LifecycleStageActive,
		ArtifactLocation: "/artifact/location",
	})
	assert.Nil(s.T(), err)

	// 2. make actual API call.
	query, err := urlquery.Marshal(request.GetExperimentRequest{
		Name: experiment.Name,
	})
	assert.Nil(s.T(), err)

	resp := response.GetExperimentResponse{}
	err = s.MlflowClient.DoGetRequest(
		fmt.Sprintf(
			"%s%s?%s", mlflow.ExperimentsRoutePrefix, mlflow.ExperimentsGetByNameRoute, query,
		),
		&resp,
	)
	assert.Nil(s.T(), err)
	// 3. check actual API response.
	assert.Equal(s.T(), fmt.Sprintf("%d", *experiment.ID), resp.Experiment.ID)
	assert.Equal(s.T(), experiment.Name, resp.Experiment.Name)
	assert.Equal(s.T(), string(experiment.LifecycleStage), resp.Experiment.LifecycleStage)
	assert.Equal(s.T(), experiment.ArtifactLocation, resp.Experiment.ArtifactLocation)
	assert.Equal(s.T(), []models.ExperimentTag{
		{
			Key:          "key1",
			Value:        "value1",
			ExperimentID: *experiment.ID,
		},
	}, experiment.Tags)
}

func (s *GetExperimentByNameTestSuite) Test_Error() {
	defer func() {
		assert.Nil(s.T(), s.NamespaceFixtures.UnloadFixtures())
	}()

	_, err := s.NamespaceFixtures.CreateNamespace(context.Background(), &models.Namespace{
		ID:                  1,
		Code:                "default",
		DefaultExperimentID: common.GetPointer(int32(0)),
	})
	assert.Nil(s.T(), err)

	testData := []struct {
		name    string
		error   *api.ErrorResponse
		request *request.GetExperimentRequest
	}{
		{
			name:  "NotFoundExperiment",
			error: api.NewResourceDoesNotExistError(`unable to find experiment 'incorrect_experiment_name'`),
			request: &request.GetExperimentRequest{
				Name: "incorrect_experiment_name",
			},
		},
		{
			name:  "EmptyExperimentName",
			error: api.NewInvalidParameterValueError(`Missing value for required parameter 'experiment_name'`),
			request: &request.GetExperimentRequest{
				Name: "",
			},
		},
	}

	for _, tt := range testData {
		s.T().Run(tt.name, func(t *testing.T) {
			query, err := urlquery.Marshal(tt.request)
			assert.Nil(s.T(), err)
			resp := api.ErrorResponse{}
			err = s.MlflowClient.DoGetRequest(
				fmt.Sprintf("%s%s?%s", mlflow.ExperimentsRoutePrefix, mlflow.ExperimentsGetByNameRoute, query),
				&resp,
			)
			assert.Nil(t, err)
			assert.Equal(s.T(), tt.error.Error(), resp.Error())
		})
	}
}
