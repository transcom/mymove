import PropTypes from 'prop-types';
import React, { Component } from 'react';
import moment from 'moment';
import { get, isNil } from 'lodash';

import LoadingPlaceholder from 'shared/LoadingPlaceholder';

export class DatesSummary extends Component {
  formatDate(date) {
    if (date) {
      return moment(date).format('ddd, MMM DD');
    }
  }

  displayDateRange(dates) {
    let span = '';
    let firstDate = '';
    if (dates.length > 1) {
      span = ` - ${this.formatDate(dates[dates.length - 1])}`;
    }
    if (dates.length >= 1) {
      firstDate = this.formatDate(dates[0]);
    }
    return firstDate + span;
  }

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
              {this.displayDateRange(packDates)}
              <span className="estimate">(estimated)</span>
            </td>
          </tr>
          <tr>
            <td aria-label="pattern">
              <div className="legend-square DayPicker-Day--pickup" />
            </td>
            <td className="legend-label"> Movers Loading Truck</td>
            <td>{this.displayDateRange(pickupDates)}</td>
          </tr>
          <tr>
            <td aria-label="pattern">
              <div className="legend-square DayPicker-Day--transit" />
            </td>
            <td className="legend-label">Moving Truck in Transit</td>
            <td>{this.displayDateRange(transitDates)}</td>
          </tr>
          <tr>
            <td aria-label="pattern">
              <div className="legend-square DayPicker-Day--delivery" />
            </td>
            <td className="legend-label">Movers Delivering</td>
            <td>
              {this.displayDateRange(deliveryDates)}
              <span className="estimate">(estimated)</span>
            </td>
          </tr>
          <tr>
            <td aria-label="pattern">
              <div className="legend-square DayPicker-Day--report" />
            </td>
            <td className="legend-label">Report By Date</td>
            <td>{this.displayDateRange(reportDates)}</td>
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
