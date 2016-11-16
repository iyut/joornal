package main

import (
  "encoding/json"
  "log"
  "os"
  "runtime"
  "github.com/iyut/joornal/app/route"
  "github.com/iyut/joornal/app/lib/jsonparser"
  "github.com/iyut/joornal/app/lib/database"
  "github.com/iyut/joornal/app/lib/email"
  "github.com/iyut/joornal/app/lib/recaptcha"
  "github.com/iyut/joornal/app/lib/server"
  "github.com/iyut/joornal/app/lib/session"
  "github.com/iyut/joornal/app/lib/view"
  "github.com/iyut/joornal/app/lib/view/plugin"
)

// ****************************
// APPLICATION CONFIG
// ****************************

// Config the settings variable
var config = &configuration{}

type configuration struct {
  Database  database.Info   `json:"Database"`
	Email     email.SMTPInfo  `json:"Email"`
	Recaptcha recaptcha.Info  `json:"Recaptcha"`
	Server    server.Server   `json:"Server"`
	Session   session.Session `json:"Session"`
	Template  view.Template   `json:"Template"`
	View      view.View       `json:"View"`
}

func (c *configuration) ParseJSON(b []byte) error{
  return json.Unmarshal(b, &c)
}

// ****************************
// APPLICATION LOGIC
// ****************************

func init(){
  // Verbose logging with filename and line number
  log.SetFlags(log.Lshortfile)

  // use all CPU cores
  runtime.GOMAXPROCS(runtime.NumCPU())
}

func main(){
  // Load the configuration file
  jsonparser.Load("src/config"+ string(os.PathSeparator)+"config.json", config)

  // Session configuration
  session.Configure(config.Session)

  // Database configuration
  database.Connect(config.Database)

  // Configure the google RECAPTCHA prior to loading view plugin
  recaptcha.Configure(config.Recaptcha)

  // Setup the views
	view.Configure(config.View)
	view.LoadTemplates(config.Template.Root, config.Template.Children)
	view.LoadPlugins(
		plugin.TagHelper(config.View),
		plugin.NoEscape(),
		plugin.PrettyTime(),
		recaptcha.Plugin())

	// Start the listener
	server.Run(route.LoadHTTP(), route.LoadHTTPS(), config.Server)

}
