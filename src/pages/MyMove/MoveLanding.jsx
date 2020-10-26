import { get } from 'lodash';
import PropTypes from 'prop-types';
import React from 'react';
import { connect } from 'react-redux';

export const MoveLanding = ({ serviceMember }) => {
  return (
    <div className="usa-grid">
      <h1>Home</h1>
      <h2>Welcome {get(serviceMember, 'first_name', '')}</h2>
    </div>
  );
};

MoveLanding.propTypes = {
  serviceMember: PropTypes.shape({
    first_name: PropTypes.string,
  }),
};

MoveLanding.defaultProps = {
  serviceMember: {},
};
function mapStateToProps(state) {
  return {
    serviceMember: get(state, 'serviceMember.currentServiceMember'),
  };
}

export default connect(mapStateToProps)(MoveLanding);
