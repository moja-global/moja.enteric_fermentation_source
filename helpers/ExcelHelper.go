package helpers

import (
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	model "github.com/MullionGroup/go-website-flintpro-example/models"
	"log"
	"sort"
	"strconv"
)

func LoadExcelIntoDataPackage(data model.DataPackage, filename string) (model.DataPackage, error) {
	data.ClearPackageData()

	f, err := excelize.OpenFile(filename)
	if err != nil {
		return model.DataPackage{}, err
	}

	// Get all the rows in the Settings.
	rows, err := f.GetRows("Settings")
	if err != nil {
		return model.DataPackage{},err
	}
	rows = rows[4:]
	defaultValue := []string{"0",""}
	for _, row := range rows {
		if len(row) <= 0 || len(row[0]) <= 0 || len(row[1]) <= 0{
			continue
		}
		if len(row) < 2 && len(rows) > 0{
			row = append(row, defaultValue[len(row):]...)
		}
		key := row[0]
		value := row[1]
		data.SettingData = append(data.SettingData, model.SettingDataItem{key, value})
	}

	rows, err = f.GetRows("System")
	if err != nil {
		return model.DataPackage{},err
	}
	rows = rows[4:]
	defaultValue = []string{"0",""}
	for _, row := range rows {
		if len(row) <= 0 || len(row[0]) <= 0 || len(row[1]) <= 0{
			continue
		}
		if len(row) < 2 && len(rows) > 0{
			row = append(row, defaultValue[len(row):]...)
		}
		id,_ := strconv.Atoi(row[0])
		name := row[1]
		data.SystemData = append(data.SystemData, model.SystemDataItem{id, name})
	}

	rows, err = f.GetRows("Location")
	if err != nil {
		return model.DataPackage{},err
	}
	rows = rows[4:]
	defaultValue = []string{"0",""}
	for _, row := range rows {
		if len(row) <= 0{
			continue
		}
		if len(row) < 2 && len(rows) > 0{
			row = append(row, defaultValue[len(row):]...)
		}
		id,_ := strconv.Atoi(row[0])
		name := row[1]
		data.LocationData = append(data.LocationData, model.LocationDataItem{id, name})
	}

	rows, err = f.GetRows("AnimalClass")
	if err != nil {
		return model.DataPackage{},err
	}
	rows = rows[4:]
	defaultValue = []string{"0","", "", ""}
	for _, row := range rows {
		if len(row) <= 0 || len(row[0]) <= 0 || len(row[1]) <= 0{
			continue
		}
		if len(row) < 4 && len(rows) > 0{
			row = append(row, defaultValue[len(row):]...)
		}
		id,_ := strconv.Atoi(row[0])
		parentClass := row[1]
		name := row[2]
		defaultEF,_ := strconv.ParseFloat(row[3], 64)
		data.AnimalClassData = append(data.AnimalClassData, model.AnimalClassDataItem{id, parentClass,name, defaultEF})
	}

	systemMap, locationsMap, animalsMap := data.GenerateMappingsForPrimaryTable()

	rows, err = f.GetRows("TemperatureLocation")
	if err != nil {
		return model.DataPackage{},err
	}
	rows = rows[4:]
	defaultValue = []string{"0","","0","0","0.0","0.0"}
	for _, row := range rows {
		if len(row) <= 0 || len(row[0]) <= 0 || len(row[1]) <= 0{
			continue
		}
		if len(row) < 5 && len(rows) > 0{
			row = append(row, defaultValue[len(row):]...)
		}
		id,_ := strconv.Atoi(row[0])
		location := row[1]
		year,_ := strconv.Atoi(row[2])
		month,_ := strconv.Atoi(row[3])
		avgTemp, _ := strconv.ParseFloat(row[4], 64)
		data.TemperatureLocationData = append(data.TemperatureLocationData, model.TemperatureLocationItem{id, locationsMap[location].Id, locationsMap[location].Name, year, month, avgTemp})
	}

	rows, err = f.GetRows("AnimalNumbers")
	if err != nil {
		return model.DataPackage{},err
	}
	rows = rows[4:]
	defaultValue = []string{"0","","","","0","0","0"}
	for _, row := range rows {
		if len(row) <= 0 || len(row[0]) <= 0 || len(row[1]) <= 0{
			continue
		}
		if len(row) < 7 && len(rows) > 0{
			row = append(row, defaultValue[len(row):]...)
		}
		id,_ := strconv.Atoi(row[0])
		location 		:= row[1]
		system 			:= row[2]
		animalClass 	:= row[3]
		year,_ 			:= strconv.Atoi(row[4])
		month,_ 		:= strconv.Atoi(row[5])
		animalNumber,_ 	:= strconv.ParseFloat(row[6],64)
		data.AnimalNumberData = append(data.AnimalNumberData, model.AnimalNumberItem{id, locationsMap[location].Id, locationsMap[location].Name, systemMap[system].Id, systemMap[system].Name,  animalsMap[animalClass].Id, animalsMap[animalClass].Name, year, month, animalNumber})
	}

	rows, err = f.GetRows("EntericFermEFParameters")
	if err != nil {
		return model.DataPackage{},err
	}
	rows = rows[4:]
	defaultValue = []string{"0","0","0","","","","0.0","0.0","0.0","0","0.0","0.0","0.0","0.0","0.0","0.0","0.0","0.0","0","0","0.0","0.0"}
	for _, row := range rows {
		if len(row) <= 0 || len(row[0]) <= 0 || len(row[1]) <= 0{
			continue
		}
		if len(row) < 22 && len(rows) > 0{
			row = append(row, defaultValue[len(row):]...)
		}
		id,_ 				:= strconv.Atoi(row[0])
		year,_ 				:= strconv.Atoi(row[1])
		month,_				:= strconv.Atoi(row[2])
		location 			:= row[3]
		system 				:= row[4]
		animalClass 		:= row[5]
		bodyWeight, _ 		:= strconv.ParseFloat(row[6], 64)
		matureWeight, _ 	:= strconv.ParseFloat(row[7], 64)
		dailyWeightGain, _ 	:= strconv.ParseFloat(row[8], 64)
		fractionOfMonthAlive,_ := strconv.ParseFloat(row[9], 64)
		cf, _ 				:= strconv.ParseFloat(row[10], 64)
		c, _ 				:= strconv.ParseFloat(row[11], 64)
		ca, _ 				:= strconv.ParseFloat(row[12], 64)
		milkProd,_			:= strconv.ParseFloat(row[13], 48)
		fatContent,_		:= strconv.ParseFloat(row[14], 64)
		cPregnancy,_		:= strconv.ParseFloat(row[15], 64)
		proportionAnimalClassPregnant,_	:= strconv.ParseFloat(row[16], 64)
		proportionAnimalClassLactating,_:= strconv.ParseFloat(row[17], 64)
		fractionOfMonthLactating,_		:= strconv.ParseFloat(row[18], 64)
		hoursWorked,_					:= strconv.ParseFloat(row[19], 64)
		de,_							:= strconv.ParseFloat(row[20], 64)
		ym,_							:= strconv.ParseFloat(row[21], 64)
		data.EntericFermEFParameterData = append(data.EntericFermEFParameterData,
			model.EntericFermEFParameterItem{
				id,
				locationsMap[location].Id,
				locationsMap[location].Name,
				systemMap[system].Id,
				systemMap[system].Name,
				animalsMap[animalClass].Id,
				animalsMap[animalClass].Name, year, month,
				bodyWeight,matureWeight,dailyWeightGain,fractionOfMonthAlive,cf,c,ca,milkProd,
				fatContent,cPregnancy,proportionAnimalClassPregnant,proportionAnimalClassLactating,
				fractionOfMonthLactating,hoursWorked,de,ym})
	}

	rows, err = f.GetRows("EntericEmissionFactors")
	if err != nil {
		//Do nothing because we could have no factor table at the first time
		//return model.DataPackage{},nil
	}
	if len(rows) > 0{
		rows = rows[4:]
		defaultValue = []string{"0","0","0","","","","0.0","0.0","0.0","0.0","0.0","0.0","0.0","0.0","0.0","0.0","0.0","0.0","0.0","0.0"}
		for _, row := range rows {
			if len(row) <= 0 || len(row[0]) <= 0 || len(row[1]) <= 0{
				continue
			}
			if len(row) < 20 && len(rows) > 0{
				row = append(row, defaultValue[len(row):]...)
			}
			id,_ 													:= strconv.Atoi(row[0])
			year,_ 													:= strconv.Atoi(row[1])
			month,_													:= strconv.Atoi(row[2])
			location 												:= row[3]
			system 													:= row[4]
			parentClass												:= row[5]
			animalClass 											:= row[6]
			calculatedEF,_											:= strconv.ParseFloat(row[7], 64)
			mapv,_ 													:= strconv.ParseFloat(row[8], 64)
			cofCalculationNEMaintenance,_							:= strconv.ParseFloat(row[9], 64)
			nEMaintenance,_											:= strconv.ParseFloat(row[10], 64)
			nEActivityCattleBuffalo,_ 								:= strconv.ParseFloat(row[11], 64)
			nEGrowthCattleBuffalo,_ 								:= strconv.ParseFloat(row[12], 64)
			nELactationBeefDairyBuffalo,_							:= strconv.ParseFloat(row[13], 64)
			nEWorkCattleBuffalo,_									:= strconv.ParseFloat(row[14], 64)
			nEPregnancyCattleBuffaloSheep,_							:= strconv.ParseFloat(row[15], 64)
			rNEAvailableDietMaintenanceDigestibleEnergyConsumed,_	:= strconv.ParseFloat(row[16], 64)
			rNEAvailableForGrowthDietDisgestibleEnergyConsumed,_	:= strconv.ParseFloat(row[17], 64)
			gECattleBuffaloSheep,_									:= strconv.ParseFloat(row[18], 64)
			entericFermentationEmissionsLivestockCategory,_			:= strconv.ParseFloat(row[19], 64)

			data.EntericEmissionFactorData = append(data.EntericEmissionFactorData, model.EntericEmissionFactorItem{
				id,
				locationsMap[location].Id,
				locationsMap[location].Name,
				systemMap[system].Id,
				systemMap[system].Name,
				parentClass,
				animalsMap[animalClass].Id,
				animalsMap[animalClass].Name, year, month,
				calculatedEF,
				mapv,
				cofCalculationNEMaintenance,
				nEMaintenance,
				nEActivityCattleBuffalo,
				nEGrowthCattleBuffalo,
				nELactationBeefDairyBuffalo,
				nEWorkCattleBuffalo,
				nEPregnancyCattleBuffaloSheep,
				rNEAvailableDietMaintenanceDigestibleEnergyConsumed,
				rNEAvailableForGrowthDietDisgestibleEnergyConsumed,
				gECattleBuffaloSheep,
				entericFermentationEmissionsLivestockCategory,
			})
		}
	}

	sort.SliceStable(data.EntericFermEFParameterData, func(i, j int) bool {
		if data.EntericFermEFParameterData[i].Locationid < data.EntericFermEFParameterData[j].Locationid {
			return true
		}
		if data.EntericFermEFParameterData[i].Locationid > data.EntericFermEFParameterData[j].Locationid {
			return false
		}
		if data.EntericFermEFParameterData[i].Systemid < data.EntericFermEFParameterData[j].Systemid {
			return true
		}
		if data.EntericFermEFParameterData[i].Systemid > data.EntericFermEFParameterData[j].Systemid {
			return false
		}
		if data.EntericFermEFParameterData[i].AnimalClassid < data.EntericFermEFParameterData[j].AnimalClassid {
			return true
		}
		if data.EntericFermEFParameterData[i].AnimalClassid > data.EntericFermEFParameterData[j].AnimalClassid {
			return false
		}
		if data.EntericFermEFParameterData[i].Year < data.EntericFermEFParameterData[j].Year {
			return true
		}
		if data.EntericFermEFParameterData[i].Year > data.EntericFermEFParameterData[j].Year {
			return false
		}
		return data.EntericFermEFParameterData[i].Month <= data.EntericFermEFParameterData[j].Month
	})

	sort.SliceStable(data.EntericEmissionFactorData, func(i, j int) bool {
		if data.EntericEmissionFactorData[i].Locationid < data.EntericEmissionFactorData[j].Locationid {
			return true
		}
		if data.EntericEmissionFactorData[i].Locationid > data.EntericEmissionFactorData[j].Locationid {
			return false
		}
		if data.EntericEmissionFactorData[i].Systemid < data.EntericEmissionFactorData[j].Systemid {
			return true
		}
		if data.EntericEmissionFactorData[i].Systemid > data.EntericEmissionFactorData[j].Systemid {
			return false
		}
		if data.EntericEmissionFactorData[i].AnimalClassid < data.EntericEmissionFactorData[j].AnimalClassid {
			return true
		}
		if data.EntericEmissionFactorData[i].AnimalClassid > data.EntericEmissionFactorData[j].AnimalClassid {
			return false
		}
		if data.EntericEmissionFactorData[i].Year < data.EntericEmissionFactorData[j].Year {
			return true
		}
		if data.EntericEmissionFactorData[i].Year > data.EntericEmissionFactorData[j].Year {
			return false
		}
		return data.EntericEmissionFactorData[i].Month <= data.EntericEmissionFactorData[j].Month
	})

	for _,record := range data.EntericEmissionFactorData{
		data.EntericEmissionFactorDataUserFriendly = append(data.EntericEmissionFactorDataUserFriendly, model.EntericEmissionFactorItemDisplay{
			locationsMap[record.Locationid].Name,
			systemMap[record.Systemid].Name,
			animalsMap[record.AnimalClassid].ParentClass,
			animalsMap[record.AnimalClassid].Name, record.Year, record.Month,
			record.CalculatedEF,
			record.MAP,
			record.CofCalculationNEMaintenance,
			record.NEMaintenance,
			record.NEActivityCattleBuffalo,
			record.NEGrowthCattleBuffalo,
			record.NELactationBeefDairyBuffalo,
			record.NEWorkCattleBuffalo,
			record.NEPregnancyCattleBuffaloSheep,
			record.RNEAvailableDietMaintenanceDigestibleEnergyConsumed,
			record.RNEAvailableForGrowthDietDisgestibleEnergyConsumed,
			record.GECattleBuffaloSheep,
			record.EntericFermentationEmissionsLivestockCategory,
		})
	}

	return data, nil
}

func LoadDataPackageIntoExcel(data model.DataPackage, filename string) ([]byte, error){
	f := excelize.NewFile()
	// Create a new sheet.
	sheet1 := f.NewSheet("Settings")
	f.NewSheet("System")
	f.NewSheet("Location")
	f.NewSheet("AnimalClass")
	f.NewSheet("TemperatureLocation")
	f.NewSheet("AnimalNumbers")
	f.NewSheet("EntericFermEFParameters")
	f.NewSheet("EntericEmissionFactors")
	f.NewSheet("ErrorRecords")


	f.SetCellValue("Settings", "A4", "Setting")
	f.SetCellValue("Settings", "B4", "Value")
	for index,record := range data.SettingData{
		f.SetCellValue("Settings", "A"+strconv.Itoa(index+5), record.Key)
		f.SetCellValue("Settings", "B"+strconv.Itoa(index+5), record.Value)
	}

	f.SetCellValue("System", "A4", "ID")
	f.SetCellValue("System", "B4", "System Name")
	for index,record := range data.SystemData{
		f.SetCellValue("System", "A"+strconv.Itoa(index+5), record.Id)
		f.SetCellValue("System", "B"+strconv.Itoa(index+5), record.Name)
	}

	f.SetCellValue("Location", "A4", "ID")
	f.SetCellValue("Location", "B4", "Location Name")
	for index,record := range data.LocationData{
		f.SetCellValue("Location", "A"+strconv.Itoa(index+5), record.Id)
		f.SetCellValue("Location", "B"+strconv.Itoa(index+5), record.Name)
	}

	f.SetCellValue("AnimalClass", "A4", "ID")
	f.SetCellValue("AnimalClass", "B4", "Parent Class")
	f.SetCellValue("AnimalClass", "C4", "Animal Class Name")
	f.SetCellValue("AnimalClass", "D4", "Default EF")
	for index,record := range data.AnimalClassData{
		f.SetCellValue("AnimalClass", "A"+strconv.Itoa(index+5), record.Id)
		f.SetCellValue("AnimalClass", "B"+strconv.Itoa(index+5), record.ParentClass)
		f.SetCellValue("AnimalClass", "C"+strconv.Itoa(index+5), record.Name)
		f.SetCellValue("AnimalClass", "D"+strconv.Itoa(index+5), record.DefaulEF)
	}

	f.SetCellValue("TemperatureLocation", "A4", "ID")
	f.SetCellValue("TemperatureLocation", "B4", "Location")
	f.SetCellValue("TemperatureLocation", "C4", "Year")
	f.SetCellValue("TemperatureLocation", "D4", "Month")
	f.SetCellValue("TemperatureLocation", "E4", "Average Temp")
	for index,record := range data.TemperatureLocationData{
		f.SetCellValue("TemperatureLocation", "A"+strconv.Itoa(index+5), record.Id)
		f.SetCellValue("TemperatureLocation", "B"+strconv.Itoa(index+5), record.Location)
		f.SetCellValue("TemperatureLocation", "C"+strconv.Itoa(index+5), record.Year)
		f.SetCellValue("TemperatureLocation", "D"+strconv.Itoa(index+5), record.Month)
		f.SetCellValue("TemperatureLocation", "E"+strconv.Itoa(index+5), record.AvgTemp)
	}

	f.SetCellValue("AnimalNumbers", "A4", "ID")
	f.SetCellValue("AnimalNumbers", "B4", "Location")
	f.SetCellValue("AnimalNumbers", "C4", "System")
	f.SetCellValue("AnimalNumbers", "D4", "Animal Class")
	f.SetCellValue("AnimalNumbers", "E4", "Year")
	f.SetCellValue("AnimalNumbers", "F4", "Month")
	f.SetCellValue("AnimalNumbers", "G4", "Animal Number")
	for index,record := range data.AnimalNumberData{
		f.SetCellValue("AnimalNumbers", "A"+strconv.Itoa(index+5), record.Id)
		f.SetCellValue("AnimalNumbers", "B"+strconv.Itoa(index+5), record.Location)
		f.SetCellValue("AnimalNumbers", "C"+strconv.Itoa(index+5), record.System)
		f.SetCellValue("AnimalNumbers", "D"+strconv.Itoa(index+5), record.AnimalClass)
		f.SetCellValue("AnimalNumbers", "E"+strconv.Itoa(index+5), record.Year)
		f.SetCellValue("AnimalNumbers", "F"+strconv.Itoa(index+5), record.Month)
		f.SetCellValue("AnimalNumbers", "G"+strconv.Itoa(index+5), record.AnimalNumber)
	}

	f.SetCellValue("EntericFermEFParameters", "A4", "ID")
	f.SetCellValue("EntericFermEFParameters", "B4", "Year")
	f.SetCellValue("EntericFermEFParameters", "C4", "Month")
	f.SetCellValue("EntericFermEFParameters", "D4", "Location")
	f.SetCellValue("EntericFermEFParameters", "E4", "System")
	f.SetCellValue("EntericFermEFParameters", "F4", "Animal Class")
	f.SetCellValue("EntericFermEFParameters", "G4", "Body Weight (by month, Kg)")
	f.SetCellValue("EntericFermEFParameters", "H4", "Mature Weight (Kg)")
	f.SetCellValue("EntericFermEFParameters", "I4", "Daily Weight Gain (Kg)")
	f.SetCellValue("EntericFermEFParameters", "J4", "Fraction of Month Alive")
	f.SetCellValue("EntericFermEFParameters", "K4", "CF (MJ day \u207b \u00b9 kg\u207b \u00b9)")
	f.SetCellValue("EntericFermEFParameters", "L4", "C")
	f.SetCellValue("EntericFermEFParameters", "M4", "Ca")
	f.SetCellValue("EntericFermEFParameters", "N4", "Milk Production")
	f.SetCellValue("EntericFermEFParameters", "O4", "Fat Content (%)")
	f.SetCellValue("EntericFermEFParameters", "P4", "CPregnancy")
	f.SetCellValue("EntericFermEFParameters", "Q4", "Proportion Animal Class Pregnant")
	f.SetCellValue("EntericFermEFParameters", "R4", "Proportion of Animal Class lactating")
	f.SetCellValue("EntericFermEFParameters", "S4", "Fraction of Month Lactating")
	f.SetCellValue("EntericFermEFParameters", "T4", "Hours Worked")
	f.SetCellValue("EntericFermEFParameters", "U4", "DE (%)")
	f.SetCellValue("EntericFermEFParameters", "V4", "Ym")
	for index,record := range data.EntericFermEFParameterData{
		f.SetCellValue("EntericFermEFParameters", "A"+strconv.Itoa(index+5), record.Id)
		f.SetCellValue("EntericFermEFParameters", "B"+strconv.Itoa(index+5), record.Year)
		f.SetCellValue("EntericFermEFParameters", "C"+strconv.Itoa(index+5), record.Month)
		f.SetCellValue("EntericFermEFParameters", "D"+strconv.Itoa(index+5), record.Location)
		f.SetCellValue("EntericFermEFParameters", "E"+strconv.Itoa(index+5), record.System)
		f.SetCellValue("EntericFermEFParameters", "F"+strconv.Itoa(index+5), record.AnimalClass)
		f.SetCellValue("EntericFermEFParameters", "G"+strconv.Itoa(index+5), record.BodyWeight)
		f.SetCellValue("EntericFermEFParameters", "H"+strconv.Itoa(index+5), record.MatureWeight)
		f.SetCellValue("EntericFermEFParameters", "I"+strconv.Itoa(index+5), record.DailyWeightGain)
		f.SetCellValue("EntericFermEFParameters", "J"+strconv.Itoa(index+5), record.FractionOfMonthAlive)
		f.SetCellValue("EntericFermEFParameters", "K"+strconv.Itoa(index+5), record.CF)
		f.SetCellValue("EntericFermEFParameters", "L"+strconv.Itoa(index+5), record.C)
		f.SetCellValue("EntericFermEFParameters", "M"+strconv.Itoa(index+5), record.CA)
		f.SetCellValue("EntericFermEFParameters", "N"+strconv.Itoa(index+5), record.MilkProd)
		f.SetCellValue("EntericFermEFParameters", "O"+strconv.Itoa(index+5), record.FatContent)
		f.SetCellValue("EntericFermEFParameters", "P"+strconv.Itoa(index+5), record.CPregnancy)
		f.SetCellValue("EntericFermEFParameters", "Q"+strconv.Itoa(index+5), record.ProportionAnimalClassPregnant)
		f.SetCellValue("EntericFermEFParameters", "R"+strconv.Itoa(index+5), record.ProportionAnimalClassLactating)
		f.SetCellValue("EntericFermEFParameters", "S"+strconv.Itoa(index+5), record.FractionOfMonthLactating)
		f.SetCellValue("EntericFermEFParameters", "T"+strconv.Itoa(index+5), record.HoursWorked)
		f.SetCellValue("EntericFermEFParameters", "U"+strconv.Itoa(index+5), record.DE)
		f.SetCellValue("EntericFermEFParameters", "V"+strconv.Itoa(index+5), record.YM)
	}

	f.SetCellValue("EntericEmissionFactors", "A4", "ID")
	f.SetCellValue("EntericEmissionFactors", "B4", "Year")
	f.SetCellValue("EntericEmissionFactors", "C4", "Month")
	f.SetCellValue("EntericEmissionFactors", "D4", "Location")
	f.SetCellValue("EntericEmissionFactors", "E4", "System")
	f.SetCellValue("EntericEmissionFactors", "F4", "Parent Class")
	f.SetCellValue("EntericEmissionFactors", "G4", "Animal Class")
	f.SetCellValue("EntericEmissionFactors", "H4", "Calculated EF (CH&X4 head\u207b \u00b9 year\u207b \u00b9)")
	f.SetCellValue("EntericEmissionFactors", "I4", "Monthly Average Population (head month\u207b \u00b9)")
	f.SetCellValue("EntericEmissionFactors", "J4", "Coefficient for calculating net energy for maintenance")
	f.SetCellValue("EntericEmissionFactors", "K4", "Net energy for maintenance (MJ day\u207b \u00b9)")
	f.SetCellValue("EntericEmissionFactors", "L4", "Net energy for activity for cattle and buffalo (MJ day\u207b \u00b9)")
	f.SetCellValue("EntericEmissionFactors", "M4", "Net energy for growth for cattle and buffalo (MJ day\u207b \u00b9)")
	f.SetCellValue("EntericEmissionFactors", "N4", "Net energy for lactation for beef, dairy and buffalo (MJ day\u207b \u00b9)")
	f.SetCellValue("EntericEmissionFactors", "O4", "Net energy for work for cattle and buffalo (MJ day&\u207b \u00b9)")
	f.SetCellValue("EntericEmissionFactors", "P4", "Net energy for pregnancy for cattle, buffalo and sheep (MJ day\u207b \u00b9)")
	f.SetCellValue("EntericEmissionFactors", "Q4", "Ratio of net energy available in a diet for maintenance to digestible energy consumed")
	f.SetCellValue("EntericEmissionFactors", "R4", "Ratio of net energy available for growth in a diet to digestible energy consumed")
	f.SetCellValue("EntericEmissionFactors", "S4", "Gross Energy for Cattle, Buffalo, Sheep (MJ day\u207b \u00b9)")
	f.SetCellValue("EntericEmissionFactors", "T4", "Enteric fermentation emissions from a livestock category (Gg CH\u2084 month\u207b \u00b9)")
	for index,record := range data.EntericEmissionFactorData{
		f.SetCellValue("EntericEmissionFactors", "A"+strconv.Itoa(index+5), record.Id)
		f.SetCellValue("EntericEmissionFactors", "B"+strconv.Itoa(index+5), record.Year)
		f.SetCellValue("EntericEmissionFactors", "C"+strconv.Itoa(index+5), record.Month)
		f.SetCellValue("EntericEmissionFactors", "D"+strconv.Itoa(index+5), record.Location)
		f.SetCellValue("EntericEmissionFactors", "E"+strconv.Itoa(index+5), record.System)
		f.SetCellValue("EntericEmissionFactors", "F"+strconv.Itoa(index+5), record.ParentClass)
		f.SetCellValue("EntericEmissionFactors", "G"+strconv.Itoa(index+5), record.AnimalClass)
		f.SetCellValue("EntericEmissionFactors", "H"+strconv.Itoa(index+5), record.CalculatedEF)
		f.SetCellValue("EntericEmissionFactors", "I"+strconv.Itoa(index+5), record.MAP)
		f.SetCellValue("EntericEmissionFactors", "J"+strconv.Itoa(index+5), record.CofCalculationNEMaintenance)
		f.SetCellValue("EntericEmissionFactors", "K"+strconv.Itoa(index+5), record.NEMaintenance)
		f.SetCellValue("EntericEmissionFactors", "L"+strconv.Itoa(index+5), record.NEActivityCattleBuffalo)
		f.SetCellValue("EntericEmissionFactors", "M"+strconv.Itoa(index+5), record.NEGrowthCattleBuffalo)
		f.SetCellValue("EntericEmissionFactors", "N"+strconv.Itoa(index+5), record.NELactationBeefDairyBuffalo)
		f.SetCellValue("EntericEmissionFactors", "O"+strconv.Itoa(index+5), record.NEWorkCattleBuffalo)
		f.SetCellValue("EntericEmissionFactors", "P"+strconv.Itoa(index+5), record.NEPregnancyCattleBuffaloSheep)
		f.SetCellValue("EntericEmissionFactors", "Q"+strconv.Itoa(index+5), record.RNEAvailableDietMaintenanceDigestibleEnergyConsumed)
		f.SetCellValue("EntericEmissionFactors", "R"+strconv.Itoa(index+5), record.RNEAvailableForGrowthDietDisgestibleEnergyConsumed)
		f.SetCellValue("EntericEmissionFactors", "S"+strconv.Itoa(index+5), record.GECattleBuffaloSheep)
		f.SetCellValue("EntericEmissionFactors", "T"+strconv.Itoa(index+5), record.EntericFermentationEmissionsLivestockCategory)
	}

	f.SetCellValue("ErrorRecords", "A4", "ID")
	f.SetCellValue("ErrorRecords", "B4", "Year")
	f.SetCellValue("ErrorRecords", "C4", "Month")
	f.SetCellValue("ErrorRecords", "D4", "Location")
	f.SetCellValue("ErrorRecords", "E4", "System")
	f.SetCellValue("ErrorRecords", "F4", "Animal Class")
	f.SetCellValue("ErrorRecords", "G4", "Error Message")
	for index,record := range data.ErrorRecordData{
		f.SetCellValue("ErrorRecords", "A"+strconv.Itoa(index+5), record.Id)
		f.SetCellValue("ErrorRecords", "B"+strconv.Itoa(index+5), record.Year)
		f.SetCellValue("ErrorRecords", "C"+strconv.Itoa(index+5), record.Month)
		f.SetCellValue("ErrorRecords", "D"+strconv.Itoa(index+5), record.Location)
		f.SetCellValue("ErrorRecords", "E"+strconv.Itoa(index+5), record.System)
		f.SetCellValue("ErrorRecords", "F"+strconv.Itoa(index+5), record.AnimalClass)
		f.SetCellValue("ErrorRecords", "G"+strconv.Itoa(index+5), record.ErrorMsg)
	}

	f.SetActiveSheet(sheet1)
	buf, err := f.WriteToBuffer()
	if  err != nil {
		log.Printf("Unable to marshal excel sheet: %v", err)
		return []byte{}, err
	}
	return buf.Bytes(), nil
}
