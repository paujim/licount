package services

import (
	"encoding/csv"
	"io"
	"log"

	"github.com/paujim/licount/entities"
)

const (
	ComputerID = iota
	UserID
	ApplicationID
	ComputerType
	Comment
)

type CsvScanner interface {
	Scan() (*entities.License, error)
	ProduceByApplicationId(applicationID string) chan *entities.Dto
}

type csvScanner struct {
	reader *csv.Reader
}

func NewScanner(o io.Reader) CsvScanner {
	csv_o := csv.NewReader(o)
	csv_o.Read() // Scanning the headers, we don't use them
	return &csvScanner{reader: csv_o}
}

func (o *csvScanner) Scan() (*entities.License, error) {
	return parse(o.reader.Read())
}

func (o *csvScanner) ProduceByApplicationId(applicationID string) chan *entities.Dto {
	outputCn := make(chan *entities.Dto)
	go func(output chan *entities.Dto) {
		l, err := o.Scan()
		for err == nil {
			// Filter by ApplicationID
			if applicationID == l.ApplicationID {
				output <- &entities.Dto{Data: l, Error: nil}
			}
			l, err = o.Scan()
		}
		if err.Error() != "EOF" {
			log.Printf("errors: %s\n", err)
			output <- &entities.Dto{Data: nil, Error: err}
		}
		close(output)
	}(outputCn)
	return outputCn
}

func parse(raw []string, err error) (*entities.License, error) {
	if err != nil {
		return nil, err
	}
	return &entities.License{
		ComputerID:    raw[ComputerID],
		UserID:        raw[UserID],
		ApplicationID: raw[ApplicationID],
		ComputerType:  raw[ComputerType],
		Comment:       raw[Comment]}, nil
}
