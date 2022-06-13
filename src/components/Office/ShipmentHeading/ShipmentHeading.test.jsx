import React from 'react';
import { shallow, mount } from 'enzyme';

import ShipmentHeading from './ShipmentHeading';

import { MockProviders } from 'testUtils';
import { permissionTypes } from 'constants/permissions';

const shipmentDestinationAddressWithPostalOnly = {
  postalCode: '98421',
};

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
};

describe('Shipment Heading with full destination address', () => {
  const wrapper = shallow(
    <ShipmentHeading
      shipmentInfo={headingInfo}
      handleUpdateMTOShipmentStatus={jest.fn()}
      handleShowCancellationModal={jest.fn()}
    />,
  );

  it('should render the data passed to it within the heading', () => {
    expect(wrapper.find('h2').text()).toEqual('Household Goods');
    expect(wrapper.find('small').text()).toContain('San Antonio, TX 98421');
    expect(wrapper.find('small').text()).toContain('Tacoma, WA 98421');
    expect(wrapper.find('small').text()).toContain('27 Mar 2020');
  });
});

describe('Shipment Heading with missing destination address', () => {
  it("only renders the postalCode of the order's new duty location", () => {
    headingInfo.destinationAddress = shipmentDestinationAddressWithPostalOnly;
    const wrapper = shallow(
      <ShipmentHeading
        shipmentInfo={headingInfo}
        handleUpdateMTOShipmentStatus={jest.fn()}
        handleShowCancellationModal={jest.fn()}
      />,
    );
    expect(wrapper.find('h2').text()).toEqual('Household Goods');
    expect(wrapper.find('small').text()).toContain('San Antonio, TX 98421');
    expect(wrapper.find('small').text()).toContain('98421');
    expect(wrapper.find('small').text()).toContain('27 Mar 2020');
  });
});

describe('Shipment Heading with diverted shipment', () => {
  it('renders the diversion tag next to the shipment type', () => {
    const wrapper = mount(
      <ShipmentHeading
        shipmentInfo={{ isDiversion: true, ...headingInfo }}
        handleUpdateMTOShipmentStatus={jest.fn()}
        handleShowCancellationModal={jest.fn()}
      />,
    );
    expect(wrapper.find('h2').text()).toEqual('Household Goods');
    expect(wrapper.find({ 'data-testid': 'tag' }).text()).toEqual('diversion');
    expect(wrapper.find('small').text()).toContain('San Antonio, TX 98421');
    expect(wrapper.find('small').text()).toContain('98421');
    expect(wrapper.find('small').text()).toContain('27 Mar 2020');
  });
});

describe('Shipment Heading with diversion requested shipment', () => {
  it('renders the diversion requested tag next to the shipment type', () => {
    const wrapper = mount(
      <ShipmentHeading
        shipmentInfo={{ isDiversion: false, ...headingInfo, shipmentStatus: 'DIVERSION_REQUESTED' }}
        handleUpdateMTOShipmentStatus={jest.fn()}
        handleShowCancellationModal={jest.fn()}
      />,
    );
    expect(wrapper.find({ 'data-testid': 'tag' }).text()).toEqual('diversion requested');
  });
});

describe('Shipment Heading with cancelled shipment', () => {
  const wrapper = mount(
    <ShipmentHeading
      shipmentInfo={{ isDiversion: false, ...headingInfo, shipmentStatus: 'CANCELED' }}
      handleUpdateMTOShipmentStatus={jest.fn()}
      handleShowCancellationModal={jest.fn()}
    />,
  );

  it('renders the cancelled tag next to the shipment type', () => {
    expect(wrapper.find({ 'data-testid': 'tag' }).text()).toEqual('cancelled');
  });

  it('hides the request cancellation button', () => {
    expect(wrapper.find('button').length).toBeFalsy();
  });
});

describe('Shipment Heading shows cancellation button with permissions', () => {
  const wrapper = mount(
    <MockProviders permissions={[permissionTypes.requestShipmentCancellation]}>
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

describe('Shipment Heading hides cancellation button without permissions', () => {
  const wrapper = mount(
    <ShipmentHeading
      shipmentInfo={headingInfo}
      handleUpdateMTOShipmentStatus={jest.fn()}
      handleShowCancellationModal={jest.fn()}
    />,
  );

  it('renders withour request shipment cancellation when user does not have permission', () => {
    expect(wrapper.find('button').length).toBeFalsy();
  });
});
