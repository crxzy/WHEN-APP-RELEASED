package task

import (
	"channel/api"
	"channel/utils"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gorhill/cronexpr"
	"gorm.io/gorm"
)

//Jobs schedule job
type Jobs struct {
	Db            *gorm.DB
	projects      []api.Project
	ScheduleTable map[int]time.Time
	Notify        chan int
}

var errorlog *os.File
var logger *log.Logger

func init() {
	errorlog, err := os.OpenFile("schedule.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
		os.Exit(1)
	}

	logger = log.New(errorlog, "scheduleLog: ", log.Lshortfile|log.LstdFlags)
}

func (p *Jobs) oldUpdateData() {
	result := p.Db.Preload("ProjectChannels.Channel").Find(&p.projects)
	if result.Error != nil {
		return
	}
	now := time.Now()
	for i, pro := range p.projects {
		express := cronexpr.MustParse(*pro.Cron)
		p.ScheduleTable[i] = express.Next(now)
		logger.Println(p.ScheduleTable[i])
	}
}

func (p *Jobs) updateData() {
	oldProject := p.projects
	newSche := make(map[int]time.Time)
	result := p.Db.Preload("ProjectChannels.Channel").Find(&p.projects)
	if result.Error != nil {
		return
	}
	now := time.Now()
	for i, pro := range p.projects {
		for oi, op := range oldProject {
			if pro.ID == op.ID {
				if pro.Cron == op.Cron {
					newSche[i] = p.ScheduleTable[oi]
				}
				break
			}
		}
		if _, ok := newSche[i]; !ok {
			express := cronexpr.MustParse(*pro.Cron)
			newSche[i] = express.Next(now)
		}
		logger.Println("schedule: ", newSche[i], *p.projects[i].Name)
	}

	p.ScheduleTable = newSche
}

func (p *Jobs) updateNextTime(i int) {
	express := cronexpr.MustParse(*p.projects[i].Cron)
	p.ScheduleTable[i] = express.Next(time.Now())
}

func (p *Jobs) do(project *api.Project) {
	logger.Println("do job:", *project.Name)

	msg := ""
	for idx := range project.ProjectChannels {
		pc := &project.ProjectChannels[idx]
		req := CommonRequest{Name: *project.Name,
			BundleID:    project.BundleID,
			PackageName: project.PackageName,
			Extra:       project.Extra,
		}

		logger.Println("do check: ", pc.Channel.Name)

		rand.Seed(time.Now().UnixNano())
		url := fmt.Sprintf("%s?rnd=%d", *pc.Channel.URL, rand.Int())
		resp, err := Check(url, req)
		if err != nil {
			logger.Println(err.Error())
			return
		}

		if resp.Version != pc.CurrentVersion {
			line := fmt.Sprintf("%s-%s: %s release at %s\n", *project.Name, *pc.Channel.Name, resp.Version, resp.ReleaseTime)
			msg += line
			logger.Println("new Version!:", msg)
			pc.CurrentVersion = resp.Version
			pc.ReleaseTime = &resp.ReleaseTime
			p.Db.Model(&pc).Updates(pc)
		} else {
			logger.Println("old version")
		}
	}
	if msg != "" {
		// send popo msg
		n := utils.GetNotify()
		n.SendPopo(msg, []string{})
	}
}

//Loop task schedule
func Loop(jobs *Jobs) {

	jobs.updateData()

	for {
		now := time.Now()

		for i, t := range jobs.ScheduleTable {
			if t.Before(now) || t.Equal(now) {
				jobs.updateNextTime(i)
				go jobs.do(&jobs.projects[i])
			}
		}

		select {
		case <-time.NewTimer(500 * time.Millisecond).C:
		case <-jobs.Notify:
			logger.Println("updateData")
			jobs.updateData()
		}
	}
}
