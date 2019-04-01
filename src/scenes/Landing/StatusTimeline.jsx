import PropTypes from 'prop-types';
import React, { PureComponent } from 'react';
import { get, indexOf } from 'lodash';
import moment from 'moment';

import { displayDateRange } from 'shared/formatters';
import './StatusTimeline.css';

class PPM {
  constructor(ppm) {
    this.ppm = ppm;
    this.statuses = [
      { name: 'Submitted', code: 'SUBMITTED', dates: null },
      { name: 'Approved', code: 'PPM_APPROVED', dates: null },
      { name: 'In progress', code: 'IN_PROGRESS', dates: null },
      { name: 'Completed', code: 'COMPLETED', dates: null },
    ];
    this.markedStatuses = this.determineStatuses();
  }

  determineStatuses() {
    let markedStatuses = ['SUBMITTED'];

    if (this.ppm.status === 'SUBMITTED') {
      return markedStatuses;
    }

    markedStatuses.push('PPM_APPROVED');
    if (this.ppm.status === 'APPROVED') {
      const moveInProgress = moment(this.ppm.original_move_date, 'YYYY-MM-DD').isSameOrBefore();
      if (moveInProgress) {
        markedStatuses.push('IN_PROGRESS');
      }
      return markedStatuses;
    }

    if (this.ppm.status === 'COMPLETED') {
      markedStatuses.push('IN_PROGRESS');
      markedStatuses.push('COMPLETED');
    }

    return markedStatuses;
  }
}

class Profile {
  constructor(profile) {
    this.profile = profile;
    this.statuses = [
      { name: 'Profile', code: 'PROFILE', dates: null },
      { name: 'Orders', code: 'ORDERS', dates: null },
      { name: 'Move Setup', code: 'MOVE_SETUP', dates: null },
      { name: 'Review', code: 'REVIEW', dates: null },
    ];
    this.markedStatuses = this.determineStatuses();
  }

  determineStatuses() {
    return ['PROFILE', 'ORDERS'];
  }
}

class Shipment {
  constructor(shipment, today) {
    this.shipment = shipment;
    this.today = today;
    this.statuses = [
      { name: 'Scheduled', code: 'SCHEDULED', dates: 'book_date' },
      { name: 'Packed', code: 'PACKED', dates: 'pack' },
      { name: 'Loaded', code: 'LOADED', dates: 'pickup' },
      { name: 'In transit', code: 'IN_TRANSIT', dates: 'transit' },
      { name: 'Delivered', code: 'DELIVERED', dates: 'delivery' },
    ];
    this.markedStatuses = this.determineStatuses(today);
  }

  determineStatuses(today) {
    const actualPackDate = get(this.shipment, 'actual_pack_date', null);
    const originalPackDate = get(this.shipment, 'original_pack_date', null);
    const pmSurveyPlannedPackDate = get(this.shipment, 'pm_survey_planned_pack_date', null);
    const actualPickupDate = get(this.shipment, 'actual_pickup_date', null);
    const pmSurveyPlannedPickupDate = get(this.shipment, 'pm_survey_planned_pickup_date', null);
    const requestedPickupDate = get(this.shipment, 'requested_pickup_date', null);
    const actualDeliveryDate = get(this.shipment, 'actual_delivery_date', null);
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
}

export class StatusTimelineContainer extends PureComponent {
  static checkIfCompleted(statuses, statusToCheck) {
    if (!statuses.includes(statusToCheck)) {
      return false;
    }
    return indexOf(statuses, statusToCheck) < statuses.length;
  }

  static checkIfCurrent(statuses, statusToCheck) {
    if (!statuses.includes(statusToCheck)) {
      return false;
    }
    return indexOf(statuses, statusToCheck) === statuses.length - 1;
  }

  static getDates(shipment, datesType) {
    if (datesType === 'book_date') {
      return [get(shipment, datesType)];
    }
    return get(shipment, datesType);
  }

  determineTimelineType() {
    if (this.props.ppm) {
      return new PPM(this.props.ppm);
    }
    if (this.props.profile) {
      return new Profile(this.props.profile);
    }
    const today = this.props.today ? moment(this.props.today) : moment();
    return new Shipment(this.props.shipment, today);
  }

  createStatusBlock = (status, markedStatuses) => {
    return (
      <StatusBlock
        name={status.name}
        code={status.code}
        key={status.code}
        shipment={this.props.shipment}
        dates={StatusTimelineContainer.getDates(this.props.shipment, status.dates)}
        formatType="condensed"
        completed={StatusTimelineContainer.checkIfCompleted(markedStatuses, status.code)}
        current={StatusTimelineContainer.checkIfCurrent(markedStatuses, status.code)}
      />
    );
  };

  render() {
    const timeline = this.determineTimelineType();
    const statusBlocks = timeline.statuses.map(status => this.createStatusBlock(status, timeline.markedStatuses));

    return (
      <div className="status_timeline">
        {statusBlocks}
        <div className="legend">* Estimated</div>
      </div>
    );
  }
}

StatusTimelineContainer.propTypes = {
  shipment: PropTypes.object,
  ppm: PropTypes.object,
  profile: PropTypes.object,
  today: PropTypes.string,
};

export default StatusTimelineContainer;

export const StatusBlock = props => {
  const classes = ['status_block', props.code.toLowerCase()];
  if (props.completed) classes.push('status_completed');
  if (props.current) classes.push('status_current');

  return (
    <div className={classes.join(' ')}>
      <div className="status_dot" />
      <div className="status_name">{props.name}</div>
      {props.shipment &&
        props.dates && <div className="status_dates">{displayDateRange(props.dates, props.formatType)}</div>}
    </div>
  );
};
