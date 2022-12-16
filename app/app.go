package app

import (
	"context"
	"fmt"
	"github.com/rayguo17/go-socks/manager"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
func (a *App) ListUser() []*manager.Display {
	resp := GetAllUser()
	if resp.GetErrCode() != 0 {
		return []*manager.Display{}
	} else {
		return resp.GetData().([]*manager.Display)
	}
}
