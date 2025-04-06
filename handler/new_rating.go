package handler

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"
)

const (
	MAX = 3300.0
)

type RatingInput struct {
	PlayerRating   *float64 `json:"player_rating,omitempty"`
	OpponentRating *float64 `json:"opponent_rating,omitempty"`
	Result         *float64 `json:"result,omitempty"`
}

type RatingResponse struct {
	NewRating      float64 `json:"new_rating"`
	GoRChange      float64 `json:"gor_change"`
	ExpectedResult float64 `json:"expected_result"`
	Con            float64 `json:"con"`
	Bonus          float64 `json:"bonus"`
	Beta           float64 `json:"beta"`
}

func RouteNewRating(w http.ResponseWriter, r *http.Request) {

	var ratingInput RatingInput

	err := json.NewDecoder(r.Body).Decode(&ratingInput)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	var missingFields []string

	if ratingInput.PlayerRating == nil {
		missingFields = append(missingFields, "player_rating")
	}
	if ratingInput.OpponentRating == nil {
		missingFields = append(missingFields, "opponent_rating")
	}
	if ratingInput.Result == nil {
		missingFields = append(missingFields, "result")
	}

	if len(missingFields) > 0 {
		errMsg := fmt.Sprintf("missing field(s):\n%s", strings.Join(missingFields, ", "))
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	// https://europeangodatabase.eu/EGD/EGF_rating_system.php

	betaPlayer := Beta(*ratingInput.PlayerRating)
	betaOpponent := Beta(*ratingInput.OpponentRating)

	expectedResult := 1.0 / (1.0 + math.Exp(betaOpponent-betaPlayer))

	con := math.Pow((MAX-*ratingInput.PlayerRating)/200, 1.6)
	bonus := math.Log(1+math.Exp((2300-*ratingInput.PlayerRating)/80)) / 5

	newRating := *ratingInput.PlayerRating + con*(*ratingInput.Result-expectedResult) + bonus

	gorChange := newRating - *ratingInput.PlayerRating

	ratingResponse := RatingResponse{
		NewRating:      newRating,
		GoRChange:      gorChange,
		ExpectedResult: expectedResult,
		Con:            con,
		Bonus:          bonus,
		Beta:           betaPlayer,
	}

	response, err := json.Marshal(ratingResponse)
	if err != nil {
		http.Error(w, "could not convert response to json.\nPLEASE NOTIFY API ADMIN\n", http.StatusInternalServerError)
		return
	}

	w.Write(response)
}

func Beta(rating float64) float64 {
	return -7.0 * math.Log(MAX-rating)
}
