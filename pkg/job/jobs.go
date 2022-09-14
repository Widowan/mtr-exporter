package job

type Jobs []*Job

func (jobs Jobs) Count() int { return len(jobs) }

func (jobs Jobs) CollectedReports() int {
	n := 0
	for _, job := range jobs {
		if job.Report != nil {
			n++
		}
	}
	return n
}
