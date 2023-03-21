import { Button, ButtonProps, Tooltip } from "antd";
import { useState } from "react";

export interface AsyncButtonProps extends ButtonProps {
  onClick: () => Promise<void>;
  tooltip?: string;
  errorTooltip?: boolean;
}

export const AsyncButton = ({
  tooltip,
  errorTooltip,
  onClick,
  ...props
}: AsyncButtonProps) => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error>();
  const tooltipTitle = errorTooltip && error?.message ? error.message : tooltip;

  const clicked = async () => {
    setLoading(true);
    try {
      await onClick();
    } catch (error) {
      console.error(error);
      setError(error as Error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Tooltip title={tooltipTitle}>
      <Button loading={loading} danger={!!error} {...props} onClick={clicked} />
    </Tooltip>
  );
};
