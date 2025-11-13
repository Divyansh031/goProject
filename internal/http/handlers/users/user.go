package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/Divyansh031/goProject/internal/storage"
	"github.com/Divyansh031/goProject/internal/types"
	"github.com/Divyansh031/goProject/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage ) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Creating user")

		var user types.User

		err := json.NewDecoder(r.Body).Decode(&user)
		if errors.Is(err, io.EOF) {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return 
		}

		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return 
		}

		// Request validation


		if err := validator.New().Struct(user); err != nil {

			validateErrs := err.(validator.ValidationErrors)
			response.WriteJSON(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return 	
		}
		
		lastId, err := storage.CreateUser(
			user.Name,
			user.Email,
			user.Age,
		)

		slog.Info("User created successfully", slog.String("UserID", fmt.Sprintf("%d", lastId)))

		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err)
			return 
		}

		response.WriteJSON(w, http.StatusCreated, map[string] int64 {"id": lastId})
	}
}