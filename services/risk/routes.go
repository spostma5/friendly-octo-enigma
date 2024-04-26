package risk

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/spostma5/friendly-octo-enigma/utils"
)

func HandleGetRisks(w http.ResponseWriter, r *http.Request) {
	if err := utils.JSONEncode(w, risks, http.StatusOK); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

func HandlePostRisk(w http.ResponseWriter, r *http.Request) {
	var risk Risk

	if err := json.NewDecoder(r.Body).Decode(&risk); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	validate := validator.New()

	if err := validate.Struct(risk); err != nil {
		errors := err.(validator.ValidationErrors)
		http.Error(w, fmt.Sprintf("Validation error: %s", errors), http.StatusBadRequest)
		return
	}

	if err := createRisk(risk); err != nil {
		http.Error(w, fmt.Sprintf("Unable to create risk: %s", err.Error()), http.StatusConflict)
		return
	}
}

func HandleGetRisk(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	risk := getRisk(id)

	if risk == nil {
		http.Error(w, "Unable to find requested risk", http.StatusNotFound)
		return
	}

	if err := utils.JSONEncode(w, risk, http.StatusOK); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}
