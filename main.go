package main

import (
	_ "host-editor/internal/logic"

	"embed"
	"host-editor/internal/view"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// 应用启动
	app := view.NewApp()
	err := wails.Run(&options.App{
		Title:     "Host Editor",
		Width:     1024,
		Height:    768,
		MinWidth:  640,
		MinHeight: 480,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		DisableResize:    false,
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		Mac: &mac.Options{
			TitleBar:    mac.TitleBarHiddenInset(),
			DisableZoom: false,
			Preferences: &mac.Preferences{
				FullscreenEnabled: mac.Enabled,
			},
		},
		OnStartup:  app.Startup,
		OnShutdown: app.Shutdown,
		Bind:       view.GetBind(),
	})
	if err != nil {
		println("Error:", err.Error())
	}
}
