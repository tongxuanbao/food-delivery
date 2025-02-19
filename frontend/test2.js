import {readFile, writeFile} from "fs";

const getRoute = async (start, end) => {
 const url = `https://api.mapbox.com/directions/v5/mapbox/driving/${start.lng},${start.lat};${end.lng},${end.lat}?geometries=geojson&access_token=pk.eyJ1IjoiYmFvdG9uZ3h1YW4iLCJhIjoiY20xZ3J2MWZlMDduaDJscHpoY2RscmZtaiJ9.qPNALF9fS6rBZJ-sxi1jFA`;

  try {
    const response = await fetch(url, {
      method: "GET",
      headers: {
        "Content-Type": "application/json"
      }
    });

    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }

    const data = await response.json();
    return data
  } catch (error) {
    console.error("Error:", error);
  }
}
function delay(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}


const result = await getRoute({ "lng": 151.20784840000113, "lat": -33.867116100001006 }, { "lng": 151.20536110000057, "lat": -33.86752950000053 })
writeFile('test.json', JSON.stringify(result), () => {})

// readFile('restaurantsCoord.json', function(err, resData) { 
//     const restaurants = JSON.parse(resData); 

//     readFile('customersCoord.json', async function(err, cusData) {
//         const customers = JSON.parse(cusData)
//         const result = []
//         for (let i = 0; i < restaurants.length; i++) {
//           for (let j = 0; j < customers.length; j++) {
//             await delay(150);
//             console.log(`restaurants: ${i}(${restaurants.length}), customers: ${j}(${customers.length})`)
//             const data = await getRoute(restaurants[i], customers[j]);
//             result.push({
//               restaurant: {
//                 index: i,
//                 coordinate: restaurants[i],
//               },
//               customer: {
//                 index: j,
//                 coordinate: customers[j]
//               },
//               result: data
//             })
//           }
//         }

//         writeFile('routeResults.json', JSON.stringify(result), () => {})
//     })
// }); 