import React from 'react';

import { storiesOf } from '@storybook/react';
import { withKnobs, boolean, text } from '@storybook/addon-knobs';

import TabNav from '../components/TabNav';

storiesOf('components', module)
  .addDecorator(withKnobs)
  .add('TabNav', () => (
    <TabNav
      options={[
        {
          title: text('Option1.title', 'Option 1', 'First Tab'),
          active: boolean('Option1.active', true, 'First Tab'),
          notice: text('Option1.notice', '2', 'First Tab'),
        },
        {
          title: text('Option2.title', 'Option 2', 'Second Tab'),
          active: boolean('Option2.active', false, 'Second Tab'),
          notice: text('Option2.notice', null, 'Second Tab'),
        },
        {
          title: text('Option3.title', 'Option 3', 'Third Tab'),
          active: boolean('Option3.active', false, 'Third Tab'),
          notice: text('Option3.notice', null, 'Third Tab'),
        },
      ]}
    />
  ));
