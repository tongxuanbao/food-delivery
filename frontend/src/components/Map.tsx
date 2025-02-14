import { useEffect, useRef, useState } from "react";

type Coordinate = {
  x: Number;
  y: Number;
};

type Restaurant = {
  id: Number;
  status: String;
  coordinate: Coordinate;
};

type Driver = {
  id: Number;
  coordinate: Coordinate;
  route: Array<Coordinate>;
  speed: Number;
  status: Number;
};

type Customer = {
  id: Number;
  coordinate: Coordinate;
};

type InitialResponse = {
  restaurants: Restaurant[];
  drivers: Driver[];
  customers: Customer[];
};

const RESTAURANT_SIZE = 150;

const MapComponent = () => {
  const customerRef = useRef<HTMLImageElement | null>(null);
  const carRef = useRef<HTMLImageElement | null>(null);
  const mapRef = useRef<HTMLImageElement | null>(null);
  const redDotRef = useRef<HTMLImageElement | null>(null);
  const greenDotRef = useRef<HTMLImageElement | null>(null);

  /* Data */
  const [restaurants, setRestaurants] = useState<Array<Restaurant>>([]);
  const [drivers, setDrivers] = useState<Array<Driver>>([]);
  const [customers, setCustomers] = useState<Array<Customer>>([]);

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

  useEffect(() => {
    // opening a connection to the server to begin receiving events from it
    const eventSource = new EventSource("http://localhost:8080/route");

    // attaching a handler to receive message events
    eventSource.onmessage = (event) => {
      const d = JSON.parse(event.data);
      setData(d);
    };

    eventSource.addEventListener("initial", (event) => {
      const data = JSON.parse(event.data) as InitialResponse;
      console.log(data);
      const r = data.restaurants.map((restaurant) => ({
        ...restaurant,
        status: "red",
      }));
      setRestaurants(r);
      setDrivers(data.drivers);
      setCustomers(data.customers);
    });

    eventSource.addEventListener("restaurant", (event) => {
      const d = JSON.parse(event.data);
      console.log(d);
      setRestaurants((restaurants) => {
        const restaurant = d.restaurant;
        return restaurants.map((r) => {
          if (r.id === restaurant.id) {
            return {
              ...r,
              status: d.event === "order_prepared" ? "green" : "red",
            };
          } else {
            return r;
          }
        });
      });
    });

    eventSource.addEventListener("driver", (event) => {
      const d = JSON.parse(event.data);
      console.log(d);
      if (d.event === "init_drivers") {
        setDrivers(d.drivers);
      }
    });

    // terminating the connection on component unmount
    return () => {
      eventSource.close();
    };
  }, []);

  useEffect(draw, [restaurants]);

  function draw() {
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
      if (restaurant.status == "ref") console.log("YESSS");
      context?.drawImage(
        restaurant.status === "red" ? redDotRef.current : greenDotRef.current,
        restaurant.coordinate.x,
        restaurant.coordinate.y,
        RESTAURANT_SIZE,
        RESTAURANT_SIZE,
      );
    });

    drivers.forEach((driver) => {
      context?.drawImage(
        customerRef.current,
        driver.coordinate.x,
        driver.coordinate.y,
        150,
        150,
      );
    });
  }

  return (
    <div className="flex-grow relative overflow-hidden rounded-xl border border-dashed border-gray-400">
      <canvas ref={canvasRef} className="h-full w-full"></canvas>
      <img src="customer.png" ref={customerRef} className="hidden" />
      <img src="car.svg" ref={carRef} className="hidden" />
      <img src="map2.png" ref={mapRef} className="hidden" />
      <img src="RedDot.svg" ref={redDotRef} className="hidden" />
      <img src="GreenDot.svg" ref={greenDotRef} className="hidden" />
    </div>
  );
};

export default MapComponent;
