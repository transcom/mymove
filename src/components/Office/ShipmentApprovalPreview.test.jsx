import React from 'react';
import { mount } from 'enzyme';

import ShipmentApprovalPreview from './ShipmentApprovalPreview';
import ShipmentContainer from './ShipmentContainer';
import AllowancesTable from './AllowancesTable';
import CustomerInfoTable from './CustomerInfoTable';

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
    id: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aea',
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
      postal_code: '94535',
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
    shipmentType: 'NTS',
    status: 'SUBMITTED',
    updatedAt: '2020-06-10T15:58:02.431995Z',
  },
];

const allowancesInfo = {
  branch: 'Navy',
  rank: 'E-6',
  weightAllowance: 11000,
  authorizedWeight: 11000,
  progear: 2000,
  spouseProgear: 500,
  storageInTransit: 90,
  dependents: true,
};

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

describe('Shipment preview modal', () => {
  it('renders the modal successfully', () => {
    const wrapper = mount(
      <ShipmentApprovalPreview
        customerInfo={customerInfo}
        mtoShipments={shipments}
        setIsModalVisible={jest.fn()}
        onSubmit={jest.fn()}
        allowancesInfo={allowancesInfo}
        counselingFee
        shipmentManagementFee
      />,
    );
    expect(wrapper.find(ShipmentApprovalPreview).exists()).toBe(true);
    expect(wrapper.find(ShipmentContainer).exists()).toBe(true);
    expect(wrapper.find(AllowancesTable).exists()).toBe(true);
    expect(wrapper.find(CustomerInfoTable).exists()).toBe(true);
  });
  it('renders the modal successfully with mtoAgents provided', () => {
    const wrapper = mount(
      <ShipmentApprovalPreview
        customerInfo={customerInfo}
        mtoShipments={shipments}
        setIsModalVisible={() => {
          return true;
        }}
        onSubmit={jest.fn()}
        allowancesInfo={allowancesInfo}
        mtoAgents={agents}
        counselingFee
        shipmentManagementFee
      />,
    );
    expect(wrapper.find(ShipmentApprovalPreview).exists()).toBe(true);
    expect(wrapper.find(ShipmentContainer).exists()).toBe(true);
    expect(wrapper.find(AllowancesTable).exists()).toBe(true);
    expect(wrapper.find(CustomerInfoTable).exists()).toBe(true);
  });

  it('renders the buttons successfully', () => {
    const wrapper = mount(
      <ShipmentApprovalPreview
        customerInfo={customerInfo}
        mtoShipments={shipments}
        setIsModalVisible={jest.fn()}
        onSubmit={jest.fn()}
        allowancesInfo={allowancesInfo}
        mtoAgents={agents}
        counselingFee
        shipmentManagementFee
      />,
    );
    expect(wrapper.find("button[type='submit']").exists()).toBe(true);
    expect(wrapper.find("button[type='reset']").exists()).toBe(true);
  });

  it('attaches onClick listeners', () => {
    const cancelClicked = jest.fn();
    const submitClicked = jest.fn();
    const wrapper = mount(
      <ShipmentApprovalPreview
        customerInfo={customerInfo}
        mtoShipments={shipments}
        setIsModalVisible={cancelClicked}
        onSubmit={submitClicked}
        allowancesInfo={allowancesInfo}
        mtoAgents={agents}
        counselingFee
        shipmentManagementFee
      />,
    );
    wrapper.find('button[type="submit"]').simulate('click');
    expect(submitClicked).toHaveBeenCalled();

    wrapper.find('button[type="reset"]').simulate('click');
    expect(cancelClicked).toHaveBeenCalledTimes(1);

    wrapper.find('[data-testid="closeShipmentApproval"]').simulate('click');
    expect(cancelClicked).toHaveBeenCalledTimes(2);
  });

  it('renders a postal only destination address', () => {
    const wrapper = mount(
      <ShipmentApprovalPreview
        customerInfo={customerInfo}
        mtoShipments={shipments}
        setIsModalVisible={jest.fn()}
        onSubmit={jest.fn()}
        allowancesInfo={allowancesInfo}
        counselingFee
        shipmentManagementFee
      />,
    );
    expect(wrapper.find('[data-testid="destinationAddress"]').at(1).text()).toEqual('94535');
  });
});
