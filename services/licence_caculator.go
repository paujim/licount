package services

import (
	"strings"

	"github.com/paujim/licount/entities"
)

type LicenceCalculator interface {
	Calculate(applicationID string) (int, error)
}

type licenceCalculator struct {
	scanner CsvScanner
}

func NewLicenceCalculator(s CsvScanner) LicenceCalculator {
	return &licenceCalculator{scanner: s}
}

func (lc *licenceCalculator) Calculate(applicationID string) (int, error) {
	dataCn := lc.scanner.ProduceByApplicationId(applicationID)
	return consumer(dataCn)
}

func consumer(input <-chan *entities.Dto) (int, error) {
	groupByUser := map[string][]*entities.License{}
	for dto := range input {
		item, err := dto.Data, dto.Error
		if err != nil {
			return -1, err
		}
		groupByUser[item.UserID] = append(groupByUser[item.UserID], item)
	}

	var sum int
	for _, computers := range groupByUser {
		sum += countDistinct(computers)
	}
	return sum, nil

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
