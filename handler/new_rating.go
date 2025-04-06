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

	/*
		https://europeangodatabase.eu/EGD/EGF_rating_system.php

		# System description
		The rating algorithm was updated starting 2021. The whole database from back in 1996 was recalculated with this algorithm. You can find the old algorithm here.

		Ratings are updated by: r' = r + con * (Sa - Se) + bonus
		r is the old EGD rating (GoR) of the player
		r' is the new EGD rating of the player
		Sa is the actual game result (1.0 = win, 0.5 = jigo, 0.0 = loss)
		Se is the expected game result as a winning probability (1.0 = 100%, 0.5 = 50%, 0.0 = 0%). See further below for its computation.
		con is a factor that determines rating volatility (similar to K in regular Elo rating systems): con = ((3300 - r) / 200)^1.6
		bonus (not found in regular Elo rating systems) is a term included to counter rating deflation: bonus = ln(1 + exp((2300 - rating) / 80)) / 5

		Se is computed by the Bradley-Terry formula: Se = 1 / (1 + exp(β(r2) - β(r1)))
		r1 is the EGD rating of the player
		r2 is the EGD rating of the opponent
		β is a mapping function for EGD ratings: β = -7 * ln(3300 - r)
	*/

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
