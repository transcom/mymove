import React from 'react';
import { act } from 'react-dom/test-utils';
import { mount, shallow } from 'enzyme';

import RequestedShipments from './RequestedShipments';

const shipments = [
  {
    approvedDate: '0001-01-01',
    createdAt: '2020-06-10T15:58:02.404029Z',
    customerRemarks: 'please treat gently',
    destinationAddress: {
      city: 'Fairfield',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zODk0MTJa',
      id: '672ff379-f6e3-48b4-a87d-796713f8f997',
      postal_code: '94535',
      state: 'CA',
      street_address_1: '987 Any Avenue',
      street_address_2: 'P.O. Box 9876',
      street_address_3: 'c/o Some Person',
    },
    eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MDQwMzFa',
    id: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aee',
    moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
    pickupAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zODQ3Njla',
      id: '1686751b-ab36-43cf-b3c9-c0f467d13c19',
      postal_code: '90210',
      state: 'CA',
      street_address_1: '123 Any Street',
      street_address_2: 'P.O. Box 12345',
      street_address_3: 'c/o Some Person',
    },
    rejectionReason: 'shipment not good enough',
    requestedPickupDate: '2018-03-15',
    scheduledPickupDate: '2018-03-16',
    secondaryDeliveryAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zOTkzMlo=',
      id: '15e8f6cc-e1d7-44b2-b1e0-fcb3d6442831',
      postal_code: '90210',
      state: 'CA',
      street_address_1: '123 Any Street',
      street_address_2: 'P.O. Box 12345',
      street_address_3: 'c/o Some Person',
    },
    secondaryPickupAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zOTM4OTZa',
      id: '9b79e0c3-8ed5-4fb8-aa36-95845707d8ee',
      postal_code: '90210',
      state: 'CA',
      street_address_1: '123 Any Street',
      street_address_2: 'P.O. Box 12345',
      street_address_3: 'c/o Some Person',
    },
    shipmentType: 'HHG',
    status: 'SUBMITTED',
    updatedAt: '2020-06-10T15:58:02.404031Z',
  },
  {
    approvedDate: '0001-01-01',
    createdAt: '2020-06-10T15:58:02.431993Z',
    customerRemarks: 'please treat gently',
    destinationAddress: {
      city: 'Fairfield',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MTcyMjZa',
      id: '00a5dfeb-c6a0-4ed8-965c-89943163fee4',
      postal_code: '94535',
      state: 'CA',
      street_address_1: '987 Any Avenue',
      street_address_2: 'P.O. Box 9876',
      street_address_3: 'c/o Some Person',
    },
    eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MzE5OTVa',
    id: 'c2f68d97-b960-4c86-a418-c70a0aeba04e',
    moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
    pickupAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MTMyNDha',
      id: '14b1d10d-b34b-4ec5-80e6-69d885206a2a',
      postal_code: '90210',
      state: 'CA',
      street_address_1: '123 Any Street',
      street_address_2: 'P.O. Box 12345',
      street_address_3: 'c/o Some Person',
    },
    rejectionReason: 'shipment not good enough',
    requestedPickupDate: '2018-03-15',
    scheduledPickupDate: '2018-03-16',
    secondaryDeliveryAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MjYxODVa',
      id: '1a4f6fec-42b9-4dd2-b205-c6770ac7ea27',
      postal_code: '90210',
      state: 'CA',
      street_address_1: '123 Any Street',
      street_address_2: 'P.O. Box 12345',
      street_address_3: 'c/o Some Person',
    },
    secondaryPickupAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MjIwNzVa',
      id: 'e188f33f-f84d-4f86-954a-938b52e38741',
      postal_code: '90210',
      state: 'CA',
      street_address_1: '123 Any Street',
      street_address_2: 'P.O. Box 12345',
      street_address_3: 'c/o Some Person',
    },
    shipmentType: 'HHG',
    status: 'SUBMITTED',
    updatedAt: '2020-06-10T15:58:02.431995Z',
  },
];

const customerInfo = {
  name: 'Smith, Kerry',
  dodId: '9999999999',
  phone: '+1 999-999-9999',
  email: 'ksmith@email.com',
  currentAddress: {
    street_address_1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postal_code: '78234',
  },
  destinationAddress: {
    street_address_1: '441 SW Rio de la Plata Drive',
    city: 'Tacoma',
    state: 'WA',
    postal_code: '98421',
  },
  backupContactName: 'Quinn Ocampo',
  backupContactPhone: '+1 999-999-9999',
  backupContactEmail: 'quinnocampo@myemail.com',
};

const agents = [
  {
    type: 'RELEASING_AGENT',
    name: 'Dorothy Lagomarsino',
    email: 'dorothyl@email.com',
    phone: '+1 999-999-9999',
    shipmentId: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aea',
  },
  {
    type: 'RECEIVING_AGENT',
    name: 'Dorothy Lagomarsino',
    email: 'dorothyl@email.com',
    phone: '+1 999-999-9999',
    shipmentId: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aea',
  },
];

const allowancesInfo = {
  branch: 'NAVY',
  rank: 'E_6',
  weightAllowance: 11000,
  authorizedWeight: 11000,
  progear: 2000,
  spouseProgear: 500,
  storageInTransit: 90,
  dependents: true,
};

const moveTaskOrder = {
  eTag: 'MjAyMC0wNi0yNlQyMDoyMjo0MS43Mjc4NTNa',
  id: '6e8c5ca4-774c-4170-934a-59d22259e480',
};

const approveMTO = jest.fn().mockResolvedValue({ response: { status: 200 } });

const requestedShipmentsComponent = (
  <RequestedShipments
    allowancesInfo={allowancesInfo}
    mtoAgents={agents}
    customerInfo={customerInfo}
    mtoShipments={shipments}
    approveMTO={approveMTO}
    shipmentsStatus="SUBMITTED"
  />
);

describe('RequestedShipments', () => {
  it('renders the container successfully', () => {
    const wrapper = shallow(requestedShipmentsComponent);
    expect(wrapper.find('div[data-cy="requested-shipments"]').exists()).toBe(true);
  });

  it('renders a shipment passed to it', () => {
    const wrapper = mount(requestedShipmentsComponent);
    expect(wrapper.find('div[data-cy="requested-shipments"]').text()).toContain('HHG');
  });

  it('renders the button', () => {
    const wrapper = mount(requestedShipmentsComponent);
    const approveButton = wrapper.find('#shipmentApproveButton');
    expect(approveButton.exists()).toBe(true);
    expect(approveButton.text()).toContain('Approve selected shipments');
    expect(approveButton.html()).toContain('disabled=""');
  });

  it('renders the checkboxes', () => {
    const wrapper = mount(requestedShipmentsComponent);
    expect(wrapper.find('div[data-testid="checkbox"]').exists()).toBe(true);
    expect(wrapper.find('div[data-testid="checkbox"]').length).toEqual(4);
  });

  it('enables the modal button when a shipment and service item are checked', async () => {
    const wrapper = mount(requestedShipmentsComponent);

    await act(async () => {
      wrapper
        .find('input[name="shipments"]')
        .at(0)
        .simulate('change', {
          target: {
            name: 'shipments',
            value: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aee',
          },
        });
    });
    wrapper.update();

    expect(wrapper.find('form button[type="button"]').prop('disabled')).toEqual(true);
    expect(wrapper.find('#approvalConfirmationModal').prop('style')).toHaveProperty('display', 'none');

    await act(async () => {
      wrapper
        .find('input[name="shipmentManagementFee"]')
        .simulate('change', { target: { name: 'shipmentManagementFee', value: true } });
    });
    wrapper.update();

    expect(wrapper.find('form button[type="button"]').prop('disabled')).toBe(false);

    await act(async () => {
      wrapper.find('form button[type="button"]').simulate('click');
    });
    wrapper.update();

    expect(wrapper.find('#approvalConfirmationModal').prop('style')).toHaveProperty('display', 'block');
  });

  it('calls approveMTO onSubmit', async () => {
    const mockOnSubmit = jest.fn((id, eTag) => {
      return new Promise((resolve) => {
        resolve({ response: { status: 200, body: { id, eTag } } });
      });
    });

    const wrapper = mount(
      <RequestedShipments
        mtoShipments={shipments}
        mtoAgents={agents}
        allowancesInfo={allowancesInfo}
        customerInfo={customerInfo}
        moveTaskOrder={moveTaskOrder}
        approveMTO={mockOnSubmit}
        shipmentsStatus="SUBMITTED"
      />,
    );

    // You could take the shortcut and call submit directly as well if providing initial values
    //  wrapper.find('form').simulate('submit');

    // When simulating change events you must pass the target with the id and
    // name for formik to know which value to update
    await act(async () => {
      wrapper
        .find('input[name="shipments"]')
        .at(0)
        .simulate('change', {
          target: {
            name: 'shipments',
            value: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aee',
          },
        });

      wrapper
        .find('input[name="shipmentManagementFee"]')
        .simulate('change', { target: { name: 'shipmentManagementFee', value: true } });

      wrapper
        .find('input[name="counselingFee"]')
        .simulate('change', { target: { name: 'counselingFee', value: true } });

      wrapper.find('form button[type="button"]').simulate('click');

      wrapper.find('button[type="submit"]').simulate('click');
    });

    expect(mockOnSubmit).toHaveBeenCalled();
    expect(mockOnSubmit.mock.calls[0]).toEqual([moveTaskOrder.id, moveTaskOrder.eTag, ['MS', 'CS']]);
  });
});
