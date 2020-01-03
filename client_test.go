package subcity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMovies(t *testing.T) {
	_, err := GetMovies(CityMoscow, "")
	assert.Nil(t, err, "got an error. expected nil")
}

func TestGetCinemas(t *testing.T) {
	_, err := GetCinema(CityMoscow, "")
	assert.Nil(t, err, "got an error. expected nil")
}

func TestGetScreenings(t *testing.T) {
	_, err := GetScreeningsByMovie(CityMoscow, 70649)
	assert.Nil(t, err, "got an error. expected nil")
}
