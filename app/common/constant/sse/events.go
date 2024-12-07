package sse

import "time"

type Event string

const (
	TickDuration                 = 15 * time.Second           // Sse推送心跳包间隔时间
	Tick                   Event = "TICK"                     // 心跳tick事件
	Message                Event = "MESSAGE"                  // 通用的消息提示
	Error                  Event = "ERROR"                    // 出错提示，客户端应该主动关闭连接
	Close                  Event = "CLOSE"                    // 服务端处理完请求，希望客户端关闭连接的事件
	ManifestEntryEvent     Event = "MANIFEST_ENTRY"           // manifest清单日志文件条目增加通知
	HostLogEvent           Event = "HOST_LOG"                 // hostLog日志文件增加通知
	TaskExecuteEndEvent    Event = "TASK_EXECUTE_END"         // 任务执行完成通知
	TaskExecuteFailEvent   Event = "TASK_EXECUTE_FAIL"        // 任务执行失败
	ScriptExecuteEndEvent  Event = "SCRIPT_TASK_EXECUTE_END"  // 脚本任务执行结束通知
	ScriptExecuteFailEvent Event = "SCRIPT_TASK_EXECUTE_FAIL" // 脚本任务执行失败
)

func (e Event) String() string {
	return string(e)
}
