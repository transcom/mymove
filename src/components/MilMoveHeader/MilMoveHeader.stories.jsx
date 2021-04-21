import React from 'react';
import { action } from '@storybook/addon-actions';

import OfficeUserInfo from './OfficeUserInfo';
import CustomerUserInfo from './CustomerUserInfo';
import LoggedOutUserInfo from './LoggedOutUserInfo';

import MilMoveHeader from './index';

import { MockProviders } from 'testUtils';

export default {
  title: 'Components/Headers/MilMove Header',
  parameters: {
    abstract: {
      url: 'https://share.goabstract.com/d9ad20e6-944c-48a2-bbd2-1c7ed8bc1315?mode=design',
    },
  },
  decorators: [
    (Story) => (
      <MockProviders>
        <Story />
      </MockProviders>
    ),
  ],
};

const props = {
  lastName: 'Baker',
  firstName: 'Riley',
  handleLogout: action('clicked'),
};

export const LoggedOutHeader = () => (
  <MilMoveHeader>
    <LoggedOutUserInfo handleLogin={action('clicked')} />
  </MilMoveHeader>
);

export const LoggedInOfficeHeader = () => (
  <MilMoveHeader>
    <OfficeUserInfo {...props} />
  </MilMoveHeader>
);

export const LoggedInOfficeHeaderWithNavigation = () => (
  <MilMoveHeader>
    <ul className="usa-nav__primary">
      <li className="usa-nav__primary-item">
        <a href="#">Navigation Link</a>
      </li>
      <li className="usa-nav__primary-item">
        <a href="#">Navigation Link</a>
      </li>
      <li className="usa-nav__primary-item">
        <a href="#">Navigation Link</a>
      </li>
    </ul>
    <OfficeUserInfo {...props} />
  </MilMoveHeader>
);

export const LoggedInCustomerHeader = () => (
  <MilMoveHeader>
    <CustomerUserInfo {...props} />
  </MilMoveHeader>
);

export const LoggedInCustomerHeaderWithProfileLink = () => (
  <MilMoveHeader>
    <CustomerUserInfo {...props} showProfileLink />
  </MilMoveHeader>
);
