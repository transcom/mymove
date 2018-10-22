import PropTypes from 'prop-types';
import React, { PureComponent } from 'react';
import { get } from 'lodash';

import { getMoveDatesSummary, selectMoveDatesSummary } from 'shared/Entities/modules/moves';
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
        <StatusBlock name="Scheduled" dates={[this.props.bookDate]} formatType={formatType} active={true} />
        <div className="status_line" />
        <StatusBlock name="Packed" dates={packDates} formatType="condensed" />
        <div className="status_line" />
        <StatusBlock name="Loaded" dates={pickupDates} formatType="condensed" />
        <div className="status_line" />
        <StatusBlock name="In transit" dates={transitDates} formatType="condensed" />
        <div className="status_line" />
        <StatusBlock name="Delivered" dates={deliveryDates} formatType="condensed" />
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
  let active = props.active ? 'status_active' : '';
  return (
    <div className="status_block">
      <div className={`status_dot ${active}`} />
      <div className="status_name">{props.name}</div>
      {props.dates && <div>{displayDateRange(props.dates, props.formatType)}</div>}
    </div>
  );
};
