package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type TestingData struct {
	Input  RatingInput    `json:"input"`
	Output RatingResponse `json:"output"`
}

func convertToBodyReader(rating RatingInput) *bytes.Reader {
	jsonData, err := json.Marshal(rating)
	if err != nil {
		panic("Invalid json")
	}

	return bytes.NewReader(jsonData)
}

func assertCorrectMessage(t *testing.T, got interface{}, want interface{}) {
	t.Helper()
	if got != want {
		t.Errorf("\nGOT:  %s\nWANT: %s\n", got, want)
	}
}

func TestPOSTNewRating(t *testing.T) {
	t.Run("empty body", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/new_rating", nil)
		response := httptest.NewRecorder()

		RouteNewRating(response, request)

		got := response.Body.String()
		want := "body is empty\n"
		assertCorrectMessage(t, got, want)
	})

	t.Run("all empty fields", func(t *testing.T) {

		rating := RatingInput{
			PlayerRating:   nil,
			OpponentRating: nil,
			Result:         nil,
		}

		body := convertToBodyReader(rating)

		request, _ := http.NewRequest(http.MethodPost, "/new_rating", body)
		response := httptest.NewRecorder()

		RouteNewRating(response, request)

		got := response.Body.String()
		want := "missing field(s):\nplayer_rating, opponent_rating, result\n"
		assertCorrectMessage(t, got, want)
	})
	t.Run("missing player_rating field", func(t *testing.T) {

		opponentRating := 2300.0
		result := 1.0

		rating := RatingInput{
			PlayerRating:   nil,
			OpponentRating: &opponentRating,
			Result:         &result,
		}

		body := convertToBodyReader(rating)

		request, _ := http.NewRequest(http.MethodPost, "/new_rating", body)
		response := httptest.NewRecorder()

		RouteNewRating(response, request)

		got := response.Body.String()
		want := "missing field(s):\nplayer_rating\n"
		assertCorrectMessage(t, got, want)
	})
	t.Run("missing opponent_rating field", func(t *testing.T) {

		playerRating := 2300.0
		result := 1.0

		rating := RatingInput{
			PlayerRating:   &playerRating,
			OpponentRating: nil,
			Result:         &result,
		}

		body := convertToBodyReader(rating)

		request, _ := http.NewRequest(http.MethodPost, "/new_rating", body)
		response := httptest.NewRecorder()

		RouteNewRating(response, request)

		got := response.Body.String()
		want := "missing field(s):\nopponent_rating\n"
		assertCorrectMessage(t, got, want)
	})
	t.Run("missing rating field", func(t *testing.T) {

		playerRating := 2300.0
		opponentRating := 2002.0

		rating := RatingInput{
			PlayerRating:   &playerRating,
			OpponentRating: &opponentRating,
			Result:         nil,
		}

		body := convertToBodyReader(rating)

		request, _ := http.NewRequest(http.MethodPost, "/new_rating", body)
		response := httptest.NewRecorder()

		RouteNewRating(response, request)

		got := response.Body.String()
		want := "missing field(s):\nresult\n"
		assertCorrectMessage(t, got, want)
	})

	t.Run("invalid rating", func(t *testing.T) {

		playerRating := 2300.0
		opponentRating := 2002.0
		result := 10.0

		rating := RatingInput{
			PlayerRating:   &playerRating,
			OpponentRating: &opponentRating,
			Result:         &result,
		}

		body := convertToBodyReader(rating)

		request, _ := http.NewRequest(http.MethodPost, "/new_rating", body)
		response := httptest.NewRecorder()

		RouteNewRating(response, request)

		got := response.Body.String()
		want := "result must be between 0.0 and 1.0\n"
		assertCorrectMessage(t, got, want)
	})

	t.Run("valid examples", func(t *testing.T) {

		jsonData, err := os.ReadFile("../test_data.json")
		if err != nil {
			t.Error(err)
			return
		}

		var testingData []TestingData
		err = json.Unmarshal(jsonData, &testingData)
		if err != nil {
			t.Error(err)
			return
		}

		for _, testData := range testingData {
			body := convertToBodyReader(testData.Input)

			request, _ := http.NewRequest(http.MethodPost, "/new_rating", body)
			response := httptest.NewRecorder()

			RouteNewRating(response, request)

			var got RatingResponse
			err := json.NewDecoder(response.Body).Decode(&got)
			if err != nil {
				t.Errorf("could not decode")
			}
			want := testData.Output
			if got != want {
				t.Errorf("\nGOT:  %v\nWANT: %v\n", got, want)
			}
		}
	})
}
