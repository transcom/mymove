import { dateSelectionWeekendHolidayCheck } from './calendar';

import { flushPromises } from 'testUtils';

describe('dateSelectionWeekendHolidayCheck', () => {
  it('calls the dateSelectionIsWeekendHolidayAPI when passed a valid date and country code', async () => {
    const mockApiResponse = {
      data: JSON.stringify({
        is_holiday: true,
        is_weekend: false,
        country_code: 'US',
        country_name: 'United States',
      }),
    };

    const mockDateSelectionIsWeekendHolidayAPI = jest.fn().mockResolvedValue(mockApiResponse);

    const mockSetAlertMessageCallback = jest.fn();
    const mockSetIsDateSelectionAlertVisibleCallback = jest.fn();
    const mockOnErrorCallback = jest.fn();

    dateSelectionWeekendHolidayCheck(
      mockDateSelectionIsWeekendHolidayAPI,
      'US',
      new Date('2025-01-01'),
      'Requested Pickup Date',
      mockSetAlertMessageCallback,
      mockSetIsDateSelectionAlertVisibleCallback,
      mockOnErrorCallback,
    );

    // ensure the mockDateSelectionIsWeekendHolidayAPI has resolved before checking the expects since dateSelectionWeekendHolidayCheck is not asynchronous
    await flushPromises();

    expect(mockDateSelectionIsWeekendHolidayAPI).toHaveBeenCalledWith('US', '2025-01-01');
    expect(mockSetAlertMessageCallback).toHaveBeenCalled();
    expect(mockSetIsDateSelectionAlertVisibleCallback).toHaveBeenCalledWith(true);
    expect(mockOnErrorCallback).not.toHaveBeenCalled();
  });

  it('passes the holiday message to setAlertMessageCallback', async () => {
    const holidayResponse = {
      data: JSON.stringify({
        is_holiday: true,
        is_weekend: false,
        country_code: 'US',
        country_name: 'United States',
      }),
    };

    const mockDateSelectionIsWeekendHolidayAPI = jest.fn().mockResolvedValue(holidayResponse);
    const mockSetAlertMessageCallback = jest.fn();
    const mockSetIsDateSelectionAlertVisibleCallback = jest.fn();
    const mockOnErrorCallback = jest.fn();

    dateSelectionWeekendHolidayCheck(
      mockDateSelectionIsWeekendHolidayAPI,
      'US',
      new Date('2025-06-02'),
      'Requested Pickup Date',
      mockSetAlertMessageCallback,
      mockSetIsDateSelectionAlertVisibleCallback,
      mockOnErrorCallback,
    );

    // ensure the mockDateSelectionIsWeekendHolidayAPI has resolved before checking the expects since dateSelectionWeekendHolidayCheck is not asynchronous
    await flushPromises();

    expect(mockDateSelectionIsWeekendHolidayAPI).toHaveBeenCalled();
    expect(mockSetAlertMessageCallback).toHaveBeenCalledWith(
      'Requested Pickup Date 02 Jun 2025 is on a holiday in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date.',
    );
    expect(mockSetIsDateSelectionAlertVisibleCallback).toHaveBeenCalledWith(true);
    expect(mockOnErrorCallback).not.toHaveBeenCalled();
  });

  it('passes the weekend message to setAlertMessageCallback', async () => {
    const weekendResponse = {
      data: JSON.stringify({
        is_holiday: false,
        is_weekend: true,
        country_code: 'US',
        country_name: 'United States',
      }),
    };

    const mockDateSelectionIsWeekendHolidayAPI = jest.fn().mockResolvedValue(weekendResponse);
    const mockSetAlertMessageCallback = jest.fn();
    const mockSetIsDateSelectionAlertVisibleCallback = jest.fn();
    const mockOnErrorCallback = jest.fn();

    dateSelectionWeekendHolidayCheck(
      mockDateSelectionIsWeekendHolidayAPI,
      'US',
      new Date('2025-06-03'),
      'Requested Pickup Date',
      mockSetAlertMessageCallback,
      mockSetIsDateSelectionAlertVisibleCallback,
      mockOnErrorCallback,
    );

    // ensure the mockDateSelectionIsWeekendHolidayAPI has resolved before checking the expects since dateSelectionWeekendHolidayCheck is not asynchronous
    await flushPromises();

    expect(mockDateSelectionIsWeekendHolidayAPI).toHaveBeenCalled();
    expect(mockSetAlertMessageCallback).toHaveBeenCalledWith(
      'Requested Pickup Date 03 Jun 2025 is on a weekend in the United States. This date may not be accepted. A government representative may not be available to provide assistance on this date.',
    );
    expect(mockSetIsDateSelectionAlertVisibleCallback).toHaveBeenCalledWith(true);
    expect(mockOnErrorCallback).not.toHaveBeenCalled();
  });

  it('passes the holiday/weekend message to setAlertMessageCallback', async () => {
    const holidayAndWeekendResponse = {
      data: JSON.stringify({
        is_holiday: true,
        is_weekend: true,
        country_code: 'GB',
        country_name: 'Great Britain',
      }),
    };
    const mockDateSelectionIsWeekendHolidayAPI = jest.fn().mockResolvedValue(holidayAndWeekendResponse);
    const mockSetAlertMessageCallback = jest.fn();
    const mockSetIsDateSelectionAlertVisibleCallback = jest.fn();
    const mockOnErrorCallback = jest.fn();

    dateSelectionWeekendHolidayCheck(
      mockDateSelectionIsWeekendHolidayAPI,
      'GB',
      new Date('2025-06-04'),
      'Requested Delivery Date',
      mockSetAlertMessageCallback,
      mockSetIsDateSelectionAlertVisibleCallback,
      mockOnErrorCallback,
    );

    // ensure the mockDateSelectionIsWeekendHolidayAPI has resolved before checking the expects since dateSelectionWeekendHolidayCheck is not asynchronous
    await flushPromises();

    expect(mockDateSelectionIsWeekendHolidayAPI).toHaveBeenCalled();
    expect(mockSetAlertMessageCallback).toHaveBeenCalledWith(
      'Requested Delivery Date 04 Jun 2025 is on a holiday and weekend in Great Britain. This date may not be accepted. A government representative may not be available to provide assistance on this date.',
    );
    expect(mockSetIsDateSelectionAlertVisibleCallback).toHaveBeenCalledWith(true);
    expect(mockOnErrorCallback).not.toHaveBeenCalled();
  });

  it('does not call dateSelectionIsWeekendHolidayAPI and hides alert when date is invalid or countryCode is falsy', () => {
    const mockDateSelectionIsWeekendHolidayAPI = jest.fn();
    const mockSetAlertMessageCallback = jest.fn();
    const mockSetIsDateSelectionAlertVisibleCallback = jest.fn();
    const mockOnErrorCallback = jest.fn();

    // Invalid date (NaN)
    dateSelectionWeekendHolidayCheck(
      mockDateSelectionIsWeekendHolidayAPI,
      'US',
      new Date('invalid date'),
      'Requested Pickup Date',
      mockSetAlertMessageCallback,
      mockSetIsDateSelectionAlertVisibleCallback,
      mockOnErrorCallback,
    );

    // empty countryCode
    dateSelectionWeekendHolidayCheck(
      mockDateSelectionIsWeekendHolidayAPI,
      '',
      new Date(),
      'Requested Pickup Date',
      mockSetAlertMessageCallback,
      mockSetIsDateSelectionAlertVisibleCallback,
      mockOnErrorCallback,
    );

    expect(mockDateSelectionIsWeekendHolidayAPI).not.toHaveBeenCalled();
    expect(mockSetAlertMessageCallback).not.toHaveBeenCalled();
    expect(mockSetIsDateSelectionAlertVisibleCallback).toHaveBeenCalledWith(false);
    expect(mockOnErrorCallback).not.toHaveBeenCalled();
  });

  it('should execute onErrorCallback if there is an api error', async () => {
    const mockDateSelectionIsWeekendHolidayAPI = jest.fn().mockRejectedValue(new Error('Invalid country code'));
    const mockSetAlertMessageCallback = jest.fn();
    const mockSetIsDateSelectionAlertVisibleCallback = jest.fn();
    const mockOnErrorCallback = jest.fn();

    dateSelectionWeekendHolidayCheck(
      mockDateSelectionIsWeekendHolidayAPI,
      'ABCDE',
      new Date('2025-01-01'),
      'Requested Pickup Date',
      mockSetAlertMessageCallback,
      mockSetIsDateSelectionAlertVisibleCallback,
      mockOnErrorCallback,
    );

    // ensure the mockDateSelectionIsWeekendHolidayAPI has resolved before checking the expects since dateSelectionWeekendHolidayCheck is not asynchronous
    await flushPromises();

    expect(mockDateSelectionIsWeekendHolidayAPI).toHaveBeenCalledWith('ABCDE', '2025-01-01');
    expect(mockSetAlertMessageCallback).not.toHaveBeenCalled();
    expect(mockSetIsDateSelectionAlertVisibleCallback).not.toHaveBeenCalled();
    expect(mockOnErrorCallback).toHaveBeenCalled();
  });
});
