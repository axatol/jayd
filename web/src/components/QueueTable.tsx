import {
  DeleteOutlined,
  ReloadOutlined,
  WarningOutlined,
} from "@ant-design/icons";
import { Table, Image, Card, Tag, Space, Button } from "antd";

import { AsyncButton } from "./AsyncButton";
import { DownloadButton } from "./DownloadButton";
import { YoutubeChannelLink, YoutubeVideoLink } from "./Links";
import { useAPI } from "../lib/api";
import { QueueItem, VideoFormat } from "../lib/types";
import { useQueuePoll } from "../lib/useQueuePoll";

export const QueueTable = (props: { queue?: QueueItem[] }) => (
  <Card hoverable style={{ cursor: "auto" }} bodyStyle={{ padding: 0 }}>
    <Table
      rowKey={(item) => item.id}
      dataSource={props.queue}
      columns={[
        {
          key: "thumbnail",
          title: "Thumbnail",
          align: "center",
          dataIndex: ["data", "thumbnail"],
          render: (thumbnail, { data }) => (
            <Image src={thumbnail} alt={data.title} style={{ maxWidth: 100 }} />
          ),
        },
        {
          key: "title",
          title: "Title",
          dataIndex: ["data", "title"],
          render: (title, { data }) => (
            <YoutubeVideoLink id={data.id} title={title} />
          ),
        },
        {
          key: "uploader",
          title: "Uploader",
          dataIndex: ["data", "uploader"],
          render: (uploader, { data }) => (
            <YoutubeChannelLink id={data.uploader_id} name={uploader} />
          ),
        },
        {
          key: "duration",
          title: "Duration",
          dataIndex: ["data", "duration_string"],
        },
        {
          key: "format",
          title: "Format",
          dataIndex: ["data", "formats"],
          render: (_, { data }) => (
            <Space direction="vertical">
              {data.formats.map((format) => (
                <Format key={format.format_id} format={format} />
              ))}
            </Space>
          ),
        },
        {
          key: "filesize",
          title: "Filesize",
          render: (_, { data }) => {
            const total = data.formats.reduce(
              (total, { filesize }) => total + filesize,
              0,
            );

            return `${(total / 1000000).toFixed(2)} MB`;
          },
        },
        {
          key: "actions",
          align: "right",
          render: (_, item) => <Actions item={item} />,
        },
      ]}
    />
  </Card>
);

const Actions = (props: { item: QueueItem }) => {
  const { refreshQueue } = useQueuePoll();
  const { completed, failed, data } = props.item;
  const { id, format_id, filename } = data;
  const api = useAPI();

  const retry = async () => {
    await api.beginDownload(
      `https://youtube.com/watch?v=${id}`,
      format_id,
      true,
    );
    await refreshQueue();
  };

  return (
    <Space>
      {failed ? (
        <>
          <Tag
            color="error"
            icon={<WarningOutlined color="danger" />}
            style={{ margin: 0 }}
          >
            Failed
          </Tag>

          <AsyncButton
            tooltip="Retry"
            shape="circle"
            icon={<ReloadOutlined />}
            onClick={retry}
          />
        </>
      ) : (
        <DownloadButton
          loading={!completed ? true : undefined}
          href={filename}
        />
      )}

      {completed && (
        <Button
          danger
          icon={<DeleteOutlined />}
          shape="circle"
          onClick={() => api.deleteQueueItem(data.id, format_id)}
        />
      )}
    </Space>
  );
};

const Format = (props: { format?: VideoFormat }) => {
  if (!props.format) {
    return <>N/A</>;
  }

  return (
    <>
      {props.format.height > 0 && (
        <FormatField name="Resolution" value={`${props.format.height}p`} />
      )}
      {props.format.video_ext !== "none" && (
        <FormatField name="Video" value={props.format.video_ext} />
      )}
      {props.format.audio_ext !== "none" && (
        <FormatField name="Audio" value={props.format.audio_ext} />
      )}
    </>
  );
};

const FormatField = (props: { name: string; value: string }) => (
  <>
    {props.name}: {props.value}
    <br />
  </>
);
