package main

import (
	"bytes"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type content struct {
	renderer  *glamour.TermRenderer
	wrapStyle lipgloss.Style
	raw       bytes.Buffer
	rendered  string
}

func newContent(wrapWidth, wordWidth int) (content, error) {
	renderer, wrapStyle, err := newRenderer(wrapWidth, wordWidth)
	if err != nil {
		return content{}, err
	}
	return content{
		renderer:  renderer,
		wrapStyle: wrapStyle,
	}, nil
}

func newRenderer(wrapWidth, wordWidth int) (*glamour.TermRenderer, lipgloss.Style, error) {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(wrapWidth),
	)
	if err != nil {
		return nil, lipgloss.NewStyle(), err
	}
	wrapStyle := lipgloss.NewStyle().Width(wordWidth)
	return renderer, wrapStyle, nil
}

func (c content) add(message string) (content, error) {
	c.raw.WriteString(message)
	var err error
	rendered, err := c.renderer.Render(c.raw.String())
	if err != nil {
		return c, err
	}
	c.rendered = c.wrapStyle.Render(rendered)
	return c, nil
}

func (c content) update(wrapWidth, wordWidth int) (content, error) {
	newC, err := newContent(wrapWidth, wordWidth)
	if err != nil {
		return newC, err
	}
	return newC.add(c.raw.String())
}
