package main

import (
	"context"
	"fmt"
	"math"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// drawConnectionIndicator draws an animated connection status indicator
func drawConnectionIndicator(x, y float32, info *ConnectionInfo, animTime float64, isHovered bool) {
	status := info.GetStatus()

	baseRadius := float32(20)
	if isHovered {
		baseRadius = 22 // Slightly bigger on hover
	}

	var primaryColor, accentColor rl.Color
	var statusText string

	switch status {
	case Connected:
		primaryColor = rl.Color{R: 46, G: 204, B: 113, A: 255} // Emerald green
		accentColor = rl.Color{R: 39, G: 174, B: 96, A: 255}   // Darker green
		statusText = "ONLINE"
	case Disconnected:
		primaryColor = rl.Color{R: 231, G: 76, B: 60, A: 255} // Red
		accentColor = rl.Color{R: 192, G: 57, B: 43, A: 255}  // Darker red
		statusText = "OFFLINE"
	case Reconnecting:
		primaryColor = rl.Color{R: 255, G: 193, B: 7, A: 255} // Orange/Yellow
		accentColor = rl.Color{R: 243, G: 156, B: 18, A: 255} // Darker orange
		statusText = "RECONNECTING"
	default: // Unknown/Starting
		primaryColor = rl.Color{R: 149, G: 165, B: 166, A: 255} // Gray
		accentColor = rl.Color{R: 127, G: 140, B: 141, A: 255}  // Darker gray
		statusText = "CHECKING"
	}

	pulseScale := float32(1.0)
	if status == Connected {
		pulseScale = 1.0 + 0.1*float32(math.Sin(animTime*3))
	}

	if isHovered {
		hoverGlow := baseRadius + 12
		hoverColor := primaryColor
		hoverColor.A = 40
		rl.DrawCircle(int32(x), int32(y), hoverGlow, hoverColor)
	}

	if status == Connected {
		glowRadius := baseRadius + 8 + 3*float32(math.Sin(animTime*2))
		glowColor := primaryColor
		glowColor.A = 60
		rl.DrawCircle(int32(x), int32(y), glowRadius, glowColor)
	}

	rl.DrawCircle(int32(x), int32(y), baseRadius*pulseScale, primaryColor)

	innerRadius := baseRadius * 0.6 * pulseScale
	rl.DrawCircle(int32(x), int32(y), innerRadius, accentColor)

	symbolSize := float32(8) * pulseScale
	switch status {
	case Connected:
		rl.DrawLineEx(
			rl.Vector2{X: x - symbolSize/2, Y: y},
			rl.Vector2{X: x - symbolSize/4, Y: y + symbolSize/2},
			3, rl.White,
		)
		rl.DrawLineEx(
			rl.Vector2{X: x - symbolSize/4, Y: y + symbolSize/2},
			rl.Vector2{X: x + symbolSize/2, Y: y - symbolSize/2},
			3, rl.White,
		)
	case Disconnected:
		rl.DrawLineEx(
			rl.Vector2{X: x - symbolSize/2, Y: y - symbolSize/2},
			rl.Vector2{X: x + symbolSize/2, Y: y + symbolSize/2},
			3, rl.White,
		)
		rl.DrawLineEx(
			rl.Vector2{X: x + symbolSize/2, Y: y - symbolSize/2},
			rl.Vector2{X: x - symbolSize/2, Y: y + symbolSize/2},
			3, rl.White,
		)
	default:
		for i := range 3 {
			angle := animTime*4 + float64(i)*2.094 // 120 degrees apart
			dotX := x + float32(math.Cos(angle))*symbolSize/2
			dotY := y + float32(math.Sin(angle))*symbolSize/2
			rl.DrawCircle(int32(dotX), int32(dotY), 2, rl.White)
		}
	}

	textWidth := rl.MeasureText(statusText, 16)
	textX := int32(x) - textWidth/2
	textY := int32(y) + int32(baseRadius) + 10
	rl.DrawText(statusText, textX, textY, 16, primaryColor)
}

// drawHoverTooltip draws a tooltip when hovering over the connection indicator
func drawHoverTooltip(x, y float32, info *ConnectionInfo) {
	status, lastError, lastPingTime := info.GetStatus(), info.GetLastError(), info.GetLastPingTime()

	tooltipWidth := int32(300)
	tooltipHeight := int32(100)
	tooltipX := int32(x) - tooltipWidth - 20 // Position to the left of the indicator
	tooltipY := int32(y) - tooltipHeight/2

	// Ensure tooltip stays on screen
	if tooltipX < 10 {
		tooltipX = int32(x) + 30 // Position to the right instead
	}
	if tooltipY < 10 {
		tooltipY = 10
	}

	// Draw tooltip background
	rl.DrawRectangle(tooltipX, tooltipY, tooltipWidth, tooltipHeight, rl.Color{R: 50, G: 50, B: 50, A: 240})
	rl.DrawRectangleLines(tooltipX, tooltipY, tooltipWidth, tooltipHeight, rl.Gray)

	// Draw tooltip content
	yOffset := tooltipY + 10
	fontSize := int32(14)
	lineHeight := int32(18)

	statusText := "Status: Unknown"
	switch status {
	case Connected:
		statusText = "Status: Connected ✓"
	case Disconnected:
		statusText = "Status: Disconnected ✗"
	case Reconnecting:
		statusText = "Status: Reconnecting..."
	}
	rl.DrawText(statusText, tooltipX+10, yOffset, fontSize, rl.White)
	yOffset += lineHeight

	if !lastPingTime.IsZero() {
		timeText := "Last Check: " + lastPingTime.Format("15:04:05")
		rl.DrawText(timeText, tooltipX+10, yOffset, fontSize, rl.LightGray)
		yOffset += lineHeight
	}

	if lastError != nil {
		errorText := "Error: " + lastError.Error()
		if len(errorText) > 35 {
			errorText = errorText[:35] + "..."
		}
		rl.DrawText(errorText, tooltipX+10, yOffset, fontSize, rl.Color{R: 255, G: 100, B: 100, A: 255})
	} else if status == 1 {
		rl.DrawText("All systems operational", tooltipX+10, yOffset, fontSize, rl.Color{R: 100, G: 255, B: 100, A: 255})
	}

	rl.DrawText("Click for details", tooltipX+10, tooltipY+tooltipHeight-20, 12, rl.Color{R: 200, G: 200, B: 200, A: 255})
}

// drawDetailWindow draws a detailed window when clicked
func drawDetailWindow(info *ConnectionInfo, mousePos rl.Vector2) (bool, bool) {
	status, lastError, lastPingTime, reconnectTime, reconnectFailed := info.GetStatus(), info.GetLastError(), info.GetLastPingTime(), info.GetReconnectTime(), info.GetAllReconnectTriesFailed()

	windowWidth := int32(500)
	windowHeight := int32(300)
	windowX := (1920 - windowWidth) / 2
	windowY := (1080 - windowHeight) / 2

	rl.DrawRectangle(windowX, windowY, windowWidth, windowHeight, rl.Color{R: 40, G: 40, B: 40, A: 250})
	rl.DrawRectangleLines(windowX, windowY, windowWidth, windowHeight, rl.Color{R: 100, G: 100, B: 100, A: 255})

	var shouldClose bool
	var shouldReconnect bool

	rl.DrawRectangle(windowX, windowY, windowWidth, 40, rl.Color{R: 60, G: 60, B: 60, A: 255})
	rl.DrawText("Connection Details", windowX+15, windowY+12, 18, rl.White)

	closeX := windowX + windowWidth - 35
	closeY := windowY + 8
	closeButtonHovered := isPointInRect(mousePos.X, mousePos.Y, closeX, closeY, 24, 24)
	closeButtonColor := rl.Color{R: 200, G: 50, B: 50, A: 255}
	if closeButtonHovered {
		closeButtonColor = rl.Color{R: 255, G: 70, B: 70, A: 255}
	}
	rl.DrawRectangle(closeX, closeY, 24, 24, closeButtonColor)
	rl.DrawText("X", closeX+8, closeY+4, 16, rl.White)

	if closeButtonHovered && rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		shouldClose = true
	}

	contentY := windowY + 50
	yOffset := contentY + 20
	lineHeight := int32(25)

	statusText := "Connection Status: "
	var statusColor rl.Color
	switch status {
	case Connected:
		statusText += "Connected ✓"
		statusColor = rl.Green
	case Disconnected:
		statusText += "Disconnected ✗"
		statusColor = rl.Red
	case Reconnecting:
		statusText += "Reconnecting..."
		statusColor = rl.Orange
	default:
		statusText += "Unknown"
		statusColor = rl.Gray
	}
	rl.DrawText(statusText, windowX+20, yOffset, 16, statusColor)
	yOffset += lineHeight

	dbText := sensoredConnString()
	rl.DrawText(dbText, windowX+20, yOffset, 14, rl.LightGray)
	yOffset += lineHeight

	if !lastPingTime.IsZero() {
		timeText := "Last Ping: " + lastPingTime.Format("2006-01-02 15:04:05")
		rl.DrawText(timeText, windowX+20, yOffset, 14, rl.LightGray)
		yOffset += lineHeight
	}

	if !reconnectTime.IsZero() {
		reconnectText := "Last Reconnect: " + reconnectTime.Format("2006-01-02 15:04:05")
		rl.DrawText(reconnectText, windowX+20, yOffset, 14, rl.LightGray)
		yOffset += lineHeight
	}

	yOffset += 10

	if lastError != nil {
		rl.DrawText("Error Details:", windowX+20, yOffset, 16, rl.Red)
		yOffset += lineHeight

		errorMsg := lastError.Error()
		maxWidth := 460
		if rl.MeasureText(errorMsg, 14) > int32(maxWidth) {
			words := []rune(errorMsg)
			line := ""
			for i, char := range words {
				line += string(char)
				if rl.MeasureText(line, 14) > int32(maxWidth) || i == len(words)-1 {
					rl.DrawText(line, windowX+20, yOffset, 14, rl.Color{R: 255, G: 150, B: 150, A: 255})
					yOffset += 20
					line = ""
					if yOffset > windowY+windowHeight-40 {
						break
					}
				}
			}
		} else {
			rl.DrawText(errorMsg, windowX+20, yOffset, 14, rl.Color{R: 255, G: 150, B: 150, A: 255})
		}
	} else if status == 1 {
		rl.DrawText("✓ All systems operational", windowX+20, yOffset, 16, rl.Green)
		yOffset += lineHeight
		rl.DrawText("Database connection is healthy and responsive.", windowX+20, yOffset, 14, rl.LightGray)
	}

	// Reconnect button (show when disconnected and auto-reconnect failed)
	if status == Disconnected && reconnectFailed {
		buttonWidth := int32(140)
		buttonHeight := int32(35)
		buttonX := windowX + (windowWidth-buttonWidth)/2 // Center the button
		buttonY := windowY + windowHeight - 80

		buttonHovered := isPointInRect(mousePos.X, mousePos.Y, buttonX, buttonY, buttonWidth, buttonHeight)
		buttonColor := rl.Color{R: 46, G: 204, B: 113, A: 255}
		if buttonHovered {
			buttonColor = rl.Color{R: 39, G: 174, B: 96, A: 255}
		}

		rl.DrawRectangle(buttonX+2, buttonY+2, buttonWidth, buttonHeight, rl.Color{R: 0, G: 0, B: 0, A: 100})

		rl.DrawRectangle(buttonX, buttonY, buttonWidth, buttonHeight, buttonColor)
		rl.DrawRectangleLines(buttonX, buttonY, buttonWidth, buttonHeight, rl.White)

		buttonText := "Manual Reconnect"
		textWidth := rl.MeasureText(buttonText, 14)
		textX := buttonX + (buttonWidth-textWidth)/2
		textY := buttonY + 10
		rl.DrawText(buttonText, textX, textY, 14, rl.White)

		if buttonHovered && rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			shouldReconnect = true
		}
	}

	rl.DrawText("Click anywhere outside this window to close", windowX+20, windowY+windowHeight-25, 12, rl.Color{R: 150, G: 150, B: 150, A: 255})

	return shouldClose, shouldReconnect
}

// isPointInCircle checks if a point is inside a circle
func isPointInCircle(px, py, cx, cy, radius float32) bool {
	dx := px - cx
	dy := py - cy
	return dx*dx+dy*dy <= radius*radius
}

// isPointInRect checks if a point is inside a rectangle
func isPointInRect(px, py float32, x, y, width, height int32) bool {
	return px >= float32(x) && px <= float32(x+width) && py >= float32(y) && py <= float32(y+height)
}

func main() {
	ctx := context.Background()
	dbManager := NewDatabaseManager(ctx)
	err := dbManager.Connect()
	if err != nil {
		fmt.Printf("Initial database connection failed: %v\n", err)
		fmt.Println("Application will start with disconnected status")
	}
	defer dbManager.Close()

	windowTitle := "notify board"

	rl.InitWindow(1920, 1080, windowTitle)
	defer rl.CloseWindow()

	targetFPS := int32(60)
	rl.SetTargetFPS(targetFPS)

	connectionInfo := &ConnectionInfo{}
	connectionInfo.SetStatus(Unknown)

	var animationTime float64 = 0
	var showDetailWindow bool
	manualReconnectSignal := make(chan struct{}, 1)

	// Connection monitoring and auto-reconnection goroutine
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		maxRetries := 3

		healthcheck := func() {
			consecutiveFailures := 0
			err := dbManager.Ping()
			connectionInfo.SetLastPingTime()
			if err != nil {
				consecutiveFailures++
				connectionInfo.SetStatus(Disconnected)
				connectionInfo.SetLastError(err)

				if consecutiveFailures < maxRetries {
					connectionInfo.SetStatus(Reconnecting)

					for attempt := range maxRetries {
						connectionInfo.SetReconnectTime()
						fmt.Printf("Auto-reconnection attempt %d/%d...\n", attempt+1, maxRetries)

						err := dbManager.Connect()
						if err == nil {
							pingErr := dbManager.Ping()
							if pingErr == nil {
								fmt.Println("Auto-reconnection successful!")
								connectionInfo.SetStatus(Connected)
								connectionInfo.SetLastError(nil)
								connectionInfo.SetAllReconnectTriesFailed(false)
								return
							}
							err = pingErr
						}

						fmt.Printf("Auto-reconnection attempt %d failed: %v\n", attempt+1, err)
						connectionInfo.SetStatus(Reconnecting)
						connectionInfo.SetLastError(err)

						// Exponential backoff
						waitTime := time.Duration(attempt+1) * 2 * time.Second
						time.Sleep(waitTime)
					}

					fmt.Println("All auto-reconnection attempts exhausted")
					connectionInfo.SetAllReconnectTriesFailed(true)
					connectionInfo.SetStatus(Disconnected)
					connectionInfo.SetLastError(err)
				}
			} else {
				if consecutiveFailures > 0 {
					fmt.Println("Connection restored")
				}
				connectionInfo.SetStatus(Connected)
				connectionInfo.SetLastError(nil)
			}
		}

		for {
			select {
			case <-ticker.C:
				if connectionInfo.GetAllReconnectTriesFailed() {
					continue
				}
				healthcheck()
			case <-manualReconnectSignal:
				connectionInfo.SetStatus(Reconnecting)
				healthcheck()
			}
		}
	}()

	for !rl.WindowShouldClose() {
		animationTime += float64(rl.GetFrameTime())

		mousePos := rl.GetMousePosition()
		indicatorX := float32(1820)
		indicatorY := float32(50)

		isHovered := isPointInCircle(mousePos.X, mousePos.Y, indicatorX, indicatorY, 30)

		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			if showDetailWindow {
				shouldClose, shouldReconnect := drawDetailWindow(connectionInfo, mousePos)
				if shouldClose {
					showDetailWindow = false
				}
				if shouldReconnect {
					select {
					case manualReconnectSignal <- struct{}{}:
					default:
					}
				}

				windowWidth := int32(500)
				windowHeight := int32(300)
				windowX := (1920 - windowWidth) / 2
				windowY := (1080 - windowHeight) / 2

				if !shouldClose && !shouldReconnect && !isPointInRect(mousePos.X, mousePos.Y, windowX, windowY, windowWidth, windowHeight) {
					showDetailWindow = false
				}
			} else if isHovered {
				showDetailWindow = true
			}
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		drawConnectionIndicator(indicatorX, indicatorY, connectionInfo, animationTime, isHovered)

		if isHovered && !showDetailWindow {
			drawHoverTooltip(indicatorX, indicatorY, connectionInfo)
		}

		if showDetailWindow {
			drawDetailWindow(connectionInfo, mousePos)
		}

		rl.EndDrawing()
	}
}
