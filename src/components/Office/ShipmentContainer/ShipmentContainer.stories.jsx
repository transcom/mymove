import React from 'react';
import { boolean, number, object, text } from '@storybook/addon-knobs';

import ShipmentHeading from '../ShipmentHeading/ShipmentHeading';

import ShipmentContainer from './ShipmentContainer';

import RequestedServiceItemsTable from 'components/Office/RequestedServiceItemsTable/RequestedServiceItemsTable';
import ImportantShipmentDates from 'components/Office/ImportantShipmentDates/ImportantShipmentDates';
import ShipmentAddresses from 'components/Office/ShipmentAddresses/ShipmentAddresses';
import ShipmentWeightDetails from 'components/Office/ShipmentWeightDetails/ShipmentWeightDetails';
import { SERVICE_ITEM_STATUS, SHIPMENT_OPTIONS } from 'shared/constants';
import { MockProviders } from 'testUtils';
import { permissionTypes } from 'constants/permissions';

export default {
  title: 'Office Components/ShipmentContainer',
  component: ShipmentContainer,
};

const shipmentInfoNoReweigh = {
  shipmentID: 'shipment1',
  ifMatchEtag: 'etag1',
  reweighID: '00000000-0000-0000-0000-000000000000',
};

export const HHG = () => (
  <MockProviders permissions={[permissionTypes.createShipmentCancellation]}>
    <ShipmentContainer shipmentType={text('ShipmentContainer.shipmentType', 'HHG')}>
      <ShipmentHeading
        shipmentInfo={{
          shipmentType: text('ShipmentInfo.shipmentType', 'Household Goods'),
          originCity: text('ShipmentInfo.originCity', 'San Antonio'),
          originState: text('ShipmentInfo.originState', 'TX'),
          originPostalCode: text('ShipmentInfo.originPostalCode', '98421'),
          destinationAddress: object('MTOShipment.destinationAddress', {
            streetAddress1: '123 Any Street',
            city: 'Tacoma',
            state: 'WA',
            postalCode: '98421',
          }),
          scheduledPickupDate: text('ShipmentInfo.destinationPostalCode', '27 Mar 2020'),
          reweigh: { id: '00000000-0000-0000-0000-000000000000' },
        }}
      />
    </ShipmentContainer>
  </MockProviders>
);

export const MTOAccessorial = () => (
  <MockProviders permissions={[permissionTypes.createShipmentCancellation]}>
    <ShipmentContainer shipmentType={text('ShipmentContainer.shipmentType', 'HHG')}>
      <ShipmentHeading
        shipmentInfo={{
          shipmentStatus: text('ShipmentInfo.shipmentStatus', 'APPROVED'),
          shipmentType: text('ShipmentInfo.shipmentType', 'Household Goods'),
          originCity: text('ShipmentInfo.originCity', 'San Antonio'),
          originState: text('ShipmentInfo.originState', 'TX'),
          originPostalCode: text('ShipmentInfo.originPostalCode', '98421'),
          destinationAddress: object('MTOShipment.destinationAddress', {
            streetAddress1: '123 Any Street',
            city: 'Tacoma',
            state: 'WA',
            postalCode: '98421',
          }),
          scheduledPickupDate: text('ShipmentInfo.destinationPostalCode', '27 Mar 2020'),
          reweigh: { id: '00000000-0000-0000-0000-000000000000' },
        }}
      />
      <ImportantShipmentDates
        requestedPickupDate={text('MTOShipment.requestedPickupDate', 'Saturday, 14 Mar 2020')}
        scheduledPickupDate={text('MTOShipment.scheduledPickupDate', 'Sunday, 15 Mar 2020')}
      />
      <ShipmentAddresses
        shipmentInfo={object('MTOShipment.shipmentInfo', {
          id: '1',
          eTag: '1',
          status: 'APPROVED',
          shipmentType: SHIPMENT_OPTIONS.HHG,
        })}
        handleDivertShipment={() => {}}
        pickupAddress={object('MTOShipment.pickupAddress', {
          streetAddress1: '123 Any Street',
          city: 'Beverly Hills',
          state: 'CA',
          postalCode: '90210',
        })}
        destinationAddress={object('MTOShipment.destinationAddress', {
          streetAddress1: '987 Any Avenue',
          city: 'Fairfield',
          state: 'CA',
          postalCode: '94535',
        })}
        originDutyLocation={object('Order.originDutyLocation', {
          streetAddress1: '',
          city: 'Fort Knox',
          state: 'KY',
          postalCode: '40121',
        })}
        destinationDutyLocation={object('Order.destinationDutyLocation', {
          streetAddress1: '',
          city: 'Fort Irwin',
          state: 'CA',
          postalCode: '92310',
        })}
      />

      <ShipmentWeightDetails
        estimatedWeight={number('ShipmentWeight.estimatedWeight', 1000)}
        actualWeight={number('ShipmentWeight.actualWeight', 999.99)}
        shipmentInfo={shipmentInfoNoReweigh}
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
  </MockProviders>
);

export const HHGDiversion = () => (
  <MockProviders permissions={[permissionTypes.createShipmentCancellation]}>
    <ShipmentContainer shipmentType={text('ShipmentContainer.shipmentType', 'HHG')}>
      <ShipmentHeading
        shipmentInfo={{
          shipmentType: text('ShipmentInfo.shipmentType', 'Household Goods'),
          isDiversion: boolean('ShipmentInfo.isDiversion', true),
          originCity: text('ShipmentInfo.originCity', 'San Antonio'),
          originState: text('ShipmentInfo.originState', 'TX'),
          originPostalCode: text('ShipmentInfo.originPostalCode', '98421'),
          destinationAddress: object('MTOShipment.destinationAddress', {
            streetAddress1: '123 Any Street',
            city: 'Tacoma',
            state: 'WA',
            postalCode: '98421',
          }),
          scheduledPickupDate: text('ShipmentInfo.destinationPostalCode', '27 Mar 2020'),
          reweigh: { id: '00000000-0000-0000-0000-000000000000' },
        }}
      />
    </ShipmentContainer>
  </MockProviders>
);
