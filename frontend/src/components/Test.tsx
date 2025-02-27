import mapboxgl from "mapbox-gl";
import "mapbox-gl/dist/mapbox-gl.css";
import { useEffect, useRef } from "react";

import data from "../../restaurantCustomerClean.json";

const TestComponent = () => {
  const mapRef = useRef();
  const mapContainerRef = useRef();
  useEffect(() => {
    mapboxgl.accessToken =
      "pk.eyJ1IjoiYmFvdG9uZ3h1YW4iLCJhIjoiY20xZ3J2MWZlMDduaDJscHpoY2RscmZtaiJ9.qPNALF9fS6rBZJ-sxi1jFA";
    mapRef.current = new mapboxgl.Map({
      container: mapContainerRef.current,
      zoom: 16.47,
      center: [151.208563, -33.865535],
      scrollZoom: false,
      boxZoom: false,
      dragRotate: false,
      dragPan: false,
      keyboard: false,
      doubleClickZoom: false,
      touchZoomRotate: false,
    });

    return () => {
      mapRef.current.remove();
    };
  }, []);

  setTimeout(() => {
    // const pixels = coords.map((restaurant) =>
    //   mapRef.current.project(restaurant),
    // );
    const result = data.map((entry) => {
      const route = entry.route.map((r) => mapRef.current.project(r));
      return {
        restaurantId: entry.restaurantId,
        customerId: entry.customerId,
        route,
      };
    });
    console.log(result);
  }, 1000);

  return (
    <div
      style={{ width: "100%", height: "calc(100vh - 61px)", overflow: "auto" }}
    >
      <div
        id="map-container"
        ref={mapContainerRef}
        style={{
          width: "1920px",
          height: "1080px",
          backgroundImage: "url(/map2.png)",
          backgroundSize: "cover",
          backgroundPosition: "center",
          backgroundRepeat: "no-repeat",
        }}
      />
    </div>
  );
};

export default TestComponent;
