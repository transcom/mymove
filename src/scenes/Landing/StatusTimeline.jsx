import PropTypes from 'prop-types';
import React, { PureComponent } from 'react';
import { get, filter, findLast, includes } from 'lodash';
import moment from 'moment';

import { displayDateRange } from 'shared/formatters';
import './StatusTimeline.css';

function getDates(source, dateType) {
  // The in progress state in PPMStatusTimeline has different expectations
  if (dateType === 'actual_move_date') {
    // if there's no approve date, then the PPM hasn't been approved yet
    // and the in progress date should not be shown
    const approveDate = get(source, 'approve_date');
    if (approveDate) {
      let date = undefined;
      // if there's an actual move date that is known and passed, show it
      // else show original move date if it has passed
      const actualMoveDate = get(source, dateType);
      const originalMoveDate = get(source, 'original_move_date');
      if (actualMoveDate && moment(actualMoveDate, 'YYYY-MM-DD').isSameOrBefore()) {
        date = actualMoveDate;
      } else if (moment(originalMoveDate, 'YYYY-MM-DD').isSameOrBefore()) {
        date = originalMoveDate;
      }
      return date;
    }

    return;
  }
  return get(source, dateType);
}

function getCurrentStatus(statuses) {
  return findLast(statuses, function(status) {
    return status.completed;
  });
}

export class PPMStatusTimeline extends React.Component {
  getStatuses() {
    return [
      { name: 'Submitted', code: 'SUBMITTED', date_type: 'submit_date' },
      { name: 'Approved', code: 'PPM_APPROVED', date_type: 'approve_date' },
      { name: 'In progress', code: 'IN_PROGRESS', date_type: 'actual_move_date' },
      { name: 'Payment requested', code: 'PAYMENT_REQUESTED' },
    ];
  }

  getCompletedStatus(status) {
    const { ppm } = this.props;

    if (status === 'SUBMITTED') {
      return true;
    }

    if (status === 'PPM_APPROVED') {
      return includes(['APPROVED', 'PAYMENT_REQUESTED', 'COMPLETED'], ppm.status);
    }

    if (status === 'IN_PROGRESS') {
      const moveInProgress = moment(ppm.original_move_date, 'YYYY-MM-DD').isSameOrBefore();
      return (moveInProgress && ppm.status === 'APPROVED') || includes(['PAYMENT_REQUESTED', 'COMPLETED'], ppm.status);
    }

    if (status === 'PAYMENT_REQUESTED') {
      return includes(['PAYMENT_REQUESTED', 'COMPLETED'], ppm.status);
    }
  }

  addDates(statuses) {
    return statuses.map(status => {
      return {
        ...status,
        dates: [getDates(this.props.ppm, status.date_type)],
      };
    });
  }

  addCompleted(statuses) {
    return statuses.map(status => {
      return {
        ...status,
        completed: this.getCompletedStatus(status.code),
      };
    });
  }

  render() {
    const statuses = this.addDates(this.addCompleted(this.getStatuses()));
    return <StatusTimeline statuses={statuses} showEstimated={false} />;
  }
}

PPMStatusTimeline.propTypes = {
  ppm: PropTypes.object.isRequired,
};

export class ShipmentStatusTimeline extends React.Component {
  getStatuses() {
    return [
      { name: 'Scheduled', code: 'SCHEDULED', date_type: 'book_date' },
      { name: 'Packed', code: 'PACKED', date_type: 'pack' },
      { name: 'Loaded', code: 'LOADED', date_type: 'pickup' },
      { name: 'In transit', code: 'IN_TRANSIT', date_type: 'transit' },
      { name: 'Delivered', code: 'DELIVERED', date_type: 'delivery' },
    ];
  }

  getCompletedStatus(status) {
    const { shipment } = this.props;
    const today = this.props.today ? moment(this.props.today) : moment();

    if (status === 'SCHEDULED') {
      return true;
    }

    if (status === 'PACKED') {
      const actualPackDate = get(shipment, 'actual_pack_date', null);
      const originalPackDate = get(shipment, 'original_pack_date', null);
      const pmSurveyPlannedPackDate = get(shipment, 'pm_survey_planned_pack_date', null);
      return actualPackDate || today.isSameOrAfter(pmSurveyPlannedPackDate) || today.isSameOrAfter(originalPackDate)
        ? true
        : false;
    }

    const actualPickupDate = get(shipment, 'actual_pickup_date', null);
    const pmSurveyPlannedPickupDate = get(shipment, 'pm_survey_planned_pickup_date', null);
    const requestedPickupDate = get(shipment, 'requested_pickup_date', null);

    if (status === 'LOADED') {
      return actualPickupDate ||
        today.isSameOrAfter(pmSurveyPlannedPickupDate, 'day') ||
        today.isSameOrAfter(requestedPickupDate, 'day')
        ? true
        : false;
    }

    if (status === 'IN_TRANSIT') {
      return (actualPickupDate && today.isAfter(actualPickupDate, 'day')) ||
        today.isAfter(pmSurveyPlannedPickupDate, 'day') ||
        today.isAfter(requestedPickupDate, 'day')
        ? true
        : false;
    }

    if (status === 'DELIVERED') {
      const actualDeliveryDate = get(shipment, 'actual_delivery_date', null);
      return actualDeliveryDate || includes(['DELIVERED', 'COMPLETED'], shipment.status) ? true : false;
    }
  }

  addDates(statuses) {
    const moveDatesSummary = get(this.props.shipment, 'move_dates_summary');
    return statuses.map(status => {
      return {
        ...status,
        dates:
          status.date_type === 'book_date'
            ? [getDates(this.props.shipment, status.date_type)]
            : getDates(moveDatesSummary, status.date_type),
      };
    });
  }

  addCompleted(statuses) {
    return statuses.map(status => {
      return {
        ...status,
        completed: this.getCompletedStatus(status.code),
      };
    });
  }

  render() {
    const statuses = this.addDates(this.addCompleted(this.getStatuses()));
    return <StatusTimeline statuses={statuses} showEstimated={true} />;
  }
}

ShipmentStatusTimeline.propTypes = {
  shipment: PropTypes.object.isRequired,
  today: PropTypes.string,
};

export class ProfileStatusTimeline extends React.Component {
  getStatuses() {
    return [
      { name: 'Profile', code: 'PROFILE', completed: true },
      { name: 'Orders', code: 'ORDERS', completed: true },
      { name: 'Move Setup', code: 'MOVE_SETUP', completed: false },
      { name: 'Review', code: 'REVIEW', completed: false },
    ];
  }

  render() {
    return <StatusTimeline statuses={this.getStatuses()} showEstimated={false} />;
  }
}

ProfileStatusTimeline.propTypes = {
  profile: PropTypes.object.isRequired,
};

class StatusTimeline extends PureComponent {
  createStatusBlock = (status, currentStatus) => {
    return (
      <StatusBlock
        name={status.name}
        code={status.code}
        key={status.code}
        dates={filter(status.dates, date => {
          return date;
        })}
        completed={status.completed}
        current={currentStatus.code === status.code}
      />
    );
  };

  render() {
    const currentStatus = getCurrentStatus(this.props.statuses);
    const statusBlocks = this.props.statuses.map(status => this.createStatusBlock(status, currentStatus));

    return (
      <div className="status_timeline">
        {statusBlocks}
        {this.props.showEstimated && <div className="legend">* Estimated</div>}
      </div>
    );
  }
}

StatusTimeline.propTypes = {
  statuses: PropTypes.array.isRequired,
};

export const StatusBlock = props => {
  const classes = ['status_block', props.code.toLowerCase()];
  if (props.completed) classes.push('status_completed');
  if (props.current) classes.push('status_current');

  return (
    <div className={classes.join(' ')}>
      <div className="status_dot" />
      <div className="status_name">{props.name}</div>
      {props.dates &&
        props.dates.length > 0 && <div className="status_dates">{displayDateRange(props.dates, 'condensed')}</div>}
    </div>
  );
};

StatusBlock.propTypes = {
  code: PropTypes.string.isRequired,
  completed: PropTypes.bool.isRequired,
  current: PropTypes.bool.isRequired,
  name: PropTypes.string.isRequired,
  dates: PropTypes.arrayOf(PropTypes.string),
};
