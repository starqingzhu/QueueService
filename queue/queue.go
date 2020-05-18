package queue

import (
	"QueueService/connManager"
	"QueueService/define"
	"QueueService/proto"
	"QueueService/utils"
	"container/list"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

var (
	EnqueueChan    chan define.ClientInfo // 缓存进入等待队列消息（仅仅是缓存并非真正缓存队列）
	QueryqueueChan chan define.ClientInfo // 查询玩家排队位置时发送
	QuitQueueChan  chan define.ClientInfo // 退出排队时发送(用户主动行为)
	//QuitGameChan   chan string            //  退出游戏时发送  时间不够先不管退出游戏的人（只关心上面退出等待队列的就ok了）
	ChangeInfoChan chan define.ChangeInfo // 在线人数变化时发送

)

var (
	WaitNumMap          sync.Map   // 缓存队列map（map[string]int32  用户id和其当前排队位置的Map）
	WaitList            *list.List // 登录等待队列 正在排队中的玩家
	OnWaitingQuePlayers int32      //正在排队中的玩家人数
	OnGamingPlayers     int32      // 正在游戏中的人数
	LoginWaitCurNum     int32      //当前队列的自增序列值
)

func Init() {
	EnqueueChan = make(chan define.ClientInfo, define.LOGIN_QUEUE_MAX_LEN)
	QueryqueueChan = make(chan define.ClientInfo, define.QUERY_LOGIN_QUEUE_MAX_LEN)
	QuitQueueChan = make(chan define.ClientInfo, define.LOGIN_QUEUE_QUIT_MAX_LEN)
	//QuitGameChan = make(chan string, define.LOGIN_GAME_QUIT_MAX_LEN)
	ChangeInfoChan = make(chan define.ChangeInfo, define.LOGIN_QUEUE_MAX_LEN)
	WaitList = list.New()
}

func OperateWaitList() {
	log.Printf("OperateWaitList enter------>>>>>>")

	//需要异步打印信息的处理
	go ListenChanges()

	//等待队列中消息任务处理
	go HandleLogin()

	//入等待队列前，进入等待队列消息的缓存处理
	go HandleEnqueueChan()

	//查询在队列在等待队列中位置 消息的任务的处理
	go HandleQueryqueueChan()

	//用户退出等待队列的处理
	go HandleQuitQueueChan()

	// 有用户退出游戏  这个预留

}

// 有玩家登陆
func HandleEnqueueChan() {
	for {
		select {
		case clientInfo := <-EnqueueChan:
			Enqueue(&clientInfo)
		}
	}
}

//玩家查询在等待队列位置
func HandleQueryqueueChan() {
	for {
		select {
		case clientInfo := <-QueryqueueChan:
			handleQuery(&clientInfo)
		}
	}
}

// 有用户退出排队
func HandleQuitQueueChan() {
	for {
		select {
		case clientInfo := <-QuitQueueChan:
			QuitQueue(&clientInfo)
		}
	}
}

func GetPlayerPosInfo(userName string) (res *define.PlayerQueInfo) {
	res = &define.PlayerQueInfo{}
	res.QueWaitPlayersNum = GetWaitingQuePlayersNum()
	absPos, ok := WaitNumMap.Load(userName)
	if ok {
		res.QuePlayerPos = getPlayerRelativeIndex(absPos.(int32))
	}
	res.PlayersGameIngNum = GetGamingPlayersNum()

	//兼容下，别个人位置比队列总长度还大了
	if res.QuePlayerPos > res.QueWaitPlayersNum {
		res.QueWaitPlayersNum = res.QuePlayerPos
	}
	res.UserName = userName

	return
}

// 新用户登陆
func Enqueue(clientInfo *define.ClientInfo) {
	if _, ok := WaitNumMap.Load(clientInfo.UserName); ok {
		WaitList.PushBack(*clientInfo)
		IncrWaitingQuePlayersNum()

		playerPosInfo := GetPlayerPosInfo(clientInfo.UserName)
		PrintWaitQueInfoChanged(playerPosInfo, define.QUE_CHANGE_REASON_PLAYER_ENTER)
	}
}

// 用户退出排队
// 从WaitList表删除此用户，原排队用户的等待位置不变
func QuitQueue(clientInfo *define.ClientInfo) {
	c, exist := connManager.ConnManager.Load(clientInfo.ConnAddr)
	if exist {
		playerPosInfo := GetPlayerPosInfo(clientInfo.UserName)
		quitRes := proto.NewQuitLoginQueRes(define.CMD_LOGIN_QUIT_RSP_NO,
			define.PROTO_VERSION,
			clientInfo.UserName,
			define.STATUS_SUCCESS,
		)
		(*(c.(*define.ConnInfo).Conn)).AsyncWrite(quitRes.ToBytes())

		if _, ok := WaitNumMap.Load(clientInfo.UserName); ok {
			DecrWaitingQuePlayersNum()
			WaitNumMap.Delete(clientInfo.UserName)
			PrintWaitQueInfoChanged(playerPosInfo, define.QUE_CHANGE_REASON_PLAYER_LEAVE)
		}

	}
}

// 向控制台打印实时更新的用户数据
func ListenChanges() {
	log.Printf("ListenChanges enter------>>>>>>")
	for {
		select {
		case info := <-ChangeInfoChan:
			log.Printf("正在游戏人数：%d  正在排队人数：%d reason：%d user：%s\n",
				info.PlayersGameIngNum,
				info.QueWaitPlayersNum,
				info.QuePlayerChangeReason,
				info.UserName,
			)
		}
	}
}

// 有用户退出游戏
//func QuitGame() {
//	userName := <-QuitGameChan
//	log.Printf("QuitQueue %s", userName)
//
//	if _, ok := WaitNumMap.Load(userName); ok {
//		WaitNumMap.Delete(userName)
//	} else {
//		DecrGamingPlayersNum()
//	}
//}

/*
处理登录等待队列中的 玩家登录消息
实现 登录逻辑处理
*/
func HandleLogin() {
	log.Printf("HandleLogin enter------>>>>>>")
	for {
		for e := WaitList.Front(); e != nil; e = e.Next() {
			clientInfo := e.Value.(define.ClientInfo)
			//先判断是否还在登录池
			if _, ok := WaitNumMap.Load(clientInfo.UserName); ok {
				//在等待缓存map删除对应玩家
				WaitNumMap.Delete(clientInfo.UserName)
				//对等待队列减1
				DecrWaitingQuePlayersNum()
				//对游戏人数加1
				IncrGamingPlayersNum()

				//打印同步信息
				playerPosInfo := GetPlayerPosInfo(clientInfo.UserName)
				PrintWaitQueInfoChanged(playerPosInfo, define.QUE_CHANGE_REASON_PLAYER_GAMING)

				//异步通知登录成功
				//先查连接是否还在
				c, exist := connManager.ConnManager.Load(clientInfo.ConnAddr)
				if exist {
					token := utils.GetToken()
					notifyInfo := proto.NewLoginNotify(define.CMD_LOGIN_NOTIFY_NO,
						define.PROTO_VERSION,
						clientInfo.UserName,
						token)

					(*(c.(*define.ConnInfo).Conn)).AsyncWrite(notifyInfo.ToBytes())
				}

			}

			WaitList.Remove(e)
		}

		time.Sleep(define.LOGIN_HANDLE_WAIT_TIME * time.Millisecond)
	}
}

func handleQuery(clientInfo *define.ClientInfo) {
	c, exist := connManager.ConnManager.Load(clientInfo.ConnAddr)
	if exist {
		playerQueInfo := GetPlayerPosInfo(clientInfo.UserName)
		queryRes := proto.NewQueryPlayerLoginQuePosRes(define.CMD_QUERY_PLAYER_LOGIN_QUE_POS_RSP_NO,
			define.PROTO_VERSION,
			playerQueInfo.QueWaitPlayersNum,
			playerQueInfo.PlayersGameIngNum,
			playerQueInfo.QuePlayerPos)

		(*(c.(*define.ConnInfo).Conn)).AsyncWrite(queryRes.ToBytes())
	}
}

func getPlayerRelativeIndex(playerLoginNum int32) int32 {
	var relativeIndex int32

	curNum := GetLoginCurNum()
	if curNum >= playerLoginNum {
		relativeIndex = curNum - playerLoginNum + 1
	} else {
		relativeIndex = define.LOGIN_QUEUE_MAX_LEN - playerLoginNum + curNum
	}
	return relativeIndex
}

/*
等待队列成员变动时打印信息
*/
func PrintWaitQueInfoChanged(playerQueInfo *define.PlayerQueInfo, reason uint16) {
	changeInfo := define.ChangeInfo{
		QueWaitPlayersNum:     playerQueInfo.QueWaitPlayersNum,
		PlayersGameIngNum:     playerQueInfo.PlayersGameIngNum,
		UserName:              playerQueInfo.UserName,
		QuePlayerChangeReason: reason,
	}
	ChangeInfoChan <- changeInfo
}

/*
检查用户是否已经在等待队列中

true 表示在队列
false 表示不在队列
*/
func CheckPlayerIsInWaitQue(userName string) (res bool) {
	if _, ok := WaitNumMap.Load(userName); ok {
		res = true
	}
	return
}

func GetGamingPlayersNum() int32 {
	n := atomic.LoadInt32(&OnGamingPlayers)
	return n
}

func IncrGamingPlayersNum() int32 {
	n := atomic.AddInt32(&OnGamingPlayers, 1)
	return n
}

func DecrGamingPlayersNum() int {
	n := atomic.AddInt32(&OnGamingPlayers, -1)
	return int(n)
}

func GetWaitingQuePlayersNum() int32 {
	n := atomic.LoadInt32(&OnWaitingQuePlayers)
	return n
}

func IncrWaitingQuePlayersNum() int32 {
	n := atomic.AddInt32(&OnWaitingQuePlayers, 1)
	return n
}

func DecrWaitingQuePlayersNum() int {
	n := atomic.AddInt32(&OnWaitingQuePlayers, -1)
	return int(n)
}

func GetLoginCurNum() int32 {
	return atomic.LoadInt32(&LoginWaitCurNum)
}

func IncrLoginCurNum() int32 {
	n := atomic.AddInt32(&LoginWaitCurNum, 1)
	if n > define.LOGIN_MAX_NUM {
		atomic.StoreInt32(&LoginWaitCurNum, int32(1))
		n = 1
	}
	return n
}
