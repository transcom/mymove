import React, { Component } from 'react';

import { HHGDetailsForm } from 'components/Customer/HHGDetailsForm';

// eslint-disable-next-line react/prefer-stateless-function
export class HHGMoveSetup extends Component {
  render() {
    return (
      <div>
        <h3>Now lets arrange details for the professional movers</h3>
        <HHGDetailsForm />
      </div>
    );
  }
}

export default HHGMoveSetup;
