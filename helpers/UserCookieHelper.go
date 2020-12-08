package helpers

import (
	"encoding/json"
	"github.com/MullionGroup/go-website-flintpro-example/cookies"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/sheets/v4"
	"net/http"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

var storeHandler = sessions.NewCookieStore(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))


func SetGoogleCookie(token oauth2.Token, response http.ResponseWriter, request *http.Request) {
	session, err := storeHandler.New(request, "cookie-name")
	if err != nil {
		logrus.Errorf("Storing Get session: %v",err)
	}
	session.Values["AccessToken"] = token.AccessToken
	session.Values["TokenType"] = token.TokenType
	session.Values["RefreshToken"] = token.RefreshToken
	err = session.Save(request, response)
	if err != nil {
		logrus.Errorf("Storing Save session: %v",err)
		return
	}
}

func DeleteGoogleCookie(token oauth2.Token, response http.ResponseWriter, request *http.Request){
	session, err := storeHandler.Get(request, "cookie-name")
	if err != nil {
		logrus.Errorf("Deleting Get session: %v",err)
	}
	session.Values["AccessToken"] = ""
	session.Values["TokenType"] = ""
	session.Values["RefreshToken"] = ""
	session.Options.MaxAge = -1

	err = session.Save(request, response)
	if err != nil {
		logrus.Errorf("Deleting Save session: %v",err)
		return
	}
}

func IsGoogleOAuthed(token oauth2.Token) bool{
	return !IsEmpty(token.AccessToken)
}

func GetGoogleCookie(request *http.Request) oauth2.Token{
	session, err := storeHandler.Get(request, "cookie-name")
	if err != nil {
		logrus.Errorf("Getting Get session: %v",err)
		return oauth2.Token{}
	}
	val1 := session.Values["AccessToken"]
	val2 := session.Values["TokenType"]
	val3 := session.Values["RefreshToken"]

	var accessToken = ""
	accessToken, ok := val1.(string)
	if !ok {
		return oauth2.Token{}
	}

	var tokenType = ""
	tokenType, ok = val2.(string)
	if !ok {
		return oauth2.Token{}
	}

	var refreshToken = ""
	refreshToken, ok = val3.(string)
	if !ok {
		return oauth2.Token{}
	}

	var token = oauth2.Token{}
	token.AccessToken = accessToken
	token.TokenType = tokenType
	token.RefreshToken = refreshToken
	return token
}


var (
	OauthConfGl = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "http://localhost:8080/login_google",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", sheets.SpreadsheetsScope, drive.DriveFileScope, drive.DriveMetadataScope},
		Endpoint:     google.Endpoint,
	}
	OauthStateStringGl = ""
)

/*
InitializeOAuthGoogle Function
*/
func InitializeOAuthGoogle() error {
	storeHandler.Options = &sessions.Options{
		MaxAge:   60 * 15,
		HttpOnly: true,
	}
	b := viper.Get("credentials_file").([]byte)
	secretFile := cookies.SecretFile{}
	if err := json.Unmarshal(b, &secretFile); err != nil {
		return err
	}
	var c *cookies.Cred
	c = secretFile.Web
	OauthConfGl.ClientID = c.ClientID
	OauthConfGl.ClientSecret = c.ClientSecret
	OauthStateStringGl = viper.GetString("oauthStateString")
	return nil
}