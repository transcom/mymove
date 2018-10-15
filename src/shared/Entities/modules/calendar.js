import { get } from 'lodash';
import { swaggerRequest } from 'shared/Swagger/request';
import { getClient } from 'shared/Swagger/api';
import { parseSwaggerDate } from 'shared/formatters';

export function getAvailableMoveDates(label, startDate) {
  return swaggerRequest(getClient, 'calendar.showAvailableMoveDates', { startDate }, { label });
}

export function selectAvailableMoveDates(state, id) {
  if (!id) {
    return null;
  }

  const availableMoveDates = get(state, 'entities.availableMoveDates.' + id);
  if (!availableMoveDates) {
    return null;
  }

  const availableLength = availableMoveDates.available.length;
  const convertedAvailable = new Array(availableLength);
  let minDate, maxDate;
  for (let i = 0; i < availableLength; i++) {
    const convertedDate = parseSwaggerDate(availableMoveDates.available[i]); // eslint-disable-line security/detect-object-injection
    if (!minDate || convertedDate < minDate) {
      minDate = convertedDate;
    }
    if (!maxDate || convertedDate > maxDate) {
      maxDate = convertedDate;
    }
    convertedAvailable[i] = convertedDate; // eslint-disable-line security/detect-object-injection
  }

  return {
    startDate: parseSwaggerDate(availableMoveDates.start_date),
    minDate: minDate,
    maxDate: maxDate,
    available: convertedAvailable,
  };
}
