import moment from 'moment';

import {
  calculateEndDate,
  calculateSitDaysAllowance,
  calculateSITEndDate,
  calculateDaysInPreviousSIT,
  calculateSITTotalDaysRemaining,
  formatSITDepartureDate,
  formatSITEntryDate,
  formatSITAuthorizedEndDate,
  getSITCurrentLocation,
  calculateApprovedAndRequestedDaysCombined,
  calculateApprovedAndRequestedDatesCombined,
} from './sitFormatters';

import { DEFAULT_EMPTY_VALUE } from 'shared/constants';
import { LOCATION_TYPES } from 'types/sitStatusShape';
import { formatDateForDatePicker } from 'shared/dates';

// ****************
// Test for calculateEndDate
// ****************
describe('calculateEndDate', () => {
  it('should calculate the correct end date', () => {
    const sitEntryDate = '2023-08-01T00:00:00Z';
    const endDate = '2023-08-10T00:00:00Z';

    const expectedEndDate = moment('2023-08-11T00:00:00Z').format('YYYY-MM-DD');

    const result = calculateEndDate(sitEntryDate, endDate).format('YYYY-MM-DD');
    expect(result).toBe(expectedEndDate);
  });
});

// ****************
// Test for calculateSitDaysAllowance
// ****************
describe('calculateSitDaysAllowance', () => {
  it('should calculate the correct SIT days allowance', () => {
    const sitEntryDate = '2023-08-01T00:00:00Z';
    const daysInPreviousSIT = 5;
    const endDate = '2023-08-10T00:00:00Z';

    const expectedAllowance = 14; // 10 days from entry to end date plus 5 previous days

    const result = calculateSitDaysAllowance(sitEntryDate, daysInPreviousSIT, endDate);
    expect(result).toBe(expectedAllowance);
  });

  it('should handle different time zones correctly', () => {
    const sitEntryDate = '2023-08-01T00:00:00+02:00';
    const daysInPreviousSIT = 3;
    const endDate = '2023-08-10T00:00:00-05:00';

    const expectedAllowance = 13; // 10 days from entry to end date plus 3 previous days

    const result = calculateSitDaysAllowance(sitEntryDate, daysInPreviousSIT, endDate);
    expect(result).toBe(expectedAllowance);
  });
});

// ****************
// Test for calculateSITEndDate
// ****************
describe('calculateSITEndDate', () => {
  it('should calculate the correct SIT end date', () => {
    const sitEntryDate = '2023-08-01T00:00:00Z';
    const daysApproved = 15;
    const daysInPreviousSIT = 5;

    const expectedEndDate = '10 Aug 2023'; // 15 days approved minus 5 previous days from the entry date

    const result = calculateSITEndDate(sitEntryDate, daysApproved, daysInPreviousSIT);
    expect(result).toBe(expectedEndDate);
  });

  it('should handle different time zones correctly', () => {
    const sitEntryDate = '2023-08-01T00:00:00+02:00';
    const daysApproved = 20;
    const daysInPreviousSIT = 10;

    const expectedEndDate = '09 Aug 2023'; // 20 days approved minus 10 previous days from the entry date

    const result = calculateSITEndDate(sitEntryDate, daysApproved, daysInPreviousSIT);
    expect(result).toBe(expectedEndDate);
  });
});

// ****************
// Test for calculateDaysInPreviousSIT
// ****************
describe('calculateDaysInPreviousSIT', () => {
  it('should calculate the correct days in previous SIT', () => {
    const totalSITDaysUsed = 15;
    const daysInSIT = 10;

    const expectedDaysInPreviousSIT = 5;

    const result = calculateDaysInPreviousSIT(totalSITDaysUsed, daysInSIT);
    expect(result).toBe(expectedDaysInPreviousSIT);
  });
});

// ****************
// Test for calculateSITTotalDaysRemaining
// ****************
describe('calculateSITTotalDaysRemaining', () => {
  it('should calculate the correct total days remaining with sitStatus', () => {
    const sitStatus = { totalDaysRemaining: 5 };
    const shipment = { sitDaysAllowance: 10 };

    const expectedDaysRemaining = 5;

    const result = calculateSITTotalDaysRemaining(sitStatus, shipment);
    expect(result).toBe(expectedDaysRemaining);
  });

  it('should calculate the correct total days remaining without sitStatus', () => {
    const sitStatus = null;
    const shipment = { sitDaysAllowance: 10 };

    const expectedDaysRemaining = 10;

    const result = calculateSITTotalDaysRemaining(sitStatus, shipment);
    expect(result).toBe(expectedDaysRemaining);
  });

  it('should return "Expired" when days remaining is less than or equal to 0', () => {
    const sitStatus = { totalDaysRemaining: -1 };
    const shipment = { sitDaysAllowance: 0 };

    const expectedDaysRemaining = 'Expired';

    const result = calculateSITTotalDaysRemaining(sitStatus, shipment);
    expect(result).toBe(expectedDaysRemaining);
  });
});

// ****************
// Test for formatSITDepartureDate
// ****************
describe('formatSITDepartureDate', () => {
  it('should format the SIT departure date correctly', () => {
    const date = '2023-08-10';
    const expectedFormattedDate = '10 Aug 2023';

    const result = formatSITDepartureDate(date);
    expect(result).toBe(expectedFormattedDate);
  });

  it('should return the default empty value if the date is invalid', () => {
    const date = null;
    const expectedFormattedDate = DEFAULT_EMPTY_VALUE;

    const result = formatSITDepartureDate(date);
    expect(result).toBe(expectedFormattedDate);
  });
});

// ****************
// Test for formatSITEntryDate
// ****************
describe('formatSITEntryDate', () => {
  it('should format the SIT entry date correctly', () => {
    const date = '2023-08-01';
    const expectedFormattedDate = '01 Aug 2023';

    const result = formatSITEntryDate(date);
    expect(result).toBe(expectedFormattedDate);
  });

  it('should return the default empty value if the date is invalid', () => {
    const date = null;
    const expectedFormattedDate = DEFAULT_EMPTY_VALUE;

    const result = formatSITEntryDate(date);
    expect(result).toBe(expectedFormattedDate);
  });
});

// ****************
// Test for formatSITAuthorizedEndDate
// ****************
describe('formatSITAuthorizedEndDate', () => {
  it('should format the SIT authorized end date correctly', () => {
    const sitStatus = { currentSIT: { sitAuthorizedEndDate: '2023-08-10T00:00:00Z' } };
    const expectedFormattedDate = moment('2023-08-09T00:00:00Z').format('YYYY-MM-DD');

    const result = formatSITAuthorizedEndDate(sitStatus).format('YYYY-MM-DD');
    expect(result).toBe(expectedFormattedDate);
  });
});

// ****************
// Test for getSITCurrentLocation
// ****************
describe('getSITCurrentLocation', () => {
  it('should return "Origin SIT" for ORIGIN location type', () => {
    const sitStatus = { currentSIT: { location: LOCATION_TYPES.ORIGIN } };
    const expectedLocation = 'Origin SIT';

    const result = getSITCurrentLocation(sitStatus);
    expect(result).toBe(expectedLocation);
  });

  it('should return "Destination SIT" for non-ORIGIN location type', () => {
    const sitStatus = { currentSIT: { location: LOCATION_TYPES.DESTINATION } };
    const expectedLocation = 'Destination SIT';

    const result = getSITCurrentLocation(sitStatus);
    expect(result).toBe(expectedLocation);
  });
});

// ****************
// Test for calculateApprovedAndRequestedDaysCombined
// ****************
describe('calculateApprovedAndRequestedDaysCombined', () => {
  it('should calculate the correct combined days', () => {
    const shipment = { sitDaysAllowance: 10 };
    const sitExtension = { requestedDays: 5 };

    const expectedCombinedDays = 15;

    const result = calculateApprovedAndRequestedDaysCombined(shipment, sitExtension);
    expect(result).toBe(expectedCombinedDays);
  });
});

// ****************
// Test for calculateApprovedAndRequestedDatesCombined
// ****************
describe('calculateApprovedAndRequestedDatesCombined', () => {
  it('should calculate the combined date correctly when both sitExtension.requestedDays and totalDaysRemaining are valid numbers', () => {
    const sitExtension = { requestedDays: 10 };
    const totalDaysRemaining = 5;

    const expectedDate = formatDateForDatePicker(moment().add(10, 'days').add(5, 'days'));

    const result = calculateApprovedAndRequestedDatesCombined(sitExtension, totalDaysRemaining);
    expect(result).toBe(expectedDate);
  });

  it('should calculate the combined date correctly when sitExtension.requestedDays is valid and totalDaysRemaining is not a number', () => {
    const sitExtension = { requestedDays: 10 };
    const totalDaysRemaining = 'invalid';

    const expectedDate = formatDateForDatePicker(moment().add(10, 'days'));

    const result = calculateApprovedAndRequestedDatesCombined(sitExtension, totalDaysRemaining);
    expect(result).toBe(expectedDate);
  });

  it('should calculate the combined date correctly when both sitExtension.requestedDays and totalDaysRemaining are zero', () => {
    const sitExtension = { requestedDays: 0 };
    const totalDaysRemaining = 0;

    const expectedDate = formatDateForDatePicker(moment());

    const result = calculateApprovedAndRequestedDatesCombined(sitExtension, totalDaysRemaining);
    expect(result).toBe(expectedDate);
  });
});
