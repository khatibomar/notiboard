package main

// Window configuration
const (
	DefaultWindowWidth  = 1920
	DefaultWindowHeight = 1080
)

type ScreenConfig struct {
	Width      int32
	Height     int32
	Fullscreen bool
	Resizable  bool
}

var (
	HD      = ScreenConfig{Width: 1280, Height: 720, Fullscreen: false, Resizable: true}
	FullHD  = ScreenConfig{Width: 1920, Height: 1080, Fullscreen: false, Resizable: true}
	QHD     = ScreenConfig{Width: 2560, Height: 1440, Fullscreen: false, Resizable: true}
	UHD     = ScreenConfig{Width: 3840, Height: 2160, Fullscreen: false, Resizable: true}
	Desktop = ScreenConfig{Width: 0, Height: 0, Fullscreen: true, Resizable: false} // Uses desktop resolution
)
