import React from 'react';
import { text, object, number } from '@storybook/addon-knobs';

import { SERVICE_ITEM_STATUS } from '../../shared/constants';

import ShipmentContainer from './ShipmentContainer';
import ShipmentHeading from './ShipmentHeading';

import RequestedServiceItemsTable from 'components/Office/RequestedServiceItemsTable/RequestedServiceItemsTable';
import ImportantShipmentDates from 'components/Office/ImportantShipmentDates';
import ShipmentAddresses from 'components/Office/ShipmentAddresses/ShipmentAddresses';
import ShipmentWeightDetails from 'components/Office/ShipmentWeightDetails/ShipmentWeightDetails';

export default {
  title: 'Office Components/ShipmentContainer',
  component: ShipmentContainer,
};

export const HHG = () => (
  <ShipmentContainer shipmentType={text('ShipmentContainer.shipmentType', 'HHG')}>
    <ShipmentHeading
      shipmentInfo={{
        shipmentType: text('ShipmentInfo.shipmentType', 'Household Goods'),
        originCity: text('ShipmentInfo.originCity', 'San Antonio'),
        originState: text('ShipmentInfo.originState', 'TX'),
        originPostalCode: text('ShipmentInfo.originPostalCode', '98421'),
        destinationAddress: object('MTOShipment.destinationAddress', {
          street_address_1: '123 Any Street',
          city: 'Tacoma',
          state: 'WA',
          postal_code: '98421',
        }),
        scheduledPickupDate: text('ShipmentInfo.destinationPostalCode', '27 Mar 2020'),
      }}
    />
  </ShipmentContainer>
);

export const MTOAccessorial = () => (
  <ShipmentContainer shipmentType={text('ShipmentContainer.shipmentType', 'HHG')}>
    <ShipmentHeading
      shipmentInfo={{
        shipmentType: text('ShipmentInfo.shipmentType', 'Household Goods'),
        originCity: text('ShipmentInfo.originCity', 'San Antonio'),
        originState: text('ShipmentInfo.originState', 'TX'),
        originPostalCode: text('ShipmentInfo.originPostalCode', '98421'),
        destinationAddress: object('MTOShipment.destinationAddress', {
          street_address_1: '123 Any Street',
          city: 'Tacoma',
          state: 'WA',
          postal_code: '98421',
        }),
        scheduledPickupDate: text('ShipmentInfo.destinationPostalCode', '27 Mar 2020'),
      }}
    />
    <ImportantShipmentDates
      requestedPickupDate={text('MTOShipment.requestedPickupDate', 'Saturday, 14 Mar 2020')}
      scheduledPickupDate={text('MTOShipment.scheduledPickupDate', 'Sunday, 15 Mar 2020')}
    />
    <ShipmentAddresses
      pickupAddress={object('MTOShipment.pickupAddress', {
        street_address_1: '123 Any Street',
        city: 'Beverly Hills',
        state: 'CA',
        postal_code: '90210',
      })}
      destinationAddress={object('MTOShipment.destinationAddress', {
        street_address_1: '987 Any Avenue',
        city: 'Fairfield',
        state: 'CA',
        postal_code: '94535',
      })}
      originDutyStation={object('Order.originDutyStation', {
        street_address_1: '',
        city: 'Fort Knox',
        state: 'KY',
        postal_code: '40121',
      })}
      destinationDutyStation={object('Order.destinationDutyStation', {
        street_address_1: '',
        city: 'Fort Irwin',
        state: 'CA',
        postal_code: '92310',
      })}
    />

    <ShipmentWeightDetails
      estimatedWeight={number('ShipmentWeight.estimatedWeight', 1000)}
      actualWeight={number('ShipmentWeight.actualWeight', 999.99)}
    />

    <RequestedServiceItemsTable
      statusForTableType={SERVICE_ITEM_STATUS.SUBMITTED}
      serviceItems={[
        object('ServiceItem.first', {
          id: '1',
          createdAt: '2020-01-10:00:00:00',
          serviceItem: 'Fuel Surcharge',
          code: 'FSC',
        }),
      ]}
    />
  </ShipmentContainer>
);
