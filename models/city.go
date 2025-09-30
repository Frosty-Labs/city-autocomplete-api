package models

// City represents a city with its information
type City struct {
	Name       string `json:"name"`
	Country    string `json:"country"`
	Subcountry string `json:"subcountry"`
	GeonameID  string `json:"geonameid"`
}
