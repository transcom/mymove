import React from 'react';
import { render } from '@testing-library/react';

import ShipmentModificationTag from './ShipmentModificationTag';

import { shipmentModificationTypes } from 'constants/shipments';

describe('ShipmentModificationTag Component', () => {
  it('renders the canceled tag', async () => {
    const { getByText } = render(
      <ShipmentModificationTag shipmentModificationType={shipmentModificationTypes.CANCELED} />,
    );
    expect(getByText('CANCELED')).toBeInTheDocument();
  });
  it('renders the diversion tag', async () => {
    const { getByText } = render(
      <ShipmentModificationTag shipmentModificationType={shipmentModificationTypes.DIVERSION} />,
    );
    expect(getByText('DIVERSION')).toBeInTheDocument();
  });
});
