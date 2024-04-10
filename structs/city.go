package structs

type City struct {
	Id         uint64 `json:"id"`
	City       string `json:"city"`
	Country    string `json:"country"`
	Population uint64 `json:"population"`
}
