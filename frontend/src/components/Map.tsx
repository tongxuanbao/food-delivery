import React, { useEffect, useRef } from "react";
import mapboxgl, { Map } from "mapbox-gl";
import "mapbox-gl/dist/mapbox-gl.css";

const MapComponent = () => {
  return (
    <div
      style={{ width: "100%", height: "calc(100vh - 61px)", overflow: "auto" }}
    >
      <div
        id="map-container"
        style={{
          width: "1920px",
          height: "1080px",
          backgroundImage: "url(/map.png)",
          backgroundSize: "cover",
          backgroundPosition: "center",
          backgroundRepeat: "no-repeat",
        }}
      />
    </div>
  );
};

export default MapComponent;
