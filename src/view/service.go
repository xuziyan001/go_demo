package main

import "github.com/go-pg/pg/orm"

/*
Some tool functions here
*/

type Conn interface {
	Insert(model ...interface{}) error
	Model(model ...interface{}) *orm.Query
}

type SimpleService struct {
	Connect Conn
}

// AddUser is for add user to database
func (s *SimpleService) AddUser(u *User) error {
	err := s.Connect.Insert(u)
	return err
}

// GetUsers is for get all users from database
func (s *SimpleService) GetUsers() (*[]User, error) {
	var users []User
	err := s.Connect.Model(&users).Select()
	return &users, err
}

// GetRelationships is for get all relationships by UserID
func (s *SimpleService) GetRelationships(id int) (*[]Relationship, error) {
	var relations []Relationship
	err := s.Connect.Model(&relations).Where("UserID=?", id).Select()
	return &relations, err
}

// GetRelationship is for get relationship between two users
func (s *SimpleService) GetRelationship(id int, otherid int) (*Relationship, error) {
	r := new(Relationship)
	err := s.Connect.Model(r).
		Where("UserID=?", id).
		Where("RelateUserID=?", otherid).
		Limit(1).Select()
	return r, err
}

// AddRelationship is for insert relationship
func (s *SimpleService) AddRelationship(r *Relationship) error {
	err := s.Connect.Insert(r)
	return err
}

// ModifyRelationshipState is for modify state of a relationship
func (s *SimpleService) ModifyRelationshipState(r *Relationship, state string) error {
	r.State = state
	_, err := s.Connect.Model(r).WherePK().Column("state").Update()
	return err
}
