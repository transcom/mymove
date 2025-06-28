import path from 'path';

import moment from 'moment';
import numeral from 'numeral';

import { ASSIGNMENT_IDS, ASSIGNMENT_NAMES } from 'constants/MoveHistory/officeUserAssignment';
import { DEPARTMENT_INDICATOR_OPTIONS } from 'constants/departmentIndicators';
import { SERVICE_MEMBER_AGENCY_LABELS } from 'content/serviceMemberAgencies';
import { ORDERS_TYPE_OPTIONS, ORDERS_TYPE_DETAILS_OPTIONS, ORDERS_TYPE, ORDERS_PAY_GRADE_TYPE } from 'constants/orders';
import { PAYMENT_REQUEST_STATUS_LABELS } from 'constants/paymentRequestStatus';
import { DEFAULT_EMPTY_VALUE } from 'shared/constants';

/**
 * Formats number into a dollar string. Eg. $1,234.12
 *
 * More info: http://numeraljs.com/
 * @param num
 * @returns {string}
 */
export function toDollarString(num) {
  return numeral(num).format('$0,0.00');
}

// Format a thousandth of an inch into an inch, e.g. 16700 -> 16.7
export function convertFromThousandthInchToInch(thousandthInch) {
  if (!Number.isFinite(thousandthInch)) {
    return null;
  }

  return thousandthInch / 1000;
}

// Format user-entered dimension into base dimension, e.g. 15.25 -> 15250
export function formatToThousandthInches(val) {
  return parseFloat(String(val).replace(',', '')) * 1000;
}

// Service Member Formatters

// Format a date in the MM-DD-YYYY format for use in the service member UI.
export function formatDateSM(date) {
  if (date) {
    return moment(date).format('MM/DD/YYYY');
  }
  return undefined;
}

// Format a date into the format required for submission as a date property in
// Swagger.
export function formatSwaggerDate(date) {
  if (date) {
    return moment(date).format('YYYY-MM-DD');
  }
  return '';
}

// Parse a date from the format used by Swagger into a Date object
export function parseSwaggerDate(dateString) {
  if (dateString) {
    return moment(dateString, 'YYYY-MM-DD').toDate();
  }
  return undefined;
}

// Format date for display of dates summaries
const formatDateForDateRange = (date, formatType) => {
  let format = '';
  switch (formatType) {
    case 'long':
      format = 'ddd, MMM DD';
      break;
    case 'condensed':
      format = 'MMM D';
      break;
    default:
      format = 'ddd, MMM DD';
  }
  if (date) {
    return moment(date).format(format);
  }
  return undefined;
};

export const displayDateRange = (dates, formatType = 'long') => {
  let span = '';
  let firstDate = '';
  if (dates.length > 1) {
    span = ` - ${formatDateForDateRange(dates[dates.length - 1], formatType)}`;
  }
  if (dates.length >= 1) {
    firstDate = formatDateForDateRange(dates[0], formatType);
  }
  return firstDate + span;
};

// Office Formatters

// Format a date and ignore any time values, e.g. 03-Jan-18
export function formatDate(date, inputFormat, outputFormat = 'DD-MMM-YY', locale = 'en', isStrict = false) {
  if (date) {
    return moment(date, inputFormat, locale, isStrict).format(outputFormat);
  }
  return undefined;
}

export function formatDateFromIso(date, outputFormat) {
  return formatDate(date, 'YYYY-MM-DDTHH:mm:ss.SSSZ', outputFormat);
}

// Format a date and include its time, e.g. 03-Jan-2018 21:23
export function formatDateTime(date) {
  if (date) {
    return moment(date).format('DD-MMM-YY HH:mm');
  }
  return undefined;
}

export function formatTimeAgo(date) {
  if (!date) return undefined;

  return moment(date)
    .fromNow()
    .replace('minute', 'min')
    .replace(/a min\s/, '1 min ');
}

// maps int to int with ordinal 1 -> 1st, 2 -> 2nd, 3rd ...
export const formatToOrdinal = (n) => {
  const s = ['th', 'st', 'nd', 'rd'];
  const v = n % 100;
  return n + (s[(v - 20) % 10] || s[v] || s[0]);
};

export const departmentIndicatorReadable = (departmentIndicator, missingText) => {
  if (!departmentIndicator) {
    return missingText;
  }
  return DEPARTMENT_INDICATOR_OPTIONS[`${departmentIndicator}`] || departmentIndicator;
};

export const serviceMemberAgencyLabel = (agency) => {
  return SERVICE_MEMBER_AGENCY_LABELS[`${agency}`] || agency;
};

export const ordersTypeReadable = (ordersType, missingText) => {
  if (!ordersType) {
    return missingText;
  }
  if (ordersType === 'SAFETY') {
    return 'Safety';
  }
  return ORDERS_TYPE_OPTIONS[`${ordersType}`] || ordersType;
};

export const ordersTypeDetailReadable = (ordersTypeDetail, missingText) => {
  if (!ordersTypeDetail) {
    return missingText;
  }
  return ORDERS_TYPE_DETAILS_OPTIONS[`${ordersTypeDetail}`] || ordersTypeDetail;
};

export const paymentRequestStatusReadable = (paymentRequestStatus) => {
  return PAYMENT_REQUEST_STATUS_LABELS[`${paymentRequestStatus}`] || paymentRequestStatus;
};

export const filenameFromPath = (filePath) => {
  return path.basename(filePath);
};

export const formatAddressShort = (address) => {
  const { city, state, postalCode } = address;
  return `${city}, ${state} ${postalCode}`;
};

export const formatPrimeAPIFullAddress = (address) => {
  const { streetAddress1, streetAddress2, streetAddress3, city, state, postalCode } = address;
  return `${streetAddress1}, ${streetAddress2}, ${streetAddress3}, ${city}, ${state} ${postalCode}`;
};

export const formatEvaluationReportShipmentAddress = (address) => {
  const { streetAddress1, city, state, postalCode } = address;
  if (!streetAddress1 || !city || !state) {
    return postalCode;
  }
  return `${streetAddress1}, ${city}, ${state} ${postalCode}`;
};

export const formatCustomerContactFullAddress = (address) => {
  let formattedAddress = '';
  if (address.streetAddress1) {
    formattedAddress += `${address.streetAddress1}`;
  }

  if (address.streetAddress2) {
    formattedAddress += `, ${address.streetAddress2}`;
  }

  if (address.streetAddress3) {
    formattedAddress += `, ${address.streetAddress3}`;
  }

  if (address.city) {
    formattedAddress += `, ${address.city}`;
  }

  if (address.state) {
    formattedAddress += `, ${address.state}`;
  }

  if (address.postalCode) {
    formattedAddress += ` ${address.postalCode}`;
  }

  if (formattedAddress[0] === ',') {
    formattedAddress = formattedAddress.substring(1);
  }

  formattedAddress = formattedAddress.trim();

  return formattedAddress;
};

export const formatMoveHistoryFullAddress = (address) => {
  let formattedAddress = '';
  if (address.street_address_1) {
    formattedAddress += `${address.street_address_1}`;
  }

  if (address.street_address_2) {
    formattedAddress += `, ${address.street_address_2}`;
  }

  if (address.city) {
    formattedAddress += `, ${address.city}`;
  }

  if (address.state) {
    formattedAddress += `, ${address.state}`;
  }

  if (address.postal_code) {
    formattedAddress += ` ${address.postal_code}`;
  }

  if (formattedAddress[0] === ',') {
    formattedAddress = formattedAddress.substring(1);
  }

  formattedAddress = formattedAddress.trim();

  return formattedAddress;
};

export const formatMoveHistoryFullAddressFromJSON = (address) => {
  return formatMoveHistoryFullAddress(JSON.parse(address));
};

export const formatMoveHistoryAgent = (agent) => {
  let agentLabel = '';
  if (agent.agent_type === 'RECEIVING_AGENT') {
    agentLabel = 'receiving_agent';
  } else if (agent.agent_type === 'RELEASING_AGENT') {
    agentLabel = 'releasing_agent';
  }

  let formattedAgent = '';

  if (agent.first_name) {
    formattedAgent += `${agent.first_name}`;
  }

  if (agent.last_name) {
    formattedAgent += ` ${agent.last_name}`;
  }

  if (agent.phone) {
    formattedAgent += `, ${agent.phone}`;
  }

  if (agent.email) {
    formattedAgent += `, ${agent.email}`;
  }

  if (formattedAgent[0] === ',') {
    formattedAgent = formattedAgent.substring(1);
  }

  formattedAgent = formattedAgent.trim();

  const formattedAgentValues = {};
  formattedAgentValues[agentLabel] = formattedAgent;

  return formattedAgentValues;
};

export const formatMoveHistoryMaxBillableWeight = (historyRecord) => {
  const { changedValues } = historyRecord;
  const newChangedValues = { ...changedValues };
  if (changedValues.authorized_weight) {
    newChangedValues.max_billable_weight = changedValues.authorized_weight;
    delete newChangedValues.authorized_weight;
  }
  return { ...historyRecord, changedValues: newChangedValues };
};

export const formatMoveHistoryGunSafe = (historyRecord) => {
  const { changedValues } = historyRecord;
  const newChangedValues = { ...changedValues };
  if (changedValues.gun_safe !== undefined) {
    newChangedValues.gun_safe_authorized = changedValues.gun_safe;
    delete newChangedValues.gun_safe;
  }
  if (changedValues.gun_safe_weight !== undefined) {
    newChangedValues.gun_safe_weight_allowance = changedValues.gun_safe_weight;
    delete newChangedValues.gun_safe_weight;
  }
  return { ...historyRecord, changedValues: newChangedValues };
};

export const dropdownInputOptions = (options) => {
  return Object.entries(options).map(([key, value]) => ({ key, value }));
};

export const formatPayGradeOptions = (payGrades) => {
  return payGrades.map((grade) => {
    return { key: grade.grade, value: grade.description };
  });
};

// Formats the numeric age input to a human readable string. Eg. 1.5 = 1 day, 2.5 = 2 days
export const formatAgeToDays = (age) => {
  if (age < 1) {
    return 'Less than 1 day';
  }
  if (age >= 1 && age < 2) {
    return '1 day';
  }
  return `${Math.floor(age)} days`;
};

/**
 * @name formatReviewShipmentWeightsDate
 * @description Format dates for Review Shipment Weights page (example: from
 * `25-Dec-23` to `Dec 25 2023`)
 * @param {string} date A string representing a date in the `DD-MMM-YY` format.
 * @return {string} A date formated into a string representing a date in the
 * following format: `Dec 25 2023`.
 */
export function formatReviewShipmentWeightsDate(date) {
  if (!date) return DEFAULT_EMPTY_VALUE;
  return moment.utc(date).format('MMM DD YYYY');
}
// Format dates for customer app (ex. 25 Dec 2020)
export function formatCustomerDate(date) {
  if (!date) return DEFAULT_EMPTY_VALUE;
  return moment.utc(date).format('DD MMM YYYY');
}
// Format dates for customer remarks in the office app (ex. 25 Dec 2020 8:00)
export function formatCustomerSupportRemarksDate(date) {
  return moment.utc(date).format('DD MMM YYYY HH:mm');
}

export function formatSignatureDate(date) {
  return moment.utc(date).format('YYYY-MM-DD');
}

// Translate boolean (true/false) into capitalized "Yes"/"No" string
export const formatYesNoMoveHistoryValue = (value) => {
  if (value === true) return 'Yes';
  if (value === false) return 'No';
  return null;
};

// Translate boolean (true/false) into "yes"/"no" string
export const formatYesNoInputValue = (value) => {
  if (value === true) return 'yes';
  if (value === false) return 'no';
  return null;
};

// Translate boolean (true/false) into "true"/"false" string
export const formatTrueFalseInputValue = (value) => {
  if (value === true) return 'true';
  if (value === false) return 'false';
  return null;
};

// Translate "yes"/"no" string into boolean (true/false)
export const formatYesNoAPIValue = (value) => {
  if (value === 'yes') return true;
  if (value === 'no') return false;
  return undefined;
};

// Translate weights from lbs to CWT
export const formatWeightCWTFromLbs = (value) => {
  return `${parseInt(value, 10) / 100} cwt`;
};

// Translate currency from millicents to dollars
export const formatDollarFromMillicents = (value, decimalPlaces = 2) => {
  return `$${(parseInt(value, 10) / 100000).toFixed(decimalPlaces)}`;
};

// Takes an whole number of day value and pluralizes with unit label
export const formatDaysInTransit = (days) => {
  if (days) {
    if (days === 1) {
      return '1 day';
    }
    return `${days} days`;
  }
  return '0 days';
};

export const formatDaysRemaining = (days) => {
  if (days) {
    if (days === 1) {
      return '1 day, ends';
    }
    if (days < 0) {
      return 'Expired, ended';
    }
    return `${days} days, ends`;
  }
  return '0 days, ends';
};

export const formatAddressForPrimeAPI = (address) => {
  return {
    streetAddress1: address.streetAddress1,
    streetAddress2: address.streetAddress2,
    streetAddress3: address.streetAddress3,
    city: address.city,
    county: address.county,
    state: address.state,
    postalCode: address.postalCode,
    usPostRegionCitiesID: address.usPostRegionCitiesID,
  };
};

export const formatExtraAddressForPrimeAPI = (address) => {
  const { streetAddress1, city, state, postalCode } = address;
  if (streetAddress1 === '' || city === '' || state === '' || postalCode === '') return null;
  return {
    streetAddress1: address.streetAddress1,
    streetAddress2: address.streetAddress2,
    streetAddress3: address.streetAddress3,
    city: address.city,
    county: address.county,
    state: address.state,
    postalCode: address.postalCode,
    usPostRegionCitiesID: address.usPostRegionCitiesID,
  };
};

const emptyAddress = {
  streetAddress1: '',
  streetAddress2: '',
  city: '',
  county: '',
  state: '',
  postalCode: '',
};

export function fromPrimeAPIAddressFormat(address) {
  if (!address) {
    return emptyAddress;
  }
  return {
    streetAddress1: address.streetAddress1,
    streetAddress2: address.streetAddress2,
    streetAddress3: address.streetAddress3,
    city: address.city,
    county: address.county,
    state: address.state,
    postalCode: address.postalCode,
    usPostRegionCitiesID: address.usPostRegionCitiesID,
  };
}

// Format a weight with lbs following, e.g. 4000 becomes 4,000 lbs
export function formatWeight(weight) {
  if (weight) {
    return `${weight.toLocaleString()} lbs`;
  }
  return '0 lbs';
}

// Format a UB allowance weight with lbs following, e.g. 4000 becomes 4,000 lbs
// if it's 0 or undefined, we'll send back a relevant string instead
export function formatUBAllowanceWeight(weight) {
  if (weight) {
    return `${weight.toLocaleString()} lbs`;
  }
  return 'your UB allowance';
}

export const formatDelimitedNumber = (number) => {
  // Fail-safe in case an actual number value is passed in
  const numberString = number.toString();
  return Number(numberString.replace(/,/g, ''));
};
/**
 * Depending on the order type, this will return:
 * Report by date (PERMANENT_CHANGE_OF_STATION)
 * Date of retirement (RETIREMENT)
 * Date of separation (SEPARATION)
 */
export const formatLabelReportByDate = (orderType) => {
  switch (orderType) {
    case 'RETIREMENT':
      return 'Date of retirement';
    case 'SEPARATION':
      return 'Date of separation';
    default:
      return 'Report by date';
  }
};

// Format a number of cents into a string, e.g. 12,345.67
export function formatCents(cents, minimumFractionDigits = 2, maximumFractionDigits = 2) {
  return (cents / 100).toLocaleString(undefined, { minimumFractionDigits, maximumFractionDigits });
}

export function formatCentsRange(min, max) {
  if (!Number.isFinite(min) || !Number.isFinite(max)) {
    return '';
  }

  return `$${formatCents(min)} - ${formatCents(max)}`;
}

// Formats a numeric value amount in the default locale with configurable options
// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Number/toLocaleString
export function formatAmount(amount, options = { minimumFractionDigits: 2, maximumFractionDigits: 2 }) {
  if (!Number.isFinite(amount)) {
    return '';
  }
  return amount.toLocaleString(undefined, options);
}

// Converts a cents value into whole dollars, rounding down.
export function convertCentsToWholeDollarsRoundedDown(cents) {
  return Math.floor(cents / 100);
}

// Converts a cents value into whole dollars, dropping the decimal precision without rounding e.g. 1234599 -> 12,345
export function formatCentsTruncateWhole(cents) {
  return formatAmount(convertCentsToWholeDollarsRoundedDown(cents), {
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  });
}

// Converts a uuid into a shortened ID that's suitable for displaying to users
function getUUIDFirstFive(uuid) {
  return uuid.substring(0, 5).toUpperCase();
}

export function formatQAReportID(uuid) {
  return `#QA-${getUUIDFirstFive(uuid)}`;
}

export function removeCommas(inputString) {
  // Use a regular expression to replace commas with an empty string
  return inputString.replace(/,/g, '');
}
export function formatEvaluationReportLocation(location) {
  switch (location) {
    case 'ORIGIN':
      return 'Origin';
    case 'DESTINATION':
      return 'Destination';
    case 'OTHER':
      return 'Other';
    default:
      return undefined;
  }
}

export function formatTimeUnitDays(days) {
  return `${days} days`;
}

export function formatDistanceUnitMiles(distance) {
  return `${distance} miles`;
}

export const constructSCOrderOconusFields = (values) => {
  const isOconus = values.originDutyLocation?.address?.isOconus || values.newDutyLocation?.address?.isOconus;
  const dependents = values.hasDependents;
  const isCivilianTDYMove =
    values.ordersType === ORDERS_TYPE.TEMPORARY_DUTY && values.grade === ORDERS_PAY_GRADE_TYPE.CIVILIAN_EMPLOYEE;
  // The `hasDependents` check within accompanied tour is due to
  // the dependents section being possible to conditionally render
  // and then un-render while still being OCONUS
  // The detailed comments make this nested ternary readable
  /* eslint-disable no-nested-ternary */
  return {
    accompaniedTour:
      isOconus && dependents
        ? // If OCONUS and dependents are present, fetch the value from the form.
          // Otherwise, default to false if OCONUS and dependents are not present
          dependents
          ? formatYesNoAPIValue(values.accompaniedTour) // Dependents are present
          : false // Dependents are not present
        : // If CONUS or no dependents, omit this field altogether
          null,
    dependentsUnderTwelve:
      isOconus && dependents
        ? // If OCONUS and dependents are present
          // then provide the number of dependents under 12. Default to 0 if not present
          Number(values.dependentsUnderTwelve) ?? 0
        : // If CONUS or no dependents, omit ths field altogether
          null,
    dependentsTwelveAndOver:
      isOconus && dependents
        ? // If OCONUS and dependents are present
          // then provide the number of dependents over 12. Default to 0 if not present
          Number(values.dependentsTwelveAndOver) ?? 0
        : // If CONUS or no dependents, omit this field altogether
          null,
    civilianTdyUbAllowance:
      isOconus && isCivilianTDYMove
        ? // If OCONUS
          // then provide the civilian TDY UB allowance. Default to 0 if not present
          Number(values.civilianTdyUbAllowance) ?? 0
        : // If CONUS, omit this field altogether
          null,
  };
};

export const formatServiceMemberNameToString = (serviceMember) => {
  let formattedUser = '';
  if (serviceMember.first_name && serviceMember.last_name) {
    formattedUser += `${serviceMember.first_name}`;
    formattedUser += ` ${serviceMember.last_name}`;
  } else {
    if (serviceMember.first_name) {
      formattedUser += `${serviceMember.first_name}`;
    }
    if (serviceMember.last_name) {
      formattedUser += `${serviceMember.last_name}`;
    }
  }
  return formattedUser;
};

export const formatAssignedOfficeUserFromContext = (historyRecord) => {
  const { changedValues, context, oldValues } = historyRecord;
  if (!context || context.length === 0) return {};

  const name = `${context[0].assigned_office_user_last_name}, ${context[0].assigned_office_user_first_name}`;
  const newValues = {};

  const assignOfficeUser = (key, assignedKey, reassignedKey) => {
    if (changedValues?.[key]) {
      newValues[oldValues[key] === null ? assignedKey : reassignedKey] = name;
    }
  };

  assignOfficeUser(
    ASSIGNMENT_IDS.SERVICES_COUNSELOR,
    ASSIGNMENT_NAMES.SERVICES_COUNSELOR.ASSIGNED,
    ASSIGNMENT_NAMES.SERVICES_COUNSELOR.RE_ASSIGNED,
  ); // counseling queue

  assignOfficeUser(
    ASSIGNMENT_IDS.CLOSEOUT_COUNSELOR,
    ASSIGNMENT_NAMES.CLOSEOUT_COUNSELOR.ASSIGNED,
    ASSIGNMENT_NAMES.CLOSEOUT_COUNSELOR.RE_ASSIGNED,
  ); // closeout queue

  assignOfficeUser(
    ASSIGNMENT_IDS.TASK_ORDERING_OFFICER,
    ASSIGNMENT_NAMES.TASK_ORDERING_OFFICER.ASSIGNED,
    ASSIGNMENT_NAMES.TASK_ORDERING_OFFICER.RE_ASSIGNED,
  ); // task order queue

  assignOfficeUser(
    ASSIGNMENT_IDS.TASK_INVOICING_OFFICER,
    ASSIGNMENT_NAMES.TASK_INVOICING_OFFICER.ASSIGNED,
    ASSIGNMENT_NAMES.TASK_INVOICING_OFFICER.RE_ASSIGNED,
  ); // payment request queue

  assignOfficeUser(
    ASSIGNMENT_IDS.TASK_ORDERING_OFFICER_DESTINATION,
    ASSIGNMENT_NAMES.DESTINATION_TASK_ORDERING_OFFICER.ASSIGNED,
    ASSIGNMENT_NAMES.DESTINATION_TASK_ORDERING_OFFICER.RE_ASSIGNED,
  ); // destination request queue
  return newValues;
};

export const userName = (user) => {
  let formattedUser = '';
  if (user.firstName && user.lastName) {
    formattedUser += `${user.lastName}, `;
    formattedUser += ` ${user.firstName}`;
  } else {
    if (user.firstName) {
      formattedUser += ` ${user.firstName}`;
    }
    if (user.lastName) {
      formattedUser += ` ${user.lastName}`;
    }
  }
  return formattedUser;
};

/**
 * @description Converts a string to title case (capitalizes the first letter of each word)
 * @param {string} str - The input string to format.
 * @returns {string} - the formatted string in the title case.
 * */
export function toTitleCase(str) {
  if (!str) return '';
  return str
    .toLowerCase()
    .split(' ')
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
    .join(' ');
}

/**
 * @description This function is used to format the port in the
 * ShipmentAddresses component.
 * It displays only the port code, port name, city, state and zip code.
 * */
export function formatPortInfo(port) {
  if (port) {
    const formattedCity = toTitleCase(port.city);
    const formattedState = toTitleCase(port.state);
    return `${port.portCode} - ${port.portName}\n${formattedCity}, ${formattedState} ${port.zip}`;
  }
  return '-';
}

export function formatFullName(firstName, middleName, lastName) {
  return [firstName, middleName, lastName].filter(Boolean).join(' ');
}

export const calculateTotal = (sectionInfo) => {
  let total = 0;

  if (sectionInfo?.haulPrice) total += sectionInfo.haulPrice;
  if (sectionInfo?.haulFSC) total += sectionInfo.haulFSC;
  if (sectionInfo?.packPrice) total += sectionInfo.packPrice;
  if (sectionInfo?.unpackPrice) total += sectionInfo.unpackPrice;
  if (sectionInfo?.dop) total += sectionInfo.dop;
  if (sectionInfo?.ddp) total += sectionInfo.ddp;
  if (sectionInfo?.intlPackingPrice) total += sectionInfo.intlPackingPrice;
  if (sectionInfo?.intlUnpackPrice) total += sectionInfo.intlUnpackPrice;
  if (sectionInfo?.intlLinehaulPrice) total += sectionInfo.intlLinehaulPrice;
  if (sectionInfo?.sitReimbursement) total += sectionInfo.sitReimbursement;

  return formatCents(total);
};
