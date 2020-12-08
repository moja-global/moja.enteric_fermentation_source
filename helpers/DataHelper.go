package helpers

import (
	model "github.com/MullionGroup/go-website-flintpro-example/models"
	"strconv"
)

func ValidateAndPreProcessData(data model.DataPackage, endYear int)(model.DataPackage, error){

	newRecordCount := len(data.AnimalNumberData)
	newAnimalNumberItem := []model.AnimalNumberItem{}
	animalNumbersMap := make(map[string] model.AnimalNumberItem)
	for _,record := range data.AnimalNumberData{
		yearStr := strconv.Itoa(record.Year)
		monthStr := strconv.Itoa(record.Month)
		locationStr := strconv.Itoa(record.Locationid)
		systemStr := strconv.Itoa(record.Systemid)
		animalClassStr := strconv.Itoa(record.AnimalClassid)
		//To check if we have duplicated record
		if (model.AnimalNumberItem{}) == animalNumbersMap[yearStr + "*" + monthStr + "*" + locationStr + "*" + systemStr + "*" + animalClassStr]{
			animalNumbersMap[yearStr + "*" + monthStr + "*" + locationStr + "*" + systemStr + "*" + animalClassStr] = record
		} else{
			errorRecord := model.ErrorRecordItem{
				Recordid: 	   record.Id,
				Locationid:    record.Locationid,
				Location:      record.Location,
				Systemid:      record.Systemid,
				System:        record.System,
				AnimalClassid: record.AnimalClassid,
				AnimalClass:   record.AnimalClass,
				Year:          record.Year,
				Month:         record.Month,
			}
			errorRecord.ErrorMsg = "Duplicated animal number data"
			data.ErrorRecordData = append(data.ErrorRecordData, errorRecord)
		}
		//To populate by month and year
		if record.Month == 0{
			for month := 1; month <= 12; month++{
				newRecordCount+=1
				copyRecord := model.AnimalNumberItem{
					Id:            newRecordCount,
					Locationid:    record.Locationid,
					Location:      record.Location,
					Systemid:      record.Systemid,
					System:        record.System,
					AnimalClassid: record.AnimalClassid,
					AnimalClass:   record.AnimalClass,
					Year:          record.Year,
					Month:         month,
					AnimalNumber:  record.AnimalNumber,
				}
				newAnimalNumberItem = append(newAnimalNumberItem, copyRecord)
			}
		} else{
			newAnimalNumberItem = append(newAnimalNumberItem, record)
		}
	}
	data.AnimalNumberData = newAnimalNumberItem

	//We already sorted this, so the first combination we get for each class will definitely be the first one
	defaultValueList := make(map[string]model.EntericFermEFParameterItem)
	entericFermEFParametersMap := make(map[string]model.EntericFermEFParameterItem)
	for _,record := range data.EntericFermEFParameterData{
		if record.Month == 0{
			// This is to auto-populate month 0 to month 1 - 12 in the later steo
			record.Month = 1
		}
		yearStr := strconv.Itoa(record.Year)
		monthStr := strconv.Itoa(record.Month)
		locationStr := strconv.Itoa(record.Locationid)
		systemStr := strconv.Itoa(record.Systemid)
		animalClassStr := strconv.Itoa(record.AnimalClassid)
		// Not all animal group starts from the same year, so we have a year as identifier too
		// We also only want the closest replacement to the missing data, so we have year inside but not month because we can just use the first representative of the yearly record
		if (model.EntericFermEFParameterItem{}) == defaultValueList[yearStr + "*" + locationStr + "*" + systemStr + "*" + animalClassStr]{
			defaultValueList[yearStr + "*" + locationStr + "*" + systemStr + "*" + animalClassStr] = record
		}
		if (model.EntericFermEFParameterItem{}) == entericFermEFParametersMap[yearStr + "*" + monthStr + "*" + locationStr + "*" + systemStr + "*" + animalClassStr]{
			entericFermEFParametersMap[yearStr + "*" + monthStr + "*" + locationStr + "*" + systemStr + "*" + animalClassStr] = record
		} else{
			errorRecord := model.ErrorRecordItem{
				Recordid: 	   record.Id,
				Locationid:    record.Locationid,
				Location:      record.Location,
				Systemid:      record.Systemid,
				System:        record.System,
				AnimalClassid: record.AnimalClassid,
				AnimalClass:   record.AnimalClass,
				Year:          record.Year,
				Month:         record.Month,
			}
			errorRecord.ErrorMsg = "Duplicated enteric ferm EF param data"
			data.ErrorRecordData = append(data.ErrorRecordData, errorRecord)
		}
	}

	dataCount := len(data.EntericFermEFParameterData)
	for _,record := range defaultValueList{
		month := record.Month
		year := record.Year
		locationStr := strconv.Itoa(record.Locationid)
		systemStr := strconv.Itoa(record.Systemid)
		animalClassStr := strconv.Itoa(record.AnimalClassid)
		replacement := record
		for year <= endYear && month <= 12{
			yearStr := strconv.Itoa(year)
			monthStr := strconv.Itoa(month)
			if (model.EntericFermEFParameterItem{}) == entericFermEFParametersMap[yearStr + "*" + monthStr + "*" + locationStr + "*" + systemStr + "*" + animalClassStr]{
				dataCount+=1
				newRecord := model.EntericFermEFParameterItem{
					Id:                             dataCount,
					Locationid:                     replacement.Locationid,
					Location:                       replacement.Location,
					Systemid:                       replacement.Systemid,
					System:                         replacement.System,
					AnimalClassid:                  replacement.AnimalClassid,
					AnimalClass:                    replacement.AnimalClass,
					Year:                           year,
					Month:                          month,
					BodyWeight:                     replacement.BodyWeight,
					MatureWeight:                   replacement.MatureWeight,
					DailyWeightGain:                replacement.DailyWeightGain,
					FractionOfMonthAlive:           replacement.FractionOfMonthAlive,
					CF:                             replacement.CF,
					C:                              replacement.C,
					CA:                             replacement.CA,
					MilkProd:                       replacement.MilkProd,
					FatContent:                     replacement.FatContent,
					CPregnancy:                     replacement.CPregnancy,
					ProportionAnimalClassPregnant:  replacement.ProportionAnimalClassPregnant,
					ProportionAnimalClassLactating: replacement.ProportionAnimalClassLactating,
					FractionOfMonthLactating:       replacement.FractionOfMonthLactating,
					HoursWorked:                    replacement.HoursWorked,
					DE:                             replacement.DE,
					YM:                             replacement.YM,
				}
				data.EntericFermEFParameterData = append(data.EntericFermEFParameterData, newRecord)
				entericFermEFParametersMap[yearStr + "*" + monthStr + "*" + locationStr + "*" + systemStr + "*" + animalClassStr] = newRecord
			} else{
				replacement = entericFermEFParametersMap[yearStr + "*" + monthStr + "*" + locationStr + "*" + systemStr + "*" + animalClassStr]
			}
			month+=1
			if month > 12{
				month = 1
				year+=1
			}
		}
	}

	// Remove duplicated temperature record
	defaultValueList2  := make(map[string]model.TemperatureLocationItem)
	temperatureLocationsMap := make(map[string]model.TemperatureLocationItem)
	for _,record := range data.TemperatureLocationData{
		if record.Month == 0{
			// This is to auto-populate month 0 to month 1 - 12 in the later steo
			record.Month = 1
		}
		yearStr := strconv.Itoa(record.Year)
		monthStr := strconv.Itoa(record.Month)
		locationStr := strconv.Itoa(record.Locationid)
		// We should have temperature all year round 24/7
		if (model.TemperatureLocationItem{}) == defaultValueList2[yearStr + "*" + monthStr + "*" + locationStr]{
			defaultValueList2[yearStr + "*" + monthStr + "*" + locationStr] = record
		}
		if (model.TemperatureLocationItem{}) == temperatureLocationsMap[yearStr + "*" + monthStr + "*" + locationStr]{
			temperatureLocationsMap[yearStr + "*" + monthStr + "*" + locationStr] = record
		} else{
			errorRecord := model.ErrorRecordItem{
				Recordid: 	   record.Id,
				Locationid:    record.Locationid,
				Location:      record.Location,
				Year:          record.Year,
				Month:         record.Month,
			}
			errorRecord.ErrorMsg = "Duplicated temperature location data"
			data.ErrorRecordData = append(data.ErrorRecordData, errorRecord)
		}
	}

	dataCount2 := len(data.TemperatureLocationData)
	for _,record := range defaultValueList2 {
		month := record.Month
		year := record.Year
		locationStr := strconv.Itoa(record.Locationid)
		replacement := record
		for year <= endYear && month <= 12{
			yearStr := strconv.Itoa(year)
			monthStr := strconv.Itoa(month)
			if (model.TemperatureLocationItem{}) == temperatureLocationsMap[yearStr + "*" + monthStr + "*" + locationStr]{
				dataCount2+=1
				newRecord := model.TemperatureLocationItem{
					Id:         dataCount2,
					Locationid: replacement.Locationid,
					Location:   replacement.Location,
					Year:       year,
					Month:      month,
					AvgTemp:    replacement.AvgTemp,
				}
				data.TemperatureLocationData = append(data.TemperatureLocationData, newRecord)
				temperatureLocationsMap[yearStr + "*" + monthStr + "*" + locationStr] = newRecord
			} else{
				replacement = temperatureLocationsMap[yearStr + "*" + monthStr + "*" + locationStr]
			}
			month+=1
			if month > 12{
				month = 1
				year+=1
			}
		}
	}
	return data, nil
}

