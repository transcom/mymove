import { formatDateForDatePicker } from 'shared/dates';

/**
 * Function to handle date selection validation for input date for particular country.
 *
 * @param {function} dateSelectionIsWeekendHolidayAPI is webservice API
 * @param {string} countryCode is country code for date
 * @param {date} date is date to be verified
 * @param {string} label is consumer label
 * @param {function} setAlertMessageCallback is callback function to set alert validation message
 * @param {function} setIsDateSelectionAlertVisibleCallback is callback function to set alert message visibility
 * @param {function} onErrorCallback is callback function for on error handling
 */
export function dateSelectionWeekendHolidayCheck(
  dateSelectionIsWeekendHolidayAPI,
  countryCode,
  date,
  label,
  setAlertMessageCallback,
  setIsDateSelectionAlertVisibleCallback,
  onErrorCallback,
) {
  if (!Number.isNaN(date.getTime()) && countryCode) {
    // Format date to yyyy-mm-dd format
    const dateSelection = date.toISOString().replace(/T.*/, '');
    dateSelectionIsWeekendHolidayAPI(countryCode, dateSelection)
      .then((response) => {
        const info = JSON.parse(response.data);
        if (info.is_holiday || info.is_weekend) {
          let type = '';
          if (info.is_holiday) {
            type = 'holiday';
          }
          if (info.is_weekend) {
            type = 'weekend';
          }
          if (info.is_holiday && info.is_weekend) {
            type = 'holiday and weekend';
          }

          const countryLabel = info.country_code === 'US' ? `the ${info.country_name}` : info.country_name;
          const message = `${label} ${formatDateForDatePicker(
            dateSelection,
          )} is on a ${type} in ${countryLabel}. This date may not be accepted. A government representative may not be available to provide assistance on this date.`;
          setAlertMessageCallback(message);
          setIsDateSelectionAlertVisibleCallback(true);
        } else if (!info.is_holiday && !info.is_weekend) {
          setIsDateSelectionAlertVisibleCallback(false);
        }
      })
      .catch((error) => {
        onErrorCallback(error);
      });
  } else {
    // Can't determine the date is a holiday/weekend without the required paramters, so remove alert message.
    setIsDateSelectionAlertVisibleCallback(false);
  }
}

export default dateSelectionWeekendHolidayCheck;
