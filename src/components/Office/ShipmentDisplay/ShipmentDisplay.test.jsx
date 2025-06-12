import React from 'react';
import { render, screen, waitFor, act } from '@testing-library/react';
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
import { PPM_TYPES, SHIPMENT_OPTIONS } from 'shared/constants';
import { isBooleanFlagEnabled } from 'utils/featureFlags';
import { shipmentStatuses } from 'constants/shipments';

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve(false)),
}));

const errorIfMissingStorageFacility = ['storageFacility'];

describe('Shipment Container', () => {
  describe('HHG Shipment', () => {
    it('renders the container successfully', () => {
      render(
        <MockProviders>
          <ShipmentDisplay
            shipmentId="1"
            displayInfo={hhgInfo}
            ordersLOA={ordersLOA}
            onChange={jest.fn()}
            isSubmitted={false}
          />
          ,
        </MockProviders>,
      );
      expect(screen.getByTestId('shipment-display')).toHaveTextContent('HHG');
      expect(screen.getByTestId('ShipmentContainer')).toHaveTextContent(hhgInfo.shipmentLocator);
    });

    it('renders the container with a heading that has a market code and shipment type', () => {
      render(
        <MockProviders>
          <ShipmentDisplay
            shipmentId="1"
            displayInfo={hhgInfo}
            ordersLOA={ordersLOA}
            onChange={jest.fn()}
            isSubmitted={false}
          />
          ,
        </MockProviders>,
      );
      expect(screen.getByTestId('shipment-display')).toHaveTextContent(`${hhgInfo.marketCode}HHG`);
    });

    it('renders the container successfully with postal only address', () => {
      render(
        <MockProviders>
          <ShipmentDisplay shipmentId="1" displayInfo={postalOnlyInfo} onChange={jest.fn()} isSubmitted={false} />
        </MockProviders>,
      );
      expect(screen.getByTestId('shipment-display')).toBeInTheDocument();
    });

    it('renders with comments', () => {
      render(
        <MockProviders>
          <ShipmentDisplay shipmentId="1" displayInfo={hhgInfo} onChange={jest.fn()} isSubmitted={false} />
        </MockProviders>,
      );
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
      render(
        <MockProviders>
          <ShipmentDisplay shipmentId="1" displayInfo={hhgInfo} onChange={jest.fn()} isSubmitted={false} />
        </MockProviders>,
      );
      expect(screen.queryByRole('button', { name: 'Edit shipment' })).not.toBeInTheDocument();
    });
    it('renders with diversion tag', () => {
      render(
        <MockProviders>
          <ShipmentDisplay shipmentId="1" displayInfo={diversionInfo} onChange={jest.fn()} isSubmitted={false} />
        </MockProviders>,
      );
      expect(screen.getByText('diversion')).toBeInTheDocument();
    });
    it('renders with canceled tag', () => {
      render(
        <MockProviders>
          <ShipmentDisplay shipmentId="1" displayInfo={canceledInfo} onChange={jest.fn()} isSubmitted={false} />
        </MockProviders>,
      );
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
    it('does not render a disabled button when shipment is not terminated', () => {
      render(
        <MockProviders permissions={[permissionTypes.updateShipment]}>
          <ShipmentDisplay shipmentId="1" displayInfo={hhgInfo} onChange={jest.fn()} isSubmitted editURL="/" />
        </MockProviders>,
      );
      expect(screen.getByTestId('shipment-display-checkbox')).toBeEnabled();
      expect(screen.queryByRole('button', { name: 'Edit shipment' })).toBeVisible();
      expect(screen.queryByRole('button', { name: 'Edit shipment' })).toBeEnabled();
    });
    it('renders a disabled button when shipment is terminated', () => {
      render(
        <MockProviders permissions={[permissionTypes.updateShipment]}>
          <ShipmentDisplay
            shipmentId="1"
            displayInfo={{ ...hhgInfo, shipmentStatus: shipmentStatuses.TERMINATED_FOR_CAUSE }}
            onChange={jest.fn()}
            isSubmitted
            editURL="/"
          />
        </MockProviders>,
      );
      expect(screen.queryByRole('button', { name: 'Edit shipment' })).toBeVisible();
      expect(screen.queryByRole('button', { name: 'Edit shipment' })).toBeDisabled();
    });
    it('renders the terminated for cause tag when shipment status allows', async () => {
      const hhgInfoTerminated = { ...hhgInfo, shipmentStatus: shipmentStatuses.TERMINATED_FOR_CAUSE };

      render(
        <MockProviders permissions={[permissionTypes.updateShipment]}>
          <ShipmentDisplay
            shipmentId="1"
            displayInfo={hhgInfoTerminated}
            shipmentType={SHIPMENT_OPTIONS.HHG}
            isSubmitted
            allowApproval={false}
          />
        </MockProviders>,
      );

      await waitFor(() => {
        const tag = screen.getByTestId('terminatedTag');
        expect(tag).toBeInTheDocument();
        expect(tag).toHaveTextContent(/terminated for cause/i);
      });
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
        <MockProviders>
          <ShipmentDisplay
            shipmentId="1"
            displayInfo={{ ...ntsInfo, usesExternalVendor: true }}
            onChange={jest.fn()}
            isSubmitted={false}
          />
        </MockProviders>,
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
      render(
        <MockProviders>
          <ShipmentDisplay shipmentId="1" displayInfo={canceledInfo} onChange={jest.fn()} isSubmitted={false} />
        </MockProviders>,
      );
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
        <MockProviders>
          <ShipmentDisplay
            shipmentId="1"
            displayInfo={{ ...ntsReleaseInfo, usesExternalVendor: true }}
            ordersLOA={ordersLOA}
            onChange={jest.fn()}
            isSubmitted
          />
        </MockProviders>,
      );
      expect(screen.queryByTestId('checkbox')).not.toBeInTheDocument();
      expect(screen.getByText('external vendor')).toBeInTheDocument();
    });

    it('renders with canceled tag', () => {
      render(
        <MockProviders>
          <ShipmentDisplay shipmentId="1" displayInfo={canceledInfo} onChange={jest.fn()} isSubmitted={false} />
        </MockProviders>,
      );
      expect(screen.getByText('canceled')).toBeInTheDocument();
    });
    it('renders with external vendor tag', () => {
      render(
        <MockProviders>
          <ShipmentDisplay
            shipmentId="1"
            displayInfo={{ ...ntsReleaseInfo, usesExternalVendor: true }}
            onChange={jest.fn()}
            isSubmitted={false}
          />
        </MockProviders>,
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
    it('renders the Send PPM to the Customer button successfully', async () => {
      await act(async () => {
        render(
          <MockProviders permissions={[permissionTypes.updateShipment]}>
            <ShipmentDisplay
              displayInfo={ppmInfo}
              sendPpmToCustomer={jest.fn()}
              counselorCanEdit={false}
              ordersLOA={ordersLOA}
              shipmentType={SHIPMENT_OPTIONS.PPM}
              isSubmitted
              allowApproval={false}
              warnIfMissing={['counselorRemarks']}
              completePpmForCustomerURL="/"
            />
          </MockProviders>,
        );
      });

      expect(screen.queryByRole('button', { name: 'Send PPM to the Customer' })).toBeVisible();
      expect(screen.getByTestId('shipment-display')).toHaveTextContent('PPM');
      expect(screen.getByTestId('ShipmentContainer')).toHaveTextContent(ppmInfo.shipmentLocator);
    });
    it('renders the Complete PPM on behalf of the Customer button successfully', async () => {
      isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
      await act(async () => {
        render(
          <MockProviders permissions={[permissionTypes.updateShipment]}>
            <ShipmentDisplay
              displayInfo={ppmInfo}
              ordersLOA={ordersLOA}
              shipmentType={SHIPMENT_OPTIONS.PPM}
              isSubmitted
              allowApproval={false}
              warnIfMissing={['counselorRemarks']}
              completePpmForCustomerURL="/"
            />
          </MockProviders>,
        );
      });

      expect(screen.queryByRole('button', { name: 'Complete PPM on behalf of the Customer' })).toBeVisible();
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
        render(
          <MockProviders>
            <ShipmentDisplay shipmentId="1" displayInfo={canceledInfo} onChange={jest.fn()} isSubmitted={false} />
          </MockProviders>,
        );
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
      isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(false));
      ppmInfo.ppmShipment.ppmType = PPM_TYPES.ACTUAL_EXPENSE;
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
      expect(screen.getByTestId('ppmStatusTag')).toBeInTheDocument();
      expect(screen.getByTestId('actualReimbursementTag')).toBeInTheDocument();
    });
    it('renders the Small Package Reimbursement (when FF is on) & PPM status tags', async () => {
      isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
      ppmInfo.ppmShipment.ppmType = PPM_TYPES.SMALL_PACKAGE;
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
      expect(screen.getByTestId('ppmStatusTag')).toBeInTheDocument();
      await waitFor(() => {
        expect(screen.getByTestId('smallPackageTag')).toBeInTheDocument();
      });
    });
  });
});
