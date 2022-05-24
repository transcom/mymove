import React from 'react';
import { screen } from '@testing-library/react';

import QAECSRMoveSearch from './QAECSRMoveSearch';

import { renderWithRouter } from 'testUtils';

describe('QAECSRMoveSearch page', () => {
  it('page loads', async () => {
    renderWithRouter(<QAECSRMoveSearch />);

    const h1 = await screen.getByRole('heading', { name: 'Search for a move', level: 1 });
    expect(h1).toBeInTheDocument();
  });
});
