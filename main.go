package main

import (
	"channel/api"
	"channel/task"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	restful "github.com/emicklei/go-restful/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

//RegisterTo register server
func RegisterTo(container *restful.Container, db *gorm.DB) {

	// rest /project
	serverProject := new(restful.WebService)
	p := &api.Project{Db: db}
	db.AutoMigrate(p)
	serverProject.
		Path("/v1/project").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	serverProject.Route(serverProject.GET("").To(p.AllProject))
	serverProject.Route(serverProject.GET("/{project-id}").To(p.FindProject))
	serverProject.Route(serverProject.POST("").To(p.AddProject))
	serverProject.Route(serverProject.PUT("").To(p.UpdateProject))
	serverProject.Route(serverProject.DELETE("/{project-id}").To(p.DeleteProject))
	container.Add(serverProject)

	serverChannel := new(restful.WebService)
	c := &api.Channel{Db: db}
	db.AutoMigrate(c)
	serverChannel.
		Path("/v1/channel").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	serverChannel.Route(serverChannel.GET("").To(c.AllChannel))
	serverChannel.Route(serverChannel.GET("/{channel-id}").To(c.FindChannel))
	serverChannel.Route(serverChannel.POST("").To(c.AddChannel))
	serverChannel.Route(serverChannel.PUT("").To(c.UpdateChannel))
	serverChannel.Route(serverChannel.DELETE("/{channel-id}").To(c.DeleteChannel))
	container.Add(serverChannel)

	//db.SetupJoinTable(p, "ProjectChannel", pc)
	projectChannel := new(restful.WebService)
	pc := &api.ProjectChannel{Db: db}
	db.AutoMigrate(pc)
	projectChannel.
		Path("/v1").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	projectChannel.Route(projectChannel.POST("/connect").To(pc.Connect))
	projectChannel.Route(projectChannel.DELETE("/connect/{connect-id}").To(pc.RemoveConnect))

	container.Add(projectChannel)
}

// NotifyTask updatedata when somgthing call
func NotifyTask(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	chain.ProcessFilter(req, resp)
	if req.Request.Method != http.MethodGet {
		notify <- 0
	}
}

var notify = make(chan int)

func main() {
	wsContainer := restful.NewContainer()

	errorlog, err := os.OpenFile("gorm.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
		os.Exit(1)
	}
	newLogger := logger.New(
		log.New(errorlog, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			LogLevel: logger.Silent, // Log level
			Colorful: true,          // Disable color
		},
	)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/channel?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_ADDR"))
	fmt.Println(dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		fmt.Println(err.Error())
		panic("connect to mysql failed")
	}

	RegisterTo(wsContainer, db)
	wsContainer.Filter(NotifyTask)

	cors := restful.CrossOriginResourceSharing{
		ExposeHeaders:  []string{"X-My-Header"},
		AllowedHeaders: []string{"Content-Type", "Accept"},
		AllowedMethods: []string{"GET", "POST", "DELETE", "PUT"},
		CookiesAllowed: false,
		Container:      wsContainer}
	wsContainer.Filter(cors.Filter)

	var jobs = &task.Jobs{ScheduleTable: make(map[int]time.Time), Notify: notify, Db: db}

	go task.Loop(jobs)
	log.Fatal(http.ListenAndServe(":62211", wsContainer))
}
