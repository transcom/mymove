import React from 'react';
import { render, screen } from '@testing-library/react';

import { futureSITShipment, noSITShipment, SITShipment } from '../ShipmentSITDisplay/ShipmentSITDisplayTestParams';

import ShipmentDetailsMain from './ShipmentDetailsMain';

import { MockProviders } from 'testUtils';

const shipmentDetailsMainParams = {
  handleDivertShipment: () => {},
  handleRequestReweighModal: () => {},
  handleReviewSITExtension: () => {},
  handleSubmitSITExtension: () => {},
  dutyLocationAddresses: {
    originDutyLocationAddress: {
      address: null,
    },
    destinationDutyLocationAddress: {
      address: null,
    },
  },
};

describe('Shipment Details Main', () => {
  it('displays SIT when there is a sitStatus', () => {
    render(
      <MockProviders>
        <ShipmentDetailsMain {...shipmentDetailsMainParams} shipment={SITShipment} />
      </MockProviders>,
    );

    expect(screen.getByText('SIT (STORAGE IN TRANSIT)')).toBeInTheDocument();
  });
  it('displays SIT when the shipment has sitDaysAllowance', () => {
    render(
      <MockProviders>
        <ShipmentDetailsMain {...shipmentDetailsMainParams} shipment={futureSITShipment} />
      </MockProviders>,
    );

    expect(screen.getByText('SIT (STORAGE IN TRANSIT)')).toBeInTheDocument();
  });
  it('does not display SIT when there is no sitStatus and no sit days allowance', () => {
    render(
      <MockProviders>
        <ShipmentDetailsMain {...shipmentDetailsMainParams} shipment={noSITShipment} />
      </MockProviders>,
    );

    expect(screen.queryByText('SIT (STORAGE IN TRANSIT)')).not.toBeInTheDocument();
  });
});
