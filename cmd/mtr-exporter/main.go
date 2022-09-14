package main

// *mtr-exporter* periodically executes *mtr* to a given host and provides the
// measured results as prometheus metrics.

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mgumz/mtr-exporter/pkg/job"

	"github.com/robfig/cron/v3"
)

func main() {

	log.SetFlags(0)

	mtrBin := flag.String("mtr", "mtr", "path to `mtr` binary")
	jobLabel := flag.String("label", "mtr-exporter-cli", "job label")
	bind := flag.String("bind", ":8080", "bind address")
	jobFile := flag.String("jobs", "", "file containing job definitions")
	schedule := flag.String("schedule", "@every 60s", "schedule at which often `mtr` is launched")
	doPrintVersion := flag.Bool("version", false, "show version")
	doPrintUsage := flag.Bool("h", false, "show help")
	doTimeStampLogs := flag.Bool("tslogs", false, "use timestamps in logs")

	flag.Usage = usage
	flag.Parse()

	if *doPrintVersion == true {
		printVersion()
		return
	}
	if *doPrintUsage == true {
		flag.Usage()
		return
	}
	if *doTimeStampLogs == true {
		log.SetFlags(log.LstdFlags | log.LUTC)
	}

	jobs, err := job.ParseJobFile(*jobFile, *mtrBin)
	if err != nil {
		log.Printf("error: parsing jobs file %q: %s", *jobFile, err)
		os.Exit(1)
	}

	if len(flag.Args()) > 0 {
		job := job.NewJob(*mtrBin, flag.Args(), *schedule)
		job.Label = *jobLabel
		jobs = append(jobs, job)
	}

	// TODO: if watching *jobsFile is implemented, having 0 jobs is ok
	if len(jobs) == 0 {
		log.Println("error: no mtr jobs defined - provide at least one via -jobfile or via arguments")
		os.Exit(1)
		return
	}

	c := cron.New()
	for _, job := range jobs {
		c.AddJob(job.Schedule, job)
	}
	c.Start()

	http.Handle("/metrics", jobs)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})

	log.Println("serving /metrics at", *bind, "...")
	log.Fatal(http.ListenAndServe(*bind, nil))
}
