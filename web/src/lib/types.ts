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
  title: string;
  formats: VideoFormat[];
  thumbnail: string;
  description: string;
  uploader: string;
  uploader_id: string;
  duration: number;
  duration_string: string;
}

export interface QueueItemWithoutFormat {
  is_completed: boolean;
  is_failed: boolean;
  selected_format_id: string;
  data: YoutubeInfoJSON;
}

export interface QueueItem extends QueueItemWithoutFormat {
  format?: VideoFormat;
}
