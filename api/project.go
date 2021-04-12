package api

import (
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"github.com/gorhill/cronexpr"
	"gorm.io/gorm"
)

// Project project model
type Project struct {
	gorm.Model
	Db              *gorm.DB `gorm:"-" json:"-"`
	Desc            string
	Name            *string `gorm:"not null"`
	BundleID        string
	PackageName     string
	Extra           string
	Cron            *string          `gorm:"not null"`
	Status          int              `json:"-" gorm:"default:0"`
	ProjectChannels []ProjectChannel `gorm:"foreignKey:ProjectID" json:"channels,omitempty"`
}

//DeleteProject del project
func (p *Project) DeleteProject(request *restful.Request, response *restful.Response) {
	rtn := Wrapper{Status: 200, Msg: "ok"}

	id := request.PathParameter("project-id")

	p.Db.Delete(&Project{}, id)
	response.WriteEntity(rtn)
}

//UpdateProject update project
func (p *Project) UpdateProject(request *restful.Request, response *restful.Response) {
	rtn := Wrapper{Status: 200, Msg: "ok"}
	project := &Project{}
	err := request.ReadEntity(&project)
	if err != nil {
		rtn.Status = http.StatusInternalServerError
		rtn.Msg = err.Error()
	} else {
		_, err := cronexpr.Parse(*project.Cron)
		if err != nil {
			rtn.Status = http.StatusBadRequest
			rtn.Msg = "cronexpr error"
		} else {
			p.Db.Model(&project).Omit("ID", "CreateAt", "UpdateAt", "DeletedAt", "Status").Updates(project)
		}
	}
	response.WriteEntity(rtn)
}

//AddProject save project
func (p *Project) AddProject(request *restful.Request, response *restful.Response) {
	rtn := Wrapper{Status: 200, Msg: "ok"}
	project := &Project{}
	err := request.ReadEntity(&project)
	if err != nil {
		rtn.Status = http.StatusInternalServerError
		rtn.Msg = err.Error()
	} else {
		_, err := cronexpr.Parse(*project.Cron)
		if err != nil {
			rtn.Status = http.StatusBadRequest
			rtn.Msg = "cronexpr error"
		} else {
			result := p.Db.Omit("ID", "CreateAt", "UpdateAt", "DeletedAt", "Status").Create(project)
			if result.Error != nil {
				rtn.Status = http.StatusInternalServerError
				rtn.Msg = result.Error.Error()
			}
		}
	}
	response.WriteEntity(rtn)
}

//AllProject project list
func (p *Project) AllProject(request *restful.Request, response *restful.Response) {
	rtn := Wrapper{Status: 200, Msg: "ok"}

	var projects []Project
	result := p.Db.Find(&projects)
	if result.Error != nil {
		rtn.Msg = result.Error.Error()
		rtn.Status = http.StatusInternalServerError
	} else {
		rtn.Data = projects
	}

	response.WriteEntity(rtn)
}

//FindProject find project
func (p *Project) FindProject(request *restful.Request, response *restful.Response) {
	rtn := Wrapper{Status: 200, Msg: "ok"}

	id := request.PathParameter("project-id")
	project := &Project{}
	result := p.Db.Where("id = ?", id).Preload("ProjectChannels.Channel").Find(&project)

	if result.Error != nil {
		rtn.Msg = result.Error.Error()
		rtn.Status = http.StatusInternalServerError
	} else {
		if result.RowsAffected > 0 {
			rtn.Data = project
		}
	}

	response.WriteEntity(rtn)
}
