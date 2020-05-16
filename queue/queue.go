package queue

import (
	"QueueService/connManager"
	"QueueService/define"
	"QueueService/proto"
	"QueueService/utils"
	"container/list"
	"github.com/panjf2000/gnet"
	"log"
	"sync/atomic"
	"time"
)

var (
	EnqueueChan    chan define.ClientInfo // 进入排队时发送
	QueryqueueChan chan define.ClientInfo // 查询玩家排队位置时发送
	QuitQueueChan  chan string            // 退出排队时发送(用户主动行为)
	QuitGameChan   chan string            //  退出游戏时发送
	ChangeInfoChan chan define.ChangeInfo // 在线人数变化时发送

)

var (
	WaitNumMap      map[string]int32 // 用户id和其当前排队位置的Map
	WaitList        *list.List       // 正在排队中的用户
	OnGamingPlayers int32            // 正在游戏中的人数
	LoginWaitCurNum int32            //当前队列的自增序列值
)

func Init() {
	EnqueueChan = make(chan define.ClientInfo, define.LOGIN_QUEUE_MAX_LEN)
	QueryqueueChan = make(chan define.ClientInfo, define.QUERY_LOGIN_QUEUE_MAX_LEN)
	QuitQueueChan = make(chan string, define.LOGIN_QUEUE_QUIT_MAX_LEN)
	QuitGameChan = make(chan string, define.LOGIN_GAME_QUIT_MAX_LEN)
	ChangeInfoChan = make(chan define.ChangeInfo, define.LOGIN_QUEUE_MAX_LEN)
	WaitNumMap = make(map[string]int32, define.LOGIN_QUEUE_MAX_LEN)
	WaitList = list.New()
}

func OperateWaitList() {
	log.Printf("OperateWaitList enter------>>>>>>")
	for {
		select {
		// 有玩家登陆
		case clientInfo := <-EnqueueChan:
			Enqueue(&clientInfo)

		//玩家查询在等待队列位置
		case clientInfo := <-QueryqueueChan:
			handleQuery(&clientInfo)
			//// 有用户退出排队
			//case userName := <-QuitQueueChan:
			//	QuitQueue(userName)
			//// 有用户退出游戏
			//case <-QuitGameChan:
			//	QuitGame()
			//// 有用户退出，具体是退出排队还是游戏在Quit(）内进一步判断
			//case userName := <-QuitChan:
			//	Quit(userName)
		}
	}
}

func GetPlayerPosInfo(userName string) (res *define.PlayerQueInfo) {
	res = &define.PlayerQueInfo{}
	res.QueWaitPlayersNum = int32(len(WaitNumMap))
	absPos, ok := WaitNumMap[userName]
	if ok {
		res.QuePlayerPos = getPlayerRelativeIndex(absPos)
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
	//log.Printf("Enqueue clientInfo:%v", clientInfo)
	WaitNumMap[clientInfo.UserName] = IncrLoginCurNum()
	WaitList.PushFront(*clientInfo)

	playerPosInfo := GetPlayerPosInfo(clientInfo.UserName)
	PrintWaitQueInfoChanged(playerPosInfo, define.QUE_CHANGE_REASON_PLAYER_ENTER)
}

// 用户退出排队
// 从WaitList表删除此用户，原排队用户的等待位置不变
func QuitQueue(userName string) {
	log.Printf("QuitQueue %s", userName)
	playerPosInfo := GetPlayerPosInfo(userName)
	PrintWaitQueInfoChanged(playerPosInfo, define.QUE_CHANGE_REASON_PLAYER_LEAVE)
	delete(WaitNumMap, userName)
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
func QuitGame() {
	userName := <-QuitGameChan
	log.Printf("QuitQueue %s", userName)
	if _, ok := WaitNumMap[userName]; ok {
		delete(WaitNumMap, userName)
	} else {
		DecrGamingPlayersNum()
	}
}

func HandleLogin() {
	log.Printf("HandleLogin enter------>>>>>>")
	for {
		for e := WaitList.Front(); e != nil; e = e.Next() {
			clientInfo := e.Value.(define.ClientInfo)
			//先判断是否还在登录池
			if _, ok := WaitNumMap[clientInfo.UserName]; ok {
				//在等待缓存map删除对应玩家
				delete(WaitNumMap, clientInfo.UserName)
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

					c.(gnet.Conn).AsyncWrite(notifyInfo.ToBytes())
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

		c.(gnet.Conn).AsyncWrite(queryRes.ToBytes())
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
