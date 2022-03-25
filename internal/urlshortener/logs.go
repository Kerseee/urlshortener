package urlshortener

// logError logs err via app.Logger.
func (app *App) logError(err error) {
	app.logger.Println(err.Error())
}

// logInfo logs msg via app.Logger
func (app *App) logInfo(msg string) {
	app.logger.Println(msg)
}
