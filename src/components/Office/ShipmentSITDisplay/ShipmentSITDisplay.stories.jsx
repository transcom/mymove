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

export const AtOriginNoPreviousSIT = () => <ShipmentSITDisplay sitStatus={SITStatusOrigin} shipment={SITShipment} />;

export const AtDestinationNoPreviousSIT = () => (
  <ShipmentSITDisplay sitStatus={SITStatusDestination} shipment={SITShipment} />
);

export const AtDestinationPreviousSITAtOrigin = () => (
  <ShipmentSITDisplay sitStatus={SITStatusWithPastSITOriginServiceItem} shipment={SITShipment} />
);

export const AtDestinationPreviousMulitpleSIT = () => (
  <ShipmentSITDisplay sitStatus={SITStatusWithPastSITServiceItems} shipment={SITShipment} />
);

export const AtDestinationPreviousSITAndExtension = () => (
  <ShipmentSITDisplay
    sitExtensions={SITExtensions}
    sitStatus={SITStatusWithPastSITServiceItems}
    shipment={SITShipment}
  />
);

export const AtDestinationPreviousSITAndExtensionWithComments = () => (
  <ShipmentSITDisplay
    sitExtensions={SITExtensionsWithComments}
    sitStatus={SITStatusWithPastSITServiceItems}
    shipment={SITShipment}
  />
);
