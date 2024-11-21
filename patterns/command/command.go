package command

type Command interface {
	Execute()
}

type Application struct {
	Enabled bool
}

func (app *Application) On() {
	app.Enabled = true
}

func (app *Application) Off() {
	app.Enabled = false
}

func (app *Application) Toggle() {
	app.Enabled = !app.Enabled
}

type OnCommand struct {
	app *Application
}

func (o *OnCommand) Execute() {
	o.app.On()
}

type OffCommand struct {
	app *Application
}

func (o *OffCommand) Execute() {
	o.app.Off()
}

type ToggleCommand struct {
	app *Application
}

func (o *ToggleCommand) Execute() {
	o.app.Toggle()
}
