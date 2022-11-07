import React from 'react';
import PropTypes from 'prop-types';
import { Route, Redirect } from 'react-router-dom';
import { connect } from 'react-redux';

import { selectGetCurrentUserIsLoading, selectIsLoggedIn } from 'store/auth/selectors';
import { selectLoggedInUser } from 'store/entities/selectors';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
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

const PrivateRoute = (props) => {
  const { loginIsLoading, userIsLoggedIn, requiredRoles, userRoles, ...routeProps } = props;
  const userRoleTypes = getRoleTypesFromRoles(userRoles);
  if (loginIsLoading) return <LoadingPlaceholder />;

  if (!userIsLoggedIn) return <Redirect to="/sign-in" />;
  if (!userIsAuthorized(userRoleTypes, requiredRoles)) {
    return <Redirect to="/invalid-permissions" />;
  }

  // eslint-disable-next-line react/jsx-props-no-spreading
  return <Route {...routeProps} />;
};

PrivateRoute.displayName = 'PrivateRoute';

PrivateRoute.propTypes = {
  loginIsLoading: PropTypes.bool,
  userIsLoggedIn: PropTypes.bool,
  requiredRoles: PropTypes.arrayOf(PropTypes.string),
  userRoles: UserRolesShape,
};

PrivateRoute.defaultProps = {
  loginIsLoading: true,
  userIsLoggedIn: false,
  requiredRoles: [],
  userRoles: [],
};

const mapStateToProps = (state) => {
  const user = selectLoggedInUser(state);

  return {
    loginIsLoading: selectGetCurrentUserIsLoading(state),
    userIsLoggedIn: selectIsLoggedIn(state),
    userRoles: user?.roles || [],
  };
};

const ConnectedPrivateRoute = connect(mapStateToProps)(PrivateRoute);

export default ConnectedPrivateRoute;
