import React from 'react';
import PropTypes from 'prop-types';
import { Navigate } from 'react-router-dom';
import { connect } from 'react-redux';

import { useTitle } from 'hooks/custom';

export function userIsAuthorized(activeRole, requiredRoles) {
  // Return true if no roles are required
  if (!requiredRoles || !requiredRoles.length) return true;

  // Return false if user has no role
  if (!activeRole) return false;

  // User active role must be defined in requiredRoles
  return requiredRoles.includes(activeRole);
}

function PrivateRoute({ requiredRoles, activeRole, children }) {
  useTitle();

  if (!userIsAuthorized(activeRole, requiredRoles)) {
    return <Navigate to="/invalid-permissions" />;
  }

  return children;
}

PrivateRoute.displayName = 'PrivateRoute';

PrivateRoute.propTypes = {
  requiredRoles: PropTypes.arrayOf(PropTypes.string),
  activeRole: PropTypes.string,
};

PrivateRoute.defaultProps = {
  requiredRoles: [],
  activeRole: null,
};

const mapStateToProps = (state) => {
  return {
    activeRole: state.auth.activeRole,
  };
};

export default connect(mapStateToProps)(PrivateRoute);
