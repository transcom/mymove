import PropTypes from 'prop-types';
import React, { Component } from 'react';
import moment from 'moment';
import { get, isNil } from 'lodash';
import { connect } from 'react-redux';

import { selectMoveDatesSummary } from 'shared/Entities/modules/moves';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';

export class DatesSummary extends Component {
  formatDate(date) {
    if (date) {
      return moment(date).format('ddd, MMM DD');
    }
  }

  displayDateRange(dates) {
    let span = '';
    if (dates.length > 1) {
      span = ` - ${this.formatDate(dates[dates.length - 1])}`;
    }
    return this.formatDate(dates[0]) + span;
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
            <th colSpan="2">Preferred Moving Dates Summary</th>
          </tr>

          <tr>
            <td>Movers Packing</td>
            <td>
              {this.displayDateRange(packDates)}
              <span className="estimate">(estimated)</span>
            </td>
          </tr>
          <tr>
            <td> Movers Loading Truck</td>
            <td>{this.displayDateRange(pickupDates)}</td>
          </tr>
          <tr>
            <td>Moving Truck in Transit</td>
            <td>{this.displayDateRange(transitDates)}</td>
          </tr>
          <tr>
            <td>Movers Delivering</td>
            <td>
              {this.displayDateRange(deliveryDates)}
              <span className="estimate">(estimated)</span>
            </td>
          </tr>
          <tr>
            <td>Report By Date</td>
            <td>{this.displayDateRange(reportDates)}</td>
          </tr>
        </tbody>
      </table>
    );
  }
}

DatesSummary.propTypes = {
  moveDate: PropTypes.string.isRequired,
};

function mapStateToProps(state, ownProps) {
  const props = {
    moveDates: selectMoveDatesSummary(state, ownProps.moveDate),
  };
  return props;
}

export default connect(mapStateToProps)(DatesSummary);
