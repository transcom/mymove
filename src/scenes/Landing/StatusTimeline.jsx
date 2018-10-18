import PropTypes from 'prop-types';
import React, { PureComponent } from 'react';
import { withRouter } from 'react-router-dom';
import connect from 'react-redux/es/connect/connect';
import { get } from 'lodash';
import { bindActionCreators } from 'redux';

import { getMoveDatesSummary, selectMoveDatesSummary } from 'shared/Entities/modules/moves';
import { displayDateRange } from '../Moves/Hhg/DatesSummary';
import './StatusTimeline.css';

const getRequestLabel = 'StatusTimeline.getMoveDatesSummary';

export class StatusTimelineContainer extends PureComponent {
  componentDidMount() {
    this.props.getMoveDatesSummary(getRequestLabel, this.props.moveId, this.props.moveDate);
  }

  render() {
    const moveDates = this.props.moveDatesSummary;
    const pickupDates = get(moveDates, 'pickup', []);
    const packDates = get(moveDates, 'pack', []);
    const deliveryDates = get(moveDates, 'delivery', []);
    const transitDates = get(moveDates, 'transit', []);
    const formatType = 'condensed';
    return (
      <div className="status_timeline">
        <StatusBlock name="Scheduled" dates={[this.props.bookDate]} formatType={formatType} active={true} />
        <StatusBlock name="Packed" dates={packDates} formatType="condensed" />
        <StatusBlock name="Loaded" dates={pickupDates} formatType="condensed" />
        <StatusBlock name="In transit" dates={transitDates} formatType="condensed" />
        <StatusBlock name="Delivered" dates={deliveryDates} formatType="condensed" />
      </div>
    );
  }
}

StatusTimelineContainer.propTypes = {
  moveDate: PropTypes.PropTypes.string,
  moveId: PropTypes.string,
  bookDate: PropTypes.string,
  moveDatesSummary: PropTypes.object,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ getMoveDatesSummary }, dispatch);
}

function mapStateToProps(state, ownProps) {
  const moveDate = get(ownProps, 'moveDate');
  const moveDatesSummary = selectMoveDatesSummary(state, ownProps.moveId, moveDate);
  return {
    moveDatesSummary: moveDatesSummary,
  };
}

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(StatusTimelineContainer));

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
