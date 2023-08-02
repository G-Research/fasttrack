package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/G-Research/fasttrackml/pkg/database"
)

var ImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Copies an input database to an output database",
	Long: `The import command will transfer the contents of the input
         database to the output database. Please make sure that the
         FasttrackML server is not currently connected to the input
         database.`,
	RunE: importCmd,
}

func importCmd(cmd *cobra.Command, args []string) error {
	// 1. init database connections.
	inputDB, outputDB, err := initDBs()
	if err != nil {
		return err
	}
	defer inputDB.Close()
	defer outputDB.Close()

	dryRun := viper.GetBool("dry-run")
	if err := database.Import(inputDB, outputDB, dryRun); err != nil {
		return err
	}

	return nil
}

// initDBs inits the input and output DB connections.
func initDBs() (input, output *database.DbInstance, err error) {
	databaseSlowThreshold := time.Second * 1
	databasePoolMax := 20
	databaseReset := false
	databaseMigrate := true
	artifactRoot := "s3://fasttrackml"
	input, err = database.ConnectDB(
		viper.GetString("input-database-uri"),
		databaseSlowThreshold,
		databasePoolMax,
		databaseReset,
		databaseMigrate,
		artifactRoot,
	)
	if err != nil {
		return input, output, fmt.Errorf("error connecting to input DB: %w", err)
	}

	output, err = database.ConnectDB(
		viper.GetString("output-database-uri"),
		databaseSlowThreshold,
		databasePoolMax,
		databaseReset,
		databaseMigrate,
		artifactRoot,
	)
	if err != nil {
		return input, output, fmt.Errorf("error connecting to output DB: %w", err)
	}
	return
}

func init() {
	RootCmd.AddCommand(ImportCmd)

	ImportCmd.Flags().StringP("input-database-uri", "i", "", "Input Database URI (eg., sqlite://fasttrackml.db)")
	ImportCmd.Flags().StringP("output-database-uri", "o", "", "Output Database URI (eg., postgres://user:psw@postgres:5432)")
	ImportCmd.Flags().BoolP("dry-run", "n", false, "Perform a dry run (will not write anything)")
	ImportCmd.MarkFlagRequired("input-database-uri")
	ImportCmd.MarkFlagRequired("output-database-uri")
}
