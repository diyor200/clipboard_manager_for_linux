package main

import (
	"context"
	"log"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"golang.design/x/clipboard"
)

// https://snapcraft.io/docs/registering-your-app-name
func main() {
	// Initialize clipboard
	if err := clipboard.Init(); err != nil {
		log.Fatalf("Failed to initialize clipboard: %v", err)
	}

	// Create Fyne app and window
	myApp := app.New()
	myWindow := myApp.NewWindow("Clipboard Manager")

	iconData := loadIconData()
	if iconData != nil {
		appIcon := fyne.NewStaticResource("icon.png", iconData)
		myApp.SetIcon(appIcon)
	}

	// Disable window resizing
	myWindow.SetFixedSize(true)
	myWindow.Resize(fyne.NewSize(500, 600))

	// Clipboard history
	clipboardHistory := []string{}

	// Container to hold clipboard items
	historyContainer := container.NewVBox()

	// Function to truncate long text for card previews
	truncateText := func(text string, maxLength int) string {
		if len(text) > maxLength {
			return text[:maxLength] + "..."
		}
		return text
	}

	// Function to create a new clipboard card
	createClipboardCard := func(text string) fyne.CanvasObject {
		// Card layout
		preview := widget.NewLabel(truncateText(text, 50))

		// View button to show full content in a dialog
		viewButton := widget.NewButton("View", func() {
			dialog.ShowCustom(
				"Clipboard Entry",
				"Close",
				widget.NewLabel(text),
				myWindow,
			)
		})

		// Copy button to copy text to clipboard
		copyButton := widget.NewButton("Copy", func() {
			clipboard.Write(clipboard.FmtText, []byte(text))
		})

		// HBox to align preview and buttons
		card := container.NewVBox(
			container.NewHBox(preview, layout.NewSpacer(), viewButton, copyButton), // Align buttons to the end
			widget.NewSeparator(),
		)
		return card
	}

	// Function to refresh the history view
	refreshHistory := func() {
		historyContainer.Objects = nil
		for _, entry := range clipboardHistory {
			historyContainer.Add(createClipboardCard(entry))
		}
		historyContainer.Refresh()
	}

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
					refreshHistory()
				}
			}
			time.Sleep(500 * time.Millisecond) // Polling interval
		}
	}()

	// Clear all button
	clearButton := widget.NewButton("Clear All", func() {
		clipboardHistory = nil
		refreshHistory()
	})

	// Layout
	scrollableHistory := container.NewScroll(historyContainer)
	content := container.NewBorder(nil, clearButton, nil, nil, scrollableHistory)
	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}

// Helper function to load icon data
func loadIconData() []byte {
	iconData, err := os.ReadFile("/usr/share/icons/clipboard.png")
	if err != nil {
		log.Fatalf("Failed to load icon: %v", err)
	}
	log.Println("Icon loaded successfully") // Add this line to confirm the icon is loaded
	return iconData
}
