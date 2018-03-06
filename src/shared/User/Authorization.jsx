import React from 'react';
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
