import PropTypes from 'prop-types';
import React, { Component } from 'react';
import connect from 'react-redux/es/connect/connect';
import { DayPicker } from 'react-day-picker';
import 'react-day-picker/lib/style.css';
import { get, isNil } from 'lodash';
import moment from 'moment';

import { formatSwaggerDate, parseSwaggerDate } from 'shared/formatters';
import { bindActionCreators } from 'redux';
import { getMoveDatesSummary } from 'shared/Entities/modules/moves';

import './DatePicker.css';
import { selectMoveDatesSummary } from 'shared/Entities/modules/moves';

const getRequestLabel = 'DatePicker.getMoveDatesSummary';

export class HHGDatePicker extends Component {
  handleDayClick = day => {
    console.log(day);
    const moveDate = day.toISOString().split('T')[0];
    this.props.input.onChange(formatSwaggerDate(day));
    this.props.getMoveDatesSummary(
      getRequestLabel,
      this.props.currentShipment.move_id,
      moveDate,
    );
  };

  componentDidMount() {
    if (this.props.currentShipment.requested_pickup_date) {
      this.props.getMoveDatesSummary(
        getRequestLabel,
        this.props.currentShipment.move_id,
        this.props.currentShipment.requested_pickup_date,
      );
    }
  }

  render() {
    function formatDate(date) {
      if (date) {
        return moment(date).format('ddd, MMMM DD');
      }
    }
    const { moveDates, currentShipment } = this.props;

    const pickupDates = get(moveDates, 'pickup', []);
    const packDates = get(moveDates, 'pack', []);
    const deliveryDates = get(moveDates, 'delivery', []);
    const transitDates = get(moveDates, 'transit', []);
    const reportDates = get(moveDates, 'report', []);

    let selectedDay =
      this.props.input.value ||
      get(currentShipment, 'requested_pickup_date', null);

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
  currentShipment: PropTypes.object,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ getMoveDatesSummary }, dispatch);
}

function mapStateToProps(state, ownProps) {
  const props = {
    moveDates: selectMoveDatesSummary(
      state,
      get(ownProps.currentShipment, 'requested_pickup_date'),
    ),
  };
  return props;
}

export default connect(mapStateToProps, mapDispatchToProps)(HHGDatePicker);
