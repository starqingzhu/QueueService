package define

const (
	LOGIN_QUEUE_MAX_LEN      = 10000 //todo 之后走配置文件
	LOGIN_QUEUE_QUIT_MAX_LEN = 100   //todo 之后走配置文件
	LOGIN_GAME_QUIT_MAX_LEN  = 10    //todo 之后走配置文件
	LOGIN_MAX_NUM            = 10000 //todo 之后走配置

	LOGIN_HANDLE_WAIT_TIME = 50 //单位毫秒
)

//等待队列玩家变动标示
const (
	QUE_CHANGE_REASON_PLAYER_ENTER = 1 //进入
	QUE_CHANGE_REASON_PLAYER_LEAVE = 2 //离开
	QUE_CHANGE_REASON_PLAYER_GAMING = 3 //开始游戏
)

type (
	ClientInfo struct {
		UserName string //用户名
		ConnAddr string //长链接地址
	}

	PlayerQueInfo struct {
		QueWaitPlayersNum int32  //排队人数
		PlayersGameIngNum int32  //玩游戏人数
		QuePlayerPos      int32  //在队列中位置
		UserName          string //用户名
	}

	ChangeInfo struct {
		QueWaitPlayersNum     int32  //排队人数
		PlayersGameIngNum     int32  //玩游戏人数
		//QuePlayerPos          int32  //在队列中位置
		UserName              string //用户名
		QuePlayerChangeReason uint16 //队列变动原因
	}
)
