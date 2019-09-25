package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.impcloud.net/RSP-Inventory-Suite/food-safety-sample/app/tag"
)

const (
	notificationEndPoint = "/api/v1/notification"
	subscriptionEndPoint = "/api/v1/subscription"
	notificationCategory = "SECURITY"
	notificationSeverity = "CRITICAL"
)

// Notification holds the body schema to post a notification event to EdgeX
type Notification struct {
	Slug     string   `json:"slug"`
	Sender   string   `json:"sender"`
	Category string   `json:"category"`
	Severity string   `json:"severity"`
	Content  string   `json:"content"`
	Labels   []string `json:"labels"`
}

// Subscriber holds the body schema to register a subscriber to EdgeX
type Subscriber struct {
	Slug                 string     `json:"slug"`
	Receiver             string     `json:"receiver"`
	SubscribedCategories []string   `json:"subscribedCategories"`
	SubscribedLabels     []string   `json:"subscribedLabels"`
	Channels             []Channels `json:"channels"`
}

// Channels holds the body schema to specify different type of notification channels (email, SMS, REST post call)
type Channels struct {
	Type          string   `json:"type"`
	URL           string   `json:"url,omitempty"`
	MailAddresses []string `json:"mailAddresses,omitempty"`
}

// PostNotification sends a notification when group of tags reach freezer area
// This leverages EdgeX Alerts & notification service
func PostNotification(content string, notificationServiceURL string) error {

	notification := Notification{
		Slug:     "freezer-arrival-notification-" + time.Now().String(),
		Labels:   []string{"RSP"},
		Sender:   "Food safety app",
		Category: notificationCategory,
		Severity: notificationSeverity,
		Content:  content}

	requestBody, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	response, err := http.Post(notificationServiceURL+notificationEndPoint, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return fmt.Errorf("POST error on notification endpoint, StatusCode %d", response.StatusCode)
	}

	return nil

}

// CreateBodyContent composes the body of the notification message
func CreateBodyContent(tags []tag.Tag, temperature float32, readerAlias string) string {

	// Extract EPC value from tags
	epcSlice := make([]string, len(tags))
	for _, val := range tags {
		epcSlice = append(epcSlice, val.Epc)
	}

	body := ` 
	Item(s) has reached the %s.
	Current temperature in the area is %d.
	EPC(s): %s
	Date: %s
	`
	content := fmt.Sprintf(body, readerAlias, temperature, strings.Join(epcSlice, ","))

	return content
}

// RegisterSubscriber registers a subscriber to EdgeX Alerts & notification service using email as channel
func RegisterSubscriber(emails []string, notificationServiceURL string) error {

	// Create requestBody
	subscriber := new(Subscriber)
	channels := Channels{Type: "EMAIL", MailAddresses: emails}

	subscriber.Slug = "freezer-arrival-notification"
	subscriber.Receiver = "USER"
	subscriber.SubscribedCategories = []string{notificationCategory}
	subscriber.SubscribedLabels = []string{notificationCategory}
	subscriber.Channels = []Channels{channels}

	requestBody, err := json.Marshal(subscriber)
	if err != nil {
		return err
	}

	response, err := http.Post(notificationServiceURL+subscriptionEndPoint, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusConflict {
		return fmt.Errorf("POST error on subscriber endpoint, StatusCode %d", response.StatusCode)
	}

	return nil

}
