import React from 'react';
import { Route } from 'react-router-dom';
import { connect } from 'react-redux';

import { selectCurrentUser, selectGetCurrentUserIsLoading } from 'shared/Data/users';
import SignIn from './SignIn';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import { Redirect } from 'react-router-dom';

// this was adapted from https://github.com/ReactTraining/react-router/blob/master/packages/react-router-redux/examples/AuthExample.js
// note that it does not work if the route is not inside a Switch
class PrivateRouteContainer extends React.Component {
  render() {
    const { loginIsLoading, userIsLoggedIn, path, requiredRole, userRoles, ...props } = this.props;

    if (userIsLoggedIn && userIsAuthorized(userRoles, requiredRole)) return <Route {...props} />;
    else if (userIsLoggedIn) return <Redirect exact to={redirectURLForRole(userRoles[0].roleType)} />;
    else if (loginIsLoading) return <LoadingPlaceholder />;
    else return <Route path={path} component={SignIn} />;
  }
}
const mapStateToProps = (state) => ({
  loginIsLoading: selectGetCurrentUserIsLoading(state),
  userIsLoggedIn: selectCurrentUser(state).isLoggedIn,
  userRoles: selectCurrentUser(state).roles,
});

function userIsAuthorized(userRoles, requiredRole) {
  return Boolean(userRoles?.map((role) => role.roleType)?.includes(requiredRole));
}

function redirectURLForRole(role) {
  switch (role) {
    case 'ppm_office_users':
      return '/queues/new';
    case 'transportation_ordering_officer':
      return '/moves/queue';
    case 'transportation_invoicing_officer':
      return '/tio/placeholder';
    default:
      return '/';
  }
}

const PrivateRoute = connect(mapStateToProps)(PrivateRouteContainer);

export default PrivateRoute;
