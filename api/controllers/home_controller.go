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
	Count int
}

type Char_Responses struct {
	Characters   []Response
	Correct      bool
	Success_Html string
	Failed       bool
	Failure_Html string
	Script       string
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

func SuccessHtml(delay string) string {
	success_html := fmt.Sprintf(`
		<div id='success'>
			<div class='wrapper'>
				<p>AIDAN</p>
				<p>OSTAP</p>
				<p>MCLEOD</p>
				<svg class='vine' viewBox='0 0 400 400' xmlns='http://www.w3.org/2000/svg'>
					<defs>
						<linearGradient id='brownGradient' x1='0%%' y1='0%%' x2='100%%' y2='100%%'>
							<stop offset='0%%' style='stop-color:#dab12b;stop-opacity:1' />
							<stop offset='50%%' style='stop-color:#d9701a;stop-opacity:1' />
							<stop offset='100%%' style='stop-color:#150d18;stop-opacity:1' />
						</linearGradient>
					</defs>
					<path id='wavePath' fill='none' stroke='url(#brownGradient)' stroke-width='3' stroke-linecap='round'>
						<animate attributeName='stroke-dashoffset' from='440' to='0' begin='%s' dur='5s' fill='freeze'/>
					</path>
					<path id='wavePathTwo' fill='none' stroke='url(#brownGradient)' stroke-width='3' stroke-linecap='round'>
						<animate attributeName='stroke-dashoffset' from='400' to='0' dur='5s' begin='%s' fill='freeze'/>
					</path>
				</svg>
			</div>
		</div>
	`, delay, delay)
	return success_html
}

func Surrender(res http.ResponseWriter, req *http.Request) {
	char_responses := Char_Responses{}
	char_responses.Correct = true
	char_responses.Failed = false
	char_responses.Success_Html = SuccessHtml("0s")
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	err := json.NewEncoder(res).Encode(char_responses)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
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
	char_responses := Char_Responses{}
	response_data := []Response{}
	green_chars := 0
	for i, char := range guess_chars {
		if char == answer[i] {
			char_data := Response{
				Character: char,
				Color:     "green",
			}
			chars_found = append(chars_found, char)
			response_data = append(response_data, char_data)
			green_chars += 1
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
	if green_chars == len(answer) {
		char_responses.Correct = true
		char_responses.Failed = false
		char_responses.Success_Html = SuccessHtml("1.1s")
	} else if guess.Count == 6 {
		char_responses.Failed = true
		char_responses.Failure_Html = `
			<p>nice try!</p>
			<div class="actions">
				<button id="cowardly-surrender" class="surrender">reveal name</button>
				<button class="restart" onclick="restart()">restart</button>
			</div>
		`
		char_responses.Script = `
			async function surrender() {
				try {
					const stored_count = guess_count
					guess_count = 6
					const result = await fetch('/cowardly-surrender', {
						method: 'GET',
						headers: {
							'Content-Type': 'application/json'
						}
					})
					if (result.status === 400) {
						guess_count = stored_count
					}
					const json_response = await result.json()
					const html = json_response.Success_Html
					displaySuccess(html, '0s')
				} catch (err) {
					console.log(err)
				}
			}
			document.getElementById('cowardly-surrender').addEventListener('click', async () => await surrender())
		`
	}
	char_responses.Characters = response_data

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	err = json.NewEncoder(res).Encode(char_responses)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
