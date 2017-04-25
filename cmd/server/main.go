package main

import (
	"flag"
	"github.com/go-kit/kit/log"
	"github.com/pmdcosta/mmm/crawler"
	"github.com/pmdcosta/mmm/database"
	"github.com/pmdcosta/mmm/download"
	"os"
	"os/signal"
	"syscall"
	"time"
)

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
		dlRate = flag.Int("rate", 10, "Minutes between updater run")
		dbHost = flag.String("db", "D:\\temp\\database.db", "Database file location")
		dlPath = flag.String("path", "C:\\Users\\pmdco\\Downloads", "Download target location")
	)
	flag.Parse()

	// create database client.
	db := database.NewClient(logger, *dbHost)
	if err := db.Open(); err != nil {
		panic(err.Error())
	}

	// create crawler client.
	cr := crawler.NewClient(logger)

	// create downloader client.
	dl := download.NewClient(logger, db.SeasonService(), cr.CrawlerService(), *dlPath, time.Duration(*dlRate))
	if err := dl.Open(); err != nil {
		panic(err.Error())
	}

	// wait for user signal.
	<-sigs
	dl.Close()
	db.Close()
	logger.Log("msg", "shutting down")
}
