package helpers

import (
	model "github.com/MullionGroup/go-website-flintpro-example/models"
	"github.com/sirupsen/logrus"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"
)

func RunSimulation(data model.DataPackage) (model.DataPackage, error){
	logrus.Info("Starting a simulation...")
	settingsMap := make(map[string]string)
	for _,record := range data.SettingData{
		settingsMap[record.Key] = record.Value
	}

	/*sY,_ := strconv.Atoi(settingsMap["Start Year"])
	sM,_ := strconv.Atoi(settingsMap["Start Month"])
	eY,_ := strconv.Atoi(settingsMap["End Year"])
	eM,_ := strconv.Atoi(settingsMap["End Month"])*/
	eDate,_ := settingsMap["End Date"]
	eDateStr := strings.Split(eDate, "/")
	eYear,err := strconv.Atoi(eDateStr[len(eDateStr)-1])
	if err != nil {
		data.ErrorGenericMsg = append(data.ErrorGenericMsg, "Error settings date format")
	}

	data, err = ValidateAndPreProcessData(data, eYear)
	if err != nil{
		return model.DataPackage{}, err
	}

	_,_,animalsMap := data.GenerateMappingsForPrimaryTable()

	temperatureLocationsMap := make(map[string]model.TemperatureLocationItem)
	for _,record := range data.TemperatureLocationData{
		yearStr := strconv.Itoa(record.Year)
		monthStr := strconv.Itoa(record.Month)
		locationStr := strconv.Itoa(record.Locationid)
		temperatureLocationsMap[yearStr + "*" + monthStr + "*" + locationStr] = record
	}

	animalNumbersMap := make(map[string]model.AnimalNumberItem)
	for _,record := range data.AnimalNumberData{
		yearStr := strconv.Itoa(record.Year)
		monthStr := strconv.Itoa(record.Month)
		locationStr := strconv.Itoa(record.Locationid)
		systemStr := strconv.Itoa(record.Systemid)
		animalClassStr := strconv.Itoa(record.AnimalClassid)
		animalNumbersMap[yearStr + "*" + monthStr + "*" + locationStr + "*" + systemStr + "*" + animalClassStr] = record
	}

	entericEmissionFactors := []model.EntericEmissionFactorItem{}
	entericEmissionFactorsMap := make(map[string]model.EntericEmissionFactorItem)
	for _,record := range data.EntericEmissionFactorData{
		yearStr := strconv.Itoa(record.Year)
		monthStr := strconv.Itoa(record.Month)
		locationStr := strconv.Itoa(record.Locationid)
		systemStr := strconv.Itoa(record.Systemid)
		animalClassStr := strconv.Itoa(record.AnimalClassid)
		entericEmissionFactorsMap[yearStr + "*" + monthStr + "*" + locationStr + "*" + systemStr + "*" + animalClassStr] = record
		entericEmissionFactors = append(entericEmissionFactors, record)
	}

	newEntericEmissionFactors := []model.EntericEmissionFactorItem{}
	for _,record := range data.EntericFermEFParameterData{
		errorneous := false
		entericEmissionFactor := model.EntericEmissionFactorItem{}
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

		entericEmissionFactor.Year = record.Year
		entericEmissionFactor.Month = record.Month
		entericEmissionFactor.Id = record.Id
		entericEmissionFactor.Locationid = record.Locationid
		entericEmissionFactor.Location = record.Location
		entericEmissionFactor.Systemid = record.Systemid
		entericEmissionFactor.System = record.System
		entericEmissionFactor.AnimalClassid = record.AnimalClassid
		entericEmissionFactor.AnimalClass = record.AnimalClass

		calculateEmissionFactor(record, &entericEmissionFactor, temperatureLocationsMap, animalNumbersMap, animalsMap)
		if math.IsNaN(entericEmissionFactor.MAP){
			errorRecord.ErrorMsg = "Not a number: montly average population"
			data.ErrorRecordData = append(data.ErrorRecordData, errorRecord)
			errorneous = true
		}
		if math.IsInf(entericEmissionFactor.MAP, 0){
			errorRecord.ErrorMsg = "Infinite: montly average population"
			data.ErrorRecordData = append(data.ErrorRecordData, errorRecord)
			errorneous = true
		}
		if math.IsNaN(entericEmissionFactor.CofCalculationNEMaintenance){
			errorRecord.ErrorMsg = "Not a number: coefficient for calculating net energy for maintenance"
			data.ErrorRecordData = append(data.ErrorRecordData, errorRecord)
			errorneous = true
		}
		if math.IsInf(entericEmissionFactor.CofCalculationNEMaintenance, 0){
			errorRecord.ErrorMsg = "Infinite: coefficient for calculating net energy for maintenance"
			data.ErrorRecordData = append(data.ErrorRecordData, errorRecord)
			errorneous = true
		}
		if math.IsNaN(entericEmissionFactor.NEMaintenance){
			errorRecord.ErrorMsg = "Not a number: net energy for maintenance"
			data.ErrorRecordData = append(data.ErrorRecordData, errorRecord)
			errorneous = true
		}
		if math.IsInf(entericEmissionFactor.NEMaintenance, 0){
			errorRecord.ErrorMsg = "Infinite: net energy for maintenance"
			data.ErrorRecordData = append(data.ErrorRecordData, errorRecord)
			errorneous = true
		}
		if math.IsNaN(entericEmissionFactor.NEActivityCattleBuffalo){
			errorRecord.ErrorMsg = "Not a number: net energy for activity for cattle and buffalo"
			data.ErrorRecordData = append(data.ErrorRecordData, errorRecord)
			errorneous = true
		}
		if math.IsInf(entericEmissionFactor.NEActivityCattleBuffalo, 0){
			errorRecord.ErrorMsg = "Infinite: net energy for activity for cattle and buffalo"
			data.ErrorRecordData = append(data.ErrorRecordData, errorRecord)
			errorneous = true
		}
		if math.IsNaN(entericEmissionFactor.NEGrowthCattleBuffalo){
			errorRecord.ErrorMsg = "Not a number: net energy for growth for cattle and buffalo"
			data.ErrorRecordData = append(data.ErrorRecordData, errorRecord)
			errorneous = true
		}
		if math.IsInf(entericEmissionFactor.NEGrowthCattleBuffalo, 0){
			errorRecord.ErrorMsg = "Infinite: net energy for growth for cattle and buffalo"
			data.ErrorRecordData = append(data.ErrorRecordData, errorRecord)
			errorneous = true
		}
		if math.IsNaN(entericEmissionFactor.NELactationBeefDairyBuffalo){
			errorRecord.ErrorMsg = "Not a number: net energy for lactation for beef, dairy and buffalo"
			data.ErrorRecordData = append(data.ErrorRecordData, errorRecord)
			errorneous = true
		}
		if math.IsInf(entericEmissionFactor.NELactationBeefDairyBuffalo, 0){
			errorRecord.ErrorMsg = "Infinite: net energy for lactation for beef, dairy and buffalo"
			data.ErrorRecordData = append(data.ErrorRecordData, errorRecord)
			errorneous = true
		}
		if math.IsNaN(entericEmissionFactor.NEWorkCattleBuffalo){
			errorRecord.ErrorMsg = "Not a number: net energy for work for cattle and buffalo"
			data.ErrorRecordData = append(data.ErrorRecordData, errorRecord)
			errorneous = true
		}
		if math.IsInf(entericEmissionFactor.NEWorkCattleBuffalo, 0){
			errorRecord.ErrorMsg = "Infinite: net energy for work for cattle and buffalo"
			data.ErrorRecordData = append(data.ErrorRecordData, errorRecord)
			errorneous = true
		}
		if math.IsNaN(entericEmissionFactor.NEPregnancyCattleBuffaloSheep){
			errorRecord.ErrorMsg = "Not a number: net energy for pregnancy for cattle, buffalo and sheep"
			data.ErrorRecordData = append(data.ErrorRecordData, errorRecord)
			errorneous = true
		}
		if math.IsInf(entericEmissionFactor.NEPregnancyCattleBuffaloSheep, 0){
			errorRecord.ErrorMsg = "Infinite: net energy for pregnancy for cattle, buffalo and sheep"
			data.ErrorRecordData = append(data.ErrorRecordData, errorRecord)
			errorneous = true
		}
		if math.IsNaN(entericEmissionFactor.RNEAvailableDietMaintenanceDigestibleEnergyConsumed){
			errorRecord.ErrorMsg = "Not a number: ratio of net energy available in a diet for maintenance to digestible energy consumed"
			data.ErrorRecordData = append(data.ErrorRecordData, errorRecord)
			errorneous = true
		}
		if math.IsInf(entericEmissionFactor.RNEAvailableDietMaintenanceDigestibleEnergyConsumed, 0){
			errorRecord.ErrorMsg = "Infinite: ratio of net energy available in a diet for maintenance to digestible energy consumed"
			data.ErrorRecordData = append(data.ErrorRecordData, errorRecord)
			errorneous = true
		}
		if math.IsNaN(entericEmissionFactor.RNEAvailableForGrowthDietDisgestibleEnergyConsumed){
			errorRecord.ErrorMsg = "Not a number: ratio of net energy available for growth in a diet to digestible energy consumed"
			data.ErrorRecordData = append(data.ErrorRecordData, errorRecord)
			errorneous = true
		}
		if math.IsInf(entericEmissionFactor.RNEAvailableForGrowthDietDisgestibleEnergyConsumed, 0){
			errorRecord.ErrorMsg = "Infinite: ratio of net energy available for growth in a diet to digestible energy consumed"
			data.ErrorRecordData = append(data.ErrorRecordData, errorRecord)
			errorneous = true
		}
		if math.IsNaN(entericEmissionFactor.GECattleBuffaloSheep){
			errorRecord.ErrorMsg = "Not a number: gross energy for cattle, buffalo, sheep"
			data.ErrorRecordData = append(data.ErrorRecordData, errorRecord)
			errorneous = true
		}
		if math.IsInf(entericEmissionFactor.GECattleBuffaloSheep, 0){
			errorRecord.ErrorMsg = "Infinite: gross energy for cattle, buffalo, sheep"
			data.ErrorRecordData = append(data.ErrorRecordData, errorRecord)
			errorneous = true
		}
		if math.IsNaN(entericEmissionFactor.CalculatedEF){
			errorRecord.ErrorMsg = "Not a number: calculated emission factor"
			data.ErrorRecordData = append(data.ErrorRecordData, errorRecord)
			errorneous = true
		}
		if math.IsInf(entericEmissionFactor.CalculatedEF, 0){
			errorRecord.ErrorMsg = "Infinite: calculated emission factor"
			data.ErrorRecordData = append(data.ErrorRecordData, errorRecord)
			errorneous = true
		}
		if math.IsNaN(entericEmissionFactor.EntericFermentationEmissionsLivestockCategory){
			errorRecord.ErrorMsg = "Not a number: enteric fermentation emissions from a livestock category"
			data.ErrorRecordData = append(data.ErrorRecordData, errorRecord)
			errorneous = true
		}
		if math.IsInf(entericEmissionFactor.EntericFermentationEmissionsLivestockCategory, 0){
			errorRecord.ErrorMsg = "Infinite: enteric fermentation emissions from a livestock category"
			data.ErrorRecordData = append(data.ErrorRecordData, errorRecord)
			errorneous = true
		}

		if errorneous == false{
			newEntericEmissionFactors = append(newEntericEmissionFactors, entericEmissionFactor)
			entericEmissionFactorDisplay := model.EntericEmissionFactorItemDisplay{
				Location:                      entericEmissionFactor.Location,
				System:                        entericEmissionFactor.System,
				ParentClass: 				   entericEmissionFactor.ParentClass,
				AnimalClass:                   entericEmissionFactor.AnimalClass,
				Year:                          entericEmissionFactor.Year,
				Month:                         entericEmissionFactor.Month,
				CalculatedEF:                  entericEmissionFactor.CalculatedEF,
				MAP:                           entericEmissionFactor.MAP,
				CofCalculationNEMaintenance:   entericEmissionFactor.CofCalculationNEMaintenance,
				NEMaintenance:                 entericEmissionFactor.NEMaintenance,
				NEActivityCattleBuffalo:       entericEmissionFactor.NEActivityCattleBuffalo,
				NEGrowthCattleBuffalo:         entericEmissionFactor.NEGrowthCattleBuffalo,
				NELactationBeefDairyBuffalo:   entericEmissionFactor.NELactationBeefDairyBuffalo,
				NEWorkCattleBuffalo:           entericEmissionFactor.NEWorkCattleBuffalo,
				NEPregnancyCattleBuffaloSheep: entericEmissionFactor.NEPregnancyCattleBuffaloSheep,
				RNEAvailableDietMaintenanceDigestibleEnergyConsumed: entericEmissionFactor.RNEAvailableDietMaintenanceDigestibleEnergyConsumed,
				RNEAvailableForGrowthDietDisgestibleEnergyConsumed:  entericEmissionFactor.RNEAvailableForGrowthDietDisgestibleEnergyConsumed,
				GECattleBuffaloSheep:                          entericEmissionFactor.GECattleBuffaloSheep,
				EntericFermentationEmissionsLivestockCategory: entericEmissionFactor.EntericFermentationEmissionsLivestockCategory,
			}
			data.EntericEmissionFactorDataUserFriendly = append(data.EntericEmissionFactorDataUserFriendly, entericEmissionFactorDisplay)
		}
	}
	data.EntericEmissionFactorData = newEntericEmissionFactors
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
	sort.SliceStable(data.AnimalNumberData, func(i, j int) bool {
		if data.AnimalNumberData[i].Locationid < data.AnimalNumberData[j].Locationid {
			return true
		}
		if data.AnimalNumberData[i].Locationid > data.AnimalNumberData[j].Locationid {
			return false
		}
		if data.AnimalNumberData[i].Systemid < data.AnimalNumberData[j].Systemid {
			return true
		}
		if data.AnimalNumberData[i].Systemid > data.AnimalNumberData[j].Systemid {
			return false
		}
		if data.AnimalNumberData[i].AnimalClassid < data.AnimalNumberData[j].AnimalClassid {
			return true
		}
		if data.AnimalNumberData[i].AnimalClassid > data.AnimalNumberData[j].AnimalClassid {
			return false
		}
		if data.AnimalNumberData[i].Year < data.AnimalNumberData[j].Year {
			return true
		}
		if data.AnimalNumberData[i].Year > data.AnimalNumberData[j].Year {
			return false
		}
		return data.AnimalNumberData[i].Month <= data.AnimalNumberData[j].Month
	})
	sort.SliceStable(data.TemperatureLocationData, func(i, j int) bool {
		if data.TemperatureLocationData[i].Locationid < data.TemperatureLocationData[j].Locationid {
			return true
		}
		if data.TemperatureLocationData[i].Locationid > data.TemperatureLocationData[j].Locationid {
			return false
		}
		if data.TemperatureLocationData[i].Year < data.TemperatureLocationData[j].Year {
			return true
		}
		if data.TemperatureLocationData[i].Year > data.TemperatureLocationData[j].Year {
			return false
		}
		return data.TemperatureLocationData[i].Month <= data.TemperatureLocationData[j].Month
	})
	logrus.Info("Finishing a simulation...")
	return data, nil
}

func calculateEmissionFactor(target model.EntericFermEFParameterItem, result *model.EntericEmissionFactorItem, tempLocationMaps map[string]model.TemperatureLocationItem, animalNumbersMap map[string]model.AnimalNumberItem, animalsMap map[interface{}]model.AnimalClassDataItem){
	yearStr := strconv.Itoa(target.Year)
	monthStr := strconv.Itoa(target.Month)
	locationStr := strconv.Itoa(target.Locationid)
	systemStr := strconv.Itoa(target.Systemid)
	animalClassStr := strconv.Itoa(target.AnimalClassid)
	time := time.Date(target.Year, time.Month(target.Month+1), 0, 0, 0, 0, 0, time.UTC) // the day before 2000-03-01
	result.MAP = CalculateMonthlyAveragePopulation(target.FractionOfMonthAlive, animalNumbersMap[yearStr + "*" + monthStr + "*" + locationStr + "*" + systemStr + "*" + animalClassStr].AnimalNumber)
	result.CofCalculationNEMaintenance = CalculateCoefficientForCalculatingNetEnergyForMaintenance(target.CF, 0.0048, tempLocationMaps[yearStr + "*" + monthStr + "*" + locationStr].AvgTemp)
	result.NEMaintenance = CalculateNetEnergyForMaintenance(result.CofCalculationNEMaintenance, 0.75, target.BodyWeight)
	result.NEActivityCattleBuffalo = CalculateNetEnergyForActivityCattleBuffalo(target.CA,result.NEMaintenance)
	result.NEGrowthCattleBuffalo = CalculateNetEnergyForGrowthCattleBuffalo(22.02, 0.75, 1.097, target.BodyWeight, target.C, target.MatureWeight, target.DailyWeightGain)
	result.NELactationBeefDairyBuffalo = CalculateNetEnergyForLactationBeefCattleDairyCattleBuffalo(1.47, 0.4, target.MilkProd, target.FatContent, target.ProportionAnimalClassLactating, target.FractionOfMonthLactating)
	result.NEWorkCattleBuffalo = CalculateNetEnergyForWorkCattleBuffalo(0.1, target.HoursWorked, result.NEMaintenance)
	result.NEPregnancyCattleBuffaloSheep = CalculateNetEnergyForPregnancyCattleBuffaloSheep(target.CPregnancy, result.NEMaintenance, target.ProportionAnimalClassPregnant)
	result.RNEAvailableDietMaintenanceDigestibleEnergyConsumed = CalculateRatioOfNetEnergyAvailableInADietForMaintenanceToDisgestibleEnergyConsumed(1.123, 4.092 * math.Pow(10, -3), target.DE, 1.126 * math.Pow(10, -5), 25.4)
	result.RNEAvailableForGrowthDietDisgestibleEnergyConsumed = CalculateRatioNetEnergyAvailableForGrowthInADietToDigestibleEnergyConsumed(1.164, 5.160 * math.Pow(10, -3), target.DE, 1.308 * math.Pow(10, -5), 37.4)
	result.GECattleBuffaloSheep = CalculateGrossEnergyCattleBuffaloSheep(result.NEMaintenance, result.NEActivityCattleBuffalo, result.NELactationBeefDairyBuffalo, result.NEWorkCattleBuffalo, 0, result.NEPregnancyCattleBuffaloSheep, result.RNEAvailableDietMaintenanceDigestibleEnergyConsumed, result.NEGrowthCattleBuffalo, result.RNEAvailableForGrowthDietDisgestibleEnergyConsumed, target.DE)
	result.CalculatedEF = CalculateEmissionFactorsForEntericFermentationFromALivestockCategory(result.GECattleBuffaloSheep, target.YM, 100, time.Day(), 55.65)
	if result.CalculatedEF <= 0 {
		result.EntericFermentationEmissionsLivestockCategory = CalculateEntericFermentationEmissionFromALivestockCategory(animalsMap[result.AnimalClassid].DefaulEF, result.MAP)
	} else{
		result.EntericFermentationEmissionsLivestockCategory = CalculateEntericFermentationEmissionFromALivestockCategory(result.CalculatedEF, result.MAP)
	}
	result.ParentClass = animalsMap[result.AnimalClassid].ParentClass
}