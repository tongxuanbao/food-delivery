import { createFileRoute } from "@tanstack/react-router";
// import "mapbox-gl/dist/mapbox-gl.css";

export const Route = createFileRoute("/monitor")({
  component: MonitorComponent,
});

function MonitorComponent() {
  // const mapRef = useRef();
  // const mapContainerRef = useRef();
  // useEffect(() => {
  //   mapboxgl.accessToken =
  //     "pk.eyJ1IjoiYmFvdG9uZ3h1YW4iLCJhIjoiY20xZ3J2MWZlMDduaDJscHpoY2RscmZtaiJ9.qPNALF9fS6rBZJ-sxi1jFA";

  //   mapRef.current = new mapboxgl.Map({
  //     container: mapContainerRef.current,
  //     zoom: 16.47,
  //     center: [151.208563, -33.865535],
  //     scrollZoom: false,
  //     boxZoom: false,
  //     dragRotate: false,
  //     dragPan: false,
  //     keyboard: false,
  //     doubleClickZoom: false,
  //     touchZoomRotate: false,
  //   });

  //   return () => {
  //     mapRef.current.remove();
  //   };
  // }, []);

  // setTimeout(() => {
  //   console.log(
  //     coordinates.map((coord) => {
  //       return mapRef.current.project(coord);
  //     }),
  //   );
  // }, 1000);

  return (
    <div
      style={{ width: "1920px", height: "1080px" }}
      id="map-container"
      // ref={mapContainerRef}
    ></div>
  );
}
