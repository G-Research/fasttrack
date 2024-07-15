package database

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"github.com/G-Research/fasttrackml/pkg/api/mlflow/common"
	"github.com/G-Research/fasttrackml/pkg/api/mlflow/dao/models"
	"github.com/G-Research/fasttrackml/pkg/database"
	"github.com/G-Research/fasttrackml/tests/integration/golang/fixtures"
	"github.com/G-Research/fasttrackml/tests/integration/golang/helpers"
)

type rowCounts struct {
	namespaces               int
	experiments              int
	runs                     int
	distinctRunExperimentIDs int
	metrics                  int
	latestMetrics            int
	tags                     int
	params                   int
	dashboards               int
	apps                     int
}

type ImportTestSuite struct {
	suite.Suite
	runs              []*models.Run
	inputRunFixtures  *fixtures.RunFixtures
	outputRunFixtures *fixtures.RunFixtures
	inputBackend      string
	outputBackend     string
	inputDB           *gorm.DB
	outputDB          *gorm.DB
}

func TestImportTestSuite(t *testing.T) {
	suite.Run(t, new(ImportTestSuite))
}

func (s *ImportTestSuite) SetupSubTest() {
	// prepare input database.
	dsn, err := helpers.GenerateDatabaseURI(s.T(), s.inputBackend)
	s.Require().Nil(err)
	db, err := database.NewDBProvider(
		dsn,
		1*time.Second,
		20,
	)
	s.Require().Nil(err)
	s.Require().Nil(database.CheckAndMigrateDB(true, db.GormDB()))
	s.Require().Nil(database.CreateDefaultNamespace(db.GormDB()))
	s.Require().Nil(database.CreateDefaultExperiment(db.GormDB(), "s3://fasttrackml"))
	s.inputDB = db.GormDB()

	inputRunFixtures, err := fixtures.NewRunFixtures(db.GormDB())
	s.Require().Nil(err)
	s.inputRunFixtures = inputRunFixtures
	s.populateDB(s.inputDB)

	// prepare output database.
	dsn, err = helpers.GenerateDatabaseURI(s.T(), s.outputBackend)
	s.Require().Nil(err)
	db, err = database.NewDBProvider(
		dsn,
		1*time.Second,
		20,
	)
	s.Require().Nil(err)
	s.Require().Nil(database.CheckAndMigrateDB(true, db.GormDB()))
	s.Require().Nil(database.CreateDefaultNamespace(db.GormDB()))
	s.Require().Nil(database.CreateDefaultExperiment(db.GormDB(), "s3://fasttrackml"))
	s.outputDB = db.GormDB()

	outputRunFixtures, err := fixtures.NewRunFixtures(db.GormDB())
	s.Require().Nil(err)
	s.outputRunFixtures = outputRunFixtures
}

func (s *ImportTestSuite) populateDB(db *gorm.DB) {
	experimentFixtures, err := fixtures.NewExperimentFixtures(db)
	s.Require().Nil(err)

	runFixtures, err := fixtures.NewRunFixtures(db)
	s.Require().Nil(err)

	namespaceFixtures, err := fixtures.NewNamespaceFixtures(db)
	s.Require().Nil(err)

	namespace, err := namespaceFixtures.CreateNamespace(context.Background(), &models.Namespace{
		Code:                "source-namespace",
		DefaultExperimentID: common.GetPointer(models.DefaultExperimentID),
	})

	// experiment 1
	experiment, err := experimentFixtures.CreateExperiment(context.Background(), &models.Experiment{
		Name:           uuid.New().String(),
		NamespaceID:    1,
		LifecycleStage: models.LifecycleStageActive,
	})
	s.Require().Nil(err)

	runs, err := runFixtures.CreateExampleRuns(context.Background(), experiment, 5)
	s.Require().Nil(err)
	s.runs = runs

	// experiment 2
	experiment, err = experimentFixtures.CreateExperiment(context.Background(), &models.Experiment{
		Name:           uuid.New().String(),
		NamespaceID:    1,
		LifecycleStage: models.LifecycleStageActive,
	})
	s.Require().Nil(err)

	runs, err = runFixtures.CreateExampleRuns(context.Background(), experiment, 5)
	s.Require().Nil(err)
	s.runs = runs

	// experiment 3
	experiment, err = experimentFixtures.CreateExperiment(context.Background(), &models.Experiment{
		Name:           uuid.New().String(),
		NamespaceID:    namespace.ID,
		LifecycleStage: models.LifecycleStageActive,
	})
	s.Require().Nil(err)

	runs, err = runFixtures.CreateExampleRuns(context.Background(), experiment, 5)
	s.Require().Nil(err)
	s.runs = runs

	dashboardFixtures, err := fixtures.NewDashboardFixtures(db)
	s.Require().Nil(err)

	// dashboard 1
	_, err = dashboardFixtures.CreateDashboard(context.Background(), &database.Dashboard{
		App: database.App{
			Type:        "mpi",
			State:       database.AppState{},
			NamespaceID: 1,
		},
		Name: uuid.NewString(),
	})
	s.Require().Nil(err)

	// dashboard 2
	_, err = dashboardFixtures.CreateDashboard(context.Background(), &database.Dashboard{
		App: database.App{
			Type:        "mpi",
			State:       database.AppState{},
			NamespaceID: 1,
		},
		Name: uuid.NewString(),
	})
	s.Require().Nil(err)

	// dashboard 3
	_, err = dashboardFixtures.CreateDashboard(context.Background(), &database.Dashboard{
		App: database.App{
			Type:        "mpi",
			State:       database.AppState{},
			NamespaceID: namespace.ID,
		},
		Name: uuid.NewString(),
	})
	s.Require().Nil(err)
}

func (s *ImportTestSuite) TearDownSubTest() {
	s.Require().Nil(s.inputRunFixtures.TruncateTables())
	s.Require().Nil(s.outputRunFixtures.TruncateTables())
}

func (s *ImportTestSuite) TestGlobalScopeImport_Ok() {

	backends := []string{"sqlite", "sqlcipher", "postgres"}
	for _, inputBackend := range backends {
		for _, outputBackend := range backends {
			s.inputBackend = inputBackend
			s.outputBackend = outputBackend
			s.Run(inputBackend+"->"+outputBackend, func() {
				// source DB should have expected
				s.validateRowCounts(s.inputDB, rowCounts{
					namespaces:               2,
					experiments:              4,
					runs:                     15,
					distinctRunExperimentIDs: 3,
					metrics:                  60,
					latestMetrics:            30,
					tags:                     15,
					params:                   30,
					dashboards:               3,
					apps:                     3,
				})

				// initially, dest DB is empty
				s.validateRowCounts(s.outputDB, rowCounts{namespaces: 1, experiments: 1})

				// invoke the Importer.Import() method
				importer := database.NewImporter(s.inputDB, s.outputDB)
				s.Require().Nil(importer.Import())

				// dest DB should now have the expected
				s.validateRowCounts(s.outputDB, rowCounts{
					namespaces:               2,
					experiments:              4,
					runs:                     15,
					distinctRunExperimentIDs: 3,
					metrics:                  60,
					latestMetrics:            30,
					tags:                     15,
					params:                   30,
					dashboards:               3,
					apps:                     3,
				})

				// invoke the Importer.Import method a 2nd time
				s.Require().Nil(importer.Import())

				// dest DB should still only have the expected (idempotent)
				s.validateRowCounts(s.outputDB, rowCounts{
					namespaces:               2,
					experiments:              4,
					runs:                     15,
					distinctRunExperimentIDs: 3,
					metrics:                  60,
					latestMetrics:            30,
					tags:                     15,
					params:                   30,
					dashboards:               3,
					apps:                     3,
				})

				// confirm row-for-row equality
				for _, table := range []string{
					"namespaces",
					"apps",
					"dashboards",
					"experiment_tags",
					"runs",
					"tags",
					"params",
					"metrics",
					"latest_metrics",
				} {
					s.validateTable(s.inputDB, s.outputDB, table)
				}
			})
		}
	}
}

func (s *ImportTestSuite) TestNamespaceScopeImport_Ok() {
	backends := []string{"sqlite", "sqlcipher", "postgres"}
	for _, inputBackend := range backends {
		for _, outputBackend := range backends {
			s.inputBackend = inputBackend
			s.outputBackend = outputBackend
			s.Run(inputBackend+"->"+outputBackend, func() {
				namespaceFixtures, err := fixtures.NewNamespaceFixtures(s.outputDB)
				s.Require().Nil(err)

				namespace, err := namespaceFixtures.CreateNamespace(context.Background(), &models.Namespace{
					Code:                "destination-namespace",
					DefaultExperimentID: common.GetPointer(models.DefaultExperimentID),
				})
				s.Require().Nil(err)

				// source DB should have expected
				s.validateRowCounts(s.inputDB, rowCounts{
					namespaces:               2,
					experiments:              4,
					runs:                     15,
					distinctRunExperimentIDs: 3,
					metrics:                  60,
					latestMetrics:            30,
					tags:                     15,
					params:                   30,
					dashboards:               3,
					apps:                     3,
				})

				// initially, dest DB is empty
				s.validateRowCounts(s.outputDB, rowCounts{namespaces: 2, experiments: 1})

				// invoke the Importer.Import() method
				importer := database.NewImporter(
					s.inputDB,
					s.outputDB,
					database.WithSourceNamespace("source-namespace"),
					database.WithDestinationNamespace(namespace.Code),
				)
				s.Require().Nil(importer.Import())

				// dest DB should now have the expected
				s.validateRowCounts(s.outputDB, rowCounts{
					namespaces:               2,
					experiments:              2,
					runs:                     5,
					distinctRunExperimentIDs: 1,
					metrics:                  20,
					latestMetrics:            10,
					tags:                     5,
					params:                   10,
					dashboards:               1,
					apps:                     1,
				})

				// invoke the Importer.Import method a 2nd time
				s.Require().Nil(importer.Import())

				// dest DB should still only have the expected (idempotent)
				s.validateRowCounts(s.outputDB, rowCounts{
					namespaces:               2,
					experiments:              2,
					runs:                     5,
					distinctRunExperimentIDs: 1,
					metrics:                  20,
					latestMetrics:            10,
					tags:                     5,
					params:                   10,
					dashboards:               1,
					apps:                     1,
				})
			})
		}
	}
}

// validateRowCounts will make assertions about the db based on the test setup.
// a db imported from the test setup db should also pass these
// assertions.
func (s *ImportTestSuite) validateRowCounts(db *gorm.DB, counts rowCounts) {
	var countVal int64
	s.Require().Nil(db.Model(&models.Namespace{}).Count(&countVal).Error)
	s.Equal(counts.namespaces, int(countVal), "Namespaces count incorrect")

	s.Require().Nil(db.Model(&models.Experiment{}).Count(&countVal).Error)
	s.Equal(counts.experiments, int(countVal), "Experiments count incorrect")

	s.Require().Nil(db.Model(&models.Run{}).Count(&countVal).Error)
	s.Equal(counts.runs, int(countVal), "Runs count incorrect")

	s.Require().Nil(db.Model(&models.Metric{}).Count(&countVal).Error)
	s.Equal(counts.metrics, int(countVal), "Metrics count incorrect")

	s.Require().Nil(db.Model(&models.LatestMetric{}).Count(&countVal).Error)
	s.Equal(counts.latestMetrics, int(countVal), "Latest metrics count incorrect")

	s.Require().Nil(db.Model(&models.Tag{}).Count(&countVal).Error)
	s.Equal(counts.tags, int(countVal), "Run tags count incorrect")

	s.Require().Nil(db.Model(&models.Param{}).Count(&countVal).Error)
	s.Equal(counts.params, int(countVal), "Run params count incorrect")

	s.Require().Nil(db.Model(&models.Run{}).Distinct("experiment_id").Count(&countVal).Error)
	s.Equal(counts.distinctRunExperimentIDs, int(countVal), "Runs experiment association incorrect")

	s.Require().Nil(db.Model(&database.App{}).Count(&countVal).Error)
	s.Equal(counts.apps, int(countVal), "Apps count incorrect")

	s.Require().Nil(db.Model(&database.Dashboard{}).Count(&countVal).Error)
	s.Equal(counts.dashboards, int(countVal), "Dashboard count incorrect")
}

// validateTable will scan source and dest table and confirm they are identical
func (s *ImportTestSuite) validateTable(source, dest *gorm.DB, table string) {
	sourceRows, err := source.Table(table).Rows()
	s.Require().Nil(err)
	s.Require().Nil(sourceRows.Err())
	destRows, err := dest.Table(table).Rows()
	s.Require().Nil(err)
	s.Require().Nil(destRows.Err())
	//nolint:errcheck
	defer sourceRows.Close()
	//nolint:errcheck
	defer destRows.Close()

	for sourceRows.Next() {
		// dest should have the same number of rows as source
		s.Require().True(destRows.Next())

		var sourceRow, destRow map[string]any
		s.Require().Nil(source.ScanRows(sourceRows, &sourceRow))
		s.Require().Nil(dest.ScanRows(destRows, &destRow))

		// translate some types to make comparison easier
		for _, row := range []map[string]any{sourceRow, destRow} {
			for k, v := range row {
				switch k {
				case "is_nan", "is_archived":
					if v, ok := v.(float64); ok {
						row[k] = v != 0
					}
				case "default_experiment_id", "experiment_id":
					if v, ok := v.(int64); ok {
						row[k] = int32(v)
					}
				case "id", "app_id":
					switch v := v.(type) {
					case *interface{}:
						switch s := (*v).(type) {
						case string:
							row[k] = s
						}
					}
				}
			}
		}

		// TODO:DSuhinin delete this fields right now, because they
		// cause comparison error when we compare `namespace` entities. Let's find smarter way to do that.
		delete(destRow, "updated_at")
		delete(destRow, "created_at")
		delete(sourceRow, "updated_at")
		delete(sourceRow, "created_at")

		s.Equal(sourceRow, destRow)
	}
	// dest should have the same number of rows as source
	s.Require().False(destRows.Next())
}
