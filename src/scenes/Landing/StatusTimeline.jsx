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

  determineCompletedAndCurrentShipmentStatuses(shipment) {
    const today = moment();
    const actualPackDate = get(shipment, 'actual_pack_date', null);
    const pmSurveyPlannedPackDate = get(shipment, 'pm_survey_planned_pack_date', null);
    const originalPackDate = get(shipment, 'original_pack_date', null);
    const actualPickupDate = get(shipment, 'actual_pickup_date', null);
    const pmSurveyPlannedPickupDate = get(shipment, 'pm_survey_planned_pickup_date', null);
    const requestedPickupDate = get(shipment, 'requested_pickup_date', null);
    const actualDeliveryDate = get(shipment, 'actual_delivery_date', null);
    let markedStatuses = ['scheduled'];

    if (actualPackDate || today.isSameOrAfter(pmSurveyPlannedPackDate) || today.isSameOrAfter(originalPackDate)) {
      markedStatuses.push('packed');
    } else {
      return markedStatuses;
    }

    if (
      actualPickupDate ||
      today.isSameOrAfter(pmSurveyPlannedPickupDate, 'day') ||
      today.isSameOrAfter(requestedPickupDate, 'day')
    ) {
      markedStatuses.push('loaded');
    } else {
      return markedStatuses;
    }

    if (
      (actualPickupDate && today.isAfter(actualPickupDate, 'day')) ||
      today.isAfter(pmSurveyPlannedPickupDate, 'day') ||
      today.isAfter(requestedPickupDate, 'day')
    ) {
      markedStatuses.push('in_transit');
    } else {
      return markedStatuses;
    }

    if (actualDeliveryDate) {
      markedStatuses.push('delivered');
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
    const HHGSTATUSES = [
      { name: 'Scheduled', dates: 'book_date' },
      { name: 'Packed', dates: 'pack' },
      { name: 'Loaded', dates: 'pickup' },
      { name: 'In transit', dates: 'transit' },
      { name: 'Delivered', dates: 'delivery' },
    ];
    const statusBlocks = [];
    const formatType = 'condensed';
    const markedStatuses = this.determineCompletedAndCurrentShipmentStatuses(this.props.shipment);

    const createStatusBlocks = status => {
      statusBlocks.push(
        <StatusBlock
          name={status.name}
          dates={this.getDates(this.props.shipment, status.dates)}
          formatType={formatType}
          completed={this.checkIfCompleted(markedStatuses, status.name.toLowerCase())}
          current={this.checkIfCurrent(markedStatuses, status.name.toLowerCase())}
        />,
      );
    };
    createStatusBlocks.bind(this);
    HHGSTATUSES.forEach(createStatusBlocks);

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
  moveDates: PropTypes.object,
};

export default StatusTimelineContainer;

const StatusBlock = props => {
  let classes = ['status_block', props.name.toLowerCase()];
  if (props.completed) classes.push('status_completed');
  if (props.current) classes.push('status_current');

  return (
    <div className={classes.join(' ')}>
      <div className="status_dot" />
      <div className="status_name">{props.name}</div>
      <div className="status_dates">
        {isEmpty(props.dates) ? 'TBD' : displayDateRange(props.dates, props.formatType)}
      </div>
    </div>
  );
};
