package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/intervention-engine/fhir/models"
	"gopkg.in/mgo.v2/bson"
)

func EligibilityResponseIndexHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var result []models.EligibilityResponse
	c := Database.C("eligibilityresponses")
	iter := c.Find(nil).Limit(100).Iter()
	err := iter.All(&result)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	var eligibilityresponseEntryList []models.EligibilityResponseBundleEntry
	for _, eligibilityresponse := range result {
		var entry models.EligibilityResponseBundleEntry
		entry.Title = "EligibilityResponse " + eligibilityresponse.Id
		entry.Id = eligibilityresponse.Id
		entry.Content = eligibilityresponse
		eligibilityresponseEntryList = append(eligibilityresponseEntryList, entry)
	}

	var bundle models.EligibilityResponseBundle
	bundle.Type = "Bundle"
	bundle.Title = "EligibilityResponse Index"
	bundle.Id = bson.NewObjectId().Hex()
	bundle.Updated = time.Now()
	bundle.TotalResults = len(result)
	bundle.Entry = eligibilityresponseEntryList

	log.Println("Setting eligibilityresponse search context")
	context.Set(r, "EligibilityResponse", result)
	context.Set(r, "Resource", "EligibilityResponse")
	context.Set(r, "Action", "search")

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(rw).Encode(bundle)
}

func LoadEligibilityResponse(r *http.Request) (*models.EligibilityResponse, error) {
	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		return nil, errors.New("Invalid id")
	}

	c := Database.C("eligibilityresponses")
	result := models.EligibilityResponse{}
	err := c.Find(bson.M{"_id": id.Hex()}).One(&result)
	if err != nil {
		return nil, err
	}

	log.Println("Setting eligibilityresponse read context")
	context.Set(r, "EligibilityResponse", result)
	context.Set(r, "Resource", "EligibilityResponse")
	return &result, nil
}

func EligibilityResponseShowHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	context.Set(r, "Action", "read")
	_, err := LoadEligibilityResponse(r)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(rw).Encode(context.Get(r, "EligibilityResponse"))
}

func EligibilityResponseCreateHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	decoder := json.NewDecoder(r.Body)
	eligibilityresponse := &models.EligibilityResponse{}
	err := decoder.Decode(eligibilityresponse)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	c := Database.C("eligibilityresponses")
	i := bson.NewObjectId()
	eligibilityresponse.Id = i.Hex()
	err = c.Insert(eligibilityresponse)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	log.Println("Setting eligibilityresponse create context")
	context.Set(r, "EligibilityResponse", eligibilityresponse)
	context.Set(r, "Resource", "EligibilityResponse")
	context.Set(r, "Action", "create")

	host, err := os.Hostname()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	rw.Header().Add("Location", "http://"+host+":3001/EligibilityResponse/"+i.Hex())
}

func EligibilityResponseUpdateHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		http.Error(rw, "Invalid id", http.StatusBadRequest)
	}

	decoder := json.NewDecoder(r.Body)
	eligibilityresponse := &models.EligibilityResponse{}
	err := decoder.Decode(eligibilityresponse)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	c := Database.C("eligibilityresponses")
	eligibilityresponse.Id = id.Hex()
	err = c.Update(bson.M{"_id": id.Hex()}, eligibilityresponse)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	log.Println("Setting eligibilityresponse update context")
	context.Set(r, "EligibilityResponse", eligibilityresponse)
	context.Set(r, "Resource", "EligibilityResponse")
	context.Set(r, "Action", "update")
}

func EligibilityResponseDeleteHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		http.Error(rw, "Invalid id", http.StatusBadRequest)
	}

	c := Database.C("eligibilityresponses")

	err := c.Remove(bson.M{"_id": id.Hex()})
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Setting eligibilityresponse delete context")
	context.Set(r, "EligibilityResponse", id.Hex())
	context.Set(r, "Resource", "EligibilityResponse")
	context.Set(r, "Action", "delete")
}