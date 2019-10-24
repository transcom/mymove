import React, { Component, Fragment } from 'react';

class InfectedUpload extends Component {
  render() {
    return (
      <Fragment>
        <div className="usa-grid">
          <h2>This file could not be saved </h2>
          <p>
            We found a possible security issue. To fix that:
            <ul>
              <li>Delete this file</li>
              <li>Upload a photo of your document</li>
            </ul>
          </p>
        </div>
      </Fragment>
    );
  }
}

export default InfectedUpload;
