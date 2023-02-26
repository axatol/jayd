import "antd/dist/reset.css";
import "./main.css";

import { Auth0Provider } from "@auth0/auth0-react";
import React from "react";
import ReactDOM from "react-dom/client";
import {
  createBrowserRouter,
  redirect,
  RouterProvider,
} from "react-router-dom";

import { config } from "./config";
import { Logout } from "./Logout";
import { Root } from "./Root";

const router = createBrowserRouter([
  { path: "/", element: <Root /> },
  { path: "/logout", element: <Logout /> },
]);

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
  <React.StrictMode>
    <Auth0Provider
      domain={config.auth0.domain}
      clientId={config.auth0.clientId}
      onRedirectCallback={(state) => redirect(state?.returnTo ?? "/")}
      authorizationParams={{
        redirect_uri: config.auth0.callbackUri,
        audience: config.api.audience,
        scope: config.api.scope,
      }}
    >
      <RouterProvider router={router} />
    </Auth0Provider>
  </React.StrictMode>,
);
