package main

// Req and Resp here

type UserPostRequest struct {
	Name string `json:"name"`
}

type UserResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type RelationPutRequest struct {
	State string `json:"state"`
}

type RelationGetResponse struct {
	UserID int64  `json:"user_id"`
	State string  `json:"state"`
	Type   string `json:"type"`
}

type RelationPutResponse RelationGetResponse

const (
	Like    string = "liked"
	Dislike string = "disliked"
	Match   string = "match"
)
