package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"web/pkg/driver/sqlite"
)

type Calculator struct {
	db sqlite.SQLite
}

func main() {
	db, err := sqlite.NewSQLite()
	if err != nil {
		log.Fatal(err)
	}
	calc := Calculator{
		db: db,
	}
	fmt.Println("Ascult 80")
	http.HandleFunc("/", serveFiles)
	http.HandleFunc("/api/solve", calc.solveIt)
	log.Fatal(http.ListenAndServe(":80", nil))
}

func (calc *Calculator) solveIt(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		writeResponse(w, "Eroare la verificare forma "+err.Error(), http.StatusBadRequest)
		return
	}
	a, err := strconv.Atoi(r.FormValue("a"))
	if err != nil {
		writeResponse(w, "Eroare la verificare A "+err.Error(), http.StatusBadRequest)
		return
	}
	b, err := strconv.Atoi(r.FormValue("b"))
	if err != nil {
		writeResponse(w, "Eroare la verificare B "+err.Error(), http.StatusBadRequest)
		return
	}
	c, err := strconv.Atoi(r.FormValue("c"))
	if err != nil {
		writeResponse(w, "Eroare la verificare C "+err.Error(), http.StatusBadRequest)
		return
	}

	result, err := calc.db.ReadCalculation(a, b, c)
	if err != nil {
		writeResponse(w, "Eroare la incarcarea calculelor "+err.Error(), http.StatusInternalServerError)
		return
	}

	if result != nil {
		//fmt.Printf("There are some datas: delta = %v", result.Delta)

		jsonBytes, err := json.Marshal(result)
		if err != nil {
			writeResponse(w, "Eroare la Marshal in JSON "+err.Error(), http.StatusInternalServerError)
			return
		}
		writeResponse(w, string(jsonBytes), http.StatusOK)
		return
	}
	newResult := calculateResult(a, b, c)
	err = calc.db.PutCalculation(newResult)
	if err != nil {
		writeResponse(w, "Eroare la incarcare rezultate in DB "+err.Error(), http.StatusInternalServerError)
		return
	}

	jsonBytes, err := json.Marshal(newResult)
	if err != nil {
		writeResponse(w, "Eroare la Marshal in JSON "+err.Error(), http.StatusInternalServerError)
		return
	}
	writeResponse(w, string(jsonBytes), http.StatusOK)
	return
}

func writeResponse(w http.ResponseWriter, reponse string, statusCode int) {
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, reponse)
}

func calculateResult(a, b, c int) sqlite.ResultCalc {
	var result sqlite.ResultCalc
	result.A = a
	result.B = b
	result.C = c
	result.Delta = b*b - 4*a*c
	result.IsValid = true
	if result.Delta > 0 {
		result.X1 = (float64(-b) + math.Sqrt(float64(result.Delta))) / (2 * float64(a))
		result.X2 = (float64(-b) - math.Sqrt(float64(result.Delta))) / (2 * float64(a))
	} else if result.Delta == 0 {
		result.X1 = float64(-b / 2 * a)
		result.X2 = result.X1
	} else {
		result.IsValid = false
	}
	return result
}

func serveFiles(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	p := "." + r.URL.Path
	if p == "./" {
		p = "./client/index.html"
	}
	http.ServeFile(w, r, p)
}
