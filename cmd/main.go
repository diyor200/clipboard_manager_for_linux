package main

import (
	"context"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"golang.design/x/clipboard"
)

//import (
//	"clipboard/internal/cache"
//	"context"
//	"fmt"
//	"time"
//
//	"golang.design/x/clipboard"
//)
//
//func init() {
//	err := clipboard.Init()
//	if err != nil {
//		panic(err)
//	}
//}
//
//func main() {
//	//println(string(data))
//	go func() {
//		var q []cache.Queue
//		for {
//			time.Sleep(1 * time.Second)
//			q = cache.ReadAll(context.TODO())
//			if len(q) > 0 {
//				fmt.Println("cache ===>", string(q[len(q)-1].Data))
//			}
//		}
//	}()
//
//	ch := clipboard.Watch(context.TODO(), clipboard.FmtText)
//	for data := range ch {
//		cache.Put(context.TODO(), cache.Queue{Data: data, Kind: clipboard.FmtText})
//	}
//}

func main() {
	// Initialize clipboard
	if err := clipboard.Init(); err != nil {
		log.Fatalf("Failed to initialize clipboard: %v", err)
	}

	// Create Fyne app and window
	myApp := app.New()
	myWindow := myApp.NewWindow("Clipboard Manager")
	myWindow.Resize(fyne.NewSize(400, 600))

	// Clipboard history slice
	clipboardHistory := []string{}

	// List widget to display clipboard history
	historyList := widget.NewList(
		func() int {
			return len(clipboardHistory)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Clipboard Entry")
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*widget.Label).SetText(clipboardHistory[id])
		},
	)

	// Function to monitor clipboard changes
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		lastClipboard := ""
		for {
			select {
			case <-ctx.Done():
				return
			default:
				// Read clipboard data
				text := string(clipboard.Read(clipboard.FmtText))
				if text != "" && text != lastClipboard {
					// Add to history if new
					clipboardHistory = append([]string{text}, clipboardHistory...)
					lastClipboard = text
					historyList.Refresh()
				}
			}
			time.Sleep(500 * time.Millisecond) // Polling interval
		}
	}()

	// Handle list item selection (copy item back to clipboard)
	historyList.OnSelected = func(id widget.ListItemID) {
		if id >= 0 && id < len(clipboardHistory) {
			text := clipboardHistory[id]
			clipboard.Write(clipboard.FmtText, []byte(text))
			log.Println("Copied to clipboard:", text)
		}
	}

	// Layout
	content := container.NewBorder(nil, nil, nil, nil, historyList)
	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}
