package main

import (
	"bufio"
	"flag"
	"os"
	"regexp"

	log "github.com/Crosse/gosimplelogger"
)

var globalInstall bool

func main() {
	var (
		fonts    []string
		filename = flag.String("fromFile", "", "text file containing fonts to install")
		debug    = flag.Bool("debug", false, "Enable debug logging")
		dryrun   = flag.Bool("dry-run", false, "Don't actually download or install anything")
		global   = flag.Bool("g", false, "Install fonts for all users (requires privileges)")
	)

	flag.Parse()

	globalInstall = *global
	if globalInstall {
		FontsDir = GlobalFontsDir
	}

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

		re := regexp.MustCompile(`^(#.*|\s*)?$`)

		scanner := bufio.NewScanner(fd)
		for scanner.Scan() {
			line := scanner.Text()
			if !re.MatchString(line) {
				fonts = append(fonts, line)
			}
		}

		if err = scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}

	fonts = append(fonts, flag.Args()...)

	for _, v := range fonts {
		if *dryrun {
			log.Infof("Would install font(s) from %v", v)
			continue
		}

		log.Debugf("Installing font from %v", v)

		if err := InstallFont(v); err != nil {
			log.Error(err)
		}
	}

	log.Infof("Installed %v fonts", installedFonts)
}
