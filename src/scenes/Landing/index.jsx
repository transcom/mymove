import React, { Component } from 'react';
import { connect } from 'react-redux';

import LoginButton from 'shared/User/LoginButton';
import { bindActionCreators } from 'redux';
import { loadUserAndToken } from 'shared/User/ducks';

export class Landing extends Component {
  componentDidMount() {
    document.title = 'Transcom PPP: Landing Page';
    this.props.loadUserAndToken();
  }

  render() {
    return (
      <div className="usa-grid">
        <h1>Welcome!</h1>
        <LoginButton />
      </div>
    );
  }
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ loadUserAndToken }, dispatch);
}

export default connect(null, mapDispatchToProps)(Landing);
