package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"text/template"
)

type Conf struct {
	Token string `json:"token"`
	URL   string `json:"url"`
	Port  string `json:"port"`
}

type Page struct {
	Title string
	Count int
}

func initConf() error {
	fileName := "../conf.json"
	conf := Conf{}
	bin, err := json.Marshal(conf)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fileName, bin, 0644)
	return err
}

func loadConf() (Conf, error) {
	fileName := "../conf.json"
	conf := Conf{}
	bin, err := ioutil.ReadFile(fileName)
	if err != nil {
		return conf, err
	}
	err = json.Unmarshal(bin, &conf)
	if err != nil {
		return conf, err
	}
	return conf, nil
}

func viewTest(w http.ResponseWriter, r *http.Request) {
	page := Page{"Hello World.", 1}
	tmpl, err := template.ParseFiles("../views/layout.html") // ParseFilesを使う
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(w, page)
	if err != nil {
		panic(err)
	}
}

func viewHome(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("../views/index.html")
	if err != nil {
		panic(err)
	}
	//err = tmpl.Execute(w, page)
	page := Page{}
	err = tmpl.Execute(w, page)
	if err != nil {
		panic(err)
	}
}

func viewInvite(w http.ResponseWriter, r *http.Request) {
	mail := r.PostFormValue("mail")
	conf, err := loadConf()
	if err != nil {
		log.Fatalln(err)
	}

	client := &http.Client{}
	data := url.Values{
		"email":      {mail},
		"token":      {conf.Token},
		"set_active": {"true"},
	}

	resp, err := client.Post(
		"https://"+conf.URL+"/api/users.admin.invite",
		"application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		log.Fatalln(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)

	var status map[string]interface{}
	json.Unmarshal(body, &status)
	fmt.Println(status["ok"])
	statusOk := fmt.Sprint(status["ok"])

	if statusOk == "true" {
		fmt.Fprintf(w, mail+"に招待状を送信しました.")
	} else {
		fmt.Fprintf(w, "失敗した。失敗した。失敗した。"+fmt.Sprint(status["error"]))
	}
}

func main() {
	conf, err := loadConf()
	if err != nil {
		log.Fatalln(err)
	}
	http.HandleFunc("/test", viewTest) // hello
	http.HandleFunc("/invite", viewInvite)
	http.HandleFunc("/", viewHome)
	http.ListenAndServe(":"+conf.Port, nil)
}
