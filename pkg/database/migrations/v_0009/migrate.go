package v_0009

import (
	"errors"
	"fmt"

	"github.com/rotisserie/eris"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/G-Research/fasttrackml/pkg/database/migrations"
)

const Version = "2c2299e4e061"

func Migrate(db *gorm.DB) error {
	return migrations.RunWithoutForeignKeyIfNeeded(db, func() error {
		return db.Transaction(func(tx *gorm.DB) error {
			// Rename the existing metricx tables and drop indexes
			tables := []string{"metrics", "latest_metrics"}
			for _, table := range tables {
				index := fmt.Sprintf("idx_%s_run_id", table)
				if err := dropIndex(tx, table, index); err != nil {
					return eris.Wrapf(err, "error dropping %s", index)
				}
				if err := tx.Migrator().RenameTable(table, backupName(table)); err != nil {
					return eris.Wrapf(err, "error renaming %s", table)
				}
			}
			index := "idx_metrics_iter"
			if err := dropIndex(tx, backupName("metrics"), index); err != nil {
				return eris.Wrapf(err, "error dropping %s", index)
			}

			// Auto-migrate the replacements and new tables
			if err := tx.Migrator().AutoMigrate(&Context{}, &Metric{}, &LatestMetric{}); err != nil {
				return eris.Wrap(err, "error automigrating new tables")
			}

			// Create the default metric context
			if err := createDefaultMetricContext(tx); err != nil {
				return eris.Wrap(err, "error creating default metric context")
			}

			// Copy the data from the old tables to the new ones with default metric context
			for _, table := range tables {
				if err := tx.Exec(fmt.Sprintf("INSERT INTO %s SELECT *, %d FROM %s",
					table,
					DefaultContext.ID,
					backupName(table))).Error; err != nil {
					return eris.Wrapf(err, "error copying data for %s", table)
				}

				// Drop the backup tables
				if err := tx.Exec("DROP TABLE " + backupName(table)).Error; err != nil {
					return eris.Wrapf(err, "error dropping %s", backupName(table))
				}
			}

			// Update the schema version
			return tx.Model(&SchemaVersion{}).
				Where("1 = 1").
				Update("Version", Version).
				Error
		})
	})
}

func backupName(name string) string {
	return fmt.Sprintf("%s_old", name)
}

func dropIndex(tx *gorm.DB, table, index string) error {
	if tx.Migrator().HasIndex(table, index) {
		if err := tx.Migrator().DropIndex(table, index); err != nil {
			return eris.Wrapf(err, "error dropping %s", index)
		}
	}
	return nil
}

// createDefaultMetricContext creates the default metric context if it doesn't exist.
func createDefaultMetricContext(db *gorm.DB) error {
	if err := db.First(&DefaultContext).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Info("Creating default context")
			if err := db.Create(&DefaultContext).Error; err != nil {
				return fmt.Errorf("error creating default context: %s", err)
			}
		} else {
			return fmt.Errorf("unable to find default context: %s", err)
		}
	}
	log.Debugf("default metric context: %v", DefaultContext)
	return nil
}
