package constant

const (

	// group admin
	//	OrdinaryMember = 0
	//	GroupOwner     = 1
	//	Administrator  = 2
	// group application
	//	Application      = 0
	//	AgreeApplication = 1

	// friend related
	BlackListFlag         = 1
	ApplicationFriendFlag = 0
	FriendFlag            = 1
	RefuseFriendFlag      = -1

	// Redis Key For Status
	UserStatusKey    = "user_status"
	AdminStatusKey   = "admin_status"
	UserIPandStatus  = "ip_status:user"
	UserGroupIDCache = "USER_GROUPID_LIST"

	// Websocket Protocol
	WSGetNewestSeq     = 1001
	WSPullMsgBySeqList = 1002
	WSSendMsg          = 1003
	WSSendSignalMsg    = 1004
	WSPushMsg          = 2001
	WSKickOnlineMsg    = 2002
	WsLogoutMsg        = 2003
	WsSyncDataMsg      = 2004
	WSDataError        = 3001

	// /ContentType
	// UserRelated
	Text                = 101 // 文本 Text
	Picture             = 102 // 图片 Picture
	Voice               = 103 // 语音 Voice
	Video               = 104 // 视频 Video
	File                = 105 // 文件 File
	AtText              = 106 // @消息 @Text
	Merger              = 107 // 合并消息 Message Merge
	Card                = 108 // 名片 Name Card
	Location            = 109 // 位置 Location
	Custom              = 110 // 自定义 Custom
	Revoke              = 111 // 撤回 Revoke
	HasReadReceipt      = 112 // 已读回执 Has Readed Receipt
	Typing              = 113 // 输入中 Typing
	Quote               = 114 // 引用 Quote
	GroupHasReadReceipt = 116 // 群消息已读回执 Group Message Has Read Receipt
	Common              = 200 // 公共消息 Common
	GroupMsg            = 201 // 群消息 Group Message
	SignalMsg           = 202 // 信号消息 Signal Message

	// SysRelated
	NotificationBegin                     = 1000 // 开始通知
	DeleteMessageNotification             = 1100 // 删除消息通知
	FriendApplicationApprovedNotification = 1201 // 接受好友 Friend Approved add_friend_response
	FriendApplicationRejectedNotification = 1202 // 拒绝好友 Friend Rejected add_friend_response
	FriendApplicationNotification         = 1203 // 申请添加好友 add_friend
	FriendAddedNotification               = 1204 // 添加好友成功 Added Friend
	FriendDeletedNotification             = 1205 // 删除好友 Delete Friend delete_friend
	FriendRemarkSetNotification           = 1206 // 好友备注 Friend Remark set_friend_remark?
	BlackAddedNotification                = 1207 // 拉黑 Add Blacklist add_black
	BlackDeletedNotification              = 1208 // 解除黑名单 Delete Blacklist remove_black

	ConversationOptChangeNotification = 1300 // 会话改变 Conversation Changed change conversation opt

	UserNotificationBegin       = 1301
	UserInfoUpdatedNotification = 1303 // 用户信息更新 User Info Updated
	UserNotificationEnd         = 1399
	OANotification              = 1400 // OA通知 OA Notification

	GroupNotificationBegin = 1500

	GroupCreatedNotification                 = 1501 // 建群 Create Group
	GroupInfoSetNotification                 = 1502 // 设置群信息 Group Info Setting
	JoinGroupApplicationNotification         = 1503 // 加群 Join Group
	MemberQuitNotification                   = 1504 // 退群 Exit Group
	GroupApplicationAcceptedNotification     = 1505 // 接受入群 Group Accepted
	GroupApplicationRejectedNotification     = 1506 // 拒绝入群 Group Rejected
	GroupOwnerTransferredNotification        = 1507 // 群主转让 Group Owner Transferred
	MemberKickedNotification                 = 1508 // 踢人 Member Kicked
	MemberInvitedNotification                = 1509 // 拉人 Member Invited
	MemberEnterNotification                  = 1510 // 进群 Member Enter
	GroupDismissedNotification               = 1511 // 群解散 Group Dismissed
	GroupMemberMutedNotification             = 1512 // 成员禁言 Member Muted
	GroupMemberCancelMutedNotification       = 1513 // 取消成员禁言 Member Cancel Muted
	GroupMutedNotification                   = 1514 // 群禁言 Group Muted
	GroupCancelMutedNotification             = 1515 // 取消群禁言 Group Cancel Muted
	GroupMemberInfoSetNotification           = 1516 // 设置成员信息 Member Info Setting
	GroupMemberSetToAdminNotification        = 1517 // 设置管理员 Admin Setting
	GroupMemberSetToOrdinaryUserNotification = 1518 // 取消管理员 Cancel Admin
	GroupAnnouncementNotification            = 1519 // 发布群公告 Post Group Announcement

	SignalingNotificationBegin = 1600
	SignalingNotification      = 1601
	SignalingNotificationEnd   = 1649

	SuperGroupNotificationBegin  = 1650
	SuperGroupUpdateNotification = 1651
	SuperGroupNotificationEnd    = 1699

	ConversationPrivateChatNotification = 1701

	OrganizationChangedNotification = 1801

	WorkMomentNotificationBegin = 1900
	WorkMomentNotification      = 1901

	NotificationEnd = 3000

	// status
	MsgNormal  = 1
	MsgDeleted = 4

	// MsgStatus
	MsgStatusDefault     = 0
	MsgStatusSending     = 1
	MsgStatusSendSuccess = 2
	MsgStatusSendFailed  = 3
	MsgStatusHasDeleted  = 4
	MsgStatusRevoked     = 5
	MsgStatusFiltered    = 6

	// MsgFrom
	UserMsgType = 100
	SysMsgType  = 200

	// SessionType
	SingleChatType       = 1
	GroupChatType        = 2
	SuperGroupChatType   = 3
	NotificationChatType = 4
	// token
	NormalToken  = 0
	InValidToken = 1
	KickedToken  = 2
	ExpiredToken = 3

	// MultiTerminalLogin
	// Full-end login, but the same end is mutually exclusive
	AllLoginButSameTermKick = 1
	// Only one of the endpoints can log in
	SingleTerminalLogin = 2
	// The web side can be online at the same time, and the other side can only log in at one end
	WebAndOther = 3
	// The PC side is mutually exclusive, and the mobile side is mutually exclusive, but the web side can be online at the same time
	PcMobileAndWeb = 4

	OnlineStatus  = "online"
	OfflineStatus = "offline"
	Registered    = "registered"
	UnRegistered  = "unregistered"

	// MsgReceiveOpt
	ReceiveMessage          = 0
	NotReceiveMessage       = 1
	ReceiveNotNotifyMessage = 2

	// OptionsKey
	IsHistory                  = "history"
	IsPersistent               = "persistent"
	IsOfflinePush              = "offlinePush"
	IsUnreadCount              = "unreadCount"
	IsConversationUpdate       = "conversationUpdate"
	IsSenderSync               = "senderSync"
	IsNotPrivate               = "notPrivate"
	IsSenderConversationUpdate = "senderConversationUpdate"
	IsSenderNotificationPush   = "senderNotificationPush"
	IsSyncToLocalDataBase      = "syncToLocalDataBase"

	// GroupStatus
	GroupOk              = 0
	GroupBanChat         = 1
	GroupStatusDismissed = 2
	GroupStatusMuted     = 3

	// GroupType
	NormalGroup     = 0
	SuperGroup      = 1
	DepartmentGroup = 2

	GroupBaned          = 3
	GroupBanPrivateChat = 4

	// UserJoinGroupSource
	JoinByAdmin = 1

	// Minio
	MinioDurationTimes = 3600

	// verificationCode used for
	VerificationCodeForRegister       = 1
	VerificationCodeForReset          = 2
	VerificationCodeForRegisterSuffix = "_forRegister"
	VerificationCodeForResetSuffix    = "_forReset"

	// callbackCommand
	CallbackBeforeSendSingleMsgCommand = "callbackBeforeSendSingleMsgCommand"
	CallbackAfterSendSingleMsgCommand  = "callbackAfterSendSingleMsgCommand"
	CallbackBeforeSendGroupMsgCommand  = "callbackBeforeSendGroupMsgCommand"
	CallbackAfterSendGroupMsgCommand   = "callbackAfterSendGroupMsgCommand"
	CallbackWordFilterCommand          = "callbackWordFilterCommand"
	CallbackUserOnlineCommand          = "callbackUserOnlineCommand"
	CallbackUserOfflineCommand         = "callbackUserOfflineCommand"
	CallbackOfflinePushCommand         = "callbackOfflinePushCommand"
	// callback actionCode
	ActionAllow     = 0
	ActionForbidden = 1
	// callback callbackHandleCode
	CallbackHandleSuccess = 0
	CallbackHandleFailed  = 1

	// minioUpload
	OtherType = 1
	VideoType = 2
	ImageType = 3

	// workMoment permission
	WorkMomentPublic            = 0
	WorkMomentPrivate           = 1
	WorkMomentPermissionCanSee  = 2
	WorkMomentPermissionCantSee = 3

	// workMoment sdk notification type
	WorkMomentCommentNotification = 0
	WorkMomentLikeNotification    = 1
	WorkMomentAtUserNotification  = 2
)
const (
	SuperGroupTableName               = "local_super_groups"
	SuperGroupErrChatLogsTableNamePre = "local_sg_err_chat_logs_"
	SuperGroupChatLogsTableNamePre    = "local_sg_chat_logs_"
)
const (
	KeywordMatchOr  = 0
	KeywordMatchAnd = 1
)
const (
	AddConOrUpLatMsg          = 2
	UnreadCountSetZero        = 3
	IncrUnread                = 5
	TotalUnreadMessageChanged = 6
	UpdateFaceUrlAndNickName  = 7
	UpdateLatestMessageChange = 8
	ConChange                 = 9
	NewCon                    = 10

	HasRead = 1
	NotRead = 0

	IsFilter  = 1
	NotFilter = 0
)
const (
	AtAllString       = "AtAllTag"
	AtNormal          = 0
	AtMe              = 1
	AtAll             = 2
	AtAllAtMe         = 3
	GroupNotification = 4
)

var ContentType2PushContent = map[int64]string{
	Picture:   "[图片]",
	Voice:     "[语音]",
	Video:     "[视频]",
	File:      "[文件]",
	Text:      "你收到了一条文本消息",
	AtText:    "[有人@你]",
	GroupMsg:  "你收到一条群聊消息",
	Common:    "你收到一条新消息",
	SignalMsg: "音视频通话邀请",
}

const (
	FieldRecvMsgOpt    = 1
	FieldIsPinned      = 2
	FieldAttachedInfo  = 3
	FieldIsPrivateChat = 4
	FieldGroupAtType   = 5
	FieldIsNotInGroup  = 6
	FieldEx            = 7
	FieldUnread        = 8
)

const (
	AppOrdinaryUsers = 1
	AppAdmin         = 2

	GroupOrdinaryUsers = 1
	GroupOwner         = 2
	GroupAdmin         = 3

	GroupResponseAgree  = 1
	GroupResponseRefuse = -1

	FriendResponseAgree  = 1
	FriendResponseRefuse = -1

	Male   = 1
	Female = 2
)

// Verification code type
const (
	SendMsgRegister      = 1
	SendMsgResetPassword = 2
)

const (
	InviteCodeStateValid   = 1
	InviteCodeStateInvalid = 2
	InviteCodeStateDelete  = 3
)
const (
	InviteChannelCodeStateValid   = 1
	InviteChannelCodeStateInvalid = 2
	InviteChannelCodeStateDelete  = 3
)

const (
	UnreliableNotification    = 1
	ReliableNotificationNoMsg = 2
	ReliableNotificationMsg   = 3
)

const (
	ConfigInviteCodeBaseLinkKey = "invite_code_base_link"

	ConfigInviteCodeIsOpenKey   = "invite_code_is_open"
	ConfigInviteCodeIsOpenTrue  = 1
	ConfigInviteCodeIsOpenFalse = 0

	ConfigInviteCodeIsLimitKey   = "invite_code_is_limit"
	ConfigInviteCodeIsLimitTrue  = 1
	ConfigInviteCodeIsLimitFalse = 0

	ConfigChannelCodeIsOpenKey    = "channel_code_is_open"
	ConfigChannelCodeIsOpenTrue   = 1
	ConfigChannelCodeIsOpenFalse  = 0
	ConfigChannelCodeIsLimitKey   = "channel_code_is_limit"
	ConfigChannelCodeIsLimitTrue  = 1
	ConfigChannelCodeIsLimitFalse = 0
)

const (
	AllowRegisterByUuid = "allow_register_by_uuid"
	AllowGuestLogin     = "allow_guest_login"
)

const (
	UserRegisterSourceTypeOfficial = 1 // 官方注册
	UserRegisterSourceTypeInvite   = 2 // 邀请注册
	UserRegisterSourceTypeChannel  = 3 // 渠道注册
)

const (
	PrivateKeyHex          = 1 // HEX
	PrivateKeyCompressed   = 2 // Compressed
	PrivateKeyUnCompressed = 3 // Un Compressed
)

const FriendAcceptTip = "You have successfully become friends, so start chatting"

func GroupIsBanChat(status int32) bool {
	if status != GroupStatusMuted {
		return false
	}
	return true
}

func GroupIsBanPrivateChat(status int32) bool {
	if status != GroupBanPrivateChat {
		return false
	}
	return true
}

const AdminAPILogFileName = "admin_api.log"
const AdminRPCLogFileName = "admin_rpc.log"
const WalletAPILogFileName = "wallet_api.log"
const WalletRPCLogFileName = "wallet_rpc.log"
const EthRPCLogFileName = "eth_rpc.log"
const BtcRPCLogFileName = "btc_rpc.log"
const TrcRPCLogFileName = "trc_rpc.log"
const SQLiteLogFileName = "sqlite.log"
const MySQLLogFileName = "mysql.log"
const CronLogFileName = "cron.log"
const PushLogFileName = "push.log"

const StatisticsTimeInterval = 60

const (
	SyncCreateGroup               = "SyncCreateGroup"
	SyncDeleteGroup               = "SyncDeleteGroup"
	SyncUpdateGroup               = "SyncUpdateGroup"
	SyncInvitedGroupMember        = "SyncInvitedGroupMember"
	SyncKickedGroupMember         = "SyncKickedGroupMember"
	SyncMuteGroupMember           = "SyncMuteGroupMember"
	SyncCancelMuteGroupMember     = "SyncCancelMuteGroupMember"
	SyncGroupMemberInfo           = "SyncGroupMemberInfo"
	SyncMuteGroup                 = "SyncMuteGroup"
	SyncCancelMuteGroup           = "SyncCancelMuteGroup"
	SyncConversation              = "SyncConversation"
	SyncGroupRequest              = "SyncGroupRequest"
	SyncAdminGroupRequest         = "SyncAdminGroupRequest"
	SyncFriendRequest             = "SyncFriendRequest"
	SyncSelfFriendRequest         = "SyncSelfFriendRequest"
	SyncFriendInfo                = "SyncFriendInfo"
	SyncAddBlackList              = "SyncAddBlackList"
	SyncDeleteBlackList           = "SyncDeleteBlackList"
	SyncUserInfo                  = "SyncUserInfo"
	SyncWelcomeMessageFromChannel = "SyncWelcomeMessageFromChannel"
)

const (
	DefaultPageNumber    = 1
	DefaultPageSize      = 10
	APITimeout           = 60
	MaxAttempts          = 10
	EmptyPassword        = ""
	BigVersion           = "v1"
	UpdateVersion        = ".1.0"
	SdkVersion           = "Share-Wallet-SDK-"
	UserAccountCreate    = 1
	UserAccountRecovered = 2
)
const (
	BTCCoin   = 1
	ETHCoin   = 2
	USDTERC20 = 3
	TRX       = 4
	USDTTRC20 = 5
	// Add other coins
)

const (
	//wallet
	GetSupportTokenAddressURL = "/api/v1/wallet/get_support_token_addresses"
	GetFundLogListURL         = "/api/v1/wallet/get_fundlog_list"
	GetAccountBalanceURL      = "/api/v1/wallet/get_balance"
	GetRecentRecordsURL       = "/api/v1/wallet/funds_log"
	GetTransactionListURL     = "/api/v1/wallet/get_transaction_list"
	GetTransactionURL         = "/api/v1/wallet/get_transaction"
	//eth
	CreateEthAccountURL              = "/api/v1/eth/create_account"
	GetETHBalanceURL                 = "/api/v1/eth/get_balance"
	ETHTransferAccountURL            = "/api/v1/eth/transfer"
	CheckBalanceAndNonceURL          = "/api/v1/eth/check_balance_nonce"
	GetGasPriceURL                   = "/api/v1/eth/get_gasprice"
	GetETHTransactionConfirmationURL = "/api/v1/eth/get_confirmation"
	//Tron
	CreateTronTransactionURL          = "/api/v1/tron/create_tx"
	TronTransferAccountURL            = "/api/v1/tron/transfer"
	GetTronTransactionConfirmationURL = "/api/v1/tron/get_confirmation"

	//sync urls
	CreateAccountInformationURL = "/api/v1/wallet/account"
	CreateFundLogURL            = "/cms/v1/admin/fund-log"
	UpdateFundLog               = "/cms/v1/admin/update-fund-log"
	UpdateLoginInformationURL   = "/api/v1/wallet/update-login-info"
	GetCoinStatusesURL          = "/api/v1/wallet/coins"
	GetCoinRatioURL             = "/api/v1/wallet/coin_ratio"
	GetUserWalletURL            = "/api/v1/wallet/get_user_wallet"

	GetBtc     = "https://api.coinbase.com/v2/exchange-rates?currency=BTC"
	GetEth     = "https://api.coinbase.com/v2/exchange-rates?currency=ETH"
	GetTrx     = "https://api.coinbase.com/v2/exchange-rates?currency=TRX"
	GetUsdt    = "https://api.coinbase.com/v2/exchange-rates?currency=USDT"
	GetTrxTemp = "https://api.coingecko.com/api/v3/simple/price?ids=tron&vs_currencies=usd,cny,eur"
	//USDT-ERC20
	GetUSDTERC20BalanceURL = "/api/v1/usdterc20/get_balance"
)

const (
	TxStatusPending        = 0
	TxStatusSuccess        = 1
	TxStatusFailed         = 2
	TxStatusExcludePending = 3

	//0 fail,1 success,2 pending
	FundlogFailed  = 0
	FundlogSuccess = 1
	FundlogPending = 2

	TransactionTypeAll     = 0
	TransactionTypeSend    = 1
	TransactionTypeReceive = 2

	TransactionTypeSendString    = "transfer"
	TransactionTypeReceiveString = "received"

	ConfirmationStatusWaitting  = 0
	ConfirmationStatusPending   = 1
	ConfirmationStatusCompleted = 2

	OrderLockDuration = 60
)

const (
	ETHGasLimit       = 21000
	USDTERC20GasLimit = 70000
)

const (
	SuperAdmin = "spswadmin"
)
