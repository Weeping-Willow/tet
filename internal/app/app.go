package app

import "context"

type App struct {
}

func MustNew(ctx context.Context) *App {
	return &App{}
}
