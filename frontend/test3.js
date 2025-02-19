import {readFile, writeFile} from "fs";

readFile('restaurantCustomer.json', function(err, data) { 
    const pairs = JSON.parse(data);
    const result = pairs.map((pair, index) => {
        console.log(index)
        return ({
            restaurantId: pair.restaurant.index,
            customerId: pair.customer.index,
            route: pair.result.routes.at(0).geometry.coordinates,
        })
    })

    
    writeFile('test.json', JSON.stringify(result), () => {})
}); 