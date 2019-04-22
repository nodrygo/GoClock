// Gtk Go Clock demo
// License MIT

package main

import (
	"fmt"
	"math"
	"time"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	"image/color"

	"golang.org/x/image/colornames"
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

// return windows size and radius
func getWinSize(win *gtk.Window) (float64, float64, float64) {
	wx, wy := win.GetSize()
	wfx := float64(wx)
	wfy := float64(wy)
	radius := math.Min(wfx, wfy) / 2
	return wfx, wfy, radius
}

// convert named color to RGB as flot 64 between 0..1 needed for cairo cr
func colorConvert(c color.RGBA) (r, g, b float64) {
	//fmt.Println("Color is ", c)
	cr, cg, cb := c.R, c.G, c.B
	//fmt.Printf("R %d, G %d B %d   ", cr, cb, cg)
	r = float64(cr) / 255.0
	g = float64(cg) / 255.0
	b = float64(cb) / 255.0
	return
}

// prepare color for stop gradian as color name
func stopColorConvertAlpha(stop float64, c color.RGBA, alpha float64) (s, r, g, b, a float64) {
	//fmt.Println("Color is ", c)
	cr, cg, cb := colorConvert(c)
	//fmt.Printf("R %d, G %d B %d   ", cr, cb, cg)
	r = cr
	g = cg
	b = cb
	s = stop
	a = alpha

	return
}

// draw neddles at angle stat at cx,cy for length cl and width siz
func drawNeddle(cr *cairo.Context, cx, cy, cl, siz, angle float64) {
	cr.Save()
	cr.MoveTo(cx, cy)
	r, g, b := colorConvert(colornames.Chocolate)
	//fmt.Println("color", r, " ", b, " ", g)
	cr.SetSourceRGBA(r, g, b, 0.8)
	cr.SetLineWidth(siz)
	cr.LineTo(cx+cl*math.Sin(angle), cy+cl*math.Cos(angle))
	cr.Stroke()
	cr.Restore()
}

// draw face minutes decorations for each 5mn center cx,cy with radius
func drawMinutes(cr *cairo.Context, cx, cy, radius float64) {
	lstart := radius - radius/3.8
	lend := lstart + radius/10
	cr.Save()
	cr.SelectFontFace("Helvetica", cairo.FONT_SLANT_NORMAL, cairo.FONT_WEIGHT_BOLD)
	cr.SetFontSize(radius / 8)
	cr.MoveTo(cx, cy)

	cr.SetLineWidth(3)
	for i := 12; i > 0; i-- {
		cr.SetSourceRGB(0, 0, 0)
		angle := (float64(i*5) * -(math.Pi / 30.0)) + math.Pi
		cr.MoveTo(cx+lstart*math.Sin(angle), cy+lstart*math.Cos(angle))
		cr.LineTo(cx+lend*math.Sin(angle),
			cy+lend*math.Cos(angle))
		cr.Stroke()
		tsize := cr.TextExtents(fmt.Sprint(i))
		tdh := tsize.Height / 2
		tdw := tsize.Width / 2
		// adapt text pos according minutes pos and text size
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
			cr.MoveTo(cx-tdw+lend*math.Sin(angle), cy+2+tdh*2+lend*math.Cos(angle))
		case 7, 8:
			cr.MoveTo(cx+(tdw*2+4.0+lend)*math.Sin(angle), cy+(tdh*2+4.0+lend)*math.Cos(angle))
		case 9:
			cr.MoveTo(cx+(tdw*2+4.0+lend)*math.Sin(angle), cy+tdh+lend*math.Cos(angle))
		case 10, 11:
			cr.MoveTo(cx+(tdw*2+4.0+lend)*math.Sin(angle), cy+(tdh+lend)*math.Cos(angle))
		case 12:
			cr.MoveTo(cx-tdw+lend*math.Sin(angle), cy+(tdh+lend)*math.Cos(angle))
		}
		cr.SetSourceRGB(colorConvert(colornames.Darkslateblue))
		cr.ShowText(fmt.Sprint(i))
	}
	cr.Restore()
}

// draw clock face
func drawFace(cr *cairo.Context, cx, cy, radius float64) {
	h := 0
	cr.SetSourceRGBA(255, 255, 255, 255)
	cr.Paint()

	cr.SetSourceRGB(colorConvert(colornames.Chocolate))
	cr.ShowText(string(h))
	// set color gradian
	p, _ := cairo.NewPatternLinear(0, 0, cx+radius, cy+radius)
	p.AddColorStopRGBA(stopColorConvertAlpha(0.0, colornames.Darkred, 0.3))
	p.AddColorStopRGBA(stopColorConvertAlpha(0.25, colornames.Red, 0.3))
	p.AddColorStopRGBA(stopColorConvertAlpha(0.45, colornames.Green, 0.3))
	p.AddColorStopRGBA(stopColorConvertAlpha(0.60, colornames.Darkgreen, 0.3))
	p.AddColorStopRGBA(stopColorConvertAlpha(0.8, colornames.Blue, 0.3))
	p.AddColorStopRGBA(stopColorConvertAlpha(1.0, colornames.Darkblue, 0.3))
	cr.SetSource(p)
	cr.Arc(cx, cy, radius, 0, math.Pi*2)
	cr.Fill()
	cr.SetLineWidth(3)
	cr.SetSourceRGB(0, 0, 0)
	cr.Arc(cx, cy, radius, 0, math.Pi*2)
	cr.Stroke()
	drawMinutes(cr, cx, cy, radius)

}

// main draw fct when canvas need redraw
func drawClock(cr *cairo.Context) {
	hour, min, sec := lastTime.Clock()
	yy, mm, dd := lastTime.Date()
	pi := math.Pi
	//halfpi := pi / 2
	hangle := pi + (float64(hour+min/90) * -(pi / 6.0))
	mangle := pi + (float64(min) * -(pi / 30.0))
	sangle := pi + (float64(sec) * -(pi / 30.0))
	maxradius := radius - 4
	cx := wfx / 2
	cy := wfy / 2
	// draw clock face
	drawFace(cr, cx, cy, maxradius)

	// draw second
	drawNeddle(cr, cx, cy, maxradius-10.0, 2.0, sangle)
	// draw min
	drawNeddle(cr, cx, cy, maxradius-(maxradius/3), 4.0, mangle)
	// draw hour
	drawNeddle(cr, cx, cy, maxradius-(maxradius/2), 8.0, hangle)
	cr.SetSourceRGB(0, 0, 0)
	//draw center point
	cr.Arc(cx, cy, maxradius/14, 0, math.Pi*2)
	cr.Fill()

	// draw digital clock inside clock
	r, g, b := colorConvert(colornames.Darkblue)
	cr.SetSourceRGB(r, g, b)
	cr.SelectFontFace("Helvetica", cairo.FONT_SLANT_NORMAL, cairo.FONT_WEIGHT_BOLD)
	cr.SetFontSize(radius / 5)
	txttime := fmt.Sprintf("%02d:%02d:%02d", hour, min, sec)
	tsize := cr.TextExtents(txttime)
	tdh := tsize.Height / 2
	tdw := tsize.Width / 2
	cr.MoveTo(cx-tdw, cy+maxradius/2.8-tdh)
	cr.ShowText(txttime)
	txtdate := fmt.Sprintf("%02d/%02d/%4d", dd, mm, yy)
	tsize = cr.TextExtents(txtdate)
	tdh = tsize.Height / 2
	tdw = tsize.Width / 2
	cr.MoveTo(cx-tdw, cy-maxradius/5)
	cr.ShowText(txtdate)
}

// hack to handle clock time seconds
func handleTick(widget *gtk.Widget, frameClock *gdk.FrameClock, userData uintptr) bool {
	_ = frameClock
	_, _, ls := lastTime.Clock()
	_, _, sec := time.Now().Clock()
	if ls != sec {
		lastTime = time.Now()
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

	// define the popp menu
	popupmenu, _ := gtk.MenuNew()
	psetAlarm, _ := gtk.MenuItemNew()
	psetAlarm.SetLabel("Set Alarm")
	pswitchDeco, _ := gtk.MenuItemNew()
	pswitchDeco.SetLabel("Switch decoration")
	pquit, _ := gtk.MenuItemNew()
	pquit.SetLabel("Quit")
	psep, _ := gtk.SeparatorMenuItemNew()
	popupmenu.Append(psetAlarm)
	popupmenu.Append(pswitchDeco)
	popupmenu.Append(psep)
	popupmenu.Append(pquit)

	// set Show for everything
	win.ShowAll()
	popupmenu.ShowAll()

	win.SetOpacity(0.9)
	clkcanvas.SetOpacity(1)
	wfx, wfy, radius = getWinSize(win)
	lastTime = time.Now()
	clkcanvas.AddTickCallback(handleTick, 1000)

	// Event handlers
	pquit.Connect("activate", func() {
		gtk.MainQuit()
	})
	pswitchDeco.Connect("activate", func() {
		fmt.Println("menuitem pswitchDeco")
		if win.GetDecorated() {
			win.SetDecorated(false)
			win.SetOpacity(0.9)
			clkcanvas.SetOpacity(1)
		} else {
			win.SetDecorated(true)
			win.SetOpacity(0.9)
			clkcanvas.SetOpacity(1)
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

	// Another way to create timer with concurrency
	// go func() {
	// 	for now := range time.Tick(time.Second) {
	// 		fmt.Println(now)
	// 		clkcanvas.Emit("draw")
	// 	}
	// }()
	gtk.Main()
}
