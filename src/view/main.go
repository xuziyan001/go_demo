package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-pg/pg"

	"github.com/gorilla/mux"
)

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
}

func InternalServerError(w http.ResponseWriter) {
	http.Error(w, "500 interal server error", http.StatusInternalServerError)
}

func BadRequest(w http.ResponseWriter) {
	http.Error(w, "400 bad request", http.StatusBadRequest)
}

func MethodNotAllowedHandler() http.Handler { return http.HandlerFunc(MethodNotAllowed) }

func main() {
	r := mux.NewRouter()
	r.MethodNotAllowedHandler = MethodNotAllowedHandler()
	r.HandleFunc("/users", UserGetHandler).Methods("GET")
	r.HandleFunc("/users", UserPostHandler).Methods("POST")
	r.HandleFunc("/users/{user_id:[0-9]+}/relationships",
		RelationshipGetHander).Methods("GET")
	r.HandleFunc("/users/{user_id:[0-9]+}/relationships/{other_user_id:[0-9]+}",
		RelationshipPutHander).Methods("PUT")
	log.Fatal(http.ListenAndServe("localhost:8081", r))
}

func UserGetHandler(w http.ResponseWriter, r *http.Request) {
	var response []UserResponse
	s := SimpleService{Connect: db}
	users, err := s.GetUsers()
	if err != nil {
		// need log here
		InternalServerError(w)
		return
	}
	for _, user := range *users {
		user := user
		response = append(response, UserResponse{user.Id, user.Name, user.Type})
	}
	if data, err := json.MarshalIndent(response, "", "  "); err != nil {
		InternalServerError(w)
	} else {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%s\n", data)
	}
}

func UserPostHandler(w http.ResponseWriter, r *http.Request) {
	var d = UserPostRequest{}
	var response UserResponse
	s := SimpleService{Connect: db}
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		BadRequest(w)
		return
	} else {
		user := &User{
			Name: d.Name,
		}
		if err := s.AddUser(user); err != nil {
			InternalServerError(w)
			return
		}
		response = UserResponse{user.Id, user.Name, user.Type}
		if data, err := json.MarshalIndent(response, "", "  "); err != nil {
			InternalServerError(w)
		} else {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "%s\n", data)
		}
	}
}

func RelationshipGetHander(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// fmt.Fprintf(w, "get relationships, user_id=%s\n", vars["user_id"])
	userID, _ := strconv.Atoi(vars["user_id"])
	s := SimpleService{Connect: db}
	userRelationResp := make([]RelationGetResponse, 0)
	relations, err := s.GetRelationships(userID)
	if err != nil {
		InternalServerError(w)
		return
	}
	for _, relation := range *relations {
		relation := relation
		u := RelationGetResponse{relation.RelateUserID, relation.State, "relationship"}
		userRelationResp = append(userRelationResp, u)
	}
	if data, err := json.MarshalIndent(userRelationResp, "", "  "); err != nil {
		InternalServerError(w)
	} else {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%s\n", data)
	}
}

func RelationshipPutHander(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["user_id"])
	otherUserID, _ := strconv.Atoi(vars["other_user_id"])
	// fmt.Fprintf(w, "put relationship, user_id=%s, other_user_id=%s\n", vars["user_id"], vars["other_user_id"])
	tx, err := db.Begin()
	if err != nil {
		InternalServerError(w)
		return
	}
	// lambda for rollback and return
	rb := func(e error, tx *pg.Tx) {
		tx.Rollback()
		InternalServerError(w)
	}
	s := SimpleService{Connect: tx}
	d := RelationPutRequest{}
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		BadRequest(w)
		return
	} else {
		response := new(RelationPutResponse)
		userRelation, err := s.GetRelationship(userID, otherUserID)
		if err != nil {
			rb(err, tx)
			return
		}
		otherUserRelation, err := s.GetRelationship(otherUserID, userID)
		if err != nil {
			rb(err, tx)
			return
		}
		if d.State == Dislike {
			if userRelation.State == "" { // empty, do create
				userRelation = &Relationship{
					UserID:       int64(userID),
					RelateUserID: int64(otherUserID),
					State:        d.State}
				if err := s.AddRelationship(userRelation); err != nil {
					rb(err, tx)
					return
				}
			} else if userRelation.State == d.State { // do nothing
			} else { //update
				if err := s.ModifyRelationshipState(userRelation, d.State); err != nil {
					rb(err, tx)
					return
				}
			}
			if otherUserRelation.State == Match { // only 'match' need to be changed
				if err := s.ModifyRelationshipState(otherUserRelation, Like); err != nil {
					rb(err, tx)
					return
				}
			}
			// construct response
			response = &RelationPutResponse{int64(otherUserID), d.State, "relationship"}
			if data, err := json.MarshalIndent(response, "", "  "); err != nil {
				InternalServerError(w)
				return
			} else {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, "%s\n", data)
			}
		} else if d.State == Like {
			state := d.State
			if otherUserRelation.State == Like {
				state = Match
				if err := s.ModifyRelationshipState(otherUserRelation, state); err != nil {
					rb(err, tx)
					return
				}
			}
			if userRelation.State == "" { // empty, do create
				userRelation = &Relationship{
					UserID:       int64(userID),
					RelateUserID: int64(otherUserID),
					State:        state}
				if err := s.AddRelationship(userRelation); err != nil {
					rb(err, tx)
					return
				}
			} else if userRelation.State == state { // do nothing
			} else { //update
				if err := s.ModifyRelationshipState(userRelation, state); err != nil {
					rb(err, tx)
					return
				}
			}
			// TODO: should we connect responses together?
			response = &RelationPutResponse{int64(otherUserID), state, "relationship"}
			if data, err := json.MarshalIndent(response, "", "  "); err != nil {
				InternalServerError(w)
			} else {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, "%s\n", data)
			}
		} else {
			BadRequest(w)
		}
	}
	if err := tx.Commit(); err != nil {
		rb(err, tx)
	}
}
