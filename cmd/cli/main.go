package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"github.com/go-kit/kit/log"
	"github.com/pmdcosta/mmm"
	"github.com/pmdcosta/mmm/database"
	"os"
)

func main() {
	// create logger.
	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	// parse flags.
	var (
		inputPath = flag.String("path", "D:\\temp\\input.txt", "Seconds between updater run")
		dbHost    = flag.String("dbhost", "D:\\temp\\database.db", "Database file location")
	)
	flag.Parse()

	// open database.
	db := database.NewClient(logger, *dbHost)
	if err := db.Open(); err != nil {
		panic(err.Error())
	}

	parseFile(logger, *inputPath, db)
}

func parseFile(logger log.Logger, input string, db *database.Client) {
	// open file.
	file, err := os.Open(input)
	if err != nil {
		logger.Log("err", err.Error())
		return
	}
	defer file.Close()

	// scan file.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		saveEntry(scanner.Text(), db)
	}

	if err := scanner.Err(); err != nil {
		logger.Log("err", err.Error())
		return
	}
}

func saveEntry(line string, db *database.Client) {
	var entry mmm.Season
	json.Unmarshal([]byte(line), &entry)
	entry.Type = mmm.TypeAnime
	entry.State = mmm.StateRunning
	entry.Index = 1

	if _, err := db.SeasonService().CreateSeason(&entry); err != nil && err != database.ErrDatabaseExists {
		panic(err.Error())
	}
}
