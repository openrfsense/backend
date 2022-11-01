var data = document.getElementById("map").dataset
var map = L.map("map", {
    center: [data.lat, data.lon], 
    minZoom: 3,
    zoom: 13
})
map.setMaxBounds([[-85.0511, -180], [85.0511, 180]])
L.marker([data.lat, data.lon]).addTo(map)
L.tileLayer(
    "https://tile.openstreetmap.org/{z}/{x}/{y}.png",
    {
        maxZoom: 19,
        noWrap: true,
        attribution: `&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>`
    }
).addTo(map)
