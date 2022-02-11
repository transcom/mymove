import classNames from 'classnames';

import styles from './index.module.scss';

let errorIfMissing = [];
let warnIfMissing = [];
let showWhenCollapsed = [];
let neverShow = [];
let shipment = null;

/*
  Set flags for a shipment based on fieldname. Null values will be set to empty arrays by default.

  Flag options:
    errorIfMissing,
    warnIfMissing,
    showWhenCollapsed,
    neverShow

  Pass in an array with string fieldnames to set any option e.g. ['counselorRemarks']
*/
export function setShipmentFlags(error, warn, show, never, ship) {
  errorIfMissing = error || [];
  warnIfMissing = warn || [];
  showWhenCollapsed = show || [];
  neverShow = never || [];
  shipment = ship;
}

/*
  Retrieve set flags for a shipment based on shipment fieldname.
  Flag options:
    errorIfMissing,
    warnIfMissing,
    showWhenCollapsed,
    neverShow
*/
export function getShipmentFlags(fieldname) {
  let alwaysShow = false;
  let classes = styles.row;
  // Hide row will override any always show that is set.
  let hideRow = false;

  if (errorIfMissing.includes(fieldname) && !shipment[fieldname]) {
    alwaysShow = true;
    classes = classNames(styles.row, styles.missingInfoError);
    return {
      alwaysShow,
      classes,
    };
  }
  if (warnIfMissing.includes(fieldname) && !shipment[fieldname]) {
    alwaysShow = true;
    classes = classNames(styles.row, styles.warning);
    return {
      alwaysShow,
      classes,
    };
  }
  if (showWhenCollapsed.includes(fieldname)) {
    alwaysShow = true;
  }

  if (neverShow.includes(fieldname)) {
    hideRow = true;
  }

  return {
    hideRow,
    alwaysShow,
    classes,
  };
}

export function getMissingOrDash(fieldName) {
  return errorIfMissing.includes(fieldName) ? 'Missing' : '—';
}
