import {
  PropsWithChildren,
  createContext,
  useCallback,
  useContext,
  useMemo,
  useRef,
  useState,
} from "react";

import { useAPI } from "./api";
import { QueueItem } from "./types";

type Interval = ReturnType<typeof setTimeout>;

const getPending = (items: QueueItem[]) =>
  items.reduce((pending, item) => pending || !item.completed_at, false);

interface QueueContextValue {
  items: QueueItem[];
  hasPending: boolean;
  refresh: () => Promise<QueueItem[]>;
  poll: () => void;
}

const QueueContext = createContext<QueueContextValue | undefined>(undefined);

export const useQueue = () => {
  const value = useContext(QueueContext);
  if (!value) {
    throw new Error("QueueContext consumer must have matching provider");
  }

  return value;
};

export const QueueProvider = (props: PropsWithChildren<{ delay?: number }>) => {
  const interval = useRef<Interval>();
  const api = useAPI();
  const [items, setItems] = useState<QueueItem[]>([]);

  const refresh = useCallback(async () => {
    const response = await api.getQueue();
    setItems(response.data);
    return response.data;
  }, []);

  const poll = useCallback(async () => {
    const newItems = await refresh();
    if (getPending(newItems)) {
      interval.current = setTimeout(poll, props.delay ?? 1000);
    }
  }, [props.delay]);

  const hasPending = useMemo(() => getPending(items), [items]);

  return (
    <QueueContext.Provider value={{ items, hasPending, refresh, poll }}>
      {props.children}
    </QueueContext.Provider>
  );
};
