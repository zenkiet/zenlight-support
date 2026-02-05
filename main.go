package main

import (
	"embed"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"zenlight-support/internal/app"
	"zenlight-support/internal/config"
	"zenlight-support/internal/service"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed version.txt
var rawVersion string
var version = strings.TrimSpace(rawVersion)

const uniqueAppID = "com.zensoftware.zenlight-support"

func main() {
	// Set up log output for debugging
	logFile, _ := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	log.SetOutput(logFile)
	defer logFile.Close()

	logger := slog.New(slog.NewJSONHandler(logFile, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	// Load configuration
	cfg, err := config.LoadConfig(version)
	if err != nil {
		slog.Error("Failed to load config:", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Initialize
	srvMgr := service.NewManager()
	a := app.NewApp(*cfg, srvMgr, version)

	// Create application menu
	// appMenu := menu.NewMenu()
	// fileMenu := appMenu.AddSubmenu("File")
	// fileMenu.AddText("Open Dashboard", keys.CmdOrCtrl("o"), func(_ *menu.CallbackData) {
	//  wailsRuntime.WindowShow(a.Ctx)
	//  wailsRuntime.WindowSetAlwaysOnTop(a.Ctx, true)
	//  wailsRuntime.WindowSetAlwaysOnTop(a.Ctx, false)
	// })
	// fileMenu.AddSeparator()
	// fileMenu.AddText("Quit", keys.CmdOrCtrl("q"), func(_ *menu.CallbackData) {
	//  a.ForceQuit()
	//  wailsRuntime.Quit(a.Ctx)
	// })

	// Create application with options
	errApp := wails.Run(&options.App{
		Title:     "Zenlight Support",
		MinWidth:  1050,
		MinHeight: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		WindowStartState: options.Normal,
		DisableResize:    false,
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        a.Startup,
		OnShutdown:       a.Shutdown,
		Bind: []interface{}{
			a,
		},
		// Menu: appMenu,
		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId: uniqueAppID,
			OnSecondInstanceLaunch: func(secondInstanceData options.SecondInstanceData) {
				wailsRuntime.WindowUnminimise(a.Ctx)
				wailsRuntime.WindowShow(a.Ctx)
				wailsRuntime.WindowSetAlwaysOnTop(a.Ctx, true)
			},
		},
		// OnBeforeClose: func(ctx context.Context) (prevent bool) {
		//  if !a.IsQuitting() {
		//    slog.Info("User clicked close, hiding window to background")

		//    wailsRuntime.WindowHide(ctx)
		//    wailsRuntime.EventsEmit(ctx, "notification", "Zenlight Support is running in background")
		//    return true
		//  }
		//  return false
		// },
		Windows: &windows.Options{
			BackdropType:        windows.Mica,
			EnableSwipeGestures: true,
			DisableWindowIcon:   false,
			Theme:               windows.SystemDefault,
		},
		Mac: &mac.Options{
			WindowIsTranslucent: true,
			About: &mac.AboutInfo{
				Title:   "Zenlight Support",
				Message: "Monitor your POS services",
			},
		},
	})
	if errApp != nil {
		println("Error:", errApp.Error())
	}
}
