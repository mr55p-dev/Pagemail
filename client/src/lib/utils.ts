export function dateToString(date: Date) {
  return `${date.getDate().toString().padStart(2, "0")}-${
    date.getMonth().toString().padStart(2, "0") + 1
  }-${date.getFullYear().toString().padStart(2, "0")}`;
}
