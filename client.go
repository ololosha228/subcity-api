package subcity

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
)

var (
	apiHostname, _ = url.Parse("https://subscity.ru")
)

type City int

const (
	CityUnknown City = iota
	CityMoscow
	CitySpb
)

type Movie struct {
	AgeRestriction int               `json:"age_restriction"`
	Cast           []string          `json:"cast"`
	Countries      []string          `json:"countries"`
	CreatedAt      time.Time         `json:"created_at"`
	Description    string            `json:"description"`
	Directors      []string          `json:"directors"`
	Duration       int               `json:"duration"`
	Genres         []string          `json:"genres"`
	ID             int               `json:"id"`
	Languages      []string          `json:"languages"`
	PosterURL      string            `json:"poster"`
	Rating         map[string]Rating `json:"rating"`
	Title          map[string]string `json:"title"`
	Year           int               `json:"year"`
}

func (e *Movie) UnmarshalJSON(b []byte) error {
	if e == nil {
		e = new(Movie)
	}

	m := make(map[string]interface{})
	err := json.Unmarshal(b, &m)
	if err != nil {
		return err
	}

	m["description"] = strings.ReplaceAll(m["description"].(string), "\u00a0", " ")

	t, err := time.Parse("2006-01-02T15:04:05-07:00", m["created_at"].(string))
	if err != nil {
		return err
	}
	m["created_at"] = t

	d, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "json",
		Result:  e,
	})

	old := e
	err = d.Decode(&m)
	if err != nil {
		e = old
		return err
	}
	return nil
}

type Rating struct {
	ID     int     `json:"id"`
	Rating float32 `json:"rating"`
	Votes  int     `json:"votes"`
}

func GetMovies(city City, sort string) ([]*Movie, error) {
	if sort == "" {
		sort = "-id"
	}

	u := new(url.URL)
	*u = *apiHostname
	switch city {
	case CityMoscow:
		u.Host = "msk." + u.Host
	case CitySpb:
		u.Host = "spb." + u.Host
	default:
		return nil, errors.New("unknown city")
	}
	u.RawQuery = "sort=" + sort
	u.Path = "/movies.json"

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}

	data := make([]*Movie, 0)

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	return data, nil
}

type Cinema struct {
	ID       int `json:"id"`
	Location struct {
		Address string   `json:"address"`
		Metro   []string `json:"metro"`
		Lat     float64  `json:"latitude"`
		Lon     float64  `json:"longitude"`
	} `json:"location"`
	Movies []int    `json:"movies"`
	Name   string   `json:"name"`
	Phones []string `json:"phones"`
	Urls   []string `json:"urls"`
}

func GetCinema(city City, sort string) ([]*Cinema, error) {
	if sort == "" {
		sort = "-id"
	}

	u := new(url.URL)
	*u = *apiHostname
	switch city {
	case CityMoscow:
		u.Host = "msk." + u.Host
	case CitySpb:
		u.Host = "spb." + u.Host
	default:
		return nil, errors.New("unknown city")
	}
	u.RawQuery = "sort=" + sort
	u.Path = "/cinemas.json"

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}

	data := make([]*Cinema, 0)

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	return data, nil
}

type Screening struct {
	CinemaID   int       `json:"cinema_id"`
	DateTime   time.Time `json:"date_time"`
	ID         int       `json:"id"`
	MovieID    int       `json:"movie_id"`
	PriceMax   int       `json:"price_max"`
	PriceMin   int       `json:"price_min"`
	TicketsUrl string    `json:"tickets_url"`
}

func GetScreeningsByMovie(city City, movieID int) ([]*Screening, error) {
	if movieID == 0 {
		return nil, errors.New("movie id is required")
	}

	u := new(url.URL)
	*u = *apiHostname
	switch city {
	case CityMoscow:
		u.Host = "msk." + u.Host
	case CitySpb:
		u.Host = "spb." + u.Host
	default:
		return nil, errors.New("unknown city")
	}
	u.Path = "/movies/screenings/" + strconv.Itoa(movieID) + ".json"

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}

	data := make([]*Screening, 0)

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	return data, nil
}
