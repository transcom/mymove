import React from 'react';
import PropTypes from 'prop-types';
import { Navigate } from 'react-router-dom';
import { connect } from 'react-redux';

import { selectLoggedInUser } from 'store/entities/selectors';
import getRoleTypesFromRoles from 'utils/user';
import { UserRolesShape } from 'types';

export function userIsAuthorized(userRoles, requiredRoles) {
  // Return true if no roles are required
  if (!requiredRoles || !requiredRoles.length) return true;

  // Return false if user has no roles
  if (!userRoles || !userRoles.length) return false;

  // User must have at least one of the roles defined in requiredRoles
  return !!userRoles?.find((r) => requiredRoles.indexOf(r) > -1);
}

function ProtectedRoute({ requiredRoles, userRoles, children }) {
  const userRoleTypes = getRoleTypesFromRoles(userRoles);

  if (!userIsAuthorized(userRoleTypes, requiredRoles)) {
    return <Navigate to="/invalid-permissions" />;
  }

  return children;
}

ProtectedRoute.displayName = 'ProtectedRoute';

ProtectedRoute.propTypes = {
  requiredRoles: PropTypes.arrayOf(PropTypes.string),
  userRoles: UserRolesShape,
};

ProtectedRoute.defaultProps = {
  requiredRoles: [],
  userRoles: [],
};

const mapStateToProps = (state) => {
  const user = selectLoggedInUser(state);

  return {
    userRoles: user?.roles || [],
  };
};

const ConnectedProtectedRoute = connect(mapStateToProps)(ProtectedRoute);

export default ConnectedProtectedRoute;
