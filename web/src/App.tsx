import { Form, Space } from "antd";
import { useState } from "react";

import { Footer } from "./components/Footer";
import { QueueTable } from "./components/QueueTable";
import { SearchForm } from "./components/SearchForm";
import { YoutubeMetadataCard } from "./components/YoutubeMetadataCard";
import { useAPI } from "./lib/api";
import { useQueue } from "./lib/QueueContext";
import { YoutubeInfoJSON } from "./lib/types";

export const App = () => {
  const api = useAPI();
  const queue = useQueue();
  const [form] = Form.useForm<{ target: string }>();
  const [metadata, setMetadata] = useState<YoutubeInfoJSON>();

  const onReset = () => {
    form.resetFields();
    setMetadata(undefined);
  };

  const onConfirm = async (videoId: string, formatId: string) => {
    onReset();
    await api.beginDownload(videoId, formatId);
    await queue.poll();
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
      <QueueTable queue={queue.items} />
      <Footer />
    </Space>
  );
};
