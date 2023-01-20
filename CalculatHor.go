package main

import (
	"fmt"
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image/color"
	"log"
	"os"
	"time"
)

var timerIncrementer chan int64

type Shift struct {
	timer          time.Time
	timerStartTime time.Time
	timerStopTime  time.Time
	started        bool
	justStarted    bool
	justStopped    bool
	lastTimer      time.Time
}

type Schedule struct {
	shifts []Shift
}

var shift Shift
var file string

func main() {
	if len(os.Args) == 2 {
		file = os.Args[1]
	}

	timerIncrementer = make(chan int64)
	go func() {
		for {
			time.Sleep(time.Second)
			timerIncrementer <- int64(time.Second)
		}
	}()
	go func() {
		w := app.NewWindow(
			app.Title("CalculatHor"),
			app.Size(unit.Dp(400), unit.Dp(450)),
		)

		err := run(w)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(w *app.Window) error {
	var ops op.Ops
	var startButton widget.Clickable
	var fileButton widget.Clickable

	th := material.NewTheme(gofont.Collection())

	for {
		select {
		case e := <-w.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				return e.Err
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)

				if startButton.Clicked() {
					shift.started = !shift.started
					if shift.started {
						shift.justStarted = !shift.justStarted
					} else {
						shift.justStopped = !shift.justStopped
						shift.lastTimer = shift.timer
						shift.timer = time.Time{}
					}
				}

				if fileButton.Clicked() {
					if !shift.started {
						saveFile(shift.timerStartTime, shift.timerStopTime, shift.lastTimer)
					} else {

					}
				}

				layout.Flex{
					Axis: layout.Vertical,
					//Spacing:   layout.SpaceStart,
					Alignment: layout.End,
				}.Layout(gtx,

					layout.Rigid(
						func(gtx layout.Context) layout.Dimensions {
							btn := material.Button(th, &fileButton, "Sauvegarder")
							_color := color.NRGBA{R: 63, G: 181, B: 181, A: 255}
							btn.Background = _color
							return btn.Layout(gtx)
						},
					),

					// Date et heure du départ du chronometre
					layout.Rigid(
						func(gtx layout.Context) layout.Dimensions {
							margins := layout.Inset{
								Top:    unit.Dp(25),
								Bottom: unit.Dp(25),
								Right:  unit.Dp(25),
								Left:   unit.Dp(25),
							}
							return margins.Layout(gtx,
								func(gtx layout.Context) layout.Dimensions {
									if shift.justStarted {
										shift.timerStartTime = time.Now()
										shift.justStarted = !shift.justStarted
									}
									beginingDate := "Début : \n" + shift.timerStartTime.Format("2006-01-02 15:04:05")
									timerCount := material.Label(th, unit.Sp(20), beginingDate)
									timerCount.Alignment = text.Middle
									return timerCount.Layout(gtx)
								},
							)
						},
					),
					// Date et heure de la fin du chronometre
					layout.Rigid(
						func(gtx layout.Context) layout.Dimensions {
							margins := layout.Inset{
								Top:    unit.Dp(25),
								Bottom: unit.Dp(25),
								Right:  unit.Dp(25),
								Left:   unit.Dp(25),
							}
							return margins.Layout(gtx,
								func(gtx layout.Context) layout.Dimensions {
									if shift.justStopped {
										shift.timerStopTime = time.Now()
										shift.justStopped = !shift.justStopped
									}
									finalDate := "Fin : \n" + shift.timerStopTime.Format("2006-01-02 15:04:05")
									timerCount := material.Label(th, unit.Sp(20), finalDate)
									timerCount.Alignment = text.Middle
									return timerCount.Layout(gtx)
								},
							)
						},
					),

					// Affichage chronometre
					layout.Rigid(
						func(gtx layout.Context) layout.Dimensions {
							margins := layout.Inset{
								Top:    unit.Dp(25),
								Bottom: unit.Dp(25),
								Right:  unit.Dp(25),
								Left:   unit.Dp(25),
							}
							return margins.Layout(gtx,
								func(gtx layout.Context) layout.Dimensions {
									timerCount := material.Label(th, unit.Sp(40), shift.timer.Format("15:04:05"))
									timerCount.Alignment = text.Middle
									return timerCount.Layout(gtx)
								},
							)
						},
					),

					// Bouton demarrage
					layout.Rigid(
						func(gtx layout.Context) layout.Dimensions {
							margins := layout.Inset{
								Top:    unit.Dp(25),
								Bottom: unit.Dp(25),
								Right:  unit.Dp(25),
								Left:   unit.Dp(25),
							}
							return margins.Layout(gtx,
								func(gtx layout.Context) layout.Dimensions {
									var text string
									var _color color.NRGBA
									if shift.started {
										text = "Finish"
										_color = color.NRGBA{R: 163, G: 81, B: 181, A: 255}
									} else {
										text = "Start"
										_color = color.NRGBA{R: 63, G: 81, B: 181, A: 255}
									}

									btn := material.Button(th, &startButton, text)
									btn.Background = _color
									return btn.Layout(gtx)
								},
							)
						},
					),
				)
				e.Frame(gtx.Ops)
			}
		case p := <-timerIncrementer:
			if shift.started {
				shift.timer = shift.timer.Add(time.Duration(p))
				w.Invalidate()
			}
		}
	}
	return nil
}

func saveFile(startTime time.Time, stopTime time.Time, timer time.Time) {

	err := os.Mkdir("Schedule", 0666)
	if err != nil {
		fmt.Println("Ne peut pas creer le dossier")
	}

	file, err := os.OpenFile("Schedule/"+startTime.Format("2006 01 02")+".txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	text := "Début : " + startTime.Format("2006-01-02 15:04:05") + "\n"
	text += "Fin : :" + stopTime.Format("2006-01-02 15:04:05") + "\n"
	text += timer.Format("15:04:05") + "\n\n"
	_, err2 := file.WriteString(text)

	if err2 != nil {
		log.Fatal(err2)
	}

	fmt.Println("done")
}

//layout.Rigid(
//	layout.Spacer{Height: unit.Dp(25)}.Layout,
//),
