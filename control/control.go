package control

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-trellis/config"
)

// configs
const (
	ConfigDatabase = "etc/database.yaml"
	ConfigApp      = "etc/app.yaml"
)

var app = AppConfig{}

// AppConfig 应用配置
type AppConfig struct {
	Gin struct {
		AppSubURL string `json:"app_sub_url" yaml:"app_sub_url"`
		AppName   string `json:"app_name" yaml:"app_name"`
		Address   string `json:"address" yaml:"address"`
		Mode      string `json:"mode" yaml:"mode"`
	} `json:"gin" yaml:"gin"`
}

func customEntry() {

	initHandlers()

	if err := config.NewSuffixReader().Read(ConfigApp, &app); err != nil {
		panic(err)
	}

	gin.SetMode(app.Gin.Mode)

	router := gin.New()

	router.Use(cors.Default())

	router.Static("/web", "./web")

	for _, entry := range GinEntries {
		for _, handler := range entry.Handlers {
			for m, hf := range handler.HandleFuncs {
				router.Handle(m, path.Join(entry.Root, entry.APIVersion, handler.Path), hf)
			}
		}
	}

	router.Run(app.Gin.Address)
}

// MainEntry 主入口
func MainEntry() {

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	go func(mainQuitChan chan os.Signal) {
		time.Sleep(10 * time.Millisecond)
		customEntry()
	}(quit)

	fmt.Println("Please Press Ctrl + C Stop api-manager Service ...")
	<-quit
}

func preflight(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
	c.JSON(http.StatusOK, struct{}{})
}
