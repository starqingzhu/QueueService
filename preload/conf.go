package preload

import (
	"QueueService/define"
	"github.com/spf13/viper"
)

var Conf define.ConfInfo

func init() {
	confFile := confPath + "/config.json"

	viper.SetConfigFile(confFile)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	Conf.TcpServer.TcpMaxConnCount = viper.GetInt(`tcpServer.tcpMaxConnCount`)
	Conf.TcpServer.TcpPort = viper.GetInt(`tcpServer.tcpPort`)

	Conf.Stress.Switch = viper.GetBool(`stress.switch`)
	Conf.Stress.HttpAddr = viper.GetString(`stress.httpAddr`)

	Conf.Server.GoMaxProcsNum = viper.GetInt(`server.goMaxProcsNum`)
	Conf.Server.Multicore = viper.GetBool(`server.multicore`)

}
