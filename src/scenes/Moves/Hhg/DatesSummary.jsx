import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { get, isNil } from 'lodash';

import { displayDateRange } from 'shared/formatters';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';

export class DatesSummary extends Component {
  render() {
    const { moveDates } = this.props;
    const pickupDates = get(moveDates, 'pickup', []);
    const packDates = get(moveDates, 'pack', []);
    const deliveryDates = get(moveDates, 'delivery', []);
    const transitDates = get(moveDates, 'transit', []);
    const reportDates = get(moveDates, 'report', []);

    return isNil(this.props.moveDates) ? (
      <LoadingPlaceholder />
    ) : (
      <table>
        <tbody>
          <tr>
            <th colSpan="3">Preferred Moving Dates Summary</th>
          </tr>

          <tr>
            <td aria-label="pattern">
              <div className="legend-square DayPicker-Day--pack" />
            </td>
            <td className="legend-label">Movers Packing</td>
            <td>
              {displayDateRange(packDates)}
              <span className="estimate">(estimated)</span>
            </td>
          </tr>
          <tr>
            <td aria-label="pattern">
              <div className="legend-square DayPicker-Day--pickup" />
            </td>
            <td className="legend-label"> Movers Loading Truck</td>
            <td>{displayDateRange(pickupDates)}</td>
          </tr>
          <tr>
            <td aria-label="pattern">
              <div className="legend-square DayPicker-Day--transit" />
            </td>
            <td className="legend-label">Moving Truck in Transit</td>
            <td>{displayDateRange(transitDates)}</td>
          </tr>
          <tr>
            <td aria-label="pattern">
              <div className="legend-square DayPicker-Day--delivery" />
            </td>
            <td className="legend-label">Movers Delivering</td>
            <td>
              {displayDateRange(deliveryDates)}
              <span className="estimate">(estimated)</span>
            </td>
          </tr>
          <tr>
            <td aria-label="pattern">
              <div className="legend-square DayPicker-Day--report" />
            </td>
            <td className="legend-label">Report By Date</td>
            <td>{displayDateRange(reportDates)}</td>
          </tr>
        </tbody>
      </table>
    );
  }
}

DatesSummary.propTypes = {
  moveDates: PropTypes.object,
};

export default DatesSummary;
