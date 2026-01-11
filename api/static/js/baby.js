async function handleFauxKey(value, guess) {
    if (guess_count >= 6) {
        return guess
    }
    if (value === 'Enter') {
        return await handleSubmit(guess)
    }
    if (value === 'Backspace') {
        if (guess.length === 0) return guess
        const used = guess_grid.querySelectorAll('div.container.used')
        const last_used = used[used.length - 1]
        last_used.classList.remove('used')
        const front = last_used.querySelector('.front')
        const back = last_used.querySelector('.back')
        front.innerHTML = ''
        back.innerHTML = ''
        guess = guess.substring(0, guess.length-1)
        return guess
    }
    if (guess.length === 5) return guess
    const first_free = guess_grid.querySelector('div.container:not(.used, .flipped)')
    const front = first_free.querySelector('.front')
    const back = first_free.querySelector('.back')
    front.innerText = value
    back.innerText = value
    first_free.classList.add('used')
    guess = guess + value
    return guess
}

function storeLocalStorageState(response) {
    const existing_storage = localStorage.getItem('baby_state')
    let json_store = []
    if (existing_storage) {
        json_store = JSON.parse(existing_storage)
    }
    json_store.push(response)
    string_store = JSON.stringify(json_store)
    localStorage.setItem('baby_state', string_store)
}

function displaySuccess(html, delay='1.1s') {
    successWrapper.innerHTML = html
    const success = document.getElementById('success')
    const svg = success.querySelector('svg')
    svg.style.visibility = 'visible'
    success.style.animation = `fade 3s ease ${delay} 1 normal forwards running`
    drawWave('wavePath')
    drawWave('wavePathTwo')
}

function loadFromLocalStorage() {
    const baby_storage = localStorage.getItem('baby_state')
    if (baby_storage) {
        const json_baby = JSON.parse(baby_storage)
        const div_cells = guess_grid.querySelectorAll('div.container')
        let row_index = 0
        let column_index = 0
        guess_count = json_baby.length
        for (let i = 0; i < div_cells.length; i++) {
            const cell = div_cells[i]
            if (column_index > 4) {
                column_index = 0
                row_index++
            }
            if (row_index > json_baby.length - 1) break
            const cell_data = json_baby[row_index].Characters[column_index]
            const back_el = cell.querySelector('.back')
            back_el.innerText = cell_data.Character
            cell.classList.add(cell_data.Color)
            cell.classList.add('flipped')
            column_index++
        }
        if (json_baby[json_baby.length - 1].Correct) {
            const html = json_baby[json_baby.length - 1].Success_Html
            displaySuccess(html)
        }
        if (json_baby[json_baby.length - 1].Failed) {
            const html = json_baby[json_baby.length - 1].Failure_Html
            completeAndUtterFailure.innerHTML = html
            const script = document.createElement('script')
            script.textContent = json_baby[json_baby.length - 1].Script
            document.body.appendChild(script)
        }
    }
}

function restart() {
    completeAndUtterFailure.innerHTML = ''
    const cells = document.querySelectorAll('.guess-row .container')
    cells.forEach(cell => {
        cell.classList.remove('flipped')
        cell.classList.remove('grey')
        cell.classList.remove('green')
        cell.classList.remove('yellow')
        const front = cell.querySelector('.front')
        const back = cell.querySelector('.back')
        front.innerHTML = ''
        back.innerHTML = ''
    })
    localStorage.removeItem('baby_state')
    guess_count = 0
    //document.getElementById('hi').addEventListener('click', async () => await surrender())
}

async function handleSubmit(guess) {
    try {
        guess_count++
        const result = await fetch('/guess', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ guess: guess, count: guess_count })
        })
        if (result.status === 400) {
            guess_count--
            return guess
        }
        const json_response = await result.json()
        storeLocalStorageState(json_response)
        const submitted = guess_grid.querySelectorAll('div.used')
        for (let i = 0; i < submitted.length; i++) {
            const square = submitted[i]
            const color = json_response.Characters[i].Color
            square.classList.add(color)
            square.classList.add('flipped')
            square.classList.remove('used')
        }
        if (json_response.Correct === true) {
            const html = json_response.Success_Html
            displaySuccess(html)
        }
        if (json_response.Failed === true) {
            const html = json_response.Failure_Html
            completeAndUtterFailure.innerHTML = html
            const script = document.createElement('script')
            script.textContent = json_response.Script
            document.body.appendChild(script)
        }
        return ''
    } catch (err) {
        console.log(err)
    }
}

const isAlpha = (char) => /^[a-zA-Z]$/.test(char);
async function handleKeyPress(e, guess) {
    if (e.ctrlKey || e.metaKey) return guess
    if (e.key === 'Backspace') {
        return await handleFauxKey(e.key, guess)
    }
    if (e.key === 'Enter') {
        return await handleSubmit(guess)
    }
    if (!isAlpha(e.key)) return guess
    return await handleFauxKey(e.key, guess)
}

function drawWave(id) {
    const path = document.getElementById(id);

    const centerX = 100;
    const centerY = 100;
    const baseRadius = 50;
    const amplitude = 10;
    const frequency = 7; // Number of wave cycles
    const points = 180;

    let pathData = '';

    for (let i = 0; i <= points; i++) {
        const angle = (i / points) * Math.PI * 2;
        const sineValue = Math.sin(angle * frequency);
        const radius = baseRadius + (sineValue * amplitude);

        const x = centerX + radius * Math.cos(angle);
        const y = centerY + radius * Math.sin(angle);

        if (i === 0) {
            pathData += `M ${x} ${y}`;
        } else {
            pathData += ` L ${x} ${y}`;
        }
    }
    path.setAttribute('d', pathData);
    const pathLength = path.getTotalLength()
    console.log(pathLength)
    path.style.strokeDasharray = pathLength
    path.style.strokeDashoffset = pathLength
}

let guess_count = 0
window.addEventListener('load', () => {
    loadFromLocalStorage()
    const faux_keys = document.querySelectorAll('.faux-key')
    let guess = ''
    faux_keys.forEach(key => {
        key.addEventListener('click', async (e) => guess = await handleFauxKey(e.target.getAttribute('data-value'), guess))
    })

    document.addEventListener('keydown', async (e) => guess = await handleKeyPress(e, guess))
})
