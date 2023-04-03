import "antd/dist/reset.css";
import "./main.css";

import { Auth0Provider } from "@auth0/auth0-react";
import React, { PropsWithChildren } from "react";
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

const Auth = (props: PropsWithChildren<{ enabled?: boolean }>) => {
  if (!config.auth0.enabled) {
    return <>{props.children}</>;
  }

  return (
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
      {props.children}
    </Auth0Provider>
  );
};

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
  <React.StrictMode>
    <Auth enabled={config.auth0.enabled}>
      <RouterProvider router={router} />
    </Auth>
  </React.StrictMode>,
);
