package prettyweather

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/rivo/tview"
	"github.com/wtfutil/wtf/view"
	"github.com/wtfutil/wtf/wtf"

	_ "github.com/schachmat/wego/backends"
	wefe "github.com/schachmat/wego/frontends"
	"github.com/schachmat/wego/iface"
)

type Widget struct {
	view.TextWidget

	result   string
	settings *Settings
}

func NewWidget(app *tview.Application, settings *Settings) *Widget {
	widget := Widget{
		TextWidget: view.NewTextWidget(app, settings.common),

		settings: settings,
	}

	return &widget
}

func (widget *Widget) Refresh() {

	widget.prettyWeather2()

	widget.Redraw(func() (string, string, bool) { return widget.CommonSettings().Title, widget.result, false })
}

func (widget *Widget) prettyWeather2() {

	be, ok := iface.AllBackends[widget.settings.backend]
	if !ok {
		log.Fatalf("Could not find selected backend \"%s\"", widget.settings.backend)
	}
	be.Setup()

	fe, _ := &wefe.aatConfig{}
	fe.Setup()

	r := be.Fetch(widget.settings.city, 0)

	unit := iface.UnitsMetric
	if widget.settings.unit == "imperial" {
		unit = iface.UnitsImperial
	} else if widget.settings.unit == "si" {
		unit = iface.UnitsSi
	} else if widget.settings.unit == "metric-ms" {
		unit = iface.UnitsMetricMs
	}

	out := fe.formatCond(make([]string, 5), r.Current, true)

	widget.result = strings.TrimSpace(wtf.ASCIItoTviewColors(strings.Join(out, "\n")))

}

//this method reads the config and calls wttr.in for pretty weather
func (widget *Widget) prettyWeather() {
	client := &http.Client{}

	city := widget.settings.city
	unit := widget.settings.unit
	view := widget.settings.view

	req, err := http.NewRequest("GET", "https://wttr.in/"+city+"?"+view+"?"+unit, nil)
	if err != nil {
		widget.result = err.Error()
		return
	}

	req.Header.Set("Accept-Language", widget.settings.language)
	req.Header.Set("User-Agent", "curl")
	response, err := client.Do(req)
	if err != nil {
		widget.result = err.Error()
		return

	}
	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		widget.result = err.Error()
		return
	}

	widget.result = strings.TrimSpace(wtf.ASCIItoTviewColors(string(contents)))
}
