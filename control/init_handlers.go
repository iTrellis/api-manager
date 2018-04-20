package control

import (
	"net/http"

	"github.com/go-trellis/api-manager/handlers"
	"github.com/go-trellis/api-manager/models"

	"github.com/gin-gonic/gin"
	"github.com/go-trellis/connector/tgorm"
)

type GinEntry struct {
	Root       string
	APIVersion string
	Handlers   []GinHandler
}

type GinHandler struct {
	Path        string
	HandleFuncs map[string]func(*gin.Context)
}

// GinEntries 处理对象
var GinEntries []*GinEntry

func initHandlers() {

	dbs, err := tgorm.NewDBs(ConfigDatabase)
	if err != nil {
		panic(err)
	}

	params := map[string]interface{}{
		"committer":                   tgorm.NewCommitter(),
		models.ModelProjectName:       models.NewMProject(dbs),
		models.ModelAPIName:           models.NewMAPI(dbs),
		models.ModelAPIParamsName:     models.NewMAPIParams(dbs),
		models.ModelProjectStatusName: models.NewMProjectStatus(dbs),
		models.ModelFieldTypeName:     models.NewMFieldType(dbs),
	}

	handlers.NewHProjects().Inject(params)
	handlers.NewHProject().Inject(params)
	handlers.NewHProjectStatus().Inject(params)
	handlers.NewHAPIs().Inject(params)
	handlers.NewHAPI().Inject(params)
	handlers.NewHAPIParams().Inject(params)
	handlers.NewHFieldType().Inject(params)

	// KEY is "ROOT/APIVersion"
	// TODO 配置获取handlers
	GinEntries = []*GinEntry{&GinEntry{
		Root: "/h", APIVersion: "/v1",
		Handlers: []GinHandler{
			GinHandler{
				Path: "/projects",
				HandleFuncs: map[string]func(*gin.Context){
					http.MethodGet: handlers.NewHProjects().Get,
				},
			},
			GinHandler{
				Path: "/project",
				HandleFuncs: map[string]func(*gin.Context){
					http.MethodGet:  handlers.NewHProject().Get,
					http.MethodPost: handlers.NewHProject().Post,
					http.MethodPut:  handlers.NewHProject().Put,
				},
			},
			GinHandler{
				Path: "/project/status",
				HandleFuncs: map[string]func(*gin.Context){
					http.MethodGet: handlers.NewHProjectStatus().Get,
				},
			},
			GinHandler{
				Path: "/apis",
				HandleFuncs: map[string]func(*gin.Context){
					http.MethodGet: handlers.NewHAPIs().Get,
				},
			},
			GinHandler{
				Path: "/api",
				HandleFuncs: map[string]func(*gin.Context){
					http.MethodGet:    handlers.NewHAPI().Get,
					http.MethodPost:   handlers.NewHAPI().Post,
					http.MethodPut:    handlers.NewHAPI().Put,
					http.MethodDelete: handlers.NewHAPI().Delete,
				},
			},
			GinHandler{
				Path: "/api/params",
				HandleFuncs: map[string]func(*gin.Context){
					http.MethodPost:   handlers.NewHAPIParams().Post,
					http.MethodDelete: handlers.NewHAPIParams().Delete,
				},
			},
			GinHandler{
				Path: "/feild_type",
				HandleFuncs: map[string]func(*gin.Context){
					http.MethodGet: handlers.NewHFieldType().Get,
				},
			},
		},
	},
	}
}
