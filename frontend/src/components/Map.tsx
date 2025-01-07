import { useEffect, useRef, useState } from "react";

type Coordinate = {
  x: Number;
  y: Number;
};

const MapComponent = () => {
  const mapRef = useRef<HTMLImageElement | null>(null);
  const dotRef = useRef<HTMLImageElement | null>(null);

  /* Canvas */
  const canvasRef = useRef<HTMLCanvasElement | null>(null);
  useEffect(() => {
    if (!canvasRef.current) return;

    const canvas = canvasRef.current;
    const context = canvasRef.current.getContext("2d");
    const { offsetWidth: width, offsetHeight: height } = canvas;

    if (!context) return;

    const zoomLevel = Math.min(6000 / width, 3375 / height);
    context.canvas.width = width * zoomLevel;
    context.canvas.height = height * zoomLevel;

    if (mapRef.current)
      context?.drawImage(mapRef.current, 0, (6000 - context.canvas.width) / 2);
  }, []);

  /* Data */
  const [data, setData] = useState<Coordinate>({});
  const [restaurants, setRestaurant] = useState<Array<Coordinate>>([]);
  useEffect(() => {
    // opening a connection to the server to begin receiving events from it
    const eventSource = new EventSource("http://localhost:8080/route");

    // attaching a handler to receive message events
    eventSource.onmessage = (event) => {
      const d = JSON.parse(event.data);
      setData(d);
    };

    eventSource.addEventListener("restaurant", (event) => {
      const d = JSON.parse(event.data);
      console.log("restaurant", d);
      setRestaurant(d);
    });

    // terminating the connection on component unmount
    return () => {
      eventSource.close();
    };
  }, []);

  useEffect(() => {
    if (!canvasRef.current) return;

    const canvas = canvasRef.current;
    const context = canvasRef.current.getContext("2d");
    const { offsetWidth: width, offsetHeight: height } = canvas;

    if (!context) return;

    const zoomLevel = Math.min(6000 / width, 3375 / height);
    context.canvas.width = width * zoomLevel;
    context.canvas.height = height * zoomLevel;

    if (mapRef.current) context?.drawImage(mapRef.current, 0, 0);

    if (dotRef.current && data.position) {
      context.save();
      context.beginPath();
      context.strokeStyle = "green";
      context.moveTo(data.position?.x, data.position?.y);
      data?.route.forEach((pos) => {
        context.lineTo(pos.x, pos.y);
      });
      context.stroke();
      context.restore();

      context?.drawImage(
        dotRef.current,
        data.position?.x,
        data.position?.y,
        40,
        40,
      );
    }
  }, [data]);

  useEffect(() => {
    if (!canvasRef.current) return;

    const canvas = canvasRef.current;
    const context = canvasRef.current.getContext("2d");
    const { offsetWidth: width, offsetHeight: height } = canvas;

    if (!context) return;

    const zoomLevel = Math.min(6000 / width, 3375 / height);
    context.canvas.width = width * zoomLevel;
    context.canvas.height = height * zoomLevel;

    if (mapRef.current) context?.drawImage(mapRef.current, 0, 0);

    restaurants.forEach((restaurant) => {
      context?.drawImage(dotRef.current, restaurant.x, restaurant.y, 70, 70);
    });
  }, [restaurants]);

  return (
    <div className="flex-grow relative overflow-hidden rounded-xl border border-dashed border-gray-400">
      <canvas ref={canvasRef} className="h-full w-full"></canvas>
      <img src="map2.png" ref={mapRef} className="hidden" />
      <img src="RedDot.svg" ref={dotRef} className="hidden" />
    </div>
  );
};

export default MapComponent;
