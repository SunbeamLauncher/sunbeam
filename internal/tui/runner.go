package tui

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"

	"github.com/acarl005/stripansi"
	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pomdtr/sunbeam/internal/utils"
	"github.com/pomdtr/sunbeam/pkg/sunbeam"
)

type Runner struct {
	embed         Page
	err           Page
	width, height int
	cancel        context.CancelFunc

	query    string
	execPath string
	args     []string
}

type ReloadMsg struct{}

func NewRunner(execPath string, args ...string) *Runner {
	detail := NewDetail("")
	detail.isLoading = true
	return &Runner{
		embed:    detail,
		execPath: execPath,
		args:     args,
	}
}

func (c *Runner) SetIsLoading(isLoading bool) tea.Cmd {
	switch page := c.embed.(type) {
	case *Detail:
		return page.SetIsLoading(isLoading)
	case *List:
		return page.SetIsLoading(isLoading)
	case *Form:
		return page.SetIsLoading(isLoading)
	}

	return nil
}

func (c *Runner) Init() tea.Cmd {
	return tea.Batch(c.embed.Init(), c.Run())
}

func (c *Runner) Focus() tea.Cmd {
	if c.embed == nil {
		return nil
	}
	return c.embed.Focus()
}

func (c *Runner) Blur() tea.Cmd {
	c.cancel()
	return nil
}

func (c *Runner) SetSize(w int, h int) {
	c.width = w
	c.height = h

	if c.err != nil {
		c.err.SetSize(w, h)
	}

	c.embed.SetSize(w, h)
}

func (me *Runner) handlePage(page sunbeam.Page) tea.Cmd {
	switch page.Type {
	case sunbeam.PageTypeList:
		list := page.List
		if embed, ok := me.embed.(*List); ok {
			embed.SetItems(list.Items...)
			embed.SetIsLoading(false)
			embed.SetEmptyText(list.EmptyText)
			embed.SetActions(list.Actions...)
			embed.SetShowDetail(list.ShowDetail)

			if list.Dynamic {
				embed.OnQueryChange = func(query string) tea.Cmd {
					me.query = query
					return me.Run()
				}

				embed.ResetSelection()
			}

			embed.SetSize(me.width, me.height)
			return embed.Init()
		}

		embed := NewList(list.Items...)
		embed.SetEmptyText(list.EmptyText)
		embed.SetActions(list.Actions...)
		embed.SetShowDetail(list.ShowDetail)
		if list.Dynamic {
			embed.OnQueryChange = func(query string) tea.Cmd {
				me.query = query
				return me.Run()
			}
		}

		me.embed = embed
		me.embed.SetSize(me.width, me.height)
		return me.embed.Init()
	case sunbeam.PageTypeDetail:
		detail := page.Detail
		if detail.Markdown != "" {
			detail := NewDetail(detail.Markdown, detail.Actions...)
			detail.Markdown = true
			detail.SetSize(me.width, me.height)
			me.embed = detail

			return me.embed.Init()
		}

		me.embed = NewDetail(detail.Text, detail.Actions...)
		me.embed.SetSize(me.width, me.height)
		return me.embed.Init()
	case sunbeam.PageTypeForm:
		form := page.Form
		me.embed = NewForm(func(flags []string) tea.Msg {
			var args []string
			args = append(args, me.args...)
			args = append(args, flags...)
			return PushPageMsg{NewRunner(me.execPath, args...)}
		}, form.Inputs...)

		me.embed.SetSize(me.width, me.height)
		return me.embed.Init()
	default:
		return func() tea.Msg {
			return fmt.Errorf("invalid page type")
		}
	}

}

func (me *Runner) handleAction(action sunbeam.ActionItem) tea.Cmd {
	return func() tea.Msg {
		switch action.Type {
		case sunbeam.ActionTypeCopy:
			if err := clipboard.WriteAll(action.Copy.Text); err != nil {
				return err
			}

			if action.Copy.Exit {
				return ExitMsg{}
			}

			return ShowNotificationMsg{"Copied!"}
		case sunbeam.ActionTypeRun:
			cmd := exec.Command(me.execPath, action.Run.Args...)
			me.SetIsLoading(true)
			output, err := cmd.Output()
			if err != nil {
				return err
			}

			var action sunbeam.ActionItem
			if err := json.Unmarshal(output, &action); err != nil {
				return err
			}

			return action
		case sunbeam.ActionTypeReload:
			if action.Reload.Args != nil {
				me.args = append(me.args, action.Reload.Args...)
			}

			return ReloadMsg{}
		case sunbeam.ActionTypePush:
			runner := NewRunner(me.execPath, action.Run.Args...)
			return PushPageCmd(runner)
		case sunbeam.ActionTypeOpen:
			if action.Open.Url != "" {
				if err := utils.Open(action.Open.Url); err != nil {
					return err
				}

				return ExitMsg{}
			} else if action.Open.Path != "" {
				if err := utils.Open(fmt.Sprintf("file://%s", action.Open.Path)); err != nil {
					return err
				}

				return ExitMsg{}
			} else {
				return fmt.Errorf("invalid target")
			}
		case sunbeam.ActionTypeExit:
			return ExitMsg{}
		default:
			return nil
		}
	}
}

func (c *Runner) Update(msg tea.Msg) (Page, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if c.err != nil {
				c.err = nil
			}
		}
	case ReloadMsg:
		return c, c.Run()
	case sunbeam.ActionItem:
		return c, c.handleAction(msg)
	case error:
		c.err = NewErrorPage(msg)
		c.err.SetSize(c.width, c.height)
		return c, c.embed.Init()
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd
	c.embed, cmd = c.embed.Update(msg)
	cmds = append(cmds, cmd)

	if c.err != nil {
		c.err, cmd = c.err.Update(msg)
		cmds = append(cmds, cmd)
	}

	return c, tea.Batch(cmds...)
}

func (c *Runner) View() string {
	if c.err != nil {
		return c.err.View()
	}

	return c.embed.View()
}

func (c *Runner) Run() tea.Cmd {
	return tea.Sequence(c.SetIsLoading(true), func() tea.Msg {
		if c.cancel != nil {
			c.cancel()
		}

		ctx, cancel := context.WithCancel(context.Background())
		c.cancel = cancel
		defer cancel()

		cmd := exec.CommandContext(ctx, c.execPath)
		cmd.Args = append(cmd.Args, c.args...)
		if c.query != "" {
			cmd.Args = append(cmd.Args, fmt.Sprintf("--query=%s", c.query))
		}

		output, err := cmd.Output()
		if err != nil {
			if errors.Is(ctx.Err(), context.Canceled) {
				return nil
			}
			var exitErr *exec.ExitError
			if errors.As(err, &exitErr) {
				return fmt.Errorf("command failed: %s", stripansi.Strip(string(exitErr.Stderr)))
			}

			return err
		}

		var page sunbeam.Page
		if err := json.Unmarshal(output, &page); err != nil {
			return err
		}

		return c.handlePage(page)
	})
}
