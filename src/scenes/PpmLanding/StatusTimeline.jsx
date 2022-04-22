import PropTypes from 'prop-types';
import React, { PureComponent } from 'react';
import { filter, findLast } from 'lodash';

import { displayDateRange } from 'utils/formatters';
import './StatusTimeline.scss';

function getCurrentStatus(statuses) {
  return findLast(statuses, function (status) {
    return status.completed;
  });
}

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

export class StatusTimeline extends PureComponent {
  createStatusBlock = (status, currentStatus) => {
    return (
      <StatusBlock
        name={status.name}
        code={status.code}
        key={status.code}
        dates={filter(status.dates, (date) => {
          return date;
        })}
        completed={status.completed}
        current={currentStatus.code === status.code}
      />
    );
  };

  render() {
    const currentStatus = getCurrentStatus(this.props.statuses);
    const statusBlocks = this.props.statuses.map((status) => this.createStatusBlock(status, currentStatus));

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

export const StatusBlock = (props) => {
  const classes = ['status_block', props.code.toLowerCase()];
  if (props.completed) classes.push('status_completed');
  if (props.current) classes.push('status_current');

  return (
    <div className={classes.join(' ')}>
      <div className="status_dot" />
      <div className="status_name">{props.name}</div>
      {props.dates && props.dates.length > 0 && (
        <div className="status_dates">{displayDateRange(props.dates, 'condensed')}</div>
      )}
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
