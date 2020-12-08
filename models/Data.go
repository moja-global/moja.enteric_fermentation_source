package models

import (
	"fmt"
	"github.com/alexeyco/simpletable"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/twinj/uuid"
)

// Data Structs

type SettingDataItem struct {
	Key 	string	`json:"key"`
	Value 	string	`json:"value"`
}

type SystemDataItem struct {
	Id 		int	`json:"id" gorm:"primary_key"`
	Name 	string	`json:"name"`
}

type LocationDataItem struct {
	Id 		int	`json:"id" gorm:"primary_key"`
	Name 	string	`json:"name"`
}

func (item *LocationDataItem) BeforeCreate(scope *gorm.Scope) (err error) {
	if item.Id == 0 {
		scope.SetColumn("id", uuid.NewV4().String())
	}
	return
}

type AnimalClassDataItem struct {
	Id 			int	`json:"id" gorm:"primary_key"`
	ParentClass string	`json:"parent_class"`
	Name 		string	`json:"name"`
	DefaulEF 	float64 `json:"default_ef"`
}


type TemperatureLocationItem struct {
	Id 			int		`json:"id" gorm:"primary_key"`
	Locationid 	int			`json:"location_id"`
	Location 	string		`json:"location" gorm:"-"`
	Year 		int			`json:"year"`
	Month 		int			`json:"month"`
	AvgTemp		float64		`json:"avg_temp"`
}

type AnimalNumberItem struct {
	Id 				int			`json:"id" gorm:"primary_key"`
	Locationid 		int			`json:"location_id"`
	Location 		string		`json:"location" gorm:"-"`
	Systemid 		int			`json:"system_id"`
	System 			string		`json:"system" gorm:"-"`
	AnimalClassid 	int			`json:"animal_class_id"`
	AnimalClass 	string		`json:"animal_class" gorm:"-"`
	Year 			int			`json:"year"`
	Month 			int			`json:"month"`
	AnimalNumber 	float64		`json:"animal_number"`
}

type EntericFermEFParameterItem struct {
	Id 								int		`json:"id" gorm:"primary_key"`
	Locationid 						int			`json:"location_id"`
	Location 						string		`json:"location" gorm:"-"`
	Systemid 						int			`json:"system_id"`
	System 							string		`json:"system" gorm:"-"`
	AnimalClassid 					int			`json:"animal_class_id"`
	AnimalClass 					string		`json:"animal_class" gorm:"-"`
	Year 							int			`json:"year"`
	Month 							int			`json:"month"`
	BodyWeight 						float64		`json:"body_weight"`
	MatureWeight 					float64		`json:"mature_weight"`
	DailyWeightGain 				float64		`json:"daily_weight_gain"`
	FractionOfMonthAlive 			float64		`json:"fraction_of_month_alive"`
	CF 								float64		`json:"cf"`
	C 								float64		`json:"c"`
	CA 								float64		`json:"ca"`
	MilkProd 						float64		`json:"milk_prod"`
	FatContent 						float64		`json:"fat_content"`
	CPregnancy 						float64		`json:"c_pregnancy"`
	ProportionAnimalClassPregnant 	float64		`json:"proportion_animal_class_pregnant"`
	ProportionAnimalClassLactating 	float64		`json:"proportion_animal_class_lactating"`
	FractionOfMonthLactating 		float64		`json:"fraction_of_month_lactating"`
	HoursWorked 					float64		`json:"hours_worked"`
	DE 								float64		`json:"de"`
	YM 								float64		`json:"ym"`
}

type EntericEmissionFactorItem struct {
	Id 													int			`json:"id" gorm:"primary_key"`
	Locationid 											int			`json:"location_id"`
	Location 											string		`json:"location" gorm:"-"`
	Systemid 											int			`json:"system_id"`
	System 												string		`json:"system" gorm:"-"`
	ParentClass											string		`json:"parent_class" gorm:"parent_class"`
	AnimalClassid 										int			`json:"animal_class_id"`
	AnimalClass 										string		`json:"animal_class" gorm:"-"`
	Year 												int			`json:"year"`
	Month 												int			`json:"month"`
	CalculatedEF 										float64		`json:"calculated_ef" gorm:"calculated_emission_factor"`
	MAP													float64		`json:"monthly_average_population" gorm:"monthly_average_population"`
	CofCalculationNEMaintenance							float64		`json:"coefficient_for_calculating_net_energy_for_maintenance" gorm:"coefficient_for_calculating_net_energy_for_maintenance"`
	NEMaintenance										float64		`json:"net_energy_for_maintenance" gorm:"net_energy_for_maintenance"`
	NEActivityCattleBuffalo								float64		`json:"net_energy_for_activity_for_cattle_and_buffalo" gorm:"net_energy_for_activity_for_cattle_and_buffalo"`
	NEGrowthCattleBuffalo								float64		`json:"net_energy_for_growth_for_cattle_and_buffalo" gorm:"net_energy_for_growth_for_cattle_and_buffalo"`
	NELactationBeefDairyBuffalo							float64		`json:"net_energy_for_lactation_for_beef_dairy_and_buffalo" gorm:"net_energy_for_lactation_for_beef_dairy_and_buffalo"`
	NEWorkCattleBuffalo									float64		`json:"net_energy_for_work_for_cattle_and_buffalo" gorm:"net_energy_for_work_for_cattle_and_buffalo"`
	NEPregnancyCattleBuffaloSheep						float64		`json:"net_energy_for_pregnancy_for_cattle_buffalo_and_sheep" gorm:"net_energy_for_pregnancy_for_cattle_buffalo_and_sheep"`
	RNEAvailableDietMaintenanceDigestibleEnergyConsumed	float64		`json:"ratio_of_net_energy_available_in_a_diet_for_maintenance_to_digestible_energy_consumed" gorm:"ratio_of_net_energy_available_in_a_diet_for_maintenance_to_digestible_energy_consumed"`
	RNEAvailableForGrowthDietDisgestibleEnergyConsumed	float64		`json:"ratio_of_net_energy_available_for_growth_in_a_diet_to_digestible_energy_consumed" gorm:"ratio_of_net_energy_available_for_growth_in_a_diet_to_digestible_energy_consumed"`
	GECattleBuffaloSheep								float64		`json:"gross_energy_for_cattle_buffalo_sheep" gorm:"gross_energy_for_cattle_buffalo_sheep"`
	EntericFermentationEmissionsLivestockCategory		float64		`json:"enteric_fermentation_emissions_from_a_livestock_category" gorm:"enteric_fermentation_emissions_from_a_livestock_category"`
}

type EntericEmissionFactorItemDisplay struct {
	Location 											string		`json:"location"`
	System 												string		`json:"system"`
	ParentClass											string		`json:"parent_class"`
	AnimalClass 										string		`json:"animal_class"`
	Year 												int			`json:"year"`
	Month 												int			`json:"month"`
	CalculatedEF 										float64		`json:"calculated_ef" gorm:"calculated_emission_factor"`
	MAP													float64		`json:"monthly_average_population" gorm:"monthly_average_population"`
	CofCalculationNEMaintenance							float64		`json:"coefficient_for_calculating_net_energy_for_maintenance" gorm:"coefficient_for_calculating_net_energy_for_maintenance"`
	NEMaintenance										float64		`json:"net_energy_for_maintenance" gorm:"net_energy_for_maintenance"`
	NEActivityCattleBuffalo								float64		`json:"net_energy_for_activity_for_cattle_and_buffalo" gorm:"net_energy_for_activity_for_cattle_and_buffalo"`
	NEGrowthCattleBuffalo								float64		`json:"net_energy_for_growth_for_cattle_and_buffalo" gorm:"net_energy_for_growth_for_cattle_and_buffalo"`
	NELactationBeefDairyBuffalo							float64		`json:"net_energy_for_lactation_for_beef_dairy_and_buffalo" gorm:"net_energy_for_lactation_for_beef_dairy_and_buffalo"`
	NEWorkCattleBuffalo									float64		`json:"net_energy_for_work_for_cattle_and_buffalo" gorm:"net_energy_for_work_for_cattle_and_buffalo"`
	NEPregnancyCattleBuffaloSheep						float64		`json:"net_energy_for_pregnancy_for_cattle_buffalo_and_sheep" gorm:"net_energy_for_pregnancy_for_cattle_buffalo_and_sheep"`
	RNEAvailableDietMaintenanceDigestibleEnergyConsumed	float64		`json:"ratio_of_net_energy_available_in_a_diet_for_maintenance_to_digestible_energy_consumed" gorm:"ratio_of_net_energy_available_in_a_diet_for_maintenance_to_digestible_energy_consumed"`
	RNEAvailableForGrowthDietDisgestibleEnergyConsumed	float64		`json:"ratio_of_net_energy_available_for_growth_in_a_diet_to_digestible_energy_consumed" gorm:"ratio_of_net_energy_available_for_growth_in_a_diet_to_digestible_energy_consumed"`
	GECattleBuffaloSheep								float64		`json:"gross_energy_for_cattle_buffalo_sheep" gorm:"gross_energy_for_cattle_buffalo_sheep"`
	EntericFermentationEmissionsLivestockCategory		float64		`json:"enteric_fermentation_emissions_from_a_livestock_category" gorm:"enteric_fermentation_emissions_from_a_livestock_category"`
}

type ErrorRecordItem struct {
	Id 													int			`json:"id" sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Recordid 											int			`json:"record_id"`
	Locationid 											int			`json:"location_id"`
	Location 											string		`json:"-" gorm:"-"`
	Systemid 											int			`json:"system_id"`
	System 												string		`json:"-" gorm:"-"`
	AnimalClassid 										int			`json:"animal_class_id"`
	AnimalClass 										string		`json:"-" gorm:"-"`
	Year 												int			`json:"year"`
	Month 												int			`json:"month"`
	ErrorMsg 											string		`json:"error_msg"`
}

type DataPackage struct {
	SheetID									string								`json:"-" gorm:"-"`
	FileName								string								`json:"-" gorm:"-"`
	SettingData								[]SettingDataItem					`json:"setting_data"`
	SystemData 								[]SystemDataItem            		`json:"system_data"`
	LocationData 							[]LocationDataItem          		`json:"location_data"`
	AnimalClassData 						[]AnimalClassDataItem       		`json:"animal_class_data"`
	TemperatureLocationData 				[]TemperatureLocationItem   		`json:"temperature_location_data"`
	AnimalNumberData 						[]AnimalNumberItem          		`json:"animal_number_data"`
	EntericFermEFParameterData 				[]EntericFermEFParameterItem 		`json:"enteric_ferm_emission_factor_parameter_data"`
	EntericEmissionFactorData				[]EntericEmissionFactorItem			`json:"enteric_emission_factor_data"`
	EntericEmissionFactorDataUserFriendly	[]EntericEmissionFactorItemDisplay	`json:"enteric_emission_factor_data_user"`
	ErrorRecordData							[]ErrorRecordItem					`json:"error_record"`
	ErrorGenericMsg							[]string							`json:"error_generic" gorm:"-"`
	SheetData								[]SheetDataItem						`json:"-" gorm:"-"`
}

type SimulationResultPackage struct {
	SettingData								[]SettingDataItem					`json:"setting_data"`
	SystemData 								[]SystemDataItem            		`json:"system_data"`
	LocationData 							[]LocationDataItem          		`json:"location_data"`
	AnimalClassData 						[]AnimalClassDataItem       		`json:"animal_class_data"`
	TemperatureLocationData 				[]TemperatureLocationItem   		`json:"temperature_location_data"`
	AnimalNumberData 						[]AnimalNumberItem          		`json:"animal_number_data"`
	EntericFermEFParameterData 				[]EntericFermEFParameterItem 		`json:"enteric_ferm_emission_factor_parameter_data"`
	EntericEmissionFactorData				[]EntericEmissionFactorItem			`json:"enteric_emission_factor_data"`
	EntericEmissionFactorDataUserFriendly	[]EntericEmissionFactorItemDisplay	`json:"enteric_emission_factor_data_user"`
	ErrorRecordData							[]ErrorRecordItem					`json:"error_record"`
	ErrorGenericMsg							[]string							`json:"-" gorm:"-"`
	GoogleAuthenticated						bool								`json:"-" gorm:"-"`
}

type SheetDataItem struct{
	SheetID		string				`json:"sheet_id"`
	SheetName	string				`json:"sheet_name"`
	Created		string				`json:"sheet_created_time"`
}

func (data DataPackage) ClearPackageData() {
	tempErrorMsg := data.ErrorGenericMsg
	tempSheetData := data.SheetData
	data = DataPackage{
		SettingData: 			 nil,
		SystemData:              nil,
		LocationData:            nil,
		AnimalClassData:         nil,
		TemperatureLocationData: nil,
		AnimalNumberData:        nil,
		EntericFermEFParameterData: nil,
		EntericEmissionFactorData: nil,
		ErrorGenericMsg: tempErrorMsg,
		SheetData: tempSheetData,
	}
}

func (data DataPackage)GenerateMappingsForPrimaryTable()(map[interface{}]SystemDataItem, map[interface{}]LocationDataItem, map[interface{}]AnimalClassDataItem){
	systemsMap := make(map[interface{}]SystemDataItem)
	for _,record := range data.SystemData{
		systemsMap[record.Id] = record
		systemsMap[record.Name] = record
	}

	locationsMap := make(map[interface{}]LocationDataItem)
	for _,record := range data.LocationData{
		locationsMap[record.Id] = record
		locationsMap[record.Name] = record
	}

	animalsMap := make(map[interface{}]AnimalClassDataItem)
	for _,record := range data.AnimalClassData{
		animalsMap[record.Id] = record
		animalsMap[record.Name] = record
	}
	return systemsMap, locationsMap, animalsMap
}

func (data DataPackage)ClearSpreadSheetPackage() {
	data.SheetData = nil
}

func OutputSettingDataTable(name string, data DataPackage) {

	logrus.Infof("===========================================================")
	logrus.Infof("Table for: %v", name)

	// SYSTEM output table
	table := simpletable.New()
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "KEY"},
			{Align: simpletable.AlignCenter, Text: "VALUE"},
		},
	}
	for _, row := range data.SettingData {
		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", row.Key)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", row.Value)},
		}

		table.Body.Cells = append(table.Body.Cells, r)
	}
	table.SetStyle(simpletable.StyleDefault)
	fmt.Println(table.String())

}

func OutputSystemDataTable(name string, data DataPackage) {

	logrus.Infof("===========================================================")
	logrus.Infof("Table for: %v", name)

	// SYSTEM output table
	systemTable := simpletable.New()
	systemTable.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "ID"},
			{Align: simpletable.AlignCenter, Text: "NAME"},
		},
	}
	for _, row := range data.SystemData {
		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%d", row.Id)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", row.Name)},
		}

		systemTable.Body.Cells = append(systemTable.Body.Cells, r)
	}
	systemTable.SetStyle(simpletable.StyleDefault)
	fmt.Println(systemTable.String())

}

func OutputLocationDataTable(name string, data DataPackage) {

	logrus.Infof("===========================================================")
	logrus.Infof("Table for: %v", name)

	// SYSTEM output table
	systemTable := simpletable.New()
	systemTable.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "ID"},
			{Align: simpletable.AlignCenter, Text: "NAME"},
		},
	}
	for _, row := range data.LocationData {
		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%d", row.Id)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", row.Name)},
		}

		systemTable.Body.Cells = append(systemTable.Body.Cells, r)
	}
	systemTable.SetStyle(simpletable.StyleDefault)
	fmt.Println(systemTable.String())

}

func OutputAnimalClassDataTable(name string, data DataPackage) {

	logrus.Infof("===========================================================")
	logrus.Infof("Table for: %v", name)

	// SYSTEM output table
	systemTable := simpletable.New()
	systemTable.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "ID"},
			{Align: simpletable.AlignCenter, Text: "PARENT CLASS"},
			{Align: simpletable.AlignCenter, Text: "NAME"},
			{Align: simpletable.AlignCenter, Text: "DEFAULT EF"},
		},
	}
	for _, row := range data.AnimalClassData {
		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%d", row.Id)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%d", row.ParentClass)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", row.Name)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", row.DefaulEF)},
		}

		systemTable.Body.Cells = append(systemTable.Body.Cells, r)
	}
	systemTable.SetStyle(simpletable.StyleDefault)
	fmt.Println(systemTable.String())

}

func OutputTLDataTable(name string, data DataPackage) {

	logrus.Infof("===========================================================")
	logrus.Infof("Table for: %v", name)
	_,locationsMap,_ := data.GenerateMappingsForPrimaryTable()

	fmt.Println(locationsMap[1].Name)

	// SYSTEM output table
	systemTable := simpletable.New()
	systemTable.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "ID"},
			{Align: simpletable.AlignCenter, Text: "LOCATION"},
			{Align: simpletable.AlignCenter, Text: "YEAR"},
			{Align: simpletable.AlignCenter, Text: "MONTH"},
			{Align: simpletable.AlignCenter, Text: "AVG TEMP"},
		},
	}
	for _, row := range data.TemperatureLocationData {
		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%d", row.Id)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", locationsMap[row.Locationid].Name)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", row.Year)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", row.Month)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%.6f", row.AvgTemp)},
		}

		systemTable.Body.Cells = append(systemTable.Body.Cells, r)
	}
	systemTable.SetStyle(simpletable.StyleDefault)
	fmt.Println(systemTable.String())

	//for n, s := range styles {
	//	fmt.Println(n)
	//
	//	systemTable.SetStyle(s)
	//	systemTable.Println()
	//
	//	fmt.Println()
	//}
}

func OutputAnimalNumberDataTable(name string, data DataPackage) {

	logrus.Infof("===========================================================")
	logrus.Infof("Table for: %v", name)
	systemMap, locationsMap, animalsMap := data.GenerateMappingsForPrimaryTable()

	// SYSTEM output table
	systemTable := simpletable.New()
	systemTable.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "ID"},
			{Align: simpletable.AlignCenter, Text: "LOCATION"},
			{Align: simpletable.AlignCenter, Text: "SYSTEM"},
			{Align: simpletable.AlignCenter, Text: "ANIMAL CLASS"},
			{Align: simpletable.AlignCenter, Text: "YEAR"},
			{Align: simpletable.AlignCenter, Text: "MONTH"},
			{Align: simpletable.AlignCenter, Text: "ANIMAL NUMBER"},
		},
	}
	for _, row := range data.AnimalNumberData {
		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%d", row.Id)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", locationsMap[row.Locationid].Name)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", systemMap[row.Systemid].Name)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", animalsMap[row.AnimalClassid].Name)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", row.Year)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", row.Month)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", row.AnimalNumber)},
		}

		systemTable.Body.Cells = append(systemTable.Body.Cells, r)
	}
	systemTable.SetStyle(simpletable.StyleDefault)
	fmt.Println(systemTable.String())
}

func OutputEntericFermEFParameterItemDataTable(name string, data DataPackage) {

	logrus.Infof("===========================================================")
	logrus.Infof("Table for: %v", name)
	systemMap, locationsMap, animalsMap := data.GenerateMappingsForPrimaryTable()

	// SYSTEM output table
	systemTable := simpletable.New()
	systemTable.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "ID"},
			{Align: simpletable.AlignCenter, Text: "LOCATION"},
			{Align: simpletable.AlignCenter, Text: "SYSTEM"},
			{Align: simpletable.AlignCenter, Text: "ANIMAL CLASS"},
			{Align: simpletable.AlignCenter, Text: "YEAR"},
			{Align: simpletable.AlignCenter, Text: "MONTH"},
			{Align: simpletable.AlignCenter, Text: "BODY WEIGHT"},
			{Align: simpletable.AlignCenter, Text: "MATURE WEIGHT"},
			{Align: simpletable.AlignCenter, Text: "DAILY WEIGHT GAIN"},
			{Align: simpletable.AlignCenter, Text: "FRACTION OF MONTH ALIVE"},
			{Align: simpletable.AlignCenter, Text: "CF"},
			{Align: simpletable.AlignCenter, Text: "C"},
			{Align: simpletable.AlignCenter, Text: "CA"},
			{Align: simpletable.AlignCenter, Text: "MILK PROD"},
			{Align: simpletable.AlignCenter, Text: "FAT CONTENT"},
			{Align: simpletable.AlignCenter, Text: "C PREGNANCY"},
			{Align: simpletable.AlignCenter, Text: "PROPORTION ANIMAL CLASS PREGNANT"},
			{Align: simpletable.AlignCenter, Text: "PROPORTION ANIMAL CLASS LACTATING"},
			{Align: simpletable.AlignCenter, Text: "FRACTION OF MONTH LACTATING"},
			{Align: simpletable.AlignCenter, Text: "DE"},
			{Align: simpletable.AlignCenter, Text: "YM"},
		},
	}
	for _, row := range data.EntericFermEFParameterData {
		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%d", row.Id)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", locationsMap[row.Locationid].Name)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", systemMap[row.Systemid].Name)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", animalsMap[row.AnimalClassid].Name)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", row.Year)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", row.Month)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", row.BodyWeight)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", row.MatureWeight)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", row.DailyWeightGain)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", row.FractionOfMonthAlive)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", row.CF)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", row.C)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", row.CA)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", row.MilkProd)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", row.FatContent)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", row.CPregnancy)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", row.ProportionAnimalClassPregnant)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", row.ProportionAnimalClassLactating)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", row.FractionOfMonthLactating)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", row.DE)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", row.YM)},
		}

		systemTable.Body.Cells = append(systemTable.Body.Cells, r)
	}
	systemTable.SetStyle(simpletable.StyleDefault)
	fmt.Println(systemTable.String())
}

func OutputEntericEmissionFactorItemDataTable(name string, data DataPackage) {

	logrus.Infof("===========================================================")
	logrus.Infof("Table for: %v", name)
	systemMap, locationsMap, animalsMap := data.GenerateMappingsForPrimaryTable()

	// SYSTEM output table
	systemTable := simpletable.New()
	systemTable.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "ID"},
			{Align: simpletable.AlignCenter, Text: "LOCATION"},
			{Align: simpletable.AlignCenter, Text: "SYSTEM"},
			{Align: simpletable.AlignCenter, Text: "ANIMAL CLASS"},
			{Align: simpletable.AlignCenter, Text: "YEAR"},
			{Align: simpletable.AlignCenter, Text: "MONTH"},
			{Align: simpletable.AlignCenter, Text: "CALCULATED EF"},
			{Align: simpletable.AlignCenter, Text: "AAP"},
			{Align: simpletable.AlignCenter, Text: "COF CALCULATION NE MAINTENANCE"},
			{Align: simpletable.AlignCenter, Text: "NE MAINTENANCE"},
			{Align: simpletable.AlignCenter, Text: "NE ACTIVITY CATTLE BUFFALO"},
			{Align: simpletable.AlignCenter, Text: "NE GROWTH CATTLE BUFFALO"},
			{Align: simpletable.AlignCenter, Text: "NE LACTATION BEEF DAIRY BUFFALO"},
			{Align: simpletable.AlignCenter, Text: "NE WORK CATTLE BUFFALO"},
			{Align: simpletable.AlignCenter, Text: "NE PRENANCY CATTLE BUFFALO SHEEP"},
			{Align: simpletable.AlignCenter, Text: "RNE AVAILABLE DIET MAINTENANCE DIGESTIBLE ENERGY CONSUMED"},
			{Align: simpletable.AlignCenter, Text: "RNE Available FOR GROWTH DIET DIGESTIBLE ENERGY CONSUMED"},
			{Align: simpletable.AlignCenter, Text: "GE CATTLE BUFFALO SHEEP"},
			{Align: simpletable.AlignCenter, Text: "ENTERIC FERMENTATION EMISSIONS LIVESTOCK CATEGORY"},
		},
	}
	for _, row := range data.EntericEmissionFactorData {
		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%d", row.Id)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", locationsMap[row.Locationid].Name)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", systemMap[row.Systemid].Name)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", animalsMap[row.AnimalClassid].Name)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", row.Year)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", row.Month)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", row.CalculatedEF)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", row.MAP)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", row.CofCalculationNEMaintenance)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", row.NEMaintenance)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", row.NEActivityCattleBuffalo)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", row.NEGrowthCattleBuffalo)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", row.NELactationBeefDairyBuffalo)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", row.NEWorkCattleBuffalo)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", row.NEPregnancyCattleBuffaloSheep)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", row.RNEAvailableDietMaintenanceDigestibleEnergyConsumed)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", row.RNEAvailableForGrowthDietDisgestibleEnergyConsumed)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", row.GECattleBuffaloSheep)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", row.EntericFermentationEmissionsLivestockCategory)},
		}

		systemTable.Body.Cells = append(systemTable.Body.Cells, r)
	}
	systemTable.SetStyle(simpletable.StyleDefault)
	fmt.Println(systemTable.String())
}


func OutputSpreadSheetsList(name string, data DataPackage) {

	logrus.Infof("===========================================================")
	logrus.Infof("Table for: %v", name)

	// SYSTEM output table
	table := simpletable.New()
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "SHEET ID"},
			{Align: simpletable.AlignCenter, Text: "SHEET NAME"},
		},
	}
	for _, row := range data.SheetData {
		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", row.SheetID)},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%v", row.SheetName)},
		}

		table.Body.Cells = append(table.Body.Cells, r)
	}
	table.SetStyle(simpletable.StyleDefault)
	fmt.Println(table.String())

}

func OutputAllTables(data DataPackage) {
	// SYSTEM output table
	fmt.Printf("Data for Sheet (%v)", data.SheetID)
	OutputSettingDataTable("Settings", data)
	OutputSystemDataTable("System", data)
	OutputLocationDataTable("Location", data)
	OutputAnimalClassDataTable("Animal Class", data)
	OutputTLDataTable("Temperature Location", data)
	OutputAnimalNumberDataTable("Animal Numbers", data)
	OutputEntericFermEFParameterItemDataTable("Enteric Ferm FE Parameters", data)
	OutputEntericEmissionFactorItemDataTable("Enteric Emission Factors", data)
}

func OutputAllSpreadSheets(data DataPackage){
	fmt.Printf("Total sheets (%v)", len(data.SheetData))
	OutputSpreadSheetsList("List of Spreadsheets", data)
}

