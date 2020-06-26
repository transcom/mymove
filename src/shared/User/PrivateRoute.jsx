import React from 'react';
import { Route } from 'react-router-dom';
import { connect } from 'react-redux';

import { selectCurrentUser, selectGetCurrentUserIsLoading } from 'shared/Data/users';
import SignIn from './SignIn';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import { Redirect, Link } from 'react-router-dom';

export function userIsAuthorized(userRoles, requiredRoles) {
  // Return true if no roles are required
  if (!requiredRoles || !requiredRoles.length) return true;

  // Return false if user has no roles
  if (!userRoles || !userRoles.length) return false;

  // User must have at least one of the roles defined in requiredRoles
  return !!userRoles?.find((r) => requiredRoles.indexOf(r) > -1);
}

const PrivateRouteContainer = (props) => {
  const { loginIsLoading, userIsLoggedIn, path, requiredRoles, userRoles, hideSwitcher, ...routeProps } = props;

  if (
    userIsLoggedIn &&
    userIsAuthorized(
      userRoles.map((role) => role.roleType),
      requiredRoles,
    )
  ) {
    // User is logged in & authorized to view the requested URL
    // If user has multiple roles, add a link to let them select which role they are using
    // TODO improve this UI

    const displaySelectApplication =
      !hideSwitcher && userRoles?.length > 1 && routeProps.location?.pathname !== '/select-application';
    return displaySelectApplication ? (
      <>
        <Link to="/select-application">Change user role</Link>
        <Route path={path} {...routeProps} />
      </>
    ) : (
      <Route path={path} {...routeProps} />
    );
  } else if (userIsLoggedIn)
    // User is logged in but not authorized to view the requested URL, redirect home
    return <Redirect to="/" />;
  else if (loginIsLoading)
    // User is still loading
    return <LoadingPlaceholder />;
  // User is not logged in, go to Sign In page
  else return <Route path={path} component={SignIn} />; // TODO - change this to a redirect
};

const mapStateToProps = (state) => ({
  loginIsLoading: selectGetCurrentUserIsLoading(state),
  userIsLoggedIn: selectCurrentUser(state).isLoggedIn,
  userRoles: selectCurrentUser(state).roles,
});

const PrivateRoute = connect(mapStateToProps)(PrivateRouteContainer);

export default PrivateRoute;
