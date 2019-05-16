import moment from 'moment';

// Calculate the total days used by an array of SITs.
export function sitTotalDaysUsed(storageInTransits) {
  if (!storageInTransits) {
    return 0;
  }

  return storageInTransits.reduce((totalDays, storageInTransit) => {
    return totalDays + sitDaysUsed(storageInTransit);
  }, 0);
}

// Calculate the days used by a single SIT.
export function sitDaysUsed(storageInTransit) {
  if (!storageInTransit || !storageInTransit.actual_start_date) {
    return 0;
  }

  const today = moment(); // TODO: What about the time zone? These moments are all local time.

  const startDate = moment(storageInTransit.actual_start_date);

  let sitDays = 0;
  switch (storageInTransit.status) {
    case 'IN_SIT':
      sitDays = today.diff(startDate, 'days') + 1;
      break;
    case 'RELEASED':
    case 'DELIVERED':
      if (storageInTransit.out_date) {
        const outDate = moment(storageInTransit.out_date);
        sitDays = outDate.diff(startDate, 'days') + 1;
      }
      break;
    default:
      sitDays = 0;
  }

  return sitDays > 0 ? sitDays : 0;
}
