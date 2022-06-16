import React from 'react';
import { Alert } from '@trussworks/react-uswds';

export default {
  title: 'Components/Alerts',
  parameters: {
    abstract: {
      url: 'https://share.goabstract.com/1be53210-d556-4f7d-9141-db2944f2e058?mode=design',
    },
  },
};

export const success = () => (
  <div>
    <Alert headingLevel="h4" heading="Success status" type="success">
      This is a succinct, helpful success message. This is a <a href="#">link</a>.
    </Alert>
    <Alert slim type="success">
      This is a succinct, helpful success message. This is a <a href="#">link</a>.
    </Alert>
  </div>
);

export const warning = () => (
  <div>
    <Alert headingLevel="h4" heading="Warning status" type="warning">
      This is a succinct, helpful warning message. This is a <a href="#">link</a>.
    </Alert>
    <Alert headingLevel="h4" slim type="warning">
      This is a succinct, helpful warning message. This is a <a href="#">link</a>.
    </Alert>
  </div>
);

export const error = () => (
  <div>
    <Alert headingLevel="h4" heading="Error status" type="error">
      This is a succinct, helpful error message. This is a <a href="#">link</a>.
    </Alert>
    <Alert headingLevel="h4" slim type="error">
      This is a succinct, helpful error message. This is a <a href="#">link</a>.
    </Alert>
  </div>
);

export const info = () => (
  <div>
    <Alert headingLevel="h4" heading="Informative status" type="info">
      This is a succinct, helpful info message. This is a <a href="#">link</a>.
    </Alert>
    <Alert headingLevel="h4" slim type="info">
      This is a succinct, helpful info message. This is a <a href="#">link</a>.
    </Alert>
  </div>
);
