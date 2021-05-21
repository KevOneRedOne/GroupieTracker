//Step 1: initialize communication with the platform
var platform = new H.service.Platform({
    apikey: "UxDC5OOIhOC4BEmxOmwUKAa8q413CJS311fmOzt3UeE"
});
var defaultLayers = platform.createDefaultLayers();

//Step 2: initialize a map - this map is centered over Europe
var map = new H.Map(document.getElementById('map'),
    defaultLayers.vector.normal.map, {
    center: { lat: 20, lng: 10 },
    zoom: 0,
    pixelRatio: window.devicePixelRatio || 1
});

//Add a resize listener to make sure that the map occupies the whole container
window.addEventListener('resize', () => map.getViewPort().resize());

//Step 3: make the map interactive
var behavior = new H.mapevents.Behavior(new H.mapevents.MapEvents(map));

// Create the default UI components
var ui = H.ui.UI.createDefault(map, defaultLayers);

// Read locations and add markers
window.onload = function () {
    let locations = document.getElementsByClassName('cities');
    for (let i = 0; i < locations.length; i++) {
        console.log(locations)
        fetch('https://geocoder.ls.hereapi.com/6.2/geocode.json?searchtext=' + locations[i].textContent + '&gen=9&apiKey=chDMQRlDRv0KngLYo3sXcOtNBESgx1eEU199e4Z1B7U')
            .then(response => response.json())
            .then(data => {
                map.addObject(new H.map.Marker({
                    lat: data.Response.View[0].Result[0].Location.DisplayPosition.Latitude,
                    lng: data.Response.View[0].Result[0].Location.DisplayPosition.Longitude
                }));
            })
    }
}