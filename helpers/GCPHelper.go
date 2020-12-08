package helpers

import (
	"encoding/json"
	"fmt"
	model "github.com/MullionGroup/go-website-flintpro-example/models"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/sheets/v4"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"
)

var (
	//Define expected posiiton of data start in Google Sheets
	ReadRangeSettings 				= "Settings!A5:B"
	ReadRangeSystem 				= "System!A5:B"
	ReadRangeLocation 				= "Location!A5:B"
	ReadRangeAnimalClass 			= "AnimalClass!A5:D"
	ReadRangeTemperatureLocation 	= "TemperatureLocation!A5:E"
	ReadRangeAnimalNumber 			= "AnimalNumbers!A5:G"
	ReadRangeEntericFermEFParameter = "EntericFermEFParameters!A5:V"
	ReadRangeEntericEmissionFactor	= "EntericEmissionFactors!A5:S"

	WriteRangeSettings 				= "Settings!A4"
	WriteRangeSystem 				= "System!A4"
	WriteRangeLocation 				= "Location!A4"
	WriteRangeAnimalClass 			= "AnimalClass!A4"
	WriteRangeTemperatureLocation 	= "TemperatureLocation!A4"
	WriteRangeAnimalNumber 			= "AnimalNumbers!A4"
	WriteRangeEntericFermEFParameter = "EntericFermEFParameters!A4"
	WriteRangeEntericEmissionFactor	= "EntericEmissionFactors!A4"
	WriteRangeErrorRecord			= "ErrorRecords!A4"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config, tokenString string) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := TokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		SaveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

func getClientPersonal(config *oauth2.Config, token oauth2.Token) *http.Client {
	if len(token.AccessToken) == 0 {
		log.Fatalf("Unable to retrieve token")
	}
	return config.Client(context.Background(), &token)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func getPermission(request *http.Request, token oauth2.Token) (*http.Client,error){
	b := viper.Get("credentials_file").([]byte)

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, sheets.SpreadsheetsScope, drive.DriveFileScope, drive.DriveMetadataScope)
	if err != nil {
		log.Printf("Unable to parse client secret file to config: %v", err)
		return nil, err
	}

	client := getClientPersonal(config, token)
	return client, nil
}

// Retrieves a token from a local file.
func TokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func SaveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func GCPServerAuthSetup(response http.ResponseWriter, request *http.Request, path string) (*oauth2.Config) {
	b := viper.Get("credentials_file").([]byte)

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, sheets.SpreadsheetsScope, drive.DriveFileScope, drive.DriveMetadataScope)
	config.RedirectURL = path
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	return config
}

func LoadEFDataFromService(srv *sheets.Service, spreadsheetId string) (model.DataPackage, error) {
	data := model.DataPackage{}
	// SETTINGS
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, ReadRangeSettings).Do()
	if err != nil {
		log.Printf("Unable to retrieve data from sheet: %v", err)
		return model.DataPackage{}, err
	}

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		for _, row := range resp.Values {
			key := row[0].(string)
			value := row[1].(string)
			data.SettingData = append(data.SettingData, model.SettingDataItem{key, value})
		}
	}

	// SYSTEM
	resp, err = srv.Spreadsheets.Values.Get(spreadsheetId, ReadRangeSystem).Do()
	if err != nil {
		log.Printf("Unable to retrieve data from sheet: %v", err)
		return model.DataPackage{}, err
	}

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		//fmt.Println("System (id):")
		defaultValue := []interface{}{"0",""}
		for _, row := range resp.Values {
			if len(row) < 2{
				row = append(row, defaultValue[len(row):]...)
			}
			id,_ := strconv.Atoi(row[0].(string))
			name := row[1].(string)
			data.SystemData = append(data.SystemData, model.SystemDataItem{id, name})
			//fmt.Printf("%s (%v)\n", row[1], row[0])
		}
	}

	// LOCATION
	//locationData := []model.SimpleDataItem{}
	resp, err = srv.Spreadsheets.Values.Get(spreadsheetId, ReadRangeLocation).Do()
	if err != nil {
		log.Printf("Unable to retrieve data from sheet: %v", err)
		return model.DataPackage{}, err
	}

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		//fmt.Println("Location (id):")
		defaultValue := []interface{}{"0",""}
		for _, row := range resp.Values {
			if len(row) < 2{
				row = append(row, defaultValue[len(row):]...)
			}
			id,_ := strconv.Atoi(row[0].(string))
			name := row[1].(string)
			data.LocationData = append(data.LocationData, model.LocationDataItem{id, name})
			//fmt.Printf("%s (%v)\n", row[1], row[0])
		}
	}

	// ANIMAL CLASS
	//animalClassData := []model.SimpleDataItem{}
	resp, err = srv.Spreadsheets.Values.Get(spreadsheetId, ReadRangeAnimalClass).Do()
	if err != nil {
		log.Printf("Unable to retrieve data from sheet: %v", err)
		return model.DataPackage{}, err
	}

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		//fmt.Println("Animal Class (id):")
		defaultValue := []interface{}{"0","","",""}
		for _, row := range resp.Values {
			if len(row) < 4{
				row = append(row, defaultValue[len(row):]...)
			}
			id,_ := strconv.Atoi(row[0].(string))
			parentClass := row[1].(string)
			name := row[2].(string)
			defaultEF,_ := strconv.ParseFloat(row[3].(string), 64)
			data.AnimalClassData = append(data.AnimalClassData, model.AnimalClassDataItem{id, parentClass, name, defaultEF})
		}
	}

	systemMap, locationsMap, animalsMap := data.GenerateMappingsForPrimaryTable()

	// Temperature Location
	resp, err = srv.Spreadsheets.Values.Get(spreadsheetId, ReadRangeTemperatureLocation).Do()
	if err != nil {
		log.Printf("Unable to retrieve data from sheet: %v", err)
		return model.DataPackage{}, err
	}

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		//fmt.Println("Temperature Location (id):")
		defaultValue := []interface{}{"0","","0","0","0.0","0.0"}
		for _, row := range resp.Values {
			if len(row) < 5{
				row = append(row, defaultValue[len(row):]...)
			}
			id,_ := strconv.Atoi(row[0].(string))
			location := row[1].(string)
			year,_ := strconv.Atoi(row[2].(string))
			month,_ := strconv.Atoi(row[3].(string))
			avgTemp, _ := strconv.ParseFloat(row[4].(string), 64)
			data.TemperatureLocationData = append(data.TemperatureLocationData, model.TemperatureLocationItem{id, locationsMap[location].Id, locationsMap[location].Name, year, month, avgTemp})
			//fmt.Printf("%s (%v)\n", row[1], row[0])
		}
	}

	// Animal Class Numbers
	resp, err = srv.Spreadsheets.Values.Get(spreadsheetId, ReadRangeAnimalNumber).Do()
	if err != nil {
		log.Printf("Unable to retrieve data from sheet: %v", err)
		return model.DataPackage{}, err
	}

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		defaultValue := []interface{}{"0","","","","0","0","0"}
		for _, row := range resp.Values {
			if len(row) < 7{
				row = append(row, defaultValue[len(row):]...)
			}
			id,_ := strconv.Atoi(row[0].(string))
			location 		:= row[1].(string)
			system 			:= row[2].(string)
			animalClass 	:= row[3].(string)
			year,_ 			:= strconv.Atoi(row[4].(string))
			month,_ 		:= strconv.Atoi(row[5].(string))
			animalNumber,_ 	:= strconv.ParseFloat(row[6].(string),64)
			data.AnimalNumberData = append(data.AnimalNumberData, model.AnimalNumberItem{id, locationsMap[location].Id, locationsMap[location].Name, systemMap[system].Id, systemMap[system].Name,  animalsMap[animalClass].Id, animalsMap[animalClass].Name, year, month, animalNumber})
		}
	}

	// Enteric Ferm EF Parameters
	resp, err = srv.Spreadsheets.Values.Get(spreadsheetId, ReadRangeEntericFermEFParameter).Do()
	if err != nil {
		log.Printf("Unable to retrieve data from sheet: %v", err)
		return model.DataPackage{}, err
	}

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		defaultValue := []interface{}{"0","0","0","","","","0.0","0.0","0.0","0","0.0","0.0","0.0","0.0","0.0","0.0","0.0","0.0","0","0","0.0","0.0"}
		for _, row := range resp.Values {
			if len(row) < 22{
				row = append(row, defaultValue[len(row):]...)
			}
			id,_ 				:= strconv.Atoi(row[0].(string))
			year,_ 				:= strconv.Atoi(row[1].(string))
			month,_				:= strconv.Atoi(row[2].(string))
			location 			:= row[3].(string)
			system 				:= row[4].(string)
			animalClass 		:= row[5].(string)
			bodyWeight, _ 		:= strconv.ParseFloat(row[6].(string), 64)
			matureWeight, _ 	:= strconv.ParseFloat(row[7].(string), 64)
			dailyWeightGain, _ 	:= strconv.ParseFloat(row[8].(string), 64)
			fractionOfMonthAlive,_ := strconv.ParseFloat(row[9].(string), 64)
			cf, _ 				:= strconv.ParseFloat(row[10].(string), 64)
			c, _ 				:= strconv.ParseFloat(row[11].(string), 64)
			ca, _ 				:= strconv.ParseFloat(row[12].(string), 64)
			milkProd,_			:= strconv.ParseFloat(row[13].(string), 48)
			fatContent,_		:= strconv.ParseFloat(row[14].(string), 64)
			cPregnancy,_		:= strconv.ParseFloat(row[15].(string), 64)
			proportionAnimalClassPregnant,_	:= strconv.ParseFloat(row[16].(string), 64)
			proportionAnimalClassLactating,_:= strconv.ParseFloat(row[17].(string), 64)
			fractionOfMonthLactating,_		:= strconv.ParseFloat(row[18].(string), 64)
			hoursWorked,_					:= strconv.ParseFloat(row[19].(string), 64)
			de,_							:= strconv.ParseFloat(row[20].(string), 64)
			ym,_							:= strconv.ParseFloat(row[21].(string), 64)
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
	}

	// Enteric Emission Factors
	resp, err = srv.Spreadsheets.Values.Get(spreadsheetId, ReadRangeEntericEmissionFactor).Do()
	if err != nil {
		log.Printf("Unable to retrieve data from sheet: %v", err)
		//return model.DataPackage{}, err
	}

	if resp == nil || len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		defaultValue := []interface{}{"0","0","0","","","","","0.0","0.0","0.0","0.0","0.0","0.0","0.0","0.0","0.0","0.0","0.0","0.0","0.0","0.0"}
		for _, row := range resp.Values {
			if len(row) < 20{
				row = append(row, defaultValue[len(row):]...)
			}
			id,_ 													:= strconv.Atoi(row[0].(string))
			year,_ 													:= strconv.Atoi(row[1].(string))
			month,_													:= strconv.Atoi(row[2].(string))
			location 												:= row[3].(string)
			system 													:= row[4].(string)
			parentClass												:= row[3].(string)
			animalClass 											:= row[5].(string)
			calculatedEF,_											:= strconv.ParseFloat(row[6].(string), 64)
			mapv,_ 													:= strconv.ParseFloat(row[7].(string), 64)
			cofCalculationNEMaintenance,_							:= strconv.ParseFloat(row[8].(string), 64)
			nEMaintenance,_											:= strconv.ParseFloat(row[9].(string), 64)
			nEActivityCattleBuffalo,_ 								:= strconv.ParseFloat(row[10].(string), 64)
			nEGrowthCattleBuffalo,_ 								:= strconv.ParseFloat(row[11].(string), 64)
			nELactationBeefDairyBuffalo,_							:= strconv.ParseFloat(row[12].(string), 64)
			nEWorkCattleBuffalo,_									:= strconv.ParseFloat(row[13].(string), 64)
			nEPregnancyCattleBuffaloSheep,_							:= strconv.ParseFloat(row[14].(string), 64)
			rNEAvailableDietMaintenanceDigestibleEnergyConsumed,_	:= strconv.ParseFloat(row[15].(string), 64)
			rNEAvailableForGrowthDietDisgestibleEnergyConsumed,_	:= strconv.ParseFloat(row[16].(string), 64)
			gECattleBuffaloSheep,_									:= strconv.ParseFloat(row[17].(string), 64)
			entericFermentationEmissionsLivestockCategory,_			:= strconv.ParseFloat(row[18].(string), 64)

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

func LoadEFSheets(request *http.Request, token oauth2.Token) (model.DataPackage, error){
	data := model.DataPackage{}
	client, err := getPermission(request, token)
	srv, err := drive.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
		return model.DataPackage{}, err
	}

	r, err := srv.Files.List().Q("mimeType='application/vnd.google-apps.spreadsheet'").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
		return model.DataPackage{}, err
	}
	if len(r.Files) == 0 {
		fmt.Println("No files found.")
		return model.DataPackage{}, err
	} else {
		for _, i := range r.Files {
			availableSheet := model.SheetDataItem{}
			availableSheet.SheetID = i.Id
			availableSheet.SheetName = i.Name
			availableSheet.Created = i.CreatedTime
			data.SheetData = append(data.SheetData, availableSheet)
		}
	}
	sort.Slice(data.SheetData, func(i, j int) bool {
		const shortDate = "2006-01-02 15:40:10"
		first, _ := time.Parse(shortDate, data.SheetData[i].Created)
		second, _ := time.Parse(shortDate, data.SheetData[j].Created)
		return second.After(first)
	})
	return data, nil
}

func LoadEFSheetIntoDataPackage(data model.DataPackage, request *http.Request, spreadsheetId string, token oauth2.Token) (model.DataPackage, error) {
	data.ClearPackageData()
	client, err := getPermission(request, token)

	srv, err := sheets.New(client)
	if err != nil {
		log.Printf("Unable to retrieve Sheets client: %v", err)
		return model.DataPackage{}, err
	}
	data, err = LoadEFDataFromService(srv, spreadsheetId)
	if err != nil {
		log.Printf("Unable to load sheets: %v", err)
		return model.DataPackage{}, err
	}
	return data, nil
}

func createNewSheet(request *http.Request, token oauth2.Token)(*sheets.Service,error){
	client, err := getPermission(request, token)
	srv, err := sheets.New(client)
	if err != nil {
		log.Printf("Unable to retrieve Sheets client: %v", err)
		return nil, err
	}
	return srv, nil
}

func LoadDataPackageIntoEFSheet(data model.DataPackage, request *http.Request, token oauth2.Token)(error){
	srv, err := createNewSheet(request, token)

	var vrWriteSettings sheets.ValueRange
	var vrWriteSystems sheets.ValueRange
	var vrWriteLocations sheets.ValueRange
	var vrWriteAnimalClasses sheets.ValueRange
	var vrWriteTemperatureLocations sheets.ValueRange
	var vrWriteAnimalNumbers sheets.ValueRange
	var vrWriteEntericFermEFParameters sheets.ValueRange
	var vrWriteEntericEmissionFactors sheets.ValueRange
	var vrWriteErrorRecords sheets.ValueRange
	rb := &sheets.ClearValuesRequest{
	}

	sheet1 := &sheets.Sheet{Properties: &sheets.SheetProperties{
		Title: "Settings",
	}};
	sheet2 := &sheets.Sheet{Properties: &sheets.SheetProperties{
		Title: "System",
	}};
	sheet3 := &sheets.Sheet{Properties: &sheets.SheetProperties{
		Title: "Location",
	}};
	sheet4 := &sheets.Sheet{Properties: &sheets.SheetProperties{
		Title: "AnimalClass",
	}};
	sheet5 := &sheets.Sheet{Properties: &sheets.SheetProperties{
		Title: "TemperatureLocation",
	}};
	sheet6 := &sheets.Sheet{Properties: &sheets.SheetProperties{
		Title: "AnimalNumbers",
	}};
	sheet7 := &sheets.Sheet{Properties: &sheets.SheetProperties{
		Title: "EntericFermEFParameters",
	}};
	sheet8 := &sheets.Sheet{Properties: &sheets.SheetProperties{
		Title: "EntericEmissionFactors",
	}};
	sheet9 := &sheets.Sheet{Properties: &sheets.SheetProperties{
		Title: "ErrorRecords",
	}};
	main := &sheets.Spreadsheet{
		Properties: &sheets.SpreadsheetProperties{
			Title: time.Now().String(),
		},
		Sheets:[]*sheets.Sheet{sheet1,sheet2,sheet3,sheet4,sheet5,sheet6,sheet7,sheet8,sheet9},
	};
	resp, err := srv.Spreadsheets.Create(main).Context(request.Context()).Do()
	if err != nil {
		log.Fatal(err)
	}


	if err != nil {
		log.Fatal(err)
	}

	// -------------------------------------- Settings ------------------------------------------
	writeRangeSettings := WriteRangeSettings
	_, err = srv.Spreadsheets.Values.Clear(resp.SpreadsheetId, writeRangeSettings, rb).Context(request.Context()).Do()
	if err != nil {
		log.Fatal(err)
	}
	header := []interface{}{"Setting", "Value"}
	vrWriteSettings.Values = append(vrWriteSettings.Values, header)
	for _,record := range data.SettingData{
		value := []interface{}{record.Key, record.Value}
		vrWriteSettings.Values = append(vrWriteSettings.Values, value)
	}
	_, err = srv.Spreadsheets.Values.Update(resp.SpreadsheetId, writeRangeSettings, &vrWriteSettings).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		log.Printf("Unable to write into sheet 1: %v", err)
		return err
	}

	// -------------------------------------- System ------------------------------------------
	writeRangeSystems := WriteRangeSystem
	_, err = srv.Spreadsheets.Values.Clear(resp.SpreadsheetId, writeRangeSystems, rb).Context(request.Context()).Do()
	if err != nil {
		log.Fatal(err)
	}
	header = []interface{}{"ID", "System Name"}
	vrWriteSystems.Values = append(vrWriteSystems.Values, header)
	for _,record := range data.SystemData{
		value := []interface{}{record.Id, record.Name}
		vrWriteSystems.Values = append(vrWriteSystems.Values, value)
	}
	_, err = srv.Spreadsheets.Values.Update(resp.SpreadsheetId, writeRangeSystems, &vrWriteSystems).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		log.Printf("Unable to write into sheet 2: %v", err)
		return err
	}

	// -------------------------------------- Location ------------------------------------------
	writeRangeLocations := WriteRangeLocation
	_, err = srv.Spreadsheets.Values.Clear(resp.SpreadsheetId, writeRangeLocations, rb).Context(request.Context()).Do()
	if err != nil {
		log.Fatal(err)
	}
	header = []interface{}{"ID", "Location Name"}
	vrWriteLocations.Values = append(vrWriteLocations.Values, header)
	for _,record := range data.LocationData{
		value := []interface{}{record.Id, record.Name}
		vrWriteLocations.Values = append(vrWriteLocations.Values, value)
	}
	_, err = srv.Spreadsheets.Values.Update(resp.SpreadsheetId, writeRangeLocations, &vrWriteLocations).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		log.Printf("Unable to write into sheet 3: %v", err)
		return err
	}

	// -------------------------------------- Animal Class ------------------------------------------
	writeRangeAnimalClasses := WriteRangeAnimalClass
	_, err = srv.Spreadsheets.Values.Clear(resp.SpreadsheetId, writeRangeAnimalClasses, rb).Context(request.Context()).Do()
	if err != nil {
		log.Fatal(err)
	}
	header = []interface{}{"ID", "Parent Class","Animal Class Name", "Default EF"}
	vrWriteAnimalClasses.Values = append(vrWriteAnimalClasses.Values, header)
	for _,record := range data.AnimalClassData{
		value := []interface{}{record.Id, record.ParentClass, record.Name, record.DefaulEF}
		vrWriteAnimalClasses.Values = append(vrWriteAnimalClasses.Values, value)
	}
	_, err = srv.Spreadsheets.Values.Update(resp.SpreadsheetId, writeRangeAnimalClasses, &vrWriteAnimalClasses).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		log.Printf("Unable to write into sheet 4: %v", err)
		return err
	}

	// -------------------------------------- Temperature Location ------------------------------------------
	writeRangeTemperatureLocations := WriteRangeTemperatureLocation
	_, err = srv.Spreadsheets.Values.Clear(resp.SpreadsheetId, writeRangeTemperatureLocations, rb).Context(request.Context()).Do()
	if err != nil {
		log.Fatal(err)
	}
	header = []interface{}{"ID", "Location", "Year", "Month", "Average Temp"}
	vrWriteTemperatureLocations.Values = append(vrWriteTemperatureLocations.Values, header)
	for _,record := range data.TemperatureLocationData{
		value := []interface{}{record.Id, record.Location, record.Year, record.Month, record.AvgTemp}
		vrWriteTemperatureLocations.Values = append(vrWriteTemperatureLocations.Values, value)
	}
	_, err = srv.Spreadsheets.Values.Update(resp.SpreadsheetId, writeRangeTemperatureLocations, &vrWriteTemperatureLocations).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		log.Printf("Unable to write into sheet 5: %v", err)
		return err
	}

	// -------------------------------------- Animal Numbers ------------------------------------------
	writeRangeAnimalNumbers := WriteRangeAnimalNumber
	_, err = srv.Spreadsheets.Values.Clear(resp.SpreadsheetId, writeRangeAnimalNumbers, rb).Context(request.Context()).Do()
	if err != nil {
		log.Fatal(err)
	}
	header = []interface{}{"ID", "Location", "System", "Animal Class", "Year", "Month", "Animal Number"}
	vrWriteAnimalNumbers.Values = append(vrWriteAnimalNumbers.Values, header)
	for _,record := range data.AnimalNumberData{
		value := []interface{}{record.Id, record.Location, record.System, record.AnimalClass, record.Year, record.Month, record.AnimalNumber}
		vrWriteAnimalNumbers.Values = append(vrWriteAnimalNumbers.Values, value)
	}
	_, err = srv.Spreadsheets.Values.Update(resp.SpreadsheetId, writeRangeAnimalNumbers, &vrWriteAnimalNumbers).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		log.Printf("Unable to write into sheet 6: %v", err)
		return err
	}

	// -------------------------------------- Enteric Fermentation Parameter ------------------------------------------
	writeRangeEntericFermEFParameter := WriteRangeEntericFermEFParameter
	_, err = srv.Spreadsheets.Values.Clear(resp.SpreadsheetId, writeRangeEntericFermEFParameter, rb).Context(request.Context()).Do()
	if err != nil {
		log.Fatal(err)
	}
	header = []interface{}{"ID", "Year", "Month", "Location", "System", "Animal Class", "Body Weight (by month)", "Mature Weight", "Daily Weight Gain", "Fraction of Month Alive", "CF  (MJ day \u207b \u00b9 kg\u207b \u00b9)", "C", "Ca", "Milk Production", "Fat Content (%)", "CPregnancy", "Proportion Animal Class Pregnant", "Proportion of Animal Class lactating", "Fraction of Month Lactating", "Hours Worked", "DE (%)", "Ym"}
	vrWriteEntericFermEFParameters.Values = append(vrWriteEntericFermEFParameters.Values, header)
	for _,record := range data.EntericFermEFParameterData{
		value := []interface{}{record.Id, record.Year, record.Month, record.Location, record.System, record.AnimalClass, record.BodyWeight, record.MatureWeight, record.DailyWeightGain, record.FractionOfMonthAlive, record.CF, record.C, record.CA, record.MilkProd, record.FatContent, record.CPregnancy, record.ProportionAnimalClassPregnant, record.ProportionAnimalClassLactating, record.FractionOfMonthLactating, record.HoursWorked, record.DE, record.YM}
		vrWriteEntericFermEFParameters.Values = append(vrWriteEntericFermEFParameters.Values, value)
	}
	_, err = srv.Spreadsheets.Values.Update(resp.SpreadsheetId, writeRangeEntericFermEFParameter, &vrWriteEntericFermEFParameters).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		log.Printf("Unable to write into sheet 7: %v", err)
		return err
	}

	// -------------------------------------- Emission Factor ------------------------------------------
	writeRangeEntericEmissionFactor := WriteRangeEntericEmissionFactor
	_, err = srv.Spreadsheets.Values.Clear(resp.SpreadsheetId, writeRangeEntericEmissionFactor, rb).Context(request.Context()).Do()
	if err != nil {
		log.Fatal(err)
	}
	header = []interface{}{"ID", "Year", "Month", "Location", "System", "Parent Class", "Animal Class", "Calculated EF (CH&X4 head\u207b \u00b9 year\u207b \u00b9)", "Monthly Average Population (head month\u207b \u00b9)",
		"Coefficient for calculating net energy for maintenance", "Net energy for maintenance (MJ day\u207b \u00b9)", "Net energy for activity for cattle and buffalo (MJ day\u207b \u00b9)", "Net energy for growth for cattle and buffalo (MJ day\u207b \u00b9)",
		"Net energy for lactation for beef, dairy and buffalo (MJ day\u207b \u00b9)", "Net energy for work for cattle and buffalo (MJ day&\u207b \u00b9)", "Net energy for pregnancy for cattle, buffalo and sheep (MJ day\u207b \u00b9)",
		"Ratio of net energy available in a diet for maintenance to digestible energy consumed", "Ratio of net energy available for growth in a diet to digestible energy consumed",
		"Gross energy for cattle, buffalo, sheep (MJ day\u207b \u00b9)", "Enteric fermentation emissions from a livestock category (Gg CH\u2084 month\u207b \u00b9)"}
	vrWriteEntericEmissionFactors.Values = append(vrWriteEntericEmissionFactors.Values, header)
	for _,record := range data.EntericEmissionFactorData{
		value := []interface{}{record.Id, record.Year, record.Month, record.Location, record.System, record.ParentClass, record.AnimalClass, record.CalculatedEF, record.MAP,
			record.CofCalculationNEMaintenance, record.NEMaintenance, record.NEActivityCattleBuffalo, record.NEGrowthCattleBuffalo,
			record.NELactationBeefDairyBuffalo, record.NEWorkCattleBuffalo, record.NEPregnancyCattleBuffaloSheep,
			record.RNEAvailableDietMaintenanceDigestibleEnergyConsumed, record.RNEAvailableForGrowthDietDisgestibleEnergyConsumed,
			record.GECattleBuffaloSheep, record.EntericFermentationEmissionsLivestockCategory}
		vrWriteEntericEmissionFactors.Values = append(vrWriteEntericEmissionFactors.Values, value)
	}
	_, err = srv.Spreadsheets.Values.Update(resp.SpreadsheetId, writeRangeEntericEmissionFactor, &vrWriteEntericEmissionFactors).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		log.Printf("Unable to write into sheet 8: %v", err)
		return err
	}

	// -------------------------------------- Error Records ------------------------------------------
	writeRangeErrorRecord := WriteRangeErrorRecord
	_, err = srv.Spreadsheets.Values.Clear(resp.SpreadsheetId, writeRangeErrorRecord, rb).Context(request.Context()).Do()
	if err != nil {
		log.Fatal(err)
	}
	header = []interface{}{"Record ID", "Year", "Month", "Location", "System", "Animal Class", "Error Message"}
	vrWriteErrorRecords.Values = append(vrWriteErrorRecords.Values, header)
	for _,record := range data.ErrorRecordData{
		value := []interface{}{record.Recordid, record.Year, record.Month, record.Location, record.System, record.AnimalClass, record.ErrorMsg}
		vrWriteErrorRecords.Values = append(vrWriteErrorRecords.Values, value)
	}
	_, err = srv.Spreadsheets.Values.Update(resp.SpreadsheetId, writeRangeErrorRecord, &vrWriteErrorRecords).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		log.Printf("Unable to write into sheet 9: %v", err)
		return err
	}
	return nil
}