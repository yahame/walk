// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type ToolButton struct {
	AssignTo           **walk.ToolButton
	Name               string
	Disabled           bool
	Hidden             bool
	Font               Font
	MinSize            Size
	MaxSize            Size
	StretchFactor      int
	Row                int
	RowSpan            int
	Column             int
	ColumnSpan         int
	ContextMenuActions []*walk.Action
	Text               string
	OnClicked          walk.EventHandler
}

func (tb ToolButton) Create(parent walk.Container) error {
	w, err := walk.NewToolButton(parent)
	if err != nil {
		return err
	}

	return InitWidget(tb, w, func() error {
		if err := w.SetText(tb.Text); err != nil {
			return err
		}

		if tb.OnClicked != nil {
			w.Clicked().Attach(tb.OnClicked)
		}

		if tb.AssignTo != nil {
			*tb.AssignTo = w
		}

		return nil
	})
}

func (tb ToolButton) WidgetInfo() (name string, disabled, hidden bool, font *Font, minSize, maxSize Size, stretchFactor, row, rowSpan, column, columnSpan int, contextMenuActions []*walk.Action) {
	return tb.Name, tb.Disabled, tb.Hidden, &tb.Font, tb.MinSize, tb.MaxSize, tb.StretchFactor, tb.Row, tb.RowSpan, tb.Column, tb.ColumnSpan, tb.ContextMenuActions
}
