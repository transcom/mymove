import classNames from 'classnames';

import flagStyles from './index.module.scss';

let errorIfMissing,
  warnIfMissing,
  showWhenCollapsed,
  neverShow,
  shipment,
  styles = null;

export function setShipmentFlags(error, warn, show, never, ship) {
  errorIfMissing = error;
  warnIfMissing = warn;
  showWhenCollapsed = show;
  neverShow = never;
  shipment = ship;

  return;
}

export function setFlagRowStyles(flagRowStyles) {
  styles = flagRowStyles;

  return;
}

export function getShipmentFlags(fieldname) {
  let alwaysShow = false;
  let classes = styles.row;
  // Hide row will override any always show that is set.
  let hideRow = false;

  if (errorIfMissing.includes(fieldname) && !shipment[fieldname]) {
    alwaysShow = true;
    classes = classNames(styles.row, flagStyles.missingInfoError);
    return {
      alwaysShow,
      classes,
    };
  }
  if (warnIfMissing.includes(fieldname) && !shipment[fieldname]) {
    alwaysShow = true;
    classes = classNames(styles.row, flagStyles.warning);
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
