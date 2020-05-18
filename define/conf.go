package define

type (
	ConfInfo struct {
		TcpServer ConfTcpServer
		Server    ConfServer
		Stress    ConfStress
	}
	// 网络相关配置
	ConfTcpServer struct {
		TcpMaxConnCount int
		TcpPort         int
	}
	//服务器相关配置
	ConfServer struct {
		GoMaxProcsNum int
		Multicore     bool
	}
	//压测 性能分析配置
	ConfStress struct {
		Switch   bool
		HttpAddr string
	}
)
