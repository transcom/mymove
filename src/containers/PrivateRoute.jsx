import React from 'react';
import PropTypes from 'prop-types';
import { Route, Redirect } from 'react-router-dom';
import { connect } from 'react-redux';

import { selectCurrentUser, selectGetCurrentUserIsLoading } from 'shared/Data/users';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import { UserRolesShape } from 'types/index';

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

  if (loginIsLoading) return <LoadingPlaceholder />;

  if (!userIsLoggedIn) return <Redirect to="/sign-in" />;
  if (
    !userIsAuthorized(
      userRoles.map((r) => r.roleType),
      requiredRoles,
    )
  )
    return <Redirect to="/" />;

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

const mapStateToProps = (state) => ({
  loginIsLoading: selectGetCurrentUserIsLoading(state),
  userIsLoggedIn: selectCurrentUser(state).isLoggedIn,
  userRoles: selectCurrentUser(state).roles,
});

const ConnectedPrivateRoute = connect(mapStateToProps)(PrivateRoute);

export default ConnectedPrivateRoute;
