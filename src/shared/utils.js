import React from 'react';
import { get, includes, find, mapValues, capitalize } from 'lodash';
import moment from 'moment';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faClock from '@fortawesome/fontawesome-free-solid/faClock';
import faCheck from '@fortawesome/fontawesome-free-solid/faCheck';
import faExclamationCircle from '@fortawesome/fontawesome-free-solid/faExclamationCircle';
import './shared.css';

export const swaggerDateFormat = 'YYYY-MM-DD';
export const defaultDateFormat = 'M/D/YYYY';

export const no_op = () => undefined;
export const no_op_action = () => {
  return function(dispatch) {
    dispatch({
      type: 'NO_OP_TYPE',
      item: null,
    });
  };
};

// Turn an array into an object with .reduce()
export const objFromArray = array =>
  array.reduce((accumulator, current) => {
    accumulator[current.id] = current;
    return accumulator;
  }, {});

export const upsert = (arr, newValue) => {
  const index = arr.findIndex(obj => obj.id === newValue.id);
  if (index !== -1) {
    arr.splice(index, 1, newValue);
  } else {
    arr.push(newValue);
  }
};

export function fetchActive(foos) {
  return find(foos, i => includes(['DRAFT', 'SUBMITTED', 'APPROVED', 'PAYMENT_REQUESTED'], get(i, 'status'))) || null;
}

export function fetchActiveShipment(shipments) {
  return (
    find(shipments, i =>
      includes(
        // For now, this include all statuses, but this may be re-evaluated in the future.
        ['DRAFT', 'SUBMITTED', 'AWARDED', 'ACCEPTED', 'APPROVED', 'IN_TRANSIT', 'DELIVERED', 'COMPLETED'],
        get(i, 'status'),
      ),
    ) || null
  );
}

export function formatDateString(dateString) {
  let parsed = moment(dateString, defaultDateFormat);
  if (parsed.isValid()) {
    return parsed.format(swaggerDateFormat);
  }

  return dateString;
}

// Formats payload values according to Swagger spec type
export function formatPayload(payload, def) {
  return mapValues(payload, (val, key) => {
    const prop = get(def.properties, key);
    const propType = get(prop, 'type');
    const propFormat = get(prop, 'format');

    if (propType === 'string') {
      if (propFormat === 'date') {
        return formatDateString(val);
      }
    }

    return val;
  });
}

export const convertDollarsToCents = dollars => Math.round(parseFloat(String(dollars).replace(',', '')) * 100);

export function renderStatusIcon(status) {
  if (!status) {
    return;
  }
  if (status === 'AWAITING_REVIEW' || status === 'DRAFT' || status === 'SUBMITTED') {
    return <FontAwesomeIcon className="icon approval-waiting" icon={faClock} />;
  }
  if (status === 'OK') {
    return <FontAwesomeIcon className="icon approval-ready" icon={faCheck} />;
  }
  if (status === 'APPROVED' || status === 'INVOICED') {
    return <FontAwesomeIcon className="icon approved" icon={faCheck} />;
  }
  if (status === 'HAS_ISSUE') {
    return <FontAwesomeIcon className="icon approval-problem" icon={faExclamationCircle} />;
  }
}

export function snakeCaseToCapitals(str) {
  return str
    .split('_')
    .map(word => capitalize(word))
    .join(' ');
}

export function humanReadableError(errors) {
  return Object.entries(errors)
    .map(error => `${snakeCaseToCapitals(error[0])} ${error[1]}`)
    .join('/n');
}
