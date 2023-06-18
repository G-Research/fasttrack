//go:build integration

package run

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/G-Research/fasttrackml/pkg/api/mlflow/api"
	"github.com/G-Research/fasttrackml/pkg/api/mlflow/api/request"
	"github.com/G-Research/fasttrackml/pkg/api/mlflow/dao/models"
	"github.com/G-Research/fasttrackml/tests/integration/golang/fixtures"
	"github.com/G-Research/fasttrackml/tests/integration/golang/helpers"
)

type DeleteRunTestSuite struct {
	suite.Suite
	client             *helpers.HttpClient
	runFixtures        *fixtures.RunFixtures
	experimentFixtures *fixtures.ExperimentFixtures
	runs               []*models.Run
}

func TestDeleteRunTestSuite(t *testing.T) {
	suite.Run(t, new(DeleteRunTestSuite))
}

func (s *DeleteRunTestSuite) SetupTest() {
	s.client = helpers.NewAimApiClient(os.Getenv("SERVICE_BASE_URL"))

	runFixtures, err := fixtures.NewRunFixtures(os.Getenv("DATABASE_DSN"))
	assert.Nil(s.T(), err)
	s.runFixtures = runFixtures

	expFixtures, err := fixtures.NewExperimentFixtures(os.Getenv("DATABASE_DSN"))
	assert.Nil(s.T(), err)
	s.experimentFixtures = expFixtures

	exp := &models.Experiment{
		Name:           uuid.New().String(),
		LifecycleStage: models.LifecycleStageActive,
	}
	_, err = s.experimentFixtures.CreateTestExperiment(context.Background(), exp)
	assert.Nil(s.T(), err)

	s.runs, err = s.runFixtures.CreateTestRuns(context.Background(), exp, 10)
	assert.Nil(s.T(), err)
}

func (s *DeleteRunTestSuite) Test_Ok() {
	defer func() {
		assert.Nil(s.T(), s.runFixtures.UnloadFixtures())
		assert.Nil(s.T(), s.experimentFixtures.UnloadFixtures())
	}()
	tests := []struct {
		name    string
		request request.DeleteRunRequest
	}{
		{
			name:    "DeleteOneRunSucceeds",
			request: request.DeleteRunRequest{RunID: s.runs[4].ID},
		},
		{
			name:    "RowNumbersAreRecalculated",
			request: request.DeleteRunRequest{RunID: s.runs[1].ID},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(T *testing.T) {
			originalMinRowNum, originalMaxRowNum, err := s.runFixtures.FindMinMaxRowNums(context.Background(), s.runs[0].ExperimentID)
			assert.NoError(s.T(), err)

			var resp any
			err = s.client.DoDeleteRequest(
				fmt.Sprintf("/%s/%s", "runs", tt.request.RunID),
				&resp,
			)
			assert.Nil(s.T(), err)

			newMinRowNum, newMaxRowNum, err := s.runFixtures.FindMinMaxRowNums(context.Background(), s.runs[0].ExperimentID)
			assert.NoError(s.T(), err)
			assert.Equal(s.T(), originalMinRowNum, newMinRowNum)
			assert.Greater(s.T(), originalMaxRowNum, newMaxRowNum)
		})
	}
}

func (s *DeleteRunTestSuite) Test_Error() {
	defer func() {
		assert.Nil(s.T(), s.runFixtures.UnloadFixtures())
		assert.Nil(s.T(), s.experimentFixtures.UnloadFixtures())
	}()
	tests := []struct {
		name    string
		request request.DeleteRunRequest
	}{
		{
			name:    "DeleteWithUnknownIDFails",
			request: request.DeleteRunRequest{RunID: "some-other-id"},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(T *testing.T) {
			originalMinRowNum, originalMaxRowNum, err := s.runFixtures.FindMinMaxRowNums(context.Background(), s.runs[0].ExperimentID)
			assert.NoError(s.T(), err)

			var resp api.ErrorResponse
			err = s.client.DoDeleteRequest(
				fmt.Sprintf("/%s/%s", "runs", tt.request.RunID),
				&resp,
			)
			assert.Nil(s.T(), err)
			assert.Equal(s.T(), "", resp.Error())

			newMinRowNum, newMaxRowNum, err := s.runFixtures.FindMinMaxRowNums(context.Background(), s.runs[0].ExperimentID)
			assert.NoError(s.T(), err)
			assert.Equal(s.T(), originalMinRowNum, newMinRowNum)
			assert.Equal(s.T(), originalMaxRowNum, newMaxRowNum)
		})
	}
}

// func (s *DeleteBatchTestSuite) Test_Error() {
// 	defer func() {
// 		assert.Nil(s.T(), s.runFixtures.UnloadFixtures())
// 	}()

// var testData = []struct {
// 	name    string
// 	error   *api.ErrorResponse
// 	request *request.DeleteBatchRequest
// }{
// 	{
// 		name:    "MissingRunIDFails",
// 		error:   api.NewInvalidParameterValueError("Missing value for required parameter 'run_id'"),
// 		request: &request.DeleteBatchRequest{},
// 	},
// 	{
// 		name:  "DuplicateKeyDifferentValueFails",
// 		error: api.NewInternalError("duplicate key"),
// 		request: &request.DeleteBatchRequest{
// 			RunID: s.run.ID,
// 			Params: []request.ParamPartialRequest{
// 				{
// 					Key:   "key1",
// 					Value: "value1",
// 				},
// 				{
// 					Key:   "key1",
// 					Value: "value2",
// 				},
// 			},
// 		},
// 	},
// }

// for _, tt := range testData {
// 	s.T().Run(tt.name, func(t *testing.T) {
// 		resp := api.ErrorResponse{}
// 		err := s.client.DoPostRequest(
// 			fmt.Sprintf("%s%s", mlflow.RunsRoutePrefix, mlflow.RunsDeleteBatchRoute),
// 			tt.request,
// 			&resp,
// 		)
// 		assert.NoError(t, err)
// 		assert.Equal(s.T(), tt.error.ErrorCode, resp.ErrorCode)
// 		assert.Contains(s.T(), resp.Error(), tt.error.Message)
// 	})
// }
// }
