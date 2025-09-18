package config

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
)

// Config is the parent configuration struct and includes fields for single
// configurations of a database, cookie, and SMTP connection.
type Config struct {
	HTTPS       bool     `json:"https"`
	Domain      string   `json:"domain"`
	ProxyDomain string   `json:"proxy_domain"`
	Port        int      `json:"port"`
	ProxyPort   int      `json:"proxy_port"`
	TemplateDir string   `json:"templates"`
	AbsPath     string   `json:"abs_path"`
	MediaDir    string   `json:"media"`
	MediaURL    string   `json:"media_url"`
	StaticDir   string   `json:"static"`
	StaticURL   string   `json:"static_url"`
	SecretKey   string   `json:"secret_key"`
	Version     string   `json:"version"`
	Database    Database `json:"database"`
	Cookie      Cookie   `json:"cookie"`
	SMTP        SMTP     `json:"smtp"`
	Metadata    Metadata `json:"metadata"`
}

// Address returns the domain:port pair.
func (c Config) Address() string {
	return fmt.Sprintf("%s:%d", c.Domain, c.Port)
}

// URL returns the domain:port scheme. Port is omitted if 80.
func (c Config) URL() (u *url.URL) {
	u = &url.URL{}
	if c.HTTPS {
		u.Scheme = "https"
	} else {
		u.Scheme = "http"
	}
	// Fallback to the non proxy domain and ports
	domain := c.ProxyDomain
	if domain == "" {
		domain = c.Domain
	}
	port := c.ProxyPort
	if port == 0 {
		port = c.Port
	}
	if port == 80 {
		u.Host = domain
	} else {
		u.Host = fmt.Sprintf("%s:%d", domain, port)
	}
	return
}

// FullAddress returns the scheme, domain, port, and host - including
// proxy info
func (c Config) FullAddress() string {
	return c.URL().String()
}

// StaticAddress adds the static URL to the full address
// TODO If staticURL is already a valid URL just return that
func (c Config) StaticAddress() string {
	url := c.URL()
	url.Path = c.StaticURL
	return url.String()
}

// Parse will create a Config using the file settings.json in the
// current directory.
func Parse() (Config, error) {
	return ParsePath("./settings.json")
}

// ParsePath will create a Config using the file at the given path.
func ParsePath(filename string) (Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return Config{}, err
	}
	return parse(f)
}

func parse(f io.Reader) (Config, error) {
	c := Config{Cookie: DefaultCookie}
	if err := json.NewDecoder(f).Decode(&c); err != nil {
		return c, err
	}
	return c, nil
}

// Default is a basic configuration with insecure values. It will return the
// Address localhost:8080
var Default = Config{
	Cookie:    DefaultCookie,
	Port:      8080,
	StaticURL: "/static/",
	Metadata:  Metadata{},
}
