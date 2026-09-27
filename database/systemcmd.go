package database

import (
	"fmt"
	"github.com/bongani-m/hardhatkv/config"
	"github.com/bongani-m/hardhatkv/interface/redis"
	"github.com/bongani-m/hardhatkv/redis/protocol"
	"github.com/bongani-m/hardhatkv/tcp"
	"os"
	"runtime"
	"strings"
	"time"
)

// Ping the server
func Ping(c redis.Connection, args [][]byte) redis.Reply {
	if len(args) == 0 {
		return &protocol.PongReply{}
	} else if len(args) == 1 {
		return protocol.MakeStatusReply(string(args[0]))
	} else {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'ping' command")
	}
}

// Info the information of the HardhatKV server returned by the INFO command
func Info(db *Server, args [][]byte) redis.Reply {
	if len(args) == 0 {
		infoCommandList := [...]string{"server", "client", "cluster", "keyspace"}
		var allSection []byte
		for _, s := range infoCommandList {
			allSection = append(allSection, GenInfoString(s, db)...)
		}
		return protocol.MakeBulkReply(allSection)
	} else if len(args) == 1 {
		section := strings.ToLower(string(args[0]))
		switch section {
		case "server":
			reply := GenInfoString("server", db)
			return protocol.MakeBulkReply(reply)
		case "client":
			return protocol.MakeBulkReply(GenInfoString("client", db))
		case "cluster":
			return protocol.MakeBulkReply(GenInfoString("cluster", db))
		case "keyspace":
			return protocol.MakeBulkReply(GenInfoString("keyspace", db))
		default:
			return protocol.MakeErrReply("Invalid section for 'info' command")
		}
	}
	return protocol.MakeArgNumErrReply("info")
}

// Auth validate client's password
func Auth(c redis.Connection, args [][]byte) redis.Reply {
	if len(args) != 1 {
		return protocol.MakeErrReply("ERR wrong number of arguments for 'auth' command")
	}
	if config.Properties.RequirePass == "" {
		return protocol.MakeErrReply("ERR Client sent AUTH, but no password is set")
	}
	passwd := string(args[0])
	c.SetPassword(passwd)
	if config.Properties.RequirePass != passwd {
		return protocol.MakeErrReply("ERR invalid password")
	}
	return &protocol.OkReply{}
}

func isAuthenticated(c redis.Connection) bool {
	if config.Properties.RequirePass == "" {
		return true
	}
	return c.GetPassword() == config.Properties.RequirePass
}

func DbSize(c redis.Connection, db *Server) redis.Reply {
	keys, _ := db.GetDBSize(c.GetDBIndex())
	return protocol.MakeIntReply(int64(keys))
}

func GenInfoString(section string, db *Server) []byte {
	startUpTimeFromNow := getRunningTime()
	switch section {
	case "server":
		s := fmt.Sprintf("# Server\r\n"+
			"hardhatkv_version:%s\r\n"+
			//"hardhatkv_git_sha1:%s\r\n"+
			//"hardhatkv_git_dirty:%d\r\n"+
			//"hardhatkv_build_id:%s\r\n"+
			"hardhatkv_mode:%s\r\n"+
			"os:%s %s\r\n"+
			"arch_bits:%d\r\n"+
			//"multiplexing_api:%s\r\n"+
			"go_version:%s\r\n"+
			"process_id:%d\r\n"+
			"run_id:%s\r\n"+
			"tcp_port:%d\r\n"+
			"uptime_in_seconds:%d\r\n"+
			"uptime_in_days:%d\r\n"+
			//"hz:%d\r\n"+
			//"lru_clock:%d\r\n"+
			"config_file:%s\r\n",
			version,
			//TODO,
			//TODO,
			//TODO,
			getRunningMode(),
			runtime.GOOS, runtime.GOARCH,
			32<<(^uint(0)>>63),
			//TODO,
			runtime.Version(),
			os.Getpid(),
			config.Properties.RunID,
			config.Properties.Port,
			startUpTimeFromNow,
			startUpTimeFromNow/time.Duration(3600*24),
			//TODO,
			//TODO,
			config.GetConfigFilePath())
		return []byte(s)
	case "client":
		s := fmt.Sprintf("# Clients\r\n"+
			"connected_clients:%d\r\n",
			//"client_recent_max_input_buffer:%d\r\n"+
			//"client_recent_max_output_buffer:%d\r\n"+
			//"blocked_clients:%d\n",
			tcp.ClientCounter,
			//TODO,
			//TODO,
			//TODO,
		)
		return []byte(s)
	case "cluster":
		if getRunningMode() == config.ClusterMode {
			s := fmt.Sprintf("# Cluster\r\n"+
				"cluster_enabled:%s\r\n",
				"1",
			)
			return []byte(s)
		} else {
			s := fmt.Sprintf("# Cluster\r\n"+
				"cluster_enabled:%s\r\n",
				"0",
			)
			return []byte(s)
		}
	case "keyspace":
		dbCount := config.Properties.Databases
		var serv []byte
		for i := 0; i < dbCount; i++ {
			keys, expiresKeys := db.GetDBSize(i)
			if keys != 0 {
				ttlSampleAverage := db.GetAvgTTL(i, 20)
				serv = append(serv, getDbSize(i, keys, expiresKeys, ttlSampleAverage)...)
			}
		}
		prefix := []byte("# Keyspace\r\n")
		keyspaceInfo := append(prefix, serv...)
		return keyspaceInfo
	}
	return []byte("")
}

// getRunningMode return the HardhatKV running mode
func getRunningMode() string {
	if config.Properties.ClusterEnable {
		return config.ClusterMode
	} else {
		return config.StandaloneMode
	}
}

// getRunningTime return the running time of HardhatKV
func getRunningTime() time.Duration {
	return time.Since(config.EachTimeServerInfo.StartUpTime) / time.Second
}

func getDbSize(dbIndex, keys, expiresKeys int, ttl int64) []byte {
	s := fmt.Sprintf("db%d:keys=%d,expires=%d,avg_ttl=%d\r\n",
		dbIndex, keys, expiresKeys, ttl)
	return []byte(s)
}
