package services

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLicenseCalculator_Calculate(t *testing.T) {
	tests := []struct {
		data  string
		total int
	}{
		{`ComputerID,UserID,ApplicationID,ComputerType,Comment 
1,1,374,LAPTOP,Exported from System A 
2,1,374,DESKTOP,Exported from System A 
3,2,374,DESKTOP,Exported from System A 
4,2,374,DESKTOP,Exported from System A
`, 3},
		{`ComputerID,UserID,ApplicationID,ComputerType,Comment 
1,1,374,LAPTOP,Exported from System A 
2,2,374,DESKTOP,Exported from System A 
2,2,374,desktop,Exported from System B `, 2},
	}

	for _, tc := range tests {
		lc := NewLicenceCalculator(NewScanner(strings.NewReader(tc.data)))
		assert := assert.New(t)
		total := lc.CalculateV2("374")
		assert.Equal(tc.total, total)
	}
}

// func TestLicenseCalculator_SmallSampleFile(t *testing.T) {
// 	f, err := os.OpenFile("../sample-small.csv", os.O_RDONLY, os.ModePerm)
// 	if err != nil {
// 		log.Printf("unable to open file: %s", err)
// 		return
// 	}
// 	defer f.Close()
// 	lc := NewLicenceCalculator(NewScanner(f))
// 	assert := assert.New(t)
// 	total := lc.CalculateV2("700")
// 	assert.Equal(231, total)
// }

// func TestLicenseCalculator_LargeSampleFile(t *testing.T) {
// 	f, err := os.OpenFile("../sample-large.csv", os.O_RDONLY, os.ModePerm)
// 	if err != nil {
// 		log.Printf("unable to open file: %s", err)
// 		return
// 	}
// 	defer f.Close()
// 	lc := NewLicenceCalculator(NewScanner(f))
// 	assert := assert.New(t)
// 	total := lc.CalculateV2("700")
// 	assert.Equal(15082, total)
// }
