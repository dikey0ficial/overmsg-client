package main

import (
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/sqweek/dialog"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image"
	"image/color"
	"math"
	"os"
	"strings"
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
	ChatList *ChatList
	ChatAct  *ChatActivity
	Size     image.Point
	Win      *app.Window
}

// NewUI is constructor for UI
func NewUI() *UI {
	ui := new(UI)
	ui.SetTheme(conf.IsDark)
	ui.ChatList = new(ChatList)
	ui.ChatAct = new(ChatActivity)
	charr := []*Chat{
		&Chat{"1", []GUIMessage{
			{"1", "QWE42"},
			{"1", "hello world"},
		}, new(widget.Clickable)},
		&Chat{"Vasyok", []GUIMessage{}, new(widget.Clickable)},
		&Chat{"Qwertyque", []GUIMessage{
			{"Qwertyque", "lorem ipsum dolor sit amet((("},
		}, new(widget.Clickable)},
		&Chat{"_", []GUIMessage{
			{"_", "Прости, что тебе спамлю, но... я шизофреник !!!!!!!!!!!!!"},
			{"admin", "да я уже понял)"},
			{"_", "Прости, что тебе спамлю, но... я шизофреник !!!!!!!!!!!!!"},
			{"admin", "да я уже понял)"},
			{"_", "Прости, что тебе спамлю, но... я шизофреник !!!!!!!!!!!!!"},
			{"admin", "да я уже понял)"},
			{"_", "Прости, что тебе спамлю, но... я шизофреник !!!!!!!!!!!!!"},
			{"admin", "да я уже понял)"},
			{"_", "Прости, что тебе спамлю, но... я шизофреник !!!!!!!!!!!!!"},
			{"admin", "да я уже понял)"},
			{"_", "Прости, что тебе спамлю, но... я шизофреник !!!!!!!!!!!!!"},
			{"admin", "да я уже понял)"},
			{"_", "Прости, что тебе спамлю, но... я шизофреник !!!!!!!!!!!!!"},
			{"admin", "да я уже понял)"},
			{"_", "Прости, что тебе спамлю, но... я шизофреник !!!!!!!!!!!!!"},
			{"admin", "да я уже понял)"},
			{"_", "Прости, что тебе спамлю, но... я шизофреник !!!!!!!!!!!!!"},
			{"admin", "да я уже понял)"},
			{"_", "Прости, что тебе спамлю, но... я шизофреник !!!!!!!!!!!!!"},
			{"admin", "да я уже понял)"},
		}, new(widget.Clickable)},
	}
	ui.ChatList.Chats = charr
	ui.ChatList.List = &layout.List{Axis: layout.Vertical}
	ui.ChatAct.List = &widget.List{List: layout.List{Axis: layout.Vertical, ScrollToEnd: true}}
	ui.ChatAct.SendBtn = material.IconButton(
		ui.Theme,
		new(widget.Clickable),
		getIcon(icons.ContentSend),
		"Send message",
	)
	ui.ChatAct.SendBtn.Size = unit.Dp(15)
	ui.ChatAct.Input = material.Editor(
		ui.Theme,
		&widget.Editor{
			SingleLine: true,
			Submit:     true,
		},
		"Type your message here...",
	)
	ui.ChatList.Selected = "_home"
	ui.ChatList.HomeTab = new(HomeTab)
	ui.ChatList.HomeTab.ListButton = material.Button(ui.Theme, new(widget.Clickable), "OVERMSg")
	ui.ChatList.HomeTab.ListButton.Font.Weight = text.Bold
	themeb := new(widget.Bool)
	themeb.Value = conf.IsDark
	ui.ChatList.HomeTab.ThemeSwitch = material.Switch(
		ui.Theme,
		themeb,
		"Dark theme",
	)
	ui.ChatList.HomeTab.NameInput = material.Editor(
		ui.Theme,
		&widget.Editor{
			SingleLine: true,
			Submit:     true,
		},
		"Type your name here...",
	)
	ui.ChatList.HomeTab.RegBtn = material.Button(ui.Theme, new(widget.Clickable), "Register")
	ui.ChatList.HomeTab.LogoutBtn = material.Button(ui.Theme, new(widget.Clickable), "Log out")
	ui.ChatAct.HomeTab = ui.ChatList.HomeTab
	return ui
}

// SetTheme sets dark or light theme
func (ui *UI) SetTheme(dark bool) {
	ui.Theme = material.NewTheme(gofont.Collection())
	if conf.IsDark {
		ui.Theme.Palette.Bg = color.NRGBA{R: 22, G: 27, B: 34, A: 255}
		ui.Theme.Palette.Fg = color.NRGBA{R: 201, G: 209, B: 217, A: 255}
		ui.Theme.Palette.ContrastFg = color.NRGBA{R: 253, G: 253, B: 253, A: 255}
	}
}

// Run starts layouting
func (ui *UI) Run(w *app.Window) error {
	ui.Win = w
	var ops op.Ops
	for {
		select {
		case e := <-w.Events():
			switch e := e.(type) {
			case system.FrameEvent:
				if e.Size.X >= 1920 {
					e.Size.X -= 80 * (1920 / e.Size.X)
				}
				if e.Size.Y >= 1080 {
					e.Size.Y -= 25 * (1080 / e.Size.Y)
				}
				gtx := layout.NewContext(&ops, e)
				if conf.IsDark {
					paint.Fill(&ops, ui.Theme.Palette.Bg)
				}
				ui.Size = e.Size
				ui.Layout(gtx)
				ui.ChatAct.Selected, ui.ChatAct.Chat = ui.ChatList.Selected, GetByPN(ui.ChatList.Chats, ui.ChatList.Selected)
				e.Frame(gtx.Ops)
			case system.DestroyEvent:
				return e.Err
			}
		}
	}
}

// Layout layouts
func (ui *UI) Layout(gtx C) D {
	ui.ChatList.MaxX = ui.Size.X
	ui.ChatAct.MaxX = ui.Size.X
	var x *int
	return layout.Inset{
		Top:    unit.Dp(stdDP.V * 2),
		Bottom: unit.Dp(stdDP.V * 2),
		Left:   stdDP,
		Right:  stdDP,
	}.Layout(gtx, func(gtx C) D {
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
							return ui.ChatAct.Layout(gtx, ui.Theme, *x, ui)
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
	HomeTab  *HomeTab
	List     *layout.List
}

// Layout _
func (cl *ChatList) Layout(gtx C, th T) D {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(
			func(gtx C) D {
				return cl.List.Layout(gtx, len(cl.Chats)+1, func(gtx C, ind int) D {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							if ind == 0 {
								return cl.HomeTab.LayoutList(gtx, th, cl)
							}
							return cl.Chats[ind-1].LayoutList(gtx, th, cl)
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
	List     *widget.List
	Input    material.EditorStyle
	SendBtn  material.IconButtonStyle
	HomeTab  *HomeTab
	Chat     *Chat
}

// Layout _
func (ca *ChatActivity) Layout(gtx C, th T, startX int, ui *UI) D {
	if ca.Selected == "" {
		ca.Selected = "_home"
	}
	return layout.Flex{
		Axis:    layout.Vertical,
		Spacing: layout.SpaceEnd,
	}.Layout(gtx,
		// header
		layout.Rigid(func(gtx C) D {
			return widget.Border{Color: th.ContrastBg, Width: unit.Dp(3.5)}.Layout(gtx,
				func(gtx C) D {
					return layout.UniformInset(unit.Dp(5)).Layout(gtx,
						func(gtx C) D {
							return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
								layout.Rigid(material.Body2(th, func() string {
									if ca.Selected == "_home" {
										return "Start page"
									}
									return "Chat with " + ca.Selected
								}()).Layout),
								layout.Rigid(layout.Spacer{Width: unit.Px(
									float32(ca.MaxX)/1.07 - // hello, hardcoded number! (got it during experiments)
										float32(
											// 20 because it is sum of spacer and insets (5 + 10 + 5)
											20+startX+material.Body2(th, func() string {
												if ca.Selected == "_home" {
													return "Start page"
												}
												return "Chat with " + ca.Selected
											}()).Layout(fgtx(gtx)).Size.X,
										),
								)}.Layout),
							)
						},
					)
				},
			)
		}),
		// messages
		layout.Rigid(func(gtx C) D {
			if ca.Selected == "_home" {
				return ca.HomeTab.Layout(gtx, th, ui)
			}
			if ca.Chat.PeerName != ca.Selected {
				return D{}
			}
			if len(ca.Chat.Messages) == 0 {
				gx := *(&gtx)
				gx.Constraints.Max.Y -= 55
				gx.Constraints.Min.Y = gx.Constraints.Max.Y
				return layout.Flex{Alignment: layout.Middle, Axis: layout.Vertical}.Layout(gx,
					layout.Rigid(layout.Spacer{Height: unit.Dp(15)}.Layout),
					layout.Rigid(material.Body2(th, "There's nothing...").Layout),
				)
			}
			return layout.UniformInset(unit.Dp(15)).Layout(gtx, func(gtx C) D {
				gx := *(&gtx)
				gx.Constraints.Max.Y -= 55
				return material.List(th, ca.List).Layout(
					gx,
					len(ca.Chat.Messages),
					func(gtx C, ind int) D { return ca.Chat.Messages[ind].Layout(gtx, th, ca.Chat.PeerName) },
				)
			},
			)
		}),
		layout.Rigid(func(gtx C) D {
			if ca.Selected == "_home" {
				return D{}
			}
			if ca.SendBtn.Button.Clicked() || func() bool {
				evs := ca.Input.Editor.Events()
				if len(evs) == 0 {
					return false
				}
				for i := range evs {
					switch evs[i].(type) {
					case widget.SubmitEvent:
						return true
					}
				}
				return false
			}() {
				txt := strings.TrimSpace(ca.Input.Editor.Text())
				if len([]rune(txt)) != 0 {
					ca.Input.Editor.SetText("")
					ca.Chat.Messages = append(ca.Chat.Messages, GUIMessage{conf.Name, txt})
				}
			}
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(
					func(gtx C) D {
						gx := *(&gtx)
						gx.Constraints.Max.X -= 90
						gx.Constraints.Min.X = gx.Constraints.Max.X
						return widget.Border{
							Width:        unit.Dp(0.5),
							Color:        th.Fg,
							CornerRadius: unit.Dp(3),
						}.Layout(gx,
							func(gtx C) D {
								return layout.UniformInset(unit.Dp(5)).Layout(gtx, ca.Input.Layout)
							},
						)
					},
				),
				layout.Rigid(layout.Spacer{Width: unit.Dp(15)}.Layout),
				layout.Rigid(ca.SendBtn.Layout),
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

// GetByPN returns chat by peername
func GetByPN(arr []*Chat, pn string) *Chat {
	for _, c := range arr {
		if c.PeerName == pn {
			return c
		}
	}
	return &Chat{}
}

func getSmallStr(c *Chat) string {
	if len(c.Messages) == 0 {
		return "[no messages]"
	}
	t := c.Messages[len(c.Messages)-1].Text
	if len([]rune(t)) > 21 {
		t = string([]rune(t)[:18]) + "..."
	}
	return t
}

// LayoutList layouts list small preview
func (c *Chat) LayoutList(gtx C, th T, cl *ChatList) D {
	maxX := cl.MaxX
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

// Layout layouts
func (g GUIMessage) Layout(gtx C, th T, chname string) D {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			gx := *(&gtx)
			gx.Constraints.Min.X = gx.Constraints.Max.X
			return layout.Flex{Axis: layout.Horizontal}.Layout(gx,
				layout.Rigid(material.Body2(func() T {
					t := *th
					t.Fg = color.NRGBA{R: 255, G: 127, A: 255}
					if g.From == conf.Name {
						t.Fg = color.NRGBA{G: 127, B: 127, A: 255}
					}
					return &t
				}(), "<"+g.From+">\t").Layout),
				layout.Rigid(material.Label(th, unit.Dp(15), g.Text).Layout),
			)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
	)
}

// HomeTab is tab which shows on start
type HomeTab struct {
	ListButton  material.ButtonStyle
	Settings    widget.Bool
	ThemeSwitch material.SwitchStyle
	NameInput   material.EditorStyle
	RegBtn      material.ButtonStyle
	LogoutBtn   material.ButtonStyle
}

// LayoutList layouts HomeTab's view in list
func (ht *HomeTab) LayoutList(gtx C, th T, cl *ChatList) D {
	gx := *(&gtx)
	gx.Constraints.Min.X = cl.MaxX/4 + 5
	gx.Constraints.Max.X = gx.Constraints.Min.X
	ht.ListButton.Inset.Top = unit.Dp(5)
	ht.ListButton.Inset.Bottom = unit.Dp(5)
	ht.ListButton.Background = th.Palette.ContrastBg
	if ht.ListButton.Button.Clicked() {
		cl.Selected = "_home"
	}
	if cl.Selected == "_home" {
		ht.ListButton.Background.G += 25
	}
	return ht.ListButton.Layout(gx)
}

const allowedSymbols = "QWERTYUIOPASDFGHJKLZXCVBNM" +
	"qwertyuiopasdfghjklzxcvbnm" +
	"0123456789" +
	"_-"

// Layout layouts HomeTab's view instead of chat
func (ht *HomeTab) Layout(gtx C, th T, ui *UI) D {
	if ht.ThemeSwitch.Switch.Changed() {
		conf.IsDark = ht.ThemeSwitch.Switch.Value
		err := saveConf()
		if err != nil {
			dialog.Message("Error: %v", err).Title("Error!!1").Error()
			errl.Println(err)
			os.Exit(1)
		}
		lastUI := *ui
		*ui = *NewUI()
		ui.ChatList.HomeTab.Settings.Value = true
		ui.Win = lastUI.Win
		ui.Win.Invalidate()
		return D{}
	}
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		hspacer,
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					var icon *widget.Icon = getIcon(icons.NavigationUnfoldMore)
					if ht.Settings.Value {
						icon = getIcon(icons.NavigationUnfoldLess)
					}
					return icon.Layout(gtx, th.Fg)
				}),
				wspacer,
				layout.Rigid(func(gtx C) D { return ht.Settings.Layout(gtx, material.H4(th, "Settings").Layout) }),
			)
		}),
		hspacer,
		layout.Rigid(func(gtx C) D {
			if ht.Settings.Value {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
							layout.Rigid(material.Label(th, unit.Dp(15), "Dark theme:\t").Layout),
							layout.Rigid(layout.Spacer{Width: unit.Dp(5)}.Layout),
							layout.Rigid(ht.ThemeSwitch.Layout),
						)
					}),
					hspacer,
					layout.Rigid(material.H5(th, "Account:\t").Layout),
					hspacer,
					layout.Rigid(func(gtx C) D {
						if conf.Name == "" {
							var warn string
							txt := strings.TrimSpace(ht.NameInput.Editor.Text())
							if strings.HasPrefix(txt, "_") {
								warn = "Nick shouldn't start with '_'"
							}
							var isValid bool = true
						BIG:
							for _, sym := range []rune(txt) {
								for _, asym := range allowedSymbols {
									if sym == asym {
										continue BIG
									}
								}
								isValid = false
								break
							}
							if !isValid {
								warn = "Nick contains illegal symbols; allowed only english alphabet, numbers, underscore and dash"
							}
							if len([]rune(ht.NameInput.Editor.Text())) > 32 {
								ht.NameInput.Editor.Delete(
									-(len([]rune(ht.NameInput.Editor.Text())) - 32),
								)
							}
							if ht.RegBtn.Button.Clicked() && warn == "" {
								conf.Name = txt
								// TODO: add server registration
								err := saveConf()
								if err != nil {
									dialog.Message("Error: %v", err).Title("Error!!1").Error()
									errl.Println(err)
									os.Exit(1)
								}
							}
							return layout.Flex{
								Axis:      layout.Horizontal,
								Alignment: layout.Middle,
							}.Layout(gtx,
								layout.Rigid(material.Label(th, unit.Dp(15), "Name:\t").Layout),
								layout.Rigid(func(gtx C) D {
									col, wid := th.Fg, unit.Dp(0.5)
									if warn != "" {
										col, wid = color.NRGBA{R: 255, A: 255}, unit.Dp(1)
									}
									return widget.Border{
										CornerRadius: unit.Dp(5),
										Color:        col,
										Width:        wid,
									}.Layout(gtx, func(gtx C) D {
										return layout.UniformInset(unit.Dp(4)).Layout(gtx, ht.NameInput.Layout) // TODO: make more strict checks
									})
								}),
								wspacer,
								layout.Rigid(func(gtx C) D {
									if warn != "" {
										return material.Label(th, unit.Dp(15), warn).Layout(gtx)
									}
									return ht.RegBtn.Layout(gtx)
								}),
							)
						}
						if ht.LogoutBtn.Button.Clicked() {
							ok := dialog.Message("Do you realy want to logout? You'll lose this nickname forever").YesNo()
							if ok {
								ok = dialog.Message("Really?").YesNo()
								if ok {
									conf.Name, conf.Token = "", ""
									ui.ChatList.Chats = []*Chat{}
									ui.Win.Invalidate()
									return D{}
								}
							}
						}
						return ht.LogoutBtn.Layout(gtx)
					}),
				)
			}
			return D{}
		}),
	)
}

func getIcon(dat []byte) *widget.Icon {
	ic, _ := widget.NewIcon(dat)
	return ic
}
