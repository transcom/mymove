// maps int to int with ordinal 1 -> 1st, 2 -> 2nd, 3rd ...
export const intToOrdinal = n => {
  const s = ['th', 'st', 'nd', 'rd'];
  const v = n % 100;
  // eslint-disable-next-line security/detect-object-injection
  return n + (s[(v - 20) % 10] || s[v] || s[0]);
};
