package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
)

var _status = StatusC{}
var _ret = Retorno{}
var _retError = RetornoError{}
var _apiret = apiCashResponse{}
var _retBox = RetornoBox{}

var config = ReadConfig()
var collection = getSession().DB(config.Database).C(config.Table)

// getSession ... < Some lines that describe your function>
func getSession() *mgo.Session {
	fmt.Println("Host: ", config.Dbhost)
	session, err := mgo.Dial(config.Dbhost)

	if err != nil {
		panic(err)
	}

	return session
}

// ResponseComplete ... < Some lines that describe your function>
func ResponseComplete(w http.ResponseWriter, code int, ret Retorno) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ret)
}

// ResponseCompleteBox ... < Some lines that describe your function>
func ResponseCompleteBox(w http.ResponseWriter, code int, ret RetornoBox) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ret)
}

// ResponseBasic ... < Some lines that describe your function>
func ResponseBasic(w http.ResponseWriter, code int, ret RetornoError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ret)
}

// Index ... < Some lines that describe your function>
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Esta API esta diseñada para ser una prueba para los nuevos candidatos al equipo.")
}

// Beers ... < Some lines that describe your function>
func beers(w http.ResponseWriter, r *http.Request) {

	/*b, err := json.MarshalIndent(_beers, "", "\t")
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)*/

	var results []Beer

	err := collection.Find(nil).Sort("-id").All(&results)

	var count = len(results)
	var code int

	if err != nil {
		log.Fatal(err)
		code = 500
		_status = StatusC{
			500, "Error contactarse con el adminsitrador"}
		_retError = RetornoError{
			_status}
	} else if count == 0 {
		code = 400
		_status = StatusC{
			400, "Sin reusltados"}
		_retError = RetornoError{
			_status}
	} else {
		_status = StatusC{
			200, "Operacion exitosa"}
		_ret = Retorno{
			_status,
			results}
		ResponseComplete(w, 200, _ret)
		return
	}

	ResponseBasic(w, code, _retError)

}

// beersIng ... < Some lines that describe your function>
func beersIng(w http.ResponseWriter, r *http.Request) {

	/*b, err := json.MarshalIndent(_beers, "", "\t")
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)*/
	decoder := json.NewDecoder(r.Body)
	var _beerdata Beer
	var code int = 0
	err4 := decoder.Decode(&_beerdata)
	if err4 != nil {
		log.Fatal(err4)
		code = 500
		_status = StatusC{
			code, "Request invalida"}
		_retError = RetornoError{
			_status}
		ResponseBasic(w, code, _retError)
		return
	}

	var results []BeerPrice
	err3 := collection.Find(bson.M{"id": _beerdata.ID}).All(&results)

	if err3 != nil || len(results) > 0 {
		code = 409
		_status = StatusC{
			code, "El ID de la cerveza ya existe"}
		_retError = RetornoError{
			_status}
		ResponseBasic(w, code, _retError)
		return
	}

	defer r.Body.Close()

	err3 = collection.Insert(_beerdata)
	if err3 != nil {
		code = 400
		_status = StatusC{
			code, "Request invalida"}
		_retError = RetornoError{
			_status}
		ResponseBasic(w, code, _retError)
	}

	if code == 0 {
		code = 200
		_status = StatusC{
			201, "Cerveza creada"}
		_retError = RetornoError{
			_status}
	}

	//_beers = append(_beers, _beerdata)
	ResponseBasic(w, code, _retError)

}

// beersBoxprice ... < Some lines that describe your function>
func beersBoxprice(w http.ResponseWriter, r *http.Request) {

	/*b, err := json.MarshalIndent(_beers, "", "\t")
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)*/
	params := mux.Vars(r)

	_beerid := params["beerID"]

	IDBeer, err1 := strconv.Atoi(_beerid)
	if err1 != nil {
		log.Fatal(err1)
	}

	var results []BeerPrice
	err := collection.Find(bson.M{"id": IDBeer}).All(&results)

	if err != nil || len(results) == 0 {
		_status = StatusC{
			404, "El Id de la cerveza no existe"}
		_retError = RetornoError{
			_status}
		ResponseBasic(w, 404, _retError)
		return
	}

	moneda := string(results[0].Currency)

	APIURL := config.Urlapi + config.Keyapi + config.Keymoneda + moneda + config.Keyformat
	req, err := http.NewRequest(http.MethodGet, APIURL, nil)
	if err != nil {
		panic(err)
	}
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	respAPI := _apiret
	json.Unmarshal(body, &respAPI)

	var cantidad float64 = 6
	total_precio := (cantidad * float64(respAPI.Quotes["USD"+moneda]))

	results[0].Quantity = cantidad
	results[0].PriceBox = total_precio

	_status = StatusC{
		200, "Operacion exitosa"}
	_retBox = RetornoBox{
		_status,
		results}

	ResponseCompleteBox(w, 200, _retBox)

	/*var data map[string]interface{}
	err2 := json.Unmarshal(body, &data)
	if err2 != nil {
		panic(err2)
	}
	//fmt.Println(data)
	//fmt.Printf("%#v", data)

	fetchValue(data)*/
}

// BeerDetail ... < Some lines that describe your function>
func BeerDetail(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	_beerid := params["beerID"]

	/*if !bson.IsObjectIdHex(_beerid) {
		_status = StatusC{
			404, "El Id de la cerveza no existe"}
		_ret = Retorno{
			_status,
			nil}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(_ret)
		return
	}

	//oid := bson.ObjectIdHex(_beerid)*/
	var results []Beer

	IDBeer, err1 := strconv.Atoi(_beerid)
	if err1 != nil {
		log.Fatal(err1)
	}

	err := collection.Find(bson.M{"id": IDBeer}).All(&results)

	if err != nil || len(results) == 0 {
		_status = StatusC{
			404, "El Id de la cerveza no existe"}
		_retError = RetornoError{
			_status}
		ResponseBasic(w, 404, _retError)
		return
	}

	_status = StatusC{
		200, "Operacion exitosa"}
	_ret = Retorno{
		_status,
		results}

	ResponseComplete(w, 200, _ret)

}
