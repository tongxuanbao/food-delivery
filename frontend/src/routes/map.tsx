import * as React from "react";
import { createFileRoute } from "@tanstack/react-router";
import Map from "../components/Map";

export const Route = createFileRoute("/map")({
  component: MapComponent,
});

function MapComponent() {
  return (
    <div className="p-2">
      <Map />
    </div>
  );
}
