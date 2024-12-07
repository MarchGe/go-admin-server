package constant

// cookie name
const (
	// LoginSession 登录成功后，保存的cookie信息
	LoginSession = "LOGIN_SESSION"
)

// session key
const (
	SessionUserId = "_userId_"
	IsRootUser    = "_isRootUser_"
)

const (
	DateFormat                   = "2006-01-02T15:04:05.999Z" // 日期格式
	RequestId                    = "X-Request-Id"             // http请求头
	Swagger                      = "/swagger"                 //swagger路径
	SshEstablishTimeoutInSeconds = 30                         // ssh建立连接的超时时间
)

const (
	ServerInternalError = "Server internal error"
)

const NewLine = "\n"
