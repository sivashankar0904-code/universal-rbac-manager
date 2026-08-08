package model

import "time"

// CType defines the contact (work, personal)
type CType string

// AType defines the address type (communication, permanent, billing)
type AType string

// VerificationStatus defines different verification status
type VerificationStatus string

type Status struct {
	State     VerificationStatus `json:"state,omitempty"`
	CreatedAt *time.Time         `json:"createdAt,omitempty"`
	UpdatedAt *time.Time         `json:"updatedAt,omitempty"`
}

// Location defines the geographic location
type Location struct {
	Lat  string `json:"lat,omitempty"`
	Long string `json:"long,omitempty"`
}

// Contact has all details required for specifying contact details
type Contact struct {
	Email       string   `json:"email,omitempty"`
	Mobile      string   `json:"mobile,omitempty"`
	Landline    string   `json:"workPhone,omitempty"`
	OtherNumber []string `json:"othersNumber,omitempty"`
	Fax         string   `json:"fax,omitempty"`
	Type        CType    `json:"type,omitempty"`
	// Verification is used to set a verification status of individual field (whichever requires verification)
	Verification map[string]Status `json:"verification,omitempty"`
}

// Address has all the fields required to represent a postal address
type Address struct {
	ID             string    `json:"id,omitempty"`
	OrganizationID string    `json:"organizationID,omitempty"`
	ClientID       string    `json:"clientID,omitempty"`
	Building       string    `json:"building,omitempty"`
	Street         string    `json:"street,omitempty"`
	Road           string    `json:"road,omitempty"`
	Line1          string    `json:"line1,omitempty"`
	Line2          string    `json:"line2,omitempty"`
	PostalCode     string    `json:"postalCode,omitempty"`
	City           string    `json:"city,omitempty"`
	State          string    `json:"state,omitempty"`
	Country        string    `json:"country,omitempty"`
	Type           AType     `json:"type,omitempty"`
	Location       *Location `json:"location,omitempty"`
	// Verification is used to set a verification status of the address
	Verification VerificationStatus `json:"verification,omitempty"`
	Status       *Status            `json:"status,omitempty"`
}
