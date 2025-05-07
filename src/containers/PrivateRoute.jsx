import React from 'react';
import PropTypes from 'prop-types';
import { Navigate } from 'react-router-dom';
import { connect } from 'react-redux';

import { selectLoggedInUser } from 'store/entities/selectors';
import { useTitle } from 'hooks/custom';
import { UserRoleShape } from 'types';

export function userIsAuthorized(userActiveRole, requiredRoles) {
  // Return true if no roles are required
  if (!requiredRoles || !requiredRoles.length) return true;

  // Return false if user has no role
  if (!userActiveRole) return false;

  // User active role must be defined in requiredRoles
  return requiredRoles.includes(userActiveRole);
}

function PrivateRoute({ requiredRoles, userActiveRole, children }) {
  useTitle();

  if (!userIsAuthorized(userActiveRole, requiredRoles)) {
    return <Navigate to="/invalid-permissions" />;
  }

  return children;
}

PrivateRoute.displayName = 'PrivateRoute';

PrivateRoute.propTypes = {
  requiredRoles: PropTypes.arrayOf(PropTypes.string),
  userActiveRole: UserRoleShape,
};

PrivateRoute.defaultProps = {
  requiredRoles: [],
  userActiveRole: null,
};

const mapStateToProps = (state) => {
  const user = selectLoggedInUser(state);

  return {
    userActiveRole: user?.userActiveRole,
  };
};

export default connect(mapStateToProps)(PrivateRoute);
