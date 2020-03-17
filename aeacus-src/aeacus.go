package main

import (
	"fmt"
	"log"
	"os"
    "runtime"
	"github.com/urfave/cli"
)

//////////////////////////////////////////////////////////////////
//  .oooo.    .ooooo.   .oooo.    .ooooo.  oooo  oooo   .oooo.o //
// `P  )88b  d88' `88b `P  )88b  d88' `"Y8 `888  `888  d88  "8  //
//  .oP"888  888ooo888  .oP"888  888        888   888  `"Y88b.  //
// d8(  888  888    .o d8(  888  888   .o8  888   888  o.  )88b //
// `Y888""8o `Y8bod8P' `Y888""8o `Y8bod8P'  `V88V"V8P' 8""888P' //
//////////////////////////////////////////////////////////////////

type metaConfig struct {
Cli        *cli.Context
TeamID  string
	ConfigName string
	DataName   string
	WebPath   string
	Config     scoringChecks
}

func main() {

    var configName string
    var dataName string
    var webName string
    if runtime.GOOS == "linux" {
        configName = "/opt/aeacus/scoring.conf"
    	dataName = "/opt/aeacus/scoring.dat"
    	webName = "/opt/aeacus/web/"
    } else if runtime.GOOS == "windows" {
        configName = "C:\\aeacus\\scoring.conf"
    	dataName = "C:\\aeacus\\scoring.dat"
    	webName = "C:\\aeacus\\web\\"
    } else {
        failPrint("This operating system (" + runtime.GOOS + ") is not supported!")
        os.Exit(1)
    }

    // read TeamID
    teamID := "B emoji"

    id := imageData{0, 0, 0, []scoreItem{}, 0, []scoreItem{}, 0, 0}

	app := &cli.App{
		UseShortOptionHandling: true,
		EnableBashCompletion:   true,
		Name:                   "aeacus",
		Usage:                  "setup and score vulnerabilities in an image",
		Action: func(c *cli.Context) error {
			mc := metaConfig{c, teamID, configName, dataName, webName, scoringChecks{}}
            checkConfig(&mc)
			scoreImage(&mc, &id)
			return nil
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "Print extra information",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "score",
				Aliases: []string{"s"},
				Usage:   "(default) Score image with current scoring config",
				Action: func(c *cli.Context) error {
					mc := metaConfig{c, teamID, configName, dataName, webName, scoringChecks{}}
                    checkConfig(&mc)
					scoreImage(&mc, &id)
					return nil
				},
			},
			{
				Name:    "simulate",
				Aliases: []string{"i"},
				Usage:   "Score image with current scoring data",
				Action: func(c *cli.Context) error {
					mc := metaConfig{c, teamID, configName, dataName, webName, scoringChecks{}}
                    parseConfig(&mc, readData(&mc))
					scoreImage(&mc, &id)
					return nil
				},
			},
            {
				Name:    "check",
				Aliases: []string{"c"},
				Usage:   "Check that the scoring config is valid",
				Action: func(c *cli.Context) error {
					mc := metaConfig{c, teamID, configName, dataName, webName, scoringChecks{}}
					checkConfig(&mc)
					return nil
				},
			},
			{
				Name:    "encrypt",
				Aliases: []string{"e"},
				Usage:   "Encrypt scoring.conf to scoring.dat",
				Action: func(c *cli.Context) error {
					mc := metaConfig{c, teamID, configName, dataName, webName, scoringChecks{}}
					writeConfig(&mc)
					return nil
				},
			},
		//	{
		//		Name:    "decrypt",
		//		Aliases: []string{"d"},
		//		Usage:   "Encrypt scoring.conf to scoring.dat",
		//		Action: func(c *cli.Context) error {
		//			mc := metaConfig{c, teamID, configName, dataName, webName, scoringChecks{}}
  //                  fmt.Println(readData(&mc))
		//			return nil
		//		},
		//	},
			{
				Name:    "createfqs",
				Aliases: []string{"f"},
				Usage:   "Create forensic question files (3 by default)",
				Action: func(c *cli.Context) error {
                    fmt.Println("todo")
					return nil
				},
			},
			{
				Name:    "release",
				Aliases: []string{"r"},
				Usage:   "Prepare the image for release",
				Action: func(c *cli.Context) error {
					mc := metaConfig{c, teamID, configName, dataName, webName, scoringChecks{}}
					releaseImage(&mc)
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

///////////////////////
// CONTROL FUNCTIONS //
///////////////////////

func scoreImage(mc *metaConfig, id *imageData) {
    connStatus := []string{"green", "OK", "green", "OK", "green", "OK"}
    if mc.Config.Remote != "" {
        connStatus, connection := checkServer(mc)
        if ! connection {
            failPrint("Can't access remote scoring server!")
            genReport(mc, id, connStatus)
            os.Exit(1)
        }
    }
    if runtime.GOOS == "linux" {
        scoreLinux(mc, id)
    } else {
        //scoreWindows(mc, id)
        fmt.Println("score wondows")
    }
    genReport(mc, id, connStatus)
}

func checkConfig(mc *metaConfig) {
    fileContent, err := readFile(mc.ConfigName)
    if err != nil {
        failPrint("Configuration file not found!")
        os.Exit(1)
    }
	parseConfig(mc, fileContent)
	if mc.Cli.Bool("v") {
		printConfig(mc)
	}
}

func releaseImage(mc *metaConfig) {
    checkConfig(mc)
	writeConfig(mc)
    genReadMe(mc)
	warnPrint("The rest of this doesn't actually do anything yet. Just pretend like it does lol")
	cleanUp(mc)
    writeDesktopFiles(mc)
    installService(mc)
	// add self to services
	// set up notifications

}
