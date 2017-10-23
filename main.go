package main

import (
	"archive/zip"
	"database/sql"
	"encoding/xml"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	dinesafe "github.com/rchristiani/dinesafe/api"
)

func downloadZipFile(filepath string, url string) (err error) {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	//Read the zip file
	r, err := zip.OpenReader(filepath)
	if err != nil {
		return err
	}
	defer r.Close()

	//Create the file for it to go into
	out, err = os.Create("dinesafe.xml")
	defer out.Close()
	if err != nil {
		return err
	}

	//Range over the content
	for _, f := range r.File {
		//Open each reader
		rc, err := f.Open()
		defer rc.Close()
		if err != nil {
			return err
		}
		//Turn that reader into bytes
		bytes, err := ioutil.ReadAll(rc)
		if err != nil {
			return err
		}
		//write to file
		out.Write(bytes)
	}

	return nil
}

func main() {
	err := downloadZipFile("dinesafe.zip", "http://opendata.toronto.ca/public.health/dinesafe/dinesafe.zip")

	if err != nil {
		log.Fatal(err)
	}

	//Open the XML file
	dinesafeXML, err := os.Open("dinesafe.xml")
	//Make sure to defer and close he files
	defer dinesafeXML.Close()
	if err != nil {
		log.Fatal(err)
	}
	//Take the file and read the bytes from it
	xmlBytes, err := ioutil.ReadAll(dinesafeXML)
	if err != nil {
		log.Fatal(err)
	}
	//Make a struct for everything
	var res dinesafe.Restaurants
	var inspections dinesafe.Inspections
	//Unmarshal the bytes into the Query
	xml.Unmarshal(xmlBytes, &res)
	xml.Unmarshal(xmlBytes, &inspections)

	db, err := sql.Open("postgres", "user=ryanchristiani dbname=dinesafe sslmode=disable")

	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	// Restaurants
	log.Println("========= Adding Restraurants =============")
	for _, restaurant := range res.Rows {
		_, err = db.Exec("INSERT INTO restaurants(establishmentID, establishmentName, establishmentType, establishmentAddress, establishmentStatus, minimumInspectionsPerYear) Values($1,$2,$3,$4,$5,$6) ON CONFLICT DO NOTHING;",
			restaurant.EstablishmentID, restaurant.EstablishmentName, restaurant.EstablishmentType, restaurant.EstablishmentAddress, restaurant.EstablishmentStatus, restaurant.MinimumInspectionsPerYear,
		)

		if err != nil {
			log.Fatalln(err)
		}

	}
	log.Println("========= Done Adding Restraurants =============")
	log.Println("========= Adding Inspections =============")

	// Inspections
	for _, inspection := range inspections.Inspections {
		_, err = db.Exec("INSERT INTO inspections(establishmentID,inspectionID,infractionDetails,inspectionDate,severity,action,courtOutcome,amountFinded) Values($1,$2,$3,$4,$5,$6,$7,$8) ON CONFLICT DO NOTHING;",
			inspection.EstablishmentID, inspection.InspectionID, inspection.InfractionDetails, inspection.InspectionDate, inspection.Severity, inspection.Action, inspection.CourtOutcome, inspection.AmountFinded,
		)

		if err != nil {
			log.Fatal(err)
		}
	}
	log.Println("========= Done Adding Inspections =============")

}
