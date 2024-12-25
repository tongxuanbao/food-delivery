import { useEffect, useState } from "react";

const MapComponent = () => {
  const [data, setData] = useState({});

  useEffect(() => {
    // opening a connection to the server to begin receiving events from it
    const eventSource = new EventSource("http://localhost:8080/route");

    // attaching a handler to receive message events
    eventSource.onmessage = (event) => {
      const d = JSON.parse(event.data);
      console.log(d);
      setData(d);
    };

    // terminating the connection on component unmount
    return () => {
      eventSource.close();
    };
  }, []);
  return (
    <div
      className="overflow-auto self-center"
      style={{
        width: "min(100vw, 178vh)",
        height: "calc(min(100vh, 56vw) - 60px)",
      }}
    >
      <div
        id="map-container"
        style={{
          width: "1920px",
          height: "1080px",
          backgroundImage: "url(/map2.png)",
          backgroundSize: "cover",
          backgroundPosition: "center",
          backgroundRepeat: "no-repeat",
        }}
      >
        <div
          style={{
            width: "10px",
            height: "10px",
            position: "relative",
            background: "red",
            left: `${data.x}px`,
            top: `${data.y}px`,
          }}
        ></div>
      </div>
    </div>
  );
};

export default MapComponent;
