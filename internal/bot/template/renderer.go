package template

import (
	"bytes"
	"embed"
	"fmt"

	"github.com/CloudyKit/jet/v6"
	"github.com/CloudyKit/jet/v6/loaders/embedfs"
)

//go:embed command/*.jet task/*.jet
var templatesFS embed.FS

var templates = newTemplateSet()

func Render(name string, variables jet.VarMap) (string, error) {
	tmpl, err := templates.GetTemplate(name)
	if err != nil {
		return "", fmt.Errorf("unable to get template %q: %w", name, err)
	}

	out := new(bytes.Buffer)
	if err := tmpl.Execute(out, variables, nil); err != nil {
		return "", fmt.Errorf("unable to execute template %q: %w", name, err)
	}

	return out.String(), nil
}

func newTemplateSet() *jet.Set {
	set := jet.NewSet(
		embedfs.NewLoader(".", templatesFS),
		jet.WithSafeWriter(nil),
	)

	set.AddGlobal("formatMoney", func(value any) string {
		return fmt.Sprintf("%.2f₽", value)
	})

	set.AddGlobal("formatMoneyDelta", func(value any) string {
		return fmt.Sprintf("%+.2f₽", value)
	})

	set.AddGlobal("formatDistance", func(value any) string {
		return fmt.Sprintf("%dкм", value)
	})

	set.AddGlobal("formatFuelConsumption", func(value any) string {
		return fmt.Sprintf("%.2fл/100км", value)
	})

	set.AddGlobal("formatLiters", func(value any) string {
		return fmt.Sprintf("%.2fл", value)
	})

	return set
}
