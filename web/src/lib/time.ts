export type TimeComponents = ReturnType<typeof parseTime>;
export const parseTime = (input: number) => {
  const hours = Math.floor(input / 3600);
  const minutes = Math.floor((input % 3600) / 60);
  const seconds = Math.floor(input % 60);
  return { hours, minutes, seconds };
};

export const toTimestamp = (input: TimeComponents) => {
  const { hours, minutes, seconds } = input;
  return [
    `${hours}`.padStart(2, "0"),
    `${minutes}`.padStart(2, "0"),
    `${seconds}`.padStart(2, "0"),
  ].join(":");
};
