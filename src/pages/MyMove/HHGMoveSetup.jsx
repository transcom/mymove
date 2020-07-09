import React, { Component } from 'react';

// import { Form } from 'components/form/Form';
import { HHGDetailsForm } from 'components/Customer/HHGDetailsForm';

// eslint-disable-next-line react/prefer-stateless-function
export class HHGMoveSetup extends Component {
  render() {
    return (
      <div>
        <h1>Now lets arrange details for the professional movers</h1>
        <HHGDetailsForm />
      </div>
    );
  }
}

export default HHGMoveSetup;
