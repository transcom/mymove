import React from 'react';
import { storiesOf } from '@storybook/react';
import { Alert } from '@trussworks/react-uswds';

storiesOf('Components|Alerts', module)
  .add('success', () => (
    <div>
      <Alert heading="Success status" type="success">
        <React.Fragment key=".0">This is a succinct, helpful success message.</React.Fragment>
      </Alert>
      <Alert type="success">
        <React.Fragment key=".0">This is a succinct, helpful success message.</React.Fragment>
      </Alert>
    </div>
  ))
  .add('warning', () => (
    <div>
      <Alert heading="Warning status" type="warning">
        <React.Fragment key=".0">This is a succinct, helpful warning message.</React.Fragment>
      </Alert>
      <Alert type="warning">
        <React.Fragment key=".0">This is a succinct, helpful warning message.</React.Fragment>
      </Alert>
    </div>
  ))
  .add('error', () => (
    <div>
      <Alert heading="Error status" type="error">
        <React.Fragment key=".0">This is a succinct, helpful error message.</React.Fragment>
      </Alert>
      <Alert type="error">
        <React.Fragment key=".0">This is a succinct, helpful error message.</React.Fragment>
      </Alert>
    </div>
  ))
  .add('info', () => (
    <div>
      <Alert heading="Informative status" type="info">
        <React.Fragment key=".0">This is a succinct, helpful info message.</React.Fragment>
      </Alert>
      <Alert type="info">
        <React.Fragment key=".0">This is a succinct, helpful info message.</React.Fragment>
      </Alert>
    </div>
  ));
