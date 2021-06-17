import React from 'react';
import { render } from '@testing-library/react';

import ImportCsvButton from './index';

import { MockProviders } from 'testUtils';

describe('ImportCsvButton component', () => {
  it('renders without error', async () => {
    const { debug } = render(
      <MockProviders>
        <ImportCsvButton resource="office_users" />
      </MockProviders>,
    );
    debug();
  });
});
