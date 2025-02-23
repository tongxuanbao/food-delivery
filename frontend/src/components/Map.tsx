import { useEffect, useRef } from "react";

type Coordinate = {
  x: number;
  y: number;
};

type Restaurant = {
  id: number;
  status: String;
  coordinate: Coordinate;
};

type Driver = {
  id: number;
  coordinate: Coordinate;
  route: Array<Coordinate>;
  speed: number;
  status: number;
  currentPosition: number;
};

type Customer = {
  id: number;
  coordinate: Coordinate;
};

type InitialResponse = {
  restaurants: Restaurant[];
  drivers: Driver[];
  customers: Customer[];
};

const ICON_SIZE = 130;

const testest = [
  {
    x: 789.3277381991164,
    y: 11.537644577026368,
  },
  {
    x: 760.7962935176897,
    y: 1.4319204495276936,
  },
  {
    x: 675.8474672953289,
    y: 184.73540463859655,
  },
  {
    x: 675.4601626172479,
    y: 202.61501483210813,
  },
  {
    x: 746.724223516017,
    y: 205.25808968720617,
  },
  {
    x: 792.4261756237644,
    y: 213.34279250627685,
  },
  {
    x: 805.8527378223562,
    y: 219.71727319411292,
  },
  {
    x: 832.8349637820398,
    y: 248.9466385076076,
  },
  {
    x: 844.7123072683077,
    y: 253.29995374326356,
  },
  {
    x: 1192.64101041629,
    y: 221.5829754440873,
  },
  {
    x: 1201.807221146572,
    y: 122.23469562765995,
  },
  {
    x: 1171.5974561997405,
    y: 116.17120966028286,
  },
  {
    x: 1063.7976539235058,
    y: -19.09044950037827,
  },
  {
    x: 1039.9138653978891,
    y: -23.59914772598832,
  },
  {
    x: 1022.8724595246493,
    y: 164.83464944156603,
  },
  {
    x: 844.8414088308076,
    y: 185.20182894483028,
  },
];

const MapComponent = () => {
  const customerRef = useRef<HTMLImageElement | null>(null);
  const restaurantRef = useRef<HTMLImageElement | null>(null);
  const carRef = useRef<HTMLImageElement | null>(null);
  const mapRef = useRef<HTMLImageElement | null>(null);

  /* Data */
  const restaurants = useRef<Map<number, Restaurant>>(new Map());
  const drivers = useRef<Map<number, Driver>>(new Map());
  const customers = useRef<Map<number, Customer>>(new Map());

  /* Canvas */
  const canvasRef = useRef<HTMLCanvasElement | null>(null);
  const drawQueued = useRef<boolean>(false);
  function drawMap(context: CanvasRenderingContext2D) {
    if (!mapRef.current) return;
    context.drawImage(mapRef.current, 0, 0);
  }
  function drawRestaurants(context: CanvasRenderingContext2D) {
    restaurants.current.forEach((restaurant) => {
      if (!restaurantRef.current) return;
      context.drawImage(
        restaurantRef.current,
        restaurant.coordinate.x - Math.round(ICON_SIZE / 2),
        restaurant.coordinate.y - ICON_SIZE,
        ICON_SIZE,
        ICON_SIZE,
      );
    });
  }
  function drawDrivers(context: CanvasRenderingContext2D) {
    drivers.current.forEach((driver) => {
      if (!carRef.current) return;

      // Draw path of the driver
      context.beginPath();
      context.moveTo(driver.coordinate.x, driver.coordinate.y);
      for (let i = driver.currentPosition + 1; i < driver.route.length; i++) {
        const nextCoordinate = driver.route.at(i);
        if (!nextCoordinate) continue;
        context.lineTo(nextCoordinate.x, nextCoordinate.y);
      }
      context.lineWidth = 10;
      context.stroke();

      // TODO: calculate rotation based on the route driver is going
      // Draw driver
      context.drawImage(
        carRef.current,
        driver.coordinate.x - 50,
        driver.coordinate.y - 50,
        100,
        100,
      );
    });
  }
  function drawCustomers(context: CanvasRenderingContext2D) {
    customers.current.forEach((customer) => {
      if (!customerRef.current) return;
      context?.drawImage(
        customerRef.current,
        customer.coordinate.x - Math.round(ICON_SIZE / 2),
        customer.coordinate.y - Math.round(ICON_SIZE / 2),
        ICON_SIZE,
        ICON_SIZE,
      );
    });
  }
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

    drawMap(context);

    // context.beginPath();
    // context.moveTo(testest[0].x * 3.125, testest[0].y * 3.125);
    // for (let i = 1; i < 16; i++) {
    //   context.lineTo(testest[i].x * 3.125, testest[i].y * 3.125);
    // }
    // context.lineWidth = 10;
    // context.stroke();

    drawDrivers(context);
    drawRestaurants(context);
    drawCustomers(context);

    drawQueued.current = false;
  }

  useEffect(() => {
    // opening a connection to the server to begin receiving events from it
    const eventSource = new EventSource("http://localhost:8080/route");

    eventSource.addEventListener("initial", (event) => {
      const data = JSON.parse(event.data) as InitialResponse;
      restaurants.current = new Map(
        data.restaurants.map((restaurant) => [restaurant.id, restaurant]),
      );
      drivers.current = new Map(
        data.drivers.map((driver) => [driver.id, driver]),
      );
      customers.current = new Map(
        data.customers.map((customer) => [customer.id, customer]),
      );
    });

    eventSource.addEventListener("restaurant", (event) => {
      const d = JSON.parse(event.data);
      if (d.event === "init_restaurants") {
        restaurantRef.current = d.restaurants;
      }
    });

    eventSource.addEventListener("driver", (event) => {
      const d = JSON.parse(event.data);
      if (d.event === "init_drivers") {
        // setDrivers(d.drivers);
        drivers.current = new Map(
          d.drivers.map((driver) => [driver.id, driver]),
        );
      }
      if (d.event === "driver") {
        console.log("set", d.driver.id, d.driver);
        drivers.current.set(d.driver.id, d.driver);
        console.log(drivers.current);
      }
    });

    // terminating the connection on component unmount
    return () => {
      eventSource.close();
    };
  }, []);

  function lerp(start: number, end: number, t: number): number {
    return start + (end - start) * t;
  }
  function calculateDriverPosition() {
    drivers.current.forEach((driver) => {
      const driverCoordinate = driver.coordinate;
      const targetCoordinate =
        driver.route.at(driver.currentPosition + 1) ?? driverCoordinate;
      driver.coordinate = {
        x: lerp(driverCoordinate.x, targetCoordinate.x, 0.01),
        y: lerp(driverCoordinate.y, targetCoordinate.y, 0.01),
      };
    });
  }

  useEffect(() => {
    let raf: number;
    const start = () => {
      draw();
      calculateDriverPosition();
      raf = requestAnimationFrame(start);
    };

    start();

    return () => {
      cancelAnimationFrame(raf);
    };
  }, []);

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
