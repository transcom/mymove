import React from 'react';
import { render, screen } from '@testing-library/react';

import OktaLogoutBanner from './index';

describe('OktaLogoutBanner component', () => {
  it('renders without crashing', () => {
    render(<OktaLogoutBanner />);
    expect(screen.getByText(/You have been logged out of Okta/)).toBeInTheDocument();
  });
});
