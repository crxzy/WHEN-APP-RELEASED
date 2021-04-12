package api

import (
	"fmt"
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"gorm.io/gorm"
)

//ProjectChannel project <=> channel
type ProjectChannel struct {
	gorm.Model
	Db             *gorm.DB `gorm:"-" json:"-"`
	Status         int      `json:"-" gorm:"default:0"`
	CurrentVersion string
	ReleaseTime    *string
	ProjectID      uint
	ChannelID      uint
	Channel        Channel
}

//RemoveConnect remove connect
func (p *ProjectChannel) RemoveConnect(request *restful.Request, response *restful.Response) {
	rtn := Wrapper{Status: 200, Msg: "ok"}
	id := request.PathParameter("connect-id")
	p.Db.Delete(&ProjectChannel{}, id)
	response.WriteEntity(rtn)
}

//Connect project and channel
func (p *ProjectChannel) Connect(request *restful.Request, response *restful.Response) {
	/*
		{
			projectid
			channelid
		}
	*/
	type R struct {
		ProjectID uint
		ChannelID uint
	}

	rtn := Wrapper{Status: 200, Msg: "ok"}
	r := &R{}
	err := request.ReadEntity(&r)
	if err != nil {
		rtn.Status = http.StatusInternalServerError
		rtn.Msg = err.Error()
		response.WriteEntity(rtn)
		return
	}

	project := &Project{}

	result := p.Db.Where("id = ?", r.ProjectID).Preload("ProjectChannels.Channel").Find(&project)
	if result.Error != nil {
		rtn.Msg = result.Error.Error()
		rtn.Status = http.StatusInternalServerError
		response.WriteEntity(rtn)
		return
	}

	if result.RowsAffected == 0 {
		rtn.Msg = "project not exists"
		rtn.Status = http.StatusBadRequest
		response.WriteEntity(rtn)
		return
	}

	fmt.Println("==============================")
	fmt.Println(project.ProjectChannels)

	for _, v := range project.ProjectChannels {
		if v.Channel.ID == r.ChannelID {
			rtn.Msg = "already exists"
			rtn.Status = http.StatusOK
			response.WriteEntity(rtn)
			return
		}
	}

	c := &Channel{}
	result = p.Db.Where("id = ?", r.ChannelID).Find(&c)
	if result.Error != nil {
		rtn.Msg = result.Error.Error()
		rtn.Status = http.StatusInternalServerError
		response.WriteEntity(rtn)
		return
	}

	if result.RowsAffected == 0 {
		rtn.Msg = "channel not exists"
		rtn.Status = http.StatusBadRequest
		response.WriteEntity(rtn)
		return
	}

	pc := &ProjectChannel{}
	pc.ChannelID = r.ChannelID
	pc.ProjectID = r.ProjectID

	//project.ProjectChannels = append(project.ProjectChannels, *pc)

	p.Db.Create(pc)
	if result.Error != nil {
		rtn.Status = http.StatusInternalServerError
		rtn.Msg = result.Error.Error()
	}

	response.WriteEntity(rtn)

}
