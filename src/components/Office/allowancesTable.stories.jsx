import React from 'react';
import { storiesOf } from '@storybook/react';
import { withKnobs, object } from '@storybook/addon-knobs';

import AllowancesTable from './AllowancesTable';

const info = {
  branch: 'NAVY',
  rank: 'E_6',
  weightAllowance: 11000,
  authorizedWeight: 11000,
  progear: 2000,
  spouseProgear: 500,
  storageInTransit: 90,
  dependents: true,
};

storiesOf('TOO/TIO Components|AllowancesTable', module)
  .addDecorator(withKnobs)
  .add('Allowances Table', () => <AllowancesTable info={object('info', info)} />);
