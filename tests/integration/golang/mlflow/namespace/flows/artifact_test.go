//go:build integration

package flows

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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

type ArtifactFlowTestSuite struct {
	suite.Suite
	helpers.BaseTestSuite
	testBuckets []string
	s3Client    *s3.Client
}

// TestArtifactFlowTestSuite tests the full `artifact` flow connected to namespace functionality.
// Flow contains next endpoints:
// - `GET /artifacts/get`
// - `GET /artifacts/list`
func TestArtifactFlowTestSuite(t *testing.T) {
	suite.Run(t, &ArtifactFlowTestSuite{
		testBuckets: []string{"bucket1", "bucket2"},
	})
}

func (s *ArtifactFlowTestSuite) SetupTest() {
	s3Client, err := helpers.NewS3Client(helpers.GetS3EndpointUri())
	assert.Nil(s.T(), err)
	s.s3Client = s3Client
}

func (s *ArtifactFlowTestSuite) TearDownTest() {
	assert.Nil(s.T(), s.NamespaceFixtures.UnloadFixtures())
}

func (s *ArtifactFlowTestSuite) Test_Ok() {
	tests := []struct {
		name           string
		setup          func() (*models.Namespace, *models.Namespace)
		namespace1Code string
		namespace2Code string
	}{
		{
			name: "TestCustomNamespaces",
			setup: func() (*models.Namespace, *models.Namespace) {
				return &models.Namespace{
						Code:                "namespace-1",
						DefaultExperimentID: common.GetPointer(int32(0)),
					}, &models.Namespace{
						Code:                "namespace-2",
						DefaultExperimentID: common.GetPointer(int32(0)),
					}
			},
			namespace1Code: "namespace-1",
			namespace2Code: "namespace-2",
		},
		{
			name: "TestObviousDefaultCustomNamespaces",
			setup: func() (*models.Namespace, *models.Namespace) {
				return &models.Namespace{
						Code:                "default",
						DefaultExperimentID: common.GetPointer(int32(0)),
					}, &models.Namespace{
						Code:                "namespace-1",
						DefaultExperimentID: common.GetPointer(int32(0)),
					}
			},
			namespace1Code: "default",
			namespace2Code: "namespace-1",
		},
		{
			name: "TestImplicitDefaultCustomNamespaces",
			setup: func() (*models.Namespace, *models.Namespace) {
				return &models.Namespace{
						Code:                "default",
						DefaultExperimentID: common.GetPointer(int32(0)),
					}, &models.Namespace{
						Code:                "namespace-1",
						DefaultExperimentID: common.GetPointer(int32(0)),
					}
			},
			namespace1Code: "",
			namespace2Code: "namespace-1",
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(T *testing.T) {
			defer func() {
				assert.Nil(s.T(), s.NamespaceFixtures.UnloadFixtures())
				assert.Nil(s.T(), helpers.RemoveS3Buckets(s.s3Client, s.testBuckets))
			}()

			// setup data under the test.
			namespace1, namespace2 := tt.setup()
			namespace1, err := s.NamespaceFixtures.CreateNamespace(context.Background(), namespace1)
			assert.Nil(s.T(), err)
			namespace2, err = s.NamespaceFixtures.CreateNamespace(context.Background(), namespace2)
			assert.Nil(s.T(), err)

			experiment1, err := s.ExperimentFixtures.CreateExperiment(context.Background(), &models.Experiment{
				Name:             "Experiment1",
				ArtifactLocation: "s3://bucket1/1",
				LifecycleStage:   models.LifecycleStageActive,
				NamespaceID:      namespace1.ID,
			})
			assert.Nil(s.T(), err)

			experiment2, err := s.ExperimentFixtures.CreateExperiment(context.Background(), &models.Experiment{
				Name:             "Experiment2",
				ArtifactLocation: "s3://bucket2/2",
				LifecycleStage:   models.LifecycleStageActive,
				NamespaceID:      namespace2.ID,
			})
			assert.Nil(s.T(), err)

			// create test buckets.
			assert.Nil(s.T(), helpers.CreateS3Buckets(s.s3Client, s.testBuckets))

			// run actual flow test over the test data.
			s.testRunArtifactFlow(tt.namespace1Code, tt.namespace2Code, experiment1, experiment2)
		})
	}
}

func (s *ArtifactFlowTestSuite) testRunArtifactFlow(
	namespace1Code, namespace2Code string, experiment1, experiment2 *models.Experiment,
) {
	// create runs and upload test artifacts
	run1ID := s.createRun(namespace1Code, &request.CreateRunRequest{
		Name:         "Run1",
		ExperimentID: fmt.Sprintf("%d", *experiment1.ID),
	})

	_, err := s.s3Client.PutObject(context.Background(), &s3.PutObjectInput{
		Key:    aws.String(fmt.Sprintf("1/%s/artifacts/artifact1.file", run1ID)),
		Body:   strings.NewReader("content1"),
		Bucket: aws.String("bucket1"),
	})
	assert.Nil(s.T(), err)

	run2ID := s.createRun(namespace2Code, &request.CreateRunRequest{
		Name:         "Run2",
		ExperimentID: fmt.Sprintf("%d", *experiment2.ID),
	})

	_, err = s.s3Client.PutObject(context.Background(), &s3.PutObjectInput{
		Key:    aws.String(fmt.Sprintf("2/%s/artifacts/artifact2.file", run2ID)),
		Body:   strings.NewReader("content2"),
		Bucket: aws.String("bucket2"),
	})
	assert.Nil(s.T(), err)

	// test `GET /artifacts/list` endpoint.
	s.listRunArtifactsAndCompare(namespace1Code, request.ListArtifactsRequest{
		RunID: run1ID,
	}, []response.FilePartialResponse{
		{
			Path:     "artifact1.file",
			IsDir:    false,
			FileSize: 8,
		},
	})

	s.listRunArtifactsAndCompare(namespace2Code, request.ListArtifactsRequest{
		RunID: run2ID,
	}, []response.FilePartialResponse{
		{
			Path:     "artifact2.file",
			IsDir:    false,
			FileSize: 8,
		},
	})

	// test `GET /artifacts/list` endpoint.
	// check that there is no intersection between runs, so when we request
	// run 1 in scope of namespace 2 and run 2 in scope of namespace 1 API will throw an error.
	resp := api.ErrorResponse{}
	assert.Nil(
		s.T(),
		s.MlflowClient.WithMethod(
			http.MethodGet,
		).WithNamespace(
			namespace2Code,
		).WithQuery(
			request.ListArtifactsRequest{
				RunID: run1ID,
			},
		).WithResponseType(
			helpers.ResponseTypeJSON,
		).WithResponse(
			&resp,
		).DoRequest(
			fmt.Sprintf("%s%s", mlflow.ArtifactsRoutePrefix, mlflow.ArtifactsListRoute),
		),
	)
	assert.Equal(s.T(), fmt.Sprintf("RESOURCE_DOES_NOT_EXIST: unable to find run '%s'", run1ID), resp.Error())
	assert.Equal(s.T(), api.ErrorCodeResourceDoesNotExist, string(resp.ErrorCode))

	resp = api.ErrorResponse{}
	assert.Nil(
		s.T(),
		s.MlflowClient.WithMethod(
			http.MethodGet,
		).WithNamespace(
			namespace1Code,
		).WithQuery(
			request.ListArtifactsRequest{
				RunID: run2ID,
			},
		).WithResponseType(
			helpers.ResponseTypeJSON,
		).WithResponse(
			&resp,
		).DoRequest(
			fmt.Sprintf("%s%s", mlflow.ArtifactsRoutePrefix, mlflow.ArtifactsListRoute),
		),
	)
	assert.Equal(s.T(), fmt.Sprintf("RESOURCE_DOES_NOT_EXIST: unable to find run '%s'", run2ID), resp.Error())
	assert.Equal(s.T(), api.ErrorCodeResourceDoesNotExist, string(resp.ErrorCode))
}

func (s *ArtifactFlowTestSuite) createRun(namespace string, req *request.CreateRunRequest) string {
	resp := response.CreateRunResponse{}
	assert.Nil(
		s.T(),
		s.MlflowClient.WithMethod(
			http.MethodPost,
		).WithNamespace(
			namespace,
		).WithRequest(
			req,
		).WithResponseType(
			helpers.ResponseTypeJSON,
		).WithResponse(
			&resp,
		).DoRequest(
			fmt.Sprintf("%s%s", mlflow.RunsRoutePrefix, mlflow.RunsCreateRoute),
		),
	)
	return resp.Run.Info.ID
}

func (s *ArtifactFlowTestSuite) listRunArtifactsAndCompare(
	namespace string, req request.ListArtifactsRequest, expectedResponse []response.FilePartialResponse,
) {
	actualResponse := response.ListArtifactsResponse{}
	assert.Nil(
		s.T(),
		s.MlflowClient.WithMethod(
			http.MethodGet,
		).WithNamespace(
			namespace,
		).WithQuery(
			req,
		).WithResponseType(
			helpers.ResponseTypeJSON,
		).WithResponse(
			&actualResponse,
		).DoRequest(
			fmt.Sprintf("%s%s", mlflow.ArtifactsRoutePrefix, mlflow.ArtifactsListRoute),
		),
	)
	assert.Equal(s.T(), expectedResponse, actualResponse.Files)
}
