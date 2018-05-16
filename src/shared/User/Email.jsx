import React from 'react';
import { connect } from 'react-redux';

export class Email extends React.Component {
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
export default connect(mapStateToProps)(Email);
