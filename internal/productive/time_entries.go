package productive

import (
	"time"
)

type TimeEntryRequest struct {
	ID     string   `jsonapi:"primary,time-entries"`
	Date   string   `jsonapi:"attr,date"`
	Note   string   `jsonapi:"attr,note"`
	Time   string   `jsonapi:"attr,time"`
	User   *Person  `jsonapi:"relation,person"`
	Budget *Service `jsonapi:"relation,service"`
}

type Person struct {
	ID string `jsonapi:"primary,people"`
}

type timeEntryService struct {
	client *Client
}

func newTimeEntryService(client *Client) *timeEntryService {
	return &timeEntryService{
		client: client,
	}
}

func (s *timeEntryService) CreateTimeEntry(
	notes string,
	time string,
	userID string,
	service *Service,
	day time.Time,
	token string,
) error {
	timeEntry := &TimeEntryRequest{
		ID:     "0",
		Date:   formatDate(day),
		Note:   notes,
		Time:   time,
		User:   &Person{ID: userID},
		Budget: service,
	}
	req, err := s.client.NewRequest("POST", "time_entries", timeEntry, getHeaders(token))
	if err != nil {
		return err
	}
	_, err = s.client.Do(req)
	if err != nil {
		return err
	}
	return nil
}
