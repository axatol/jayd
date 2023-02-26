import { Form, Space } from "antd";
import { useEffect, useState } from "react";

import { QueueTable } from "./components/QueueTable";
import { SearchForm } from "./components/SearchForm";
import { YoutubeMetadataCard } from "./components/YoutubeMetadataCard";
import { useAPI } from "./lib/api";
import { QueueItem, YoutubeInfoJSON } from "./lib/types";

let intervalFn: (() => void) | undefined = undefined;
setInterval(() => {
  intervalFn?.();
}, 5000);

export const App = () => {
  const api = useAPI();
  const [form] = Form.useForm<{ target: string }>();
  const [metadata, setMetadata] = useState<YoutubeInfoJSON>();
  const [queue, setQueue] = useState<QueueItem[]>();

  useEffect(() => {
    refreshQueue();
    intervalFn = refreshQueue;
    return () => (intervalFn = undefined);
  }, []);

  const refreshQueue = async () => {
    const queue = await api.getQueue();
    setQueue(queue);
  };

  const onReset = () => {
    form.resetFields();
    setMetadata(undefined);
  };

  const onConfirm = async (videoId: string, formatId: string) => {
    onReset();
    await api.beginDownload(videoId, formatId);
    await refreshQueue();
  };

  return (
    <Space direction="vertical" style={{ padding: 8, width: "100%" }}>
      <SearchForm form={form} onMetadata={setMetadata} />

      {metadata && (
        <YoutubeMetadataCard
          metadata={metadata}
          onConfirm={onConfirm}
          onReset={onReset}
        />
      )}

      <QueueTable queue={queue} />
    </Space>
  );
};
