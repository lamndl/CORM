package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"ChessRepertoire/backend"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {

    // app := NewApp(db)
	// // Create application with options
	// // err := wails.Run(&options.App{
	// // 	Title:  "ChessRepertoire",
	// // 	Width:  1024,
	// // 	Height: 768,
	// // 	AssetServer: &assetserver.Options{
	// // 		Assets: assets,
	// // 	},
	// // 	BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
	// // 	OnStartup:        app.startup,
	// // 	Bind: []interface{}{
	// // 		app,
	// // 	},
	// // })

    //    db, err := storage.Open("file:repertoire.db?_foreign_keys=on")
    // if err != nil {
    //     log.Fatal(err)
    // }
    // defer db.Close()

    // init DB
    db, err := backend.Open("file:repertoire.db?_foreign_keys=on")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    // wrap in App
    app := NewApp(db)

    // bind both App and its RepMgr
    if err := wails.Run(&options.App{
        Title:  "Repertoire Manager",
        Width:  1024,
        Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
        Bind:   []interface{}{app, app.RepMgr},
    }); err != nil {
        log.Fatal(err)
    }
}
