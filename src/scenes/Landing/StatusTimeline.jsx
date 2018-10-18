import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { withRouter } from 'react-router-dom';

import { getMoveDatesSummary, selectMoveDatesSummary } from 'shared/Entities/modules/moves';
import { bindActionCreators } from 'redux';
import connect from 'react-redux/es/connect/connect';
import { get } from 'lodash';
import { displayDateRange } from '../Moves/Hhg/DatesSummary';

const getRequestLabel = 'StatusTimeline.getMoveDatesSummary';

export class StatusTimeline extends Component {
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
        <div className="status_block">
          <div className="status_dot status_active" />
          <div className="status_name">Scheduled</div>
          <div>date requested</div>
        </div>
        <div className="status_block">
          <div className="status_dot" />
          <div className="status_name">Packed</div>
          <div>{displayDateRange(packDates, formatType)}</div>
        </div>
        <div className="status_block">
          <div className="status_dot" />
          <div className="status_name">Loaded</div>
          <div>{displayDateRange(pickupDates, formatType)}</div>
        </div>
        <div className="status_block">
          <div className="status_dot" />
          <div className="status_name">In transit</div>
          <div>{displayDateRange(transitDates, formatType)}</div>
        </div>
        <div className="status_block">
          <div className="status_dot" />
          <div className="status_name">Delivered</div>
          <div>{displayDateRange(deliveryDates, formatType)}</div>
        </div>
      </div>
    );
  }
}

StatusTimeline.propTypes = {
  moveDate: PropTypes.PropTypes.string,
  moveId: PropTypes.string,
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

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(StatusTimeline));
