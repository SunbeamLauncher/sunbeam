package sunbeam

import (
	"encoding/json"
	"fmt"
)

type Page struct {
	Type  PageType `json:"type,omitempty"`
	Title string   `json:"title,omitempty"`

	List   List   `json:"list,omitempty"`
	Detail Detail `json:"detail,omitempty"`
	Form   Form   `json:"form,omitempty"`
}

func (p *Page) UnmarshalJSON(data []byte) error {
	var aux struct {
		Type  PageType `json:"type,omitempty"`
		Title string   `json:"title,omitempty"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	p.Type = aux.Type
	p.Title = aux.Title
	switch aux.Type {
	case PageTypeList:
		return json.Unmarshal(data, &p.List)
	case PageTypeDetail:
		return json.Unmarshal(data, &p.Detail)
	case PageTypeForm:
		return json.Unmarshal(data, &p.Form)
	default:
		return fmt.Errorf("unsupported page type: %s", aux.Type)
	}
}

type PageType string

const (
	PageTypeList   PageType = "list"
	PageTypeDetail PageType = "detail"
	PageTypeForm   PageType = "form"
)

type List struct {
	Items      []ListItem   `json:"items,omitempty"`
	Dynamic    bool         `json:"dynamic,omitempty"`
	EmptyText  string       `json:"emptyText,omitempty"`
	ShowDetail bool         `json:"showDetail,omitempty"`
	Actions    []ActionItem `json:"actions,omitempty"`
}

type ListItem struct {
	Id          string         `json:"id,omitempty"`
	Title       string         `json:"title"`
	Subtitle    string         `json:"subtitle,omitempty"`
	Detail      ListItemDetail `json:"detail,omitempty"`
	Accessories []string       `json:"accessories,omitempty"`
	Actions     []ActionItem   `json:"actions,omitempty"`
}

type ListItemDetail struct {
	Markdown string `json:"markdown,omitempty"`
	Text     string `json:"text,omitempty"`
}

type Detail struct {
	Actions  []ActionItem `json:"actions,omitempty"`
	Markdown string       `json:"markdown,omitempty"`
	Text     string       `json:"text,omitempty"`
}

type Form struct {
	Inputs []Input `json:"inputs,omitempty"`
}

type InputType string

const (
	InputTextField InputType = "textfield"
	InputTextArea  InputType = "textarea"
	InputCheckbox  InputType = "checkbox"
	InputNumber    InputType = "number"
	InputPassword  InputType = "password"
)

type Input struct {
	Type     InputType `json:"type"`
	Name     string    `json:"name"`
	Title    string    `json:"title"`
	Optional bool      `json:"optional,omitempty"`
	Default  any       `json:"default,omitempty"`
}
