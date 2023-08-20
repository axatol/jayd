import {
  DeleteOutlined,
  ReloadOutlined,
  WarningOutlined,
} from "@ant-design/icons";
import { Table, Image, Card, Tag, Space } from "antd";

import { AsyncButton } from "./AsyncButton";
import { DownloadButton } from "./DownloadButton";
import { YoutubeChannelLink, YoutubeVideoLink } from "./Links";
import { useAPI } from "../lib/api";
import { useQueue } from "../lib/QueueContext";
import { QueueItem, VideoFormat } from "../lib/types";

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
          dataIndex: ["info", "thumbnail"],
          render: (thumbnail, { info: data }) => (
            <Image src={thumbnail} alt={data.title} style={{ maxWidth: 100 }} />
          ),
        },
        {
          key: "title",
          title: "Title",
          dataIndex: ["info", "title"],
          render: (title, { info: data }) => (
            <YoutubeVideoLink id={data.id} title={title} />
          ),
        },
        {
          key: "uploader",
          title: "Uploader",
          dataIndex: ["info", "uploader"],
          render: (uploader, { info: data }) => (
            <YoutubeChannelLink id={data.uploader_id} name={uploader} />
          ),
        },
        {
          key: "duration",
          title: "Duration",
          dataIndex: ["info", "duration_string"],
        },
        {
          key: "format",
          title: "Format",
          dataIndex: ["info", "formats"],
          render: (_, { info: data }) => (
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
          render: (_, { info: data }) => {
            const total = data.formats.reduce(
              (total, { filesize }) => total + filesize,
              0,
            );

            return `${(total / 1000000).toFixed(2)} MB`;
          },
        },
        {
          key: "added_at",
          title: "Added at",
          render: ({ added_at }) => {
            const date = new Date(added_at);
            return (
              <>
                {date.toLocaleDateString()} {date.toLocaleTimeString()}
              </>
            );
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
  const queue = useQueue();
  const { completed_at: completed, failed_at: failed, info: data } = props.item;
  const { id, format_id, filename } = data;
  const api = useAPI();

  const retryDownload = async () => {
    await api.beginDownload(
      `https://youtube.com/watch?v=${id}`,
      format_id,
      true,
    );
    queue.poll();
  };

  const deleteItem = (target: string, format: string) => async () => {
    await api.deleteQueueItem(target, format);
    queue.poll();
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
            onClick={retryDownload}
          />
        </>
      ) : (
        <DownloadButton
          loading={!completed ? true : undefined}
          href={filename}
        />
      )}

      {completed && (
        <AsyncButton
          danger
          icon={<DeleteOutlined />}
          shape="circle"
          onClick={deleteItem(data.id, format_id)}
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
