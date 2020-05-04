import React from 'react';

import { storiesOf } from '@storybook/react';
import { withKnobs, text } from '@storybook/addon-knobs';

import SomethingWentWrong from '../shared/SomethingWentWrong';

// Left Nav

storiesOf('SomethingWentWrong', module)
  .addDecorator(withKnobs)
  .add('component', () => (
    <div className="usa-grid">
      <div style={{ textAlign: 'center' }}>
        <SomethingWentWrong
          error={text('SomethingWentWrong.error', 'error')}
          info={text('SomethingWentWrong.info', 'info')}
        />
      </div>
    </div>
  ));
