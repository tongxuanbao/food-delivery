import {readFile, writeFile} from "fs";

const toString = (coord) => `${coord[0]}_${coord[1]}`;

readFile('restaurants.geojson', function(err, data) { 
    const geojson = JSON.parse(data); 

	const list = geojson.elements.map(element => ([element.lon, element.lat]));
	

    writeFile('restaurant.json', JSON.stringify(list), () => {})
}); 