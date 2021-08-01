package main

import (
	"flag"
	"log"
	"os"

	"github.com/paujim/licount/services"
)

func main() {

	filename := flag.String("file", "sample-large.csv", "filename")
	applicationId := flag.String("app-id", "700", "application id")

	flag.Parse()

	if *filename == "" || *applicationId == "" {
		flag.Usage()
		return
	}

	licenseFile, err := os.OpenFile(*filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Printf("unable to open file: %s", err)
		return
	}
	defer licenseFile.Close()
	log.Printf("Calculating for %s\n", *applicationId)
	total, err := services.
		NewLicenceCalculator(services.NewScanner(licenseFile)).
		Calculate(*applicationId)

	if err != nil {
		log.Printf("error:, %s", err)
		return
	}

	log.Printf("The number of liceses for %s is %v", *applicationId, total)

}
