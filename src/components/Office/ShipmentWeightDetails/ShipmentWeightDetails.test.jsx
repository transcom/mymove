import React from 'react';
import { fireEvent, render, screen } from '@testing-library/react';
import { mount } from 'enzyme';

import ShipmentWeightDetails from './ShipmentWeightDetails';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { MockProviders } from 'testUtils';
import { permissionTypes } from 'constants/permissions';
import { shipmentStatuses } from 'constants/shipments';

const emDash = '\u2014';

const shipmentInfoReweighRequested = {
  shipmentID: 'shipment1',
  ifMatchEtag: 'etag1',
  reweighID: 'reweighRequestID',
  shipmentType: SHIPMENT_OPTIONS.HHG,
  shipmentActualProGearWeight: 800,
  shipmentActualSpouseProGearWeight: 200,
};

const shipmentInfoNoReweigh = {
  shipmentID: 'shipment1',
  ifMatchEtag: 'etag1',
  shipmentType: SHIPMENT_OPTIONS.HHG,
};

const shipmentInfoNoReweighPPM = {
  shipmentID: 'shipment1',
  ifMatchEtag: 'etag1',
  shipmentType: SHIPMENT_OPTIONS.PPM,
};

const shipmentInfoReweigh = {
  shipmentID: 'shipment1',
  ifMatchEtag: 'etag1',
  reweighID: 'reweighRequestID',
  reweighWeight: 1000,
  shipmentType: SHIPMENT_OPTIONS.HHG,
  shipmentActualProGearWeight: 100,
  shipmentActualSpouseProGearWeight: 50,
};

const handleRequestReweighModal = jest.fn();

describe('ShipmentWeightDetails', () => {
  it('renders without crashing', async () => {
    render(
      <MockProviders permissions={[permissionTypes.createReweighRequest, permissionTypes.updateMTOPage]}>
        <ShipmentWeightDetails
          estimatedWeight={4500}
          initialWeight={5000}
          shipmentInfo={shipmentInfoNoReweigh}
          handleRequestReweighModal={handleRequestReweighModal}
        />
      </MockProviders>,
    );

    const estimatedWeight = await screen.findByText('Estimated weight');
    expect(estimatedWeight).toBeTruthy();

    const initialWeight = await screen.findByText('Initial weight');
    expect(initialWeight).toBeTruthy();

    const reweighButton = await screen.findByText('Request reweigh');
    expect(reweighButton).toBeTruthy();

    const actualWeight = await screen.findByText('Actual shipment weight');
    expect(actualWeight).toBeTruthy();

    const actualProGearWeight = await screen.findByText('Actual pro gear weight');
    expect(actualProGearWeight).toBeTruthy();

    const actualSpouseProGearWeight = await screen.findByText('Actual spouse pro gear weight');
    expect(actualSpouseProGearWeight).toBeTruthy();
  });

  it('does not render pro gear for PPMs', async () => {
    render(
      <MockProviders permissions={[permissionTypes.createReweighRequest, permissionTypes.updateMTOPage]}>
        <ShipmentWeightDetails
          estimatedWeight={4500}
          initialWeight={5000}
          shipmentInfo={shipmentInfoNoReweighPPM}
          handleRequestReweighModal={handleRequestReweighModal}
        />
      </MockProviders>,
    );

    const estimatedWeight = await screen.findByText('Estimated weight');
    expect(estimatedWeight).toBeTruthy();

    const initialWeight = await screen.findByText('Initial weight');
    expect(initialWeight).toBeTruthy();

    const reweighButton = screen.queryByText('Request reweigh');
    expect(reweighButton).toBeNull();

    const actualWeight = await screen.findByText('Actual shipment weight');
    expect(actualWeight).toBeTruthy();

    const actualProGearWeight = await screen.queryByText('Actual pro gear weight');
    expect(actualProGearWeight).toBeNull();

    const actualSpouseProGearWeight = await screen.queryByText('Actual spouse pro gear weight');
    expect(actualSpouseProGearWeight).toBeNull();
  });

  it('renders with estimated weight if not an NTSR', async () => {
    render(
      <ShipmentWeightDetails
        estimatedWeight={11000}
        initialWeight={12000}
        shipmentInfo={shipmentInfoReweighRequested}
        handleRequestReweighModal={handleRequestReweighModal}
      />,
    );

    const estWeight = await screen.findByText('11,000 lbs');
    const actualProGearWeight = await screen.findByText('800 lbs');
    const actualSpouseProGearWeight = await screen.findByText('200 lbs');
    expect(estWeight).toBeTruthy();
    expect(actualProGearWeight).toBeTruthy();
    expect(actualSpouseProGearWeight).toBeTruthy();
  });

  it('renders without estimated weight if an NTSR', async () => {
    const wrapper = mount(
      <ShipmentWeightDetails
        estimatedWeight={null}
        initialWeight={12000}
        shipmentInfo={{ ...shipmentInfoReweighRequested, shipmentType: SHIPMENT_OPTIONS.NTSR }}
        handleRequestReweighModal={handleRequestReweighModal}
      />,
    );

    expect(wrapper.find('td').at(0).text()).toEqual(emDash);
  });

  it('renders with shipment weight', async () => {
    render(
      <ShipmentWeightDetails
        estimatedWeight={11000}
        initialWeight={12000}
        shipmentInfo={shipmentInfoReweighRequested}
        handleRequestReweighModal={handleRequestReweighModal}
      />,
    );

    const shipWeight = await screen.findAllByText('12,000 lbs');
    const actualProGearWeight = await screen.findByText('800 lbs');
    const actualSpouseProGearWeight = await screen.findByText('200 lbs');
    expect(shipWeight).toBeTruthy();
    expect(actualProGearWeight).toBeTruthy();
    expect(actualSpouseProGearWeight).toBeTruthy();
  });

  it('calls the submit function when submit button is clicked', async () => {
    render(
      <MockProviders permissions={[permissionTypes.createReweighRequest, permissionTypes.updateMTOPage]}>
        <ShipmentWeightDetails
          estimatedWeight={11000}
          initialWeight={12000}
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
      <MockProviders permissions={[permissionTypes.createReweighRequest, permissionTypes.updateMTOPage]}>
        <ShipmentWeightDetails
          estimatedWeight={11000}
          initialWeight={12000}
          shipmentInfo={shipmentInfoReweighRequested}
          handleRequestReweighModal={handleRequestReweighModal}
        />
      </MockProviders>,
    );

    const reweighButton = await screen.queryByText('Request reweigh');
    const reweighRequestedLabel = await screen.queryByText('reweigh requested');
    const reweighedLabel = await screen.queryByText('reweighed');
    const actualProGearWeight = await screen.findByText('800 lbs');
    const actualSpouseProGearWeight = await screen.findByText('200 lbs');

    expect(reweighButton).toBeFalsy();
    expect(reweighRequestedLabel).toBeTruthy();
    expect(reweighedLabel).toBeFalsy();
    expect(actualProGearWeight).toBeTruthy();
    expect(actualSpouseProGearWeight).toBeTruthy();
  });

  it('renders without the rewiegh button when the user does not have permission', async () => {
    render(
      <ShipmentWeightDetails
        estimatedWeight={11000}
        initialWeight={12000}
        shipmentInfo={shipmentInfoNoReweigh}
        handleRequestReweighModal={handleRequestReweighModal}
      />,
    );

    const reweighButton = await screen.queryByText('Request reweigh');
    expect(reweighButton).toBeFalsy();
  });

  it('renders without the rewiegh button when the user does not have updateMTOPage permission', async () => {
    render(
      <MockProviders permissions={[permissionTypes.createReweighRequest]}>
        <ShipmentWeightDetails
          estimatedWeight={11000}
          initialWeight={12000}
          shipmentInfo={shipmentInfoNoReweigh}
          handleRequestReweighModal={handleRequestReweighModal}
        />
      </MockProviders>,
    );

    const reweighButton = await screen.queryByText('Request reweigh');
    expect(reweighButton).toBeFalsy();
  });

  it('only renders the reweighed label if a shipment has been reweighed', async () => {
    render(
      <ShipmentWeightDetails
        estimatedWeight={11000}
        initialWeight={12000}
        shipmentInfo={shipmentInfoReweigh}
        handleRequestReweighModal={handleRequestReweighModal}
      />,
    );

    const reweighRequestedLabel = await screen.queryByText('reweigh requested');
    const reweighedLabel = await screen.queryByText('reweighed');
    const actualProGearWeight = await screen.findByText('100 lbs');
    const actualSpouseProGearWeight = await screen.findByText('50 lbs');

    expect(reweighRequestedLabel).toBeFalsy();
    expect(reweighedLabel).toBeTruthy();
    expect(actualProGearWeight).toBeTruthy();
    expect(actualSpouseProGearWeight).toBeTruthy();
  });

  it('renders the lowest of either reweight or actual weight', async () => {
    const wrapper = mount(
      <ShipmentWeightDetails
        estimatedWeight={11000}
        initialWeight={12000}
        shipmentInfo={shipmentInfoReweigh}
        handleRequestReweighModal={handleRequestReweighModal}
      />,
    );

    expect(wrapper.find('td').at(2).text()).toEqual('1,000 lbs');
  });

  it('renders with request reweigh button disabled when move is locked', async () => {
    const isMoveLocked = true;
    render(
      <MockProviders permissions={[permissionTypes.createReweighRequest, permissionTypes.updateMTOPage]}>
        <ShipmentWeightDetails
          estimatedWeight={11000}
          initialWeight={12000}
          shipmentInfo={shipmentInfoNoReweigh}
          handleRequestReweighModal={handleRequestReweighModal}
          isMoveLocked={isMoveLocked}
        />
      </MockProviders>,
    );

    expect(screen.getByText('Request reweigh')).toBeDisabled();
  });
  it('renders with request reweigh button disabled when shipment is terminated', async () => {
    render(
      <MockProviders permissions={[permissionTypes.createReweighRequest, permissionTypes.updateMTOPage]}>
        <ShipmentWeightDetails
          estimatedWeight={11000}
          initialWeight={12000}
          shipmentInfo={{ status: shipmentStatuses.TERMINATED_FOR_CAUSE, ...shipmentInfoNoReweigh }}
          handleRequestReweighModal={handleRequestReweighModal}
        />
      </MockProviders>,
    );

    expect(screen.getByText('Request reweigh')).toBeDisabled();
  });
});
