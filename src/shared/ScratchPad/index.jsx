import React, { Component } from 'react';

class ScratchPad extends Component {
  render() {
    return (
      <div className="usa-grid grid-wide panels-body">
        <div className="usa-width-one-whole">
          <div className="usa-width-one-third">
            <button className="usa-button">Click Me (I do nothing)</button>
          </div>
        </div>
      </div>
    );
  }
}
export default ScratchPad;
