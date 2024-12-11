import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import {
  hhgInfo,
  ntsInfo,
  ntsMissingInfo,
  postalOnlyInfo,
  diversionInfo,
  canceledInfo,
  ntsReleaseInfo,
  ntsReleaseMissingInfo,
  ordersLOA,
  ppmInfo,
  ppmInfoApprovedOrExcluded,
  ppmInfoRejected,
} from './ShipmentDisplayTestData';
import ShipmentDisplay from './ShipmentDisplay';

import { MockProviders } from 'testUtils';
import { permissionTypes } from 'constants/permissions';
import { SHIPMENT_OPTIONS } from 'shared/constants';

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
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
      expect(screen.getByTestId('shipment-display')).toHaveTextContent('HHG');
      expect(screen.getByTestId('ShipmentContainer')).toHaveTextContent(hhgInfo.shipmentLocator);
    });

    it('renders the container with a heading that has a market code and shipment type', () => {
      render(
        <ShipmentDisplay
          shipmentId="1"
          displayInfo={hhgInfo}
          ordersLOA={ordersLOA}
          onChange={jest.fn()}
          isSubmitted={false}
        />,
      );
      expect(screen.getByTestId('shipment-display')).toHaveTextContent(`${hhgInfo.marketCode}HHG`);
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
        expect(mockNavigate).toHaveBeenCalledWith('/');
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
    it('renders with canceled tag', () => {
      render(<ShipmentDisplay shipmentId="1" displayInfo={canceledInfo} onChange={jest.fn()} isSubmitted={false} />);
      expect(screen.getByText('canceled')).toBeInTheDocument();
    });
    it('renders a disabled button when move is locked', () => {
      render(
        <MockProviders permissions={[permissionTypes.updateShipment]}>
          <ShipmentDisplay
            shipmentId="1"
            displayInfo={hhgInfo}
            onChange={jest.fn()}
            isSubmitted
            editURL="/"
            isMoveLocked
          />
        </MockProviders>,
      );
      expect(screen.getByTestId('shipment-display-checkbox')).toBeDisabled();
      expect(screen.queryByRole('button', { name: 'Edit shipment' })).toBeVisible();
      expect(screen.queryByRole('button', { name: 'Edit shipment' })).toBeDisabled();
    });
  });

  describe('NTS shipment', () => {
    it('renders the container successfully', () => {
      render(
        <MockProviders permissions={[permissionTypes.updateShipment]}>
          <ShipmentDisplay shipmentId="1" displayInfo={ntsInfo} onChange={jest.fn()} isSubmitted editURL="/" />
        </MockProviders>,
      );
      expect(screen.queryByTestId('checkbox')).toBeVisible();
      expect(screen.queryByRole('button', { name: 'Edit shipment' })).toBeVisible();
      expect(screen.getByTestId('shipment-display-checkbox')).not.toBeDisabled();
      expect(screen.getByTestId('shipment-display')).toHaveTextContent('NTS');
      expect(screen.getByTestId('ShipmentContainer')).toHaveTextContent(ntsInfo.shipmentLocator);
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
    it('renders with canceled tag', () => {
      render(<ShipmentDisplay shipmentId="1" displayInfo={canceledInfo} onChange={jest.fn()} isSubmitted={false} />);
      expect(screen.getByText('canceled')).toBeInTheDocument();
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

      expect(screen.getByTestId('shipment-display')).toBeVisible();
      expect(screen.getByTestId('shipment-display-checkbox')).not.toBeDisabled();
      expect(screen.queryByRole('button', { name: 'Edit shipment' })).toBeVisible();
      expect(screen.getByTestId('shipment-display')).toHaveTextContent('NTS-release');
      expect(screen.getByTestId('ShipmentContainer')).toHaveTextContent(ntsReleaseInfo.shipmentLocator);
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

    it('renders with canceled tag', () => {
      render(<ShipmentDisplay shipmentId="1" displayInfo={canceledInfo} onChange={jest.fn()} isSubmitted={false} />);
      expect(screen.getByText('canceled')).toBeInTheDocument();
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

  describe('PPM shipment', () => {
    it('renders the container successfully', () => {
      render(
        <MockProviders permissions={[permissionTypes.updateShipment]}>
          <ShipmentDisplay
            displayInfo={ppmInfo}
            ordersLOA={ordersLOA}
            shipmentType={SHIPMENT_OPTIONS.PPM}
            isSubmitted
            allowApproval={false}
            warnIfMissing={['counselorRemarks']}
            reviewURL="/"
          />
        </MockProviders>,
      );

      expect(screen.queryByRole('button', { name: 'Review documents' })).toBeVisible();
      expect(screen.getByTestId('shipment-display')).toHaveTextContent('PPM');
      expect(screen.getByTestId('ShipmentContainer')).toHaveTextContent(ppmInfo.shipmentLocator);
    });
    it('renders aoa packet link when approved', () => {
      render(
        <MockProviders permissions={[permissionTypes.updateShipment]}>
          <ShipmentDisplay
            displayInfo={{ ...ppmInfo }}
            ordersLOA={ordersLOA}
            shipmentType={SHIPMENT_OPTIONS.PPM}
            isSubmitted
            allowApproval={false}
            warnIfMissing={['counselorRemarks']}
            reviewURL="/"
          />
        </MockProviders>,
      );
      expect(screen.getByTestId('aoaPacketDownload')).toBeInTheDocument();
    });
    it('renders the view documents button successfully', () => {
      render(
        <MockProviders permissions={[permissionTypes.updateShipment]}>
          <ShipmentDisplay
            displayInfo={ppmInfo}
            ordersLOA={ordersLOA}
            shipmentType={SHIPMENT_OPTIONS.PPM}
            isSubmitted
            allowApproval={false}
            warnIfMissing={['counselorRemarks']}
            viewURL="/"
          />
        </MockProviders>,
      );

      expect(screen.queryByRole('button', { name: 'View documents' })).toBeVisible();
      expect(screen.getByTestId('shipment-display')).toHaveTextContent('PPM');
      expect(screen.getByTestId('ShipmentContainer')).toHaveTextContent(ppmInfo.shipmentLocator);
    });
    describe("renders the 'packet ready for download' tag when", () => {
      it('approved', () => {
        render(
          <MockProviders permissions={[permissionTypes.updateShipment]}>
            <ShipmentDisplay
              displayInfo={{ ...ppmInfoApprovedOrExcluded }}
              ordersLOA={ordersLOA}
              shipmentType={SHIPMENT_OPTIONS.PPM}
              isSubmitted
              allowApproval={false}
              warnIfMissing={['counselorRemarks']}
              reviewURL="/"
            />
          </MockProviders>,
        );
        expect(screen.getByTestId('ppmStatusTag')).toBeInTheDocument();
      });
      it('renders with canceled tag', () => {
        render(<ShipmentDisplay shipmentId="1" displayInfo={canceledInfo} onChange={jest.fn()} isSubmitted={false} />);
        expect(screen.getByText('canceled')).toBeInTheDocument();
      });
      it('excluded', () => {
        render(
          <MockProviders permissions={[permissionTypes.updateShipment]}>
            <ShipmentDisplay
              displayInfo={{ ...ppmInfoApprovedOrExcluded }}
              ordersLOA={ordersLOA}
              shipmentType={SHIPMENT_OPTIONS.PPM}
              isSubmitted
              allowApproval={false}
              warnIfMissing={['counselorRemarks']}
              reviewURL="/"
            />
          </MockProviders>,
        );
        expect(screen.getByTestId('ppmStatusTag')).toBeInTheDocument();
      });
      it('rejected', () => {
        render(
          <MockProviders permissions={[permissionTypes.updateShipment]}>
            <ShipmentDisplay
              displayInfo={{ ...ppmInfoRejected }}
              ordersLOA={ordersLOA}
              shipmentType={SHIPMENT_OPTIONS.PPM}
              isSubmitted
              allowApproval={false}
              warnIfMissing={['counselorRemarks']}
              reviewURL="/"
            />
          </MockProviders>,
        );
        expect(screen.getByTestId('ppmStatusTag')).toBeInTheDocument();
      });
    });
    it('renders the Actual Expense Reimbursement & PPM status tags', () => {
      render(
        <MockProviders permissions={[permissionTypes.updateShipment]}>
          <ShipmentDisplay
            displayInfo={{ isActualExpenseReimbursement: true, ...ppmInfo }}
            ordersLOA={ordersLOA}
            shipmentType={SHIPMENT_OPTIONS.PPM}
            isSubmitted
            allowApproval={false}
            warnIfMissing={['counselorRemarks']}
            reviewURL="/"
          />
        </MockProviders>,
      );
      expect(screen.getByTestId('actualReimbursementTag')).toBeInTheDocument();
      expect(screen.getByTestId('ppmStatusTag')).toBeInTheDocument();
    });
  });
});
