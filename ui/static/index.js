// Toggle all checkboxes using the "toggle-all" checkbox on top of the table
var selectAll = document.querySelector("input[name=sensors-all]")
selectAll.addEventListener("change", () => {
    document.querySelectorAll(".sensor-checkbox").forEach(sc => {
        sc.checked = selectAll.checked
        sc.dispatchEvent(new Event("change"))
    })
})

// Enable or disable the table control buttons (any button with `.sensor-table-button`)
// if no checkboxes are ticked
var checkboxes = document.querySelectorAll(".sensor-checkbox")
checkboxes.forEach(sc => {
    sc.addEventListener("change", () => {
        var checked = Array.from(checkboxes).filter(c => c.checked).length
        var buttons = document.querySelectorAll(".sensor-table-button")
        if (checked > 0) {
            buttons.forEach(b => b.disabled = false)
        } else {
            buttons.forEach(b => b.disabled = true)
        }
    })
})

function formSubmit(event) {
    fetch(event.target.action, {
        method: event.target.method,
        body: new FormData(event.target),
    })
    event.preventDefault()
}

// Prevent campaign form from redirecting, use REST instead
document.querySelector("#campaign-form").addEventListener("submit", event => {
    var data = Object.fromEntries(new FormData(event.target))
    console.log(data)

    event.preventDefault()
})

document.querySelector("input[value=raw]").addEventListener("change", () => {
    document.querySelectorAll("input.raw-disable").forEach(i => {
        i.disabled = true
        i.required = false
    })
})

document.querySelector("input[value=aggregated]").addEventListener("change", () => {
    document.querySelectorAll("input.raw-disable").forEach(i => {
        i.disabled = false
        i.required = true
    })
})