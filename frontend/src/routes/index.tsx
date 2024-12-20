import { useEffect, useState } from "react";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/")({
  component: HomeComponent,
});

function HomeComponent() {
  const [data, setData] = useState("");
  const [data2, setData2] = useState("");

  useEffect(() => {
    // opening a connection to the server to begin receiving events from it
    const eventSource = new EventSource("http://localhost:8080/rates");
    const eventSource2 = new EventSource("http://localhost:8080/rates");

    // attaching a handler to receive message events
    eventSource.onmessage = (event) => {
      setData(event.data);
    };

    // attaching a handler to receive message events
    eventSource2.onmessage = (event) => {
      setData2(event.data);
    };

    // terminating the connection on component unmount
    return () => {
      eventSource.close();
      eventSource2.close();
    };
  }, []);
  return (
    <div className="p-2">
      <h3>{`source1: ${data}`}</h3>
      <h3>{`source2: ${data2}`}</h3>
    </div>
  );
}
