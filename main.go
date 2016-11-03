package main

import (
	"bufio"
	"flag"
	"os"
	"runtime"

	log "github.com/Crosse/gosimplelogger"
)

func main() {
	var fonts []string

	var filename = flag.String("fromFile", "", "text file containing fonts to install")
	var debug = flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()

	if *filename == "" && len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	if *debug {
		log.LogLevel = log.LogDebug
	} else {
		log.LogLevel = log.LogInfo
	}

	if *filename != "" {
		fd, err := os.Open(*filename)
		if err != nil {
			log.Fatal(err)
		}
		scanner := bufio.NewScanner(fd)
		for scanner.Scan() {
			fonts = append(fonts, scanner.Text())
		}
		if err = scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}

	for _, v := range flag.Args() {
		fonts = append(fonts, v)
	}

	for _, v := range fonts {
		log.Debugf("Installing font from %v", v)
		if err := InstallFont(v); err != nil {
			log.Error(err)
		}
	}

	if runtime.GOOS == "windows" {
		log.Info("You will need to logoff and logon before the installed font(s) will be available.")
	}
}
