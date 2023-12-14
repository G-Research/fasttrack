package database

import (
	"database/sql"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gorm.io/datatypes"

	"github.com/G-Research/fasttrackml/pkg/database"
	"github.com/G-Research/fasttrackml/tests/integration/golang/helpers"
)

type Metric struct {
	Key           string  `gorm:"type:varchar(250);not null;primaryKey"`
	Value         float64 `gorm:"type:double precision;not null;primaryKey"`
	ContextNullID *uint
	ContextID     uint `gorm:"not null"`
}

type Context struct {
	ID   uint           `gorm:"primaryKey;autoIncrement"`
	Json datatypes.JSON `gorm:"not null;unique;index"`
}

type JsonTestSuite struct {
	suite.Suite
	db *sql.DB
}

func (s *JsonTestSuite) SetupSuite() {
	// setup db
	dsn, err := helpers.GenerateDatabaseURI(s.T(), helpers.GetDatabaseBackend())
	s.Require().Nil(err)

	dbProvider, err := database.NewDBProvider(
		dsn,
		1*time.Second,
		20,
	)
	s.Require().Nil(err)

	// use simplified schema
	s.Require().Nil(dbProvider.GormDB().AutoMigrate(&Context{}))
	s.Require().Nil(dbProvider.GormDB().AutoMigrate(&Metric{}))

	// Begin a transaction
	s.db, err = dbProvider.GormDB().DB()
	s.Require().Nil(err)
	tx, err := s.db.Begin()
	s.Require().Nil(err)

	// Prepare a statement for inserting data
	contextStmt, err := tx.Prepare("INSERT INTO contexts(id, json) VALUES($1, $2)")
	s.Require().Nil(err)
	//nolint:errcheck
	defer contextStmt.Close()

	// Prepare a statement for inserting data into the 'metrics' table
	stmtMetrics, err := tx.Prepare("INSERT INTO metrics(key, value, context_id, context_null_id) VALUES($1, $2, $3, $4)")
	s.Require().Nil(err)
	//nolint:errcheck
	defer stmtMetrics.Close()

	// Create a default/empty json context
	_, err = contextStmt.Exec(1, "{}") // Insert the JSON document
	s.Require().Nil(err)
	defaultContextId := int64(1)
	s.Require().Nil(err)

	// Insert a large number of rows
	for i := 2; i < 1000002; i++ {
		// Create a JSON document with small variations
		jsonDoc := fmt.Sprintf(`{"key": "key%d", "value": "value%d"}`, i, i)
		id := int64(i)
		_, err := contextStmt.Exec(i, jsonDoc) // Insert the JSON document
		s.Require().Nil(err)

		// Randomly decide whether to insert null/defuult or the current context id
		contextNullId := sql.NullInt64{Int64: id, Valid: true}
		//nolint:gosec
		if rand.Intn(2) == 0 { // 50% chance of being true
			// when true, put empty doc reference and null
			contextNullId = sql.NullInt64{Int64: 0, Valid: false}
			id = defaultContextId
		}

		// Insert into the 'metrics' table
		key := fmt.Sprintf("key%d", i)
		//nolint:gosec
		value := rand.Float64()
		_, err = stmtMetrics.Exec(key, value, id, contextNullId) // Replace 'i' with the actual value you want to insert
		s.Require().Nil(err)
	}

	// Commit the transaction
	err = tx.Commit()
	s.Require().Nil(err)
}

func (s *JsonTestSuite) TearDownSuite() {
	// Close the database connection
	s.Require().Nil(s.db.Close())
}

func TestJsonTestSuite(t *testing.T) {
	suite.Run(t, new(JsonTestSuite))
}

func (s *JsonTestSuite) TestJson() {
	tests := []struct {
		name       string
		joinColumn string
		key        string
		value      string
	}{
		{
			name:       "TestNullable",
			joinColumn: "context_null_id",
		},
		{
			name:       "TestNotNullable",
			joinColumn: "context_id",
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Begin a transaction
			tx, err := s.db.Begin()
			s.Require().Nil(err)
			//nolint:errcheck
			defer tx.Commit()

			pathOperator := `->>`
			keyPattern := `key%d`
			if helpers.GetDatabaseBackend() == `postgres` {
				pathOperator = `#>>`
				keyPattern = `{key%d}`
			}

			// Prepare a statement for selecting data using the join column
			// and a json path expression
			sql := `SELECT * FROM metrics LEFT JOIN contexts ON metrics.` +
				tt.joinColumn +
				` = contexts.id WHERE contexts.json` +
				pathOperator +
				`$1 = $2`

			//nolint:gosec
			contextStmt, err := tx.Prepare(sql)
			s.Require().Nil(err)

			for i := 0; i < 1000; i++ {
				key := fmt.Sprintf(keyPattern, i)
				value := fmt.Sprintf("value%d", i)
				_, err = contextStmt.Exec(key, value)
				s.Require().Nil(err)
			}

			//nolint:errcheck
			defer contextStmt.Close()
		})
	}
}
