import { CloseOutlined, DownloadOutlined } from "@ant-design/icons";
import { Button, Card, Descriptions, Image, Select, Space } from "antd";
import { useState } from "react";

import { ExternalLink } from "./Links";
import { parseTime, toTimestamp } from "../lib/time";
import { VideoFormat, YoutubeInfoJSON } from "../lib/types";

const sortFormats = (a: VideoFormat, b: VideoFormat) => {
  // biggest resolution first
  const height = b.height - a.height;
  if (height !== 0) {
    return height;
  }

  // smallest filesize first
  const filesize = a.filesize - b.filesize;
  if (filesize !== 0) {
    return filesize;
  }

  const video = a.video_ext.localeCompare(b.video_ext);
  if (video !== 0) {
    return video;
  }

  return a.audio_ext.localeCompare(b.audio_ext);
};

const buildFormatOptions = (formats: VideoFormat[]) => {
  const video: { label: string; value: string }[] = [];
  const audio: { label: string; value: string }[] = [];

  formats
    .sort(sortFormats)
    .filter(
      (format) =>
        format.filesize > 0 &&
        (format.audio_ext !== "none" || format.video_ext !== "none"),
    )
    .forEach((format) => {
      const size = (format.filesize / 1000000).toFixed(2);

      if (format.video_ext !== "none") {
        video.push({
          value: format.format_id,
          label: `${format.height}p ${format.video_ext} (${size}MB)`,
        });
      }

      if (format.audio_ext !== "none" && format.video_ext === "none") {
        audio.push({
          value: format.format_id,
          label: `${format.audio_ext} (${size}MB)`,
        });
      }
    });

  return { video, audio };
};

export interface YoutubeMetadataCardProps {
  metadata: YoutubeInfoJSON;
  onConfirm: (videoId: string, formatId: string) => void;
  onReset?: () => void;
}

export const YoutubeMetadataCard = (props: YoutubeMetadataCardProps) => {
  const { metadata, onConfirm, onReset } = props;
  const [formatId, setFormatId] = useState<string>();
  const options = buildFormatOptions(metadata.formats);

  const selectFormat = (selected: string) => {
    setFormatId(selected);
  };

  const beginDownload = () => {
    if (formatId) {
      onConfirm(metadata.id, formatId);
    }
  };

  return (
    <Card
      style={{ cursor: "auto" }}
      hoverable
      title={
        <ExternalLink to={`https://youtube.com/watch?v=${metadata.id}`}>
          {metadata.title}
        </ExternalLink>
      }
      extra={
        <Space>
          <Select
            showSearch
            placeholder="Select a format"
            optionFilterProp="label"
            style={{ width: 300 }}
            onChange={selectFormat}
            options={[
              { label: "Video", options: options.video },
              { label: "Audio only", options: options.audio },
            ]}
          />

          <Button
            type="primary"
            icon={<DownloadOutlined />}
            disabled={!formatId}
            onClick={beginDownload}
            shape="round"
          >
            Begin download
          </Button>

          <Button onClick={onReset} icon={<CloseOutlined />} shape="circle" />
        </Space>
      }
    >
      <Card.Meta
        avatar={
          <Image
            src={metadata.thumbnail}
            alt={metadata.title}
            style={{ maxWidth: 200, maxHeight: 200 }}
          />
        }
        description={
          <Descriptions bordered size="small">
            <Descriptions.Item span={3} label="Duration">
              {toTimestamp(parseTime(metadata.duration))}
            </Descriptions.Item>

            <Descriptions.Item span={3} label="Channel">
              <ExternalLink to={`https://youtube.com/${metadata.uploader_id}`}>
                {metadata.uploader}
              </ExternalLink>
            </Descriptions.Item>

            <Descriptions.Item
              span={3}
              label="Description"
              style={{ whiteSpace: "pre-wrap" }}
            >
              {metadata.description}
            </Descriptions.Item>
          </Descriptions>
        }
      />
    </Card>
  );
};
