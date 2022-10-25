// Toggle all checkboxes using the "toggle-all" checkbox on top of the table
var selectAll = document.querySelector("input[name=sensors-all]")
var checkboxes = document.querySelectorAll("input[name=sensor-checkbox]")

selectAll.addEventListener("change", () => {
    checkboxes.forEach(sc => {
        sc.checked = selectAll.checked
        sc.dispatchEvent(new Event("change"))
    })
})

// Enable or disable the table control buttons (any button with `.sensor-table-button`)
// if no checkboxes are ticked
checkboxes.forEach(sc => {
    sc.addEventListener("change", () => {
        var checked = Array.from(checkboxes).filter(c => c.checked).length
        var buttons = document.querySelectorAll(".sensor-table-button")
        if (checked > 0) {
            buttons.forEach(b => b.disabled = false)
        } else {
            buttons.forEach(b => b.disabled = true)
        }

        // If all checkboxes are ticked, also tick the "select all"
        selectAll.checked = (checked === checkboxes.length)
    })
})

// Prevent campaign form from redirecting, use REST instead
var form = document.getElementById("campaign-form")
form.addEventListener("submit", event => {
    // Hide the modal if all is well (required fields will be handled
    // by the browser)
    bootstrap.Modal.getInstance("#modal-campaign").hide()

    var data = Object.fromEntries(new FormData(event.target))

    // Get selected/checked sensors from table
    data.sensors = []
    checkboxes.forEach(cb => {
        if (cb.checked) data.sensors.push(cb.value)
    })
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