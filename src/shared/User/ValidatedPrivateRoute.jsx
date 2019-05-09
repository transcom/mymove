import React from 'react';
import { Route } from 'react-router-dom';
import { connect } from 'react-redux';

import { selectCurrentUser } from 'shared/Data/users';
import { get } from 'lodash';
import SignIn from './SignIn';

// this was adapted from https://github.com/ReactTraining/react-router/blob/master/packages/react-router-redux/examples/AuthExample.js
// note that it does not work if the route is not inside a Switch
class ValidatedPrivateRouteContainer extends React.Component {
  render() {
    const { isLoggedIn, requiresAccessCode, accessCode, path, ...props } = this.props;
    console.log('Requires access code', requiresAccessCode);
    console.log('Access code', accessCode);
    if (isLoggedIn && (!requiresAccessCode || (requiresAccessCode && accessCode !== undefined)))
      return <Route {...props} />;
    else return <Route path={path} component={SignIn} />;
  }
}
const mapStateToProps = state => {
  const user = selectCurrentUser(state);
  const serviceMember = get(user, 'service_member');
  return {
    isLoggedIn: user.isLoggedIn,
    requiresAccessCode: get(serviceMember, 'requires_access_code'),
    accessCode: get(serviceMember, 'access_code'),
  };
};
const ValidatedPrivateRoute = connect(mapStateToProps)(ValidatedPrivateRouteContainer);

export default ValidatedPrivateRoute;
