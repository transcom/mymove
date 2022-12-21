import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { mount } from 'enzyme';

import ShipmentDetailsSidebar from './ShipmentDetailsSidebar';

import { MockProviders } from 'testUtils';
import { LOA_TYPE } from 'shared/constants';
import { formatAccountingCode } from 'utils/shipmentDisplay';
import { permissionTypes } from 'constants/permissions';

const mockPush = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useHistory: () => ({
    push: mockPush,
  }),
  useParams: () => ({
    moveCode: 'TestCode',
  }),
}));

const shipment = {
  mtoAgents: [
    {
      agentType: 'RELEASING_AGENT',
      firstName: 'Quinn',
      lastName: 'Ocampo',
      phone: '999-999-9999',
      email: 'quinnocampo@myemail.com',
    },
    {
      agentType: 'RECEIVING_AGENT',
      firstName: 'Kate',
      lastName: 'Smith',
      phone: '419-555-9999',
      email: 'ksmith@email.com',
    },
  ],
  secondaryAddresses: {
    secondaryPickupAddress: {
      streetAddress1: '444 S 131st St',
      city: 'San Antonio',
      state: 'TX',
      postalCode: '78234',
    },
    secondaryDeliveryAddress: {
      streetAddress1: '7 Q St',
      city: 'Austin',
      state: 'TX',
      postalCode: '78722',
    },
  },
  storageFacility: {
    address: {
      streetAddress1: '456 S 131st St',
      city: 'San Antonio',
      state: 'TX',
      postalCode: '78212',
    },
    lotNumber: '654321',
    facilityName: 'Some storage facility',
  },
  serviceOrderNumber: '1234',
  tacType: LOA_TYPE.HHG,
  sacType: LOA_TYPE.NTS,
  shipmentType: 'HHG',
};

const ordersLOA = {
  tac: '1234',
  sac: '567',
  ntsTac: '8912',
  ntsSac: '345',
};

const headers = [
  'Releasing agent',
  'Receiving agent',
  'Service order number',
  'Facility info and address',
  'Secondary addresses',
  'Pickup',
  'Destination',
  'Accounting codes',
];

describe('Shipment Details Sidebar', () => {
  it('renders all fields when provided', () => {
    render(
      <MockProviders>
        <ShipmentDetailsSidebar
          shipment={shipment}
          ordersLOA={ordersLOA}
          handleEditFacilityInfo={() => {}}
          handleEditAccountingCodes={() => {}}
        />
      </MockProviders>,
    );

    headers.forEach((header) => {
      expect(screen.getByText(header)).toBeInTheDocument();
    });

    expect(screen.getByText(shipment.mtoAgents[0].email)).toBeInTheDocument();
    expect(screen.getByText(shipment.mtoAgents[1].email)).toBeInTheDocument();
    expect(
      screen.getByText(shipment.secondaryAddresses.secondaryPickupAddress.streetAddress1, { exact: false }),
    ).toBeInTheDocument();
    expect(
      screen.getByText(shipment.secondaryAddresses.secondaryDeliveryAddress.streetAddress1, { exact: false }),
    ).toBeInTheDocument();
    expect(screen.getByText(shipment.storageFacility.address.streetAddress1, { exact: false })).toBeInTheDocument();
    expect(screen.getByText(`Lot ${shipment.storageFacility.lotNumber}`)).toBeInTheDocument();
    expect(screen.getByText(shipment.serviceOrderNumber)).toBeInTheDocument();
    expect(screen.getByText(`TAC: ${formatAccountingCode(ordersLOA.tac, shipment.tacType)}`)).toBeInTheDocument();
    expect(screen.getByText(`SAC: ${formatAccountingCode(ordersLOA.ntsSac, shipment.sacType)}`)).toBeInTheDocument();
  });

  it('renders nothing with no info passed in', () => {
    render(
      <MockProviders>
        <ShipmentDetailsSidebar handleEditFacilityInfo={() => {}} handleEditAccountingCodes={() => {}} />
      </MockProviders>,
    );

    headers.forEach((header) => {
      expect(screen.queryByText(header)).toBeNull();
    });
  });

  it('shows edit facility info modal on edit button click', () => {
    render(
      <MockProviders permissions={[permissionTypes.updateShipment]}>
        <ShipmentDetailsSidebar
          shipment={shipment}
          ordersLOA={ordersLOA}
          handleEditFacilityInfo={() => {}}
          handleEditAccountingCodes={() => {}}
        />
      </MockProviders>,
    );

    const openEditFacilityModalButton = screen.getByTestId('edit-facility-info-modal-open');

    userEvent.click(openEditFacilityModalButton);

    // This text is in the edit facility info modal
    expect(screen.getByText('Edit facility info and address')).toBeInTheDocument();
  });

  it('shows edit service order number modal on edit button click', () => {
    render(
      <MockProviders permissions={[permissionTypes.updateShipment]}>
        <ShipmentDetailsSidebar
          shipment={shipment}
          ordersLOA={ordersLOA}
          handleEditFacilityInfo={() => {}}
          handleEditAccountingCodes={() => {}}
        />
      </MockProviders>,
    );

    const openSonModalButton = screen.getByTestId('service-order-number-modal-open');

    userEvent.click(openSonModalButton);

    // This text is in the edit facility info modal
    expect(screen.getByRole('heading', { name: 'Edit service order number' })).toBeInTheDocument();
  });

  it('shows accounting codes modal on edit button click', () => {
    render(
      <MockProviders permissions={[permissionTypes.updateShipment]}>
        <ShipmentDetailsSidebar
          shipment={shipment}
          ordersLOA={ordersLOA}
          handleEditFacilityInfo={() => {}}
          handleEditAccountingCodes={() => {}}
        />
      </MockProviders>,
    );

    const openAccountingCodesModalButton = screen.getByTestId('edit-accounting-codes-modal-open');

    userEvent.click(openAccountingCodesModalButton);

    // This text is in the accounting codes modal
    expect(screen.getByRole('heading', { name: 'Edit accounting codes' })).toBeInTheDocument();
  });

  it('hides edit buttons when permissions are missing', () => {
    const wrapper = mount(
      <MockProviders>
        <ShipmentDetailsSidebar
          shipment={shipment}
          ordersLOA={ordersLOA}
          handleEditFacilityInfo={() => {}}
          handleEditAccountingCodes={() => {}}
        />
      </MockProviders>,
    );

    expect(wrapper.find('button[data-testid="service-order-number-modal-open"]').length).toBeFalsy();
    expect(wrapper.find('button[data-testid="edit-facility-info-modal-open"]').length).toBeFalsy();
    expect(wrapper.find('button[data-testid="edit-accounting-codes-modal-open"]').length).toBeFalsy();
  });
});
