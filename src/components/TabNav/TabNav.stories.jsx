import React from 'react';
import { withKnobs } from '@storybook/addon-knobs';
import { Tag } from '@trussworks/react-uswds';

import TabNav from './index';

export default {
  title: 'Components/Tab Navigation',
  decorators: [withKnobs],
  parameters: {
    abstract: {
      url: 'https://share.goabstract.com/d23132ee-a6ce-451e-95f9-0a4ef0882ace?mode=design',
    },
  },
};

export const Default = () => (
  <TabNav
    items={[
      <a href="#" className="usa-current usa-nav__link" role="tab">
        <span className="tab-title">Move details</span>
      </a>,
      <a href="#" className="usa-nav__link" role="tab">
        <span className="tab-title">Move task order</span>
      </a>,
      <a href="#" className="usa-nav__link" role="tab">
        <span className="tab-title">Payment requests</span>
      </a>,
    ]}
    role="navigation"
  />
);

export const withTag = () => (
  <TabNav
    items={[
      <a href="#" className="usa-nav__link" role="tab">
        <span className="tab-title">Move details</span>
        <Tag>2</Tag>
      </a>,
      <a href="#" className="usa-current usa-nav__link" role="tab">
        <span className="tab-title">Move task order</span>
      </a>,
      <a href="#" className="usa-nav__link" role="tab">
        <span className="tab-title">Payment requests</span>
      </a>,
    ]}
    role="navigation"
  />
);
