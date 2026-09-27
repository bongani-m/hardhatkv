package main

import (
	"fmt"
	"os"

	"github.com/bongani-m/hardhatkv/cluster"
	"github.com/bongani-m/hardhatkv/config"
	"github.com/bongani-m/hardhatkv/database"
	idatabase "github.com/bongani-m/hardhatkv/interface/database"
	"github.com/bongani-m/hardhatkv/lib/logger"
	"github.com/bongani-m/hardhatkv/lib/utils"
	"github.com/bongani-m/hardhatkv/redis/server/gnet"
	stdserver "github.com/bongani-m/hardhatkv/redis/server/std"
)

var banner = `
 _   _               _ _           _   _  ___     __
| | | | __ _ _ __ __| | |__   __ _| |_| |/ \ \   / /
| |_| |/ _' | '__/ _' | '_ \ / _' | __| ' / \ \ / /
|  _  | (_| | | | (_| | | | | (_| | |_| . \  \ V /
|_| |_|\__,_|_|  \__,_|_| |_|\__,_|\__|_|\_\  \_/
`

var defaultProperties = &config.ServerProperties{
	Bind:           "0.0.0.0",
	Port:           6399,
	AppendOnly:     false,
	AppendFilename: "",
	MaxClients:     1000,
	RunID:          utils.RandString(40),
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	return err == nil && !info.IsDir()
}

func main() {
	print(banner)
	logger.Setup(&logger.Settings{
		Path:       "logs",
		Name:       "hardhatkv",
		Ext:        "log",
		TimeFormat: "2006-01-02",
	})
	configFilename := os.Getenv("CONFIG")
	if configFilename == "" {
		if fileExists("redis.conf") {
			config.SetupConfig("redis.conf")
		} else {
			config.Properties = defaultProperties
		}
	} else {
		config.SetupConfig(configFilename)
	}
	listenAddr := fmt.Sprintf("%s:%d", config.Properties.Bind, config.Properties.Port)

	var err error
	if config.Properties.UseGnet {
		var db idatabase.DB
		if config.Properties.ClusterEnable {
			db = cluster.MakeCluster()
		} else {
			db = database.NewStandaloneServer()
		}
		server := gnet.NewGnetServer(db)
		err = server.Run(listenAddr)
	} else {
		handler := stdserver.MakeHandler()
		err = stdserver.Serve(listenAddr, handler)
	}
	if err != nil {
		logger.Errorf("start server failed: %v", err)
	}
}
