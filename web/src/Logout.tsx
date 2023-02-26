import { CheckCircleOutlined } from "@ant-design/icons";
import { useAuth0 } from "@auth0/auth0-react";
import { Button, Space, Typography } from "antd";
import { useNavigate } from "react-router-dom";

import { config } from "./config";

export const Logout = () => {
  const navigate = useNavigate();
  const { isAuthenticated, logout } = useAuth0();

  if (isAuthenticated) {
    logout({ logoutParams: { returnTo: config.auth0.logoutUri } });
  }

  return (
    <Space
      direction="vertical"
      align="center"
      style={{ width: "100%", height: "100%", margin: 8 }}
    >
      <Typography.Text>
        <CheckCircleOutlined /> You are logged out
      </Typography.Text>

      <Button
        shape="round"
        type="primary"
        onClick={() => navigate({ pathname: "/", search: "" })}
      >
        Return to home
      </Button>
    </Space>
  );
};
