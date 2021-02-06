import React from 'react';
import { bindActionCreators } from 'redux';
import { Route } from 'react-router-dom';
import { connect } from 'react-redux';

import { selectIsLoggedIn } from 'store/entities/selectors';

import { get } from 'lodash';
import SignIn from './SignIn';
import AccessCode from './AccessCode';

import { fetchAccessCode } from 'shared/Entities/modules/accessCodes';
import { selectServiceMemberFromLoggedInUser } from 'store/entities/selectors';

// this was adapted from https://github.com/ReactTraining/react-router/blob/master/packages/react-router-redux/examples/AuthExample.js
// note that it does not work if the route is not inside a Switch
class ValidatedPrivateRouteContainer extends React.Component {
  componentDidMount() {
    this.props.fetchAccessCode();
  }

  render() {
    const { isLoggedIn, requiresAccessCode, accessCode, path, ...props } = this.props;
    if (!isLoggedIn) return <Route path={path} component={SignIn} />;
    if (isLoggedIn && requiresAccessCode && !accessCode) return <Route path={path} component={AccessCode} />;
    return <Route {...props} />;
  }
}

const mapStateToProps = (state) => {
  const serviceMember = selectServiceMemberFromLoggedInUser(state);
  const accessCodes = get(state, 'entities.accessCodes');

  return {
    isLoggedIn: selectIsLoggedIn(state),
    requiresAccessCode: serviceMember?.requires_access_code,
    accessCode: accessCodes && Object.values(accessCodes).length > 0 ? Object.values(accessCodes)[0].code : null,
  };
};

const mapDispatchToProps = (dispatch) => bindActionCreators({ fetchAccessCode }, dispatch);

const ValidatedPrivateRoute = connect(mapStateToProps, mapDispatchToProps)(ValidatedPrivateRouteContainer);

export default ValidatedPrivateRoute;
