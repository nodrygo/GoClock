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

func drawHMneddle(cr *cairo.Context, cx, cy, cl, siz float64, col color.RGBA, angle float64) {
	cr.Save()
	cr.Translate(cx, cy)
	cr.Rotate(angle)
	r, g, b := colorConvert(col)
	cr.SetSourceRGBA(r, g, b, 0.55)
	cr.SetLineWidth(siz)
	cr.MoveTo(0, 0)
	cr.LineTo(cl/3, cl/10)
	cr.LineTo(cl, 0)
	cr.LineTo(cl/3, -cl/10)
	cr.LineTo(0, 0)
	cr.Fill()
	cr.SetSourceRGBA(0, 0, 1, 0.5)
	cr.MoveTo(0, 0)
	cr.LineTo(cl/3, cl/10)
	cr.LineTo(cl, 0)
	cr.LineTo(cl/3, -cl/10)
	cr.LineTo(0, 0)
	cr.Stroke()
	cr.Restore()
}

// draw neddles at angle stat at cx,cy for length cl and width siz
func drawSneddle(cr *cairo.Context, cx, cy, cl, siz float64, col color.RGBA, angle float64) {
	cr.Save()
	cr.Translate(cx, cy)
	cr.Rotate(angle)
	r, g, b := colorConvert(col)
	cr.SetSourceRGBA(r, g, b, 0.45)
	cr.SetLineWidth(siz)
	cr.MoveTo(0, 0)
	cr.LineTo(cl, 0)
	cr.Stroke()
	cr.Arc(cl/1.4, 0, cl/30+siz, 0, math.Pi*2)
	cr.Fill()
	cr.Restore()
}

// draw graduations and texts faces
func drawGraduations(cr *cairo.Context, cx, cy, radius float64) {
	lstart := radius - radius/3.8
	lend := lstart + radius/10
	fsize := radius / 8
	cr.SelectFontFace("Helvetica", cairo.FONT_SLANT_NORMAL, cairo.FONT_WEIGHT_BOLD)
	cr.SetFontSize(fsize)
	cr.SetLineWidth(3)

	for i := 1; i < 13; i++ {
		cr.Save()
		angle := (float64(i*5) * (math.Pi / 30.0)) - math.Pi/2
		cr.Translate(cx, cy)
		cr.Rotate(angle)
		cr.MoveTo(lstart, 0)
		cr.LineTo(lend, 0)
		cr.Stroke()

		txt := fmt.Sprint(i)
		tsize := cr.TextExtents(txt)
		tdh := tsize.Height
		tdw := tsize.Width

		// adapt text pos according minutes pos and text size
		switch i {
		default:
			cr.MoveTo(lend+tdw, tdh)
		case 1, 2:
			cr.MoveTo(lend+tdw/2, tdh/2)
		case 3:
			cr.MoveTo(lend+tdw/2, tdh/2)
		case 4:
			cr.MoveTo(lend+tdw, tdh/2)
		case 5:
			cr.MoveTo(lend+tdw*1.5, tdh/3)
		case 6:
			cr.MoveTo(lend+tdw+fsize/2, tdh/2)
		case 7:
			cr.MoveTo(lend+tdw*2, tdh/3)
		case 8:
			cr.MoveTo(lend+tdw*2, -tdh/3)
		case 9:
			cr.MoveTo(lend+tdw*1.5, -tdh/2)
		case 10:
			cr.MoveTo(lend+tdw, -tdh)
		case 11:
			cr.MoveTo(lend+tdw/1.5, -tdh/1.5)
		case 12:
			cr.MoveTo(lend+tdw/3, -tdh/3)
		}
		cr.Rotate(-angle)
		cr.SetSourceRGB(colorConvert(colornames.Darkslateblue))
		cr.ShowText(txt)
		cr.Restore()
	}
}

// draw clock face
func drawFace(cr *cairo.Context, cx, cy, radius float64) {
	cr.SetSourceRGBA(255, 255, 255, 1)
	cr.Arc(cx, cy, radius, 0, math.Pi*2)
	cr.Paint()
	h := 0
	cr.SetSourceRGB(colorConvert(colornames.Chocolate))
	cr.ShowText(string(h))
	// set color gradian
	p, _ := cairo.NewPatternLinear(0, 0, cx+radius, cy+radius)
	p.AddColorStopRGBA(stopColorConvertAlpha(0, colornames.Blueviolet, 0.5))
	p.AddColorStopRGBA(stopColorConvertAlpha(0.30, colornames.Lightsteelblue, 0.5))
	p.AddColorStopRGBA(stopColorConvertAlpha(0.50, colornames.Lightgreen, 0.5))
	p.AddColorStopRGBA(stopColorConvertAlpha(0.80, colornames.Pink, 0.5))
	p.AddColorStopRGBA(stopColorConvertAlpha(1.0, colornames.Red, 0.5))
	cr.SetSource(p)
	cr.Arc(cx, cy, radius, 0, math.Pi*2)
	cr.Fill()
	cr.SetLineWidth(3)
	cr.SetSourceRGB(colorConvert(colornames.Darkslateblue))
	cr.Arc(cx, cy, radius, 0, math.Pi*2)
	cr.Stroke()
	drawGraduations(cr, cx, cy, radius)
}

// main draw fct when canvas need redraw
func drawClock(cr *cairo.Context) {

	//cr.SetSourceRGBA(255, 255, 255, 0)
	//cr.Paint()
	hour, min, sec := lastTime.Clock()
	yy, mm, dd := lastTime.Date()
	hangle := (float64(hour)+float64(min)/90.0)*(math.Pi/6.0) - math.Pi/2
	mangle := float64(min)*(math.Pi/30.0) - math.Pi/2
	sangle := float64(sec)*(math.Pi/30.0) - math.Pi/2
	maxradius := radius - 4
	cx := wfx / 2
	cy := wfy / 2

	// draw clock face
	drawFace(cr, cx, cy, maxradius)

	// draw second
	drawSneddle(cr, cx, cy, maxradius-10.0, 2.0, colornames.Chocolate, sangle)

	// draw min
	drawHMneddle(cr, cx, cy, maxradius-(maxradius/3), 3.0, colornames.Darkgoldenrod, mangle)

	// draw hour
	drawHMneddle(cr, cx, cy, maxradius-(maxradius/2), 2.0, colornames.Firebrick, hangle)
	cr.SetSourceRGB(0, 0, 0)

	//draw center point
	cr.Arc(cx, cy, maxradius/14, 0, math.Pi*2)
	cr.Fill()

	// draw digital clock inside the clock
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
	cr.MoveTo(cx-tdw, cy-maxradius/6)
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
	//win.SetDecorated(false)
	// define the popp menu
	popupmenu, _ := gtk.MenuNew()
	psetAlarm, _ := gtk.MenuItemNew()
	psetAlarm.SetLabel("Set Alarm")
	pswitchDeco, _ := gtk.MenuItemNew()
	pswitchDeco.SetLabel("Switch decoration")
	pabout, _ := gtk.MenuItemNew()
	pabout.SetLabel("About")
	pquit, _ := gtk.MenuItemNew()
	pquit.SetLabel("Quit")
	psep, _ := gtk.SeparatorMenuItemNew()
	popupmenu.Append(psetAlarm)
	popupmenu.Append(pswitchDeco)
	popupmenu.Append(psep)
	popupmenu.Append(pabout)
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
	pabout.Connect("activate", func() {
		setAboutDlg()
	})
	psetAlarm.Connect("activate", func() {
		fmt.Println("Activate Setalarm")
		setAlarmDlg()
	})
	pquit.Connect("activate", func() {
		gtk.MainQuit()
	})
	pswitchDeco.Connect("activate", func() {
		/// try some hack to handle the transparency problem after change deco!
		/// seem to be a GTK BUG
		if win.GetDecorated() {
			win.SetDecorated(false)
			win.SetOpacity(0.9)
			clkcanvas.SetOpacity(0.8)
			win.ShowAll()

		} else {
			win.SetDecorated(true)
			win.SetOpacity(0.9)
			clkcanvas.SetOpacity(0.8)
			win.ShowAll()
		}
	})

	win.Connect("check_resize", func(win *gtk.Window) {
		wfx, wfy, radius = getWinSize(win)
	})

	clkcanvas.Connect("draw", func(da *gtk.DrawingArea, cr *cairo.Context) {
		drawClock(cr)
	})
	win.Connect("button-press-event", func(win *gtk.Window, ev *gdk.Event) {
		//fmt.Println("Event", ev)
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
