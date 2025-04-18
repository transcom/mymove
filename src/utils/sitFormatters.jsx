import React from 'react';
import moment from 'moment';

import styles from './sitFormatters.module.scss';

import { DatePickerInput } from 'components/form/fields';
import { DEFAULT_EMPTY_VALUE } from 'shared/constants';
import { formatDateForDatePicker, swaggerDateFormat } from 'shared/dates';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import { LOCATION_TYPES } from 'types/sitStatusShape';
import DataTable from 'components/DataTable';

// ****************
// ****************
// formatting dates
// ****************
// ****************

// takes in a date, returns it in DD MMM YYYY format
export function formatSITDepartureDate(date) {
  return formatDateForDatePicker(date) || DEFAULT_EMPTY_VALUE;
}

// takes in a date, returns it in DD MMM YYYY format
export function formatSITEntryDate(date) {
  return formatDateForDatePicker(date) || DEFAULT_EMPTY_VALUE;
}

// takes in the sitStatus object, returns a date in moment format YYYY-MM-DD
export const formatSITAuthorizedEndDate = (sitStatus) => {
  return moment(sitStatus.currentSIT.sitAuthorizedEndDate, swaggerDateFormat).subtract(1, 'days');
};

// takes in a date and adds the days provided
export const formatEndDate = (date, days) => {
  return moment(date, swaggerDateFormat).add(days, 'days');
};

// ****************
// ****************
// date calculations
// ****************
// ****************

export function calculateDaysInPreviousSIT(totalSITDaysUsed, daysInSIT) {
  return totalSITDaysUsed - daysInSIT;
}

// returns the end date after making sure it is in the same time zone as the SIT entry date
// adds one day since we include the end date in the calculation
export const calculateEndDate = (sitEntryDate, endDate) => {
  const sitEntryMoment = moment(sitEntryDate);
  const sitEndMoment = moment(endDate);
  const endDay = moment(sitEndMoment).utcOffset(sitEntryMoment.utcOffset(), true).startOf('day').add(1, 'day');
  return endDay;
};

// calculating the total number of SIT days allowed, taking into account the entry date, the end date, and any previous SIT days used
// normalizes the dates and accounts for time zones
export const calculateSitDaysAllowance = (sitEntryDate, daysInPreviousSIT, endDate) => {
  const sitEntryMoment = moment(sitEntryDate);
  const sitEndMoment = moment(endDate);
  const endDay = moment(sitEndMoment).utcOffset(sitEntryMoment.utcOffset(), true).startOf('day');
  const startDay = sitEntryMoment.startOf('day');
  const sitDurationDays = moment.duration(endDay.diff(startDay)).asDays();
  const calculatedSitDaysAllowance = sitDurationDays + daysInPreviousSIT;
  return calculatedSitDaysAllowance;
};

// determining the authorized end date for SIT based on the entry date and approved days, adjusting for days already used in previous SIT periods
export const calculateSITEndDate = (sitEntryDate, daysApproved, daysInPreviousSIT) => {
  const sitEntryMoment = moment(sitEntryDate);
  return formatDateForDatePicker(sitEntryMoment.add(daysApproved - daysInPreviousSIT, 'days').subtract(1, 'day'));
};

export const calculateSITTotalDaysRemaining = (sitStatus, shipment) => {
  const daysRemaining = sitStatus ? sitStatus.totalDaysRemaining : shipment.sitDaysAllowance;
  if (daysRemaining > 0) {
    return daysRemaining;
  }
  return 'Expired';
};

export const calculateApprovedAndRequestedDaysCombined = (shipment, sitExtension) => {
  return shipment.sitDaysAllowance + sitExtension.requestedDays;
};

// adds the requested number of days & total days remaining
export const calculateApprovedAndRequestedDatesCombined = (sitExtension, totalDaysRemaining) => {
  return formatDateForDatePicker(
    moment()
      .add(sitExtension.requestedDays, 'days')
      .add(Number.isInteger(totalDaysRemaining) ? totalDaysRemaining : 0, 'days'),
  );
};

// ************************
// ************************
// components & UI elements
// ************************
// ************************

export const SitEndDateForm = ({ onChange }) => (
  <div className={styles.sitDatePicker} data-testid="sitEndDate">
    <DatePickerInput name="sitEndDate" label="" id="sitEndDate" onChange={onChange} />
  </div>
);

export const CurrentSITDateData = ({ currentLocation, daysInSIT, sitDepartureDate }) => {
  const currentDaysInSit = <p>{daysInSIT}</p>;

  return (
    <DataTable
      testID="currentSITDateData"
      columnHeaders={[`Total days in ${currentLocation}`, `SIT departure date`]}
      dataRow={[currentDaysInSit, sitDepartureDate]}
    />
  );
};

export const SitDaysAllowanceForm = ({ onChange }) => (
  <div className={styles.sitDatePicker}>
    <MaskedTextField
      data-testid="daysApproved"
      defaultValue="1"
      id="daysApproved"
      name="daysApproved"
      mask={Number}
      lazy={false}
      scale={0}
      signed={false}
      inputClassName={styles.approvedDaysInput}
      errorClassName={styles.errors}
      onChange={onChange}
      label="Days approved"
      labelClassName={styles.label}
    />
  </div>
);

export const SITHistoryItemHeaderDays = ({ title, approved, requested, value }) => {
  return (
    <div data-happo-hide className={styles.sitHistoryItemHeader}>
      {title}
      <span className={styles.hintText}>
        Previously approved ({approved}) + <br />
        Requested ({requested}) = {value}
      </span>
    </div>
  );
};

export const SITHistoryItemHeaderDate = ({ title, endDate, requested, value }) => {
  return (
    <div data-happo-hide className={styles.sitHistoryItemHeader}>
      {title}
      <span className={styles.hintText}>
        Previously authorized end date
        <br />({formatDateForDatePicker(endDate)}) + <br />
        days requested ({requested}) =<br /> {value}
      </span>
    </div>
  );
};

export function getSITCurrentLocation(sitStatus) {
  return sitStatus.currentSIT.location === LOCATION_TYPES.ORIGIN ? 'Origin SIT' : 'Destination SIT';
}
