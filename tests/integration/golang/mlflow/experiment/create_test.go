//go:build integration

package experiment

import (
	"context"
	"fmt"
	"net/http"
	"testing"

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

type CreateExperimentTestSuite struct {
	suite.Suite
	helpers.BaseTestSuite
}

func TestCreateExperimentTestSuite(t *testing.T) {
	suite.Run(t, new(CreateExperimentTestSuite))
}

func (s *CreateExperimentTestSuite) SetupTest() {
	s.BaseTestSuite.SetupTest(s.T())
}

func (s *CreateExperimentTestSuite) Test_Ok() {
	defer func() {
		assert.Nil(s.T(), s.NamespaceFixtures.UnloadFixtures())
	}()

	_, err := s.NamespaceFixtures.CreateNamespace(context.Background(), &models.Namespace{
		ID:                  1,
		Code:                "default",
		DefaultExperimentID: common.GetPointer(int32(0)),
	})
	assert.Nil(s.T(), err)

	req := request.CreateExperimentRequest{
		Name:             "ExperimentName",
		ArtifactLocation: "/artifact/location",
		Tags: []request.ExperimentTagPartialRequest{
			{
				Key:   "key1",
				Value: "value1",
			},
			{
				Key:   "key2",
				Value: "value2",
			},
		},
	}
	resp := response.CreateExperimentResponse{}
	assert.Nil(
		s.T(),
		s.MlflowClient.WithMethod(
			http.MethodPost,
		).WithRequest(
			req,
		).WithResponse(
			&resp,
		).DoRequest(
			fmt.Sprintf("%s%s", mlflow.ExperimentsRoutePrefix, mlflow.ExperimentsCreateRoute),
		),
	)
	assert.NotEmpty(s.T(), resp.ID)
}

func (s *CreateExperimentTestSuite) Test_Error() {
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
		request *request.CreateExperimentRequest
	}{
		{
			name:    "EmptyNameProperty",
			error:   api.NewInvalidParameterValueError("Missing value for required parameter 'name'"),
			request: &request.CreateExperimentRequest{},
		},
		{
			name: "EmptyArtifactLocationProperty",
			error: api.NewInvalidParameterValueError(
				`Invalid value for parameter 'artifact_location': error parsing artifact location: parse "incorrect-protocol,:/incorrect-location": first path segment in URL cannot contain colon`,
			),
			request: &request.CreateExperimentRequest{
				Name:             "name",
				ArtifactLocation: "incorrect-protocol,:/incorrect-location",
			},
		},
	}

	for _, tt := range testData {
		s.T().Run(tt.name, func(t *testing.T) {
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
					fmt.Sprintf("%s%s", mlflow.ExperimentsRoutePrefix, mlflow.ExperimentsCreateRoute),
				),
			)
			assert.Equal(s.T(), tt.error.Error(), resp.Error())
		})
	}
}
