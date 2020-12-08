package helpers

import (
	"context"
	"fmt"
	model "github.com/MullionGroup/go-website-flintpro-example/models"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"time"
)

type IGormClient interface {
	QuerySettingData				(ctx context.Context) ([]model.SettingDataItem, error)
	QuerySystemData					(ctx context.Context) ([]model.SystemDataItem, error)
	QueryLocationData				(ctx context.Context) ([]model.LocationDataItem, error)
	QueryAnimalClassData			(ctx context.Context) ([]model.AnimalClassDataItem, error)
	QueryTemperatureLocationData	(ctx context.Context) ([]model.TemperatureLocationItem, error)
	QueryAnimalNumberData			(ctx context.Context) ([]model.AnimalNumberItem, error)
	QueryEntericFermEFParameterData	(ctx context.Context) ([]model.EntericFermEFParameterItem, error)
	QueryEntericEmissionFactorData	(ctx context.Context) ([]model.EntericEmissionFactorItem, error)

	SeedSettingDataItem				(item model.SettingDataItem) error
	SeedSystemDataItem				(item model.SystemDataItem) error
	SeedLocationDataItem			(item model.LocationDataItem) error
	SeedAnimalClassDataItem			(item model.AnimalClassDataItem) error
	SeedTemperatureLocationItem		(item model.TemperatureLocationItem) error
	SeedAnimalNumberItem			(item model.AnimalNumberItem) error
	SeedEntericFermEFParameterItem	(item model.EntericFermEFParameterItem) error
	SeedEntericEmissionFactorItem	(item model.EntericEmissionFactorItem) error
	SeedErrorRecordItem				(item model.ErrorRecordItem) error

	OpenDB(addr string, debug bool)
	SetupDB(addr string, debug bool)
	DeleteTableData() error
	Check() bool
	Close()
}

var(
	dbMaxopenconns = 10
	dbConnmaxlifetime_hours = 1
)


var DBClient IGormClient = &GormClient{}

// =====================================================================================================================
// Utilities
// =====================================================================================================================

func LoadSQLiteIntoDataPackage(data model.DataPackage, filename string) (model.DataPackage, error) {
	data.ClearPackageData()

	DBClient.SetupDB(filename, false)

	ctx := context.Background()
	data.SettingData, _ 			= DBClient.QuerySettingData(ctx)
	data.SystemData, _ 				= DBClient.QuerySystemData(ctx)
	data.LocationData, 	_ 			= DBClient.QueryLocationData(ctx)
	data.AnimalClassData, 	_ 		= DBClient.QueryAnimalClassData(ctx)
	data.TemperatureLocationData,_ 	= DBClient.QueryTemperatureLocationData(ctx)
	data.AnimalNumberData, 	_ 		= DBClient.QueryAnimalNumberData(ctx)
	data.EntericFermEFParameterData, _ = DBClient.QueryEntericFermEFParameterData(ctx)
	data.EntericEmissionFactorData, _ = DBClient.QueryEntericEmissionFactorData(ctx)

	systemMap, locationsMap, animalsMap := data.GenerateMappingsForPrimaryTable()

	for index,_ := range data.TemperatureLocationData {
		data.TemperatureLocationData[index].Location = locationsMap[data.TemperatureLocationData[index].Locationid].Name
	}

	for index,_ := range data.AnimalNumberData {
		data.AnimalNumberData[index].Location = locationsMap[data.AnimalNumberData[index].Locationid].Name
		data.AnimalNumberData[index].System = systemMap[data.AnimalNumberData[index].Systemid].Name
		data.AnimalNumberData[index].AnimalClass = animalsMap[data.AnimalNumberData[index].AnimalClassid].Name
	}

	for index,_ := range data.EntericFermEFParameterData {
		data.EntericFermEFParameterData[index].Location = locationsMap[data.EntericFermEFParameterData[index].Locationid].Name
		data.EntericFermEFParameterData[index].System = systemMap[data.EntericFermEFParameterData[index].Systemid].Name
		data.EntericFermEFParameterData[index].AnimalClass = animalsMap[data.EntericFermEFParameterData[index].AnimalClassid].Name
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

	for index,_ := range data.EntericEmissionFactorData {
		data.EntericEmissionFactorData[index].Location = locationsMap[data.EntericEmissionFactorData[index].Locationid].Name
		data.EntericEmissionFactorData[index].System = systemMap[data.EntericEmissionFactorData[index].Systemid].Name
		data.EntericEmissionFactorData[index].AnimalClass = animalsMap[data.EntericEmissionFactorData[index].AnimalClassid].Name

		data.EntericEmissionFactorDataUserFriendly = append(data.EntericEmissionFactorDataUserFriendly, model.EntericEmissionFactorItemDisplay{
			Location:                      data.EntericEmissionFactorData[index].Location,
			System:                        data.EntericEmissionFactorData[index].System,
			ParentClass:				   data.EntericEmissionFactorData[index].ParentClass,
			AnimalClass:                   data.EntericEmissionFactorData[index].AnimalClass,
			Year:                          data.EntericEmissionFactorData[index].Year,
			Month:                         data.EntericEmissionFactorData[index].Month,
			CalculatedEF:                  data.EntericEmissionFactorData[index].CalculatedEF,
			MAP:                           data.EntericEmissionFactorData[index].MAP,
			CofCalculationNEMaintenance:   data.EntericEmissionFactorData[index].CofCalculationNEMaintenance,
			NEMaintenance:                 data.EntericEmissionFactorData[index].NEMaintenance,
			NEActivityCattleBuffalo:       data.EntericEmissionFactorData[index].NEActivityCattleBuffalo,
			NEGrowthCattleBuffalo:         data.EntericEmissionFactorData[index].NEGrowthCattleBuffalo,
			NELactationBeefDairyBuffalo:   data.EntericEmissionFactorData[index].NELactationBeefDairyBuffalo,
			NEWorkCattleBuffalo:           data.EntericEmissionFactorData[index].NEWorkCattleBuffalo,
			NEPregnancyCattleBuffaloSheep: data.EntericEmissionFactorData[index].NEPregnancyCattleBuffaloSheep,
			RNEAvailableDietMaintenanceDigestibleEnergyConsumed: data.EntericEmissionFactorData[index].RNEAvailableDietMaintenanceDigestibleEnergyConsumed,
			RNEAvailableForGrowthDietDisgestibleEnergyConsumed:  data.EntericEmissionFactorData[index].RNEAvailableForGrowthDietDisgestibleEnergyConsumed,
			GECattleBuffaloSheep:                          data.EntericEmissionFactorData[index].GECattleBuffaloSheep,
			EntericFermentationEmissionsLivestockCategory: data.EntericEmissionFactorData[index].EntericFermentationEmissionsLivestockCategory,
		})
	}

	defer DBClient.Close()
	return data, nil
}

func LoadDataPackageIntoSQLite(data model.DataPackage) ([]byte, error){
	tmpDB, err := ioutil.TempFile("", data.SheetID+".sqlite3")
	if err != nil {
		fmt.Printf("Create temp db with error (%v)", err.Error())
		return []byte{}, err
	}

	DBClient.SetupDB(tmpDB.Name(), false)
	DBClient.DeleteTableData()

	for _,item := range data.SettingData {
		DBClient.SeedSettingDataItem(item)
	}
	for _,item := range data.SystemData {
		DBClient.SeedSystemDataItem(item)
	}
	for _,item := range data.LocationData {
		DBClient.SeedLocationDataItem(item)
	}
	for _,item := range data.AnimalClassData {
		DBClient.SeedAnimalClassDataItem(item)
	}
	for _,item := range data.TemperatureLocationData {
		DBClient.SeedTemperatureLocationItem(item)
	}
	for _,item := range data.AnimalNumberData {
		DBClient.SeedAnimalNumberItem(item)
	}
	for _,item := range data.EntericFermEFParameterData {
		DBClient.SeedEntericFermEFParameterItem(item)
	}
	for _,item := range data.EntericEmissionFactorData {
		DBClient.SeedEntericEmissionFactorItem(item)
	}
	for _,item := range data.ErrorRecordData {
		DBClient.SeedErrorRecordItem(item)
	}

	defer DBClient.Close()

	// Remove temp file after on exit.
	defer func() {
		tmpDB.Close()
		err := os.Remove(tmpDB.Name())
		if err != nil {
			log.Print(err)
		}
	}()

	sqlByteData, err := ioutil.ReadFile(tmpDB.Name())
	if err != nil {
		fmt.Printf("Read temp db with error (%v)", err.Error())
		return []byte{}, err
	}
	return sqlByteData, nil
}



// =====================================================================================================================
// GORM Client
// =====================================================================================================================

type GormClient struct {
	crDB *gorm.DB
}

func (gc *GormClient) SetupDB(addr string, debug bool) {
	logrus.Infof("Connecting with connection string: '%v'", addr)
	var err error
	gc.crDB, err = gorm.Open("sqlite3", addr)
	if err != nil {
		// pause and try again
		time.Sleep(5 * time.Second)

		gc.crDB, err = gorm.Open("sqlite3", addr)
		if err != nil {
			panic("failed to connect database: " + err.Error())
		}
	}
	err = gc.crDB.DB().Ping() // Send a ping to make sure the database connection is alive.
	if err != nil {
		panic("failed to ping database: " + err.Error())
	}

	// Too much logging for now
	//gc.crDB.LogMode(debug)

	// Migrate the schema
	gc.crDB.AutoMigrate(
		&model.SettingDataItem{},
		&model.SystemDataItem{},
		&model.LocationDataItem{},
		&model.AnimalClassDataItem{},
		&model.TemperatureLocationItem{},
		&model.AnimalNumberItem{},
		&model.EntericFermEFParameterItem{},
		&model.EntericEmissionFactorItem{},
		&model.ErrorRecordItem{},
	)
}

func (gc *GormClient) DeleteTableData() error {
	if !gc.Check() {
		return fmt.Errorf("Connection error")
	}

	tx := gc.crDB.Begin()
	tx = tx.Delete(&model.SettingDataItem{})
	if tx.Error != nil {
		return fmt.Errorf("Server error")
	}
	tx.Commit()

	tx = gc.crDB.Begin()
	tx = tx.Delete(&model.SystemDataItem{})
	if tx.Error != nil {
		return fmt.Errorf("Server error")
	}
	tx.Commit()

	tx = gc.crDB.Begin()
	tx = tx.Delete(&model.LocationDataItem{})
	if tx.Error != nil {
		return fmt.Errorf("Server error")
	}
	tx.Commit()

	tx = gc.crDB.Begin()
	tx = tx.Delete(&model.AnimalClassDataItem{})
	if tx.Error != nil {
		return fmt.Errorf("Server error")
	}
	tx.Commit()

	tx = gc.crDB.Begin()
	tx = tx.Delete(&model.TemperatureLocationItem{})
	if tx.Error != nil {
		return fmt.Errorf("Server error")
	}
	tx.Commit()

	tx = gc.crDB.Begin()
	tx = tx.Delete(&model.AnimalNumberItem{})
	if tx.Error != nil {
		return fmt.Errorf("Server error")
	}
	tx.Commit()

	tx = gc.crDB.Begin()
	tx = tx.Delete(&model.EntericFermEFParameterItem{})
	if tx.Error != nil {
		return fmt.Errorf("Server error")
	}
	tx.Commit()

	tx = gc.crDB.Begin()
	tx = tx.Delete(&model.EntericEmissionFactorItem{})
	if tx.Error != nil {
		return fmt.Errorf("Server error")
	}
	tx.Commit()

	tx = gc.crDB.Begin()
	tx = tx.Delete(&model.ErrorRecordItem{})
	if tx.Error != nil {
		return fmt.Errorf("Server error")
	}
	tx.Commit()

	if tx.Error != nil {
		return fmt.Errorf("Server error")
	}
	return nil
}

func (gc *GormClient) OpenDB(addr string, debug bool) {
	logrus.Debugf("Connecting with connection string: '%v'", addr)
	var err error
	gc.crDB, err = gorm.Open("sqlite3", addr)
	if err != nil {
		// pause and try again
		time.Sleep(5 * time.Second)

		gc.crDB, err = gorm.Open("sqlite3", addr)
		if err != nil {
			panic("failed to connect database: " + err.Error())
		}
	}
	err = gc.crDB.DB().Ping() // Send a ping to make sure the database connection is alive.
	if err != nil {
		panic("failed to ping database: " + err.Error())
	}
}


func (gc *GormClient) Check() bool {
	logrus.Debug("Checking connection to DB")
	return gc.crDB != nil
}

func (gc *GormClient) Close() {
	logrus.Info("Closing connection to DB")
	gc.crDB.Close()
}

func (gc *GormClient) QuerySettingData(ctx context.Context) ([]model.SettingDataItem, error) {
	acc := []model.SettingDataItem{}

	if !gc.Check() {
		return acc, fmt.Errorf("Connection error")
	}

	tx := gc.crDB.Begin()

	tx = tx.
		Find(&acc).
		Order("id");

	if tx.Error != nil {
		return acc, fmt.Errorf("Server error")
	}

	tx.Commit()
	if tx.Error != nil {
		return acc, fmt.Errorf("Server error")
	}
	return acc , nil
}

func (gc *GormClient) QuerySystemData(ctx context.Context) ([]model.SystemDataItem, error) {
	acc := []model.SystemDataItem{}

	if !gc.Check() {
		return acc, fmt.Errorf("Connection error")
	}

	tx := gc.crDB.Begin()

	tx = tx.
		Find(&acc).
		Order("id");

	if tx.Error != nil {
		return acc, fmt.Errorf("Server error")
	}

	tx.Commit()
	if tx.Error != nil {
		return acc, fmt.Errorf("Server error")
	}
	return acc , nil
}

func (gc *GormClient) QueryLocationData(ctx context.Context) ([]model.LocationDataItem, error) {
	acc := []model.LocationDataItem{}

	if !gc.Check() {
		return acc, fmt.Errorf("Connection error")
	}

	tx := gc.crDB.Begin()

	tx = tx.
		Find(&acc).
		Order("id");

	if tx.Error != nil {
		return acc, fmt.Errorf("Server error")
	}

	tx.Commit()
	if tx.Error != nil {
		return acc, fmt.Errorf("Server error")
	}
	return acc , nil
}

func (gc *GormClient) QueryAnimalClassData(ctx context.Context) ([]model.AnimalClassDataItem, error) {
	acc := []model.AnimalClassDataItem{}

	if !gc.Check() {
		return acc, fmt.Errorf("Connection error")
	}

	tx := gc.crDB.Begin()

	tx = tx.
		Find(&acc).
		Order("id");

	if tx.Error != nil {
		return acc, fmt.Errorf("Server error")
	}

	tx.Commit()
	if tx.Error != nil {
		return acc, fmt.Errorf("Server error")
	}
	return acc , nil
}

func (gc *GormClient) QueryTemperatureLocationData(ctx context.Context) ([]model.TemperatureLocationItem, error) {
	acc := []model.TemperatureLocationItem{}

	if !gc.Check() {
		return acc, fmt.Errorf("Connection error")
	}

	tx := gc.crDB.Begin()

	tx = tx.
		Find(&acc).
		Order("id");

	if tx.Error != nil {
		return acc, fmt.Errorf("Server error")
	}

	tx.Commit()
	if tx.Error != nil {
		return acc, fmt.Errorf("Server error")
	}
	return acc , nil
}

func (gc *GormClient) QueryAnimalNumberData(ctx context.Context) ([]model.AnimalNumberItem, error) {
	acc := []model.AnimalNumberItem{}

	if !gc.Check() {
		return acc, fmt.Errorf("Connection error")
	}

	tx := gc.crDB.Begin()

	tx = tx.
		Find(&acc).
		Order("id");

	if tx.Error != nil {
		return acc, fmt.Errorf("Server error")
	}

	tx.Commit()
	if tx.Error != nil {
		return acc, fmt.Errorf("Server error")
	}
	return acc , nil
}

func (gc *GormClient) QueryEntericFermEFParameterData(ctx context.Context) ([]model.EntericFermEFParameterItem, error) {
	acc := []model.EntericFermEFParameterItem{}

	if !gc.Check() {
		return acc, fmt.Errorf("Connection error")
	}

	tx := gc.crDB.Begin()

	tx = tx.
		Find(&acc)

	if tx.Error != nil {
		return acc, fmt.Errorf("Server error")
	}

	tx.Commit()
	if tx.Error != nil {
		return acc, fmt.Errorf("Server error")
	}
	return acc , nil
}

func (gc *GormClient) QueryEntericEmissionFactorData(ctx context.Context) ([]model.EntericEmissionFactorItem, error) {
	acc := []model.EntericEmissionFactorItem{}

	if !gc.Check() {
		return acc, fmt.Errorf("Connection error")
	}

	tx := gc.crDB.Begin()

	tx = tx.
		Find(&acc)

	if tx.Error != nil {
		return acc, fmt.Errorf("Server error")
	}

	tx.Commit()
	if tx.Error != nil {
		return acc, fmt.Errorf("Server error")
	}
	return acc , nil
}

func (gc *GormClient) SeedSettingDataItem(item model.SettingDataItem) error {
	tx := gc.crDB.Begin()
	tx = tx.Create(&item)
	if tx.Error != nil {
		tx.Rollback()
		return fmt.Errorf("Create Error in seeding SettingData: %v", tx.Error.Error())
	}
	tx = tx.Commit()
	if tx.Error != nil {
		return fmt.Errorf("Commit Error in seeding SettingData: %v", tx.Error.Error())
	}
	return nil
}

func (gc *GormClient) SeedSystemDataItem(item model.SystemDataItem) error {
	tx := gc.crDB.Begin()
	tx = tx.Create(&item)
	if tx.Error != nil {
		tx.Rollback()
		return fmt.Errorf("Create Error in seeding SystemDataItem: %v", tx.Error.Error())
	}
	tx = tx.Commit()
	if tx.Error != nil {
		return fmt.Errorf("Commit Error in seeding SystemDataItem: %v", tx.Error.Error())
	}
	return nil
}

func (gc *GormClient) SeedLocationDataItem(item model.LocationDataItem) error {
	tx := gc.crDB.Begin()
	tx = tx.Create(&item)
	if tx.Error != nil {
		tx.Rollback()
		return fmt.Errorf("Create Error in seeding LocationDataItem: %v", tx.Error.Error())
	}
	tx = tx.Commit()
	if tx.Error != nil {
		return fmt.Errorf("Commit Error in seeding LocationDataItem: %v", tx.Error.Error())
	}
	return nil
}

func (gc *GormClient) SeedAnimalClassDataItem(item model.AnimalClassDataItem) error {
	tx := gc.crDB.Begin()
	tx = tx.Create(&item)
	if tx.Error != nil {
		tx.Rollback()
		return fmt.Errorf("Create Error in seeding AnimalClassDataItem: %v", tx.Error.Error())
	}
	tx = tx.Commit()
	if tx.Error != nil {
		return fmt.Errorf("Commit Error in seeding AnimalClassDataItem: %v", tx.Error.Error())
	}
	return nil
}

func (gc *GormClient) SeedTemperatureLocationItem(item model.TemperatureLocationItem) error {
	tx := gc.crDB.Begin()
	tx = tx.Create(&item)
	if tx.Error != nil {
		tx.Rollback()
		return fmt.Errorf("Create Error in seeding TemperatureLocationItem: %v", tx.Error.Error())
	}
	tx = tx.Commit()
	if tx.Error != nil {
		return fmt.Errorf("Commit Error in seeding TemperatureLocationItem: %v", tx.Error.Error())
	}
	return nil
}

func (gc *GormClient) SeedAnimalNumberItem(item model.AnimalNumberItem) error {
	tx := gc.crDB.Begin()
	tx = tx.Create(&item)
	if tx.Error != nil {
		tx.Rollback()
		return fmt.Errorf("Create Error in seeding AnimalNumberItem: %v", tx.Error.Error())
	}
	tx = tx.Commit()
	if tx.Error != nil {
		return fmt.Errorf("Commit Error in seeding AnimalNumberItem: %v", tx.Error.Error())
	}
	return nil
}

func (gc *GormClient) SeedEntericFermEFParameterItem(item model.EntericFermEFParameterItem) error {
	tx := gc.crDB.Begin()
	tx = tx.Create(&item)
	if tx.Error != nil {
		tx.Rollback()
		return fmt.Errorf("Create Error in seeding EntericFermEFParameterItem: %v", tx.Error.Error())
	}
	tx = tx.Commit()
	if tx.Error != nil {
		return fmt.Errorf("Commit Error in seeding EntericFermEFParameterItem: %v", tx.Error.Error())
	}
	return nil
}

func (gc *GormClient) SeedEntericEmissionFactorItem(item model.EntericEmissionFactorItem) error {
	tx := gc.crDB.Begin()
	tx = tx.Create(&item)
	if tx.Error != nil {
		tx.Rollback()
		return fmt.Errorf("Create Error in seeding EntericEmissionFactorItem: %v", tx.Error.Error())
	}
	tx = tx.Commit()
	if tx.Error != nil {
		return fmt.Errorf("Commit Error in seeding EntericEmissionFactorItem: %v", tx.Error.Error())
	}
	return nil
}

func (gc *GormClient) SeedErrorRecordItem(item model.ErrorRecordItem) error {
	tx := gc.crDB.Begin()
	tx = tx.Create(&item)
	if tx.Error != nil {
		tx.Rollback()
		return fmt.Errorf("Create Error in seeding ErrorRecordItem: %v", tx.Error.Error())
	}
	tx = tx.Commit()
	if tx.Error != nil {
		return fmt.Errorf("Commit Error in seeding ErrorRecordItem: %v", tx.Error.Error())
	}
	return nil
}

