package client

type OrganizationMembership struct {
	ID   string  `jsonapi:"primary,organization_memberships"`
	User *Person `jsonapi:"relation,person"`
}

type Deal struct {
	ID   string `jsonapi:"primary,deals"`
	Name string `jsonapi:"attr,name"`
}

type Service struct {
	ID   string `jsonapi:"primary,services"`
	Name string `jsonapi:"attr,name"`
}

type Project struct {
	deal    *Deal
	service *Service
}

type TimeEntry struct {
	ID     string   `jsonapi:"primary,time-entries"`
	Date   string   `jsonapi:"attr,date"`
	Note   string   `jsonapi:"attr,note"`
	Time   string   `jsonapi:"attr,time"`
	User   *Person  `jsonapi:"relation,person"`
	Budget *Service `jsonapi:"relation,service"`
}

type SessionRequest struct {
	ID       string `jsonapi:"primary,sessions"`
	Email    string `jsonapi:"attr,email"`
	Password string `jsonapi:"attr,password"`
}

type SessionResponse struct {
	ID                  string `jsonapi:"primary,sessions"`
	Token               string `jsonapi:"attr,token"`
	TokenExpirationDate string `jsonapi:"attr,token_expires_at"`
}

func NewTimeEntry(notes string, duration string, userID string, service *Service, day string) *TimeEntry {
	return &TimeEntry{
		ID:     "0",
		Date:   day,
		Note:   notes,
		Time:   duration,
		User:   &Person{ID: userID},
		Budget: service,
	}
}

type Person struct {
	ID string `jsonapi:"primary,people"`
}
