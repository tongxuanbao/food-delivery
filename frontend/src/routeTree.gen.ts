/* prettier-ignore-start */

/* eslint-disable */

// @ts-nocheck

// noinspection JSUnusedGlobalSymbols

// This file is auto-generated by TanStack Router

// Import Routes

import { Route as rootRoute } from './routes/__root'
import { Route as MapImport } from './routes/map'
import { Route as IndexImport } from './routes/index'

// Create/Update Routes

const MapRoute = MapImport.update({
  path: '/map',
  getParentRoute: () => rootRoute,
} as any)

const IndexRoute = IndexImport.update({
  path: '/',
  getParentRoute: () => rootRoute,
} as any)

// Populate the FileRoutesByPath interface

declare module '@tanstack/react-router' {
  interface FileRoutesByPath {
    '/': {
      id: '/'
      path: '/'
      fullPath: '/'
      preLoaderRoute: typeof IndexImport
      parentRoute: typeof rootRoute
    }
    '/map': {
      id: '/map'
      path: '/map'
      fullPath: '/map'
      preLoaderRoute: typeof MapImport
      parentRoute: typeof rootRoute
    }
  }
}

// Create and export the route tree

export interface FileRoutesByFullPath {
  '/': typeof IndexRoute
  '/map': typeof MapRoute
}

export interface FileRoutesByTo {
  '/': typeof IndexRoute
  '/map': typeof MapRoute
}

export interface FileRoutesById {
  __root__: typeof rootRoute
  '/': typeof IndexRoute
  '/map': typeof MapRoute
}

export interface FileRouteTypes {
  fileRoutesByFullPath: FileRoutesByFullPath
  fullPaths: '/' | '/map'
  fileRoutesByTo: FileRoutesByTo
  to: '/' | '/map'
  id: '__root__' | '/' | '/map'
  fileRoutesById: FileRoutesById
}

export interface RootRouteChildren {
  IndexRoute: typeof IndexRoute
  MapRoute: typeof MapRoute
}

const rootRouteChildren: RootRouteChildren = {
  IndexRoute: IndexRoute,
  MapRoute: MapRoute,
}

export const routeTree = rootRoute
  ._addFileChildren(rootRouteChildren)
  ._addFileTypes<FileRouteTypes>()

/* prettier-ignore-end */

/* ROUTE_MANIFEST_START
{
  "routes": {
    "__root__": {
      "filePath": "__root.tsx",
      "children": [
        "/",
        "/map"
      ]
    },
    "/": {
      "filePath": "index.tsx"
    },
    "/map": {
      "filePath": "map.tsx"
    }
  }
}
ROUTE_MANIFEST_END */
