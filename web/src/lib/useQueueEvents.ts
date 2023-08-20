import { useEffect, useReducer, useRef } from "react";
import useWebSocket from "react-use-websocket";

import { QueueEvent, QueueItem } from "./types";

const reducer = (
  queue: QueueItem[],
  event: QueueEvent | { action: "set"; items: QueueItem[] },
) => {
  switch (event.action) {
    case "set":
      return event.items;

    case "added":
      if (queue.find((item) => item.id === event.item.id)) {
        return queue.map((item) =>
          item.id === event.item.id ? event.item : item,
        );
      }

      return [...queue, event.item];

    case "completed":
    case "failed":
      return queue.map((item) =>
        item.id === event.item.id ? event.item : item,
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
  const wsUrl = `${window.location.origin.replace(/^http/, "ws")}/ws/queue`;
  const { lastJsonMessage } = useWebSocket(wsUrl, {
    onError: console.error,
    shouldReconnect: () => mounted.current,
  });

  // TODO ping

  useEffect(() => {
    if (import.meta.env.DEV) {
      console.log("useEffect", { wsUrl, lastJsonMessage });
    }

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
