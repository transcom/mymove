import React from 'react';
import { storiesOf } from '@storybook/react';
import { withKnobs, text } from '@storybook/addon-knobs';

import TabNav from '../components/TabNav';
import TabNavPanel from '../components/TabNavPanel';

storiesOf('Components|TabNav', module)
  .addDecorator(withKnobs)
  .add('default', () => (
    <TabNav
      options={[
        {
          title: text('Option1.title', 'Option 1', 'First Tab'),
          notice: text('Option1.notice', null, 'First Tab'),
        },
        {
          title: text('Option2.title', 'Option 2', 'Second Tab'),
          notice: text('Option2.notice', null, 'Second Tab'),
        },
        {
          title: text('Option3.title', 'Option 3', 'Third Tab'),
          notice: text('Option3.notice', null, 'Third Tab'),
        },
      ]}
    >
      <TabNavPanel>
        Body Of Tab
        {text('Option1.title', 'Option 1', 'First Tab')}
      </TabNavPanel>
      <TabNavPanel>
        Body Of Tab
        {text('Option2.title', 'Option 2', 'Second Tab')}
      </TabNavPanel>
      <TabNavPanel>
        Body Of Tab
        {text('Option3.title', 'Option 3', 'Third Tab')}
      </TabNavPanel>
    </TabNav>
  ))
  .add('withTag', () => (
    <TabNav
      options={[
        {
          title: text('Option1.title', 'Option 1', 'First Tab'),
          notice: text('Option1.notice', '2', 'First Tab'),
        },
        {
          title: text('Option2.title', 'Option 2', 'Second Tab'),
          notice: text('Option2.notice', null, 'Second Tab'),
        },
        {
          title: text('Option3.title', 'Option 3', 'Third Tab'),
          notice: text('Option3.notice', null, 'Third Tab'),
        },
      ]}
    >
      <TabNavPanel>
        Body Of Tab
        {text('Option1.title', 'Option 1', 'First Tab')}
      </TabNavPanel>
      <TabNavPanel>
        Body Of Tab
        {text('Option2.title', 'Option 2', 'Second Tab')}
      </TabNavPanel>
      <TabNavPanel>
        Body Of Tab
        {text('Option3.title', 'Option 3', 'Third Tab')}
      </TabNavPanel>
    </TabNav>
  ));
