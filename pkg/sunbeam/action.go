package sunbeam

import (
	"encoding/json"
	"fmt"
)

type ActionItem struct {
	Title string     `json:"title,omitempty"`
	Key   string     `json:"key,omitempty"`
	Type  ActionType `json:"type,omitempty"`

	Copy   CopyAction   `json:"-"`
	Run    RunAction    `json:"-"`
	Push   PushAction   `json:"-"`
	Open   OpenAction   `json:"-"`
	Reload ReloadAction `json:"-"`
}

func (a *ActionItem) UnmarshalJSON(data []byte) error {
	var aux struct {
		Type  ActionType `json:"type,omitempty"`
		Title string     `json:"title,omitempty"`
		Key   string     `json:"key,omitempty"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	a.Title = aux.Title
	a.Type = aux.Type
	a.Key = aux.Key

	switch aux.Type {
	case ActionTypeCopy:
		return json.Unmarshal(data, &a.Copy)
	case ActionTypeRun:
		return json.Unmarshal(data, &a.Run)
	case ActionTypeOpen:
		return json.Unmarshal(data, &a.Open)
	case ActionTypeReload:
		return json.Unmarshal(data, &a.Reload)
	case ActionTypePush:
		return json.Unmarshal(data, &a.Push)
	case ActionTypeExit:
		return nil
	default:
		return fmt.Errorf("unsupported action type: %s", aux.Type)
	}

}

type ReloadAction struct {
	Args []string `json:"params,omitempty"`
}

type RunAction struct {
	Args []string `json:"args,omitempty"`
}

type PushAction struct {
	Args []string `json:"args,omitempty"`
}

type CopyAction struct {
	Text string `json:"text,omitempty"`
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
	ActionTypePush   ActionType = "push"
	ActionTypeReload ActionType = "reload"
)
