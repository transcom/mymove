// @ts-ignore

const format = function (formattingString, ...args) {
  return formattingString.replace(/{(\d+)}/g, (match, index) => (args[index] !== undefined ? args[index] : match));
};

const trimOuterSymbols = function (stringValue, symbol) {
  return stringValue.replace(new RegExp(`^\\${symbol}+|\\${symbol}+$`, 'g'), '');
};

const filterYear = function (yearString) {
  return yearString.replace(/.*?(\b\d{4}\b).*?/g, '$1');
};

const trimUrlOperation = function (stringSegment) {
  return new Set(stringSegment).size === stringSegment.length;
};

const GetURLOrigin = function (stringUrl) {
  return new URL(stringUrl).origin;
};

const concatUrlSegments = (stringPathSegments) => {
  const trimSegments = (path) => (typeof path === 'string' ? trimOuterSymbols(path, trimOuterSymbols('/')) : []);
  const result = stringPathSegments.flatMap(trimSegments).join('/');
  return result;
};

export const stringHelpers = {
  format,
  trimOuterSymbols,
  filterYear,
  trimUrlOperation,
  GetURLOrigin,
  concatUrlSegments,
};

export const formatRelativeDate = (daysInterval) => {
  const relativeDate = new Date();
  relativeDate.setDate(relativeDate.getDate() + daysInterval);
  const formattedDay = relativeDate.toLocaleDateString(undefined, { day: '2-digit' });
  const formattedMonth = relativeDate.toLocaleDateString(undefined, {
    month: 'short',
  });
  const formattedYear = relativeDate.toLocaleDateString(undefined, {
    year: 'numeric',
  });
  const formattedDate = `${formattedDay} ${formattedMonth} ${formattedYear}`;

  return {
    relativeDate,
    formattedDate,
  };
};

export const formatNumericDate = (date) => {
  const formattedDay = date.toLocaleDateString(undefined, { day: '2-digit' });
  const formattedMonth = date.toLocaleDateString(undefined, {
    month: '2-digit',
  });
  const formattedYear = date.toLocaleDateString(undefined, {
    year: 'numeric',
  });

  return [formattedYear, formattedMonth, formattedDay].join('-');
};

export const textWithNoTrailingNumbers = (value) => RegExp(`${value}$`);

const DAY_PICKER_YEAR_CLASS = 'div.DayPicker-Caption';
const DAY_PICKER_PREVIOUS = '.DayPicker-NavButton--prev';
const DAY_PICKER_NEXT = '.DayPicker-NavButton--next';
const DAY_PICKER_DAYS = 'div.DayPicker-Week';

export const dateInputOperator = async (page, locator, dateValue) => {
  const clickDate = async (page, dayOfMonthNumber) => {
    const dayElementLocator = page.locator(DAY_PICKER_DAYS, { state: 'attached' });
    const filterOnlyCompleteMatch = new RegExp(`^${dayOfMonthNumber}$`);
    const targetDate = await dayElementLocator.getByText(filterOnlyCompleteMatch, { state: 'attached' });
    await targetDate.click();
  };

  /** calculate how many clicks to the correct month */
  const getMonthsDifference = (originDate, targetDate) =>
    targetDate.getMonth() - originDate.getMonth() + 12 * (targetDate.getFullYear() - originDate.getFullYear());

  const datesMatch = (dateA, dateB) =>
    dateA.getMonth() === dateB.getMonth() && dateA.getFullYear() === dateB.getFullYear();

  await locator.focus();
  const yearElement = await page.locator(DAY_PICKER_YEAR_CLASS, { state: 'attached' });
  const monthYearText = await yearElement.textContent();
  const firstDate = new Date(['1', monthYearText].join(' '));

  const startingPageIsTargetPage = datesMatch(dateValue, firstDate);

  if (startingPageIsTargetPage) {
    await clickDate(page, dateValue.getDate());
    return;
  }

  const lastDate = new Date(firstDate);
  lastDate.setMonth(1);
  lastDate.setHours(-24);

  /** The amount of months is how many times we click to to get to the target month */
  const months = getMonthsDifference(dateValue, firstDate);

  /** Set the direction going backwards or forwards based on month's sign. */
  const directionAction = months > 0 ? DAY_PICKER_PREVIOUS : DAY_PICKER_NEXT;

  let currentMonth = monthYearText;
  for (let i = 0; i < Math.abs(months); i++) {
    await page.locator(directionAction).click();
    let changedMonth = await page.locator(DAY_PICKER_YEAR_CLASS, { hasNotText: currentMonth, state: 'attached' });
    currentMonth = await changedMonth.textContent();
  }

  const firstDateOfCurrentMonth = new Date(['1', currentMonth].join(' '));
  const traversedPageIsTargetPage = datesMatch(dateValue, firstDateOfCurrentMonth);

  if (traversedPageIsTargetPage) {
    await clickDate(page, dateValue.getDate());
    return;
  }
};

export {};
