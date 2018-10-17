import PropTypes from 'prop-types';
import React, { Component } from 'react';

import { getMoveDatesSummary } from 'shared/Entities/modules/moves';
import { bindActionCreators } from 'redux';
import connect from 'react-redux/es/connect/connect';

const getRequestLabel = 'statusTimeline.getMoveDatesSummary';

export class StatusTimeline extends Component {
  componentDidMount() {
    getMoveDatesSummary(getRequestLabel, 'moveId', this.props.moveDate);
  }

  render() {
    return (
      <div className="status_timeline">
        <div className="status_block">
          <div className="status_dot status_active" />
          <div className="status_name">Scheduled</div>
          <div>Date</div>
        </div>
        <div className="status_block">
          <div className="status_dot" />
          <div className="status_name">Packed</div>
          <div>Date</div>
        </div>
        <div className="status_block">
          <div className="status_dot" />
          <div className="status_name">Loaded</div>
          <div>date</div>
        </div>
        <div className="status_block">
          <div className="status_dot" />
          <div className="status_name">In transit</div>
          <div>date range</div>
        </div>
        <div className="status_block">
          <div className="status_dot" />
          <div className="status_name">Delivered</div>
          <div>date</div>
        </div>
      </div>
    );
  }
}

StatusTimeline.propTypes = {
  getMoveDatesSummary: PropTypes.object,
  moveDate: PropTypes.date,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ getMoveDatesSummary }, dispatch);
}

export default connect(() => ({}), mapDispatchToProps)(StatusTimeline);
