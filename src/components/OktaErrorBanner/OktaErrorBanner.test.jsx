import React from 'react';
import { render, screen } from '@testing-library/react';

import OktaErrorBanner from './OktaErrorBanner';

describe('OktaNeedsLoggedOutBanner component', () => {
  it('renders without crashing', () => {
    render(<OktaErrorBanner />);

    // Check if the component renders
    const oktaErrorBannerElement = screen.getByTestId('okta-error-banner');
    expect(oktaErrorBannerElement).toBeInTheDocument();

    // Check the content of the component
    expect(oktaErrorBannerElement).toHaveTextContent('You must use a different e-mail when authenticating with Okta.');
    expect(oktaErrorBannerElement).toHaveTextContent(
      'Access to this application is denied with the previously used authentication method.',
    );
  });
});
