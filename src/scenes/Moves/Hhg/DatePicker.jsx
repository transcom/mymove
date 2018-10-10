import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { DayPicker } from 'react-day-picker';
import 'react-day-picker/lib/style.css';
import { get } from 'lodash';
import moment from 'moment';

import { formatSwaggerDate, parseSwaggerDate } from 'shared/formatters';
import { bindActionCreators } from 'redux';
import { getMoveDatesSummary } from 'shared/Entities/modules/moves';
import DatesSummary from 'scenes/Moves/Hhg/DatesSummary.jsx';

import './DatePicker.css';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';

const getRequestLabel = 'DatePicker.getMoveDatesSummary';

export class HHGDatePicker extends Component {
  handleDayClick = (day, { disabled }) => {
    if (disabled) {
      return;
    }
    const moveDate = day.toISOString().split('T')[0];
    this.props.input.onChange(formatSwaggerDate(day));
    this.props.getMoveDatesSummary(
      getRequestLabel,
      this.props.currentShipment.move_id,
      moveDate,
    );
  };

  isDayDisabled = day => {
    const availableMoveDates = this.props.availableMoveDates;
    if (!availableMoveDates) {
      return true;
    }

    const momentDay = moment(day);
    if (
      momentDay.isBefore(availableMoveDates.minDate, 'day') ||
      momentDay.isAfter(availableMoveDates.maxDate, 'day')
    ) {
      return true;
    }

    return !availableMoveDates.available.find(element =>
      momentDay.isSame(element, 'day'),
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

    const availableMoveDates = this.props.availableMoveDates;
    return (
      <div className="form-section">
        <h3 className="instruction-heading">
          Great! Let's find a date for a moving company to move your stuff.
        </h3>
        {availableMoveDates ? (
          <div className="usa-grid">
            <h4>Select a move date</h4>
            <div className="usa-width-one-third">
              <DayPicker
                onDayClick={this.handleDayClick}
                selectedDays={parseSwaggerDate(selectedDay)}
                disabledDays={this.isDayDisabled}
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
        ) : (
          <LoadingPlaceholder />
        )}
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
