package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Request struct {
	ID         int64               `json:"id,omitempty" bson:"-"`
	HttpMethod string              `json:"HttpMethod,omitempty"`
	Path       string              `json:"path,omitempty"`
	Headers    map[string][]string `json:"headers,omitempty"`
	Cookie     map[string]string   `json:"cookies,omitempty"`
	GetParams  map[string][]string `json:"get-params,omitempty"`
	Body       string              `json:"body,omitempty"`
	PostParams map[string]string   `json:"post-params,omitempty"`
}

type ReplyRequest struct {
	ID         primitive.ObjectID  `bson:"_id" json:"id,omitempty"`
	HttpMethod string              `json:"HttpMethod,omitempty"`
	Path       string              `json:"path,omitempty"`
	Headers    map[string][]string `json:"headers,omitempty"`
	Cookie     map[string]string   `json:"cookies,omitempty"`
	GetParams  map[string][]string `json:"get-params,omitempty"`
	Body       string              `json:"body,omitempty"`
	PostParams map[string]string   `json:"post-params,omitempty"`
}
