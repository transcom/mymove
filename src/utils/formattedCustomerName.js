export default function formattedCustomerName(last, first, suffix = '', middle = '') {
  if (suffix.length > 0 && middle.length > 0) {
    return `${last} ${suffix}, ${first} ${middle}`;
  }
  if (suffix.length > 0 && middle.length === 0) {
    return `${last} ${suffix}, ${first}`;
  }
  if (suffix.length === 0 && middle.length > 0) {
    return `${last}, ${first} ${middle}`;
  }
  return `${last}, ${first}`;
}
