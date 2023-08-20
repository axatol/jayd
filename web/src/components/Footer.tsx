import { GithubFilled } from "@ant-design/icons";
import { Button, Space, Typography } from "antd";

import { config } from "../config";

export const Footer = () => (
  <Space>
    <Button
      type="text"
      icon={<GithubFilled />}
      rel="noopener noreferrer"
      target="_blank"
      href="https://github.com/axatol/jayd"
    />

    <Typography.Text>{config.commitSha}</Typography.Text>
  </Space>
);
