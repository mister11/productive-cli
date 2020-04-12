package model

import (
	"strconv"
)

type OrganizationMembership struct {
	ID   string  `jsonapi:"primary,organization_memberships"`
	User *Person `jsonapi:"relation,person"`
}

type Deal struct {
	ID      string `jsonapi:"primary,deals"`
	Name    string `jsonapi:"attr,name"`
	EndDate string `jsonapi:"attr,end_date"`
}

type Service struct {
	ID   string `jsonapi:"primary,services"`
	Name string `jsonapi:"attr,name"`
}

type TimeEntry struct {
	ID     string   `jsonapi:"primary,time-entries"`
	Date   string   `jsonapi:"attr,date"`
	Note   string   `jsonapi:"attr,note"`
	Time   string   `jsonapi:"attr,time"`
	User   *Person  `jsonapi:"relation,person"`
	Budget *Service `jsonapi:"relation,service"`
}

func NewTimeEntry(notes string, duration int, userID string, service *Service, date string) *TimeEntry {
	return &TimeEntry{
		ID:     "0",
		Date:   date,
		Note:   notes,
		Time:   strconv.Itoa(duration),
		User:   &Person{ID: userID},
		Budget: service,
	}
}

type Person struct {
	ID string `jsonapi:"primary,people"`
}
