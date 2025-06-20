export const findOptionWithinOpenedDropdown = (page, searchTerm) => {
  return page.locator('[id^="react-select"][id*="listbox"]').locator(`[id*="option"]:has(:text("${searchTerm}"))`);
};

/**
 * @returns a filename prefix appended with timestamp
 *
 * @param {string} prefix // filename prefix to append timestamp with '-{10 digit ISO timestamp}'
 */
export const appendTimestampToFilenamePrefix = (prefix) => {
  const timestamp = new Date()
    .toISOString()
    .replace(/[-T:.Z]/g, '')
    .slice(0, 10); // take first 10 digits (YYYYMMDDHH) instead of 14 (YYYYMMDDHHmmss) to give us a general timeframe to account for processing delays
  return `${prefix}-${timestamp}`;
};

/**
 * @param {Date} date
 */
export function formatDate(date) {
  const day = date.toLocaleString('default', { day: '2-digit' });
  const month = date.toLocaleString('default', { month: 'short' });
  const year = date.toLocaleString('default', { year: 'numeric' });
  return `${day} ${month} ${year}`;
}

export function getFutureDate() {
  const tomorrow = new Date();
  tomorrow.setDate(tomorrow.getDate() + 1);
  return formatDate(tomorrow);
}

export default findOptionWithinOpenedDropdown;
