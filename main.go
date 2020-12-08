package main

import (
	"flag"
	"fmt"
	common "github.com/MullionGroup/go-website-flintpro-example/common"
	"github.com/MullionGroup/go-website-flintpro-example/helpers"
	model "github.com/MullionGroup/go-website-flintpro-example/models"
	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/twinj/uuid"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
)

var router = mux.NewRouter()
var sensitiveBox = packr.New("Sensitive", "./sensitive")
func layoutFiles() []string {
	files, err := filepath.Glob("templates/*.html")
	if err != nil {
		panic(err)
	}
	return files
}

func ExitErrorMessage(flagset...*flag.FlagSet) {
	fmt.Println("Expected: 'help', 'webserver' or 'cli' command")
	for _, val := range flagset {
		val.Usage()
	}
}

func ShowUsage(flagset...*flag.FlagSet) {
	fmt.Println("Allowed commands: 'help', 'webserver' or 'cli'")
	for _, val := range flagset {
		val.Usage()
	}
}

func init() {
	uuid.SwitchFormat(uuid.FormatHex)

	common.TemplateFiles = layoutFiles()

	helpCmd 			:= flag.NewFlagSet("help", flag.ExitOnError)
	helpSubject			:= helpCmd.String("subject", "all", "Subject to give help on")

	webserverCmd 		:= flag.NewFlagSet("webserver", flag.ExitOnError)
	webserverPort 		:= webserverCmd.Int("port", 8080, "Port for webservice to listen on")
	debugMode			:= webserverCmd.Bool("debug", false, "Run in debug mode")

	cliCmd 				:= flag.NewFlagSet("cli", flag.ExitOnError)
	cliInputFile		:= cliCmd.String("input_file", "", "Input file name (.json or .db or .sqlite3")
	cliRunSimulation	:= cliCmd.Bool("run_simulation", false, "to start the simulation")
	cliOutputFile		:= cliCmd.String("output_file", "", "Output file name (.json or .db or .sqlite3")

	//SetCustomUsage(webserverCmd, cliCmd)
	var isWebserver = false
	var isCli = false
	runtime.GOMAXPROCS(runtime.NumCPU())

	if len(os.Args) < 2 {
		isWebserver = true
	} else{
		switch os.Args[1] {
		case "webserver":
			webserverCmd.Parse(os.Args[2:])
			isWebserver = true
		case "cli":
			cliCmd.Parse(os.Args[2:])
			isCli = true
		case "help":
			helpCmd.Parse(os.Args[2:])
			switch *helpSubject {
			case "all":
				ShowUsage(helpCmd, webserverCmd, cliCmd)
			case "webserver":
				webserverCmd.Usage()
			case "cli":
				cliCmd.Usage()
			}
		default:
			ExitErrorMessage(webserverCmd, cliCmd)
			os.Exit(1)
		}
	}

	// We box the sensitive file as a binary and only put it back in during runtime
	sensitiveFile, err := sensitiveBox.Find("ef_credentials.json")
	if err != nil {
		logrus.Errorf("Error with parsing Templates in sensitive files")
	}
	// Global settings
	//flag.Parse()

	// Help Settings
	viper.SetDefault("subject", 				*helpSubject)

	// WebServer Settings
	viper.SetDefault("webserver", 				isWebserver)
	viper.SetDefault("webserver_port", 		*webserverPort)
	viper.SetDefault("debug",	 				*debugMode)

	// CLI Settings
	viper.SetDefault("cli", 					isCli)
	viper.SetDefault("input_file", 			*cliInputFile)
	viper.SetDefault("run_simulation", 		*cliRunSimulation)
	viper.SetDefault("output_file", 			*cliOutputFile)
	viper.SetDefault("credentials_file",sensitiveFile)
	viper.SetDefault("redirect_uri","http://localhost:8080/login_google")
	viper.SetDefault("redirect_uri_data","http://localhost:8080/index")
	viper.SetDefault("read_me","README.md")

	viper.AutomaticEnv()

	err = helpers.InitializeOAuthGoogle()
	if err != nil {
		logrus.Errorf("Error with parsing secret file in sensitive files")
	}
}


//TODO: allow export data in different format in cli (.db,.json)
//TODO: run simulation
//TODO: show about us
func main() {
	if viper.GetBool("webserver") {
		common.StartWebServer(viper.GetString("webserver_port"))
		handleSigterm(func() {
			fmt.Println("Captured Ctrl+C")
		})
	} else if viper.GetBool("cli")  {
		// Run a simulation directly
		// This method will load the same Google Sheet info, but from a JSON file containing the same info
		// Data sources could be SQLite, JSON, GoogleSheets, csv, etc...
		fmt.Printf("Loading input file (%v - %v)\n", viper.GetString("input_type"), viper.GetString("input_file"))

		filename := viper.GetString("input_file")
		var data model.DataPackage
		var err error
		if len(filename) > 0{
			extension := filepath.Ext(filename)
			if extension!=".db" && extension!=".sqlite3" && extension!=".json"{
				// This is an error
				fmt.Printf("Input file not a valid filetype (%v) - needed .json, .db or .sqlite3 - exiting\n", extension)
				os.Exit(1)
			}

			switch extension {
			case ".json":
				handleSigterm(func() {
					fmt.Println("Captured Ctrl+C")
				})
				data, err := helpers.LoadJSONIntoDataPackage(model.DataPackage{}, filename)
				if err != nil {
					// This is an error
					fmt.Printf("Error (%v) - exiting\n", err)
					os.Exit(1)
				}
				model.OutputAllTables(data)
			case ".sqlite3":
				handleSigterm(func() {
					fmt.Println("Captured Ctrl+C")
					helpers.DBClient.Close()
				})
				helpers.DBClient = &helpers.GormClient{}
				data, err := helpers.LoadSQLiteIntoDataPackage(model.DataPackage{},filename)
				if err != nil {
					// This is an error
					fmt.Printf("Error (%v) - exiting\n", err)
					os.Exit(1)
				}
				model.OutputAllTables(data)
			case ".db":
				handleSigterm(func() {
					fmt.Println("Captured Ctrl+C")
					helpers.DBClient.Close()
				})
				helpers.DBClient = &helpers.GormClient{}
				data, err := helpers.LoadSQLiteIntoDataPackage(model.DataPackage{}, filename)
				if err != nil {
					// This is an error
					fmt.Printf("Error (%v) - exiting\n", err)
					os.Exit(1)
				}
				model.OutputAllTables(data)
			default:
				// This is an error
				fmt.Printf("Input file not a valid type (%v) - exiting\n", viper.GetString("input_file"))
				os.Exit(1)
			}
		} else{
			// This is an error
			fmt.Printf("Input file not found")
			os.Exit(1)
		}

		runSimulation := viper.GetBool("run_simulation")
		if runSimulation == true{
			data, err = helpers.RunSimulation(data)
			if err != nil {
				// This is an error
				fmt.Printf("Error (%v) - exiting\n", err)
				os.Exit(1)
			}
		}

		outputFilename := viper.GetString("output_file")
		if len(outputFilename) > 0{
			outputExtension := filepath.Ext(outputFilename)
			var dataByte []byte
			var err	error
			var newFilename string
			switch outputExtension {
			case ".json":
				handleSigterm(func() {
					fmt.Println("Captured Ctrl+C")
				})
				dataByte, err = helpers.LoadDataPackageIntoJSON(data)
				if err != nil {
					// This is an error
					fmt.Printf("Error (%v) - exiting\n", err)
					os.Exit(1)
				}
				newFilename = data.FileName+".json"
			case ".sqlite3":
				handleSigterm(func() {
					fmt.Println("Captured Ctrl+C")
					helpers.DBClient.Close()
				})
				dataByte, err = helpers.LoadDataPackageIntoSQLite(data)
				if err != nil {
					// This is an error
					fmt.Printf("Error (%v) - exiting\n", err)
					os.Exit(1)
				}
				newFilename = data.FileName+".sqlite3"
			case ".db":
				handleSigterm(func() {
					fmt.Println("Captured Ctrl+C")
					helpers.DBClient.Close()
				})
				dataByte, err = helpers.LoadDataPackageIntoSQLite(data)
				if err != nil {
					// This is an error
					fmt.Printf("Error (%v) - exiting\n", err)
					os.Exit(1)
				}
				newFilename = data.FileName+".sqlite3"
			default:
				// This is an error
				fmt.Printf("Input file not a valid type (%v) - exiting\n", viper.GetString("input_file"))
				os.Exit(1)
			}
			err = ioutil.WriteFile(newFilename, dataByte, 0644)
			if err != nil {
				// This is an error
				fmt.Printf("Error (%v) - exiting\n", err)
				os.Exit(1)
			}
			model.OutputAllTables(data)
		}
	} else {
		// Assume help or error was handled already
	}
}

// Handles Ctrl+C or most other means of "controlled" shutdown gracefully. Invokes the supplied func before exiting.
func handleSigterm(handleExit func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		handleExit()
		os.Exit(1)
	}()
}
