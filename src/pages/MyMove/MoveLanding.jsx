import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';

import { selectServiceMemberFromLoggedInUser } from 'shared/Entities/modules/serviceMembers';

export const MoveLanding = ({ serviceMember }) => {
  return (
    <div className="usa-grid">
      <h1 className="sm-heading">Home</h1>
      {/* eslint-disable-next-line camelcase */}
      <h2>Welcome {serviceMember?.first_name}</h2>
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
    serviceMember: selectServiceMemberFromLoggedInUser(state),
  };
}

export default connect(mapStateToProps)(MoveLanding);
