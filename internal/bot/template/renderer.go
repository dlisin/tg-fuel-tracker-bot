package template

import (
	"bytes"
	"embed"
	"fmt"
	"reflect"
	"strings"

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

	set.AddGlobal("formatYearMoney", func(value any) string {
		return formatNumber(value, 0) + " ₽"
	})

	set.AddGlobal("formatYearPrice", func(value any) string {
		return formatNumber(value, 2) + " ₽"
	})

	set.AddGlobal("formatYearPercent", func(value any) string {
		return formatNumber(value, 1) + "%"
	})

	set.AddGlobal("formatYearDistance", func(value any) string {
		return formatNumber(value, 0) + " км"
	})

	set.AddGlobal("formatYearLiters", func(value any) string {
		return formatNumber(value, 1) + " л"
	})

	set.AddGlobal("formatYearDays", func(value any) string {
		return formatNumber(value, 0)
	})

	return set
}

func formatNumber(value any, precision int) string {
	number, ok := numericValue(value)
	if !ok {
		return fmt.Sprint(value)
	}

	parts := strings.SplitN(fmt.Sprintf("%.*f", precision, number), ".", 2)
	integer := groupDigits(parts[0])
	if precision == 0 {
		return integer
	}

	return integer + "," + parts[1]
}

func numericValue(value any) (float64, bool) {
	if value == nil {
		return 0, false
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(v.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(v.Uint()), true
	case reflect.Float32, reflect.Float64:
		return v.Float(), true
	default:
		return 0, false
	}
}

func groupDigits(value string) string {
	sign := ""
	if strings.HasPrefix(value, "-") {
		sign = "-"
		value = strings.TrimPrefix(value, "-")
	}

	if len(value) <= 3 {
		return sign + value
	}

	firstGroupLength := len(value) % 3
	if firstGroupLength == 0 {
		firstGroupLength = 3
	}

	var result strings.Builder
	result.WriteString(sign)
	result.WriteString(value[:firstGroupLength])

	for i := firstGroupLength; i < len(value); i += 3 {
		result.WriteByte(' ')
		result.WriteString(value[i : i+3])
	}

	return result.String()
}
