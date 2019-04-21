package main

import (
	"fmt"
	"math"
	"time"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

const (
	KEY_LEFT  uint = 65361
	KEY_UP    uint = 65362
	KEY_RIGHT uint = 65363
	KEY_DOWN  uint = 65364
)

// global var
var wfx float64
var wfy float64
var radius float64
var lastTime time.Time

func getWinSize(win *gtk.Window) (float64, float64, float64) {
	wx, wy := win.GetSize()
	wfx := float64(wx)
	wfy := float64(wy)
	radius := math.Min(wfx, wfy)
	return wfx, wfy, radius
}

func drawNeddle(cr *cairo.Context, cx, cy, cl, siz, angle float64) {
	cr.Save()
	cr.MoveTo(cx, cy)
	cr.SetSourceRGB(0, 0, 0)
	cr.SetLineWidth(siz)

	cr.LineTo(cx+cl*math.Sin(angle),
		cy+cl*math.Cos(angle))
	cr.Stroke()
	cr.Restore()
}
func drawMinutes(cr *cairo.Context, cx, cy, halfrad float64) {
	lstart := halfrad - halfrad/4
	lend := lstart + 14
	cr.Save()
	cr.SelectFontFace("Helvetica", cairo.FONT_SLANT_NORMAL, cairo.FONT_WEIGHT_BOLD)
	cr.SetFontSize(halfrad / 10)
	cr.MoveTo(cx, cy)
	cr.SetSourceRGB(0, 0, 0)
	cr.SetLineWidth(2)
	for i := 12; i > 0; i-- {
		angle := (float64(i*5) * -(math.Pi / 30.0)) + math.Pi
		cr.MoveTo(cx+lstart*math.Sin(angle), cy+lstart*math.Cos(angle))
		cr.LineTo(cx+lend*math.Sin(angle),
			cy+lend*math.Cos(angle))
		cr.Stroke()
		tsize := cr.TextExtents(fmt.Sprint(i))
		tdh := tsize.Height / 2
		tdw := tsize.Width / 2
		switch i {
		default:
			cr.MoveTo(cx+(lend-tdw)*math.Sin(angle), cy+(lend+tdh)*math.Cos(angle))
		case 1, 2:
			cr.MoveTo(cx+tdw+lend*math.Sin(angle), cy+(tdh+lend)*math.Cos(angle))
		case 3:
			cr.MoveTo(cx+tdw+lend*math.Sin(angle), cy+tdh+lend*math.Cos(angle))
		case 4, 5:
			cr.MoveTo(cx+tdw+lend*math.Sin(angle), cy+4.0+tdh+lend*math.Cos(angle))
		case 6:
			cr.MoveTo(cx-tdw+lend*math.Sin(angle), cy+tdh*2+lend*math.Cos(angle))
		case 7, 8:
			cr.MoveTo(cx+(tdw*2+4.0+lend)*math.Sin(angle), cy+(tdh*2+4.0+lend)*math.Cos(angle))
		case 9:
			cr.MoveTo(cx+(tdw*2+4.0+lend)*math.Sin(angle), cy+tdh+lend*math.Cos(angle))
		case 10, 11:
			cr.MoveTo(cx+(tdw*2+4.0+lend)*math.Sin(angle), cy+(tdh+lend)*math.Cos(angle))
		case 12:
			cr.MoveTo(cx-tdw+lend*math.Sin(angle), cy+(tdh+lend)*math.Cos(angle))
		}

		cr.ShowText(fmt.Sprint(i))
	}
	cr.Restore()
}
func drawFace(cr *cairo.Context, cx, cy, halfrad float64) {
	h := 0
	cr.SetSourceRGB(255, 255, 255)
	cr.Paint()
	cr.ShowText(string(h))
	cr.SetSourceRGB(255, 0, 0)
	//p, _ := cairo.NewPatternLinear()l(cx, cy, 0, cx+halfrad, cy+halfrad, halfrad)
	p, _ := cairo.NewPatternLinear(0, 0, cx+halfrad, cy+halfrad)
	p.AddColorStopRGBA(0, 255, 0, 0, 0.3)
	p.AddColorStopRGBA(0.5, 255, 255, 0, 0.3)
	p.AddColorStopRGBA(1, 255, 0, 255, 0.3)
	cr.SetSource(p)
	cr.Arc(cx, cy, halfrad, 0, math.Pi*2)
	cr.Fill()
	cr.SetLineWidth(3)
	cr.SetSourceRGB(0, 0, 0)
	cr.Arc(cx, cy, halfrad, 0, math.Pi*2)
	cr.Stroke()
	drawMinutes(cr, cx, cy, halfrad)

}

func drawClock(cr *cairo.Context) {
	hour, min, sec := lastTime.Clock()
	pi := math.Pi
	//halfpi := pi / 2
	hangle := pi + (float64(hour+min/90) * -(pi / 6.0))
	mangle := pi + (float64(min) * -(pi / 30.0))
	sangle := pi + (float64(sec) * -(pi / 30.0))
	halfrad := (radius / 2) - 4
	cx := wfx / 2
	cy := wfy / 2
	// draw clock face
	drawFace(cr, cx, cy, halfrad)
	// draw second
	drawNeddle(cr, cx, cy, halfrad-10.0, 2.0, sangle)
	// draw min
	drawNeddle(cr, cx, cy, halfrad-(halfrad/3), 4.0, mangle)
	// draw hour
	drawNeddle(cr, cx, cy, halfrad-(halfrad/2), 8.0, hangle)
	cr.SetSourceRGB(0, 0, 0)
	cr.Arc(cx, cy, halfrad/14, 0, math.Pi*2)
	cr.Fill()
}

func handleTick(widget *gtk.Widget, frameClock *gdk.FrameClock, userData uintptr) bool {
	_ = frameClock
	_, _, ls := lastTime.Clock()
	hour, min, sec := time.Now().Clock()
	if ls != sec {
		lastTime = time.Now()
		fmt.Printf("Time:  %d:%d:%d\n", hour, min, sec)
		widget.QueueDraw()
	}
	return true
}

func main() {
	gtk.Init(nil)
	win, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	win.SetDefaultSize(300, 300)
	clkcanvas, _ := gtk.DrawingAreaNew()
	win.Add(clkcanvas)
	win.SetTitle("Clock")
	win.Connect("destroy", gtk.MainQuit)

	popupmenu, _ := gtk.MenuNew()
	psetAlarm, _ := gtk.MenuItemNew()
	psetAlarm.SetLabel("Set Alarm")
	psetAlarm.SetName("SetAlarm")
	pswitchDeco, _ := gtk.MenuItemNew()
	pswitchDeco.SetLabel("Switch decoration")
	popupmenu.Append(psetAlarm)
	popupmenu.Append(pswitchDeco)

	win.ShowAll()
	popupmenu.ShowAll()

	win.SetOpacity(0.8)
	clkcanvas.SetOpacity(1)
	wfx, wfy, radius = getWinSize(win)
	lastTime = time.Now()
	clkcanvas.AddTickCallback(handleTick, 1000)
	// Event handlers

	pswitchDeco.Connect("activate", func() {
		fmt.Println("menuitem pswitchDeco")
		if win.GetDecorated() {
			win.SetDecorated(false)
			win.SetOpacity(0.8)
		} else {
			win.SetDecorated(true)
			win.SetOpacity(0.8)
		}
	})

	win.Connect("check_resize", func(win *gtk.Window) {
		wfx, wfy, radius = getWinSize(win)
	})

	clkcanvas.Connect("draw", func(da *gtk.DrawingArea, cr *cairo.Context) {
		drawClock(cr)
	})
	win.Connect("button-press-event", func(win *gtk.Window, ev *gdk.Event) {
		fmt.Println("Event", ev)
		button := &gdk.EventButton{ev}
		if button.Button() == 3 {
			// posx := button.X()
			// posy := button.Y()
			//fmt.Println("Button 3 detected ", posx, ":", posy)
			popupmenu.PopupAtPointer(ev)
			win.QueueDraw()
		}
	})

	// Another way to create timer
	// go func() {
	// 	for now := range time.Tick(time.Second) {
	// 		fmt.Println(now)
	// 		clkcanvas.Emit("draw")
	// 	}
	// }()
	gtk.Main()
}
