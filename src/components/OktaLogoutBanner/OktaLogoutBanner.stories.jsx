import React from 'react';

import { OktaLoggedOutBanner, OktaNeedsLoggedOutBanner } from './index';

import { MockRouterProvider } from 'testUtils';

export default {
  title: 'Components / Okta Logout Banners',
};

export const OktaLoggedOutBannerComponent = () => (
  <MockRouterProvider>
    <OktaLoggedOutBanner />
  </MockRouterProvider>
);

export const OktaNeedsLoggedOutBannerComponent = () => (
  <MockRouterProvider>
    <OktaNeedsLoggedOutBanner />
  </MockRouterProvider>
);
