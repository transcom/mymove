import React from 'react';
import { connect } from 'react-redux';
// this is from https://hackernoon.com/role-based-authorization-in-react-c70bb7641db4
// as of 3/8 we are not using this yet, but it seems like it would work once we have a need to include roles in user state
const AuthorizationContainer = (WrappedComponent, allowedRoles) =>
  function WithAuthorization(props) {
    const { role } = props;
    if (allowedRoles.includes(role)) {
      return <WrappedComponent {...props} />;
    } else {
      return (
        <div className="usa-grid">
          <h1>You are not authorized to view this page</h1>
        </div>
      );
    }
  };

const Authorization = connect(state => ({
  role: state.user.role,
}))(AuthorizationContainer);
export default Authorization;
