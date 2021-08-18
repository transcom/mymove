export default function returnLowestValue(val1, val2) {
  if (val1 && val2) {
    return Math.min(val1, val2);
  }
  if (val1) {
    return val1;
  }
  if (val2) {
    return val2;
  }
  return null;
}
