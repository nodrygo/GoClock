package main

import (
	"fmt"
	_ "math"
	_ "time"

	"github.com/gotk3/gotk3/gtk"
)

// create About DLG
func setAboutDlg() {
	about, _ := gtk.AboutDialogNew()
	about.Activate()
	about.SetComments("Go Gtk CLOCK DEMO\nset not done")
	about.SetVersion("0.1")
	about.SetName("GoClock")
	about.SetCopyright("MIT")
	about.AddCreditSection("CLOCK", []string{"nodryo"})
	about.SetModal(true)
	about.ShowNow()
	about.Run()
	about.Destroy()
}
func setAlarmDlg() {
	fmt.Println("SetAlarmDLG in own pkg")

}
