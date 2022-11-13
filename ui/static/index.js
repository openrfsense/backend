// Toggle all checkboxes using the "toggle-all" checkbox on top of the table
var selectAll = document.querySelector("input[name=sensors-all]")
var checkboxes = document.querySelectorAll("input[name=sensor-checkbox]")

selectAll.addEventListener("click", () => {
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
    })
})

// Prevent campaign form from redirecting, use REST instead
var form = document.getElementById("campaign-form")
form.addEventListener("submit", event => {
    event.preventDefault()

    var data = Object.fromEntries(new FormData(event.target))
    data.begin = Math.floor(
        new Date(`${data.startDate}T${data.startTime}`).getTime() / 1000
    )
    data.end = Math.floor(
        new Date(`${data.endDate}T${data.endTime}`).getTime() / 1000
    )
    if (data.freqCenter)
        data.freqCenter = parseInt(data.freqCenter, 10)
    if (data.freqMin)
        data.freqMin = parseInt(data.freqMin, 10)
    if (data.freqMax)
        data.freqMax = parseInt(data.freqMax, 10)
    if (data.freqRes)
        data.freqRes = parseInt(data.freqRes, 10)
    if (data.timeRes)
        data.timeRes = parseInt(data.timeRes, 10)
    data.sensors = []
    // Get selected/checked sensors from table
    checkboxes.forEach(cb => {
        if (cb.checked) data.sensors.push(cb.value)
    })

    fetch(
        data.measurementType === "raw" ? "/api/v1/raw" : "/api/v1/aggregated",
        {
            headers: { "Content-Type": "application/json" },
            method: "post",
            body: JSON.stringify(data),
        }
    ).then(response => {
        if (!response.ok) {
            document.getElementById("alert-error").classList.toggle("show", true)
            document.getElementById("alert-success").classList.toggle("show", false)
            return
        }

        // Hide the modal if all is well (required fields will be handled
        // by the browser)
        bootstrap.Modal.getInstance("#modal-campaign").hide()
        document.getElementById("alert-error").classList.toggle("show", false)
        document.getElementById("alert-success").classList.toggle("show", true)
    })
})

// Raw measurement radio toggle in modal
document.querySelector("input[value=raw]").addEventListener("change", () => {
    document.querySelectorAll(".aggregated-vanish").forEach(i => i.style.display = "")
    document.querySelectorAll(".raw-vanish").forEach(i => i.style.display = "none")
    document.querySelectorAll("input.raw-disable").forEach(i => {
        i.disabled = true
        i.required = false
    })
    document.querySelectorAll("input.aggregated-disable").forEach(i => {
        i.disabled = false
        i.required = true
    })
})

// Sampled measurement radio toggle in modal
document.querySelector("input[value=aggregated]").addEventListener("change", () => {
    document.querySelectorAll(".aggregated-vanish").forEach(i => i.style.display = "none")
    document.querySelectorAll(".raw-vanish").forEach(i => i.style.display = "")
    document.querySelectorAll("input.raw-disable").forEach(i => {
        i.disabled = false
        i.required = true
    })
    document.querySelectorAll("input.aggregated-disable").forEach(i => {
        i.disabled = true
        i.required = false
    })
})

// Initialize the Leaflet map
var map = L.map("map", {
    center: [0, 0],
    zoom: 3,
    minZoom: 3
})
map.setMaxBounds([[-85.0511, -180], [85.0511, 180]])
// Add markers to map with simple DOM trick
document.querySelectorAll("marker").forEach(m => {
    var data = m.dataset
    L.marker([data.lat, data.lon]).addTo(map).bindPopup(m.innerHTML)
    m.remove()
})
L.tileLayer(
    "https://tile.openstreetmap.org/{z}/{x}/{y}.png",
    {
        maxZoom: 19,
        noWrap: true,
        attribution: `&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>`
    }
).addTo(map)