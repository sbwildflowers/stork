async function handleFauxKey(value, guess) {
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

function loadFromLocalStorage() {
    const baby_storage = localStorage.getItem('baby_state')
    if (baby_storage) {
        const json_baby = JSON.parse(baby_storage)
        const div_cells = guess_grid.querySelectorAll('div.container')
        let row_index = 0
        let column_index = 0
        for (let i = 0; i < div_cells.length; i++) {
            const cell = div_cells[i]
            if (column_index > 4) {
                column_index = 0
                row_index++
            }
            if (row_index > json_baby.length - 1) break
            const cell_data = json_baby[row_index][column_index]
            const back_el = cell.querySelector('.back')
            back_el.innerText = cell_data.Character
            cell.classList.add(cell_data.Color)
            cell.classList.add('flipped')
            column_index++
        }
    }
}

async function handleSubmit(guess) {
    try {
        const result = await fetch('/guess', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ guess: guess })
        })
        if (result.status === 400) {
            return guess
        }
        const json_response = await result.json()
        storeLocalStorageState(json_response)
        console.log(json_response)
        const submitted = guess_grid.querySelectorAll('div.used')
        for (let i = 0; i < submitted.length; i++) {
            const square = submitted[i]
            const color = json_response[i].Color
            square.classList.add(color)
            square.classList.add('flipped')
            square.classList.remove('used')
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

window.addEventListener('load', () => {
    loadFromLocalStorage()
    const faux_keys = document.querySelectorAll('.faux-key')
    let guess = ''
    faux_keys.forEach(key => {
        key.addEventListener('click', async (e) => guess = await handleFauxKey(e.target.getAttribute('data-value'), guess))
    })

    document.addEventListener('keydown', async (e) => guess = await handleKeyPress(e, guess))
})
