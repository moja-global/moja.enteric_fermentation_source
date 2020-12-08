package helpers

import (
	"encoding/json"
	"fmt"
	model "github.com/MullionGroup/go-website-flintpro-example/models"
	"io/ioutil"
	"sort"
)

func LoadJSONIntoDataPackage(data model.DataPackage, filename string) (model.DataPackage, error) {
	data.ClearPackageData()
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Reading inputfile (%v) failed with error (%v)", filename, err.Error())
		return model.DataPackage{}, err
	}

	err = json.Unmarshal(b, &data)
	if err != nil {
		fmt.Printf("Unmarshal of inputfile (%v) failed with error (%v)", filename, err.Error())
		return model.DataPackage{}, err
	}

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

	for index,_ := range data.EntericEmissionFactorData {
		data.EntericEmissionFactorData[index].Location = locationsMap[data.EntericEmissionFactorData[index].Locationid].Name
		data.EntericEmissionFactorData[index].System = systemMap[data.EntericEmissionFactorData[index].Systemid].Name
		data.EntericEmissionFactorData[index].AnimalClass = animalsMap[data.EntericEmissionFactorData[index].AnimalClassid].Name
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

	sort.SliceStable(data.EntericEmissionFactorDataUserFriendly, func(i, j int) bool {
		if locationsMap[data.EntericEmissionFactorDataUserFriendly[i].Location].Id < locationsMap[data.EntericEmissionFactorDataUserFriendly[j].Location].Id {
			return true
		}
		if locationsMap[data.EntericEmissionFactorDataUserFriendly[i].Location].Id > locationsMap[data.EntericEmissionFactorDataUserFriendly[j].Location].Id {
			return false
		}
		if systemMap[data.EntericEmissionFactorDataUserFriendly[i].System].Id < systemMap[data.EntericEmissionFactorDataUserFriendly[j].System].Id {
			return true
		}
		if systemMap[data.EntericEmissionFactorDataUserFriendly[i].System].Id > systemMap[data.EntericEmissionFactorDataUserFriendly[j].System].Id {
			return false
		}
		if animalsMap[data.EntericEmissionFactorDataUserFriendly[i].AnimalClass].Id < animalsMap[data.EntericEmissionFactorDataUserFriendly[j].AnimalClass].Id {
			return true
		}
		if animalsMap[data.EntericEmissionFactorDataUserFriendly[i].AnimalClass].Id > animalsMap[data.EntericEmissionFactorDataUserFriendly[j].AnimalClass].Id {
			return false
		}
		if data.EntericEmissionFactorDataUserFriendly[i].Year < data.EntericEmissionFactorDataUserFriendly[j].Year {
			return true
		}
		if data.EntericEmissionFactorDataUserFriendly[i].Year > data.EntericEmissionFactorDataUserFriendly[j].Year {
			return false
		}
		return data.EntericEmissionFactorDataUserFriendly[i].Month <= data.EntericEmissionFactorDataUserFriendly[j].Month
	})
	return data, nil
}

func LoadDataPackageIntoJSON(data model.DataPackage) ([]byte, error){
	records, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Unmarshal of inputfile failed with error (%v)", err.Error())
		return []byte{}, err
	}
	return records, nil
}