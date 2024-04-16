import moment from 'moment';

import { SERVICE_ITEM_STATUS } from 'shared/constants';

export function numOfDaysBetweenDates(date1, date2) {
  const mDate1 = moment.utc(date1);
  const mDate2 = moment.utc(date2);

  return mDate2.diff(mDate1, 'days');
}

export const selectDatePrefixByStatus = (status) => {
  switch (status) {
    case SERVICE_ITEM_STATUS.APPROVED:
      return 'Date approved';
    case SERVICE_ITEM_STATUS.REJECTED:
      return 'Date rejected';
    case 'Move Task Order Approved':
      return 'Date approved';
    case 'Move Task Order Rejected':
      return 'Date rejected';
    case SERVICE_ITEM_STATUS.SUBMITTED:
    default:
      return 'Date requested';
  }
};

export const selectDateFieldByStatus = (status) => {
  switch (status) {
    case SERVICE_ITEM_STATUS.APPROVED:
      return 'approvedAt';
    case SERVICE_ITEM_STATUS.REJECTED:
      return 'rejectedAt';
    case 'Move Task Order Approved':
      return 'approvedAt';
    case 'Move Task Order Rejected':
      return 'rejectedAt';
    case SERVICE_ITEM_STATUS.SUBMITTED:
    default:
      return 'createdAt';
  }
};
