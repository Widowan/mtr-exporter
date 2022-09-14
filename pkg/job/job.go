package job

import (
	"bytes"
	"os/exec"
	"strings"
	"sync"
	"time"

    "github.com/mgumz/mtr-exporter/pkg/mtr"
)

type Job struct {
	Report   *mtr.Report
	Launched time.Time
	Duration time.Duration
	Schedule string
	Label    string

	mtrBinary string
	args      []string
	cmdLine   string

	sync.Mutex
}

func NewJob(mtr string, args []string, schedule string) *Job {
	extra := []string{
		"-j", // json output
	}
	args = append(extra, args...)
	cmd := exec.Command(mtr, args...)

	return &Job{
		Schedule:  schedule,
		args:      args,
		mtrBinary: mtr,
		cmdLine:   strings.Join(cmd.Args, " "),
	}
}

func (job *Job) Launch() error {

	// TODO: maybe use CommandContext to have an upper limit in the execution

	cmd := exec.Command(job.mtrBinary, job.args...)

	// launch mtr
	buf := bytes.Buffer{}
	cmd.Stdout = &buf
	launched := time.Now()
	if err := cmd.Run(); err != nil {
		return err
	}
	duration := time.Since(launched)

	// decode the report
	report := &mtr.Report{}
	if err := report.Decode(&buf); err != nil {
		return err
	}

	// copy the report into the job
	job.Lock()
	job.Report = report
	job.Launched = launched
	job.Duration = duration
	job.Unlock()

	// done.
	return nil
}

