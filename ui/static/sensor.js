var data = document.getElementById("map").dataset
var map = L.map("map").setView([data.lat, data.lon], 13)
L.marker([data.lat, data.lon]).addTo(map)
L.tileLayer(
    "https://tile.openstreetmap.org/{z}/{x}/{y}.png",
    {
        maxZoom: 19,
        noWrap: true,
        attribution: `&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>`
    }
).addTo(map)
