import React from 'react';
import { storiesOf } from '@storybook/react';
import { withKnobs, text, object } from '@storybook/addon-knobs';
import CustomerRemarksAgentsDetails from 'components/CustomerRemarksAgentsDetails';

storiesOf('TOO/TIO Components|CustomerRemarksAgentsDetails', module)
  .addDecorator(withKnobs)
  .add('with no details', () => {
    return (
      <div id="containers" style={{ padding: '20px' }}>
        <CustomerRemarksAgentsDetails />
      </div>
    );
  })
  .add('with details', () => {
    const customerRemarks = 'This is a customer remark.';
    const releasingAgent = {
      firstName: 'FirstName',
      lastName: 'LastName',
      phone: '(111) 111-1111',
      email: 'test@test.com',
    };
    const receivingAgent = {
      firstName: 'FirstName',
      lastName: 'LastName',
      phone: '(111) 111-1111',
      email: 'test@test.com',
    };

    return (
      <div id="containers" style={{ padding: '20px' }}>
        <CustomerRemarksAgentsDetails
          customerRemarks={text('customerRemarks', customerRemarks)}
          releasingAgent={object('releasingAgent', releasingAgent)}
          receivingAgent={object('receivingAgent', receivingAgent)}
        />
      </div>
    );
  });
