package cluster

import (
	"strings"

	v1 "github.com/rancher/rancher/pkg/apis/catalog.cattle.io/v1"
)

type App struct {
	Name     string
	Obj      *v1.App
	Migrated bool
	Diff     bool
}

func newApp(obj v1.App) (*App, bool) {
	if strings.HasPrefix(obj.Name, "rancher-") {
		return nil, true
	}
	return &App{
		Name:     obj.Name,
		Obj:      obj.DeepCopy(),
		Migrated: false,
		Diff:     false,
	}, false
}

func (a *App) normalize() {
}

func (a *App) mutate() {
	a.Obj.SetFinalizers(nil)
	a.Obj.SetResourceVersion("")
}
