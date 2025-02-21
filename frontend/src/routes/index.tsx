import { createFileRoute } from "@tanstack/react-router";
import Map from "../components/Map";

export const Route = createFileRoute("/")({
  component: HomeComponent,
});

function HomeComponent() {
  return (
    <>
      <main className="lg:pl-20 h-screen">
        <div className="xl:pl-96 h-full">
          <div className="h-full flex px-4 py-10 sm:px-6 lg:px-8 lg:py-6">
            <Map />
            {/* <Test /> */}
          </div>
        </div>
      </main>

      <aside className="fixed inset-y-0 left-20 hidden w-96 overflow-y-auto border-r border-gray-200 px-4 py-6 sm:px-6 lg:px-8 xl:block bg-green">
        <div className="h-full flex">
          <div className="flex-grow relative overflow-hidden rounded-xl border border-dashed border-gray-400 opacity-75 bg-green-100">
            home
          </div>
        </div>
      </aside>
    </>
  );
}
