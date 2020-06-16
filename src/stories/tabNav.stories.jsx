import React from 'react';

import { storiesOf } from '@storybook/react';
import { withKnobs } from '@storybook/addon-knobs';

import { Tag } from '@trussworks/react-uswds';

import TabNav from '../components/TabNav';

storiesOf('Components|TabNav', module)
  .addDecorator(withKnobs)
  .add('default', () => (
    <TabNav
      items={[
        <a href="#" className="usa-current usa-nav__link">
          <span className="tab-title">Move Details</span>
        </a>,
        <a href="#" className="usa-nav__link">
          <span className="tab-title">Move Task Order</span>
        </a>,
        <a href="#" className="usa-nav__link">
          <span className="tab-title">Payment requests</span>
        </a>,
      ]}
    />
  ))
  .add('withTag', () => (
    <TabNav
      items={[
        <>
          <a href="#" className="usa-nav__link">
            <span className="tab-title">Move Details</span>
            <Tag>2</Tag>
          </a>
        </>,
        <a href="#" className="usa-current usa-nav__link">
          <span className="tab-title">Move Task Order</span>
        </a>,
        <a href="#" className="usa-nav__link">
          <span className="tab-title">Payment requests</span>
        </a>,
      ]}
    />
  ));
