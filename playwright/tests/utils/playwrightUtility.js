export const findOptionWithinOpenedDropdown = (page, searchTerm) => {
  return page.locator('[id^="react-select"][id*="listbox"]').locator(`[id*="option"]:has(:text("${searchTerm}"))`);
};

export default findOptionWithinOpenedDropdown;
