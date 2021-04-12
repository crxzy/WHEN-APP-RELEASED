package api

import (
	"net/http"
	"strings"

	"github.com/emicklei/go-restful/v3"
	"gorm.io/gorm"
)

//Channel msg
type Channel struct {
	gorm.Model
	Db        *gorm.DB `gorm:"-" json:"-"`
	Desc      string
	Name      *string `gorm:"not null"`
	URL       *string `gorm:"not null"`
	Author    string
	ProjectID uint `gorm:"not null,default:0"`
	Status    int  `json:"-" gorm:"default:0"`
}

//DeleteChannel del channel
func (p *Channel) DeleteChannel(request *restful.Request, response *restful.Response) {
	rtn := Wrapper{Status: 200, Msg: "ok"}

	id := request.PathParameter("channel-id")

	p.Db.Delete(&Channel{}, id)
	response.WriteEntity(rtn)
}

//UpdateChannel update Channel
func (p *Channel) UpdateChannel(request *restful.Request, response *restful.Response) {
	rtn := Wrapper{Status: 200, Msg: "ok"}
	channel := &Channel{}
	err := request.ReadEntity(&channel)
	if err != nil {
		rtn.Status = http.StatusInternalServerError
		rtn.Msg = err.Error()
	} else {
		p.Db.Model(&channel).Omit("ID", "CreateAt", "UpdateAt", "DeletedAt", "Status").Updates(channel)
	}
	response.WriteEntity(rtn)
}

//AddChannel save Channel
func (p *Channel) AddChannel(request *restful.Request, response *restful.Response) {
	rtn := Wrapper{Status: 200, Msg: "ok"}
	channel := &Channel{}
	err := request.ReadEntity(&channel)
	if err != nil {
		rtn.Status = http.StatusInternalServerError
		rtn.Msg = err.Error()
	} else {
		if strings.Trim(*channel.Name, " ") == "" {
			rtn.Status = http.StatusNotAcceptable
			rtn.Msg = "Parameter error"
		} else {
			result := p.Db.Omit("ID", "CreateAt", "UpdateAt", "DeletedAt", "Status").Create(channel)
			if result.Error != nil {
				rtn.Status = http.StatusInternalServerError
				rtn.Msg = result.Error.Error()
			}
		}
	}
	response.WriteEntity(rtn)
}

//AllChannel channel list
func (p *Channel) AllChannel(request *restful.Request, response *restful.Response) {
	rtn := Wrapper{Status: 200, Msg: "ok"}

	var channels []Channel
	result := p.Db.Find(&channels)
	if result.Error != nil {
		rtn.Msg = result.Error.Error()
		rtn.Status = http.StatusInternalServerError
	} else {
		rtn.Data = channels
	}

	response.WriteEntity(rtn)
}

//FindChannel find Channel
func (p *Channel) FindChannel(request *restful.Request, response *restful.Response) {
	rtn := Wrapper{Status: 200, Msg: "ok"}

	id := request.PathParameter("channel-id")
	channel := &Channel{}
	result := p.Db.Where("id = ?", id).Find(&channel)

	if result.Error != nil {
		rtn.Msg = result.Error.Error()
		rtn.Status = http.StatusInternalServerError
	} else {
		if result.RowsAffected > 0 {
			rtn.Data = channel
		}
	}

	response.WriteEntity(rtn)
}
