import React from 'react';
import { render, screen } from '@testing-library/react';
import { mount } from 'enzyme';

import OktaLogoutBanner from './index';

describe('OktaLogoutBanner component', () => {
  it('renders without crashing', () => {
    render(<OktaLogoutBanner />);
    const wrapper = mount(<OktaLogoutBanner />);

    // Check if the component renders
    const oktaLogoutBannerElement = screen.getByTestId('okta-logout-banner');
    expect(oktaLogoutBannerElement).toBeInTheDocument();

    // Check the content of the component
    expect(oktaLogoutBannerElement).toHaveTextContent('You have been logged out of Okta.');
    expect(oktaLogoutBannerElement).toHaveTextContent('Sign in');
    expect(oktaLogoutBannerElement).toHaveTextContent('troubleshooting page');

    // Check the presence of links
    const troubleshootingLink = wrapper.find('a.link[href*="okta-troubleshooting"]');
    expect(troubleshootingLink.exists()).toBe(true);

    // Dynamically generate expected href for oktaSettingsLink
    const hostname = window && window.location && window.location.hostname;
    const expectedOktaSettingsHref =
      hostname === 'office.move.mil' || hostname === 'admin.move.mil'
        ? 'https://milmove.okta.mil/enduser/settings'
        : 'https://test-milmove.okta.mil/enduser/settings';

    const oktaSettingsLink = wrapper.find(`a.link[href*="${expectedOktaSettingsHref}"]`);
    expect(oktaSettingsLink.exists()).toBe(true);
  });
});
