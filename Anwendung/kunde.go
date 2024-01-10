package Anwendung

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Kunde struct {
	Collection             *mongo.Collection
	ProduktMicroserviceURL string
}

type KundeModel struct {
	ID               primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Vorname          string             `json:"vorname,omitempty" bson:"vorname,omitempty"`
	Nachname         string             `json:"nachname,omitempty" bson:"nachname,omitempty"`
	Email            string             `json:"email,omitempty" bson:"email,omitempty"`
	GekaufteProdukte []GekauftesProdukt `json:"gekaufteProdukte,omitempty" bson:"gekaufteProdukte,omitempty"`
}
type ProduktMenge struct {
	Quantity int `json:"quantity,omitempty" bson:"quantity,omitempty"`
}

type GekauftesProdukt struct {
	ProduktID primitive.ObjectID `json:"produktID,omitempty" bson:"produktID,omitempty"`
	Quantity  int                `json:"quantity,omitempty" bson:"quantity,omitempty"`
}

func (p *Kunde) Create(w http.ResponseWriter, r *http.Request) {
	var kunde KundeModel
	if err := json.NewDecoder(r.Body).Decode(&kunde); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Platzhalter für gekaufte Produkte
	kunde.GekaufteProdukte = []GekauftesProdukt{}

	result, err := p.Collection.InsertOne(context.TODO(), kunde)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result.InsertedID)
}

func (p *Kunde) List(w http.ResponseWriter, r *http.Request) {
	cursor, err := p.Collection.Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var kunde []KundeModel
	if err := cursor.All(context.TODO(), &kunde); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(kunde)
}

func (p *Kunde) GetByID(w http.ResponseWriter, r *http.Request) {
	kundeID := chi.URLParam(r, "id")
	objID, err := primitive.ObjectIDFromHex(kundeID)
	if err != nil {
		http.Error(w, "Ungültige Kunden ID", http.StatusBadRequest)
		return
	}

	var kunde KundeModel
	err = p.Collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&kunde)
	if err != nil {
		http.Error(w, "Kunde konnte nicht gefunden werden", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(kunde)
}

func (p *Kunde) UpdateByID(w http.ResponseWriter, r *http.Request) {
	kundeID := chi.URLParam(r, "id")
	objID, err := primitive.ObjectIDFromHex(kundeID)
	if err != nil {
		http.Error(w, "Ungültige Kunden ID", http.StatusBadRequest)
		return
	}

	var updatedKunde KundeModel
	if err := json.NewDecoder(r.Body).Decode(&updatedKunde); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = p.Collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": objID},
		bson.M{"$set": updatedKunde},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (p *Kunde) DeleteByID(w http.ResponseWriter, r *http.Request) {
	kundeID := chi.URLParam(r, "id")
	objID, err := primitive.ObjectIDFromHex(kundeID)
	if err != nil {
		http.Error(w, "Ungültige Kunden ID", http.StatusBadRequest)
		return
	}

	_, err = p.Collection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (k *Kunde) Kaufen(w http.ResponseWriter, r *http.Request) {
	kundeID := chi.URLParam(r, "id")
	objID, err := primitive.ObjectIDFromHex(kundeID)
	if err != nil {
		http.Error(w, "Ungültige Kunden ID", http.StatusBadRequest)
		return
	}

	var kaufanfrage struct {
		ProduktID   primitive.ObjectID `json:"produktID,omitempty" bson:"produktID,omitempty"`
		Quantity    int                `json:"quantity,omitempty" bson:"quantity,omitempty"`
		ProduktName string             `json:"name,omitempty" bson:"name,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&kaufanfrage); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Rufe die aktuelle Produktmenge ab
	productMicroserviceURL := k.ProduktMicroserviceURL + "/product/" + kaufanfrage.ProduktID.Hex()
	response, err := http.Get(productMicroserviceURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		http.Error(w, "Fehler beim Abrufen der Produktmenge", http.StatusInternalServerError)
		return
	}

	var produktMenge struct {
		Quantity int `json:"quantity,omitempty" bson:"quantity,omitempty"`
	}

	if err := json.NewDecoder(response.Body).Decode(&produktMenge); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if kaufanfrage.Quantity > produktMenge.Quantity {
		http.Error(w, "Nicht genügend Produktbestand", http.StatusBadRequest)
		return
	}

	aktualisierteMenge := produktMenge.Quantity - kaufanfrage.Quantity

	// Rufe den "Produktkauf" im Produkt-Microservice auf
	updateRequest := struct {
		Quantity int    `json:"quantity,omitempty" bson:"quantity,omitempty"`
		Name     string `json:"name,omitempty" bson:"name,omitempty"`
	}{
		Quantity: aktualisierteMenge,
		Name:     kaufanfrage.ProduktName,
	}

	updateRequestBody, err := json.Marshal(updateRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	productMicroserviceURL = k.ProduktMicroserviceURL + "/product/" + kaufanfrage.ProduktID.Hex()
	req, err := http.NewRequest("PUT", productMicroserviceURL, bytes.NewBuffer(updateRequestBody))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err = client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		http.Error(w, "Fehler beim Produktkauf", http.StatusInternalServerError)
		return
	}

	// Aktualisiere die Liste der gekauften Produkte im Kunden
	gekauftesProdukt := GekauftesProdukt{
		ProduktID: kaufanfrage.ProduktID,
		Quantity:  kaufanfrage.Quantity,
	}

	_, err = k.Collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": objID},
		bson.M{"$push": bson.M{"gekaufteProdukte": gekauftesProdukt}},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
