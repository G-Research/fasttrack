package artifact

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/G-Research/fasttrackml/pkg/api/mlflow"
	"github.com/G-Research/fasttrackml/pkg/api/mlflow/api"
	"github.com/G-Research/fasttrackml/pkg/api/mlflow/api/request"
	"github.com/G-Research/fasttrackml/pkg/api/mlflow/api/response"
	"github.com/G-Research/fasttrackml/pkg/common/dao/models"
	"github.com/G-Research/fasttrackml/tests/integration/golang/helpers"
)

type ListArtifactGSTestSuite struct {
	helpers.GSTestSuite
}

func TestListArtifactGSTestSuite(t *testing.T) {
	suite.Run(t, &ListArtifactGSTestSuite{
		helpers.NewGSTestSuite("bucket1", "bucket2"),
	})
}

func (s *ListArtifactGSTestSuite) Test_Ok() {
	tests := []struct {
		name   string
		bucket string
	}{
		{
			name:   "TestWithBucket1",
			bucket: "bucket1",
		},
		{
			name:   "TestWithBucket2",
			bucket: "bucket2",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			// 1. create test experiment.
			experiment, err := s.ExperimentFixtures.CreateExperiment(context.Background(), &models.Experiment{
				Name:             fmt.Sprintf("Test Experiment In Bucket %s", tt.bucket),
				NamespaceID:      s.DefaultNamespace.ID,
				LifecycleStage:   models.LifecycleStageActive,
				ArtifactLocation: fmt.Sprintf("gs://%s/1", tt.bucket),
			})
			s.Require().Nil(err)

			// 2. create test run.
			runID := strings.ReplaceAll(uuid.New().String(), "-", "")
			run, err := s.RunFixtures.CreateRun(context.Background(), &models.Run{
				ID:             runID,
				Status:         models.StatusRunning,
				SourceType:     "JOB",
				ExperimentID:   *experiment.ID,
				ArtifactURI:    fmt.Sprintf("%s/%s/artifacts", experiment.ArtifactLocation, runID),
				LifecycleStage: models.LifecycleStageActive,
			})
			s.Require().Nil(err)

			// 3. upload artifact objects to GS.
			writer := s.Client.Bucket(
				tt.bucket,
			).Object(
				fmt.Sprintf("1/%s/artifacts/artifact.txt", runID),
			).NewWriter(
				context.Background(),
			)
			_, err = writer.Write([]byte("contentX"))
			s.Require().Nil(err)
			s.Require().Nil(writer.Close())

			writer = s.Client.Bucket(
				tt.bucket,
			).Object(
				fmt.Sprintf("1/%s/artifacts/artifact/artifact.txt", runID),
			).NewWriter(
				context.Background(),
			)
			_, err = writer.Write([]byte("contentXX"))
			s.Require().Nil(err)
			s.Require().Nil(writer.Close())

			// 4. make actual API call for root dir.
			rootDirQuery := request.ListArtifactsRequest{
				RunID: run.ID,
			}

			rootDirResp := response.ListArtifactsResponse{}
			s.Require().Nil(
				s.MlflowClient().WithQuery(
					rootDirQuery,
				).WithResponse(
					&rootDirResp,
				).DoRequest(
					"%s%s", mlflow.ArtifactsRoutePrefix, mlflow.ArtifactsListRoute,
				),
			)

			s.Equal(run.ArtifactURI, rootDirResp.RootURI)
			s.Equal(2, len(rootDirResp.Files))
			s.Equal([]response.FilePartialResponse{
				{
					Path:     "artifact",
					IsDir:    true,
					FileSize: 0,
				},
				{
					Path:     "artifact.txt",
					IsDir:    false,
					FileSize: 8,
				},
			}, rootDirResp.Files)
			s.Require().Nil(err)

			// 5. make actual API call for sub dir.
			subDirQuery := request.ListArtifactsRequest{
				RunID: run.ID,
				Path:  "artifact",
			}

			subDirResp := response.ListArtifactsResponse{}
			s.Require().Nil(
				s.MlflowClient().WithQuery(
					subDirQuery,
				).WithResponse(
					&subDirResp,
				).DoRequest(
					"%s%s", mlflow.ArtifactsRoutePrefix, mlflow.ArtifactsListRoute,
				),
			)

			s.Equal(run.ArtifactURI, subDirResp.RootURI)
			s.Equal(1, len(subDirResp.Files))
			s.Equal(response.FilePartialResponse{
				Path:     "artifact/artifact.txt",
				IsDir:    false,
				FileSize: 9,
			}, subDirResp.Files[0])
			s.Require().Nil(err)

			// 6. make actual API call for non-existing dir.
			nonExistingDirQuery := request.ListArtifactsRequest{
				RunID: run.ID,
				Path:  "non-existing-dir",
			}
			s.Require().Nil(err)

			nonExistingDirResp := response.ListArtifactsResponse{}
			s.Require().Nil(
				s.MlflowClient().WithQuery(
					nonExistingDirQuery,
				).WithResponse(
					&nonExistingDirResp,
				).DoRequest(
					"%s%s", mlflow.ArtifactsRoutePrefix, mlflow.ArtifactsListRoute,
				),
			)

			s.Equal(run.ArtifactURI, nonExistingDirResp.RootURI)
			s.Equal(0, len(nonExistingDirResp.Files))
			s.Require().Nil(err)
		})
	}
}

func (s *ListArtifactGSTestSuite) Test_Error() {
	tests := []struct {
		name    string
		error   *api.ErrorResponse
		request request.ListArtifactsRequest
	}{
		{
			name:    "EmptyOrIncorrectRunIDOrRunUUID",
			error:   api.NewInvalidParameterValueError("Missing value for required parameter 'run_id'"),
			request: request.ListArtifactsRequest{},
		},
		{
			name:  "IncorrectPathProvidedCase1",
			error: api.NewInvalidParameterValueError("Invalid path"),
			request: request.ListArtifactsRequest{
				RunID: "run_id",
				Path:  "..",
			},
		},
		{
			name:  "IncorrectPathProvidedCase2",
			error: api.NewInvalidParameterValueError("Invalid path"),
			request: request.ListArtifactsRequest{
				RunID: "run_id",
				Path:  "./..",
			},
		},
		{
			name:  "IncorrectPathProvidedCase3",
			error: api.NewInvalidParameterValueError("Invalid path"),
			request: request.ListArtifactsRequest{
				RunID: "run_id",
				Path:  "./../",
			},
		},
		{
			name:  "IncorrectPathProvidedCase4",
			error: api.NewInvalidParameterValueError("Invalid path"),
			request: request.ListArtifactsRequest{
				RunID: "run_id",
				Path:  "foo/../bar",
			},
		},
		{
			name:  "IncorrectPathProvidedCase5",
			error: api.NewInvalidParameterValueError("Invalid path"),
			request: request.ListArtifactsRequest{
				RunID: "run_id",
				Path:  "/foo/../bar",
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			resp := api.ErrorResponse{}
			s.Require().Nil(
				s.MlflowClient().WithQuery(
					tt.request,
				).WithResponse(
					&resp,
				).DoRequest(
					"%s%s", mlflow.ArtifactsRoutePrefix, mlflow.ArtifactsListRoute,
				),
			)
			s.Equal(tt.error.Error(), resp.Error())
		})
	}
}
