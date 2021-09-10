import React from 'react';
import { render, screen } from '@testing-library/react';

import MarkerIO from './MarkerIO';

describe('components/ThirdParty/MarkerIO', () => {
  it('renders a script tag within the page', async () => {
    render(<MarkerIO />);

    expect(await screen.findByTestId('markerio script tag')).toBeInTheDocument();
  });
});
