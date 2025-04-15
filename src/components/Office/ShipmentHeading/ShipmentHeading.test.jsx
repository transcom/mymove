import React from 'react';
import { mount } from 'enzyme';
import { render, screen, waitFor } from '@testing-library/react';

import ShipmentHeading from './ShipmentHeading';

import { MockProviders } from 'testUtils';
import { permissionTypes } from 'constants/permissions';
import { shipmentStatuses } from 'constants/shipments';
import { isBooleanFlagEnabled } from 'utils/featureFlags';
import { roleTypes } from 'constants/userRoles';

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve(false)),
}));

const shipmentDestinationAddress = {
  streetAddress1: '123 Main St',
  city: 'Tacoma',
  state: 'WA',
  postalCode: '98421',
};

const headingInfo = {
  shipmentID: '1',
  moveTaskOrderID: '2',
  shipmentType: 'Household Goods',
  originCity: 'San Antonio',
  originState: 'TX',
  originPostalCode: '98421',
  destinationAddress: shipmentDestinationAddress,
  scheduledPickupDate: '27 Mar 2020',
  shipmentStatus: 'SUBMITTED',
  ifMatchEtag: '1234',
  shipmentLocator: 'EVLRPT-01',
};

describe('ShipmentHeading', () => {
  describe('Shipment Heading with diversion requested shipment', () => {
    it('renders the diversion requested tag next to the shipment type', () => {
      const wrapper = mount(
        <MockProviders>
          <ShipmentHeading
            shipmentInfo={{ isDiversion: false, ...headingInfo, shipmentStatus: 'DIVERSION_REQUESTED' }}
            handleUpdateMTOShipmentStatus={jest.fn()}
            handleShowCancellationModal={jest.fn()}
          />
          ,
        </MockProviders>,
      );
      expect(wrapper.find({ 'data-testid': 'tag' }).text()).toEqual('diversion requested');
    });
  });

  describe('Shipment Heading with canceled shipment', () => {
    const wrapper = mount(
      <MockProviders permissions={[permissionTypes.createShipmentCancellation, permissionTypes.updateMTOPage]}>
        <ShipmentHeading
          shipmentInfo={{ isDiversion: false, ...headingInfo, shipmentStatus: 'CANCELED' }}
          handleUpdateMTOShipmentStatus={jest.fn()}
          handleShowCancellationModal={jest.fn()}
        />
      </MockProviders>,
    );

    it('renders the canceled tag next to the shipment type', () => {
      expect(wrapper.find({ 'data-testid': 'tag' }).text()).toEqual('canceled');
    });

    it('hides the request cancellation button', () => {
      expect(wrapper.find({ 'data-testid': 'requestCancellationBtn' }).length).toBeFalsy();
    });
  });

  describe('Shipment Heading with shipment cancellation requested', () => {
    const wrapper = mount(
      <MockProviders permissions={[permissionTypes.createShipmentCancellation, permissionTypes.updateMTOPage]}>
        <ShipmentHeading
          shipmentInfo={{ isDiversion: false, ...headingInfo, shipmentStatus: shipmentStatuses.CANCELLATION_REQUESTED }}
          handleUpdateMTOShipmentStatus={jest.fn()}
          handleShowCancellationModal={jest.fn()}
        />
      </MockProviders>,
    );

    it('renders the cancellation requested tag next to the shipment type', () => {
      expect(wrapper.find({ 'data-testid': 'tag' }).text()).toEqual('Cancellation Requested');
    });

    it('hides the request cancellation button', async () => {
      expect(wrapper.find({ 'data-testid': 'requestCancellationBtn' }).length).toBeFalsy();
    });
  });

  describe('Shipment Heading shows cancellation button with permissions', () => {
    const wrapper = mount(
      <MockProviders permissions={[permissionTypes.createShipmentCancellation, permissionTypes.updateMTOPage]}>
        <ShipmentHeading
          shipmentInfo={headingInfo}
          handleUpdateMTOShipmentStatus={jest.fn()}
          handleShowCancellationModal={jest.fn()}
        />
      </MockProviders>,
    );

    it('renders with request shipment cancellation when user has permission', () => {
      expect(wrapper.find('button').length).toEqual(1);
    });
  });

  describe('Shipment Heading hides cancellation button without any permissions', () => {
    const wrapper = mount(
      <MockProviders>
        <ShipmentHeading
          shipmentInfo={headingInfo}
          handleUpdateMTOShipmentStatus={jest.fn()}
          handleShowCancellationModal={jest.fn()}
        />
        ,
      </MockProviders>,
    );

    it('renders without request shipment cancellation when user does not have any permissions', () => {
      expect(wrapper.find('button').length).toBeFalsy();
    });
  });

  describe('Shipment Heading hides cancellation button when user is missing permission(s)', () => {
    let wrapper = mount(
      <MockProviders permissions={[permissionTypes.createShipmentCancellation]}>
        <ShipmentHeading
          shipmentInfo={headingInfo}
          handleUpdateMTOShipmentStatus={jest.fn()}
          handleShowCancellationModal={jest.fn()}
        />
      </MockProviders>,
    );

    it('renders withour request shipment cancellation when user is missing one permissions', () => {
      expect(wrapper.find({ 'data-testid': 'requestCancellationBtn' }).length).toBeFalsy();
    });

    wrapper = mount(
      <MockProviders permissions={[permissionTypes.updateMTOPage]}>
        <ShipmentHeading
          shipmentInfo={headingInfo}
          handleUpdateMTOShipmentStatus={jest.fn()}
          handleShowCancellationModal={jest.fn()}
        />
      </MockProviders>,
    );

    it('renders withour request shipment cancellation when user does not have both permissions', () => {
      expect(wrapper.find({ 'data-testid': 'requestCancellationBtn' }).length).toBeFalsy();
    });
  });

  describe('Shipment Heading shows cancellation button but disabled when move is locked', () => {
    const isMoveLocked = true;
    const wrapper = mount(
      <MockProviders permissions={[permissionTypes.createShipmentCancellation, permissionTypes.updateMTOPage]}>
        <ShipmentHeading
          shipmentInfo={headingInfo}
          handleUpdateMTOShipmentStatus={jest.fn()}
          handleShowCancellationModal={jest.fn()}
          isMoveLocked={isMoveLocked}
        />
      </MockProviders>,
    );

    it('renders with disabled request shipment cancellation button', () => {
      expect(wrapper.find('button').length).toEqual(1);
      expect(wrapper.find('button[data-testid="requestCancellationBtn"]').prop('disabled')).toBe(true);
    });

    it('renders the terminate shipment button when permissions allow', async () => {
      isBooleanFlagEnabled.mockResolvedValue(true);
      const shipmentInfo = { ...headingInfo, shipmentStatus: shipmentStatuses.APPROVED, actualPickupdate: null };

      const corState = {
        auth: {
          activeRole: roleTypes.CONTRACTING_OFFICER,
        },
      };

      render(
        <MockProviders permissions={[permissionTypes.createShipmentTermination]} initialState={corState}>
          <ShipmentHeading
            shipmentInfo={shipmentInfo}
            handleUpdateMTOShipmentStatus={jest.fn()}
            handleShowCancellationModal={jest.fn()}
            isMoveLocked={isMoveLocked}
            activeRole={roleTypes.CONTRACTING_OFFICER}
          />
        </MockProviders>,
      );

      await waitFor(() => {
        expect(screen.queryByRole('button', { name: 'Terminate Shipment' })).toBeVisible();
      });
    });

    it('does NOT render the terminate shipment button when values meet criteria but roleType is not COR', async () => {
      isBooleanFlagEnabled.mockResolvedValue(true);
      const shipmentInfo = { ...headingInfo, shipmentStatus: shipmentStatuses.APPROVED, actualPickupdate: null };

      const tooState = {
        auth: {
          activeRole: roleTypes.TOO,
        },
      };

      render(
        <MockProviders permissions={[permissionTypes.createShipmentTermination]} initialState={tooState}>
          <ShipmentHeading
            shipmentInfo={shipmentInfo}
            handleUpdateMTOShipmentStatus={jest.fn()}
            handleShowCancellationModal={jest.fn()}
            isMoveLocked={isMoveLocked}
            activeRole={roleTypes.CONTRACTING_OFFICER}
          />
        </MockProviders>,
      );

      await waitFor(() => {
        expect(screen.queryByRole('button', { name: 'Terminate Shipment' })).not.toBeInTheDocument();
      });
    });
  });
});
