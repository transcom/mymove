import React from 'react';
import { screen, render, waitFor } from '@testing-library/react';

import FeatureFlag, { BOOLEAN_FLAG_TYPE, VARIANT_FLAG_TYPE } from './FeatureFlag';

import { getBooleanFeatureFlagForUser, getVariantFeatureFlagForUser } from 'services/internalApi';

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  getBooleanFeatureFlagForUser: jest.fn(),
  getVariantFeatureFlagForUser: jest.fn(),
}));

describe('FeatureFlag', () => {
  const featureFlagRender = (flagValue) => {
    if (flagValue === 'true') {
      return <div>Yes</div>;
    }
    if (flagValue === 'false') {
      return <div>Nope</div>;
    }
    if (flagValue === '') {
      return <div>Missing</div>;
    }
    return <div>{flagValue}</div>;
  };

  it('should render yes if boolean enabled', async () => {
    getBooleanFeatureFlagForUser.mockResolvedValue({ match: true });

    render(<FeatureFlag flagType={BOOLEAN_FLAG_TYPE} flagKey="key" render={featureFlagRender} />);
    await waitFor(() => {
      expect(screen.getByText('Yes')).toBeInTheDocument();
    });
  });

  it('should render nope if boolean disabled', async () => {
    getBooleanFeatureFlagForUser.mockResolvedValue({ match: false });

    render(<FeatureFlag flagType={BOOLEAN_FLAG_TYPE} flagKey="key" render={featureFlagRender} />);
    await waitFor(() => {
      expect(screen.getByText('Nope')).toBeInTheDocument();
    });
  });

  it('should render missing if variant has no match', async () => {
    getVariantFeatureFlagForUser.mockResolvedValue({ match: false, variant: '' });

    render(<FeatureFlag flagType={VARIANT_FLAG_TYPE} flagKey="key" render={featureFlagRender} />);
    await waitFor(() => {
      expect(screen.getByText('Missing')).toBeInTheDocument();
    });
  });

  it('should render value if variant has match', async () => {
    const myVariant = 'my_variant';
    getVariantFeatureFlagForUser.mockResolvedValue({ match: true, variant: myVariant });

    render(<FeatureFlag flagType={VARIANT_FLAG_TYPE} flagKey="key" render={featureFlagRender} />);
    await waitFor(() => {
      expect(screen.getByText(myVariant)).toBeInTheDocument();
    });
  });
});
