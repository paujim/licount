package services

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCsvScanner_Scan(t *testing.T) {
	testData := `ComputerID,UserID,ApplicationID,ComputerType,Comment 
1,1,374,LAPTOP,Exported from System A 
4,2,375,DESKTOP,Exported from System A`

	csvs := NewScanner(strings.NewReader(testData))
	assert := assert.New(t)
	l, err := csvs.Scan()

	assert.NoError(err)
	assert.NotNil(l)
	assert.Equal("1", l.ComputerID)
	assert.Equal("1", l.UserID)
	assert.Equal("374", l.ApplicationID)
	assert.Equal("LAPTOP", l.ComputerType)

	l, err = csvs.Scan()

	assert.NoError(err)
	assert.NotNil(l)
	assert.Equal("4", l.ComputerID)
	assert.Equal("2", l.UserID)
	assert.Equal("375", l.ApplicationID)
	assert.Equal("DESKTOP", l.ComputerType)

	l, err = csvs.Scan()
	assert.Error(err, "EOF")
	assert.Nil(l)
}

func TestCsvScanner_ScanWithError(t *testing.T) {
	testData := `ComputerID,UserID,ApplicationID,ComputerType,Comment 
4,2,375,DESKTOP`

	csvs := NewScanner(strings.NewReader(testData))
	assert := assert.New(t)
	l, err := csvs.Scan()
	assert.Error(err)
	assert.Nil(l)

}
