package main

import (
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/pmdcosta/mmm/crawler"
	"github.com/pmdcosta/mmm/database"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

var wg sync.WaitGroup

func main() {
	// create logger.
	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	// gracefully shutdown.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// parse flags.
	var (
		cronTimer = flag.Int("time", 60, "Seconds between updater run")
		dbHost    = flag.String("dbhost", "database.db", "Database file location")
		dlPath    = flag.String("dlPath", filepath.Join("C:", "Users", "pmdco", "Downloads"), "Download target location")
	)
	flag.Parse()

	// open database.
	db := database.NewClient(logger, *dbHost)
	err := db.Open()
	if err != nil {
		panic(err.Error())
	}

	// open crawler.
	c := crawler.NewClient(logger)

	// run cron job
	quit := make(chan bool)
	wg.Add(1)
	go func() {
		ticker := time.NewTicker(time.Duration(*cronTimer) * time.Second)
		defer wg.Done()
		for {
			select {
			case <-ticker.C:
				cron(logger, *db, *c, *dlPath)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	// wait for user signal.
	<-sigs
	quit <- true
	logger.Log("msg", "shutting down")
	wg.Wait()
}

func cron(logger log.Logger, db database.Client, c crawler.Client, dlPath string) error {
	// list active Seasons from DB.
	seasons, err := db.SeasonService().ListCompleteSeasons(false)
	if err != nil {
		return err
	}

	// get available episodes for each season.
	for _, s := range seasons {
		ts, err := c.Search(s)
		if err != nil {
			return err
		}
		fmt.Printf("%+v", ts)
		// TODO - check if it has already been downloaded.
		// TODO - download file.
		// TODO - set torrent as downloaded.
	}
	return nil
}
