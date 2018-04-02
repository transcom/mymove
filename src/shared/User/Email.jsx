import React, { Component } from 'react';
import { connect } from 'react-redux';

import { bindActionCreators } from 'redux';
import { loadUserAndToken } from 'shared/User/ducks';

export class Email extends Component {
  componentDidMount() {
    this.props.loadUserAndToken();
  }
  render() {
    return (
      <span>{this.props.isLoggedIn && <span>{this.props.email}</span>}</span>
    );
  }
}

function mapStateToProps(state) {
  return {
    isLoggedIn: state.user.isLoggedIn,
    email: state.user.email,
  };
}
function mapDispatchToProps(dispatch) {
  return bindActionCreators({ loadUserAndToken }, dispatch);
}
export default connect(mapStateToProps, mapDispatchToProps)(Email);
