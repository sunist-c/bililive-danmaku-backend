package message

type Command string

const (
	CommandNewDanmaku        Command = "DANMU_MSG"                     // 普通弹幕信息
	CommandNewGift           Command = "SEND_GIFT"                     // 普通的礼物，不包含礼物连击
	CommandWelcome           Command = "WELCOME"                       // 欢迎VIP
	CommandWelcomeGuard      Command = "WELCOME_GUARD"                 // 欢迎房管
	CommandWelcomeVip        Command = "ENTRY_EFFECT"                  // 欢迎舰长等头衔
	CommandRoomFocusedChange Command = "ROOM_REAL_TIME_MESSAGE_UPDATE" // 房间关注数变动
)
