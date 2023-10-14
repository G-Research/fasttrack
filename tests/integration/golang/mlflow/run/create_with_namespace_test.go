//go:build integration

package run

import (
	"github.com/stretchr/testify/suite"

	"github.com/G-Research/fasttrackml/tests/integration/golang/helpers"
)

type CreateRunWithNamespaceTestSuite struct {
	suite.Suite
	helpers.BaseTestSuite
}

// TODO enable once ns-admin-ui is merged
/*
func TestCreateRunWithNamespaceTestSuite(t *testing.T) {
	suite.Run(t, new(CreateRunWithNamespaceTestSuite))
}

func (s *CreateRunWithNamespaceTestSuite) SetupTest() {
	s.BaseTestSuite.SetupTest(s.T())
	s.BaseTestSuite.MlflowClient = helpers.NewClient(helpers.GetServiceUri(), "/ns/custom-ns/api/2.0/mlflow")
}

func (s *CreateRunWithNamespaceTestSuite) Test_Ok() {
	defer func() {
		assert.Nil(s.T(), s.NamespaceFixtures.UnloadFixtures())
	}()

	// create test namespace and experiment.
	namespace, err := s.NamespaceFixtures.CreateNamespace(context.Background(), &models.Namespace{
		ID:                  1,
		Code:                "custom-ns",
		DefaultExperimentID: common.GetPointer(int32(0)),
	})
	assert.Nil(s.T(), err)

	experiment, err := s.ExperimentFixtures.CreateExperiment(context.Background(), &models.Experiment{
		Name:           uuid.New().String(),
		NamespaceID:    namespace.ID,
		LifecycleStage: models.LifecycleStageActive,
	})
	assert.Nil(s.T(), err)

	// set namespace default experiment.
	namespace.DefaultExperimentID = experiment.ID
	_, err = s.NamespaceFixtures.UpdateNamespace(context.Background(), namespace)
	assert.Nil(s.T(), err)

	// send request with "0" for experiment ID, which is default provided by mlflow client.
	req := request.CreateRunRequest{
		Name: "TestRun",
		Tags: []request.RunTagPartialRequest{
			{
				Key:   "key1",
				Value: "value1",
			},
			{
				Key:   "key2",
				Value: "value2",
			},
		},
		StartTime:    1234567890,
		ExperimentID: "0",
	}

	resp := response.CreateRunResponse{}
	assert.Nil(
		s.T(),
		s.MlflowClient.WithMethod(
			http.MethodPost,
		).WithRequest(
			req,
		).WithResponse(
			&resp,
		).DoRequest(
			fmt.Sprintf("%s%s", mlflow.RunsRoutePrefix, mlflow.RunsCreateRoute),
		),
	)
	assert.NotEmpty(s.T(), resp.Run.Info.ID)
	assert.NotEmpty(s.T(), resp.Run.Info.UUID)
	assert.Equal(s.T(), "TestRun", resp.Run.Info.Name)
	assert.Equal(s.T(), fmt.Sprintf("%d", *experiment.ID), resp.Run.Info.ExperimentID)
	assert.Equal(s.T(), int64(1234567890), resp.Run.Info.StartTime)
	assert.Equal(s.T(), int64(0), resp.Run.Info.EndTime)
	assert.Equal(s.T(), string(models.StatusRunning), resp.Run.Info.Status)
	assert.NotEmpty(s.T(), resp.Run.Info.ArtifactURI)
	assert.Equal(s.T(), string(models.LifecycleStageActive), resp.Run.Info.LifecycleStage)
	assert.Equal(s.T(), []response.RunTagPartialResponse{
		{
			Key:   "key1",
			Value: "value1",
		},
		{
			Key:   "key2",
			Value: "value2",
		},
	}, resp.Run.Data.Tags)
}

func (s *CreateRunWithNamespaceTestSuite) Test_Error() {
	_, err := s.NamespaceFixtures.CreateNamespace(context.Background(), &models.Namespace{
		ID:                  1,
		Code:                "custom-ns",
		DefaultExperimentID: common.GetPointer(int32(0)),
	})
	assert.Nil(s.T(), err)

	tests := []struct {
		name    string
		error   *api.ErrorResponse
		request request.CreateRunRequest
	}{
		{
			name: "CreateRunWithInvalidExperimentID",
			request: request.CreateRunRequest{
				ExperimentID: "invalid_experiment_id",
			},
			error: api.NewBadRequestError(
				`unable to parse experiment id 'invalid_experiment_id': strconv.ParseInt: ` +
					`parsing "invalid_experiment_id": invalid syntax`,
			),
		},
		{
			name: "CreateRunWithNotExistingExperiment",
			request: request.CreateRunRequest{
				ExperimentID: "1",
			},
			error: api.NewResourceDoesNotExistError(
				`unable to find experiment for namespace with id '1': error getting experiment by id: 1: record not found`,
			),
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(T *testing.T) {
			resp := api.ErrorResponse{}
			assert.Nil(
				s.T(),
				s.MlflowClient.WithMethod(
					http.MethodPost,
				).WithRequest(
					tt.request,
				).WithResponse(
					&resp,
				).DoRequest(
					fmt.Sprintf("%s%s", mlflow.RunsRoutePrefix, mlflow.RunsCreateRoute),
				),
			)
			assert.Equal(s.T(), tt.error.Error(), resp.Error())
		})
	}
}
*/
