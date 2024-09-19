package tui

import "github.com/pomdtr/sunbeam/pkg/sunbeam"

func NewErrorPage(err error, additionalActions ...sunbeam.ActionItem) *Detail {
	var actions []sunbeam.ActionItem
	actions = append(actions, sunbeam.ActionItem{
		Title: "Copy error",
		Type:  sunbeam.ActionTypeCopy,
		Copy:  sunbeam.CopyAction{Text: err.Error(), Exit: true},
	})
	actions = append(actions, additionalActions...)

	detail := NewDetail("Error", err.Error(), actions...)

	return detail
}
