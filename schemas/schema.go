package schemas

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v5"
	"gopkg.in/yaml.v3"

	_ "embed"

	_ "github.com/santhosh-tekuri/jsonschema/v5/httploader"
)

//go:embed schema.json
var schemaString string

var Schema *jsonschema.Schema

func init() {
	var err error

	compiler := jsonschema.NewCompiler()

	if err = compiler.AddResource("https://pomdtr.github.io/sunbeam/schemas/page.json", strings.NewReader(schemaString)); err != nil {
		panic(err)
	}
	Schema, err = compiler.Compile("https://pomdtr.github.io/sunbeam/schemas/page.json")
	if err != nil {
		panic(err)
	}
}

func Validate(bytes []byte) error {
	var v any
	if err := json.Unmarshal(bytes, &v); err != nil {
		return fmt.Errorf("invalid json: %w", err)
	}

	return Schema.Validate(v)
}

type PageType int

const (
	Unknown PageType = iota
	DetailPage
	ListPage
)

func (p *PageType) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	switch s {
	case "detail":
		*p = DetailPage
	case "list":
		*p = ListPage
	default:
		return fmt.Errorf("unknown page type: %s", s)
	}

	return nil
}

func (p PageType) MarshalJSON() ([]byte, error) {
	switch p {
	case DetailPage:
		return json.Marshal("detail")
	case ListPage:
		return json.Marshal("list")
	default:
		return nil, fmt.Errorf("unknown page type: %d", p)
	}
}

func (p *PageType) UnmarshalYAML(node *yaml.Node) error {
	var s string
	if err := node.Decode(&s); err != nil {
		return err
	}

	switch s {
	case "detail":
		*p = DetailPage
	case "list":
		*p = ListPage
	default:
		return fmt.Errorf("unknown page type: %s", s)
	}

	return nil
}

type Page struct {
	Type    PageType `json:"type"`
	Title   string   `json:"title,omitempty"`
	Actions []Action `json:"actions,omitempty"`

	*Detail `yaml:",inline"`
	*List   `yaml:",inline"`
}

type Detail struct {
	Text     string `json:"text"`
	Command  string `json:"command"`
	Language string `json:"language"`
}

type List struct {
	ShowDetail    bool       `json:"showDetail,omitempty" yaml:"showDetail"`
	GenerateItems bool       `json:"generateItems,omitempty" yaml:"generateItems"`
	EmptyText     string     `json:"emptyText,omitempty" yaml:"emptyText"`
	Items         []ListItem `json:"items"`
}

type ListItem struct {
	Id          string   `json:"id,omitempty"`
	Title       string   `json:"title"`
	Subtitle    string   `json:"subtitle,omitempty"`
	Detail      *Detail  `json:"detail,omitempty"`
	Accessories []string `json:"accessories,omitempty"`
	Actions     []Action `json:"actions,omitempty"`
}

type FormInputType int

const (
	UnknownFormInput FormInputType = iota
	TextField
	TextArea
	DropDown
)

func (input *FormInputType) UnmarshalJSON(bytes []byte) error {
	var s string
	if err := json.Unmarshal(bytes, &s); err != nil {
		return fmt.Errorf("unable to unmarshal form input type: %w", err)
	}

	switch s {
	case "textfield":
		*input = TextField
	case "textarea":
		*input = TextArea
	case "dropdown":
		*input = DropDown
	default:
		return fmt.Errorf("unknown form input type: %s", s)
	}

	return nil
}

func (input *FormInputType) UnmarshalYAML(value *yaml.Node) error {
	var s string
	if err := value.Decode(&s); err != nil {
		return fmt.Errorf("unable to unmarshal form input type: %w", err)
	}

	switch s {
	case "textfield":
		*input = TextField
	case "textarea":
		*input = TextArea
	case "dropdown":
		*input = DropDown
	default:
		return fmt.Errorf("unknown form input type: %s", s)
	}

	return nil
}

func (input FormInputType) MarshalJSON() ([]byte, error) {
	switch input {
	case TextField:
		return json.Marshal("textfield")
	case TextArea:
		return json.Marshal("textarea")
	case DropDown:
		return json.Marshal("dropdown")
	default:
		return nil, fmt.Errorf("unknown form input type: %d", input)
	}
}

type FormInput struct {
	Name        string        `json:"name"`
	Type        FormInputType `json:"type"`
	Title       string        `json:"title"`
	Placeholder string        `json:"placeholder,omitempty"`
	Default     string        `json:"default,omitempty"`

	// Only for dropdown
	Choices []string `json:"choices,omitempty"`
}

type ActionType int

const (
	UnknownAction ActionType = iota
	CopyAction
	OpenAction
	ReadAction
	EditAction
	RunAction
	ReloadAction
)

func (a *ActionType) UnmarshalJSON(bytes []byte) error {
	var rawType string
	if err := json.Unmarshal(bytes, &rawType); err != nil {
		return err
	}

	switch rawType {
	case "copy":
		*a = CopyAction
	case "open":
		*a = OpenAction
	case "read":
		*a = ReadAction
	case "run":
		*a = RunAction
	case "edit":
		*a = EditAction
	case "reload":
		*a = ReloadAction
	default:
		return fmt.Errorf("unknown action type: %s", rawType)
	}

	return nil
}

func (a ActionType) MarshalJSON() ([]byte, error) {
	switch a {
	case CopyAction:
		return json.Marshal("copy")
	case OpenAction:
		return json.Marshal("open")
	case ReadAction:
		return json.Marshal("read")
	case EditAction:
		return json.Marshal("edit")
	case RunAction:
		return json.Marshal("run")
	case ReloadAction:
		return json.Marshal("reload")
	default:
		return nil, fmt.Errorf("unknown action type: %d", a)
	}
}

func (a *ActionType) UnmarshalYAML(node *yaml.Node) error {
	var rawType string
	if err := node.Decode(&rawType); err != nil {
		return err
	}

	switch rawType {
	case "copy":
		*a = CopyAction
	case "open":
		*a = OpenAction
	case "read":
		*a = ReadAction
	case "run":
		*a = RunAction
	case "edit":
		*a = EditAction
	case "reload":
		*a = ReloadAction
	default:
		return fmt.Errorf("unknown action type: %s", rawType)
	}

	return nil
}

type OnSuccessType int

const (
	ExitOnSuccess OnSuccessType = iota
	PushOnSuccess
	ReloadOnSuccess
)

func (o *OnSuccessType) UnmarshalJSON(bytes []byte) error {
	var rawType string
	if err := json.Unmarshal(bytes, &rawType); err != nil {
		return err
	}

	switch rawType {
	case "push":
		*o = PushOnSuccess
	case "reload":
		*o = ReloadOnSuccess
	default:
		return fmt.Errorf("unknown on success type: %s", rawType)
	}

	return nil
}

func (o OnSuccessType) MarshalJSON() ([]byte, error) {
	switch o {
	case PushOnSuccess:
		return json.Marshal("push")
	case ReloadOnSuccess:
		return json.Marshal("reload")
	default:
		return nil, fmt.Errorf("unknown on success type: %d", o)
	}
}

func (o *OnSuccessType) UnmarshalYAML(node *yaml.Node) error {
	var rawType string
	if err := node.Decode(&rawType); err != nil {
		return err
	}

	switch rawType {
	case "push":
		*o = PushOnSuccess
	case "reload":
		*o = ReloadOnSuccess
	default:
		return fmt.Errorf("Unknown onSuccess type: %s", rawType)
	}

	return nil
}

type Action struct {
	RawTitle string     `json:"title,omitempty" yaml:"title"`
	Shortcut string     `json:"shortcut,omitempty"`
	Type     ActionType `json:"type"`

	// open
	Target string `json:"target,omitempty"`

	// copy
	Text string `json:"text,omitempty"`

	// read
	Page string `json:"page,omitempty"`

	// edit
	Path string `json:"path,omitempty"`

	// run
	Command   string        `json:"command,omitempty"`
	Inputs    []FormInput   `json:"inputs,omitempty"`
	OnSuccess OnSuccessType `json:"onSuccess,omitempty" yaml:"onSuccess"`
}

func (a Action) Title() string {
	if a.RawTitle != "" {
		return a.RawTitle
	}
	switch a.Type {
	case CopyAction:
		return "Copy"
	case OpenAction:
		return "Open"
	case ReadAction:
		return "Read"
	case RunAction:
		return "Run"
	case ReloadAction:
		return "Reload"
	default:
		return "Unknown"
	}
}
