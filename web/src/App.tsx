import { Form, Space } from "antd";
import { useEffect, useState } from "react";

import { QueueTable } from "./components/QueueTable";
import { SearchForm } from "./components/SearchForm";
import { YoutubeMetadataCard } from "./components/YoutubeMetadataCard";
import { useAPI } from "./lib/api";
import { YoutubeInfoJSON } from "./lib/types";
import { useQueuePoll } from "./lib/useQueuePoll";

export const App = () => {
  const api = useAPI();
  const [form] = Form.useForm<{ target: string }>();
  const [metadata, setMetadata] = useState<YoutubeInfoJSON>();
  const { queue, refreshQueue, pollIfPending } = useQueuePoll();

  useEffect(() => {
    refreshQueue();
  }, []);

  useEffect(() => {
    pollIfPending();
  }, [queue]);

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
