import React from 'react';
import { render, screen } from '@testing-library/react';

import {
  futureSITShipment,
  noSITShipment,
  SITShipment,
  futureSITShipmentSITExtension,
} from '../ShipmentSITDisplay/ShipmentSITDisplayTestParams';

import ShipmentDetailsMain from './ShipmentDetailsMain';

import { createPPMShipmentWithFinalIncentive } from 'utils/test/factories/ppmShipment';
import { MockProviders } from 'testUtils';
import { formatDateWithUTC } from 'shared/dates';
import { permissionTypes } from 'constants/permissions';

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

const shipmentDetailsMainParamsSITExtension = {
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
  it('displays SIT when there are service items which contain SIT', () => {
    render(
      <MockProviders>
        <ShipmentDetailsMain {...shipmentDetailsMainParams} shipment={SITShipment} />
      </MockProviders>,
    );

    expect(screen.getByText('SIT (STORAGE IN TRANSIT)')).toBeInTheDocument();
  });
  it('displays SIT when there are service items which contain SIT in the future', () => {
    render(
      <MockProviders>
        <ShipmentDetailsMain {...shipmentDetailsMainParams} shipment={futureSITShipment} />
      </MockProviders>,
    );

    expect(screen.getByText('SIT (STORAGE IN TRANSIT)')).toBeInTheDocument();
  });
  it('does not display SIT when there are no service items which contain SIT', () => {
    render(
      <MockProviders>
        <ShipmentDetailsMain {...shipmentDetailsMainParams} shipment={noSITShipment} />
      </MockProviders>,
    );

    expect(screen.queryByText('SIT (STORAGE IN TRANSIT)')).not.toBeInTheDocument();
  });
  it('renders disabled edit button on SIT panel when move is locked', () => {
    const isMoveLocked = true;
    render(
      <MockProviders permissions={[permissionTypes.updateSITExtension]}>
        <ShipmentDetailsMain {...shipmentDetailsMainParams} shipment={futureSITShipment} isMoveLocked={isMoveLocked} />
      </MockProviders>,
    );

    expect(screen.getByRole('button', { name: 'Edit' })).toBeDisabled();
  });
  it('renders disabled review request button on SIT panel when move is locked', () => {
    const isMoveLocked = true;
    render(
      <MockProviders permissions={[permissionTypes.createSITExtension]}>
        <ShipmentDetailsMain
          {...shipmentDetailsMainParamsSITExtension}
          shipment={futureSITShipmentSITExtension}
          isMoveLocked={isMoveLocked}
        />
      </MockProviders>,
    );
    expect(screen.getByRole('button', { name: 'Review request' })).toBeDisabled();
  });
});
it('does display PPM shipment', () => {
  const ppmShipment = createPPMShipmentWithFinalIncentive({
    ppmShipment: {
      expectedDepartureDate: '2024-01-01',
      actualMoveDate: '2024-02-22',
      estimatedWeight: 100,
      shipment: {
        estimatedIncentive: 2000,
      },
    },
  });
  render(
    <MockProviders>
      <ShipmentDetailsMain {...shipmentDetailsMainParams} shipment={ppmShipment} />
    </MockProviders>,
  );

  expect(screen.queryByText(formatDateWithUTC('2024-01-01'))).toBeInTheDocument();
});

describe('Shipment Details Main - PortTable', () => {
  const poeLocation = {
    portCode: 'PDX',
    portName: 'PORTLAND INTL',
    city: 'PORTLAND',
    state: 'OREGON',
    zip: '97220',
  };

  const podLocation = {
    portCode: 'SEA',
    portName: 'SEATTLE TACOMA INTL',
    city: 'SEATTLE',
    state: 'WASHINGTON',
    zip: '98158',
  };

  it('displays PortTable when poeLocation is provided', () => {
    const shipmentWithPOELocation = {
      ...noSITShipment,
      poeLocation,
      podLocation: null,
    };

    render(
      <MockProviders>
        <ShipmentDetailsMain {...shipmentDetailsMainParams} shipment={shipmentWithPOELocation} />
      </MockProviders>,
    );

    expect(screen.getByText('Port of Embarkation')).toBeInTheDocument();
    expect(screen.getByText('Port of Debarkation')).toBeInTheDocument();
  });

  it('displays PortTable when podLocation is provided', () => {
    const shipmentWithPODLocation = {
      ...noSITShipment,
      poeLocation: null,
      podLocation,
    };

    render(
      <MockProviders>
        <ShipmentDetailsMain {...shipmentDetailsMainParams} shipment={shipmentWithPODLocation} />
      </MockProviders>,
    );

    expect(screen.getByText('Port of Embarkation')).toBeInTheDocument();
    expect(screen.getByText('Port of Debarkation')).toBeInTheDocument();
  });

  it('does not display PortTable when poeLocation and podLocation are not provided', () => {
    const shipmentWithNoPortLocation = {
      ...noSITShipment,
      poeLocation: null,
      podLocation: null,
    };

    render(
      <MockProviders>
        <ShipmentDetailsMain {...shipmentDetailsMainParams} shipment={shipmentWithNoPortLocation} />
      </MockProviders>,
    );

    expect(screen.queryByText('Port of Embarkation')).not.toBeInTheDocument();
    expect(screen.queryByText('Port of Debarkation')).not.toBeInTheDocument();
  });
});
