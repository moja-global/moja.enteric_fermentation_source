package handlers

import (
	"bytes"
	"fmt"
	"github.com/MullionGroup/flintpro-common/v2/tracing"
	"github.com/MullionGroup/go-website-flintpro-example/models"
	"github.com/hashicorp/go-uuid"
	"github.com/opentracing/opentracing-go"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func reset() {
	tracing.SetTracer(opentracing.NoopTracer{})
	viper.Set("debug", true)
	uid,_ := uuid.GenerateUUID()
	viper.SetDefault("cookies_value_test", uid)
	credential := `{"web":{"client_id":"test.apps.googleusercontent.com","project_id":"aciar-sleek","auth_uri":"https://accounts.google.com/o/oauth2/auth","token_uri":"https://oauth2.googleapis.com/token","auth_provider_x509_cert_url":"https://www.googleapis.com/oauth2/v1/certs","client_secret":"V5jHmP9c9Nrmvog3bswLtxIh","redirect_uris":["http://localhost:8080/index","http://localhost:8080/login_google","http://127.0.0.1/login_google"],"javascript_origins":["http://localhost:8080","http://127.0.0.1:8080"]}}`
	viper.SetDefault("credentials_file",[]byte(credential))
	viper.SetDefault("read_me","../README.md")
}

func TestLogin(t *testing.T) {
	reset()
	Convey("Given a HTTP request for login page with cookies",t, func(){
		req:= httptest.NewRequest("GET", "/", nil)
		resp:= httptest.NewRecorder()
		Convey("When the request is handled by the router", func(){
			//Getting response from the HTTP request, escape JWT auth
			NewRouter().ServeHTTP(resp,req)
			Convey("The response should be a 200",func(){
				So(resp.Code, ShouldEqual, 200)
			})
		})
	})
}

func TestIndex(t *testing.T) {
	reset()
	Convey("Given a HTTP request for index page with cookies",t, func(){
		req:= httptest.NewRequest("GET", "/index", nil)
		resp:= httptest.NewRecorder()
		req.AddCookie(&http.Cookie{
			Name:       "cookie_google_test",
			Value:      viper.GetString("cookies_value_test"),
			Path:       "",
			Domain:     "",
			Expires:    time.Time{},
			RawExpires: "",
			MaxAge:     0,
			Secure:     false,
			HttpOnly:   false,
			SameSite:   0,
			Raw:        "",
			Unparsed:   nil,
		})
		currentDataPackage := models.DataPackage{
			SheetData: nil,
		}
		currentSheetData := models.SheetDataItem{
			SheetID:   "123",
			SheetName: "testing",
			Created:   "2019-01-01",
		}

		currentDataPackage.SheetData = append(currentDataPackage.SheetData, currentSheetData)
		Convey("When the request is handled by the router", func(){
			//Getting response from the HTTP request, escape JWT auth
			NewRouter().ServeHTTP(resp,req)
			Convey("The response should be a 200",func(){
				So(resp.Code, ShouldEqual, 200)
				//This should be 0 because we have a ClearSpreadSheetPackage()
				//So(len(models.CurrentDataPackage.SheetData), ShouldEqual,0)
			})
		})
	})

	Convey("Given a HTTP request for index page with no cookies",t, func(){
		req:= httptest.NewRequest("GET", "/index", nil)
		resp:= httptest.NewRecorder()
		Convey("When the request is handled by the router", func(){
			//Getting response from the HTTP request, escape JWT auth
			NewRouter().ServeHTTP(resp,req)
			Convey("The response should be a 200",func(){
				So(resp.Code, ShouldEqual, 200)
			})
		})
	})
}

func TestLogout(t *testing.T) {
	reset()
	Convey("Given a HTTP request for logout page with cookies",t, func(){
		req:= httptest.NewRequest("POST", "/logout_google", nil)
		resp:= httptest.NewRecorder()
		req.AddCookie(&http.Cookie{
			Name:       "cookie_google_test",
			Value:      viper.GetString("cookies_value_test"),
			Path:       "",
			Domain:     "",
			Expires:    time.Time{},
			RawExpires: "",
			MaxAge:     0,
			Secure:     false,
			HttpOnly:   false,
			SameSite:   0,
			Raw:        "",
			Unparsed:   nil,
		})
		Convey("When the request is handled by the router", func(){
			//Getting response from the HTTP request, escape JWT auth
			NewRouter().ServeHTTP(resp,req)
			Convey("The response should be a 302",func(){
				So(resp.Code, ShouldEqual, 302)
			})
		})
	})
}

func TestAboutUs(t *testing.T) {
	reset()
	Convey("Given a HTTP request for about us page with cookies",t, func(){
		req:= httptest.NewRequest("GET", "/about_us", nil)
		resp:= httptest.NewRecorder()
		req.AddCookie(&http.Cookie{
			Name:       "cookie_google_test",
			Value:      viper.GetString("cookies_value_test"),
			Path:       "",
			Domain:     "",
			Expires:    time.Time{},
			RawExpires: "",
			MaxAge:     0,
			Secure:     false,
			HttpOnly:   false,
			SameSite:   0,
			Raw:        "",
			Unparsed:   nil,
		})
		Convey("When the request is handled by the router", func(){
			//Getting response from the HTTP request, escape JWT auth
			NewRouter().ServeHTTP(resp,req)
			Convey("The response should be a 200",func(){
				So(resp.Code, ShouldEqual, 200)
			})
		})
	})

	Convey("Given a HTTP request for about us page with no cookies",t, func(){
		req:= httptest.NewRequest("GET", "/about_us", nil)
		resp:= httptest.NewRecorder()
		Convey("When the request is handled by the router", func(){
			//Getting response from the HTTP request, escape JWT auth
			NewRouter().ServeHTTP(resp,req)
			Convey("The response should be a 200",func(){
				So(resp.Code, ShouldEqual, 200)
			})
		})
	})
}

//TOOD: I have no idea how to test the google sheet loading yet
func TestSimulationData(t *testing.T) {
	reset()
	Convey("Given a HTTP request for load data page with cookies",t, func(){
		file,err := os.Open("../test.json")
		if err!=nil{
			fmt.Println(err)
		}
		fileContents,err := ioutil.ReadAll(file)
		if err!=nil{
			fmt.Println(err)
		}
		fi,err := file.Stat()
		if err!=nil{
			fmt.Println(err)
		}
		file.Close()

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part,err := writer.CreateFormFile("uploadfile", fi.Name())
		if err!=nil{
			fmt.Println(err)
		}
		part.Write(fileContents)
		writer.Close()
		req:= httptest.NewRequest("POST", "/load_data?type=json", body)
		resp:= httptest.NewRecorder()
		req.Header.Add("Content-Type", writer.FormDataContentType())
		Convey("When the request is handled by the router", func(){
			//Getting response from the HTTP request, escape JWT auth
			NewRouter().ServeHTTP(resp,req)
			Convey("The response should be a 200",func(){
				So(resp.Code, ShouldEqual, 200)
				//So(len(models.CurrentDataPackage.SettingData),ShouldBeGreaterThan, 0)
			})
		})
	})

	Convey("Given a HTTP request for load data page with cookies",t, func(){
		file,err := os.Open("../test.sqlite3")
		if err!=nil{
			fmt.Println(err)
		}
		fileContents,err := ioutil.ReadAll(file)
		if err!=nil{
			fmt.Println(err)
		}
		fi,err := file.Stat()
		if err!=nil{
			fmt.Println(err)
		}
		file.Close()

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part,err := writer.CreateFormFile("uploadfile", fi.Name())
		if err!=nil{
			fmt.Println(err)
		}
		part.Write(fileContents)
		writer.Close()
		req:= httptest.NewRequest("POST", "/load_data?type=db", body)
		resp:= httptest.NewRecorder()
		req.AddCookie(&http.Cookie{
			Name:       "cookie_google_test",
			Value:      viper.GetString("cookies_value_test"),
			Path:       "",
			Domain:     "",
			Expires:    time.Time{},
			RawExpires: "",
			MaxAge:     0,
			Secure:     false,
			HttpOnly:   false,
			SameSite:   0,
			Raw:        "",
			Unparsed:   nil,
		})
		req.Header.Add("Content-Type", writer.FormDataContentType())
		Convey("When the request is handled by the router", func(){
			//Getting response from the HTTP request, escape JWT auth
			NewRouter().ServeHTTP(resp,req)
			Convey("The response should be a 200",func(){
				So(resp.Code, ShouldEqual, 200)
				//So(len(models.CurrentDataPackage.SettingData),ShouldBeGreaterThan, 0)
			})
		})
	})

	Convey("Given a HTTP request for load data page with cookies",t, func(){
		file,err := os.Open("../test.xlsx")
		if err!=nil{
			fmt.Println(err)
		}
		fileContents,err := ioutil.ReadAll(file)
		if err!=nil{
			fmt.Println(err)
		}
		fi,err := file.Stat()
		if err!=nil{
			fmt.Println(err)
		}
		file.Close()

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part,err := writer.CreateFormFile("uploadfile", fi.Name())
		if err!=nil{
			fmt.Println(err)
		}
		part.Write(fileContents)
		writer.Close()
		req:= httptest.NewRequest("POST", "/load_data?type=excel", body)
		resp:= httptest.NewRecorder()
		req.Header.Add("Content-Type", writer.FormDataContentType())
		Convey("When the request is handled by the router", func(){
			//Getting response from the HTTP request, escape JWT auth
			NewRouter().ServeHTTP(resp,req)
			Convey("The response should be a 200",func(){
				So(resp.Code, ShouldEqual, 200)
				//So(len(models.CurrentDataPackage.SettingData),ShouldBeGreaterThan, 0)
			})
		})
	})
}

func TestRunSimulation(t *testing.T) {
	reset()
	Convey("Given a HTTP request for create simulation page with no cookies",t, func(){
		file,err := os.Open("../test.json")
		if err!=nil{
			fmt.Println(err)
		}
		fileContents,err := ioutil.ReadAll(file)
		if err!=nil{
			fmt.Println(err)
		}
		body := bytes.NewBuffer(fileContents)
		req:= httptest.NewRequest("POST", "/run_simulation", body)
		resp:= httptest.NewRecorder()
		Convey("When the request is handled by the router", func(){
			//Getting response from the HTTP request, escape JWT auth
			NewRouter().ServeHTTP(resp,req)
			Convey("The response should be a 200",func(){
				So(resp.Code, ShouldEqual, 200)
			})
		})
	})
}

func TestExportSimulation(t *testing.T) {
	reset()
	Convey("Given a HTTP request for export simulation page with cookies",t, func(){
		file,err := os.Open("../test.json")
		if err!=nil{
			fmt.Println(err)
		}
		fileContents,err := ioutil.ReadAll(file)
		if err!=nil{
			fmt.Println(err)
		}
		body := bytes.NewBuffer(fileContents)
		req:= httptest.NewRequest("POST", "/export_simulation?type=json", body)
		resp:= httptest.NewRecorder()
		Convey("When the request is handled by the router", func(){
			//Getting response from the HTTP request, escape JWT auth
			NewRouter().ServeHTTP(resp,req)
			Convey("The response should be a 200",func(){
				So(resp.Code, ShouldEqual, 200)
			})
		})
	})

	Convey("Given a HTTP request for export simulation page with cookies",t, func(){
		file,err := os.Open("../test.json")
		if err!=nil{
			fmt.Println(err)
		}
		fileContents,err := ioutil.ReadAll(file)
		if err!=nil{
			fmt.Println(err)
		}
		body := bytes.NewBuffer(fileContents)
		req:= httptest.NewRequest("POST", "/export_simulation?type=db", body)
		resp:= httptest.NewRecorder()
		Convey("When the request is handled by the router", func(){
			//Getting response from the HTTP request, escape JWT auth
			NewRouter().ServeHTTP(resp,req)
			Convey("The response should be a 200",func(){
				So(resp.Code, ShouldEqual, 200)
			})
		})
	})

	Convey("Given a HTTP request for export simulation page with cookies",t, func(){
		file,err := os.Open("../test.json")
		if err!=nil{
			fmt.Println(err)
		}
		fileContents,err := ioutil.ReadAll(file)
		if err!=nil{
			fmt.Println(err)
		}
		body := bytes.NewBuffer(fileContents)
		req:= httptest.NewRequest("POST", "/export_simulation?type=excel", body)
		resp:= httptest.NewRecorder()
		Convey("When the request is handled by the router", func(){
			//Getting response from the HTTP request, escape JWT auth
			NewRouter().ServeHTTP(resp,req)
			Convey("The response should be a 200",func(){
				So(resp.Code, ShouldEqual, 200)
			})
		})
	})

	Convey("Given a HTTP request for export simulation page with no cookies",t, func(){
		file,err := os.Open("../test.json")
		if err!=nil{
			fmt.Println(err)
		}
		fileContents,err := ioutil.ReadAll(file)
		if err!=nil{
			fmt.Println(err)
		}
		body := bytes.NewBuffer(fileContents)
		req:= httptest.NewRequest("POST", "/export_simulation?type=json", body)
		resp:= httptest.NewRecorder()
		Convey("When the request is handled by the router", func(){
			//Getting response from the HTTP request, escape JWT auth
			NewRouter().ServeHTTP(resp,req)
			Convey("The response should be a 200",func(){
				So(resp.Code, ShouldEqual, 200)
			})
		})
	})
}