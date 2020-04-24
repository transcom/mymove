import React from 'react';

import { storiesOf } from '@storybook/react';
import { withKnobs, text } from '@storybook/addon-knobs';

import SomethingWentWrong from '../shared/SomethingWentWrong';

// Left Nav

storiesOf('SomethingWentWrong', module)
  .addDecorator(withKnobs)
  .add('component', () => (
    <div id="l-nav" style={{ padding: '20px', background: '#f0f0f0' }}>
      <SomethingWentWrong
        error={text('SomethingWentWrong.error', 'error')}
        info={text('SomethingWentWrong.info', 'info')}
      />
    </div>
  ));
