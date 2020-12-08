package handlers

import (
	helpers "github.com/MullionGroup/go-website-flintpro-example/helpers"
	model "github.com/MullionGroup/go-website-flintpro-example/models"
	"github.com/gobuffalo/packr/v2"
	"github.com/gomarkdown/markdown"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

type View struct {
	Template *template.Template
	//Report 	model.ReportData
	Layout   string
}

var TemplateFiles []string
var templatesBox = packr.New("Templates", "../templates")
var readmeBox = packr.New("ReadMe","./"+viper.GetString("read_me"))

var FuncMap = template.FuncMap{
	"GetDataLayerJSON": helpers.GetDataLayerJSON,
	"IsFeatureCollection": helpers.IsFeatureCollection,
}

var TPLLogin = View{}
var TPLDataLayer = View{}
var TPLSimulationData = View{}
var TPLIndex = View{}
var TPLAboutUs = View{}
var DBClient helpers.IGormClient = &helpers.GormClient{}
// Handlers

func AboutUsPageHandler(response http.ResponseWriter, request *http.Request) {
	b,err := readmeBox.Find("README.md")

	content := markdown.ToHTML(b, nil, nil)
	templateAboutUs, err := templatesBox.FindString("about_us.html")
	if err != nil {
		logrus.Errorf("Error with parsing Templates in AboutUsPageHandler")
		http.Redirect(response, request, "/index", 302)
		return
	}
	tmplAboutUs := template.New("").Funcs(FuncMap)
	tmplAboutUs, err = tmplAboutUs.Parse(templateAboutUs)
	if err != nil {
		logrus.Errorf("Error with Templates in AboutUsPageHandler")
		http.Redirect(response, request, "/index", 302)
		return
	}
	TPLAboutUs.Template = tmplAboutUs
	err = TPLAboutUs.Template.Execute(response, template.HTML(content))
	if err != nil {
		logrus.Errorf("ErrorCode is -11 : Template %v.", err.Error())
		http.Redirect(response, request, "/index", 500)
		return
	}
}

func PostSimulationDataPageHandler(response http.ResponseWriter, request *http.Request) {
	var data model.DataPackage
	var err error

	token := helpers.GetGoogleCookie(request)
	format := request.FormValue("type")
	if format == "google"{
		sheet_id := request.FormValue("sheet_id")
		if helpers.IsGoogleOAuthed(token) {
			tokenString := helpers.GetGoogleCookie(request)
			data, err = helpers.LoadEFSheetIntoDataPackage(data, request, sheet_id, tokenString)
			if err != nil {
				data.ErrorGenericMsg = append(data.ErrorGenericMsg, "Error loading file")
				logrus.Errorf("Error creating temp file, %v", err)
				http.Redirect(response, request, "/index", 302)
				return
			}
			data.SheetID = sheet_id
		} else {
			http.Redirect(response, request, "/load_google", 302)
			return
		}
	} else if format == "excel"{
		request.ParseMultipartForm(32 << 20)
		file, handler, err := request.FormFile("uploadfile")
		if err != nil {
			data.ErrorGenericMsg = append(data.ErrorGenericMsg, "Error getting file")
			logrus.Errorf("Error getting file, %v", err)
			http.Redirect(response, request, "/index", 302)
			return
		}
		extension := filepath.Ext(handler.Filename)
		if extension!=".xlsx" && extension!=".xls"{
			data.ErrorGenericMsg = append(data.ErrorGenericMsg, "Error file type")
			logrus.Errorf("Error file type")
			http.Redirect(response, request, "/index", 302)
			return
		}
		f, err := ioutil.TempFile("", strconv.FormatInt(time.Now().UnixNano(), 10)+extension)
		if err != nil {
			data.ErrorGenericMsg = append(data.ErrorGenericMsg, "Error opening file")
			logrus.Errorf("Error creating temp file, %v", err)
			http.Redirect(response, request, "/index", 302)
			return
		}
		defer func() {
			f.Close()
			file.Close()
		}()
		io.Copy(f, file)
		data, err = helpers.LoadExcelIntoDataPackage(data, f.Name())
		if err != nil {
			data.ErrorGenericMsg = append(data.ErrorGenericMsg, "Error loading file")
			logrus.Errorf("Error creating temp file, %v", err)
			http.Redirect(response, request, "/index", 302)
			return
		}
	} else if format == "db"{
		request.ParseMultipartForm(32 << 20)
		file, handler, err := request.FormFile("uploadfile")
		if err != nil {
			data.ErrorGenericMsg = append(data.ErrorGenericMsg, "Error getting file")
			logrus.Errorf("Error getting file, %v", err)
			http.Redirect(response, request, "/index", 302)
			return
		}
		extension := filepath.Ext(handler.Filename)
		if extension!=".db" && extension!=".sqlite3"{
			data.ErrorGenericMsg = append(data.ErrorGenericMsg, "Error file type")
			logrus.Errorf("Error file type")
			http.Redirect(response, request, "/index", 302)
			return
		}
		f, err := ioutil.TempFile("", strconv.FormatInt(time.Now().UnixNano(), 10)+extension)
		if err != nil {
			data.ErrorGenericMsg = append(data.ErrorGenericMsg, "Error opening file")
			logrus.Errorf("Error creating temp file, %v", err)
			http.Redirect(response, request, "/index", 302)
			return
		}
		defer func() {
			f.Close()
			file.Close()
		}()
		io.Copy(f, file)
		DBClient = &helpers.GormClient{}
		data, err = helpers.LoadSQLiteIntoDataPackage(data, f.Name())
		if err != nil {
			data.ErrorGenericMsg = append(data.ErrorGenericMsg, "Error loading file")
			logrus.Errorf("Error creating temp file, %v", err)
			http.Redirect(response, request, "/index", 302)
			return
		}
	} else{
		request.ParseMultipartForm(32 << 20)
		file, handler, err := request.FormFile("uploadfile")
		if err != nil {
			data.ErrorGenericMsg = append(data.ErrorGenericMsg, "Error getting file")
			logrus.Errorf("Error creating temp file, %v", err)
			http.Redirect(response, request, "/index", 302)
			return
		}
		extension := filepath.Ext(handler.Filename)
		if extension!=".json"{
			data.ErrorGenericMsg = append(data.ErrorGenericMsg, "Error file type")
			logrus.Errorf("Error file type")
			http.Redirect(response, request, "/index", 302)
			return
		}
		f, err := ioutil.TempFile("", strconv.FormatInt(time.Now().UnixNano(), 10)+extension)
		if err != nil {
			data.ErrorGenericMsg = append(data.ErrorGenericMsg, "Error opening file")
			logrus.Errorf("Error opening file, %v", err)
			http.Redirect(response, request, "/index", 302)
			return
		}
		defer func() {
			f.Close()
			file.Close()
		}()
		io.Copy(f, file)
		data, err = helpers.LoadJSONIntoDataPackage(data, f.Name())
		if err != nil {
			data.ErrorGenericMsg = append(data.ErrorGenericMsg, "Error loading file")
			logrus.Errorf("Error creating temp file, %v", err)
			http.Redirect(response, request, "/index", 302)
			return
		}
	}

	templateSheet, err := templatesBox.FindString("simulation_data.html")
	if err != nil {
		logrus.Errorf("Error with parsing Templates in SimulationDataPageHandler")
		http.Redirect(response, request, "/index", 302)
		return
	}
	tmplSheet := template.New("").Funcs(FuncMap)
	tmplSheet, err = tmplSheet.Parse(templateSheet)
	if err != nil {
		logrus.Errorf("Error with Templates in SimulationDataPageHandler")
		http.Redirect(response, request, "/index", 302)
		return
	}

	result := model.SimulationResultPackage{
		SettingData:						   data.SettingData,
		SystemData:						       data.SystemData,
		LocationData:						   data.LocationData,
		AnimalClassData:					   data.AnimalClassData,
		TemperatureLocationData:			   data.TemperatureLocationData,
		AnimalNumberData:					   data.AnimalNumberData,
		EntericFermEFParameterData:            data.EntericFermEFParameterData,
		ErrorRecordData:                       data.ErrorRecordData,
		ErrorGenericMsg:                       data.ErrorGenericMsg,
		GoogleAuthenticated:				   helpers.IsGoogleOAuthed(token),
	}
	TPLSimulationData.Template = tmplSheet
	err = TPLSimulationData.Template.Execute(response, result)
	if err != nil {
		logrus.Errorf("ErrorCode is -12 : Template %v.", err.Error())
		http.Redirect(response, request, "/index", 302)
		return
	}
}

//func LoadSimulationDataPageHandler(response http.ResponseWriter, request *http.Request) {
//	//defer request.Body.Close()
//	//
//	//type parameters struct {
//	//	GoogleJWT *string
//	//}
//	//
//	//decoder := json.NewDecoder(request.Body)
//	//params := parameters{}
//	//redirectTarget := "/"
//	//err := decoder.Decode(&params)
//	//if err != nil {
//	//	http.Redirect(response, request, redirectTarget, 302)
//	//	return
//	//}
//	//
//	//// Validate the JWT is valid
//	//claims, err := helpers.ValidateGoogleJWT(*params.GoogleJWT)
//	//if err != nil {
//	//	http.Redirect(response, request, redirectTarget, 302)
//	//	return
//	//}
//	//// create a JWT for OUR app and give it back to the client for future requests
//	//helpers.SetGoogleCookie(claims, response)
//
//	templateSheet, err := templatesBox.FindString("simulation_data.html")
//	if err != nil {
//		logrus.Errorf("Error with parsing Templates in SimulationDataPageHandler")
//		http.Redirect(response, request, "/index", 302)
//		return
//	}
//	tmplSheet := template.New("").Funcs(FuncMap)
//	tmplSheet, err = tmplSheet.Parse(templateSheet)
//	if err != nil {
//		logrus.Errorf("Error with Templates in SimulationDataPageHandler")
//		http.Redirect(response, request, "/index", 302)
//		return
//	}
//	http.Redirect(response, request, "/index", 200)
//	return
//}

//-----------------------------------------------------------------------------
// for GET
func LoginGooglePageHandler(response http.ResponseWriter, request *http.Request) {
	token := helpers.GetGoogleCookie(request)
	// Authenticate whether user has gain access from Google
	if (helpers.IsGoogleOAuthed(token)) {
		http.Redirect(response, request, "/index", 302)
	} else {
		// Setup the login_google template and load the variables into template file
		templateLogin, err := templatesBox.FindString("login_google.html")
		if err != nil {
			logrus.Errorf("Error with parsing Templates in LoginGooglePageHandler")
			http.Redirect(response, request, "/index", 302)
			return
		}
		tmplLogin := template.New("").Funcs(FuncMap)
		tmplLogin, err = tmplLogin.Parse(templateLogin)
		if err != nil {
			logrus.Errorf("Error with Templates in LoginGooglePageHandler")
			http.Redirect(response, request, "/index", 302)
			return
		}

		TPLLogin.Template = tmplLogin

		// Load the login_google template
		err = TPLLogin.Template.Execute(response, nil)
		if err != nil {
			logrus.Errorf("ErrorCode is -14 : Template %v.", err.Error())
			http.Redirect(response, request, "/index", 500)
			return
		}
	}
}

func IndexPageHandler(response http.ResponseWriter, request *http.Request) {
	templateIndex, err := templatesBox.FindString("index.html")
	if err != nil {
		logrus.Errorf("Error with parsing Templates in IndexPageHandler")
		http.Redirect(response, request, "/index", 302)
		return
	}
	tmplIndex := template.New("").Funcs(FuncMap)
	tmplIndex, err = tmplIndex.Parse(templateIndex)
	if err != nil {
		logrus.Errorf("Error with Templates in IndexPageHandler")
		http.Redirect(response, request, "/index", 302)
		return
	}
	TPLIndex.Template = tmplIndex
	data := model.DataPackage{}
	token := helpers.GetGoogleCookie(request)
	if helpers.IsGoogleOAuthed(token) {
		data, err = helpers.LoadEFSheets(request,token)
		if err != nil {
			logrus.Errorf("Error with Templates in IndexPageHandler")
			http.Redirect(response, request, "/index", 302)
			return
		}
	} else{
		data.ClearSpreadSheetPackage()
	}
	err = TPLIndex.Template.Execute(response, data)
	if err != nil {
		logrus.Errorf("ErrorCode is -15 : Template %v.", err.Error())
		http.Redirect(response, request, "/index", 500)
		return
	}

}

// for POST
func LogoutGooglePageHandler(response http.ResponseWriter, request *http.Request) {
	logrus.Infof("Logout successful!")
	tokenString := helpers.GetGoogleCookie(request)
	helpers.DeleteGoogleCookie(tokenString, response, request)
	http.Redirect(response, request, "/", 302)
}