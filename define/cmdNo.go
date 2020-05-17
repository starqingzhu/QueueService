package define

const (
	CMD_LOGIN_REQ_NO    = 10001 //登录请求
	CMD_LOGIN_RES_NO    = 20001 //登录返回
	CMD_LOGIN_NOTIFY_NO = 30001 //登录异步通知

	CMD_QUERY_PLAYER_LOGIN_QUE_POS_REQ_NO = 10002 //查询玩家在登录队列中的位置 请求
	CMD_QUERY_PLAYER_LOGIN_QUE_POS_RSP_NO = 20002 //查询玩家在登录队列中的位置 返回

	CMD_LOGIN_QUIT_REQ_NO = 10003 //玩家退出等待队列请求
	CMD_LOGIN_QUIT_RSP_NO = 20003 //玩家退出等待队列回复
)
