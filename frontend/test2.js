import {readFile, writeFile} from "fs";

const toString = (coord) => `${coord[0]}_${coord[1]}`;

readFile('customers.geojson', function(err, data) { 
    const geojson = JSON.parse(data); 

	// const list = geojson.elements.map(element => ([element.lon, element.lat]));
	const list = geojson.elements.map(element => {
        if (element.type == "node") {
            return [element.lon, element.lat];
        } else {
            return [element.center.lon, element.center.lat];
        }
    });
	

    writeFile('customers.json', JSON.stringify(list), () => {})
}); 