package main

import (
	"embed"
	"github.com/wailsapp/wails/v2/pkg/options/windows"

	//"github.com/wailsapp/wails/v2/pkg/menu"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
)

//go:embed all:frontend/dist
var assets embed.FS

//wails build -ldflags="-H windowsgui -w -s"

func main() {
	// Create an instance of the app structure
	app := NewApp()
	//menu := menu.Menu{}

	// Create application with options
	err := wails.Run(&options.App{
		Title:            "goscrcpy",
		Width:            500,
		Height:           120,
		Assets:           assets,
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		DisableResize:    true,
		AlwaysOnTop:      true,

		//Menu:             &menu,
		Frameless: true,
		Bind: []interface{}{
			app,
		},
		OnBeforeClose: app.onBeforeClose,
		Windows: &windows.Options{
			//WebviewIsTransparent:              false,
			//WindowIsTranslucent:               false,
			//DisableWindowIcon:                 false,
			//DisableFramelessWindowDecorations: false,
			WebviewUserDataPath: "./data",
			//Theme:                             windows.SystemDefault,
			//CustomTheme: &windows.ThemeSettings{
			//	DarkModeTitleBar:   windows.RGB(20, 20, 20),
			//	DarkModeTitleText:  windows.RGB(200, 200, 200),
			//	DarkModeBorder:     windows.RGB(20, 0, 20),
			//	LightModeTitleBar:  windows.RGB(200, 200, 200),
			//	LightModeTitleText: windows.RGB(20, 20, 20),
			//	LightModeBorder:    windows.RGB(200, 200, 200),
			//},
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
