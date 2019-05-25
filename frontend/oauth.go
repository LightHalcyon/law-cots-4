package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	// "strconv"
	// "time"

	"os"
	"strings"
)

// ErrorResponse error response struct
type ErrorResponse struct {
	Status      string `json:"status"`
	Description string `json:"error_description"`
}

// OAuthData response struct
type OAuthData struct {
	AccessToken string `json:"access_token,omitempty"`
	ClientID    string `json:"client_id,omitempty"`
	UserID      string `json:"user_id,omitempty"`
}

// HeaderWriter helper method for response returns
func HeaderWriter(w http.ResponseWriter, err ErrorResponse) {
	if err.Description == "500 Internal Server Error" {
		w.WriteHeader(http.StatusInternalServerError)
	} else if err.Description == "401 Unauthorized" {
		w.WriteHeader(http.StatusUnauthorized)
	} else if err.Description == "200 OK" {
		w.WriteHeader(http.StatusOK)
	} else if err.Description == "403 Forbidden" {
		w.WriteHeader(http.StatusForbidden)
	}
}

// Login authentication to https://oauth.infralabs.cs.ui.ac.id/
func Login(w http.ResponseWriter, r *http.Request) {
	type jsonresp struct {
		Status      string `json:"status"`
		AccessToken string `json:"token"`
	}

	type oauthResponse struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    int    `json:"expires_in"`
		TokenType    string `json:"token_type"`
		Scope        string `json:"scope"`
		RefreshToken string `json:"refresh_token"`
	}

	type receivedData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var params receivedData
	receiveddata := json.NewDecoder(r.Body)
	receiveddata.Decode(&params)
	log.Println(params)
	oauthURL := "https://oauth.infralabs.cs.ui.ac.id"
	tokenPath := "/oauth/token"
	// verificationPath := "/oauth/resource"

	data := url.Values{}
	data.Set("username", params.Username)
	data.Set("password", params.Password)
	data.Set("grant_type", "password")
	data.Set("client_id", os.Getenv("CLIENTID"))
	data.Set("client_secret", os.Getenv("CLIENTSECRET"))
	log.Println(data)

	u, _ := url.ParseRequestURI(oauthURL)
	u.Path = tokenPath
	urlStr := u.String()

	client := &http.Client{}
	req, err := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Request Error")
		errorresp := ErrorResponse{
			Status:      "error",
			Description: "500 Internal Server Error",
		}
		log.Println(err)
		json.NewEncoder(w).Encode(errorresp)
	} else {
		defer resp.Body.Close()

		var oautherror ErrorResponse
		var oauthresp oauthResponse
		// var jsonstr jsonresp
		// var errorresp ErrorResponse
		dec := json.NewDecoder(resp.Body)
		if resp.StatusCode == http.StatusOK {
			decodeError := dec.Decode(&oauthresp)
			if decodeError != nil {
				log.Println("Decode Data Error")
				errorresp := ErrorResponse{
					Status:      "error",
					Description: "500 Internal Server Error",
				}
				log.Println(decodeError)
				HeaderWriter(w, errorresp)
				json.NewEncoder(w).Encode(errorresp)
			} else {
				jsonstr := jsonresp{
					Status:      resp.Status,
					AccessToken: oauthresp.AccessToken,
				}

				json.NewEncoder(w).Encode(jsonstr)
			}
		} else {
			errorresp := ErrorResponse{
				Status:      "error",
				Description: resp.Status,
			}
			HeaderWriter(w, errorresp)
			json.NewEncoder(w).Encode(errorresp)
			dec.Decode(&oautherror)
			log.Println(oautherror.Description)
		}
	}
}

// Authenticate validate access token given
func Authenticate(token string) (ErrorResponse, OAuthData) {
	var oauthresp OAuthData
	var oautherror ErrorResponse
	var errorresp ErrorResponse
	oauthURL := "https://oauth.infralabs.cs.ui.ac.id/oauth/token"
	verificationPath := "/oauth/resource"

	u, _ := url.ParseRequestURI(oauthURL)
	u.Path = verificationPath
	urlStr := u.String()

	client := &http.Client{}
	req, err := http.NewRequest("GET", urlStr, nil)
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Request Error")
		errorresp = ErrorResponse{
			Status:      "error",
			Description: "500 Internal Server Error",
		}
		log.Println(err)
	} else {
		defer resp.Body.Close()
		log.Println(resp.Body)
		dec := json.NewDecoder(resp.Body)
		log.Println(resp.Status)
		if resp.StatusCode == http.StatusOK {
			decodeError := dec.Decode(&oauthresp)
			log.Println(oauthresp)
			if decodeError != nil {
				log.Println("Decode Data Error")
				errorresp = ErrorResponse{
					Status:      "error",
					Description: "500 Internal Server Error",
				}
			} else {
				if oauthresp.AccessToken == token {
					oauthresp = OAuthData{
						AccessToken: oauthresp.AccessToken,
						UserID:      oauthresp.UserID,
					}
					errorresp = ErrorResponse{
						Status:      "OK",
						Description: "200 OK",
					}
				} else {
					errorresp = ErrorResponse{
						Status:      "error",
						Description: "401 Unauthorized",
					}
				}
			}
		} else {
			errorresp = ErrorResponse{
				Status:      "error",
				Description: resp.Status,
			}
			dec.Decode(&oautherror)
			log.Println(oautherror.Description)
		}
	}
	return errorresp, oauthresp
}
