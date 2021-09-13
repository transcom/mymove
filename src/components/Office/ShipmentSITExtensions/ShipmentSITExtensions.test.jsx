import React from 'react';
import { render, screen } from '@testing-library/react';

import ShipmentSITExtensions from './ShipmentSITExtensions';
import testProps from './ShipmentSITExtensionsTestParams';

describe('ShipmentSITExtensions', () => {
  it('renders the Shipment SIT Extensions', async () => {
    render(<ShipmentSITExtensions sitExtensions={testProps} />);
    expect(screen.getByText('SIT (STORAGE IN TRANSIT)')).toBeTruthy();
  });
});
