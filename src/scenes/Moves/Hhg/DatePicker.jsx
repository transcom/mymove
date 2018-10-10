import PropTypes from 'prop-types';
import React, { Component } from 'react';
import connect from 'react-redux/es/connect/connect';
import { DayPicker } from 'react-day-picker';
import 'react-day-picker/lib/style.css';
import { get } from 'lodash';

import { formatSwaggerDate, parseSwaggerDate } from 'shared/formatters';
import { bindActionCreators } from 'redux';
import { getMoveDatesSummary } from 'shared/Entities/modules/moves';
import DatesSummary from 'scenes/Moves/Hhg/DatesSummary.jsx';

import './DatePicker.css';

const getRequestLabel = 'DatePicker.getMoveDatesSummary';

export class HHGDatePicker extends Component {
  handleDayClick = day => {
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
    let selectedDay =
      this.props.input.value ||
      get(this.props.currentShipment, 'requested_pickup_date');

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
              <DatesSummary
                moveDate={
                  this.props.input.value ||
                  get(this.props.currentShipment, 'requested_pickup_date')
                }
              />
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

export default connect(() => ({}), mapDispatchToProps)(HHGDatePicker);
