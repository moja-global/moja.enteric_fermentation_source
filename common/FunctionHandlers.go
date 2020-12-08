package handlers

import (
	"encoding/json"
	"github.com/MullionGroup/go-website-flintpro-example/helpers"
	model "github.com/MullionGroup/go-website-flintpro-example/models"
	"github.com/gobuffalo/packr/v2"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

var assetsBox = packr.New("Assets", "../assets")

//TODO: we don't have a button for them on this yet
//func EFSheetDownloadHandler(response http.ResponseWriter, request *http.Request) {
//	cookieData := cookies.GetGoogleCookie(request)
//	data := model.DataPackage{}
//	body, err := ioutil.ReadAll(request.Body)
//	err = json.Unmarshal(body, &data)
//	if err != nil {
//		data.ErrorGenericMsg = append(data.ErrorGenericMsg, "Error reading raw data: "+err.Error())
//		logrus.Errorf("ErrorCode is 3-1 : Template %v.", err.Error())
//		http.Redirect(response, request, "/index", 302)
//		return
//	}
//
//	// Check if it is debug mode
//	if viper.GetBool("debug") && len(request.Cookies()) > 0 && request.Cookies()[0].Name == "cookie_google_test" && request.Cookies()[0].Value == viper.GetString("cookies_value_test") {
//		cookieData = cookies.GoogleCookieContents{
//			AuthURL:    "test",
//			AuthCode:   "test",
//			OAuthToken: "test",
//			Msg:        "test",
//		}
//	}
//
//	if helpers.IsGoogleOAuthed(cookieData) {
//		vars := mux.Vars(request)
//		sheet_id, _ := vars["sheet_id"]
//		data, err := json.Marshal(data)
//		if err != nil {
//			logrus.Errorf("Error Marshalling data (%v)", err.Error())
//			http.Redirect(response, request, "/index", 302)
//			return
//		}
//
//		filename := fmt.Sprintf("ef_inputfile.%v.json", sheet_id)
//		response.Header().Set("Content-Disposition", "attachment; filename="+filename)
//		response.Header().Set("Content-Type", "application/json")
//		response.Header().Set("Content-Length", strconv.Itoa(len(data)))
//		response.Header().Set("Connection", "close")
//		response.WriteHeader(http.StatusOK)
//		response.Write(data)
//	} else {
//		http.Redirect(response, request, "/index", 302)
//	}
//}

// This is to redirect user to Google login page
func RedirectHandlerGoogleAuth(response http.ResponseWriter, request *http.Request) {
	token := helpers.GetGoogleCookie(request)
	if helpers.IsGoogleOAuthed(token){
		http.Redirect(response, request, "/index", http.StatusTemporaryRedirect)
	} else{
		URL, err := url.Parse(helpers.OauthConfGl.Endpoint.AuthURL)
		if err != nil {
			http.Redirect(response, request, "/", http.StatusTemporaryRedirect)
		}
		parameters := url.Values{}
		parameters.Add("client_id", helpers.OauthConfGl.ClientID)
		parameters.Add("scope", strings.Join(helpers.OauthConfGl.Scopes, " "))
		parameters.Add("redirect_uri", helpers.OauthConfGl.RedirectURL)
		parameters.Add("response_type", "code")
		parameters.Add("state", helpers.OauthStateStringGl)
		URL.RawQuery = parameters.Encode()
		url := URL.String()
		http.Redirect(response, request, url, http.StatusTemporaryRedirect)
	}
}

// Callback for Google to redirect to our endpoint
func LoginHandlerGoogleAuth(response http.ResponseWriter, request *http.Request) {
	state := request.FormValue("state")
	if state != helpers.OauthStateStringGl {
		logrus.Error("invalid oauth state, expected " + helpers.OauthStateStringGl + ", got " + state + "\n")
		http.Redirect(response, request, "/", http.StatusTemporaryRedirect)
		return
	}

	code := request.FormValue("code")
	if code == "" {
		logrus.Error("Code not found..")
		reason := request.FormValue("error_reason")
		if reason == "user_denied" {
			logrus.Error("User has denied Permission..")
			http.Redirect(response, request, "/", http.StatusTemporaryRedirect)
		}
	} else {
		token, err := helpers.OauthConfGl.Exchange(oauth2.NoContext, code)
		if err != nil {
			logrus.Error("oauthConfGl.Exchange() failed with " + err.Error() + "\n")
			http.Redirect(response, request, "/", http.StatusTemporaryRedirect)
			return
		}
		helpers.SetGoogleCookie(*token, response, request)
		resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(token.AccessToken))
		if err != nil {
			logrus.Error("Get: " + err.Error() + "\n")
			http.Redirect(response, request, "/", http.StatusTemporaryRedirect)
			return
		}
		defer resp.Body.Close()
		redirectTarget := "/index"
		http.Redirect(response, request, redirectTarget, 302)
		return
	}
}

func ExportSimulationHandler(response http.ResponseWriter, request *http.Request){
	data := model.DataPackage{}
	body, err := ioutil.ReadAll(request.Body)
	err = json.Unmarshal(body, &data)
	if err != nil {
		data.ErrorGenericMsg = append(data.ErrorGenericMsg, "Error reading raw data: "+err.Error())
		logrus.Errorf("ErrorCode is 3-6 : Template %v.", err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte("Marshalling error: "+err.Error()))
		return
	}

	format := request.FormValue("type")
	data.FileName = time.Now().String()
	var byteData []byte
	var filename string
	if format == "google"{
		token := helpers.GetGoogleCookie(request)
		if helpers.IsGoogleOAuthed(token) {
			tokenString := helpers.GetGoogleCookie(request)
			err = helpers.LoadDataPackageIntoEFSheet(data, request, tokenString)
			if err != nil {
				data.ErrorGenericMsg = append(data.ErrorGenericMsg, "Error downloading: "+err.Error())
				logrus.Errorf("ErrorCode is 3-3 : Template %v.", err.Error())
				response.WriteHeader(http.StatusInternalServerError)
				response.Write([]byte("Can't export into Google drive file: "+err.Error()))
				return
			}
			filename = data.FileName
			filenameJson, err := json.Marshal(filename)
			if err != nil {
				data.ErrorGenericMsg = append(data.ErrorGenericMsg, "Error downloading: "+err.Error())
				logrus.Errorf("ErrorCode is 3-3 : Template %v.", err.Error())
				response.WriteHeader(http.StatusInternalServerError)
				response.Write([]byte("Can't export into Google drive file: "+err.Error()))
				return
			}
			response.WriteHeader(http.StatusOK)
			response.Write(filenameJson)
			return
		} else{
			logrus.Errorf("ErrorCode is 3-3-1 : Not Authenticated")
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte("You are not signed in via Google login to use this service."))
			return
		}
	} else if format == "excel"{
		filename = data.FileName+".xlsx"
		byteData, err = helpers.LoadDataPackageIntoExcel(data, filename)
	} else if format == "db"{
		filename = data.FileName+".sqlite3"
		byteData, err = helpers.LoadDataPackageIntoSQLite(data)
	} else{
		filename = data.FileName+".json"
		byteData, err = helpers.LoadDataPackageIntoJSON(data)
	}

	if err != nil {
		data.ErrorGenericMsg = append(data.ErrorGenericMsg, "Error downloading: "+err.Error())
		logrus.Errorf("ErrorCode is 3-3 : Template %v.", err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte("Unable to export file: "+err.Error()))
		return
	}

	if strings.HasSuffix(filename, ".json") {
		response.Header().Set("Content-Type", "application/json; charset=UTF-8")
	} else if strings.HasSuffix(filename, ".xlsx") {
		response.Header().Set("Content-Type", "application/octet-stream")
	} else if strings.HasSuffix(filename, ".sqlite3") {
		response.Header().Set("Content-Type", "application/vnd.sqlite3")
	} else {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte("Unknown filename error"))
		return
	}
	response.Header().Set("Content-Disposition", "attachment; filename="+filename)
	response.Header().Set("Connection", "close")
	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	response.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
	if request.Method == "OPTIONS" {
		return
	}
	response.WriteHeader(http.StatusOK)
	response.Write(byteData)
}

func AssetsHandler(response http.ResponseWriter, request *http.Request){
	file :=  strings.TrimPrefix(request.RequestURI, "/assets/")
	data, err := assetsBox.Find(file)
	if err != nil {
		logrus.Errorf("Error with parsing item %v", err)
		http.Redirect(response, request, "/index", 302)
		return
	}

	contentType := http.DetectContentType(data)
	extension := filepath.Ext(file)
	if strings.Contains(contentType,"text/plain") && (extension == ".css" || extension == ".scss"){
		contentType = "text/css"
	} else if strings.Contains(contentType,"text/plain") && extension == ".js"{
		contentType = "text/javascript"
	}

	response.Header().Set("Connection", "close")
	response.Header().Add("Content-Type", contentType)
	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	response.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
	if request.Method == "OPTIONS" {
		return
	}
	response.WriteHeader(http.StatusOK)
	response.Write(data)
}

func RunSimulationHandler(response http.ResponseWriter, request *http.Request) {
	data := model.DataPackage{}
	body, err := ioutil.ReadAll(request.Body)
	err = json.Unmarshal(body, &data)
	if err != nil {
		data.ErrorGenericMsg = append(data.ErrorGenericMsg, "Error reading raw data: "+err.Error())
		logrus.Errorf("ErrorCode is 3-6 : Template %v.", err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte("Marshaling error"))
		return
	}
	data, err = helpers.RunSimulation(data)
	if err != nil {
		logrus.Errorf("Error simulation: %v", err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte("Unknown simulation run error: "+err.Error()))
		return
	}
	result := model.SimulationResultPackage{
		TemperatureLocationData:			   data.TemperatureLocationData,
		AnimalNumberData:					   data.AnimalNumberData,
		EntericFermEFParameterData:            data.EntericFermEFParameterData,
		EntericEmissionFactorData:             data.EntericEmissionFactorData,
		EntericEmissionFactorDataUserFriendly: data.EntericEmissionFactorDataUserFriendly,
		ErrorRecordData:                       data.ErrorRecordData,
		ErrorGenericMsg:                       data.ErrorGenericMsg,
	}
	dataJson, err := json.Marshal(&result)
	if err != nil {
		logrus.Errorf("Error Marshalling data (%v)", err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte("Marshalling error"))
		return
	}
	response.Header().Set("Connection", "close")
	response.Header().Set("Content-Type", "application/json")
	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	response.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
	response.WriteHeader(http.StatusOK)
	response.Write(dataJson)
}