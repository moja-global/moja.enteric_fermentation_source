package helpers

import (
	model "github.com/MullionGroup/go-website-flintpro-example/models"
	"math"
)

func CalculateMonthlyAveragePopulation(fractionOfMonthAlive float64, populationPerAnimalClassPerSystem float64) float64{
	return fractionOfMonthAlive * populationPerAnimalClassPerSystem
}

func CalculateCoefficientForCalculatingNetEnergyForMaintenance(factor float64, setValue float64, averageTemperatureByMonth float64) float64{
	return factor + (setValue * (20-averageTemperatureByMonth))
}

func CalculateNetEnergyForMaintenance(factor float64, setValue float64, animalWeightPerClassPerSystem float64) float64{
	return factor * math.Pow(animalWeightPerClassPerSystem,setValue)
}

func CalculateNetEnergyForActivityCattleBuffalo(factor float64, netEnergyForMaintenance float64) float64{
	return factor * netEnergyForMaintenance
}

func CalculateNetEnergyForGrowthCattleBuffalo(setValue1 float64, setValue2 float64, setValue3 float64, liveBodyWeight float64, animalTypeCoefficient float64, matureWeight float64, weightGain float64) float64{
	return setValue1 * (math.Pow(liveBodyWeight/(animalTypeCoefficient * matureWeight),setValue2)) * math.Pow(weightGain,setValue3)
}

func CalculateNetEnergyForLactationBeefCattleDairyCattleBuffalo(setValue1 float64, setValue2 float64, milkProduced float64, fatContent float64, proportionAnimalClassLactating float64, fractionOfMonthLactating float64) float64{
	return (milkProduced * proportionAnimalClassLactating * fractionOfMonthLactating * (setValue1 + (setValue2 * fatContent)))
}

func CalculateNetEnergyForWorkCattleBuffalo(setValue1 float64, hours float64, netEnergyMaintenance float64) float64{
	return setValue1 * netEnergyMaintenance * hours
}

func CalculateNetEnergyForPregnancyCattleBuffaloSheep(cPregnancy float64, netEnergyMaintenance float64, proportionOfAnimalClassPregnant float64) float64{
	return (cPregnancy * netEnergyMaintenance) * proportionOfAnimalClassPregnant
}

func CalculateRatioOfNetEnergyAvailableInADietForMaintenanceToDisgestibleEnergyConsumed(setValue1 float64, setValue2 float64, digestibleEnergy float64, setValue3 float64, setValue4 float64) float64{
	return (setValue1 - (setValue2 * digestibleEnergy)) + (setValue3 * math.Pow(digestibleEnergy, 2)) - (setValue4 / digestibleEnergy)
}

func CalculateRatioNetEnergyAvailableForGrowthInADietToDigestibleEnergyConsumed(setValue1 float64, setValue2 float64, digestibleEnergy float64, setValue3 float64, setValue4 float64) float64{
	return (setValue1 - (setValue2 * digestibleEnergy)) + (setValue3 * math.Pow(digestibleEnergy, 2)) - (setValue4 / digestibleEnergy)
}

func CalculateGrossEnergyCattleBuffaloSheep(NEm float64, NEa float64, NEi float64, NEwork float64, NEwool float64, NEp float64, REM float64, NEg float64, REG float64, digestibleEnergy float64) float64{
	equation1 := (NEm + NEa + NEi + NEwork + NEp) / REM
	equation2 := (NEg + NEwool) / REG
	return (equation1 + equation2) / (digestibleEnergy/100)
}

func CalculateEmissionFactorsForEntericFermentationFromALivestockCategory(GE float64, Ym float64, setValue1 float64, setValue2 int, setValue3 float64) float64{
	return (GE * (Ym/setValue1) * float64(setValue2)) / setValue3
}

func CalculateEntericFermentationEmissionFromALivestockCategory(EF float64, livestockNumbersPerClass float64) float64{
	return float64(EF) * float64(livestockNumbersPerClass / math.Pow(10, 6))
}

func CalculateTotalEmissionsFromLivestockEntericFermentation(records []model.EntericEmissionFactorItem) float64{
	result := 0.0
	for _, value := range records{
		result += value.CalculatedEF
	}
	return result
}