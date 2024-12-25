import { useEffect, useState } from "react";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/")({
  component: HomeComponent,
});

function HomeComponent() {
  const [data, setData] = useState("");

  useEffect(() => {
    // opening a connection to the server to begin receiving events from it
    const eventSource = new EventSource("http://localhost:8080/rates");

    // attaching a handler to receive message events
    eventSource.onmessage = (event) => {
      console.log(event);
      setData(event.data);
    };

    // terminating the connection on component unmount
    return () => {
      eventSource.close();
    };
  }, []);
  return (
    <div className="p-2">
      <h3>{`source1: ${data}`}</h3>
    </div>
  );
}
