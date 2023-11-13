package constant

import "errors"

// key = errCode, string = errMsg
type ErrInfo struct {
	ErrCode int32
	ErrMsg  string
}

var (
	OK              = ErrInfo{0, ""}
	ErrServer       = ErrInfo{500, "server error"}
	ErrUnauthorized = ErrInfo{401, "Unauthorized"}

	//	ErrMysql             = ErrInfo{100, ""}
	//	ErrMongo             = ErrInfo{110, ""}
	//	ErrRedis             = ErrInfo{120, ""}
	ErrParseToken = ErrInfo{700, ParseTokenMsg.Error()}
	//	ErrCreateToken       = ErrInfo{201, "Create token failed"}
	//	ErrAppServerKey      = ErrInfo{300, "key error"}
	ErrTencentCredential           = ErrInfo{400, ThirdPartyMsg.Error()}
	ErrTokenExpired                = ErrInfo{701, TokenExpiredMsg.Error()}
	ErrTokenInvalid                = ErrInfo{702, TokenInvalidMsg.Error()}
	ErrTokenMalformed              = ErrInfo{703, TokenMalformedMsg.Error()}
	ErrTokenNotValidYet            = ErrInfo{704, TokenNotValidYetMsg.Error()}
	ErrTokenUnknown                = ErrInfo{705, TokenUnknownMsg.Error()}
	ErrTokenKicked                 = ErrInfo{706, TokenUserKickedMsg.Error()}
	ErrorUserLoginNetDisconnection = ErrInfo{707, "Network disconnection"}

	ErrAccess                 = ErrInfo{ErrCode: 801, ErrMsg: AccessMsg.Error()}
	ErrDB                     = ErrInfo{ErrCode: 802, ErrMsg: DBMsg.Error()}
	ErrArgs                   = ErrInfo{ErrCode: 803, ErrMsg: ArgsMsg.Error()}
	ErrStatus                 = ErrInfo{ErrCode: 804, ErrMsg: StatusMsg.Error()}
	ErrCallback               = ErrInfo{ErrCode: 809, ErrMsg: CallBackMsg.Error()}
	ErrSendLimit              = ErrInfo{ErrCode: 810, ErrMsg: "send msg limit, to many request, try again later"}
	ErrMessageHasReadDisable  = ErrInfo{ErrCode: 811, ErrMsg: "message has read disable"}
	ErrInternal               = ErrInfo{ErrCode: 812, ErrMsg: "internal error"}
	ErrRPC                    = ErrInfo{ErrCode: 813, ErrMsg: "rpc failed"}
	ErrEmptyResponse          = ErrInfo{ErrCode: 814, ErrMsg: "response is empty"}
	LimitExceeded             = ErrInfo{ErrCode: 815, ErrMsg: "Limit Exceeded"}
	ErrGetVersion             = ErrInfo{ErrCode: 816, ErrMsg: "get version failed"}
	ErrAddFriendStoped        = ErrInfo{ErrCode: 819, ErrMsg: AccessMsg.Error()}
	ErrUserPhoneAlreadyExsist = ErrInfo{ErrCode: 820, ErrMsg: PhoneAlreadyExsist.Error()}
	ErrUserIDAlreadyExsist    = ErrInfo{ErrCode: 821, ErrMsg: UserIDAlreadyExsist.Error()}
	AddFriendNotSuperUserErr  = ErrInfo{ErrCode: 822, ErrMsg: AddFriendNotSuper.Error()}
	ErrWrongPassword          = ErrInfo{ErrCode: 826, ErrMsg: "Wrong password"}

	ErrUserBanned = ErrInfo{ErrCode: 900, ErrMsg: "The operation was rejected!"}

	ErrInviteCode            = ErrInfo{ErrCode: 901, ErrMsg: "Invite code error!"}
	ErrInviteCodeInexistence = ErrInfo{ErrCode: 902, ErrMsg: "Invite code existence!"}

	ErrCaptchaError     = ErrInfo{ErrCode: 903, ErrMsg: "Captcha incorrect!"}
	ErrChannelCodeError = ErrInfo{ErrCode: 904, ErrMsg: "Invitation code incorrect"}

	ErrChannelCodeInexistence  = ErrInfo{ErrCode: 905, ErrMsg: "Channel code inexistence!"}
	ErrChannelCodeIsNull       = ErrInfo{ErrCode: 906, ErrMsg: "Please enter invitation code"}
	OnlyOneOfficialChannelCode = ErrInfo{ErrCode: 907, ErrMsg: "Only one official channel code"}
	ErrChannelCodeIsDelete     = ErrInfo{ErrCode: 908, ErrMsg: "The channel code has been deleted!"}

	ErrInviteCodeLimit = ErrInfo{ErrCode: 1001, ErrMsg: "Invite code limit is failed!"}

	ErrInviteCodeSwitch = ErrInfo{ErrCode: 1101, ErrMsg: "Invite code switch is failed!"}

	ErrInviteCodeMultiDelete = ErrInfo{ErrCode: 1201, ErrMsg: "Invite code Multi delete is failed!"}

	ErrChannelCodeLimit = ErrInfo{ErrCode: 1301, ErrMsg: "Channel code limit is failed!"}

	ErrChannelCodeSwitch = ErrInfo{ErrCode: 1401, ErrMsg: "Channel code switch is failed!"}

	ErrChannelCodeMultiDelete = ErrInfo{ErrCode: 1501, ErrMsg: "Channel code Multi delete is failed!"}

	ErrAddInviteCodeIsExist      = ErrInfo{ErrCode: 1601, ErrMsg: "Add invite code is exist!"}
	ErrAddInviteCodeUserIsExist  = ErrInfo{ErrCode: 1602, ErrMsg: "user is exist!"}
	ErrAddInviteCodeUserNotExist = ErrInfo{ErrCode: 1602, ErrMsg: "user not exist!"}
	ErrAddInviteCodeUserHasExist = ErrInfo{ErrCode: 1602, ErrMsg: "user already has invite code!"}
	ErrAddInviteCodeIsDelete     = ErrInfo{ErrCode: 1603, ErrMsg: "The share code has been deleted!"}

	ErrEditInviteCodeIsNotExist = ErrInfo{ErrCode: 1701, ErrMsg: "Edit invite code is not exist!"}
	ErrEditInviteCodeFailed     = ErrInfo{ErrCode: 1702, ErrMsg: "Edit invite code failed!"}

	ErrRevokeMessageTypeError = ErrInfo{ErrCode: 1801, ErrMsg: "Message type error!"}
	ErrNotASuperUser          = ErrInfo{ErrCode: 1805, ErrMsg: "You are not a super user"}

	ErrUserExistArg              = ErrInfo{ErrCode: 1901, ErrMsg: "Args Error!"}
	ErrUserExistUserIdExist      = ErrInfo{ErrCode: 1902, ErrMsg: "User Id Exist!"}
	ErrUserExistPhoneNumberExist = ErrInfo{ErrCode: 1903, ErrMsg: "Phone Number Exist!"}
	ErrChannelCodeNotExist       = ErrInfo{ErrCode: 1904, ErrMsg: "Channel Code Not Exist!"}
	ErrChannelCodeInvalid        = ErrInfo{ErrCode: 1905, ErrMsg: "Channel Code Invalid!"}
	ErrInviteCodeInvalid         = ErrInfo{ErrCode: 1906, ErrMsg: "Invite Code Invalid!"}
	ErrUserNotExist              = ErrInfo{ErrCode: 1907, ErrMsg: "User doesn't exists"}

	ErrNotAllowGuestLogin  = ErrInfo{ErrCode: 2001, ErrMsg: "Not allow guest login!"}
	ErrRegisterByUuidLimit = ErrInfo{ErrCode: 2002, ErrMsg: "Register by uuid limit!"}

	// ErrRoleNameAlreadyExsist
	ErrRoleNameAlreadyExist = ErrInfo{ErrCode: 823, ErrMsg: roleAlreadyExist}
	ErrRoleNameDoesntExist  = ErrInfo{ErrCode: 824, ErrMsg: "role name doesnt exist!"}
	ErrRoleIsInUse          = ErrInfo{ErrCode: 825, ErrMsg: "The role has been applied"}

	//ETH errors
	ErrEthInstance               = ErrInfo{ErrCode: 3000, ErrMsg: "Failed in dialing ETH node"}
	ErrEthBalanceFailed          = ErrInfo{ErrCode: 3001, ErrMsg: "Failed in getting eth balance"}
	ErrEthBalanceNil             = ErrInfo{ErrCode: 3002, ErrMsg: "Your ETH balance is nil"}
	ErrEthBalanceZero            = ErrInfo{ErrCode: 3003, ErrMsg: "Your ETH balance is zero"}
	ErrEthBalanceNotEnough       = ErrInfo{ErrCode: 3004, ErrMsg: "You ETH balance is not enough to transact"}
	ErrEthBalanceLessThanFee     = ErrInfo{ErrCode: 3005, ErrMsg: "You ETH balance is less then transaction fee"}
	ErrEthGettingNonceFailed     = ErrInfo{ErrCode: 3006, ErrMsg: "Failed in getting nonce"}
	ErrUSDTERC20BalanceFailed    = ErrInfo{ErrCode: 3007, ErrMsg: "Failed in getting USDT-ERC20 balance"}
	ErrUSDTERC20BalanceNil       = ErrInfo{ErrCode: 3008, ErrMsg: "Your USDT-ERC20 balance is nil"}
	ErrUSDTERC20BalanceZero      = ErrInfo{ErrCode: 3009, ErrMsg: "Your USDT-ERC20 balance is zero"}
	ErrUSDTERC20BalanceNotEnough = ErrInfo{ErrCode: 3010, ErrMsg: "You USDT-ERC20 balance is not enough to transact"}

	//Transfer
	ErrTransactAmountZero        = ErrInfo{ErrCode: 3011, ErrMsg: "The transaction amount cannot be zero"}
	ErrTransactMerchantIncorrect = ErrInfo{ErrCode: 3012, ErrMsg: "Your seed phrase has been bound by another user, you cannot transfer"}
	ErrTransactFailed            = ErrInfo{ErrCode: 3013, ErrMsg: "The transaction is failed"}

	//Wallet
	ErrSendingAddressIncorrect  = ErrInfo{ErrCode: 3014, ErrMsg: "The sender address is not correct"}
	ErrReceiverAddressIncorrect = ErrInfo{ErrCode: 3015, ErrMsg: "The receiver address is not correct"}
	ErrAddressAlreadyExists     = ErrInfo{ErrCode: 3016, ErrMsg: "This address already exists"}
	ErrNetworkError             = ErrInfo{ErrCode: 3017, ErrMsg: "Network error"}

	ErrNotYourWallet   = ErrInfo{ErrCode: 3018, ErrMsg: "Your seed phrase has been bound by another user"}
	ErrYouBoundAlready = ErrInfo{ErrCode: 3019, ErrMsg: "You have a seed phrase already, you need to unbind it before you bind another one"}

	//Tron
	ErrCreateTronTransaction = ErrInfo{ErrCode: 3100, ErrMsg: "create transaction failed"}
	ErrTronInstance          = ErrInfo{ErrCode: 3101, ErrMsg: "Failed in dialing Tron node"}

	ErrIncorrectAddress = ErrInfo{ErrCode: 3018, ErrMsg: "The address is incorrect, please input again"}
)

// Constants to store error messages
const (
	roleAlreadyExist = "role name is already used by another role"
)

var (
	ParseTokenMsg       = errors.New("parse token failed")
	TokenExpiredMsg     = errors.New("token is timed out, please log in again")
	TokenInvalidMsg     = errors.New("token has been invalidated")
	TokenNotValidYetMsg = errors.New("token not active yet")
	TokenMalformedMsg   = errors.New("that's not even a token")
	TokenUnknownMsg     = errors.New("couldn't handle this token")
	TokenUserKickedMsg  = errors.New("user has been kicked")
	AccessMsg           = errors.New("no permission")
	StatusMsg           = errors.New("status is abnormal")
	DBMsg               = errors.New("db failed")
	ArgsMsg             = errors.New("args failed")
	CallBackMsg         = errors.New("callback failed")
	AddFriendNotSuper   = errors.New("sorry, only super user can add friends")

	ThirdPartyMsg = errors.New("third party error")

	PhoneAlreadyExsist  = errors.New("phone number is already used by a user")
	UserIDAlreadyExsist = errors.New("user id is already used by a user")
)

const (
	NoError              = 0
	FormattingError      = 10001
	HasRegistered        = 10002
	NotRegistered        = 10003
	PasswordErr          = 10004
	GetIMTokenErr        = 10005
	RepeatSendCode       = 10006
	MailSendCodeErr      = 10007
	SmsSendCodeErr       = 10008
	CodeInvalidOrExpired = 10009
	RegisterFailed       = 10010
	ResetPasswordFailed  = 10011
	NotAllowRegisterType = 10012
	DatabaseError        = 10002
	ServerError          = 10004
	HttpError            = 10005
	IoError              = 10006
	IntentionalError     = 10007
)

func (e ErrInfo) Error() string {
	return e.ErrMsg
}

func (e *ErrInfo) Code() int32 {
	return e.ErrCode
}

func NewErrInfo(code int32, msg string) *ErrInfo {
	return &ErrInfo{
		code,
		msg,
	}
}
