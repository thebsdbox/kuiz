package ui

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// ErrorMessage - displays an error modal dialog
func ErrorMessage(err error) error {
	app := tview.NewApplication()
	button := tview.NewButton(err.Error()).SetSelectedFunc(func() {
		app.Stop()
	})
	button.SetBackgroundColor(tcell.ColorRed)
	button.SetBorder(true).SetRect(5, 5, 75, 3)
	if err := app.SetRoot(button, false).Run(); err != nil {
		return err
	}
	return nil
}
