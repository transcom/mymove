import React from 'react';
import { Tag } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import LeftNav from './index';

// Left Nav
export default {
  title: 'Components/Left Nav',
  component: LeftNav,
  parameters: {
    abstract: {
      url: 'https://share.goabstract.com/6e8668b7-5562-4894-a661-648ab4883d8f?mode=design',
    },
  },
};

export const Basic = () => (
  <div id="l-nav" style={{ padding: '20px', background: '#f0f0f0' }}>
    <LeftNav>
      <a href="#">Default</a>

      <a href="#" className="active">
        Allowances
      </a>
      <a href="#">
        Requested Shipments
        <Tag className="usa-tag usa-tag--alert">
          <FontAwesomeIcon icon="exclamation" />
        </Tag>
      </a>
      <a href="#orders-anchor">
        Orders
        <Tag className="usa-tag--teal">INTL</Tag>
      </a>

      <a href="#">
        Customer Info
        <Tag>3</Tag>
      </a>
    </LeftNav>
  </div>
);
