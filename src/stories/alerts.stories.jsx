import React from 'react';
import { storiesOf } from '@storybook/react';
import { Alert } from '@trussworks/react-uswds';

storiesOf('Components|Alerts', module)
  .add('success', () => (
    <div>
      <Alert heading="Success status" type="success">
        This is a succinct, helpful success message.
      </Alert>
      <Alert slim type="success">
        This is a succinct, helpful success message.
      </Alert>
    </div>
  ))
  .add('warning', () => (
    <div>
      <Alert heading="Warning status" type="warning">
        This is a succinct, helpful warning message.
      </Alert>
      <Alert slim type="warning">
        This is a succinct, helpful warning message.
      </Alert>
    </div>
  ))
  .add('error', () => (
    <div>
      <Alert heading="Error status" type="error">
        This is a succinct, helpful error message.
      </Alert>
      <Alert slim type="error">
        This is a succinct, helpful error message.
      </Alert>
    </div>
  ))
  .add('info', () => (
    <div>
      <Alert heading="Informative status" type="info">
        This is a succinct, helpful info message.
      </Alert>
      <Alert slim type="info">
        This is a succinct, helpful info message.
      </Alert>
    </div>
  ));
