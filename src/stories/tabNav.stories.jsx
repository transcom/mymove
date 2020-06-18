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
  ))
  .add('withTag', () => (
    <TabNav
      items={[
        <>
          <a href="#" className="usa-nav__link" role="tab">
            <span className="tab-title">Move details</span>
            <Tag>2</Tag>
          </a>
        </>,
        <a href="#" className="usa-current usa-nav__link" role="tab">
          <span className="tab-title">Move task order</span>
        </a>,
        <a href="#" className="usa-nav__link" role="tab">
          <span className="tab-title">Payment requests</span>
        </a>,
      ]}
      role="navigation"
    />
  ));
