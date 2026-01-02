package controllers

import (
	"encoding/json"
	"fmt"
	"gotemplate/templates"
	"net/http"
	"os"
	"slices"
	"strings"
)

func GetHome(res http.ResponseWriter, req *http.Request) {
	component := templates.HomePage()
	cssFiles := []string{""}
	jsFiles := []string{"baby.js"}
	page := templates.Html(component, cssFiles, jsFiles)
	page.Render(req.Context(), res)
}

type Guess struct {
	Guess string
}

type Response struct {
	Character string
	Color     string
}

func CountOccurences(chars []string, char string) int {
	count := 0
	for _, n := range chars {
		if n == char {
			count++
		}
	}
	return count
}

func ProcessGuess(res http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var guess Guess
	err := decoder.Decode(&guess)
	if err != nil {
		fmt.Println("Could not decode guess")
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	answer_string := os.Getenv("ANSWER")
	answer := strings.Split(answer_string, "")
	if len(guess.Guess) != len(answer) {
		fmt.Println("Guess incorrect length")
		http.Error(res, "incorrect guess length", http.StatusBadRequest)
		return
	}
	guess_chars := strings.Split(guess.Guess, "")
	chars_found := []string{}
	response_data := []Response{}
	for i, char := range guess_chars {
		if char == answer[i] {
			char_data := Response{
				Character: char,
				Color:     "green",
			}
			chars_found = append(chars_found, char)
			response_data = append(response_data, char_data)
			continue
		}
		if slices.Contains(answer, char) {
			if !slices.Contains(chars_found, char) || (slices.Contains(chars_found, char) && CountOccurences(chars_found, char) < CountOccurences(answer, char)) {
				char_data := Response{
					Character: char,
					Color:     "yellow",
				}
				chars_found = append(chars_found, char)
				response_data = append(response_data, char_data)
				continue
			}
		}
		char_data := Response{
			Character: char,
			Color:     "grey",
		}
		response_data = append(response_data, char_data)
	}
	for i, char := range slices.Backward(response_data) {
		if char.Color == "yellow" && CountOccurences(chars_found, char.Character) > CountOccurences(answer, char.Character) {
			char_found_index := slices.Index(chars_found, char.Character)
			chars_found = slices.Delete(chars_found, char_found_index, char_found_index+1)
			response_data[i] = Response{
				Character: char.Character,
				Color:     "grey",
			}
		}
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	err = json.NewEncoder(res).Encode(response_data)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
