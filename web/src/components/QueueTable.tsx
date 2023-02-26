import { DeleteOutlined, WarningOutlined } from "@ant-design/icons";
import { Table, Image, Card, Tag, Tooltip, Space, Button } from "antd";

import { DownloadButton } from "./DownloadButton";
import { YoutubeChannelLink, YoutubeVideoLink } from "./Links";
import { useAPI } from "../lib/api";
import { QueueItem, VideoFormat } from "../lib/types";

export const QueueTable = (props: { queue?: QueueItem[] }) => (
  <Card hoverable style={{ cursor: "auto" }} bodyStyle={{ padding: 0 }}>
    <Table
      rowKey={(item) => item.data.id + item.selected_format_id}
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
          dataIndex: ["format"],
          render: (format) => <Format format={format} />,
        },
        {
          key: "filesize",
          title: "Filesize",
          dataIndex: ["format", "filesize"],
          render: (value) => `${(value / 1000000).toFixed(2)} MB`,
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
  const { is_completed, is_failed, format, data } = props.item;
  const api = useAPI();

  if (is_failed || !format) {
    return (
      <Tooltip title="An error occurred">
        <Tag color="error" icon={<WarningOutlined color="danger" />}>
          Failed
        </Tag>
      </Tooltip>
    );
  }

  const { audio_ext, video_ext, format_id } = format;
  const ext = video_ext === "none" ? audio_ext : video_ext;
  const filename = `${data.id}_${format_id}.${ext}`;

  return (
    <Space>
      <DownloadButton
        loading={!is_completed ? true : undefined}
        href={filename}
      />

      <Button
        danger
        icon={<DeleteOutlined />}
        shape="circle"
        onClick={() => api.deleteQueueItem(data.id, format.format_id)}
      />
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
