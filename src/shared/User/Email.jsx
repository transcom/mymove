import React, { Component } from 'react';
import { connect } from 'react-redux';
import { NavLink } from 'react-router-dom';

import { bindActionCreators } from 'redux';
import { loadUserAndToken } from 'shared/User/ducks';

export class Email extends Component {
  componentDidMount() {
    this.props.loadUserAndToken();
  }
  render() {
    const { isLoggedIn, email, userId } = this.props;
    return (
      <span>
        {isLoggedIn && (
          <NavLink to={`/service-member/${userId}/create`}>
            <span>{email}</span>
          </NavLink>
        )}
      </span>
    );
  }
}

function mapStateToProps(state) {
  return {
    isLoggedIn: state.user.isLoggedIn,
    email: state.user.email,
    userId: state.user.userId,
  };
}
function mapDispatchToProps(dispatch) {
  return bindActionCreators({ loadUserAndToken }, dispatch);
}
export default connect(mapStateToProps, mapDispatchToProps)(Email);
