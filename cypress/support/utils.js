import moment from 'moment';

// convert nth day of current month to date string YYYY-MM-DD
export const nthDayOfCurrentMonth = n => {
  const today = moment();
  const day = String('0' + n).slice(-2);
  return `${today.format('YYYY')}-${today.format('MM')}-${day}`;
};
