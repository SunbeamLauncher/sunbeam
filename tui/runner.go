package tui

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pomdtr/sunbeam/app"
)

type RunContainer struct {
	width, height int
	currentView   string

	extension app.Extension
	with      map[string]any

	form   *Form
	list   *List
	detail *Detail

	Script app.Script
}

func NewRunContainer(manifest app.Extension, script app.Script, with map[string]any) *RunContainer {
	params := make(map[string]any)
	for k, v := range with {
		params[k] = v
	}

	return &RunContainer{
		extension: manifest,
		Script:    script,
		with:      params,
	}
}

func (c *RunContainer) Init() tea.Cmd {
	return c.Run()
}

type CommandOutput string

func (c RunContainer) ScriptCmd() tea.Msg {
	input := app.CommandParams{
		With: c.with,
	}
	if c.currentView == "list" {
		input.Query = c.list.Query()
	}

	command, err := c.Script.Cmd(input)
	if err != nil {
		return err
	}

	switch c.Script.Cwd {
	case "homeDir":
		command.Dir, err = os.UserHomeDir()
		if err != nil {
			return err
		}
	case "extensionDir":
		command.Dir = c.extension.Dir()
	case "currentDir":
		command.Dir, err = os.Getwd()
		if err != nil {
			return err
		}
	default:
		command.Dir = c.extension.Dir()
	}

	if c.Script.OnSuccess == "" {
		return ExecCommandMsg{Command: command}
	}

	res, err := command.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if ok := errors.As(err, &exitErr); ok {
			return NewErrorCmd(fmt.Errorf("command failed with exit code %d, error: %s", exitErr.ExitCode(), exitErr.Stderr))
		}
		return NewErrorCmd(err)
	}

	return CommandOutput(string(res))
}

func (c *RunContainer) Run() tea.Cmd {
	missing := c.Script.CheckMissingParams(c.with)

	if len(missing) > 0 {
		var err error
		items := make([]FormItem, len(missing))
		for i, param := range missing {
			items[i], err = NewFormItem(param)
			if err != nil {
				return NewErrorCmd(err)
			}
		}

		c.currentView = "form"
		c.form = NewForm(c.Script.Page.Title, items, func(values map[string]any) tea.Cmd {
			for k, v := range values {
				c.with[k] = v
			}

			return c.Run()
		})
		c.form.SetSize(c.width, c.height)
		return c.form.Init()
	}

	switch c.Script.Page.Type {
	case "list":
		c.currentView = "list"
		if c.list != nil {
			c.list.SetIsLoading(true)
			return c.ScriptCmd
		}
		c.list = NewList(c.Script.Page.Title)
		if c.Script.Page.Mode == "generator" {
			c.list.Dynamic = true
		}
		if c.Script.Page.ShowPreview {
			c.list.ShowPreview = true
		}
		c.list.SetSize(c.width, c.height)
		cmd := c.list.SetIsLoading(true)
		return tea.Batch(c.ScriptCmd, c.list.Init(), cmd)
	default:
		c.currentView = "detail"
		if c.detail != nil {
			c.detail.SetIsLoading(true)
			return c.ScriptCmd
		}

		c.detail = NewDetail(c.Script.Page.Title)
		c.detail.SetSize(c.width, c.height)
		cmd := c.detail.SetIsLoading(true)
		return tea.Batch(c.ScriptCmd, cmd, c.detail.Init())
	}
}

func (c *RunContainer) SetSize(width, height int) {
	c.width, c.height = width, height
	switch c.currentView {
	case "list":
		c.list.SetSize(width, height)
	case "detail":
		c.detail.SetSize(width, height)
	case "form":
		c.form.SetSize(width, height)
	}
}

func (c *RunContainer) Update(msg tea.Msg) (Container, tea.Cmd) {
	switch msg := msg.(type) {
	case CommandOutput:
		switch c.Script.OnSuccess {
		// Reload page is used to reload the *previous* page
		case "reload-page":
			return c, tea.Sequence(PopCmd, NewReloadPageCmd(nil))
		case "copy-to-clipboard":
			return c, NewCopyTextCmd(string(msg))
		case "open-url":
			return c, NewOpenUrlCmd(string(msg))
		case "push-page":
			switch c.currentView {
			case "detail":
				c.detail.SetContent(string(msg))
				c.detail.SetIsLoading(false)
				c.detail.SetActions(
					Action{Title: "Quit", Shortcut: "enter", Cmd: tea.Quit},
					Action{Title: "Copy Output", Shortcut: "ctrl+y", Cmd: func() tea.Msg {
						return CopyTextMsg{
							Text: string(msg),
						}
					}},
					Action{Title: "Reload", Shortcut: "ctrl+r", Cmd: NewReloadPageCmd(nil)},
				)
				return c, nil
			case "list":
				scriptItems, err := app.ParseListItems(string(msg))
				if err != nil {
					return c, NewErrorCmd(err)
				}
				listItems := make([]ListItem, len(scriptItems))

				for i, scriptItem := range scriptItems {
					scriptItem := scriptItem
					actions := make([]Action, len(scriptItem.Actions))
					for i, scriptAction := range scriptItem.Actions {
						if i == 0 {
							scriptAction.Shortcut = "enter"
						}
						if scriptAction.Extension == "" {
							scriptAction.Extension = c.extension.Name
						}
						actions[i] = NewAction(scriptAction)
					}
					if scriptItem.Id == "" {
						scriptItem.Id = strconv.Itoa(i)
					}
					listItems[i] = ListItem{
						id:       scriptItem.Id,
						Title:    scriptItem.Title,
						Subtitle: scriptItem.Subtitle,
						Preview:  scriptItem.Preview,
						PreviewCmd: func() string {
							cmd := scriptItem.PreviewCommand()
							if cmd == nil {
								return "No preview command"
							}
							out, err := cmd.Output()
							if err != nil {
								var exitErr *exec.ExitError
								if errors.As(err, &exitErr) {
									return string(exitErr.Stderr)
								}
								return err.Error()
							}

							return string(out)
						},
						Accessories: scriptItem.Accessories,
						Actions:     actions,
					}
				}

				cmd := c.list.SetItems(listItems)
				c.list.SetIsLoading(false)
				return c, cmd
			}
		}
	case ReloadPageMsg:
		for k, v := range msg.Params {
			c.with[k] = v
		}
		var cmd tea.Cmd
		if c.currentView == "list" {
			cmd = c.list.header.SetIsLoading(true)
		} else if c.currentView == "detail" {
			cmd = c.detail.header.SetIsLoading(true)
		}
		return c, tea.Batch(cmd, c.Run())
	}

	var cmd tea.Cmd
	var container Container

	switch c.currentView {
	case "list":
		container, cmd = c.list.Update(msg)
		c.list, _ = container.(*List)
	case "detail":
		container, cmd = c.detail.Update(msg)
		c.detail, _ = container.(*Detail)
	case "form":
		container, cmd = c.form.Update(msg)
		c.form, _ = container.(*Form)
	}
	return c, cmd
}

func (c *RunContainer) View() string {
	switch c.currentView {
	case "form":
		return c.form.View()
	case "list":
		return c.list.View()
	case "detail":
		return c.detail.View()
	default:
		return ""
	}
}
