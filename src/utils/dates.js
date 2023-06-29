import moment from 'moment';

import { SERVICE_ITEM_STATUS } from 'shared/constants';

export function numOfDaysBetweenDates(date1, date2) {
  const mDate1 = moment.utc(date1);
  const mDate2 = moment.utc(date2);

  return mDate2.diff(mDate1, 'days');
}

export function selectDateFieldByStatus(status) {
  let dateField;

  switch (status) {
    case SERVICE_ITEM_STATUS.SUBMITTED:
      dateField = 'createdAt';
      break;
    case SERVICE_ITEM_STATUS.APPROVED:
      dateField = 'approvedAt';
      break;
    case SERVICE_ITEM_STATUS.REJECTED:
      dateField = 'rejectedAt';
      break;
    default:
      dateField = 'createdAt';
  }

  return dateField;
}
