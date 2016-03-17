package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// SMTP contains the fields needed to connect to a SMTP server.
type SMTP struct {
	Port     int64  `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	From     string `json:"from"`
	Alias    string `json:"alias"`
}

// FromAddress creates a string suitable for use in an Email's From header.
func (c SMTP) FromAddress() string {
	if c.Alias != "" {
		return fmt.Sprintf(`"%s" <%s>`, c.Alias, c.From)
	}
	return fmt.Sprintf("<%s>", c.From)
}

// Address will return a string of the host and port separated by a colon.
func (c SMTP) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// ParseSMTPPath will create an SMTP using the given filepath.
func ParseSMTPPath(path string) (c SMTP, err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	err = json.NewDecoder(f).Decode(&c)
	return
}
