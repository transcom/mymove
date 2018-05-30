import React from 'react';
import { connect } from 'react-redux';

export class Email extends React.Component {
  render() {
    if (!this.props.isLoggedIn) return <span />;
    else return <span>{this.props.email}</span>;
  }
}

function mapStateToProps(state) {
  return {
    isLoggedIn: state.user.isLoggedIn,
    email: state.user.email,
  };
}
export default connect(mapStateToProps)(Email);
