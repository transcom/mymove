import React from 'react';
import { Route } from 'react-router-dom';
import { connect } from 'react-redux';

import { selectCurrentUser, selectGetCurrentUserIsLoading } from 'shared/Data/users';
import SignIn from './SignIn';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import { Redirect, Link } from 'react-router-dom';

// this was adapted from https://github.com/ReactTraining/react-router/blob/master/packages/react-router-redux/examples/AuthExample.js
// note that it does not work if the route is not inside a Switch
class PrivateRouteContainer extends React.Component {
  render() {
    const { loginIsLoading, userIsLoggedIn, path, requiredRoles, userRoles, ...props } = this.props;

    if (
      userIsLoggedIn &&
      userIsAuthorized(
        userRoles.map((role) => role.roleType),
        requiredRoles,
      )
    )
      return wrapRouteForMultipleRoles(<Route {...props} />, userRoles);
    else if (userIsLoggedIn) return <Redirect exact to={redirectURLForRole(userRoles[0].roleType)} />;
    else if (loginIsLoading) return <LoadingPlaceholder />;
    else return <Route path={path} component={SignIn} />; // TODO - change this to a redirect
  }
}
const mapStateToProps = (state) => ({
  loginIsLoading: selectGetCurrentUserIsLoading(state),
  userIsLoggedIn: selectCurrentUser(state).isLoggedIn,
  userRoles: selectCurrentUser(state).roles,
});

function userIsAuthorized(userRoles, requiredRoles) {
  return !!userRoles.find((r) => requiredRoles.indexOf(r) > -1);
}

function redirectURLForRole(role) {
  switch (role) {
    case 'ppm_office_users':
      return '/queues/new';
    case 'transportation_ordering_officer':
      return '/moves/queue';
    case 'transportation_invoicing_officer':
      return '/invoicing/queue';
    default:
      return '/';
  }
}

function wrapRouteForMultipleRoles(route, userRoles) {
  if (userRoles?.length > 1 && route.props.location.pathname !== '/select-application') {
    return (
      <div>
        <Link to="/select-application">Select application</Link>
        {route}
      </div>
    );
  } else {
    return route;
  }
}

const PrivateRoute = connect(mapStateToProps)(PrivateRouteContainer);

export default PrivateRoute;
