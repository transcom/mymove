import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import {
  hhgInfo,
  ntsInfo,
  ntsMissingInfo,
  postalOnlyInfo,
  diversionInfo,
  cancelledInfo,
  ntsReleaseInfo,
  ntsReleaseMissingInfo,
  ordersLOA,
} from './ShipmentDisplayTestData';
import ShipmentDisplay from './ShipmentDisplay';

import { MockProviders } from 'testUtils';
import { permissionTypes } from 'constants/permissions';

const mockPush = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useHistory: () => ({
    push: mockPush,
  }),
}));

const errorIfMissingStorageFacility = ['storageFacility'];

describe('Shipment Container', () => {
  describe('HHG Shipment', () => {
    it('renders the container successfully', () => {
      render(
        <ShipmentDisplay
          shipmentId="1"
          displayInfo={hhgInfo}
          ordersLOA={ordersLOA}
          onChange={jest.fn()}
          isSubmitted={false}
        />,
      );
      expect(screen.getByTestId('shipment-display')).toBeInTheDocument();
    });

    it('renders the container successfully with postal only address', () => {
      render(<ShipmentDisplay shipmentId="1" displayInfo={postalOnlyInfo} onChange={jest.fn()} isSubmitted={false} />);
      expect(screen.getByTestId('shipment-display')).toBeInTheDocument();
    });

    it('renders with comments', () => {
      render(<ShipmentDisplay shipmentId="1" displayInfo={hhgInfo} onChange={jest.fn()} isSubmitted={false} />);
      expect(screen.getByText('Counselor remarks')).toBeInTheDocument();
    });

    it('renders with edit button when user has permission', async () => {
      render(
        <MockProviders permissions={[permissionTypes.updateShipment]}>
          <ShipmentDisplay shipmentId="1" displayInfo={hhgInfo} onChange={jest.fn()} isSubmitted={false} editURL="/" />
        </MockProviders>,
      );

      const button = screen.getByRole('button', { name: 'Edit shipment' });
      expect(button).toBeInTheDocument();
      await userEvent.click(button);
      await waitFor(() => {
        expect(mockPush).toHaveBeenCalledWith('/');
      });
    });
    it('renders without edit button when user does not have permissions', () => {
      render(<ShipmentDisplay shipmentId="1" displayInfo={hhgInfo} onChange={jest.fn()} isSubmitted={false} />);
      expect(screen.queryByRole('button', { name: 'Edit shipment' })).not.toBeInTheDocument();
    });
    it('renders with diversion tag', () => {
      render(<ShipmentDisplay shipmentId="1" displayInfo={diversionInfo} onChange={jest.fn()} isSubmitted={false} />);
      expect(screen.getByText('diversion')).toBeInTheDocument();
    });
    it('renders with cancelled tag', () => {
      render(<ShipmentDisplay shipmentId="1" displayInfo={cancelledInfo} onChange={jest.fn()} isSubmitted={false} />);
      expect(screen.getByText('cancelled')).toBeInTheDocument();
    });
  });

  describe('NTS shipment', () => {
    it('renders the container successfully', () => {
      render(
        <MockProviders permissions={[permissionTypes.updateShipment]}>
          <ShipmentDisplay shipmentId="1" displayInfo={ntsInfo} onChange={jest.fn()} isSubmitted editURL="/" />
        </MockProviders>,
      );
      expect(screen.getByTestId('shipment-display')).toBeInTheDocument();
      expect(screen.queryByTestId('checkbox')).toBeInTheDocument();
      expect(screen.queryByRole('button', { name: 'Edit shipment' })).toBeInTheDocument();
      expect(screen.getByTestId('shipment-display-checkbox')).not.toBeDisabled();
    });
    it('renders without the approval checkbox for external vendor shipments', () => {
      render(
        <ShipmentDisplay
          shipmentId="1"
          displayInfo={{ ...ntsInfo, usesExternalVendor: true }}
          onChange={jest.fn()}
          isSubmitted={false}
        />,
      );
      expect(screen.queryByTestId('checkbox')).not.toBeInTheDocument();
      expect(screen.getByText('external vendor')).toBeInTheDocument();
    });
    it('checkbox is disabled when information is missing', () => {
      render(
        <MockProviders permissions={[permissionTypes.updateShipment]}>
          <ShipmentDisplay
            shipmentId="1"
            displayInfo={{ ...ntsMissingInfo }}
            onChange={jest.fn()}
            isSubmitted
            errorIfMissing={errorIfMissingStorageFacility}
          />
        </MockProviders>,
      );
      expect(screen.getByTestId('shipment-display-checkbox')).toBeDisabled();
    });
  });

  describe('NTS-release shipment', () => {
    it('renders the container successfully', () => {
      render(
        <MockProviders permissions={[permissionTypes.updateShipment]}>
          <ShipmentDisplay
            shipmentId="1"
            displayInfo={ntsReleaseInfo}
            ordersLOA={ordersLOA}
            onChange={jest.fn()}
            isSubmitted
            editURL="/"
          />
        </MockProviders>,
      );

      expect(screen.getByTestId('shipment-display')).toBeInTheDocument();
      expect(screen.getByTestId('shipment-display-checkbox')).not.toBeDisabled();
      expect(screen.queryByRole('button', { name: 'Edit shipment' })).toBeInTheDocument();
    });
    it('renders without the approval checkbox for external vendor shipments', () => {
      render(
        <ShipmentDisplay
          shipmentId="1"
          displayInfo={{ ...ntsReleaseInfo, usesExternalVendor: true }}
          ordersLOA={ordersLOA}
          onChange={jest.fn()}
          isSubmitted
        />,
      );
      expect(screen.queryByTestId('checkbox')).not.toBeInTheDocument();
      expect(screen.getByText('external vendor')).toBeInTheDocument();
    });

    it('renders with external vendor tag', () => {
      render(
        <ShipmentDisplay
          shipmentId="1"
          displayInfo={{ ...ntsReleaseInfo, usesExternalVendor: true }}
          onChange={jest.fn()}
          isSubmitted={false}
        />,
      );
      expect(screen.getByText('external vendor')).toBeInTheDocument();
    });
    it('checkbox is disabled when information is missing', () => {
      render(
        <MockProviders permissions={[permissionTypes.updateShipment]}>
          <ShipmentDisplay
            shipmentId="1"
            displayInfo={{ ...ntsReleaseMissingInfo }}
            ordersLOA={ordersLOA}
            onChange={jest.fn()}
            isSubmitted
            errorIfMissing={errorIfMissingStorageFacility}
          />
        </MockProviders>,
      );
      expect(screen.getByTestId('shipment-display-checkbox')).toBeDisabled();
    });
  });
});
