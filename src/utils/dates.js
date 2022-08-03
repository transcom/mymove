import moment from 'moment';

export default function numOfDaysBetweenDates(date1, date2) {
  const mDate1 = moment(date1);
  const mDate2 = moment(date2);

  return mDate2.diff(mDate1, 'days');
}
