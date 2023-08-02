//go:build integration

package run

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/hetiansu5/urlquery"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/G-Research/fasttrackml/pkg/api/mlflow"
	"github.com/G-Research/fasttrackml/pkg/api/mlflow/api"
	"github.com/G-Research/fasttrackml/pkg/api/mlflow/api/request"
	"github.com/G-Research/fasttrackml/pkg/api/mlflow/api/response"
	"github.com/G-Research/fasttrackml/pkg/api/mlflow/dao/models"
	"github.com/G-Research/fasttrackml/tests/integration/golang/fixtures"
	"github.com/G-Research/fasttrackml/tests/integration/golang/helpers"
)

type SearchTestSuite struct {
	suite.Suite
	client             *helpers.HttpClient
	runFixtures        *fixtures.RunFixtures
	tagFixtures        *fixtures.TagFixtures
	metricFixtures     *fixtures.MetricFixtures
	experimentFixtures *fixtures.ExperimentFixtures
}

func TestSearchTestSuite(t *testing.T) {
	suite.Run(t, new(SearchTestSuite))
}

func (s *SearchTestSuite) SetupTest() {
	s.client = helpers.NewMlflowApiClient(helpers.GetServiceUri())
	runFixtures, err := fixtures.NewRunFixtures(helpers.GetDatabaseUri())
	assert.Nil(s.T(), err)
	s.runFixtures = runFixtures
	tagFixtures, err := fixtures.NewTagFixtures(helpers.GetDatabaseUri())
	assert.Nil(s.T(), err)
	s.tagFixtures = tagFixtures
	metricFixtures, err := fixtures.NewMetricFixtures(helpers.GetDatabaseUri())
	assert.Nil(s.T(), err)
	s.metricFixtures = metricFixtures
	expFixtures, err := fixtures.NewExperimentFixtures(helpers.GetDatabaseUri())
	assert.Nil(s.T(), err)
	s.experimentFixtures = expFixtures
}

func (s *SearchTestSuite) Test_Ok() {
	defer func() {
		assert.Nil(s.T(), s.experimentFixtures.UnloadFixtures())
	}()

	// create test experiment.
	experiment, err := s.experimentFixtures.CreateExperiment(context.Background(), &models.Experiment{
		Name:           uuid.New().String(),
		LifecycleStage: models.LifecycleStageActive,
	})
	assert.Nil(s.T(), err)

	// create 3 different test runs and attach tags, metrics, params, etc.
	run1, err := s.runFixtures.CreateRun(context.Background(), &models.Run{
		ID:         "id1",
		Name:       "TestRun1",
		UserID:     "1",
		Status:     models.StatusRunning,
		SourceType: "JOB",
		StartTime: sql.NullInt64{
			Int64: 123456789,
			Valid: true,
		},
		EndTime: sql.NullInt64{
			Int64: 123456789,
			Valid: true,
		},
		ExperimentID:   *experiment.ID,
		ArtifactURI:    "artifact_uri1",
		LifecycleStage: models.LifecycleStageActive,
	})
	assert.Nil(s.T(), err)
	_, err = s.tagFixtures.CreateTag(context.Background(), &models.Tag{
		Key:   "mlflow.runName",
		Value: "TestRunTag1",
		RunID: run1.ID,
	})
	assert.Nil(s.T(), err)
	_, err = s.metricFixtures.CreateLatestMetric(context.Background(), &models.LatestMetric{
		Key:       "run1",
		Value:     1.1,
		Timestamp: 1234567890,
		Step:      1,
		IsNan:     false,
		RunID:     run1.ID,
		LastIter:  1,
	})
	assert.Nil(s.T(), err)

	run2, err := s.runFixtures.CreateRun(context.Background(), &models.Run{
		ID:         "id2",
		Name:       "TestRun2",
		UserID:     "2",
		Status:     models.StatusScheduled,
		SourceType: "JOB",
		StartTime: sql.NullInt64{
			Int64: 111111111,
			Valid: true,
		},
		EndTime: sql.NullInt64{
			Int64: 222222222,
			Valid: true,
		},
		ExperimentID:   *experiment.ID,
		ArtifactURI:    "artifact_uri2",
		LifecycleStage: models.LifecycleStageDeleted,
	})
	assert.Nil(s.T(), err)
	_, err = s.tagFixtures.CreateTag(context.Background(), &models.Tag{
		Key:   "mlflow.runName",
		Value: "TestRunTag2",
		RunID: run2.ID,
	})
	assert.Nil(s.T(), err)
	_, err = s.metricFixtures.CreateLatestMetric(context.Background(), &models.LatestMetric{
		Key:       "run2",
		Value:     2.1,
		Timestamp: 1234567890,
		Step:      1,
		IsNan:     false,
		RunID:     run2.ID,
		LastIter:  1,
	})
	assert.Nil(s.T(), err)

	run3, err := s.runFixtures.CreateRun(context.Background(), &models.Run{
		ID:         "id3",
		Name:       "TestRun3",
		UserID:     "3",
		Status:     models.StatusRunning,
		SourceType: "JOB",
		StartTime: sql.NullInt64{
			Int64: 333444444,
			Valid: true,
		},
		EndTime: sql.NullInt64{
			Int64: 444555555,
			Valid: true,
		},
		ExperimentID:   *experiment.ID,
		ArtifactURI:    "artifact_uri3",
		LifecycleStage: models.LifecycleStageActive,
	})
	assert.Nil(s.T(), err)
	_, err = s.tagFixtures.CreateTag(context.Background(), &models.Tag{
		Key:   "mlflow.runName",
		Value: "TestRunTag3",
		RunID: run3.ID,
	})
	assert.Nil(s.T(), err)
	_, err = s.metricFixtures.CreateLatestMetric(context.Background(), &models.LatestMetric{
		Key:       "run3",
		Value:     3.1,
		Timestamp: 1234567890,
		Step:      1,
		IsNan:     false,
		RunID:     run3.ID,
		LastIter:  1,
	})
	assert.Nil(s.T(), err)

	tests := []struct {
		name     string
		error    *api.ErrorResponse
		request  *request.SearchRunsRequest
		response *response.SearchRunsResponse
	}{
		{
			name: "SearchWithViewTypeAllParameter3RunsShouldBeReturned",
			request: &request.SearchRunsRequest{
				ViewType:      request.ViewTypeAll,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run1.ID,
							Name:           "TestRunTag1",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "1",
							Status:         string(models.StatusRunning),
							StartTime:      123456789,
							EndTime:        123456789,
							ArtifactURI:    "artifact_uri1",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
					{
						Info: response.RunInfoPartialResponse{
							ID:             run2.ID,
							Name:           "TestRunTag2",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "2",
							Status:         string(models.StatusScheduled),
							StartTime:      111111111,
							EndTime:        222222222,
							ArtifactURI:    "artifact_uri2",
							LifecycleStage: string(models.LifecycleStageDeleted),
						},
					},
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithViewTypeActiveOnlyParameter2RunsShouldBeReturned",
			request: &request.SearchRunsRequest{
				ViewType:      request.ViewTypeActiveOnly,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run1.ID,
							Name:           "TestRunTag1",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "1",
							Status:         string(models.StatusRunning),
							StartTime:      123456789,
							EndTime:        123456789,
							ArtifactURI:    "artifact_uri1",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithViewTypeDeletedOnlyParameter1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				ViewType:      request.ViewTypeDeletedOnly,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run2.ID,
							Name:           "TestRunTag2",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "2",
							Status:         string(models.StatusScheduled),
							StartTime:      111111111,
							EndTime:        222222222,
							ArtifactURI:    "artifact_uri2",
							LifecycleStage: string(models.LifecycleStageDeleted),
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeStartTimeOperationGrater1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.start_time > 123456789`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeStartTimeOperationGraterOrEqual2RunsShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.start_time >= 123456789`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run1.ID,
							Name:           "TestRunTag1",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "1",
							Status:         string(models.StatusRunning),
							StartTime:      123456789,
							EndTime:        123456789,
							ArtifactURI:    "artifact_uri1",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeStartTimeOperationNotEqual1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.start_time != 123456789`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeStartTimeOperationEqual1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.start_time = 123456789`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run1.ID,
							Name:           "TestRunTag1",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "1",
							Status:         string(models.StatusRunning),
							StartTime:      123456789,
							EndTime:        123456789,
							ArtifactURI:    "artifact_uri1",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeStartTimeOperationLess1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.start_time < 333444444`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run1.ID,
							Name:           "TestRunTag1",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "1",
							Status:         string(models.StatusRunning),
							StartTime:      123456789,
							EndTime:        123456789,
							ArtifactURI:    "artifact_uri1",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeStartTimeOperationLessOrEqual2RunsShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.start_time <= 333444444`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run1.ID,
							Name:           "TestRunTag1",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "1",
							Status:         string(models.StatusRunning),
							StartTime:      123456789,
							EndTime:        123456789,
							ArtifactURI:    "artifact_uri1",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeEndTimeOperationGrater1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.end_time > 123456789`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeEndTimeOperationGraterOrEqual2RunsShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.end_time >= 123456789`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run1.ID,
							Name:           "TestRunTag1",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "1",
							Status:         string(models.StatusRunning),
							StartTime:      123456789,
							EndTime:        123456789,
							ArtifactURI:    "artifact_uri1",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeEndTimeOperationNotEqual1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.end_time != 123456789`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeEndTimeOperationEqual1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.end_time = 123456789`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run1.ID,
							Name:           "TestRunTag1",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "1",
							Status:         string(models.StatusRunning),
							StartTime:      123456789,
							EndTime:        123456789,
							ArtifactURI:    "artifact_uri1",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeEndTimeOperationLess1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.end_time < 444555555`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run1.ID,
							Name:           "TestRunTag1",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "1",
							Status:         string(models.StatusRunning),
							StartTime:      123456789,
							EndTime:        123456789,
							ArtifactURI:    "artifact_uri1",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeEndTimeOperationLessOrEqual2RunsShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.end_time <= 444555555`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run1.ID,
							Name:           "TestRunTag1",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "1",
							Status:         string(models.StatusRunning),
							StartTime:      123456789,
							EndTime:        123456789,
							ArtifactURI:    "artifact_uri1",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeRunNameOperationNotEqual1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.run_name != "TestRunTag1"`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeRunNameOperationEqual1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.run_name = "TestRunTag1"`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run1.ID,
							Name:           "TestRunTag1",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "1",
							Status:         string(models.StatusRunning),
							StartTime:      123456789,
							EndTime:        123456789,
							ArtifactURI:    "artifact_uri1",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeRunNameOperationLike1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.run_name LIKE "TestRunTag1"`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run1.ID,
							Name:           "TestRunTag1",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "1",
							Status:         string(models.StatusRunning),
							StartTime:      123456789,
							EndTime:        123456789,
							ArtifactURI:    "artifact_uri1",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeRunNameOperationILike1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.run_name ILIKE "testruntag1"`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run1.ID,
							Name:           "TestRunTag1",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "1",
							Status:         string(models.StatusRunning),
							StartTime:      123456789,
							EndTime:        123456789,
							ArtifactURI:    "artifact_uri1",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeStatusOperationNotEqualNoRunsShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.status != "RUNNING"`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{},
		},
		{
			name: "SearchWithAttributeStatusOperationEqual2RunsShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.status = "RUNNING"`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run1.ID,
							Name:           "TestRunTag1",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "1",
							Status:         string(models.StatusRunning),
							StartTime:      123456789,
							EndTime:        123456789,
							ArtifactURI:    "artifact_uri1",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeStatusOperationLike2RunsShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.status LIKE "RUNNING"`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run1.ID,
							Name:           "TestRunTag1",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "1",
							Status:         string(models.StatusRunning),
							StartTime:      123456789,
							EndTime:        123456789,
							ArtifactURI:    "artifact_uri1",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeStatusOperationILike2RunsShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.status ILIKE "running"`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run1.ID,
							Name:           "TestRunTag1",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "1",
							Status:         string(models.StatusRunning),
							StartTime:      123456789,
							EndTime:        123456789,
							ArtifactURI:    "artifact_uri1",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeUserIDOperationNotEqual1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.user_id != 1`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeUserIDOperationEqual1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.user_id = 3`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeUserIDOperationLike1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.user_id LIKE "3"`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeUserIDOperationILike1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.user_id ILIKE "3"`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeArtifactURIOperationNotEqual1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.artifact_uri != "artifact_uri1"`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeArtifactURIOperationEqual1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.artifact_uri = "artifact_uri3"`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeArtifactURIOperationLike1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.artifact_uri LIKE "artifact_uri3"`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeArtifactURIOperationILike1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `attributes.artifact_uri ILIKE "ArTiFaCt_UrI3"`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeRunIDOperationNotEqual1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        fmt.Sprintf(`attributes.run_id != "%s"`, run1.ID),
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeRunIDOperationEqual1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        fmt.Sprintf(`attributes.run_id = "%s"`, run3.ID),
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeRunIDOperationLike1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        fmt.Sprintf(`attributes.run_id LIKE "%s"`, run3.ID),
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeRunIDOperationILike1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        fmt.Sprintf(`attributes.run_id ILIKE "%s"`, strings.ToUpper(run3.ID)),
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeRunIDOperationIN1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        fmt.Sprintf(`attributes.run_id IN ('%s')`, run3.ID),
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeRunIDOperationIN1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        fmt.Sprintf(`attributes.run_id NOT IN ('%s')`, run1.ID),
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeMetricsOperationGrater1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `metrics.run3 > 1.1`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeMetricsOperationGraterOrEqual1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `metrics.run3 >= 1.1`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeMetricsOperationNotEqual1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `metrics.run3 != 1.1`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeMetricsOperationEqual1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `metrics.run3 = 3.1`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
				},
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeMetricsOperationLess0RunsShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `metrics.run3 < 3.1`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				NextPageToken: "",
			},
		},
		{
			name: "SearchWithAttributeMetricsOperationLessOrEqual1RunShouldBeReturned",
			request: &request.SearchRunsRequest{
				Filter:        `metrics.run3 <= 3.1`,
				ExperimentIDs: []string{fmt.Sprintf("%d", *experiment.ID)},
			},
			response: &response.SearchRunsResponse{
				Runs: []*response.RunPartialResponse{
					{
						Info: response.RunInfoPartialResponse{
							ID:             run3.ID,
							Name:           "TestRunTag3",
							ExperimentID:   fmt.Sprintf("%d", *experiment.ID),
							UserID:         "3",
							Status:         string(models.StatusRunning),
							StartTime:      333444444,
							EndTime:        444555555,
							ArtifactURI:    "artifact_uri3",
							LifecycleStage: string(models.LifecycleStageActive),
						},
					},
				},
				NextPageToken: "",
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(T *testing.T) {
			query, err := urlquery.Marshal(tt.request)
			assert.Nil(s.T(), err)
			resp := &response.SearchRunsResponse{}
			err = s.client.DoPostRequest(
				fmt.Sprintf("%s%s?%s", mlflow.RunsRoutePrefix, mlflow.RunsSearchRoute, query),
				tt.request,
				&resp,
			)
			assert.Nil(s.T(), err)
			helpers.CompareExpectedSearchRunsResponseWithActualSearchRunsResponse(s.T(), tt.response, resp)
		})
	}
}
