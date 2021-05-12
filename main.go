package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
	// "path/filepath"
	// "net/http"
	// "os"
	// "io/ioutil"
	// "text/template"
	// "net/http"
	// "log"
	// "os"
)

/*--------------------------------------------------------------------------------------------
-------------------------------------- Type Struct -------------------------------------------
----------------------------------------------------------------------------------------------*/

// Datas allows...
type Datas struct {
	ArtistsData   []Artists   `json:"Artists_Datas"`
	LocationsData []Locations `json:"Locations_Datas"`
	RelationData  []Relations `json:"Relation_Datas"`
	DatesData     []Dates     `json:"Dates_Datas"`
}

// Artistists struct allows to put on map the data of artists.json
type Artists struct {
	ID           int      `json:"id"`
	IMAGE        string   `json:"image"`
	NAME         string   `json:"name"`
	MEMBERS      []string `json:"members"`
	CREA_DATE    int      `json:"creationDate"`
	FIRST_ALBUM  string   `json:"firstAlbum"`
	LOCATIONS    string   `json:"locations"`
	CONCERT_DATE string   `json:"concertDates"`
	RELATION     string   `json:"relations"`
}

// Locations struct allows to ...
type Locations struct {
	ID        int      `json:"id"`
	LOCATIONS []string `json:"locations"`
	DATES     string   `json:"dates"`
}

// Dates struct allows to ..
type Dates struct {
	ID    int      `json:"id"`
	DATES []string `json:"dates"`
}

// Realtions struct allows to ..
type Relations struct {
	ID         int                 `json:"id"`
	DATESLOCAT map[string][]string `json:"datesLocations"`
}

/*--------------------------------------------------------------------------------------------
----------------------------------------------------------------------------------------------
----------------------------------------------------------------------------------------------*/
// var ART []Artists

// var locations []Locations
// var dates []Dates
// var relations []Relations

/*--------------------------------------------------------------------------------------------
------------------------------------------ Func ----------------------------------------------
----------------------------------------------------------------------------------------------*/
func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmp, error := template.ParseFiles("./templates/index.html")

	if error != nil {
		log.Fatal(error)
	}

	if error := tmp.Execute(w, tmp); error != nil {
		log.Fatal(error)
	}

}

func mainPageHandler(w http.ResponseWriter, r *http.Request) {
	file, error := template.ParseFiles("./templates/mainPage.html")

	if error != nil {
		log.Fatal(error)
	}

	jsonfile, _ := os.Open("./data/artists.json")

	defer jsonfile.Close()

	byteValue, _ := ioutil.ReadAll(jsonfile)

	var data Datas

	json.Unmarshal(byteValue, &data)

	var test string
	var id int
	var test2 string

	var toPrint string

	for i := 0; i < len(data.ArtistsData); i++ {

		id = data.ArtistsData[i].ID
		test = data.ArtistsData[i].NAME
		test2 = data.ArtistsData[i].IMAGE
		// test2 = "Artists ID : " + strconv.Itoa(data.ArtistsData[i].ID)

		toPrint += "<div class='artists' style='text-align:center; background-image: url(" + test2 + ");background-position: center center;background-repeat: no-repeat;background-size: cover;border-radius: 30px;box-shadow: 0 4px 10px 0 rgb(14 14 14);'>"
		toPrint += "<a href='/artists/" + strconv.Itoa(id) + "'/>"
		// toPrint += strconv.Itoa(id)
		toPrint += "<h3>" + strconv.Itoa(id) + " - " + test + "</h3>"
		// toPrint += "<img src=\"" + test2 + "\"/>"

		toPrint += "</div>"

	}
	// toPrint += "</div>"
	file.Execute(w, toPrint)
	// file.ExecuteTemplate(w, toPrint)

	// Writer := data.ArtistsData

	// if error := file.Execute(w, Writer); error != nil {
	// 	log.Fatal(error)
	// }

}

func artistHandler(w http.ResponseWriter, r *http.Request) {
	// var toto []Artists
	tmp, error := template.ParseFiles("./templates/artists.html")

	jsonfile, _ := os.Open("./data/artists.json")

	defer jsonfile.Close()

	byteValue, _ := ioutil.ReadAll(jsonfile)

	var patate Datas

	json.Unmarshal(byteValue, &patate)

	id := r.URL.Path[9:]

	p, _ := strconv.Atoi(id)

	if error != nil {
		log.Fatal(error)
	}

	if error := tmp.Execute(w, patate.ArtistsData[p-1]); error != nil {
		log.Fatal(error)
	}

}

/*--------------------------------------------------------------------------------------------
--------------------------------------Main Func-----------------------------------------------
----------------------------------------------------------------------------------------------*/

func main() {
	// Serving templates files
	filesServer := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", filesServer))

	// Index handler
	http.HandleFunc("/", indexHandler)

	http.HandleFunc("/mainPage", mainPageHandler)

	http.HandleFunc("/artists/", artistHandler)

	fmt.Println("Server is starting...")
	fmt.Println()
	fmt.Println("Go on http://localhost:8080/")
	fmt.Println()
	fmt.Println("To shut down the server press CTRL + C")

	// Starting serveur
	http.ListenAndServe(":8080", nil)

}
