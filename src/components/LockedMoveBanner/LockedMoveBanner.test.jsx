import React from 'react';
import { render } from '@testing-library/react';

import LockedMoveBanner from './LockedMoveBanner';

describe('LockedMoveBanner', () => {
  it('renders children with a lock icon', () => {
    const { getByTestId, getByText } = render(<LockedMoveBanner>Some random text</LockedMoveBanner>);

    const banner = getByTestId('locked-move-banner');
    expect(banner).toBeInTheDocument();

    const childText = getByText('Some random text');
    expect(childText).toBeInTheDocument();
  });
});
