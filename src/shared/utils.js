import React from 'react';
import { capitalize, find, get, includes, mapValues } from 'lodash';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import { formatDateForSwagger } from './dates';

export const no_op = () => undefined;
export const no_op_action = () => {
  return function (dispatch) {
    dispatch({
      type: 'NO_OP_TYPE',
      item: null,
    });
  };
};

// Turn an array into an object with .reduce()
export const objFromArray = (array) =>
  array.reduce((accumulator, current) => {
    accumulator[current.id] = current;
    return accumulator;
  }, {});

export const upsert = (arr, newValue) => {
  const index = arr.findIndex((obj) => obj.id === newValue.id);
  if (index !== -1) {
    arr.splice(index, 1, newValue);
  } else {
    arr.push(newValue);
  }
};

export function fetchActive(foos) {
  return find(foos, (i) => includes(['DRAFT', 'SUBMITTED', 'APPROVED', 'PAYMENT_REQUESTED'], get(i, 'status'))) || null;
}

export function fetchActivePPM(foos) {
  return (
    find(foos, (i) =>
      includes(['DRAFT', 'SUBMITTED', 'APPROVED', 'PAYMENT_REQUESTED', 'COMPLETED'], get(i, 'status')),
    ) || null
  );
}

export function fetchActiveShipment(shipments) {
  return (
    find(shipments, (i) =>
      includes(
        // For now, this include all statuses, but this may be re-evaluated in the future.
        ['DRAFT', 'SUBMITTED', 'AWARDED', 'ACCEPTED', 'APPROVED', 'IN_TRANSIT', 'DELIVERED'],
        get(i, 'status'),
      ),
    ) || null
  );
}

// Formats payload values according to Swagger spec type
export function formatPayload(payload, def) {
  return mapValues(payload, (val, key) => {
    const prop = get(def.properties, key);
    const propType = get(prop, 'type');
    const propFormat = get(prop, 'format');

    if (propType === 'string') {
      if (propFormat === 'date') {
        return formatDateForSwagger(val);
      }
    }

    return val;
  });
}

export const convertDollarsToCents = (dollars) => {
  if (!dollars && dollars !== 0) {
    return undefined;
  }

  return Math.round(parseFloat(String(dollars).replace(/,/g, '')) * 100);
};

export function renderStatusIcon(status) {
  if (!status) {
    return undefined;
  }
  if (status === 'AWAITING_REVIEW' || status === 'DRAFT' || status === 'SUBMITTED') {
    return <FontAwesomeIcon className="icon approval-waiting" icon="clock" />;
  }
  if (status === 'OK') {
    return <FontAwesomeIcon className="icon approval-ready" icon="check" />;
  }
  if (status === 'APPROVED' || status === 'INVOICED' || status === 'CONDITIONALLY_APPROVED') {
    return <FontAwesomeIcon className="icon approved" icon="check" />;
  }
  if (status === 'HAS_ISSUE') {
    return <FontAwesomeIcon className="icon approval-problem" icon="exclamation-circle" />;
  }
  return undefined;
}

export function snakeCaseToCapitals(str) {
  return str
    .split('_')
    .map((word) => capitalize(word))
    .join(' ');
}

export function humanReadableError(errors) {
  return Object.entries(errors)
    .map((error) => `${snakeCaseToCapitals(error[0])} ${error[1]}`)
    .join('/n');
}

export function detectIE11() {
  let sAgent = window.navigator.userAgent;
  let Idx = sAgent.indexOf('Trident');
  if (Idx > -1) {
    return true;
  }
  return false;
}

export function detectFirefox() {
  if (typeof InstallTrigger !== 'undefined') {
    return true;
  }
  return false;
}

export function openLinkInNewWindow(url, windowName, window, relativeSize) {
  // eslint-disable-next-line security/detect-non-literal-fs-filename
  window
    .open(
      url,
      windowName,
      `resizable,scrollbars,status,noopener=true,noreferrer=true,width=${window.outerWidth * relativeSize},height=${
        window.outerHeight * relativeSize
      }`,
    )
    .focus(); // required in IE to put re-used window on top
  return false;
}

// Sort ascending by objects with string iso timestamps
export function dateSort(field, direction) {
  if (direction === 'desc') {
    return (a, b) => {
      return Date.parse(b[`${field}`]) - Date.parse(a[`${field}`]);
    };
  } else {
    return (a, b) => {
      return Date.parse(a[`${field}`]) - Date.parse(b[`${field}`]);
    };
  }
}

export function isValidWeight(weight) {
  if (weight !== 'undefined' && weight && weight > 0) {
    return true;
  }
  return false;
}

export function isEmpty(obj) {
  let empty = true;
  Object.keys(obj).forEach((key) => {
    if (typeof obj[key] !== 'undefined' && obj[key]) {
      empty = false;
    }
  });
  return empty;
}
