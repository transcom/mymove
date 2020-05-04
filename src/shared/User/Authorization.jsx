import React from 'react';
import { connect } from 'react-redux';
import { compose } from 'redux';
// this is from https://hackernoon.com/role-based-authorization-in-react-c70bb7641db4
// as of 3/8 we are not using this yet, but it seems like it would work once we have a need to include roles in user state
const AuthorizationContainer = (WrappedComponent, role) =>
  function WithAuthorization(props) {
    const { userInfo } = props;
    const userRoles = userInfo.roles?.map((role) => role.roleType);

    if (!isAuthorized(userRoles, role)) {
      return (
        <div className="usa-grid">
          <h1>You are not authorized to view this page</h1>
        </div>
      );
    }

    return <WrappedComponent {...props} />;
  };

const mapStateToProps = (state) => ({
  userInfo: state.user.userInfo,
});

function isAuthorized(userRoles, authorizedRole) {
  return userRoles?.includes(authorizedRole);
}

// see https://medium.com/practo-engineering/connected-higher-order-component-hoc-93ee63c91526
const Authorization = compose(connect(mapStateToProps, null), AuthorizationContainer);
export default Authorization;
