import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { DayPicker } from 'react-day-picker';
import 'react-day-picker/lib/style.css';

import { formatSwaggerDate } from 'shared/formatters';
import './DatePicker.css';

export class HHGDatePicker extends Component {
  state = { showInfo: false };

  constructor(props) {
    super(props);
    this.handleDayClick = this.handleDayClick.bind(this);
    this.state = {
      selectedDay: undefined,
    };
  }

  handleDayClick(day) {
    this.setState({ selectedDay: day });
    this.setState({ showInfo: true });
    this.props.setDate(formatSwaggerDate(day));
  }

  render() {
    return (
      <div className="form-section">
        <h3 className="instruction-heading usa-heading">
          Great! Let's find a date for a moving company to move your stuff.
        </h3>
        <h4 className="usa-heading">Select a move date</h4>
        <div className="usa-grid">
          <div className="usa-width-one-third">
            <DayPicker
              onDayClick={this.handleDayClick}
              selectedDays={this.state.selectedDay}
            />
          </div>

          <div className="usa-width-two-thirds">
            {this.state.showInfo && (
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
  error: PropTypes.object,
};

export default HHGDatePicker;
