// @ts-ignore

const format = function (formattingString, ...args) {
  return formattingString.replace(/{(\d+)}/g, (match, index) => (args[index] !== undefined ? args[index] : match));
};

const trimOuterSymbols = function(stringValue, symbol) {
  return stringValue.replace(new RegExp(`^\\${symbol}+|\\${symbol}+$`, 'g'), '');
}

const filterYear = function(yearString) {
  return yearString.replace(/.*?(\b\d{4}\b).*?/g, '$1');
}

const trimUrlOperation = function(stringSegment) {
  return new Set(stringSegment).size === stringSegment.length;
};

const GetURLOrigin = function(stringUrl) {
  return new URL(stringUrl).origin;
}

const concatUrlSegments = (stringPathSegments) => {
  const trimSegments = (path) => typeof path === 'string' ? trimOuterSymbols(path, trimOuterSymbols('/')) : [];
  const result = stringPathSegments.flatMap(trimSegments).join('/');
  return result;
}

export const stringHelpers = {
  format,
  trimOuterSymbols,
  filterYear,
  trimUrlOperation,
  GetURLOrigin,
  concatUrlSegments,
}

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

const DAY_PICKER_PARENT_CLASS = 'div.DayPicker-wrapper';
const DAY_PICKER_YEAR_CLASS = 'div.DayPicker-Caption';
const DAY_PICKER_PREVIOUS = '.DayPicker-NavButton--prev';
const DAY_PICKER_NEXT = '.DayPicker-NavButton--next';
const DAY_PICKER_DAYS = 'div.DayPicker-Week'

export const dateInputOperator = async (page, locator, dateValue) => {
  await locator.focus();
  const navPrevious = await page.locator(DAY_PICKER_PREVIOUS, { state: 'attached'});
  const navNext = await page.locator(DAY_PICKER_NEXT, { state: 'attached'});

  const yearElement = await page.locator(DAY_PICKER_YEAR_CLASS, {  state: 'attached'});
  const monthYearText = await yearElement.textContent();
  const startYear = filterYear(monthYearText);

  const attemptThisMonth = page.locator(DAY_PICKER_DAYS);

  // const dayElementsQuery = await dayElementLocator.getByText(/\d/, { state: 'attached' }).all();

  const firstDate = new Date(['1', monthYearText].join(' '));

  const compareDate = dateValue;

  let matching = {
    month: compareDate.getMonth() === firstDate.getMonth(),
    year: compareDate.getYear() === firstDate.getYear(),
  }

  let isOnTargetDate = matching.month && matching.year;

  if(isOnTargetDate){
    const dayElementLocator = page.locator(DAY_PICKER_DAYS, { state: 'attached' });
    const targetDate = await dayElementLocator.getByText(compareDate.getDate(), { state: 'attached' });
    await targetDate.click();
    return;
  }

  const lastDate = new Date(firstDate);
  lastDate.setMonth(1);
  lastDate.setHours(-24);

  //calculate how many clicks to the correct month

  const getDateDifference = (originDate, targetDate) => 
  targetDate.getMonth() - originDate.getMonth() + 
    (12 * (targetDate.getFullYear() - originDate.getFullYear()));
  
  /** The amount of months is how many times we click to to get to the target month */
  const months = getDateDifference(dateValue, firstDate)
  const directionAction = months > 0 ? DAY_PICKER_PREVIOUS : DAY_PICKER_NEXT;

  let currentMonth = monthYearText;
  for(let i = 0; i < Math.abs(months); i++){
    await page.locator(directionAction).click();
    let changedMonth = await page.locator(DAY_PICKER_YEAR_CLASS, { hasNotText: currentMonth, state: 'attached'});
    currentMonth = await changedMonth.textContent();
  }
  
  const firstDateOfCurrentMonth = new Date(['1', currentMonth].join(' '));

  matching = {
    month: compareDate.getMonth() === firstDateOfCurrentMonth.getMonth(),
    year: compareDate.getYear() === firstDateOfCurrentMonth.getYear(),
  }

  isOnTargetDate = matching.month && matching.year;

  if(isOnTargetDate){
    const dayElementLocator = locator.locator(DAY_PICKER_DAYS, { state: 'attached' });
    const targetDate = await dayElementLocator.getByText(compareDate.getDate(), { state: 'attached' });
    await targetDate.click();
    return;
  }
}

export {};
