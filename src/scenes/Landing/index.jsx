import React, { Component } from 'react';
import { parse } from 'qs';

import Alert from 'shared/Alert';
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
        <h1>Welcome!</h1>
        {query.err && (
          <Alert type="error" heading="Server Error">
            Sorry, something went wrong
          </Alert>
        )}
        <LoginButton />
      </div>
    );
  }
}

export default Landing;
