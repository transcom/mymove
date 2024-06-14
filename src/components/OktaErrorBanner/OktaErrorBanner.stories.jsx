import React from 'react';

import OktaErrorBanner from './OktaErrorBanner';

import { MockRouterProvider } from 'testUtils';

export default {
  title: 'Components / Okta Error Banner',
};

export const OktaErrorBannerComponent = () => (
  <MockRouterProvider>
    <OktaErrorBanner />
  </MockRouterProvider>
);
