import React, { Component } from 'react';

class ProcessingUpload extends Component {
  render() {
    return (
      <div className="usa-grid">
        <div className="grid-container">
          <h2>Your file is being scanned for viruses </h2>
          <p>It will be available within a few minutes.</p>
        </div>
      </div>
    );
  }
}

export default ProcessingUpload;
