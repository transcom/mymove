export default function formattedCustomerName(last, first, suffix = '', middle = '') {
  const lastFirst = [`${last}, ${first}`];
  const lastFirstMiddle = [lastFirst, middle].filter(Boolean).join(' ');

  return suffix.length > 0 ? `${lastFirstMiddle}, ${suffix}` : lastFirstMiddle;
}
