package main

import (
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image"
	"image/color"
	"math"
)

type (
	// T is for Theme
	T = *material.Theme
	// C is for Context
	C = layout.Context
	// D is for Dimensions
	D = layout.Dimensions
)

var (
	stdDP   = unit.Dp(10)
	hspacer = layout.Rigid(layout.Spacer{Height: stdDP}.Layout)
	wspacer = layout.Rigid(layout.Spacer{Width: stdDP}.Layout)
)

// UI _
type UI struct {
	Theme    *material.Theme
	IsDark   bool
	ChatList *ChatList
	ChatAct  *ChatActivity
	Size     image.Point
	Win      *app.Window
}

// NewUI is constructor for UI
func NewUI() *UI {
	ui := new(UI)
	ui.Theme = material.NewTheme(gofont.Collection())
	if ui.IsDark {
		ui.Theme.Palette.Bg = color.NRGBA{R: 22, G: 27, B: 34, A: 255}
		ui.Theme.Palette.Fg = color.NRGBA{R: 201, G: 209, B: 217, A: 255}
		ui.Theme.Palette.ContrastFg = color.NRGBA{R: 253, G: 253, B: 253, A: 255}
	}
	ui.ChatList = new(ChatList)
	ui.ChatAct = new(ChatActivity)
	charr := []*Chat{
		&Chat{"1", []GUIMessage{
			{"1", "hello world!"},
			{"1", "QWE42"},
		}, new(widget.Clickable)},
		&Chat{"Vasyok", []GUIMessage{}, new(widget.Clickable)},
		&Chat{"Qwertyque", []GUIMessage{
			{"Qwertyque", "lol !!!! sus amogus(((((((((((((((((((((("},
		}, new(widget.Clickable)},
		&Chat{"_", []GUIMessage{
			{"_", "Прости, что тебе спамлю, но... я шизофреник !!! !!! !!!! !!!"},
		}, new(widget.Clickable)},
	}
	ui.ChatList.Chats = charr
	ui.ChatAct.Chats = charr
	return ui
}

// Run starts layouting
func (ui *UI) Run(w *app.Window) error {
	ui.Win = w
	var ops op.Ops
	for e := range w.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			if ui.IsDark {
				paint.Fill(&ops, ui.Theme.Palette.Bg)
			}
			ui.Size = e.Size
			ui.Layout(gtx)
			ui.ChatAct.Selected = ui.ChatList.Selected
			e.Frame(gtx.Ops)
		case system.DestroyEvent:
			return e.Err
		}
	}
	return nil
}

// Layout layouts
func (ui *UI) Layout(gtx C) D {
	ui.ChatList.MaxX = ui.Size.X
	ui.ChatAct.MaxX = ui.Size.X
	var x *int
	return layout.UniformInset(stdDP).Layout(gtx, func(gtx C) D {
		return layout.Flex{
			Axis: layout.Vertical,
		}.Layout(gtx,
			layout.Rigid(
				func(gtx C) D {
					return layout.Flex{
						Axis: layout.Horizontal,
					}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							dims := ui.ChatList.Layout(gtx, ui.Theme)
							x = &dims.Size.X
							return dims
						}),
						layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
						layout.Rigid(func(gtx C) D {
							return ui.ChatAct.Layout(gtx, ui.Theme, *x)
						}),
					)
				},
			),
		)
	})
}

func fgtx(gtx C) C {
	fgx := *(&gtx) // to safely copy data (don't know why)
	fgx.Ops = new(op.Ops)
	return fgx
}

// ChatList _
type ChatList struct {
	MaxX     int
	Selected string
	Chats    []*Chat
}

// Layout _
func (cl *ChatList) Layout(gtx C, th T) D {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(
			func(gtx C) D {
				return (&layout.List{Axis: layout.Vertical}).
					Layout(gtx, len(cl.Chats), func(gtx C, ind int) D {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								return cl.Chats[ind].LayoutList(gtx, th, cl.MaxX, cl)
							}),
							layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
						)
					})
			},
		),
	)
}

// ChatActivity _
type ChatActivity struct {
	MaxX     int
	Selected string
	Chats    []*Chat
}

// Layout _
func (ca *ChatActivity) Layout(gtx C, th T, startX int) D {
	if ca.Selected == "" {
		ca.Selected = "OVERMSg"
	}
	return layout.Flex{
		Axis:    layout.Vertical,
		Spacing: layout.SpaceEnd,
	}.Layout(gtx,
		// header...
		layout.Rigid(func(gtx C) D {
			return widget.Border{Color: th.ContrastBg, Width: unit.Dp(3.5)}.Layout(gtx,
				func(gtx C) D {
					return layout.UniformInset(unit.Dp(5)).Layout(gtx,
						func(gtx C) D {
							return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
								layout.Rigid(material.Body2(th, ca.Selected).Layout),
								layout.Rigid(layout.Spacer{Width: unit.Px(
									float32(ca.MaxX)/1.07 - // hello, hardcoded number! (got it during experiments)
										float32(
											// 20 because it is sum of spacer and insets (5 + 10 + 5)
											20+startX+material.Body2(th, ca.Selected).Layout(fgtx(gtx)).Size.X,
										),
								)}.Layout),
							)
						},
					)
				},
			)
		}),
		layout.Rigid(func(gtx C) D {
			return (&layout.List{Axis: layout.Vertical}).Layout(
				gtx,
				len(ca.Chats),
				func(gtx C, ind int) D { return D{} },
			)
		}),
	)
}
func th2w(next func(C, T) D, th T) func(C) D {
	return func(gtx C) D {
		return next(gtx, th)
	}
}

// Chat is chat
type Chat struct {
	PeerName string
	Messages []GUIMessage
	Button   *widget.Clickable
}

// LayoutList layouts list small preview
func (c *Chat) LayoutList(gtx C, th T, maxX int, cl *ChatList) D {
	if c.Button.Clicked() {
		cl.Selected = c.PeerName
	}
	return c.Button.Layout(gtx, func(gtx C) D {
		return func() widget.Border {
			b := widget.Border{
				CornerRadius: unit.Dp(5),
				Color:        th.Fg,
				Width:        unit.Dp(0.5),
			}
			if c.PeerName == cl.Selected {
				b.Color = color.NRGBA{R: 255, A: 255}
				b.Width = unit.Dp(1.5)
			}

			return b
		}().Layout(gtx, func(gtx C) D {
			return layout.Inset{
				Top:    unit.Dp(5),
				Bottom: unit.Dp(5),
				Left:   unit.Dp(5),
				Right: unit.Px(float32(maxX)/4 - float32(math.Max(
					float64(material.Label(th, unit.Dp(12.5), getSmallStr(c)).Layout(fgtx(gtx)).Size.X),
					float64(material.Body2(th, c.PeerName).Layout(fgtx(gtx)).Size.X),
				))),
			}.Layout(gtx,
				func(gtx C) D {
					return layout.Flex{
						Axis:      layout.Vertical,
						Alignment: layout.Start,
					}.Layout(gtx,
						layout.Rigid(material.Body2(th, c.PeerName).Layout),
						layout.Rigid(layout.Spacer{Height: unit.Dp(7.5)}.Layout),
						layout.Rigid(material.Label(th, unit.Dp(12.5), getSmallStr(c)).Layout),
					)
				},
			)
		})
	})
}

// GUIMessage is message
type GUIMessage struct {
	From string
	Text string
}

func getSmallStr(c *Chat) string {
	if len(c.Messages) == 0 {
		return "[нет сообщений]"
	}
	t := c.Messages[len(c.Messages)-1].Text
	if len([]rune(t)) > 21 {
		t = string([]rune(t)[:18]) + "..."
	}
	return t
}
