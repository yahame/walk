// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"math"
	"strconv"
	"syscall"
	"unsafe"
)

import . "github.com/lxn/go-winapi"

const numberEditWindowClass = `\o/ Walk_NumberEdit_Class \o/`

func init() {
	MustRegisterWindowClass(numberEditWindowClass)
}

type NumberEdit struct {
	WidgetBase
	edit                  *LineEdit
	hWndUpDown            HWND
	bindingMember         string
	increment             float64
	oldValue              float64
	valueChangedPublisher EventPublisher
}

func NewNumberEdit(parent Container) (*NumberEdit, error) {
	ne := &NumberEdit{increment: 1}

	if err := InitChildWidget(
		ne,
		parent,
		numberEditWindowClass,
		WS_VISIBLE,
		WS_EX_CONTROLPARENT); err != nil {
		return nil, err
	}

	var succeeded bool
	defer func() {
		if !succeeded {
			ne.Dispose()
		}
	}()

	var err error
	ne.edit, err = newLineEdit(ne)
	if err != nil {
		return nil, err
	}
	if err = ne.edit.setAndClearStyleBits(ES_RIGHT, ES_LEFT|ES_CENTER); err != nil {
		return nil, err
	}
	nv := NewNumberValidator()
	ne.edit.SetValidator(nv)
	nv.SetDecimals(2)
	nv.SetRange(0, 100)

	ne.hWndUpDown = CreateWindowEx(
		0, syscall.StringToUTF16Ptr("msctls_updown32"), nil,
		WS_CHILD|WS_VISIBLE|UDS_ALIGNRIGHT|UDS_ARROWKEYS|UDS_HOTTRACK,
		0, 0, 16, 20, ne.hWnd, 0, 0, nil)
	if ne.hWndUpDown == 0 {
		return nil, lastError("CreateWindowEx")
	}

	SendMessage(ne.hWndUpDown, UDM_SETBUDDY, uintptr(ne.edit.hWnd), 0)

	if err = ne.SetValue(0); err != nil {
		return nil, err
	}

	succeeded = true

	return ne, nil
}

func (ne *NumberEdit) Enabled() bool {
	return ne.WidgetBase.Enabled()
}

func (ne *NumberEdit) SetEnabled(value bool) {
	ne.edit.SetEnabled(value)
	ne.WidgetBase.SetEnabled(value)
}

func (ne *NumberEdit) Font() *Font {
	var f *Font
	if ne.edit != nil {
		f = ne.font
	}

	if f != nil {
		return f
	} else if ne.parent != nil {
		return ne.parent.Font()
	}

	return defaultFont
}

func (ne *NumberEdit) SetFont(value *Font) {
	ne.edit.SetFont(value)
}

func (*NumberEdit) LayoutFlags() LayoutFlags {
	return ShrinkableHorz | GrowableHorz
}

func (ne *NumberEdit) MinSizeHint() Size {
	return ne.dialogBaseUnitsToPixels(Size{20, 12})
}

func (ne *NumberEdit) SizeHint() Size {
	s := ne.dialogBaseUnitsToPixels(Size{50, 12})
	return Size{s.Width, maxi(s.Height, 22)}
}

func (ne *NumberEdit) BindingMember() string {
	return ne.bindingMember
}

func (ne *NumberEdit) SetBindingMember(member string) error {
	if err := validateBindingMemberSyntax(member); err != nil {
		return err
	}

	ne.bindingMember = member

	return nil
}

func (ne *NumberEdit) BindingValue() interface{} {
	return ne.Value()
}

func (ne *NumberEdit) SetBindingValue(value interface{}) error {
	return ne.SetValue(value.(float64))
}

func (ne *NumberEdit) BindingValueChanged() *Event {
	return ne.ValueChanged()
}

func (ne *NumberEdit) Decimals() int {
	return ne.edit.Validator().(*NumberValidator).Decimals()
}

func (ne *NumberEdit) SetDecimals(value int) error {
	if err := ne.edit.Validator().(*NumberValidator).SetDecimals(value); err != nil {
		return err
	}

	return ne.SetValue(ne.oldValue)
}

func (ne *NumberEdit) Increment() float64 {
	return ne.increment
}

func (ne *NumberEdit) SetIncrement(value float64) error {
	ne.increment = value

	return nil
}

func (ne *NumberEdit) MinValue() float64 {
	return ne.edit.Validator().(*NumberValidator).MinValue()
}

func (ne *NumberEdit) MaxValue() float64 {
	return ne.edit.Validator().(*NumberValidator).MaxValue()
}

func (ne *NumberEdit) SetRange(min, max float64) error {
	return ne.edit.Validator().(*NumberValidator).SetRange(min, max)
}

func (ne *NumberEdit) Value() float64 {
	val, _ := parseFloat(ne.edit.Text())
	return val
}

func (ne *NumberEdit) SetValue(value float64) (err error) {
	var text string
	prec := ne.Decimals()

	if prec == 0 {
		text = strconv.Itoa(int(value))
	} else {
		text, err = formatFloat(value, prec)
		if err != nil {
			return
		}
	}

	if err = ne.edit.SetText(text); err != nil {
		return
	}

	return
}

func (ne *NumberEdit) ValueChanged() *Event {
	return ne.valueChangedPublisher.Event()
}

func (ne *NumberEdit) SetFocus() error {
	if SetFocus(ne.edit.hWnd) == 0 {
		return lastError("SetFocus")
	}

	return nil
}

func (ne *NumberEdit) TextSelection() (start, end int) {
	return ne.edit.TextSelection()
}

func (ne *NumberEdit) SetTextSelection(start, end int) {
	ne.edit.SetTextSelection(start, end)
}

func (ne *NumberEdit) WndProc(hwnd HWND, msg uint32, wParam, lParam uintptr) uintptr {
	if ne.hWndUpDown != 0 {
		switch msg {
		case WM_COMMAND:
			switch HIWORD(uint32(wParam)) {
			case EN_CHANGE:
				value := ne.Value()
				if math.Abs(value-ne.oldValue) < math.SmallestNonzeroFloat64 {
					break
				}

				ne.oldValue = value

				ne.valueChangedPublisher.Publish()
			}

		case WM_NOTIFY:
			switch ((*NMHDR)(unsafe.Pointer(lParam))).Code {
			case UDN_DELTAPOS:
				nmud := (*NMUPDOWN)(unsafe.Pointer(lParam))
				val := ne.Value()
				val -= float64(nmud.IDelta) * ne.increment
				ne.SetValue(val)
			}

		case WM_SIZE, WM_SIZING:
			cb := ne.ClientBounds()
			if err := ne.edit.SetBounds(cb); err != nil {
				break
			}
			SendMessage(ne.hWndUpDown, UDM_SETBUDDY, uintptr(ne.edit.hWnd), 0)
		}
	}

	return ne.WidgetBase.WndProc(hwnd, msg, wParam, lParam)
}
