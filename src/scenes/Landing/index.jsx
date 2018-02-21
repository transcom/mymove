import React, { Component } from 'react';

import LoginButton from 'scenes/Landing/LoginButton';

export class Landing extends Component {
  componentDidMount() {
    document.title = 'Transcom PPP: Landing Page';
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

export default Landing;
