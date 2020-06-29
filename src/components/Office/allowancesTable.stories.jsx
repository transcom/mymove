import React from 'react';
import { storiesOf } from '@storybook/react';
import { withKnobs, object } from '@storybook/addon-knobs';

import AllowancesTable from './AllowancesTable';

const info = {
  branch: 'NAVY',
  rank: 'E_6',
  weightAllowance: '11,000 lbs',
  authorizedWeight: '11,000 lbs',
  progear: '2,000 lbs',
  spouseProgear: '500 lbs',
  storageInTransit: '90 days',
  dependents: 'Authorized',
};

storiesOf('TOO/TIO Components|AllowancesTable', module)
  .addDecorator(withKnobs)
  .add('Allowances Table', () => <AllowancesTable info={object('info', info)} />);
