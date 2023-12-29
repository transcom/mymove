import React from 'react';
import { render, screen } from '@testing-library/react';
import { mount } from 'enzyme';

import { OktaLoggedOutBanner, OktaNeedsLoggedOutBanner } from './index';

describe('OktaLoggedOutBanner component', () => {
  it('renders without crashing', () => {
    render(<OktaLoggedOutBanner />);
    const wrapper = mount(<OktaLoggedOutBanner />);

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

describe('OktaNeedsLoggedOutBanner component', () => {
  it('renders without crashing', () => {
    render(<OktaNeedsLoggedOutBanner />);
    const wrapper = mount(<OktaNeedsLoggedOutBanner />);

    // Check if the component renders
    const oktaLogoutBannerElement = screen.getByTestId('okta-logout-banner');
    expect(oktaLogoutBannerElement).toBeInTheDocument();

    // Check the content of the component
    expect(oktaLogoutBannerElement).toHaveTextContent(
      'You have an existing Okta session. Please log out of Okta completely.',
    );
    expect(oktaLogoutBannerElement).toHaveTextContent(
      "You can access your Okta dashboard by following this link. In the top-right corner, you can click the drop down where it displays your name and select 'Sign Out'. Once you sign out of Okta, you should be able to sign into MilMove.If you have issues logging in or authenticating with Okta, please refer to our troubleshooting page.",
    );

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
