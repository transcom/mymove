import React from 'react';
import { Route } from 'react-router-dom';
import { connect } from 'react-redux';

import { selectCurrentUser } from 'shared/Data/users';
import { get } from 'lodash';
import SignIn from './SignIn';
import AccessCode from './AccessCode';

// this was adapted from https://github.com/ReactTraining/react-router/blob/master/packages/react-router-redux/examples/AuthExample.js
// note that it does not work if the route is not inside a Switch
class ValidatedPrivateRouteContainer extends React.Component {
  render() {
    const { isLoggedIn, requiresAccessCode, hasMove, path, ...props } = this.props;
    if (!isLoggedIn) return <Route path={path} component={SignIn} />;
    if (isLoggedIn && requiresAccessCode && !hasMove) return <Route path={path} component={AccessCode} />;
    return <Route {...props} />;
  }
}
const mapStateToProps = state => {
  const user = selectCurrentUser(state);
  const serviceMember = get(state, 'serviceMember.currentServiceMember');
  console.log(state);
  return {
    isLoggedIn: user.isLoggedIn,
    requiresAccessCode: get(serviceMember, 'requires_access_code'),
    hasMove: get(state, 'moves.currentMove'),
  };
};
const ValidatedPrivateRoute = connect(mapStateToProps)(ValidatedPrivateRouteContainer);

export default ValidatedPrivateRoute;
