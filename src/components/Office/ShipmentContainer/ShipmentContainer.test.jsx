import React from 'react';
import { render, screen } from '@testing-library/react';

import ShipmentHeading from '../ShipmentHeading/ShipmentHeading';

import ShipmentContainer from './ShipmentContainer';

import { shipmentStatuses } from 'constants/shipments';
import { SHIPMENT_OPTIONS } from 'shared/constants';

const headingInfo = {
  shipmentInfo: {
    shipmentID: '1',
    shipmentStatus: shipmentStatuses.SUBMITTED,
    shipmentType: 'Household Goods',
    originCity: 'San Antonio',
    originState: 'TX',
    originPostalCode: '98421',
    destinationAddress: {
      city: 'Tacoma',
      state: 'WA',
      postalCode: '98421',
    },
    scheduledPickupDate: '27 Mar 2020',
    reweigh: { id: '00000000-0000-0000-0000-000000000000' },
    ifMatchEtag: 'etag',
    moveTaskOrderID: 'mtoID',
  },
  handleShowCancellationModal: jest.fn(),
};

describe('Shipment Container', () => {
  it('renders the container successfully', async () => {
    render(
      <ShipmentContainer>
        <ShipmentHeading {...headingInfo} />
      </ShipmentContainer>,
    );

    const shipmentContainer = await screen.findByTestId('ShipmentContainer');

    expect(shipmentContainer).toBeInTheDocument();

    expect(shipmentContainer.className).toContain('container--accent--default');
  });

  it('renders a child component passed to it', async () => {
    render(
      <ShipmentContainer>
        <ShipmentHeading {...headingInfo} />
      </ShipmentContainer>,
    );

    const childHeading = await screen.findByRole('heading', { level: 2, name: headingInfo.shipmentInfo.shipmentType });
    expect(childHeading).toBeInTheDocument();
  });

  it.each([
    [SHIPMENT_OPTIONS.HHG, 'container--accent--hhg'],
    [SHIPMENT_OPTIONS.NTS, 'container--accent--nts'],
    [SHIPMENT_OPTIONS.NTSR, 'container--accent--ntsr'],
    [SHIPMENT_OPTIONS.UNACCOMPANIED_BAGGAGE, 'container--accent--ub'],
  ])('renders a container for a shipment (%s) with className %s ', async (shipmentType, expectedClass) => {
    const newHeadingInfo = {
      ...headingInfo,
      shipmentInfo: { ...headingInfo.shipmentInfo, shipmentType },
    };

    render(
      <ShipmentContainer shipmentType={shipmentType}>
        <ShipmentHeading {...newHeadingInfo} />
      </ShipmentContainer>,
    );

    const shipmentContainer = await screen.findByTestId('ShipmentContainer');

    expect(shipmentContainer.className).toContain(expectedClass);
  });
});
