import classNames from 'classnames';

let errorIfMissing = [];
let warnIfMissing = [];
let showWhenCollapsed = [];
let neverShow = [];
let flaggedItem = null;
let flagStyles = null;

/*
  Set CSS styles for the flags
*/
export function setFlagStyles(incomingStyles) {
  if (incomingStyles) {
    flagStyles = {
      row: incomingStyles.row || null,
      warning: incomingStyles.warning || null,
      missingInfoError: incomingStyles.missingInfoError || null,
    };
  } else {
    flagStyles = {};
  }
}

/*
  Set flags for a item based on fieldname. Null values will be set to empty arrays by default.

  Flag options:
    errorIfMissing,
    warnIfMissing,
    showWhenCollapsed,
    neverShow

  Pass in an array with string fieldnames to set any option e.g. ['counselorRemarks']
*/
export function setDisplayFlags(error, warn, show, never, item) {
  errorIfMissing = error || [];
  warnIfMissing = warn || [];
  showWhenCollapsed = show || [];
  neverShow = never || [];
  flaggedItem = item;
}

/*
  Retrieve set flags for a shipment based on shipment fieldname.
  Flag options:
    errorIfMissing,
    warnIfMissing,
    showWhenCollapsed,
    neverShow
*/
export function getDisplayFlags(fieldname) {
  if (!flagStyles) {
    return {};
  }

  let alwaysShow = false;
  let classes = flagStyles.row;
  // Hide row will override any always show that is set.
  let hideRow = false;

  if (errorIfMissing.includes(fieldname) && !flaggedItem[fieldname]) {
    alwaysShow = true;
    classes = classNames(flagStyles.row, flagStyles.missingInfoError);
    return {
      alwaysShow,
      classes,
    };
  }
  if (warnIfMissing.includes(fieldname) && !flaggedItem[fieldname]) {
    alwaysShow = true;
    classes = classNames(flagStyles.row, flagStyles.warning);
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
