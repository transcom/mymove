import PropTypes from 'prop-types';
import React, { Component } from 'react';
import connect from 'react-redux/es/connect/connect';
import { DayPicker } from 'react-day-picker';
import 'react-day-picker/lib/style.css';
import { get, isNil } from 'lodash';

import { formatSwaggerDate, parseSwaggerDate } from 'shared/formatters';
import './DatePicker.css';
import { bindActionCreators } from 'redux';
import { getMoveDatesSummary } from './ducks';
import moment from 'moment';

export class HHGDatePicker extends Component {
  handleDayClick = day => {
    this.props.input.onChange(formatSwaggerDate(day));
    this.props.getMoveDatesSummary(this.props.shipment.move_id, day);
  };

  componentDidMount() {
    // TODO: make this actually work
    if (!isNil(this.props.requestedPickupDate)) {
      this.props.getMoveDatesSummary(
        this.props.shipment.move_id,
        this.props.requestedPickupDate,
      );
    }
  }

  render() {
    function formatDate(date) {
      if (date) {
        return moment(date).format('ddd, MMMM DD');
      }
    }

    const pickupDates = get(this.props.moveDates, 'pickup', []);
    const packDates = get(this.props.moveDates, 'pack', []);
    const deliveryDates = get(this.props.moveDates, 'delivery', []);
    const transitDates = get(this.props.moveDates, 'transit', []);
    const reportDates = get(this.props.moveDates, 'report', []);

    const selectedDay = this.props.input.value;

    return (
      <div className="form-section">
        <h3 className="instruction-heading">
          Great! Let's find a date for a moving company to move your stuff.
        </h3>
        <div className="usa-grid">
          <h4>Select a move date</h4>
          <div className="usa-width-one-third">
            <DayPicker
              onDayClick={this.handleDayClick}
              selectedDays={parseSwaggerDate(selectedDay)}
            />
          </div>

          <div className="usa-width-two-thirds">
            {selectedDay &&
              !isNil(this.props.moveDates) && (
                <table>
                  <tbody>
                    <tr>
                      <th>Preferred Moving Dates Summary</th>
                    </tr>
                    <tr>
                      <td>Movers Packing</td>
                      <td>
                        {formatDate(packDates[0])} -{' '}
                        {formatDate(packDates[packDates.length - 1])}
                        <span className="estimate">*estimated</span>
                      </td>
                    </tr>
                    <tr>
                      <td>Movers Loading Truck</td>
                      <td>{formatDate(pickupDates[0])}</td>
                    </tr>
                    <tr>
                      <td>Moving Truck in Transit</td>
                      <td>
                        {formatDate(transitDates[0])} -{' '}
                        {formatDate(transitDates[transitDates.length - 1])}
                      </td>
                    </tr>
                    <tr>
                      <td>Movers Delivering</td>
                      <td>
                        {formatDate(deliveryDates[0])}{' '}
                        <span className="estimate">*estimated</span>
                      </td>
                    </tr>
                    <tr>
                      <td>Report By Date</td>
                      <td>{formatDate(reportDates[0])}</td>
                    </tr>
                  </tbody>
                </table>
              )}
          </div>
        </div>
      </div>
    );
  }
}
HHGDatePicker.propTypes = {
  input: PropTypes.object.isRequired,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ getMoveDatesSummary }, dispatch);
}

function mapStateToProps(state) {
  const props = {
    requestedPickupDate: get(state.hhg, 'currentHhg.requestedPickupDate', null),
    shipment: get(state.hhg, 'currentHhg', {}),
    moveDates: get(state.hhg, 'moveDates', {}),
  };
  return props;
}

export default connect(mapStateToProps, mapDispatchToProps)(HHGDatePicker);
