import React from 'react';
import { shallow } from 'enzyme';

import ShipmentHeading from '../ShipmentHeading/ShipmentHeading';

import ShipmentContainer from './ShipmentContainer';

import { shipmentStatuses } from 'constants/shipments';
import { SHIPMENT_OPTIONS } from 'shared/constants';

const headingInfo = {
  shipmentInfo: {
    shipmentID: '1',
    shipmentStatus: shipmentStatuses.SUBMITTED,
    shipmentType: 'Household Goods',
    originCity: 'San Antonio',
    originState: 'TX',
    originPostalCode: '98421',
    destinationCity: 'Tacoma',
    destinationState: 'WA',
    destinationPostalCode: '98421',
    scheduledPickupDate: '27 Mar 2020',
    reweigh: { id: '00000000-0000-0000-0000-000000000000' },
    ifMatchEtag: 'etag',
    moveTaskOrderID: 'mtoID',
  },
  handleShowCancellationModal: jest.fn(),
};

describe('Shipment Container', () => {
  it('renders the container successfully', () => {
    const wrapper = shallow(
      <ShipmentContainer>
        <ShipmentHeading {...headingInfo} />
      </ShipmentContainer>,
    );
    expect(wrapper.find('[data-testid="ShipmentContainer"]').exists()).toBe(true);
  });
  it('renders a child component passed to it', () => {
    const wrapper = shallow(
      <ShipmentContainer>
        <ShipmentHeading {...headingInfo} />
      </ShipmentContainer>,
    );
    expect(wrapper.find(ShipmentHeading).length).toBe(1);
  });
  it('renders a container with className container--accent--hhg', () => {
    let wrapper = shallow(
      <ShipmentContainer shipmentType={SHIPMENT_OPTIONS.HHG}>
        <ShipmentHeading {...headingInfo} />
      </ShipmentContainer>,
    );
    expect(wrapper.find('.container--accent--hhg').length).toBe(1);

    wrapper = shallow(
      <ShipmentContainer shipmentType={SHIPMENT_OPTIONS.HHG_SHORTHAUL_DOMESTIC}>
        <ShipmentHeading {...headingInfo} />
      </ShipmentContainer>,
    );
    expect(wrapper.find('.container--accent--hhg').length).toBe(1);

    wrapper = shallow(
      <ShipmentContainer shipmentType={SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC}>
        <ShipmentHeading {...headingInfo} />
      </ShipmentContainer>,
    );
    expect(wrapper.find('.container--accent--hhg').length).toBe(1);
  });
  it('renders a container with className container--accent--nts', () => {
    const wrapper = shallow(
      <ShipmentContainer shipmentType={SHIPMENT_OPTIONS.NTS}>
        <ShipmentHeading {...headingInfo} />
      </ShipmentContainer>,
    );
    expect(wrapper.find('.container--accent--nts').length).toBe(1);
  });
  it('renders a container with className container--accent--ntsr', () => {
    const wrapper = shallow(
      <ShipmentContainer shipmentType={SHIPMENT_OPTIONS.NTSR}>
        <ShipmentHeading {...headingInfo} />
      </ShipmentContainer>,
    );
    expect(wrapper.find('.container--accent--ntsr').length).toBe(1);
  });
});
