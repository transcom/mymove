import React from 'react';
import { mount } from 'enzyme';

import ShipmentContainer from '../ShipmentContainer/ShipmentContainer';
import ShipmentInfoList from '../DefinitionLists/ShipmentInfoList';
import AllowancesList from '../DefinitionLists/AllowancesList';
import CustomerInfoList from '../DefinitionLists/CustomerInfoList';

import ShipmentApprovalPreview from './ShipmentApprovalPreview';

import { SHIPMENT_OPTIONS } from 'shared/constants';

const shipments = [
  {
    approvedDate: '0001-01-01',
    createdAt: '2020-06-10T15:58:02.404029Z',
    customerRemarks: 'please treat gently',
    counselorRemarks: 'all good',
    destinationAddress: {
      city: 'Fairfield',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zODk0MTJa',
      id: '672ff379-f6e3-48b4-a87d-796713f8f997',
      postalCode: '94535',
      state: 'CA',
      streetAddress1: '987 Any Avenue',
      streetAddress2: 'P.O. Box 9876',
      streetAddress3: 'c/o Some Person',
    },
    eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MDQwMzFa',
    id: 'ce01a5b8-9b44-4511-8a8d-edb60f2a4aea',
    moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
    pickupAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zODQ3Njla',
      id: '1686751b-ab36-43cf-b3c9-c0f467d13c19',
      postalCode: '90210',
      state: 'CA',
      streetAddress1: '123 Any Street',
      streetAddress2: 'P.O. Box 12345',
      streetAddress3: 'c/o Some Person',
    },
    rejectionReason: 'shipment not good enough',
    requestedPickupDate: '2018-03-15',
    scheduledPickupDate: '2018-03-16',
    secondaryDeliveryAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zOTkzMlo=',
      id: '15e8f6cc-e1d7-44b2-b1e0-fcb3d6442831',
      postalCode: '90210',
      state: 'CA',
      streetAddress1: '123 Any Street',
      streetAddress2: 'P.O. Box 12345',
      streetAddress3: 'c/o Some Person',
    },
    secondaryPickupAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi4zOTM4OTZa',
      id: '9b79e0c3-8ed5-4fb8-aa36-95845707d8ee',
      postalCode: '90210',
      state: 'CA',
      streetAddress1: '123 Any Street',
      streetAddress2: 'P.O. Box 12345',
      streetAddress3: 'c/o Some Person',
    },
    shipmentType: SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
    status: 'SUBMITTED',
    updatedAt: '2020-06-10T15:58:02.404031Z',
  },
  {
    approvedDate: '0001-01-01',
    createdAt: '2020-06-10T15:58:02.431993Z',
    customerRemarks: 'please treat gently',
    counselorRemarks: 'all good',
    eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MzE5OTVa',
    id: 'c2f68d97-b960-4c86-a418-c70a0aeba04e',
    moveTaskOrderID: '9c7b255c-2981-4bf8-839f-61c7458e2b4d',
    pickupAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MTMyNDha',
      id: '14b1d10d-b34b-4ec5-80e6-69d885206a2a',
      postalCode: '90210',
      state: 'CA',
      streetAddress1: '123 Any Street',
      streetAddress2: 'P.O. Box 12345',
      streetAddress3: 'c/o Some Person',
    },
    rejectionReason: 'shipment not good enough',
    requestedPickupDate: '2018-03-15',
    scheduledPickupDate: '2018-03-16',
    secondaryDeliveryAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MjYxODVa',
      id: '1a4f6fec-42b9-4dd2-b205-c6770ac7ea27',
      postalCode: '90210',
      state: 'CA',
      streetAddress1: '123 Any Street',
      streetAddress2: 'P.O. Box 12345',
      streetAddress3: 'c/o Some Person',
    },
    secondaryPickupAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMC0wNi0xMFQxNTo1ODowMi40MjIwNzVa',
      id: 'e188f33f-f84d-4f86-954a-938b52e38741',
      postalCode: '90210',
      state: 'CA',
      streetAddress1: '123 Any Street',
      streetAddress2: 'P.O. Box 12345',
      streetAddress3: 'c/o Some Person',
    },
    shipmentType: SHIPMENT_OPTIONS.NTSR,
    status: 'SUBMITTED',
    updatedAt: '2020-06-10T15:58:02.431995Z',
  },
];

const ordersInfo = {
  newDutyLocation: {
    address: {
      city: 'Augusta',
      country: 'United States',
      eTag: 'MjAyMC0wOC0wNlQxNDo1Mjo0MS45NDQ0ODla',
      id: '5ac95be8-0230-47ea-90b4-b0f6f60de364',
      postalCode: '30813',
      state: 'GA',
      streetAddress1: 'Fort Gordon',
    },
    address_id: '5ac95be8-0230-47ea-90b4-b0f6f60de364',
    eTag: 'MjAyMC0wOC0wNlQxNDo1Mjo0MS45NDQ0ODla',
    id: '2d5ada83-e09a-47f8-8de6-83ec51694a86',
    name: 'Fort Gordon',
  },
  currentDutyLocation: {
    address: {
      city: 'Des Moines',
      country: 'US',
      eTag: 'MjAyMC0wOC0wNlQxNDo1MzozMC42NjEwODFa',
      id: '37880d6d-2c78-47f1-a71b-53c0ea1a0107',
      postalCode: '50309',
      state: 'IA',
      streetAddress1: '987 Other Avenue',
      streetAddress2: 'P.O. Box 1234',
      streetAddress3: 'c/o Another Person',
    },
    address_id: '37880d6d-2c78-47f1-a71b-53c0ea1a0107',
    eTag: 'MjAyMC0wOC0wNlQxNDo1MzozMC42Njg5MDFa',
    id: '07282a8f-a496-4648-ae24-119775eef57d',
    name: 'vC6w22RPYC',
  },
  issuedDate: '2018-03-15',
  reportByDate: '2018-08-01',
  departmentIndicator: 'COAST_GUARD',
  ordersNumber: 'ORDER3',
  ordersType: 'PERMANENT_CHANGE_OF_STATION',
  ordersTypeDetail: 'TBD',
  tacMDC: '',
  sacSDN: '',
};

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
    streetAddress1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
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
        ordersInfo={ordersInfo}
        allowancesInfo={allowancesInfo}
        counselingFee
        shipmentManagementFee
      />,
    );
    expect(wrapper.find(ShipmentApprovalPreview).exists()).toBe(true);
    expect(wrapper.find(ShipmentContainer).exists()).toBe(true);
    expect(wrapper.find(AllowancesList).exists()).toBe(true);
    expect(wrapper.find(CustomerInfoList).exists()).toBe(true);

    expect(wrapper.find('h3').at(0).text()).toEqual('Household goods');
    expect(wrapper.find('h3').at(1).text()).toEqual('Non-temp storage release');
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
        ordersInfo={ordersInfo}
        allowancesInfo={allowancesInfo}
        mtoAgents={agents}
        counselingFee
        shipmentManagementFee
      />,
    );
    expect(wrapper.find(ShipmentApprovalPreview).exists()).toBe(true);
    expect(wrapper.find(ShipmentContainer).exists()).toBe(true);
    expect(wrapper.find(ShipmentInfoList).exists()).toBe(true);
    expect(wrapper.find(AllowancesList).exists()).toBe(true);
    expect(wrapper.find(CustomerInfoList).exists()).toBe(true);
  });

  it('renders the buttons successfully', () => {
    const wrapper = mount(
      <ShipmentApprovalPreview
        customerInfo={customerInfo}
        mtoShipments={shipments}
        setIsModalVisible={jest.fn()}
        onSubmit={jest.fn()}
        ordersInfo={ordersInfo}
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
        ordersInfo={ordersInfo}
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
        ordersInfo={ordersInfo}
        customerInfo={customerInfo}
        mtoShipments={shipments}
        setIsModalVisible={jest.fn()}
        onSubmit={jest.fn()}
        allowancesInfo={allowancesInfo}
        counselingFee
        shipmentManagementFee
      />,
    );
    expect(wrapper.find('[data-testid="destinationAddress"]').at(0).text()).toEqual(
      '987 Any Avenue, Fairfield, CA 94535',
    );
    expect(wrapper.find('[data-testid="destinationAddress"]').at(1).text()).toEqual(
      ordersInfo.newDutyLocation.address.postalCode,
    );
  });

  it('renders the customer and counselor remarks', () => {
    const wrapper = mount(
      <ShipmentApprovalPreview
        ordersInfo={ordersInfo}
        customerInfo={customerInfo}
        mtoShipments={shipments}
        setIsModalVisible={jest.fn()}
        onSubmit={jest.fn()}
        allowancesInfo={allowancesInfo}
        counselingFee
        shipmentManagementFee
      />,
    );

    const customerRemarks = wrapper.find('[data-testid="customerRemarks"]');
    const counselorRemarks = wrapper.find('[data-testid="counselorRemarks"]');

    expect(customerRemarks.at(0).text()).toEqual('please treat gently');
    expect(customerRemarks.at(1).text()).toEqual('please treat gently');

    expect(counselorRemarks.at(0).text()).toEqual('all good');
    expect(counselorRemarks.at(1).text()).toEqual('all good');
  });
});
