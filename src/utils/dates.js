export default function numOfDaysBetweenDates(date1, date2) {
  const jsDate1 = new Date(date1);
  const jsDate2 = new Date(date2);

  // To calculate the time difference of two dates
  const timeDiff = jsDate2.getTime() - jsDate1.getTime();

  // To calculate the no. of days between two dates
  return timeDiff / (1000 * 3600 * 24);
}
