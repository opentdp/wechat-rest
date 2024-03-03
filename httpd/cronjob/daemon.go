package cronjob

import (
	"github.com/opentdp/go-helper/command"
	"github.com/opentdp/go-helper/logman"
	"github.com/robfig/cron/v3"

	"github.com/opentdp/wechat-rest/dbase/cronjob"
	"github.com/opentdp/wechat-rest/dbase/tables"
)

var crontab *cron.Cron
var logger *logman.Logger

func Daemon() {

	logger = logman.Named("cronjob")
	logger.Info("cron:daemon start")

	crontab = cron.New(cron.WithSeconds())

	jobs, err := cronjob.FetchAll(&cronjob.FetchAllParam{})
	if err != nil || len(jobs) == 0 {
		return
	}

	for _, job := range jobs {
		AttachJob(job)
	}

	crontab.Start()

}

func AttachJob(job *tables.Cronjob) error {

	sepc := job.Second + " " + job.Minute + " " + job.Hour + " " + job.DayOfMonth + " " + job.Month + " " + job.DayOfWeek

	task := func() {
		logger.Info("cron:run "+job.Name, "entryId", job.EntryId)
		output, err := command.Exec(&command.ExecPayload{
			Name:          "cron: " + job.Name,
			CommandType:   job.Type,
			WorkDirectory: job.Directory,
			Content:       job.Content,
			Timeout:       job.Timeout,
		})
		logger.Warn("cron:run "+job.Name, "output", output, "error", err)
		// 将输出发送到指定的通道
		if err == nil && output != "" && job.Deliver != "-" {
			MsgDeliver(job.Deliver, output)
		}
	}

	entryId, err := crontab.AddFunc(sepc, task)
	if err != nil {
		return err
	}

	logger.Info("cron:attach "+job.Name, "entryId", entryId)
	err = cronjob.Update(&cronjob.UpdateParam{
		Rd:      job.Rd,
		EntryId: int64(entryId),
	})

	return err

}

func NewById(rd uint) {

	job, err := cronjob.Fetch(&cronjob.FetchParam{Rd: rd})

	if err == nil && job.Rd > 0 {
		AttachJob(job)
	}

}

func UndoById(rd uint) {

	job, err := cronjob.Fetch(&cronjob.FetchParam{Rd: rd})

	if err == nil && job.Rd > 0 {
		logger.Info("cron:remove "+job.Name, "entryId", job.EntryId)
		crontab.Remove(cron.EntryID(job.EntryId))
	}

}

func RedoById(rd uint) {

	job, err := cronjob.Fetch(&cronjob.FetchParam{Rd: rd})

	if err == nil && job.Rd > 0 {
		logger.Info("cron:update "+job.Name, "entryId", job.EntryId)
		crontab.Remove(cron.EntryID(job.EntryId))
		AttachJob(job)
	}

}

// 获取任务执行状态

type JobStatus struct {
	EntryId  int64 `json:"entry_id"`
	NextTime int64 `json:"next_time"`
	PrevTime int64 `json:"prev_time"`
}

func GetEntries() map[uint]JobStatus {

	list := map[uint]JobStatus{}

	jobs, err := cronjob.FetchAll(&cronjob.FetchAllParam{})
	if err != nil || len(jobs) == 0 {
		return list
	}

	for _, job := range jobs {
		entry := crontab.Entry(cron.EntryID(job.EntryId))
		list[job.Rd] = JobStatus{
			EntryId:  int64(entry.ID),
			NextTime: entry.Next.Unix(),
			PrevTime: entry.Prev.Unix(),
		}
	}

	return list

}