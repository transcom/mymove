import PropTypes from 'prop-types';
import React, { PureComponent } from 'react';
import { get, indexOf } from 'lodash';
import moment from 'moment';

import { displayDateRange } from 'shared/formatters';
import './StatusTimeline.css';

export class StatusTimelineContainer extends PureComponent {
  checkIfCompleted(statuses, statusToCheck) {
    return indexOf(statuses, statusToCheck) < statuses.length;
  }
  checkIfCurrent(statuses, statusToCheck) {
    return indexOf(statuses, statusToCheck) === statuses.length - 1;
  }

  render() {
    const moveDates = this.props.moveDates;
    const pickupDates = get(moveDates, 'pickup', []);
    const packDates = get(moveDates, 'pack', []);
    const deliveryDates = get(moveDates, 'delivery', []);
    const transitDates = get(moveDates, 'transit', []);
    const formatType = 'condensed';

    const shipment = this.props.shipment;
    // pack dates
    const actualPackDate = get(shipment, 'actual_pack_date');
    const pmSurveyPlannedPackDate = get(shipment, 'pm_survey_planned_pack_date');
    // pickup dates
    const actualPickupDate = get(shipment, 'actual_pickup_date');
    const pmSurveyPlannedPickupDate = get(shipment, 'pm_survey_planned_pickup_date');

    const actualDeliveryDate = get(shipment, 'actual_delivery_date');
    // const pmSurveyDeliveryDate = get(shipment, 'pm_survey_planned_delivery_date');

    // Create an array to push in completed UI statuses (current at -1 index)
    let markedStatuses = ['scheduled'];
    // define each UI status and date range if needed
    if (actualPackDate || moment().isSameOrAfter(pmSurveyPlannedPackDate)) {
      markedStatuses.push('packed');
    }
    if (actualPickupDate || moment().isSameOrAfter(pmSurveyPlannedPickupDate, 'day')) {
      markedStatuses.push('loaded');
    }
    if (
      (actualPickupDate && moment().isAfter(actualPickupDate, 'day')) ||
      moment().isAfter(pmSurveyPlannedPickupDate, 'day')
    ) {
      markedStatuses.push('in_transit');
    }
    if (actualDeliveryDate) {
      markedStatuses.push('delivered');
    }

    return (
      <div className="status_timeline">
        <StatusBlock
          name="Scheduled"
          dates={[this.props.bookDate]}
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
  bookDate: PropTypes.string,
  moveDatesSummary: PropTypes.object,
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
