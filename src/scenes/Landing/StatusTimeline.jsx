import PropTypes from 'prop-types';
import React, { PureComponent } from 'react';
import { get, findLast } from 'lodash';
import moment from 'moment';

import { displayDateRange } from 'shared/formatters';
import './StatusTimeline.css';

function getDates(source, dateType) {
  if (dateType === 'book_date') {
    return [get(source, dateType)];
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
      { name: 'Submitted', code: 'SUBMITTED' },
      { name: 'Approved', code: 'PPM_APPROVED' },
      { name: 'In progress', code: 'IN_PROGRESS' },
      { name: 'Payment requested', code: 'PAYMENT_REQUESTED' },
      { name: 'Payment approved', code: 'PAYMENT_APPROVED' },
    ];
  }

  getCompletedStatus(status) {
    const { ppm } = this.props;

    if (status === 'SUBMITTED') {
      return true;
    }

    if (status === 'PPM_APPROVED') {
      return ['APPROVED', 'PAYMENT_REQUESTED', 'COMPLETED'].includes(ppm.status);
    }

    if (status === 'IN_PROGRESS') {
      const moveInProgress = moment(ppm.original_move_date, 'YYYY-MM-DD').isSameOrBefore();
      return moveInProgress;
    }

    if (status === 'PAYMENT_REQUESTED') {
      return ppm.status === 'PAYMENT_REQUESTED';
    }

    if (status === 'PAYMENT_APPROVED') {
      return ppm.status === 'COMPLETED';
    }
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
    const statuses = this.addCompleted(this.getStatuses());
    return <StatusTimeline statuses={statuses} />;
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
      return actualDeliveryDate ? true : false;
    }
  }

  addDates(statuses) {
    const moveDatesSummary = get(this.props.shipment, 'move_dates_summary');
    return statuses.map(status => {
      return {
        ...status,
        dates:
          status.date_type === 'book_date'
            ? getDates(this.props.shipment, status.date_type)
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
    const statuses = this.addCompleted(this.addDates(this.getStatuses()));
    return <StatusTimeline statuses={statuses} />;
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
    return <StatusTimeline statuses={this.getStatuses()} />;
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
        dates={status.dates}
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
        <div className="legend">* Estimated</div>
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
