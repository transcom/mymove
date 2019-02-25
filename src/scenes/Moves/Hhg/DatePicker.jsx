import PropTypes from 'prop-types';
import React, { Component, Fragment } from 'react';
import { connect } from 'react-redux';
import { DayPicker } from 'react-day-picker';
import { withRouter } from 'react-router-dom';
import { get } from 'lodash';
import 'react-day-picker/lib/style.css';
import moment from 'moment';

import { formatSwaggerDate, parseSwaggerDate } from 'shared/formatters';
import { bindActionCreators } from 'redux';
import { getMoveDatesSummary, selectMoveDatesSummary } from 'shared/Entities/modules/moves';
import DatesSummary from 'scenes/Moves/Hhg/DatesSummary.jsx';

import './DatePicker.css';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import { loadEntitlementsFromState } from 'shared/entitlements';
import { getAvailableMoveDates, selectAvailableMoveDates } from 'shared/Entities/modules/calendar';

const getAvailableMoveDatesLabel = 'MoveDate.getAvailableMoveDates';

function createModifiers(moveDates) {
  if (!moveDates) {
    return null;
  }

  return {
    pack: convertDateStringArray(moveDates.pack),
    pickup: convertDateStringArray(moveDates.pickup),
    transit: convertDateStringArray(moveDates.transit),
    delivery: convertDateStringArray(moveDates.delivery),
    report: convertDateStringArray(moveDates.report),
  };
}

function convertDateStringArray(dateStrings) {
  return dateStrings && dateStrings.map(dateString => parseSwaggerDate(dateString));
}

export class HHGDatePicker extends Component {
  constructor(props) {
    super(props);
    this.state = {
      selectedDay: null,
    };
  }

  handleDayClick = (day, { disabled }) => {
    if (disabled) {
      return;
    }
    const moveDate = formatSwaggerDate(day);
    this.props.input.onChange(moveDate);
    this.props.getMoveDatesSummary(this.props.moveID, moveDate);
    this.setState({
      selectedDay: moveDate,
    });
  };

  isDayDisabled = day => {
    const availableMoveDates = this.props.availableMoveDates;
    if (!availableMoveDates) {
      return true;
    }

    const momentDay = moment(day);
    if (momentDay.isBefore(availableMoveDates.minDate, 'day') || momentDay.isAfter(availableMoveDates.maxDate, 'day')) {
      return true;
    }

    return !availableMoveDates.available.find(element => momentDay.isSame(element, 'day'));
  };

  componentDidUpdate(prevProps) {
    const moveDateIsSavedDate =
      get(this.props, 'currentShipment.requested_pickup_date') &&
      this.props.currentShipment.requested_pickup_date === this.props.input.value;
    if (this.props.currentShipment !== prevProps.currentShipment && this.props.currentShipment.requested_pickup_date) {
      if (!moveDateIsSavedDate) {
        this.props.getMoveDatesSummary(this.props.moveID, this.props.currentShipment.requested_pickup_date);
      }
      this.setState({
        selectedDay: this.props.input.value || this.props.currentShipment.requested_pickup_date,
      });
    }
  }

  componentDidMount() {
    this.props.getAvailableMoveDates(getAvailableMoveDatesLabel, this.props.today);
    this.setState({
      selectedDay: this.props.input.value || get(this.props, 'currentShipment.requested_pickup_date'),
    });
  }

  render() {
    const availableMoveDates = this.props.availableMoveDates;
    const parsedSelectedDay = parseSwaggerDate(this.state.selectedDay);
    const entitlementSum = get(this.props, 'entitlement.sum');
    return (
      <div className="form-section hhg-date-picker">
        {availableMoveDates ? (
          <Fragment>
            <div className="usa-grid">
              <div className="usa-width-one-whole">
                <h3 className="instruction-heading">Pick a moving date.</h3>
              </div>
            </div>
            <div className="usa-grid">
              <div className="usa-width-one-third">
                <DayPicker
                  onDayClick={this.handleDayClick}
                  month={parsedSelectedDay || (availableMoveDates && availableMoveDates.minDate)}
                  disabledDays={this.isDayDisabled}
                  modifiers={this.props.modifiers}
                  showOutsideDays
                />
              </div>
              <div className="usa-width-two-thirds">
                {this.state.selectedDay && <DatesSummary moveDates={this.props.moveDates} />}
              </div>
            </div>
          </Fragment>
        ) : (
          <LoadingPlaceholder />
        )}
        <div className="usa-grid">
          <div className="usa-width-one-whole pack-days-notice">
            <p>Can't find a date that works? Talk with a move counselor in your local Transportation office (PPPO).</p>
            {entitlementSum ? (
              <div className="weight-info-box">
                * It takes 1 day for every 5,000 lbs of stuff movers need to pack. You have an allowance of{' '}
                {entitlementSum.toLocaleString()} lbs, so we estimate it will take {Math.ceil(entitlementSum / 5000)}{' '}
                days to pack.
              </div>
            ) : (
              <LoadingPlaceholder />
            )}
          </div>
        </div>
      </div>
    );
  }
}

HHGDatePicker.propTypes = {
  input: PropTypes.object.isRequired,
  moveID: PropTypes.string.isRequired,
  currentShipment: PropTypes.object,
  moveDates: PropTypes.object,
  availableMoveDates: PropTypes.object,
  today: PropTypes.string,
};

function mapStateToProps(state, ownProps) {
  const moveDate = ownProps.input.value || get(ownProps, 'currentShipment.requested_pickup_date');
  const today = new Date().toISOString().split('T')[0];
  const moveDateIsSavedDate =
    get(ownProps, 'currentShipment.requested_pickup_date') &&
    ownProps.currentShipment.requested_pickup_date === ownProps.input.value;
  const moveDates = moveDateIsSavedDate
    ? ownProps.currentShipment.move_dates_summary
    : selectMoveDatesSummary(state, ownProps.moveID, moveDate);
  return {
    moveDates: moveDates,
    modifiers: createModifiers(moveDates),
    entitlement: loadEntitlementsFromState(state),
    availableMoveDates: selectAvailableMoveDates(state, today),
    today: today,
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ getMoveDatesSummary, getAvailableMoveDates }, dispatch);
}

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(HHGDatePicker));
