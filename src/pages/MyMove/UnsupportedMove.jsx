import { get } from 'lodash';
import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';

export const UnsupportedMove = ({ serviceMember }) => (
  <div className="usa-grid">
    <h1 className="sm-heading">MilMove doesn’t support OCONUS moves yet</h1>
    <p>
      For now, MilMove only supports CONUS moves. Contact the {get(serviceMember.current_station, 'name', '')}{' '}
      transportation office for help setting up your move.
    </p>
  </div>
);

UnsupportedMove.propTypes = {
  serviceMember: PropTypes.shape({
    current_station: PropTypes.shape({
      name: PropTypes.string,
    }),
  }),
};

UnsupportedMove.defaultProps = {
  serviceMember: {
    current_station: {},
  },
};

function mapStateToProps(state) {
  return {
    serviceMember: get(state, 'serviceMember.currentServiceMember'),
  };
}

export default connect(mapStateToProps)(UnsupportedMove);
