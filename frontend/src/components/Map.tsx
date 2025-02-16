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

const ICON_SIZE = 130;

// Hook
function useWindowSize() {
  // Initialize state with undefined width/height so server and client renders match
  // Learn more here: https://joshwcomeau.com/react/the-perils-of-rehydration/
  const [windowSize, setWindowSize] = useState({
    width: undefined,
    height: undefined,
  });
  useEffect(() => {
    // Handler to call on window resize
    function handleResize() {
      // Set window width/height to state
      setWindowSize({
        width: window.innerWidth,
        height: window.innerHeight,
      });
    }
    // Add event listener
    window.addEventListener("resize", handleResize);
    // Call handler right away so state gets updated with initial window size
    handleResize();
    // Remove event listener on cleanup
    return () => window.removeEventListener("resize", handleResize);
  }, []); // Empty array ensures that effect is only run on mount
  return windowSize;
}

const MapComponent = () => {
  const customerRef = useRef<HTMLImageElement | null>(null);
  const restaurantRef = useRef<HTMLImageElement | null>(null);
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
  function draw() {
    if (!canvasRef.current) return;

    const canvas = canvasRef.current;
    const context = canvasRef.current.getContext("2d");
    const { offsetWidth: width, offsetHeight: height } = canvas;
    if (!context) return;

    const zoomLevel = Math.min(6000 / width, 3375 / height);
    context.canvas.width = width * zoomLevel;
    context.canvas.height = height * zoomLevel;
    context.imageSmoothingEnabled = false;

    if (mapRef.current) context?.drawImage(mapRef.current, 0, 0);

    restaurants.forEach((restaurant) => {
      context?.drawImage(
        restaurantRef.current,
        restaurant.coordinate.x - Math.round(ICON_SIZE / 2),
        restaurant.coordinate.y - ICON_SIZE,
        ICON_SIZE,
        ICON_SIZE,
      );
    });

    drivers.forEach((driver) => {
      context?.drawImage(
        carRef.current,
        driver.coordinate.x - 50,
        driver.coordinate.y - 50,
        100,
        100,
      );
    });

    customers.forEach((customer) => {
      context?.drawImage(
        customerRef.current,
        customer.coordinate.x - Math.round(ICON_SIZE / 2),
        customer.coordinate.y - Math.round(ICON_SIZE / 2),
        ICON_SIZE,
        ICON_SIZE,
      );
    });
  }

  useEffect(() => {
    // opening a connection to the server to begin receiving events from it
    const eventSource = new EventSource("http://localhost:8080/route");

    eventSource.addEventListener("initial", (event) => {
      const data = JSON.parse(event.data) as InitialResponse;
      console.log("initial", data);
      const r = data.restaurants.map((restaurant) => ({
        ...restaurant,
        status: "red",
      }));
      setRestaurants(r);
      setDrivers(data.drivers);
      setCustomers(data.customers);
      console.log(data.restaurants.map((restaurant) => restaurant.coordinate));
      // console.log(data.customers.map((restaurant) => restaurant.coordinate));
    });

    eventSource.addEventListener("restaurant", (event) => {
      const d = JSON.parse(event.data);
      console.log("restaurant", d);
      if (d.event === "init_restaurants") {
        setRestaurants(d.restaurants);
      } else {
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
      }
    });

    eventSource.addEventListener("driver", (event) => {
      const d = JSON.parse(event.data);
      console.log("driver", d);
      if (d.event === "init_drivers") {
        setDrivers(d.drivers);
      }
    });

    // terminating the connection on component unmount
    return () => {
      eventSource.close();
    };
  }, []);

  useEffect(() => {
    requestAnimationFrame(draw);
    return;
  }, [restaurants, drivers, customers]);

  useEffect(() => {
    // Handler to call on window resize
    const handler = () => requestAnimationFrame(draw);
    // Add event listener
    window.addEventListener("resize", handler);
    // Call handler right away so state gets updated with initial window size
    handler();
    // Remove event listener on cleanup
    return () => window.removeEventListener("resize", handler);
  }, []); // Empty array ensures that effect is only run on mount

  return (
    <div className="flex-grow relative overflow-hidden rounded-xl border border-dashed border-gray-400">
      <canvas ref={canvasRef} className="h-full w-full"></canvas>
      <img src="home.svg" ref={customerRef} className="hidden" />
      <img src="car.svg" ref={carRef} className="hidden" />
      <img src="map.png" ref={mapRef} className="hidden" />
      <img src="restaurant.svg" ref={restaurantRef} className="hidden" />
    </div>
  );
};

export default MapComponent;
