package services

import (
	"encoding/csv"
	"io"

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
	ProduceByApplicationId(applicationID string) chan *entities.License
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

func (o *csvScanner) ProduceByApplicationId(applicationID string) chan *entities.License {
	link := make(chan *entities.License)
	go func() {
		l, err := o.Scan()
		for err == nil {
			if applicationID == l.ApplicationID {
				link <- l
			}
			l, err = o.Scan()
		}
		close(link)
	}()
	return link
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
