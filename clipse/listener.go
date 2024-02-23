package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/atotto/clipboard"
)

/* runListener is essentially a while loop to be created as a system background process on boot.
   can be stopped at any time with:
   	clipse -kill
   	pkill -f clipse
   	killall clipse
*/

func runListener(historyFilePath, clipsDir, displayServer string, imgEnabled bool) error {
	// Listen for SIGINT (Ctrl+C) and SIGTERM signals to properly close the program
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	// Load existing data from file, if any

	go func() { // go routine necessary to acheive desired CTRL+C behavior
		for {
			// Get the current clipboard content
			input, err := clipboard.ReadAll()
			handleError(err)

			dt := dataType(input)

			switch dt {
			case "text":
				if input != "" && !contains(input) {
					err := addClipboardItem(historyFilePath, input, "null")
					handleError(err)
				}
			case "png", "jpeg":
				if imgEnabled {
					file := fmt.Sprintf("%s.%s", strconv.Itoa(len(input)), dt)
					filePath := filepath.Join(clipsDir, tmpDir, file)
					title := fmt.Sprintf("<BINARY FILE> %s", file)
					if !contains(title) {
						err = saveImage(filePath, displayServer)
						handleError(err)
						err = addClipboardItem(historyFilePath, title, filePath)
						handleError(err)
					}
				}
			}
			time.Sleep(pollInterval) // pollInterval defined in constants.go
		}
	}()

	<-interrupt
	return nil
}
