package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kingpin"
)

const applicationName = "alexaops"
const applicationVersion = "v0.1.0-alpha"

const cultureEnglish = "en"
const cultureGerman = "de"
const defaultCulture = cultureEnglish

const defaultConfigFilePath = "alexaops.conf"

var (
	app = kingpin.New(applicationName, fmt.Sprintf(`「 %s 」%s is your 24/7 endpoint for your "Amazon Echo" based DevOPS skills.

    🌈 https://github.com/andreaskoch/alexa-ops
`, applicationName, applicationVersion))

	// listen
	listenCommand    = app.Command("listen", "Listen for requests from Alexa")
	listenAddress    = listenCommand.Flag("address", "The address/port to listen on").Default(":33011").Envar("ALEXAOPS_LISTEN_ADDRESS").Short('a').String()
	listenConfigPath = listenCommand.Flag("config", "The config file path").Default(defaultConfigFilePath).Envar("ALEXAOPS_CONFIG").Short('c').String()

	// dump sample config
	dumpSampleConfigCommand = app.Command("dump-sample-config", "Dump a sample configuration file to the current directory")
	dumpSampleConfigPath    = dumpSampleConfigCommand.Flag("config", "The config file path").Default(defaultConfigFilePath).Envar("ALEXAOPS_CONFIG").Short('c').String()
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

	case dumpSampleConfigCommand.FullCommand():
		sampleConfig := newSampleConfig()
		saveError := saveConfigToFile(*dumpSampleConfigPath, sampleConfig)
		if saveError != nil {
			fmt.Fprintf(os.Stderr, "Failed to write the sample config to %q: %s\n", *dumpSampleConfigPath, saveError.Error())
			os.Exit(1)
		}

		os.Exit(0)

	case listenCommand.FullCommand():
		config, configError := readConfigFromFile(*listenConfigPath)
		if configError != nil {
			if os.IsNotExist(configError) {
				fmt.Fprintf(os.Stderr, "The config file %q was not found\n", *listenConfigPath)
			} else {
				fmt.Fprintf(os.Stderr, "Failed to read the config file %q: %s\n", *listenConfigPath, configError.Error())
			}

			os.Exit(1)
		}

		intendHandlerProvider := newIntendHandlerProvider(config)
		server, serverError := NewServer(*listenAddress, config, intendHandlerProvider)
		if serverError != nil {
			fmt.Fprintf(os.Stderr, "Failed to create the server: %s\n", serverError.Error())
			os.Exit(1)
		}

		if runError := server.Run(); runError != nil {
			fmt.Fprintf(os.Stderr, "The server returned an error: %s\n", runError)
			os.Exit(1)
		}

		os.Exit(0)
	}

}
