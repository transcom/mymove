import React, { Component } from 'react';
import { parse } from 'qs';

import LoginButton from 'scenes/Landing/LoginButton';

export class Landing extends Component {
  componentDidMount() {
    document.title = 'Transcom PPP: Landing Page';
  }

  render() {
    const location = this.props.location;
    const query = parse(location.search.substr(1));

    return (
      <div className="usa-grid">
        {query.email && <h1>Welcome {query.email}!</h1>}
        {!query.email && <h1>Welcome!</h1>}
        <LoginButton />
      </div>
    );
  }
}

export default Landing;
