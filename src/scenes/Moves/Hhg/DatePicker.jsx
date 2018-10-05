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

export class HHGDatePicker extends Component {
  handleDayClick = day => {
    this.props.input.onChange(formatSwaggerDate(day));
    this.props.getMoveDatesSummary(day);
  };

  componentDidMount() {
    if (!isNil(this.props.requestedPickupDate)) {
      this.props.getMoveDatesSummary(this.props.requestedPickupDate);
    }
  }

  render() {
    // const {
    //     packDates,
    //     pickupDates,
    //     transitDates,
    //     deliveryDates,
    //     reportDates,
    // } = this.props.moveDates;

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
            {selectedDay && (
              <table className="Todo-phase2">
                <tbody>
                  <tr>
                    <th className="Todo-phase2">
                      Preferred Moving Dates Summary
                    </th>
                  </tr>
                  <tr>
                    <td>Movers Packing</td>
                    <td className="Todo-phase2">
                      Wed, June 6 - Thur, June 7{' '}
                      <span className="estimate">*estimated</span>
                    </td>
                  </tr>
                  <tr>
                    <td>Movers Loading Truck</td>
                    <td className="Todo-phase2">Fri, June 8</td>
                  </tr>
                  <tr>
                    <td>Moving Truck in Transit</td>
                    <td className="Todo-phase2">Fri, June 8 - Mon, June 11</td>
                  </tr>
                  <tr>
                    <td>Movers Delivering</td>
                    <td className="Todo-phase2">
                      Tues, June 12 <span className="estimate">*estimated</span>
                    </td>
                  </tr>
                  <tr>
                    <td>Report By Date</td>
                    <td className="Todo-phase2">Monday, July 16</td>
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
    moveDates: state.moveDates,
  };
  return props;
}

export default connect(mapStateToProps, mapDispatchToProps)(HHGDatePicker);
