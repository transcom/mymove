import React from 'react';
import { fireEvent, render, screen } from '@testing-library/react';

import ShipmentWeightDetails from './ShipmentWeightDetails';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { MockProviders } from 'testUtils';
import { permissionTypes } from 'constants/permissions';

const shipmentInfoReweighRequested = {
  shipmentID: 'shipment1',
  ifMatchEtag: 'etag1',
  reweighID: 'reweighRequestID',
  shipmentType: SHIPMENT_OPTIONS.HHG,
};

const shipmentInfoNoReweigh = {
  shipmentID: 'shipment1',
  ifMatchEtag: 'etag1',
  shipmentType: SHIPMENT_OPTIONS.HHG,
};

const shipmentInfoReweigh = {
  shipmentID: 'shipment1',
  ifMatchEtag: 'etag1',
  reweighID: 'reweighRequestID',
  reweighWeight: 1000,
  shipmentType: SHIPMENT_OPTIONS.HHG,
};

const handleRequestReweighModal = jest.fn();

describe('ShipmentWeightDetails', () => {
  it('renders without crashing', async () => {
    render(
      <MockProviders permissions={[permissionTypes.createReweighRequest]}>
        <ShipmentWeightDetails
          estimatedWeight={4500}
          actualWeight={5000}
          shipmentInfo={shipmentInfoNoReweigh}
          handleRequestReweighModal={handleRequestReweighModal}
        />
      </MockProviders>,
    );

    const estWeight = await screen.findByText('Estimated weight');
    expect(estWeight).toBeTruthy();

    const shipWeight = await screen.findByText('Shipment weight');
    expect(shipWeight).toBeTruthy();

    const reweighButton = await screen.findByText('Request reweigh');
    expect(reweighButton).toBeTruthy();
  });

  it('renders with estimated weight if not an NTSR', async () => {
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

  it('renders without estimated weight if an NTSR', async () => {
    render(
      <ShipmentWeightDetails
        estimatedWeight={null}
        actualWeight={12000}
        shipmentInfo={{ ...shipmentInfoReweighRequested, shipmentType: SHIPMENT_OPTIONS.NTSR }}
        handleRequestReweighModal={handleRequestReweighModal}
      />,
    );

    expect(screen.queryByText('Estimated weight')).not.toBeInTheDocument();
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
      <MockProviders permissions={[permissionTypes.createReweighRequest]}>
        <ShipmentWeightDetails
          estimatedWeight={11000}
          actualWeight={12000}
          shipmentInfo={shipmentInfoNoReweigh}
          handleRequestReweighModal={handleRequestReweighModal}
        />
      </MockProviders>,
    );

    await fireEvent.click(screen.getByText('Request reweigh'));
    expect(handleRequestReweighModal).toHaveBeenCalled();
  });

  it('renders without the reweigh button if a reweigh has been requested', async () => {
    render(
      <MockProviders permissions={[permissionTypes.createReweighRequest]}>
        <ShipmentWeightDetails
          estimatedWeight={11000}
          actualWeight={12000}
          shipmentInfo={shipmentInfoReweighRequested}
          handleRequestReweighModal={handleRequestReweighModal}
        />
      </MockProviders>,
    );

    const reweighButton = await screen.queryByText('Request reweigh');
    const reweighRequestedLabel = await screen.queryByText('reweigh requested');
    const reweighedLabel = await screen.queryByText('reweighed');

    expect(reweighButton).toBeFalsy();
    expect(reweighRequestedLabel).toBeTruthy();
    expect(reweighedLabel).toBeFalsy();
  });

  it('renders without the rewiegh button when the user does not have permission', async () => {
    render(
      <ShipmentWeightDetails
        estimatedWeight={11000}
        actualWeight={12000}
        shipmentInfo={shipmentInfoNoReweigh}
        handleRequestReweighModal={handleRequestReweighModal}
      />,
    );

    const reweighButton = await screen.queryByText('Request reweigh');
    expect(reweighButton).toBeFalsy();
  });

  it('only renders the reweighed label if a shipment has been reweighed', async () => {
    render(
      <ShipmentWeightDetails
        estimatedWeight={11000}
        actualWeight={12000}
        shipmentInfo={shipmentInfoReweigh}
        handleRequestReweighModal={handleRequestReweighModal}
      />,
    );

    const reweighRequestedLabel = await screen.queryByText('reweigh requested');
    const reweighedLabel = await screen.queryByText('reweighed');

    expect(reweighRequestedLabel).toBeFalsy();
    expect(reweighedLabel).toBeTruthy();
  });

  it('renders the lowest of either reweight or actual weight', async () => {
    render(
      <ShipmentWeightDetails
        estimatedWeight={11000}
        actualWeight={12000}
        shipmentInfo={shipmentInfoReweigh}
        handleRequestReweighModal={handleRequestReweighModal}
      />,
    );

    expect(screen.getByText('1,000 lbs')).toBeInTheDocument();
  });
});
