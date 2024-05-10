package constant

import "time"

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

const (
	SseTickDuration         time.Duration = 15 * time.Second    // Sse推送心跳包间隔时间
	SseTickEvent                          = "TICK"              // 心跳tick事件
	SseErrorEvent                         = "ERROR"             // 通用的出错提示
	SseCloseEvent                         = "CLOSE"             // 服务端处理完请求，希望客户端关闭连接的事件
	SseManifestEntryEvent                 = "MANIFEST_ENTRY"    // manifest清单日志文件条目增加通知
	SseHostLogEvent                       = "HOST_LOG"          // hostLog日志文件增加通知
	SseTaskExecuteEndEvent                = "TASK_EXECUTE_END"  // 任务执行完成通知
	SseTaskExecuteFailEvent               = "TASK_EXECUTE_FAIL" // 任务执行失败
)

const NewLine = "\n"
