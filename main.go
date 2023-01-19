package main

import (
	"TVHelper/global"
	"TVHelper/internal/routers"
	"TVHelper/pkg/logging"
	"TVHelper/pkg/setting"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

var dir string

func init() {
	err := setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}
	logging.Init()
	currDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	flag.StringVar(&dir, "d", currDir, "配置目录")
	flag.Parse()
	err = os.Chdir(dir)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	global.Logger.Info("TVHelper starting...",
		zap.String("port", global.ServerSetting.HttpPort),
		zap.String("dir", dir))

	gin.SetMode(global.ServerSetting.RunMode)
	router := routers.NewRouter()
	s := &http.Server{
		Addr:           ":" + global.ServerSetting.HttpPort,
		Handler:        router,
		ReadTimeout:    global.ServerSetting.ReadTimeout,
		WriteTimeout:   global.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	if err := s.ListenAndServe(); err != nil {
		global.Logger.Fatal("startup service failed...", zap.Error(err))
	}
}

func setupSetting() error {
	newSetting, err := setting.NewSetting()
	if err != nil {
		return err
	}
	err = newSetting.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		return err
	}
	err = newSetting.ReadSection("Log", &global.LogSetting)
	if err != nil {
		return err
	}

	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second
	return nil
}
