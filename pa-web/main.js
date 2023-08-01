import Map from 'ol/Map.js';
import TileLayer from 'ol/layer/Tile.js';
import View from 'ol/View.js';
import { OSM, Vector as sourceVector } from 'ol/source.js';
import { Vector as layerVector } from 'ol/layer.js';
import { fromLonLat, toLonLat } from 'ol/proj.js';
import { Point } from 'ol/geom';;
import Feature from 'ol/Feature';
import Icon from 'ol/style/Icon';
import Style from 'ol/style/Style';
import axios from 'axios';
import {v4 as uuidv4} from 'uuid';


const markerSource = new sourceVector();




let socket;
axios.post('http://localhost:8000/google/oauth2/token', {
    "token": window.localStorage.getItem("googleToken"),
    "client_type":"web"
}).then(response => {
    const paToken = response.data.token;
    socket = new WebSocket(`ws://localhost:8000/socket/connect?Authorization=${paToken}`);
    socket.onmessage = function (e) {
        //console.log(e.data);

        var message = JSON.parse(e.data);
        var matchMapping = mapping[message.operation];

        if (matchMapping != null && matchMapping != undefined) {
            matchMapping(message);
        }

    };
})



var mapping = {}
mapping["get_locations_nearby"] = function (message) {

    var fruits = document.querySelector('#locations');

    var i, L = fruits.options.length - 1;
   for(i = L; i >= 0; i--) {
    fruits.remove(i);
   }
    

    message.data.locations.forEach(function (item, index) {
    var orange = document.createElement('option');
    orange.text = `lat : ${item.latitude}, lon: ${item.longitude}, distance_to : ${item.distance_to}` ;
    orange.value = item.id;
    fruits.add(orange);

        addMarker(item.longitude, item.latitude);

    });

};

mapping["get_reserved_park_locations"] = function (message) {

    var fruits = document.querySelector('#locations');


    var i, L = fruits.options.length - 1;
   for(i = L; i >= 0; i--) {
    fruits.remove(i);
   }

    message.data.locations.forEach(function (item, index) {
    var orange = document.createElement('option');
    orange.text = `lat : ${item.latitude}, lon: ${item.longitude}` ;
    orange.value = item.id;
    fruits.add(orange);

        addMarker(item.longitude, item.latitude);

    });

};







var markerStyle = new Style({
    image: new Icon(({
        scale: 0.3,
        src: 'https://raw.githubusercontent.com/openlayers/openlayers/v3.20.1/examples/resources/logo-70x70.png',
        anchor: [0.5, 46],
        anchorXUnits: 'fraction',
        anchorYUnits: 'pixels',
        opacity: 0.75,
    }))
});

const map = new Map({
    layers: [
        new TileLayer({
            source: new OSM(),
        }),
        new layerVector({
            source: markerSource,
            style: markerStyle,
        })
    ],
    target: 'map',
    view: new View({
        center: fromLonLat([29.0202249, 41.0906537]),
        zoom: 15,
    }),
});

var markerGeometry = new Point(fromLonLat([29.0202249, 41.0906537]));
var markerFeature = new Feature({
    geometry: markerGeometry
});

var markerStyle = new Icon(({
    src: 'https://github.com/openlayers/openlayers/blob/v3.20.1/examples/resources/logo-70x70.png'
}));

markerFeature.setStyle(new Style({
    image: markerStyle,
}));

var vectorSource = new sourceVector({
    features: [markerFeature]
});

var markerLayer = new layerVector({
    title: "RoutePoint",
    visible: true,
    source: vectorSource
});

//map.layers.addLayer(markerLayer);


function addMarker(lon, lat) {
    console.log('lon:', lon);
    console.log('lat:', lat);

    var iconFeatures = [];

    var iconFeature = new Feature({
        geometry: new Point(fromLonLat([lon, lat])),
        name: 'Null Island',
        population: 4000,
        rainfall: 500
    });

    markerSource.addFeature(iconFeature);
    let trxId = uuidv4();
    // var data = {"operation": "create_park_location", "transaction_id": trxId, "data": { "longitude": lon, "latitude": lat } };
    // socket.send(JSON.stringify(data));

}
var lastLon, lastLat
map.on('singleclick', function (event) {
    var lonLat = toLonLat(event.coordinate);
    lastLon = lonLat[0];
    lastLat = lonLat[1]
    addMarker(lonLat[0], lonLat[1]);
});

map.on('dblclick', function (event) {
    var lonLat = toLonLat(event.coordinate);
    //addMarker(lonLat[0], lonLat[1]);

    var lon = lonLat[0];
    var lat = lonLat[1];

    let trxId = uuidv4();
    var data = {"operation": "get_locations_nearby", "transaction_id": trxId, "timeout":5, "data": { "longitude": lon, "latitude": lat, "distance": 5000,"count":10,"location_types":[0,1,2],"vehicle_types":[1,2,3] } }

    socket.send(JSON.stringify(data));

});

function parseJwt(token) {
    var base64Url = token.split('.')[1];
    var base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    return JSON.parse(window.atob(base64));
}



var btnCreateParkLocation = document.querySelector('#btnCreateParkLocation');

btnCreateParkLocation.onclick = function(){

    let trxId = uuidv4();
    var data = {"operation": "create_park_location", "transaction_id": trxId,"timeout":5, "data": { "longitude": lastLon, "latitude": lastLat } };
    socket.send(JSON.stringify(data));
};

var btnReserveParkLocation = document.querySelector('#btnReserveParkLocation');

btnReserveParkLocation.onclick = function(){

    let trxId = uuidv4();

    var data = {"operation":"reserve_park_location","transaction_id":trxId,"timeout":5,"data":{"id":document.querySelector('#locations').value}};
    socket.send(JSON.stringify(data));
};


var btnScheduleParkLocationAvailability = document.querySelector('#btnScheduleParkLocationAvailability');

btnScheduleParkLocationAvailability.onclick = function(){

    let trxId = uuidv4();    
    //var data = {"operation":"schedule_park_location_availability","transaction_id":trxId,"data":{"longitude": lon, "latitude": lat,"schedule_type":0,"scheduled_time":Date.now()+3600}};

    var data = {"operation":"schedule_park_location_availability","transaction_id":trxId,"timeout":5,"data":{"id":document.querySelector('#locations').value,"schedule_type":1,"scheduled_time":Date.now()+3600}};

    socket.send(JSON.stringify(data));
};

var btnGetReservedParkLocation = document.querySelector('#btnGetReservedParkLocation');

btnGetReservedParkLocation.onclick = function(){

    let trxId = uuidv4();    
    //var data = {"operation":"schedule_park_location_availability","transaction_id":trxId,"data":{"longitude": lon, "latitude": lat,"schedule_type":0,"scheduled_time":Date.now()+3600}};

    var data = {"operation":"get_reserved_park_locations","transaction_id":trxId,"timeout":5,"data":null};

    socket.send(JSON.stringify(data));
};

var btnRegister = document.querySelector('#btnRegister');

btnRegister.onclick = function(){

    axios.post('http://localhost:8000/google/oauth2/register', {
    "token": window.localStorage.getItem("googleToken"),
    "client_type":"web"
}).then(response => {
    const paToken = response.data.token;
    socket = new WebSocket(`ws://localhost:8000/socket/connect?Authorization=${paToken}`);
    socket.onmessage = function (e) {
        //console.log(e.data);

        var message = JSON.parse(e.data);
        var matchMapping = mapping[message.operation];

        if (matchMapping != null && matchMapping != undefined) {
            matchMapping(message);
        }

    };
})
};