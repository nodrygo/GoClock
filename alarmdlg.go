package main

import (
	"fmt"
	_ "math"
	"strconv"
	_ "time"

	"github.com/gotk3/gotk3/gtk"
)

type Alarm struct {
	hour string
	min  string
	//	days      []string
	activated bool
}

// check
func (al *Alarm) checkAlarm(hour, min, second int) bool {
	hh, _ := strconv.ParseInt(alarm.hour, 10, 0)
	mm, _ := strconv.ParseInt(alarm.min, 10, 0)
	if (int(hh) == hour) && (int(mm) == min) && second == 0 {
		return true
	}
	return false
}

// create new button with callback
func createButton(label string, t *gtk.Entry, maxval int64) *gtk.Button {
	b, _ := gtk.ButtonNew()
	b.SetLabel(label)
	b.Connect("clicked", func() {
		txt, _ := t.GetText()
		if label == "UP" {
			t.SetText(inc(txt, maxval))
		} else {
			t.SetText(dec(txt, maxval))
		}
	})
	return b
}

// cyclic increment 0 to max
func inc(txt string, max int64) string {
	val, _ := strconv.ParseInt(txt, 10, 64)
	if val < max {
		val++
	} else {
		val = 0
	}
	return fmt.Sprintf("%02d", val)
}

// cyclic decrement 0 to max
func dec(txt string, max int64) string {
	val, _ := strconv.ParseInt(txt, 10, 64)
	if val > 0 {
		val--
	} else {
		val = max
	}
	return fmt.Sprintf("%02d", val)
}

// create txt
func createEntry(txt string) *gtk.Entry {
	e, _ := gtk.EntryNew()
	e.SetMaxLength(2)
	e.SetText("00")
	e.SetWidthChars(2)
	return e
}

// set Alarm Dlg button up down for H M activate
// days select
func (al *Alarm) openAlarmDlg(win *gtk.Window) {
	// Dialog init/destroy
	dlg, _ := gtk.DialogNew()
	dlg.SetHExpand(true)
	dlg.SetVExpand(true)
	dlg.SetParent(win)
	defer dlg.Destroy()

	grid, _ := gtk.GridNew()
	grid.SetOrientation(gtk.ORIENTATION_VERTICAL)
	dca, _ := dlg.GetContentArea()
	dca.Add(grid)

	// my own spinner
	thour := createEntry("HOUR")
	thour.SetEditable(false)
	tmin := createEntry("MIN")
	tmin.SetEditable(false)
	lblH, _ := gtk.LabelNew("HOUR")
	lblM, _ := gtk.LabelNew("MIN")
	//MAX SPINNER
	grid.Attach(lblH, 0, 0, 1, 1)
	grid.Attach(createButton("UP", thour, 23), 1, 0, 1, 1)
	grid.Attach(thour, 2, 0, 1, 1)
	grid.Attach(createButton("DOWN", thour, 23), 3, 0, 1, 1)
	// MIN SPINNER
	grid.Attach(lblM, 0, 1, 1, 1)
	grid.Attach(createButton("UP", tmin, 59), 1, 1, 1, 1)
	grid.Attach(tmin, 2, 1, 1, 1)
	grid.Attach(createButton("DOWN", tmin, 59), 3, 1, 1, 1)

	dlg.AddButton("CANCEL", gtk.RESPONSE_CANCEL)
	dlg.AddButton("OK", gtk.RESPONSE_ACCEPT)

	dlg.SetModal(true)
	dlg.ShowAll()
	resp := dlg.Run()
	hh, _ := thour.GetText()
	mm, _ := tmin.GetText()

	if resp == gtk.RESPONSE_ACCEPT {
		alarm.activated = true
	} else {
		alarm.activated = false
	}
	alarm.hour = hh
	alarm.min = mm

	fmt.Printf("HOUR:%s  MIN:%s\n", hh, mm)
	fmt.Println("ALARM STAT IS", alarm.activated)
}

// create About DLG
func setAboutDlg() {
	about, _ := gtk.AboutDialogNew()
	about.Activate()
	about.SetComments("Go Gtk(gotk3) CLOCK DEMO\nset not done")
	about.SetVersion("0.1")
	about.SetName("GoClock")
	about.SetCopyright("License MIT")
	about.AddCreditSection("", []string{"nodryo"})
	about.SetModal(true)
	about.ShowNow()
	about.Run()
	about.Destroy()
}
