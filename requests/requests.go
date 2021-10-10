package requests

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"proxy/models"
	"proxy/proxy"
	"proxy/saver"
	"strings"
)

type Handlers struct {
	Proxy proxy.Proxy
	Saver *saver.Saver
}

func Respond(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func (h *Handlers) GetRequests(w http.ResponseWriter, r *http.Request) {
	collection := h.Saver.GetClient().Database("test").Collection("requests")

	var results []*models.ReplyRequest

	findOptions := options.Find()

	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {
		var elem models.ReplyRequest
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, &elem)
	}

	Respond(w, 200, results)
}

func (h *Handlers) GetRequestByID(w http.ResponseWriter, r *http.Request) {
	collection := h.Saver.GetClient().Database("test").Collection("requests")

	var result models.ReplyRequest

	params := mux.Vars(r)
	ID, _ := params["id"]

	objectId, _ := primitive.ObjectIDFromHex(ID)

	filter := bson.D{{"_id", objectId}}

	_ = collection.FindOne(context.TODO(), filter).Decode(&result)

	Respond(w, 200, result)
}

func (h *Handlers) RepeatRequest(w http.ResponseWriter, r *http.Request) {collection := h.Saver.GetClient().Database("test").Collection("requests")

	var result models.ReplyRequest

	params := mux.Vars(r)
	ID, _ := params["id"]

	objectId, _ := primitive.ObjectIDFromHex(ID)

	filter := bson.D{{"_id", objectId}}

	_ = collection.FindOne(context.TODO(), filter).Decode(&result)

	var repeatReq *http.Request

	var url string
	url = result.Path
	if result.HttpMethod == "GET" {
		if len(result.GetParams) != 0 {
			url += "?"
		}
		for key, value := range result.GetParams {
			url += key + "=" + strings.Join(value, ",")
		}
	} else {
		for key, value := range result.PostParams {
			repeatReq.PostForm.Set(key, value)
		}
	}

	repeatReq, _ = http.NewRequest(result.HttpMethod, url, strings.NewReader(result.Body))

	for key, value := range result.Headers {
		repeatReq.Header.Set(key, strings.Join(value, ""))
	}

	for key, value := range result.Cookie {
		repeatReq.AddCookie(&http.Cookie{
			Name:  key,
			Value: value,
		})
	}

	if repeatReq.Method == "CONNECT" {
		h.Proxy.HttpsHandler(w, repeatReq)
	} else {
		h.Proxy.HttpHandler(w, repeatReq)
	}
}
