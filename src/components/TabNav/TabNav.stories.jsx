import React from 'react';
import { Tag } from '@trussworks/react-uswds';

import TabNav from './index';

export default {
  title: 'Components/Tab Navigation',
};

export const Default = () => (
  <TabNav
    items={[
      <a href="#" className="usa-current usa-nav__link">
        <span className="tab-title">Move Details</span>
      </a>,
      <a href="#" className="usa-nav__link">
        <span className="tab-title">Move Task Order</span>
      </a>,
      <a href="#" className="usa-nav__link">
        <span className="tab-title">Payment Requests</span>
      </a>,
    ]}
    role="navigation"
  />
);

export const withTag = () => (
  <TabNav
    items={[
      <a href="#" className="usa-nav__link">
        <span className="tab-title">Move Details</span>
        <Tag>2</Tag>
      </a>,
      <a href="#" className="usa-current usa-nav__link">
        <span className="tab-title">Move Task Order</span>
      </a>,
      <a href="#" className="usa-nav__link">
        <span className="tab-title">Payment Requests</span>
      </a>,
    ]}
    role="navigation"
  />
);
