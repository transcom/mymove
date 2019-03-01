package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"

	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/models"
)

func getCookie(name string, cookies []*http.Cookie) (*http.Cookie, error) {
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie, nil
		}
	}
	return nil, errors.Errorf("Unable to find cookie: %s", name)
}

func main() {
	jar, _ := cookiejar.New(nil)
	client := http.Client{
		Jar: jar,
	}

	req, err := http.NewRequest("GET", "http://milmovelocal:8080/", nil)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	csrf, err := getCookie("masked_gorilla_csrf", resp.Cookies())
	if err != nil {
		log.Fatal(err)
	}
	csrfToken := csrf.Value

	req, err = http.NewRequest("POST", "http://milmovelocal:8080/devlocal-auth/create", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("x-csrf-token", csrfToken)

	resp, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	user := models.User{}
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(user)

	sessionCookie, err := getCookie("mil_session_token", resp.Cookies())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(sessionCookie)

	req, err = http.NewRequest("GET", "http://milmovelocal:8080/internal/users/logged_in", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("x-csrf-token", csrfToken)
	req.AddCookie(sessionCookie)

	resp, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))
}
