package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"runtime"

	log "github.com/Crosse/gosimplelogger"
)

func main() {
	var filename = flag.String("j", "", "JSON file describing fonts to install")
	/*
		var url = flag.String("url", "", "URL of font file to download and install")
		var localFile = flag.String("file", "", "Local font file to install")
		var isZip = flag.Bool("zip", false, "File or URL is a ZIP file")
	*/
	var debug = flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()

	if len(*filename) == 0 {
		log.Fatal("No filename!")
	}

	if *debug {
		log.LogLevel = log.LogDebug
	} else {
		log.LogLevel = log.LogInfo
	}

	// Read in a JSON-encoded file containing the names and URLs for
	// fonts the user wants to install.
	var fonts = Fonts{}

	log.Debugf("Reading fonts from file %v", *filename)
	f, err := ioutil.ReadFile(*filename)
	if err != nil {
		log.Fatalf("while reading file: %v", err)
	}
	if err = json.Unmarshal(f, &fonts); err != nil {
		log.Fatalf("while unmarshalling JSON: %v", err)
	}
	for _, v := range fonts {
		log.Debugf("Downloading and installing font %v from %v", v.Name, v.URL)
		if err = v.Install(); err != nil {
			log.Error(err)
		}
	}
	if runtime.GOOS == "windows" {
		log.Info("You will need to logoff and logon before the installed fonts will be available.")
	}
}
