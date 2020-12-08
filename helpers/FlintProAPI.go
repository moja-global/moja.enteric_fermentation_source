package helpers

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type CommonFields struct {
	ID        			string	`json:"id" firestore:"id" gorm:"primary_key"`
	UserId 				string	`json:"user_id" firestore:"user_id"`
	GroupId   			string	`json:"group_id" firestore:"group_id"`
	OrgId   			string	`json:"org_id" firestore:"org_id"`
	Permissions			uint32	`json:"permissions" firestore:"permissions"`
	Created   			string 	`json:"created" firestore:"created"`
}

type CommonListFields struct {
	Page		int 					`json:"page" firestore:"page"`
	PerPage		int 					`json:"per_page" firestore:"per_page"`
	Total		int 					`json:"total" firestore:"total"`
	TotalPages	int 					`json:"total_pages" firestore:"total_pages"`
	Offset		int 					`json:"offset" firestore:"offset"`
	Limit		int 					`json:"limit" firestore:"limit"`
}


type DataLayerCategories struct {
	ID        	string 		`json:"id" firestore:"id" gorm:"primary_key"`
	CategoryID  string 		`json:"category_id" firestore:"category_id"`
	DataLayerID string   	`json:"-" firestore:"-" gorm:"index"` // Don't serialize + index which is very important for performance.
	Name   		string      `json:"name" firestore:"name"`
}

type DataLayerData struct {
	CommonFields

	Type   				string      			`json:"type" firestore:"type"`			// "Raster", "RasterStack", "Feature Collection"
	Version   			string      			`json:"version" firestore:"version"`		// Extra info on Type used to inform processing

	Name   				string      			`json:"name" firestore:"name"`
	Desc   				string      			`json:"description" firestore:"description" gorm:"column:description"`
	DefaultLayer		bool					`json:"default_layer" firestore:"default_layer" gorm:"column:default_layer"`

	Status   			string      			`json:"status" firestore:"status"`
	StatusMessage		string      			`json:"status_message" firestore:"status_message"`

	DataType			string      			`json:"data_type" firestore:"data_type"`
	NoDataValueInt   	int64      				`json:"no_data_int" firestore:"no_data_int" gorm:"column:no_data_int"`
	NoDataValueFloat  	float64      			`json:"no_data_float" firestore:"no_data_float" gorm:"column:no_data_float"`
	NumberLayers   		int      				`json:"number_layers" firestore:"number_layers" gorm:"column:number_layers"`

	FileURI   			string      			`json:"file_uri" firestore:"file_uri"`
	MetaDataURI			string      			`json:"metadata_uri" firestore:"metadata_uri"`

	CategoryTitle		string      			`json:"category_title" firestore:"category_title"`
	CategoryProp		string      			`json:"category_prop" firestore:"category_prop"`
	CategoryCode		string      			`json:"category_code" firestore:"category_code"`

	Classification  	string 					`json:"classification" firestore:"classification"`
	Tags	 			string   				`json:"tags" firestore:"tags"`

	FeatureCollectionID string 					`json:"feature_collection_id" firestore:"feature_collection_id"`
	Categories			[]DataLayerCategories	`json:"categories" firestore:"categories" gorm:"ForeignKey:DataLayerID"`
}

type DataLayerDTO struct {
	DataLayerData
	AccessTokens []string `json:"access_tokens"`
}

type DataLayerDTOList struct {
	CommonListFields
	Data []DataLayerDTO		`json:"data"`
}

func GetDataLayerJSON(datalayer DataLayerDTO) string {
	xJson, _ := json.MarshalIndent(datalayer, "", "    ")
	//xJson, _ := json.Marshal(datalayer)
	return string(xJson)
}

func IsFeatureCollection(datalayer DataLayerDTO) bool {
	return datalayer.Type == "Feature Collection"
}

func LoadDataLayers(jwtToken string) DataLayerDTOList {

	urlStr := "https://staging.flintproapi.com/v1/datalayers?verbose"
	req, err := http.NewRequest("GET", urlStr, nil)

	if err != nil {
		logrus.Errorf("ErrorCode is -18 : Login failed %v.", err.Error())
		return DataLayerDTOList{}
	}

	var transport http.RoundTripper = &http.Transport{
		DisableKeepAlives: true,
	}
	//timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Transport: transport,
		//Timeout:   timeout,
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", jwtToken)

	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("ErrorCode is -19 : Login failed %v.", err.Error())
		return DataLayerDTOList{}
	}

	if resp.StatusCode != 200 {
		logrus.Errorf("Login failed! %v", resp.Status)
		//fmt.Fprintf(w, "ErrorCode is -13 : Login failed %v.", resp.Status)
		return DataLayerDTOList{}
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("Login response error! %v", err.Error())
		return DataLayerDTOList{}
	}

	var list = DataLayerDTOList{}
	err = json.Unmarshal(responseBody, &list)
	if err != nil {
		logrus.Errorf("Unmarshal error! %v", err.Error())
		return DataLayerDTOList{}
	}

	return list
}

