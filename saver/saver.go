package saver

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
	"net/http"
	"proxy/models"
)

type Saver struct {
	client *mongo.Client
}

func (s *Saver) MongoConnect() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	s.client = client
	fmt.Println(err)

	collection := s.client.Database("test").Collection("requests")
	collection.Drop(context.TODO())
	return
}

func (s *Saver) SaveRequest(req *http.Request) {
	var parseReq models.Request

	parseReq.HttpMethod = req.Method

	parseReq.Path = req.URL.Host + req.URL.Path

	parseReq.GetParams = make(map[string][]string)

	if req.Method == "GET" {
		for key, value := range req.URL.Query() {
			parseReq.GetParams[key] = value
		}
	}

	parseReq.Headers = make(map[string][]string)

	for key, value := range req.Header {
		parseReq.Headers[key] = value
	}

	parseReq.Cookie = make(map[string]string)

	for _, cookie := range req.Cookies() {
		parseReq.Cookie[cookie.Name] = cookie.Value
	}

	byt := []byte{}
	read, _ := req.Body.Read(byt)

	var dat map[string]string

	if read != 0 {
		if err := json.Unmarshal(byt, &dat); err != nil {
			panic(err)
		}
	}

	if parseReq.PostParams == nil {
		parseReq.PostParams = make(map[string]string)
	}

	for key, value := range dat {
		parseReq.PostParams[key] = value
	}

	collection := s.client.Database("test").Collection("requests")

	_, err := collection.InsertOne(context.TODO(), parseReq)

	if err != nil {
		log.Fatal(err)
	}

}

func (s *Saver) SaveResponse(resp *http.Response) {
	var parseResp models.Response

	parseResp.Code = resp.StatusCode

	parseResp.Message = resp.Status

	parseResp.Headers = make(map[string][]string)
	for key, value := range resp.Header {
		parseResp.Headers[key] = value
	}

	if resp.Body != nil {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		parseResp.Body = string(bodyBytes)
	}

	collection := s.client.Database("test").Collection("responses")

	_, err := collection.InsertOne(context.TODO(), parseResp)

	if err != nil {
		log.Fatal(err)
	}
}

func (s *Saver) GetClient() *mongo.Client {
	return s.client
}
