package services

import (
	"strings"

	"github.com/paujim/licount/entities"
)

type LicenceCalculator interface {
	Calculate(applicationID string) int
	CalculateV2(applicationID string) int
}

type licenceCalculator struct {
	scanner CsvScanner
}

func NewLicenceCalculator(s CsvScanner) LicenceCalculator {
	return &licenceCalculator{scanner: s}
}

func (lc *licenceCalculator) Calculate(applicationID string) int {
	data := []*entities.License{}
	l, err := lc.scanner.Scan()
	for err == nil {
		data = append(data, l)
		l, err = lc.scanner.Scan()
	}
	return filterProcess(data, applicationID)
}

func (lc *licenceCalculator) CalculateV2(applicationID string) int {
	link := lc.scanner.ProduceByApplicationId(applicationID)
	sumChan := consumer(link)
	sum := <-sumChan
	return sum
}

func consumer(link <-chan *entities.License) chan int {
	total := make(chan int)
	go func() {
		groupByUser := map[string][]*entities.License{}
		for item := range link {
			groupByUser[item.UserID] = append(groupByUser[item.UserID], item)
		}
		var sum int
		for _, computers := range groupByUser {
			sum += countDistinct(computers)
		}
		total <- sum
	}()
	return total
}

func filterProcess(data []*entities.License, applicationID string) int {
	groupByUser := map[string][]*entities.License{}
	for _, item := range data {
		// FilterBy ApplicationID
		if item.ApplicationID == applicationID {
			// GroupBy UserID
			groupByUser[item.UserID] = append(groupByUser[item.UserID], item)
		}
	}
	var sum int
	for _, computers := range groupByUser {
		sum += countDistinct(computers)
	}
	return sum
}

func countDistinct(groupedByUser []*entities.License) int {
	uniqueDesktop := map[string]bool{}
	uniqueLaptop := map[string]bool{}
	for _, item := range groupedByUser {
		compType := strings.ToUpper(item.ComputerType)
		key := item.UserID + "-" + item.ComputerID
		if compType == "DESKTOP" {
			uniqueDesktop[key] = true

		} else {
			uniqueLaptop[key] = true
		}
	}
	totalDesktop := len(uniqueDesktop)
	totalLaptop := len(uniqueLaptop)
	if totalDesktop > totalLaptop {
		return totalDesktop
	}
	return totalLaptop
}
