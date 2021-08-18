import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';

import ShipmentWeightDetails from './ShipmentWeightDetails';

const shipmentInfoReweighRequested = {
  shipmentID: 'shipment1',
  ifMatchEtag: 'etag1',
  reweighID: 'reweighRequestID',
};

const shipmentInfoNoReweigh = {
  shipmentID: 'shipment1',
  ifMatchEtag: 'etag1',
};

const handleRequestReweighModal = jest.fn();

describe('ShipmentWeightDetails', () => {
  it('renders without crashing', async () => {
    render(
      <ShipmentWeightDetails
        estimatedWeight={4500}
        actualWeight={5000}
        shipmentInfo={shipmentInfoNoReweigh}
        handleRequestReweighModal={handleRequestReweighModal}
      />,
    );

    const estWeight = await screen.findByText('Estimated weight');
    expect(estWeight).toBeTruthy();

    const shipWeight = await screen.findByText('Shipment weight');
    expect(shipWeight).toBeTruthy();

    const reweighButton = await screen.findByText('Request reweigh');
    expect(reweighButton).toBeTruthy();
  });

  it('renders with estimated weight', async () => {
    render(
      <ShipmentWeightDetails
        estimatedWeight={11000}
        actualWeight={12000}
        shipmentInfo={shipmentInfoReweighRequested}
        handleRequestReweighModal={handleRequestReweighModal}
      />,
    );

    const estWeight = await screen.findByText('11,000 lbs');
    expect(estWeight).toBeTruthy();
  });

  it('renders with shipment weight', async () => {
    render(
      <ShipmentWeightDetails
        estimatedWeight={11000}
        actualWeight={12000}
        shipmentInfo={shipmentInfoReweighRequested}
        handleRequestReweighModal={handleRequestReweighModal}
      />,
    );

    const shipWeight = await screen.findByText('12,000 lbs');
    expect(shipWeight).toBeTruthy();
  });

  it('calls the submit function when submit button is clicked', async () => {
    render(
      <ShipmentWeightDetails
        estimatedWeight={11000}
        actualWeight={12000}
        shipmentInfo={shipmentInfoNoReweigh}
        handleRequestReweighModal={handleRequestReweighModal}
      />,
    );

    await fireEvent.click(screen.getByText('Request reweigh'));
    expect(handleRequestReweighModal).toHaveBeenCalled();
  });

  it('renders without the reweigh button if a reweigh has been requested', async () => {
    render(
      <ShipmentWeightDetails
        estimatedWeight={11000}
        actualWeight={12000}
        shipmentInfo={shipmentInfoReweighRequested}
        handleRequestReweighModal={handleRequestReweighModal}
      />,
    );

    const reweighButton = await screen.queryByText('Request reweigh');
    expect(reweighButton).toBeFalsy();
  });
});
