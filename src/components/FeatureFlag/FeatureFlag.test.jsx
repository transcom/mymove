import React from 'react';
import { screen, render, waitFor } from '@testing-library/react';

import FeatureFlag, { featureIsEnabled, DISABLED_VALUE, ENABLED_VALUE } from './FeatureFlag';

import { getFeatureFlagForUser } from 'services/internalApi';

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  getFeatureFlagForUser: jest.fn(),
}));

describe('FeatureFlag', () => {
  const featureFlagRender = (flagValue) => {
    if (featureIsEnabled(flagValue)) {
      return <div>Yes</div>;
    }
    return <div>Nope</div>;
  };

  it('should render enabled if enabled', async () => {
    getFeatureFlagForUser.mockResolvedValue({ match: true, value: ENABLED_VALUE });

    render(<FeatureFlag flagKey="key" render={featureFlagRender} />);
    await waitFor(() => {
      expect(screen.getByText('Yes')).toBeInTheDocument();
    });
  });

  it('should render disabled if disabled', async () => {
    getFeatureFlagForUser.mockResolvedValue({ match: true, value: DISABLED_VALUE });

    render(<FeatureFlag flagKey="key" render={featureFlagRender} />);
    await waitFor(() => {
      expect(screen.getByText('Nope')).toBeInTheDocument();
    });
  });
  it('should render disabled if no match', async () => {
    getFeatureFlagForUser.mockResolvedValue({ match: false, value: '' });

    render(<FeatureFlag flagKey="key" render={featureFlagRender} />);
    await waitFor(() => {
      expect(screen.getByText('Nope')).toBeInTheDocument();
    });
  });
});
