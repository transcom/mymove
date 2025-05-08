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
    .slice(0, 10); // 10 instead of 14 ignores time to give us a general timeframe to account for processing delays
  return `${prefix}-${timestamp}`;
};

/**
 * @returns tomorrow as a formatted "dd mmm yyyy" date string in UTC
 */
export function getTomorrowUTC() {
  const tomorrow = new Date();
  tomorrow.setUTCDate(new Date().getUTCDate() + 1);

  const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];

  const day = String(tomorrow.getUTCDate()).padStart(2, '0');
  const month = months[tomorrow.getUTCMonth()];
  const year = tomorrow.getUTCFullYear();

  return `${day} ${month} ${year}`;
}

export default findOptionWithinOpenedDropdown;
