package main

import (
	"API/internal/storage"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type PatchExerciseRequest struct {
	Name   *string  `json:"name"`
	Repeat *int     `json:"reps"`
	Set    *int     `json:"sets"`
	Weight *float32 `json:"weight"`
}

func main() {

	mux := http.NewServeMux()
	Exercises := make([]storage.Exercise, 0, 10)
	countID := 3

	Exercises = append(Exercises, storage.Exercise{
		ID:     1,
		Name:   "Присед",
		Repeat: 10,
		Set:    3,
		Weight: 120.0,
	}, storage.Exercise{
		ID:     2,
		Name:   "Жим лежа",
		Repeat: 6,
		Set:    4,
		Weight: 100.0,
	}, storage.Exercise{
		ID:     3,
		Name:   "Становая тяга",
		Repeat: 6,
		Set:    3,
		Weight: 160.0,
	})

	AddExercise := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		exercise := storage.Exercise{}
		w.Header().Set("Content-Type", "application/json")

		err := json.NewDecoder(r.Body).Decode(&exercise)
		if err != nil {
			http.Error(w, "Error decoding JSON", http.StatusBadRequest)
			return
		}

		if strings.TrimSpace(exercise.Name) == "" || exercise.Repeat <= 0 || exercise.Set <= 0 || exercise.Weight < 0 {
			http.Error(w, "Unknow data", http.StatusBadRequest)
			return
		}

		if len(Exercises) == 0 {
			exercise.ID = 1
		} else {
			countID++
			exercise.ID = countID
		}

		err = json.NewEncoder(w).Encode(&exercise)
		if err != nil {
			http.Error(w, "Error decoding JSON", http.StatusBadRequest)
			return
		}
		Exercises = append(Exercises, exercise)
	}

	GetExercises := func(w http.ResponseWriter, r *http.Request) {
		// for _, exercise := range Exercises {
		// 	fmt.Fprintf(w, "%d. %s %d * %d - %.2f кг\n",
		// 		exercise.ID, exercise.Name, exercise.Repeat, exercise.Set, exercise.Weight)
		// }
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(Exercises)
		if err != nil {
			http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
			return
		}
	}

	DelExercises := func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || id <= 0 {
			http.Error(w, "Error uncorrect id", http.StatusBadRequest)
			return
		}

		for i, ex := range Exercises {
			if id == ex.ID {
				Exercises = append(Exercises[:i], Exercises[i+1:]...)
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
		http.Error(w, "Exercise Not Found", http.StatusNotFound)

	}

	GetExerciseById := func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || id <= 0 {
			http.Error(w, "Error uncorrect id", http.StatusBadRequest)
			return
		}

		for _, ex := range Exercises {
			if id == ex.ID {
				w.Header().Set("Content-Type", "application/json")
				err := json.NewEncoder(w).Encode(ex)
				if err != nil {
					http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
					return
				}
				return
			}

		}
		http.Error(w, "Exercise Not Found", http.StatusNotFound)

	}

	UpdateExercise := func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || id <= 0 {
			http.Error(w, "Error uncorrect id", http.StatusBadRequest)
			return
		}
		var exercise storage.Exercise

		err = json.NewDecoder(r.Body).Decode(&exercise)
		if err != nil {
			http.Error(w, "Json Decode Error", http.StatusBadRequest)
			return
		}
		if strings.TrimSpace(exercise.Name) == "" || exercise.Repeat <= 0 || exercise.Set <= 0 || exercise.Weight < 0 {
			http.Error(w, "Unknow data", http.StatusBadRequest)
			return
		}

		for i, ex := range Exercises {
			if ex.ID == id {
				Exercises[i] = exercise
				Exercises[i].ID = id

				w.Header().Set("Content-Type", "application/json")
				err = json.NewEncoder(w).Encode(Exercises[i])
				if err != nil {
					http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
					return
				}
				return
			}
		}
		http.Error(w, "Exercise Not Found", http.StatusNotFound)

	}

	PatchExersise := func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || id <= 0 {
			http.Error(w, "Error uncorrect ID", http.StatusBadRequest)
			return
		}
		var exercise PatchExerciseRequest

		err = json.NewDecoder(r.Body).Decode(&exercise)
		if err != nil {
			http.Error(w, "Json Decode Error", http.StatusBadRequest)
			return
		}

		for i, ex := range Exercises {
			if ex.ID == id {
				Exercises[i].ID = id

				if exercise.Weight != nil && *exercise.Weight < 0 {
					http.Error(w, "Unknow data", http.StatusBadRequest)
					return
				}
				if exercise.Name != nil && strings.TrimSpace(*exercise.Name) == "" {
					http.Error(w, "Unknow data", http.StatusBadRequest)
					return
				}
				if exercise.Set != nil && *exercise.Set <= 0 {
					http.Error(w, "Unknow data", http.StatusBadRequest)
					return
				}
				if exercise.Repeat != nil && *exercise.Repeat <= 0 {
					http.Error(w, "Unknow data", http.StatusBadRequest)
					return
				}

				if exercise.Repeat != nil {
					Exercises[i].Repeat = *exercise.Repeat
				}
				if exercise.Set != nil {
					Exercises[i].Set = *exercise.Set
				}
				if exercise.Weight != nil {
					Exercises[i].Weight = *exercise.Weight
				}
				if exercise.Name != nil {
					Exercises[i].Name = *exercise.Name
				}

				w.Header().Set("Content-Type", "application/json")
				err = json.NewEncoder(w).Encode(Exercises[i])
				if err != nil {
					http.Error(w, "Json Decode Error", http.StatusBadRequest)
					return
				}

				return
			}
		}
		http.Error(w, "Exercise Not Found", http.StatusNotFound)
	}

	mux.HandleFunc("PATCH /exercises/{id}", PatchExersise)
	mux.HandleFunc("PUT /exercises/{id}", UpdateExercise)
	mux.HandleFunc("GET /exercises/{id}", GetExerciseById)
	mux.HandleFunc("DELETE /exercises/{id}", DelExercises)
	mux.HandleFunc("GET /exercises", GetExercises)
	mux.HandleFunc("POST /exercises", AddExercise)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
	log.Println("Server started on port 8080")
}
