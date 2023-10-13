package main

import (
	"redis-based-on-go/config"
	"redis-based-on-go/database"
	"redis-based-on-go/log"
	"redis-based-on-go/redis/server"
)

func main() {
	config.InitConfig()
	mylog.InitLog()
	database.InitDBCommand()
	err := server.StartServer(config.Conf.Server)
	if err == nil {
		mylog.Logger.Info("退出成功")
	} else {
		mylog.Logger.Error(err)
	}

}
