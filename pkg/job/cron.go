package job

import (
    "log"
)

// cron.v3 interface
func (job *Job) Run() {

	log.Println("launching", job.Label, job.cmdLine)
	if err := job.Launch(); err != nil {
		log.Println("failed:", err)
		return
	}
	log.Println("done: ",
		len(job.Report.Hubs), "hops in", job.Duration, ".")
}
