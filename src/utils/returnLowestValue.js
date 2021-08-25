export default function returnLowestValue(val1, val2) {
  if (val1 != null && val2 != null) {
    return Math.min(val1, val2);
  }

  return val1 != null ? val1 : val2;
}
