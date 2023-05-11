import React from 'react';
import { render, screen } from '@testing-library/react';

import QAECSRMoveSearch from './QAECSRMoveSearch';

import { MockProviders } from 'testUtils';

describe('QAECSRMoveSearch page', () => {
  it('page loads', async () => {
    render(
      <MockProviders>
        <QAECSRMoveSearch />
      </MockProviders>,
    );

    const h1 = await screen.getByRole('heading', { name: 'Search for a move', level: 1 });
    expect(h1).toBeInTheDocument();

    const results = screen.queryByText(/Results/);
    expect(results).not.toBeInTheDocument();
  });
});
