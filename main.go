package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	// Create a new HTTP client with a 5-second timeout for both the request and response
	client := &http.Client{
		Timeout: time.Second * 5,
	}

	// Create a new GET request
	req, err := http.NewRequest("GET", "https://api.rubyonrails.org/", nil)
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}

	// Send the request using the client
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading HTTP response:", err)
		return
	}

	// Convert the response body to a string and split it into lines
	lines := strings.Split(string(body), "\n")

	// Create a new Tview application
	app := tview.NewApplication()

	// Create a new text view to display the lines
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	// Add the lines to the text view
	for i, line := range lines {
		fmt.Fprintf(textView, "[%d:%d] %s\n", i, len(line), line)
	}

	// Create a new input field for fuzzy searching
	inputField := tview.NewInputField().
		SetLabel("Search").
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetFieldTextColor(tcell.ColorWhite).
		SetLabelColor(tcell.ColorLightGray)

	// Set the input field's changed function to update the text view based on the search input
	inputField.SetChangedFunc(func(search string) {
		// Clear the current text view content
		textView.Clear()

		// Add the lines matching the search input to the text view
		for i, line := range lines {
			if strings.Contains(line, search) {
				fmt.Fprintf(textView, "[%d:%d] %s\n", i, len(line), line)
			}
		}

		// Force the text view to update
		app.Draw()
	})

	// Create a new flex layout and add the input field and text view
	flex := tview.NewFlex().
		AddItem(inputField, 0, 1, false).
		AddItem(textView, 0, 10, true)

	// Set the root of the application to the flex layout
	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
