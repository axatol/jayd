import { useAuth0 } from "@auth0/auth0-react";
import { Button, Spin, Typography } from "antd";
import { useEffect } from "react";
import { useNavigate } from "react-router-dom";

import { App } from "./App";
import { Centre } from "./components/Centre";
import { config } from "./config";
import { QueueProvider } from "./lib/QueueContext";

export const Root = () => {
  const navigate = useNavigate();
  const { isAuthenticated, isLoading, error, loginWithRedirect } = useAuth0();

  useEffect(() => {
    if (!config.auth0.enabled) {
      return;
    }

    if (!isLoading && !isAuthenticated && !error) {
      loginWithRedirect({ appState: { returnTo: window.location.pathname } });
    }
  }, [isLoading, isAuthenticated, error]);

  if (config.auth0.enabled && isLoading) {
    return (
      <Centre>
        <Spin size="large" />
        <Typography.Text strong>Logging you in...</Typography.Text>
      </Centre>
    );
  }

  if (config.auth0.enabled && error) {
    return (
      <Centre>
        <Typography.Text strong style={{ paddingBottom: 8 }}>
          Something went wrong: {error.message}
        </Typography.Text>

        <Button
          type="primary"
          shape="round"
          onClick={() => {
            navigate({ pathname: "/", search: "" });
            window.location.reload();
          }}
        >
          Refresh
        </Button>
      </Centre>
    );
  }

  return (
    <QueueProvider>
      <App />
    </QueueProvider>
  );
};
