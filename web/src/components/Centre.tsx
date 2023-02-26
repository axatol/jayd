import { CSSProperties, PropsWithChildren } from "react";

const raisedStyle: CSSProperties = {
  padding: 8,
  borderRadius: 8,
};

export interface CentreProps {
  style?: CSSProperties;
  raised?: boolean;
}

export const Centre = (props: PropsWithChildren<CentreProps>) => (
  <div
    {...props}
    style={{
      height: "100%",
      width: "100%",
      display: "flex",
      justifyContent: "center",
      alignItems: "center",
      flexDirection: "column",
      ...(props.raised && raisedStyle),
      ...props.style,
    }}
  >
    {props.children}
  </div>
);
