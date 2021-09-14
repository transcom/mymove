import React from 'react';
import { render, screen } from '@testing-library/react';

import ShipmentSITExtensions from './ShipmentSITExtensions';
import { testProps, testPropsWithComments } from './ShipmentSITExtensionsTestParams';

describe('ShipmentSITExtensions', () => {
  it('renders the Shipment SIT Extensions', async () => {
    render(<ShipmentSITExtensions sitExtensions={testProps} />);
    expect(screen.getByText('SIT (STORAGE IN TRANSIT)')).toBeTruthy();

    expect(await screen.queryByText('Office remarks:')).toBeFalsy();
  });

  it('renders the Shipment SIT Extensions with comments', async () => {
    render(<ShipmentSITExtensions sitExtensions={testPropsWithComments} />);
    expect(screen.getByText('SIT (STORAGE IN TRANSIT)')).toBeTruthy();

    await expect(screen.getByText('Office remarks:')).toBeTruthy();
    await expect(screen.getByText('Contractor remarks:')).toBeTruthy();
  });
});
