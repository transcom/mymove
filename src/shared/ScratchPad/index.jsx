import React, { Component } from 'react';
import BasicPanel from 'shared/BasicPanel';

class ScratchPad extends Component {
  render() {
    return (
      <div className="usa-grid">
        <div className="usa-width-one-whole">
          <BasicPanel title={'TEST TITLE'}>HERE ARE MY BABIES</BasicPanel>
        </div>
      </div>
    );
  }
}

export default ScratchPad;
