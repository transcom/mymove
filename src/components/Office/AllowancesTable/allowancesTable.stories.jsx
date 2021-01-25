import React from 'react';
import { object, boolean } from '@storybook/addon-knobs';

import AllowancesTable from './AllowancesTable';

import { MockProviders } from 'testUtils';

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

export default {
  title: 'Office Components/AllowancesTable',
  decorators: [
    (Story) => (
      <div style={{ 'max-width': '800px' }}>
        <MockProviders initialEntries={[`/moves/1000/details`]}>
          <Story />
        </MockProviders>
      </div>
    ),
  ],
};

export const Default = () => <AllowancesTable info={object('info', info)} />;

export const HasEditBtn = () => (
  <AllowancesTable info={object('info', info)} showEditBtn={boolean('Show Edit Btn', true)} />
);
