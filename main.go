package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Questionnaire struct {
	ID        int
	Name      string
	Questions []Question
}

type Question struct {
	ID              int
	QuestionnaireID int
	Question        string
}

type QuestionnaireResponse struct {
	QuestionnaireID int
	Responses       []QuestionResponse
}

type QuestionResponse struct {
	QuestionID int
	Response   string
}

func main() {
	r := chi.NewRouter()
	questionnaireTpl := template.Must(template.ParseFiles("questionnaire.gohtml"))

	q := lookupQuestionnaire(1)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		questionnaireTpl.Execute(w, q)
	})
	r.Post("/responses", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		var response QuestionnaireResponse
		id, err := strconv.Atoi(r.FormValue("questionnaire_id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		response.QuestionnaireID = id
		questionnaire := lookupQuestionnaire(id)
		for _, question := range questionnaire.Questions {
			response.Responses = append(response.Responses, QuestionResponse{
				QuestionID: question.ID,
				Response:   r.FormValue(fmt.Sprintf("questions[%d]", question.ID)),
			})
		}

		// Just a quick way to render the response that we are parsing.
		enc := json.NewEncoder(w)
		enc.Encode(response)
	})

	http.ListenAndServe(":3000", r)
}

// Using this to fake some stuff
func lookupQuestionnaire(id int) Questionnaire {
	return Questionnaire{
		ID:   id,
		Name: "Test Questionnaire",
		Questions: []Question{
			{
				ID:              1,
				QuestionnaireID: 1,
				Question:        "What is your name?",
			},
			{
				ID:              6,
				QuestionnaireID: 1,
				Question:        "What is your quest?",
			},
			{
				ID:              91,
				QuestionnaireID: 1,
				Question:        "What is your favorite color?",
			},
		},
	}
}
