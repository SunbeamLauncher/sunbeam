package sunbeam

type ActionItem struct {
	Title string     `json:"title,omitempty"`
	Key   string     `json:"key,omitempty"`
	Type  ActionType `json:"type,omitempty"`

	Copy   CopyAction   `json:"copy,omitempty"`
	Run    RunAction    `json:"run,omitempty"`
	Open   OpenAction   `json:"open,omitempty"`
	Reload ReloadAction `json:"reload,omitempty"`
}

type ReloadAction struct {
	Args []string `json:"params,omitempty"`
}

type RunAction struct {
	Args   []string `json:"args,omitempty"`
	Reload bool     `json:"reload,omitempty"`
}

type CopyAction struct {
	Text string `json:"text,omitempty"`
	Exit bool   `json:"exit,omitempty"`
}

type OpenAction struct {
	Url  string `json:"url,omitempty"`
	Path string `json:"path,omitempty"`
}

type ActionType string

const (
	ActionTypeRun    ActionType = "run"
	ActionTypeOpen   ActionType = "open"
	ActionTypeCopy   ActionType = "copy"
	ActionTypeExit   ActionType = "exit"
	ActionTypeReload ActionType = "reload"
)
