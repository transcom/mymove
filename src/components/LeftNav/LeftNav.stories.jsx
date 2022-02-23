import React from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import LeftNav from './LeftNav';

import LeftNavTag from 'components/LeftNavTag/LeftNavTag';

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
    <LeftNav sections={['allowances', 'requested-shipments', 'orders', 'customer-info']}>
      <LeftNavTag associatedSectionName="requested-shipments" showTag className="usa-tag usa-tag--alert">
        <FontAwesomeIcon icon="exclamation" />
      </LeftNavTag>

      <LeftNavTag associatedSectionName="orders" showTag className="usa-tag--teal">
        INTL
      </LeftNavTag>

      <LeftNavTag associatedSectionName="customer-info" showTag>
        3
      </LeftNavTag>
    </LeftNav>
  </div>
);
