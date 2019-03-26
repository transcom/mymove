import PropTypes from 'prop-types';
import React, { PureComponent } from 'react';
import { get, indexOf, isEmpty } from 'lodash';
import moment from 'moment';

import { displayDateRange } from 'shared/formatters';
import './StatusTimeline.css';

export class StatusTimelineContainer extends PureComponent {
  checkIfCompleted(statuses, statusToCheck) {
    if (!statuses.includes(statusToCheck)) {
      return false;
    }
    return indexOf(statuses, statusToCheck) < statuses.length;
  }

  checkIfCurrent(statuses, statusToCheck) {
    if (!statuses.includes(statusToCheck)) {
      return false;
    }
    return indexOf(statuses, statusToCheck) === statuses.length - 1;
  }

  determineCompletedAndCurrentProfileStatuses() {
    return ['PROFILE', 'ORDERS'];
  }

  determineCompletedAndCurrentPPMStatuses(ppm) {
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

  determineCompletedAndCurrentShipmentStatuses(shipment) {
    const today = moment();
    const actualPackDate = get(shipment, 'actual_pack_date', null);
    const pmSurveyPlannedPackDate = get(shipment, 'pm_survey_planned_pack_date', null);
    const originalPackDate = get(shipment, 'original_pack_date', null);
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

  getDates(shipment, datesType) {
    if (datesType === 'book_date') {
      return [get(shipment, datesType)];
    }
    return get(shipment, datesType);
  }

  render() {
    const PPMSTATUSES = [
      { name: 'Submitted', code: 'SUBMITTED', dates: null },
      { name: 'Approved', code: 'PPM_APPROVED', dates: null },
      { name: 'In progress', code: 'IN_PROGRESS', dates: null },
      { name: 'Completed', code: 'COMPLETED', dates: null },
    ];
    const PROFILESTATUSES = [
      { name: 'Profile', code: 'PROFILE', dates: null },
      { name: 'Orders', code: 'ORDERS', dates: null },
      { name: 'Move Setup', code: 'MOVE_SETUP', dates: null },
      { name: 'Review', code: 'REVIEW', dates: null },
    ];
    const HHGSTATUSES = [
      { name: 'Scheduled', code: 'SCHEDULED', dates: 'book_date' },
      { name: 'Packed', code: 'PACKED', dates: 'pack' },
      { name: 'Loaded', code: 'LOADED', dates: 'pickup' },
      { name: 'In transit', code: 'IN_TRANSIT', dates: 'transit' },
      { name: 'Delivered', code: 'DELIVERED', dates: 'delivery' },
    ];
    const statusBlocks = [];
    const formatType = 'condensed';

    let statuses = HHGSTATUSES;
    if (this.props.ppm) {
      statuses = PPMSTATUSES;
    } else if (this.props.profile) {
      statuses = PROFILESTATUSES;
    }
    let markedStatuses = [];
    if (this.props.ppm) {
      markedStatuses = this.determineCompletedAndCurrentPPMStatuses(this.props.ppm);
    } else if (this.props.profile) {
      markedStatuses = this.determineCompletedAndCurrentProfileStatuses(this.props.profile);
    } else {
      markedStatuses = this.determineCompletedAndCurrentShipmentStatuses(this.props.shipment);
    }

    const createStatusBlocks = status => {
      statusBlocks.push(
        <StatusBlock
          name={status.name}
          dates={this.getDates(this.props.shipment, status.dates)}
          formatType={formatType}
          completed={this.checkIfCompleted(markedStatuses, status.code)}
          current={this.checkIfCurrent(markedStatuses, status.code)}
          shipment={this.props.shipment}
          code={status.code}
        />,
      );
    };
    createStatusBlocks.bind(this);
    statuses.forEach(createStatusBlocks);

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
};

export default StatusTimelineContainer;

const StatusBlock = props => {
  let classes = ['status_block', props.code.toLowerCase()];
  if (props.completed) classes.push('status_completed');
  if (props.current) classes.push('status_current');

  return (
    <div className={classes.join(' ')}>
      <div className="status_dot" />
      <div className="status_name">{props.name}</div>
      {props.shipment && (
        <div className="status_dates">
          {isEmpty(props.dates) ? 'TBD' : displayDateRange(props.dates, props.formatType)}
        </div>
      )}
    </div>
  );
};
