import React from 'react';
import { Tag } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import LeftNav from './index';

import LeftNavSection from 'components/LeftNavSection/LeftNavSection';

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
      <LeftNavSection sectionName="default">Default</LeftNavSection>

      <LeftNavSection sectionName="allowances" isActive>
        Allowances
      </LeftNavSection>
      <LeftNavSection sectionName="requestedShipments">
        Requested Shipments
        <Tag className="usa-tag usa-tag--alert">
          <FontAwesomeIcon icon="exclamation" />
        </Tag>
      </LeftNavSection>
      <LeftNavSection sectionName="orders-anchor">
        Orders
        <Tag className="usa-tag--teal">INTL</Tag>
      </LeftNavSection>

      <LeftNavSection sectionName="customerInfo">
        Customer Info
        <Tag>3</Tag>
      </LeftNavSection>
    </LeftNav>
  </div>
);
