export interface APIResponse<T> {
  message: string;
  data: T;
}

export interface YoutubeVideoMetadata {
  id: string;
  channel_id: string;
  title: string;
  description: string;
  thumnail_url: string;
  channel_title: string;
  duration: number;
}

export interface VideoFormat {
  filesize: number;
  format_id: string;
  width: number;
  height: number;
  fps: number;
  audio_ext: string;
  video_ext: string;
}

export interface YoutubeInfoJSON {
  id: string;
  format_id: string;
  title: string;
  formats: VideoFormat[];
  thumbnail: string;
  description: string;
  uploader: string;
  uploader_id: string;
  duration: number;
  duration_string: string;
  ext: string;
  filename: string;
}

export interface QueueItem {
  id: string;
  info: YoutubeInfoJSON;
  added_at: number;
  failed_at: number;
  completed_at: number;
}

export interface QueueEvent {
  action: "added" | "completed" | "failed" | "removed";
  item: QueueItem;
}
