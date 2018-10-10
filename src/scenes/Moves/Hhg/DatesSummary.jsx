import PropTypes from 'prop-types';
import React, { Component } from 'react';
import moment from 'moment';
import { get } from 'lodash';

import { selectMoveDatesSummary } from 'shared/Entities/modules/moves';
import { connect } from 'react-redux';

export class DatesSummary extends Component {
  formatDate(date) {
    if (date) {
      return moment(date).format('ddd, MMM DD');
    }
  }
  render() {
    const { moveDates } = this.props;
    const pickupDates = get(moveDates, 'pickup', []);
    const packDates = get(moveDates, 'pack', []);
    const deliveryDates = get(moveDates, 'delivery', []);
    const transitDates = get(moveDates, 'transit', []);
    const reportDates = get(moveDates, 'report', []);

    return (
      <table>
        <tbody>
          <tr>
            <th colSpan="2">Preferred Moving Dates Summary</th>
          </tr>
          <tr>
            <td>Movers Packing</td>
            <td>
              {this.formatDate(packDates[0])} -{' '}
              {this.formatDate(packDates[packDates.length - 1])}*
              <span className="estimate">(estimated)</span>
            </td>
          </tr>
          <tr>
            <td>Movers Loading Truck</td>
            <td>{this.formatDate(pickupDates[0])}</td>
          </tr>
          <tr>
            <td>Moving Truck in Transit</td>
            <td>
              {this.formatDate(transitDates[0])} -{' '}
              {this.formatDate(transitDates[transitDates.length - 1])}
            </td>
          </tr>
          <tr>
            <td>Movers Delivering</td>
            <td>
              {this.formatDate(deliveryDates[0])}{' '}
              <span className="estimate">(estimated)</span>
            </td>
          </tr>
          <tr>
            <td>Report By Date</td>
            <td>{this.formatDate(reportDates[0])}</td>
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
