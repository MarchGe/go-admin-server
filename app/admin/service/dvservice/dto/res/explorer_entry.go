package res

type EntryType string

const (
	EntryTypeDefault     EntryType = "" // 其他未具体处理的类型，统一归到这个类型下
	EntryTypeDir         EntryType = "d"
	EntryTypeLink        EntryType = "l"
	EntryTypeCharDevice  EntryType = "c"
	EntryTypeBlockDevice EntryType = "b"
	EntryTypeSocket      EntryType = "s"
	EntryTypeNamedPipe   EntryType = "p"
)

type ExplorerEntry struct {
	Name string    `json:"name"`
	Type EntryType `json:"type"`
}
