import PropTypes from 'prop-types';
import React, { PureComponent } from 'react';
import { get } from 'lodash';

import { displayDateRange } from 'shared/formatters';
import './StatusTimeline.css';

export class StatusTimelineContainer extends PureComponent {
  render() {
    const moveDates = this.props.moveDates;
    const pickupDates = get(moveDates, 'pickup', []);
    const packDates = get(moveDates, 'pack', []);
    const deliveryDates = get(moveDates, 'delivery', []);
    const transitDates = get(moveDates, 'transit', []);
    const formatType = 'condensed';

    return (
      <div className="status_timeline">
        <StatusBlock name="Scheduled" dates={[this.props.bookDate]} formatType={formatType} completed={true} />
        <StatusBlock name="Packed" dates={packDates} formatType="condensed" current={true} />
        <StatusBlock name="Loaded" dates={pickupDates} formatType="condensed" />
        <StatusBlock name="In transit" dates={transitDates} formatType="condensed" />
        <StatusBlock name="Delivered" dates={deliveryDates} formatType="condensed" />
        <div className="legend">* Estimated</div>
      </div>
    );
  }
}

StatusTimelineContainer.propTypes = {
  bookDate: PropTypes.string,
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
