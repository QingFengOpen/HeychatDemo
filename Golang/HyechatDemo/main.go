package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/viper"

	"HyechatDemo/domain/eventhandler"
	"HyechatDemo/heybotclient"
	"fmt"
)

type Config struct {
	Token string `yaml:"token"`
}

func initConfig() *Config {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("无法读取配置文件: %s", err))
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("无法解析配置文件: %s", err))
	}
	return &config
}

func main() {
	config := initConfig()
	token := config.Token
	handler := eventhandler.New(token)
	client := heybotclient.NewWebSocketClient(context.Background(), token, nil, handler)
	if err := client.Connect(); err != nil {
		log.Fatalf("Failed to connect to WebSocket server: %v", err)
	}
	defer client.Close()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	// 监听退出信号
	<-quit
	// 当接收到退出信号时，关闭 WebSocket 连接
	log.Println("Received exit signal, closing WebSocket connection...")
}
