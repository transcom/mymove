export default function formattedCustomerName(last, first, suffix = '', middle = '') {
  if (suffix.length > 0 && middle.length > 0) {
    return `${last}, ${first} ${middle}, ${suffix},`;
  }
  if (suffix.length > 0 && middle.length === 0) {
    return `${last}, ${first}, ${suffix}`;
  }
  if (suffix.length === 0 && middle.length > 0) {
    return `${last}, ${first} ${middle}`;
  }
  return `${last}, ${first}`;
}
