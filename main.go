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
)

/*--------------------------------------------------------------------------------------------
-------------------------------------- Type Struct -------------------------------------------
----------------------------------------------------------------------------------------------*/

// Datas type allows store all the datas of the structs
type Datas struct {
	ArtistsData   []Artists   `json:"Artists_Datas"`
	LocationsData []Locations `json:"Locations_Datas"`
	RelationData  []Relations `json:"Relation_Datas"`
	DatesData     []Dates     `json:"Dates_Datas"`
}

// Artists struct match artists.json var type and logic
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

// Locations struct locations.json
type Locations struct {
	ID        int      `json:"id"`
	LOCATIONS []string `json:"locations"`
	DATES     string   `json:"dates"`
}

// Dates struct match dates.json var type and logic
type Dates struct {
	ID    int      `json:"id"`
	DATES []string `json:"dates"`
}

// Realtions struct match relations.json var type and logic
type Relations struct {
	ID         int                 `json:"id"`
	DATESLOCAT map[string][]string `json:"datesLocations"`
}

/*--------------------------------------------------------------------------------------------
------------------------------ Func Handler Index and MainPage -------------------------------
----------------------------------------------------------------------------------------------*/
// index.html
func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Set index.html as template
	file, error := template.ParseFiles("./templates/index.html")

	if error != nil {
		log.Fatal(error)
	}

	// Render template
	if error := file.Execute(w, file); error != nil {
		log.Fatal(error)
	}

}

//mainPage.html
func mainPageHandler(w http.ResponseWriter, r *http.Request) {
	file, error := template.ParseFiles("./templates/mainPage.html")

	if error != nil {
		log.Fatal(error)
	}

	// data var of type DATAS struct
	var data Datas

	// Open JSON files
	data = openJSON("artists.json", data)
	data = openJSON("locations.json", data)

	// Render template with value stored on data variable
	file.Execute(w, data)
}

/*--------------------------------------------------------------------------------------------
-------------------------------Read and store Data from JSON----------------------------------
----------------------------------------------------------------------------------------------*/

// openJSON function open a json file from the data folder, and store into the Datas struct
func openJSON(file string, data Datas) Datas {
	jsonfile, _ := os.Open("./data/" + file)
	defer jsonfile.Close()

	// read the data from the json file
	byteValue, _ := ioutil.ReadAll(jsonfile)

	//encoding  of the data into byte
	json.Unmarshal(byteValue, &data)

	return data
}

/*--------------------------------------------------------------------------------------------
-------------------------------Func Handler ArtistPage----------------------------------------
----------------------------------------------------------------------------------------------*/

func artistHandler(w http.ResponseWriter, r *http.Request) {
	tmp, error := template.ParseFiles("./templates/artists.html")

	// ART var of type DATAS struct
	var ART Datas

	// open JSON files
	ART = openJSON("artists.json", ART)
	ART = openJSON("relation.json", ART)

	// Get the ID of the Artist
	id := r.URL.Path[9:]
	p, _ := strconv.Atoi(id)

	if error != nil {
		log.Fatal(error)
	}

	var infoArt = make(map[string]interface{})

	// Store usefull information in infoArt map
	infoArt["IMAGE"] = ART.ArtistsData[p-1].IMAGE
	infoArt["NAME"] = ART.ArtistsData[p-1].NAME
	infoArt["MEMBERS"] = ART.ArtistsData[p-1].MEMBERS
	infoArt["CREA_DATE"] = ART.ArtistsData[p-1].CREA_DATE
	infoArt["FIRST_ALBUM"] = ART.ArtistsData[p-1].FIRST_ALBUM
	infoArt["DATESLOCAT"] = ART.RelationData[p-1].DATESLOCAT

	if error := tmp.Execute(w, infoArt); error != nil {
		log.Fatal(error)
	}
}

/*--------------------------------------------------------------------------------------------
-------------------------------Functions for the Filter Part----------------------------------
----------------------------------------------------------------------------------------------*/

// getFilters function recovers, from the filters on the MainPage HTML, the various data and stores them into a map (filters)
func getFilters(r *http.Request) map[string]string {
	filters := make(map[string]string)

	filters["members"] = r.URL.Query().Get("members")
	filters["firstAlbum"] = r.URL.Query().Get("firstAlbum")
	filters["creationDate"] = r.URL.Query().Get("creationDate")
	filters["CitySearch"] = r.URL.Query().Get("CitySearch")

	return filters
}

// /filter.html
func filterHandler(w http.ResponseWriter, r *http.Request) {
	file, error := template.ParseFiles("./templates/filter.html")

	if error != nil {
		log.Fatal(error)
	}

	// Get filters
	filters := getFilters(r)

	var data Datas

	// Open JSON files
	data = openJSON("artists.json", data)
	data = openJSON("locations.json", data)

	// check if there are datas into filters on the HTMLpage
	if filters["creationDate"] == "" && filters["firstAlbum"] == "" && filters["members"] == "" && filters["CitySearch"] == "" {
		fmt.Println("No filter")
		// if there isn't filters, we execute the templates with the data from Json files (render like mainPage)
		file.Execute(w, data.ArtistsData)
	} else {
		// if there are filters, we call the function dataToPush and execute the template with filter
		toPush := dataToPush(filters, data)

		// Render template with artist matching filters
		file.Execute(w, toPush)
	}

}

// dataToPush function test filters on every artist and save matching artist in okFilters
func dataToPush(filters map[string]string, data Datas) []Artists {
	var okFilters Datas
	// browses the data from ArtistsData and compares with all the filters
	for id, arti := range data.ArtistsData {
		// try for every artist if it match with filters
		if testMembers, _ := strconv.Atoi(filters["members"]); len(arti.MEMBERS) == testMembers {
			// next filter (Number of members)
			if testDate, _ := strconv.Atoi(filters["creationDate"]); arti.CREA_DATE == testDate || testDate == 0 {
				// next filter (creationDate)
				if arti.FIRST_ALBUM == filters["firstAlbum"] || filters["firstAlbum"] == "" {
					//next filter (FirstAlbum)
					for _, City := range data.LocationsData[id].LOCATIONS {
						// next filter (Locations)
						if filters["CitySearch"] == City || filters["CitySearch"] == "" {
							// Add the result to okFilters of all the filters
							okFilters.ArtistsData = append(okFilters.ArtistsData, data.ArtistsData[id])
							// print if artist matched
							fmt.Println(arti.NAME)
							break
						}
					}
				}
			}
		}
	}

	return okFilters.ArtistsData
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
	// mainPage Handler
	http.HandleFunc("/mainPage", mainPageHandler)
	// artist handler
	http.HandleFunc("/artists/", artistHandler)
	// fliter handler
	http.HandleFunc("/filter", filterHandler)

	fmt.Println("Server is starting...\n")
	fmt.Println("Go on http://localhost:8080/\n")
	fmt.Println("To shut down the server press CTRL + C")

	// Starting serveur
	http.ListenAndServe(":8080", nil)
}
