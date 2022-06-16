import React from 'react';
import MockDate from 'mockdate';
import addons from '@storybook/addons';

import ShipmentSITDisplay from './ShipmentSITDisplay';
import {
  SITExtensions,
  SITStatusOrigin,
  SITStatusDestination,
  SITShipment,
  SITStatusWithPastSITOriginServiceItem,
  SITStatusWithPastSITServiceItems,
  SITExtensionsWithComments,
} from './ShipmentSITDisplayTestParams';

import { MockProviders } from 'testUtils';
import { permissionTypes } from 'constants/permissions';

const mockedDate = '2020-12-08T00:00:00.000Z';

export default {
  title: 'Office Components/Shipment SIT',
  decorators: [
    (Story) => {
      MockDate.set(mockedDate);
      addons.getChannel().on('storyRendered', MockDate.reset);
      return (
        <div style={{ padding: '1em', backgroundColor: '#f9f9f9' }}>
          <Story />
        </div>
      );
    },
  ],
};

export const AtOriginNoPreviousSIT = () => (
  <MockProviders permissions={[permissionTypes.updateSITExtension]}>
    <ShipmentSITDisplay sitStatus={SITStatusOrigin} shipment={SITShipment} />
  </MockProviders>
);

export const AtDestinationNoPreviousSIT = () => (
  <MockProviders permissions={[permissionTypes.updateSITExtension]}>
    <ShipmentSITDisplay sitStatus={SITStatusDestination} shipment={SITShipment} />
  </MockProviders>
);

export const AtDestinationPreviousSITAtOrigin = () => (
  <MockProviders permissions={[permissionTypes.updateSITExtension]}>
    <ShipmentSITDisplay sitStatus={SITStatusWithPastSITOriginServiceItem} shipment={SITShipment} />
  </MockProviders>
);

export const AtDestinationPreviousMulitpleSIT = () => (
  <MockProviders permissions={[permissionTypes.updateSITExtension]}>
    <ShipmentSITDisplay sitStatus={SITStatusWithPastSITServiceItems} shipment={SITShipment} />
  </MockProviders>
);

export const AtDestinationPreviousSITAndExtension = () => (
  <MockProviders permissions={[permissionTypes.updateSITExtension]}>
    <ShipmentSITDisplay
      sitExtensions={SITExtensions}
      sitStatus={SITStatusWithPastSITServiceItems}
      shipment={SITShipment}
    />
  </MockProviders>
);

export const AtDestinationPreviousSITAndExtensionWithComments = () => (
  <MockProviders permissions={[permissionTypes.updateSITExtension]}>
    <ShipmentSITDisplay
      sitExtensions={SITExtensionsWithComments}
      sitStatus={SITStatusWithPastSITServiceItems}
      shipment={SITShipment}
    />{' '}
  </MockProviders>
);
