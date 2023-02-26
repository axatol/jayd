import { DownloadOutlined } from "@ant-design/icons";
import { Button, ButtonProps, Progress } from "antd";
import { AxiosInstance } from "axios";
import { useState } from "react";

import { useAPI } from "../lib/api";

const downloadFile = async (
  api: AxiosInstance,
  href: string,
  setProgress: (progress?: number) => void,
) => {
  const response = await api.get(`/static/${href}`, {
    responseType: "blob",
    onDownloadProgress: (event) => {
      console.log("setProgress", event.progress);
      setProgress(event.progress);
    },
  });

  const data = URL.createObjectURL(response.data);
  const anchor = document.createElement("a");
  anchor.href = data;
  anchor.download = href;
  anchor.click();
  URL.revokeObjectURL(data);
  anchor.remove();
};

export interface DownloadButtonProps extends ButtonProps {
  href: string;
}

export const DownloadButton = ({
  href,
  children,
  ...props
}: DownloadButtonProps) => {
  const api = useAPI();
  const [loading, setLoading] = useState(false);
  const [progress, setProgress] = useState<number>();

  const download = async () => {
    setLoading(true);
    await downloadFile(api.api, href, setProgress);
    setLoading(false);
  };

  if (progress !== undefined && progress < 1) {
    return <Progress percent={progress} style={{ width: 120 }} />;
  }

  return (
    <Button
      {...props}
      loading={loading}
      onClick={download}
      type="primary"
      shape="round"
      icon={<DownloadOutlined />}
      download
    >
      {children ?? "Download"}
    </Button>
  );
};
