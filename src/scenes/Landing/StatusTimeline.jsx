import PropTypes from 'prop-types';
import React, { PureComponent } from 'react';
import { get, indexOf } from 'lodash';
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

  determineCompletedAndCurrentStatuses(shipment) {
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

  render() {
    const bookDate = get(this.props.shipment, 'book_date');
    const moveDates = this.props.shipment.move_dates_summary;
    const pickupDates = get(moveDates, 'pickup', []);
    const packDates = get(moveDates, 'pack', []);
    const deliveryDates = get(moveDates, 'delivery', []);
    const transitDates = get(moveDates, 'transit', []);

    const formatType = 'condensed';

    const markedStatuses = this.determineCompletedAndCurrentStatuses(this.props.shipment);

    return (
      <div className="status_timeline">
        <StatusBlock
          name="Scheduled"
          dates={[bookDate]}
          formatType={formatType}
          completed={true}
          current={this.checkIfCurrent(markedStatuses, 'scheduled')}
        />
        <StatusBlock
          name="Packed"
          dates={packDates}
          formatType={formatType}
          completed={this.checkIfCompleted(markedStatuses, 'packed')}
          current={this.checkIfCurrent(markedStatuses, 'packed')}
        />
        <StatusBlock
          name="Loaded"
          dates={pickupDates}
          formatType={formatType}
          completed={this.checkIfCompleted(markedStatuses, 'loaded')}
          current={this.checkIfCurrent(markedStatuses, 'loaded')}
        />
        <StatusBlock
          name="In transit"
          dates={transitDates}
          formatType={formatType}
          completed={this.checkIfCompleted(markedStatuses, 'in_transit')}
          current={this.checkIfCurrent(markedStatuses, 'in_transit')}
        />
        <StatusBlock
          name="Delivered"
          dates={deliveryDates}
          formatType={formatType}
          completed={this.checkIfCompleted(markedStatuses, 'delivered')}
          current={this.checkIfCurrent(markedStatuses, 'delivered')}
        />
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
      {props.dates && <div className="status_dates">{displayDateRange(props.dates, props.formatType)}</div>}
    </div>
  );
};
