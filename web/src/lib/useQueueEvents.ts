import { useEffect, useReducer, useRef } from "react";
import useWebSocket from "react-use-websocket";

import { QueueEvent, QueueItem } from "./types";
import { config } from "../config";

const baseUrl = config.api.baseUrl.replace(/^http/, "ws");

const reducer = (
  queue: QueueItem[],
  event: QueueEvent | { action: "set"; items: QueueItem[] },
) => {
  if (event.action === "set") {
    return event.items;
  }

  switch (event.action) {
    case "added":
      return [...queue, event.item];
    case "completed":
      return queue.map((item) =>
        item.id === event.item.id ? { ...item, completed: true } : item,
      );
    case "failed":
      return queue.map((item) =>
        item.id === event.item.id ? { ...item, failed: true } : item,
      );
    case "removed":
      return queue.filter((item) => item.id !== event.item.id);

    default:
      return queue;
  }
};

export const useQueueEvents = () => {
  const mounted = useRef(true);
  const [items, dispatch] = useReducer(reducer, []);
  const { lastJsonMessage } = useWebSocket(`${baseUrl}/ws/queue`, {
    onError: (e) => console.error(e),
    shouldReconnect: () => mounted.current,
  });

  useEffect(() => {
    console.log("useEffect", lastJsonMessage);

    if (!lastJsonMessage) {
      return;
    }

    const event: QueueEvent = lastJsonMessage as any;
    dispatch(event);
  }, [lastJsonMessage]);

  const set = (items: QueueItem[]) => {
    dispatch({ action: "set", items });
  };

  const append = (item: QueueItem) => {
    dispatch({ action: "added", item });
  };

  useEffect(() => {
    return () => {
      mounted.current = false;
    };
  });

  return { items, set, append, dispatch };
};
