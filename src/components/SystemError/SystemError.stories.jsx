import React from 'react';

import SystemError from './index';

export default {
  title: 'Components/Alerts/System Error',
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
