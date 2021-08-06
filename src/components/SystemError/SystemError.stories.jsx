import React from 'react';

import SystemError from './index';

export default {
  title: 'Components/Alerts/System Error',
  parameters: {
    abstract: {
      url: 'https://share.goabstract.com/d9ad20e6-944c-48a2-bbd2-1c7ed8bc1315?mode=design',
    },
  },
};

export const SystemErrorComponent = () => (
  <SystemError>
    This is a succinct, helpful error message. Also inlcuded is an example of some&nbsp;
    <a href="#">link text</a>
    .
    <br />
    This is a second line to test the line height.
  </SystemError>
);
