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

const markerSource = new sourceVector();

const socket = new WebSocket("ws://localhost:8000/socket/connect");

var mapping = {}
mapping["get_locations_nearby"] = function(message){

    

   message.data.forEach(function(item,index){
    addMarker(item.longitude,item.latitude);

   });

};

socket.onmessage = function (e) {
    //console.log(e.data);

    var message =JSON.parse(e.data);
    var matchMapping = mapping[message.operation];

    if (matchMapping != null && matchMapping != undefined){
        matchMapping(message);
    }
   
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

    //var data = {"client_id":"1234","operation":"create_park_location","transaction_id":"123e4567-e89b-12d3-a456-426655440000","data":{"longitude":lon,"latitude":lat}};
    //socket.send(JSON.stringify(data));

}

map.on('singleclick', function (event) {
    var lonLat = toLonLat(event.coordinate);
    addMarker(lonLat[0], lonLat[1]);
});

map.on('dblclick', function (event) {
    var lonLat = toLonLat(event.coordinate);
    //addMarker(lonLat[0], lonLat[1]);

    var lon = lonLat[0];
    var lat = lonLat[1];


    var data ={"client_id":"1234","operation":"get_locations_nearby","transaction_id":"123e4567-e89b-12d3-a456-426655440000","data":{"longitude":lon,"latitude":lat,"distance":5000}}

    socket.send(JSON.stringify(data));

});