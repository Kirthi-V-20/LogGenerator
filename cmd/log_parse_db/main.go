package main

import (
	databasemodel "LogGenerator/pkg/db_model"
	dbmodel "LogGenerator/pkg/db_model"
	"LogGenerator/pkg/parser"
	ginhandler "LogGenerator/pkg/web"
	"fmt"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
)

const dbUrl = "postgresql:///log_analyzer?host=/var/run/postgresql/"

func handleCommand(args []string) error {
	db, err := dbmodel.CreateDB(dbUrl)
	if err != nil {
		return err
	}
	switch args[0] {
	case "init":
		err := dbmodel.InitDb(db)
		if err != nil {
			return err
		}
	case "add":
		dirpath := args[1]
		if dirpath == "" {
			slog.Error("Specify directory!")
		}
		entries, err := parser.ParseLogFiles(dirpath)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			databasemodel.AddEntry(db, entry)
		}
		return nil
	case "query":
		query := args[1:]
		fmt.Println(query)

		entries, err := databasemodel.Query(db, query)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			fmt.Println(entry)
		}
		slog.Info("Filtering successful!", "no. of entries:", len(entries))
		return nil
	case "web":
		r := gin.Default()
		r.LoadHTMLGlob("pkg/web/templates/*")
		ginhandler.DBRef = db
		ginhandler.SetupRoutes(r)
		r.Run(":8081")
	default:
		return fmt.Errorf("unknown command: %s (expected: init | add | query)", args[0])

	}
	return nil

}

func main() {
	err := handleCommand(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in invocation %v", err)
		os.Exit(-1)
	}

}
