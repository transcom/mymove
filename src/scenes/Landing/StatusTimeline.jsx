import PropTypes from 'prop-types';
import React, { PureComponent } from 'react';
import { get, indexOf } from 'lodash';
import moment from 'moment';

import { displayDateRange } from 'shared/formatters';
import './StatusTimeline.css';

function getDates(shipment, dateType) {
  if (dateType === 'book_date') {
    return [get(shipment, dateType)];
  }
  return get(shipment, dateType);
}

function checkIfCompleted(statuses, statusToCheck) {
  if (!statuses.includes(statusToCheck)) {
    return false;
  }
  return indexOf(statuses, statusToCheck) < statuses.length;
}

function checkIfCurrent(statuses, statusToCheck) {
  if (!statuses.includes(statusToCheck)) {
    return false;
  }
  return indexOf(statuses, statusToCheck) === statuses.length - 1;
}

export class PPMStatusTimeline extends React.Component {
  statuses() {
    return [
      { name: 'Submitted', code: 'SUBMITTED' },
      { name: 'Approved', code: 'PPM_APPROVED' },
      { name: 'In progress', code: 'IN_PROGRESS' },
      { name: 'Completed', code: 'COMPLETED' },
    ];
  }

  markedStatuses() {
    const { ppm } = this.props;

    let markedStatuses = ['SUBMITTED'];

    if (ppm.status === 'SUBMITTED') {
      return markedStatuses;
    }

    markedStatuses.push('PPM_APPROVED');
    if (ppm.status === 'APPROVED') {
      const moveInProgress = moment(ppm.original_move_date, 'YYYY-MM-DD').isSameOrBefore();
      if (moveInProgress) {
        markedStatuses.push('IN_PROGRESS');
      }
      return markedStatuses;
    }

    if (ppm.status === 'COMPLETED') {
      markedStatuses.push('IN_PROGRESS');
      markedStatuses.push('COMPLETED');
    }

    return markedStatuses;
  }

  render() {
    return <StatusTimelineContainer statuses={this.statuses()} markedStatuses={this.markedStatuses()} />;
  }
}

PPMStatusTimeline.propTypes = {
  ppm: PropTypes.object.isRequired,
};

export class ShipmentStatusTimeline extends React.Component {
  statuses() {
    return [
      { name: 'Scheduled', code: 'SCHEDULED', date_type: 'book_date' },
      { name: 'Packed', code: 'PACKED', date_type: 'pack' },
      { name: 'Loaded', code: 'LOADED', date_type: 'pickup' },
      { name: 'In transit', code: 'IN_TRANSIT', date_type: 'transit' },
      { name: 'Delivered', code: 'DELIVERED', date_type: 'delivery' },
    ];
  }

  markedStatuses() {
    const { shipment } = this.props;
    const today = this.props.today ? moment(this.props.today) : moment();

    const actualPackDate = get(shipment, 'actual_pack_date', null);
    const originalPackDate = get(shipment, 'original_pack_date', null);
    const pmSurveyPlannedPackDate = get(shipment, 'pm_survey_planned_pack_date', null);
    const actualPickupDate = get(shipment, 'actual_pickup_date', null);
    const pmSurveyPlannedPickupDate = get(shipment, 'pm_survey_planned_pickup_date', null);
    const requestedPickupDate = get(shipment, 'requested_pickup_date', null);
    const actualDeliveryDate = get(shipment, 'actual_delivery_date', null);
    let markedStatuses = ['SCHEDULED'];

    if (actualPackDate || today.isSameOrAfter(pmSurveyPlannedPackDate) || today.isSameOrAfter(originalPackDate)) {
      markedStatuses.push('PACKED');
    } else {
      return markedStatuses;
    }

    if (
      actualPickupDate ||
      today.isSameOrAfter(pmSurveyPlannedPickupDate, 'day') ||
      today.isSameOrAfter(requestedPickupDate, 'day')
    ) {
      markedStatuses.push('LOADED');
    } else {
      return markedStatuses;
    }

    if (
      (actualPickupDate && today.isAfter(actualPickupDate, 'day')) ||
      today.isAfter(pmSurveyPlannedPickupDate, 'day') ||
      today.isAfter(requestedPickupDate, 'day')
    ) {
      markedStatuses.push('IN_TRANSIT');
    } else {
      return markedStatuses;
    }

    if (actualDeliveryDate) {
      markedStatuses.push('DELIVERED');
    } else {
      return markedStatuses;
    }
    return markedStatuses;
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

  render() {
    const statuses = this.addDates(this.statuses());
    return <StatusTimelineContainer statuses={statuses} markedStatuses={this.markedStatuses()} />;
  }
}

ShipmentStatusTimeline.propTypes = {
  shipment: PropTypes.object.isRequired,
  today: PropTypes.string,
};

export class ProfileStatusTimeline extends React.Component {
  statuses() {
    return [
      { name: 'Profile', code: 'PROFILE' },
      { name: 'Orders', code: 'ORDERS' },
      { name: 'Move Setup', code: 'MOVE_SETUP' },
      { name: 'Review', code: 'REVIEW' },
    ];
  }

  markedStatuses() {
    return ['PROFILE', 'ORDERS'];
  }

  render() {
    return <StatusTimelineContainer statuses={this.statuses()} markedStatuses={this.markedStatuses()} />;
  }
}

ProfileStatusTimeline.propTypes = {
  profile: PropTypes.object.isRequired,
};

class StatusTimelineContainer extends PureComponent {
  createStatusBlock = (status, markedStatuses) => {
    return (
      <StatusBlock
        name={status.name}
        code={status.code}
        key={status.code}
        dates={status.dates}
        completed={checkIfCompleted(markedStatuses, status.code)}
        current={checkIfCurrent(markedStatuses, status.code)}
      />
    );
  };

  render() {
    const statusBlocks = this.props.statuses.map(status => this.createStatusBlock(status, this.props.markedStatuses));

    return (
      <div className="status_timeline">
        {statusBlocks}
        <div className="legend">* Estimated</div>
      </div>
    );
  }
}

StatusTimelineContainer.propTypes = {
  statuses: PropTypes.array.isRequired,
  markedStatuses: PropTypes.array.isRequired,
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
