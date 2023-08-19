import { useRef, useState } from "react";

import { useAPI } from "./api";
import { QueueItem } from "./types";

export const useQueuePoll = () => {
  const timeout = useRef<ReturnType<typeof setTimeout>>();
  const api = useAPI();
  const [polling, setPolling] = useState(false);
  const [pending, setPending] = useState(false);
  const [queue, setQueue] = useState<QueueItem[]>([]);

  const hasPending = () =>
    queue.reduce((pending, item) => pending || !item.completed, false);

  const refreshQueue = async (poll = false) => {
    try {
      setPending(true);
      const items = await api.getQueue();
      setQueue(items.data);
      if (poll) {
        scheduleNext();
      }
    } catch (error) {
      console.error("failed to poll the queue", error);
    } finally {
      setPending(false);
    }
  };

  const scheduleNext = (interval = 1000) => {
    if (!polling || timeout.current) {
      return;
    }

    timeout.current = setTimeout(refreshQueue, interval);
  };

  const pollIfPending = () => {
    if (hasPending()) {
      poll();
    }
  };

  const poll = () => {
    setPolling(true);
    refreshQueue(true);
  };

  const stop = () => {
    setPolling(false);
    clearTimeout(timeout.current);
    timeout.current = undefined;
  };

  return {
    pending,
    queue,
    refreshQueue,
    pollIfPending,
    poll,
    stop,
  };
};
