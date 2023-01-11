package main

import (
	"context"
	"embed"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
	"path"
	goRuntime "runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

var app *App
var api *Api

func main() {

	// Create an instance of the app structure
	app = NewApp()
	api = NewApi()

	AppMenu := menu.NewMenu()
	FileMenu := AppMenu.AddSubmenu("File")
	FileMenu.AddText("Quit", keys.CmdOrCtrl("q"), func(_ *menu.CallbackData) {
		runtime.Quit(app.ctx)
	})
	if goRuntime.GOOS == "darwin" {
		AppMenu.Append(menu.EditMenu()) // on macos platform, we should append EditMenu to enable Cmd+C,Cmd+V,Cmd+Z... shortcut
	}

	logFile := os.Getenv("HOME") + "/.goarchitect/goarchitect.log"
	checkForLoggingFile(logFile)
	fileLogger := logger.NewFileLogger(logFile)

	// Create application with options
	err := wails.Run(&options.App{
		Title:            "Go Architect",
		WindowStartState: options.Maximised,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        onStartup,
		Bind: []interface{}{
			app,
			api,
		},
		Menu:               AppMenu,
		Logger:             fileLogger,
		LogLevel:           logger.DEBUG,
		LogLevelProduction: logger.ERROR,
		Mac: &mac.Options{
			About: &mac.AboutInfo{
				Title:   "My Application",
				Message: "© 2021 Me",
				Icon:    icon,
			},
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

func checkForLoggingFile(filename string) {
	dir := path.Dir(filename)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		_, err = os.Create(filename)
		if err != nil {
			panic(err)
		}
	}
}

func onStartup(ctx context.Context) {
	app.SetContext(ctx)
	api.SetContext(ctx)
}
