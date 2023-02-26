import { useAuth0 } from "@auth0/auth0-react";
import axios, { AxiosError } from "axios";
import { useEffect } from "react";

import {
  APIResponse,
  QueueItem,
  QueueItemWithoutFormat,
  YoutubeInfoJSON,
} from "./types";
import { config } from "../config";

const api = axios.create({ baseURL: config.api.baseUrl });
let tokenInterceptorId: number | undefined;
let sanitiseInterceptorId: number | undefined;

const sanitise = (error: AxiosError) => {
  const data = error.toJSON();
  const sanitised = JSON.stringify(data).replace(
    /"Authorization":(\s?)".*?"/,
    '"Authorization":$1"***"',
  );
  return JSON.parse(sanitised);
};

export const useAPI = () => {
  const { getAccessTokenSilently } = useAuth0();

  useEffect(() => {
    if (tokenInterceptorId === undefined) {
      tokenInterceptorId = api.interceptors.request.use(async (request) => {
        const token = await getAccessTokenSilently();
        request.headers.set("Authorization", `Bearer ${token}`);
        return request;
      });
    }

    if (sanitiseInterceptorId === undefined) {
      sanitiseInterceptorId = api.interceptors.response.use(
        undefined,
        (error) => {
          if (axios.isAxiosError(error)) {
            return Promise.reject(sanitise(error));
          }
          return Promise.reject(error);
        },
      );
    }
  }, [getAccessTokenSilently]);

  const getMetadata = (target: string) =>
    api
      .get<APIResponse<YoutubeInfoJSON>>("/api/youtube/metadata", {
        params: { target },
      })
      .then((result) => result.data.data);

  const beginDownload = (target: string, format: string) =>
    api
      .post("/api/youtube", null, { params: { target, format } })
      .then((result) => result.data);

  const getQueueItem = (target: string) =>
    api
      .get<APIResponse<QueueItemWithoutFormat[]>>("/api/queue", {
        params: { target },
      })
      .then((result) => result.data.data.map(mapItemFormat));

  const getQueue = () =>
    api
      .get<APIResponse<QueueItemWithoutFormat[]>>("/api/queue")
      .then((result) => result.data.data.map(mapItemFormat));

  const staticFile = (file: string) => api.get(`/static/${file}`);

  return {
    api,
    getMetadata,
    beginDownload,
    getQueue,
    getQueueItem,
    staticFile,
  };
};

const mapItemFormat = (item: QueueItemWithoutFormat): QueueItem => {
  const selected = item.data.formats.find(
    (format) => format.format_id === item.selected_format_id,
  );

  return { ...item, format: selected };
};
