package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kingpin"
)

const applicationName = "alexaops"
const applicationVersion = "v0.1.0-alpha"

var (
	app = kingpin.New(applicationName, fmt.Sprintf(`「 %s 」%s is your 24/7 endpoint for your "Amazon Echo" based DevOPS skills.

    🌈 https://github.com/andreaskoch/alexa-ops
`, applicationName, applicationVersion))

	// crawl
	listenCommand = app.Command("listen", "Listen for incoming requests")
	listenAddress = listenCommand.Flag("address", "The address/port to listen on").Default(":33011").Envar("ALEXAOPS_LISTEN_ADDRESS").Short('a').String()
	listenConfig  = listenCommand.Flag("config", "The config file path").Default("alexaops.conf").Envar("ALEXAOPS_CONFIG").Short('c').String()
)

func init() {
	app.Version(applicationVersion)
	app.Author("Andreas Koch <andy@ak7.io>")
}

func main() {
	handleCommandlineArgument(os.Args[1:])
}

func handleCommandlineArgument(arguments []string) {

	switch kingpin.MustParse(app.Parse(arguments)) {

	case listenCommand.FullCommand():
		server, serverError := NewServer()
		if serverError != nil {
			fmt.Fprintf(os.Stderr, "%s", serverError.Error())
			os.Exit(1)
		}

		if runError := server.Run(); runError != nil {
			fmt.Fprintf(os.Stderr, "%s", runError)
			os.Exit(1)
		}

		os.Exit(0)
	}

}
