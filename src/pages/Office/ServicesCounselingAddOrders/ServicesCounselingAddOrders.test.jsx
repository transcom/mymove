/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';

import ServicesCounselingAddOrders from './ServicesCounselingAddOrders';

import { MockProviders } from 'testUtils';

const renderWithMocks = () => {
  render(
    <MockProviders>
      <ServicesCounselingAddOrders />
    </MockProviders>,
  );
};

describe('ServicesCounselingAddOrders component', () => {
  it('renders the Services Counseling Add Orders Form', async () => {
    renderWithMocks();

    const h1 = await screen.getByRole('heading', { name: 'Tell us about the orders', level: 1 });
    await waitFor(() => {
      expect(h1).toBeInTheDocument();
    });
  });
});
