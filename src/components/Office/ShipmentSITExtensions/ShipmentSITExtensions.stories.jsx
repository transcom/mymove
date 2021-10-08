import React from 'react';
import MockDate from 'mockdate';
import addons from '@storybook/addons';
import { isHappoRun } from 'happo-plugin-storybook/register';

import ShipmentSITExtensions from './ShipmentSITExtensions';
import {
  SITExtensions,
  SITStatusOrigin,
  SITStatusDestination,
  SITShipment,
  SITStatusWithPastSITOriginServiceItem,
  SITStatusWithPastSITServiceItems,
  SITExtensionsWithComments,
} from './ShipmentSITExtensionsTestParams';

const mockedDate = '2020-12-08T00:00:00.000Z';

export default {
  title: 'Office Components/Shipment SIT',
  decorators: [
    (Story) => {
      if (isHappoRun()) {
        MockDate.set(mockedDate);
        addons.getChannel().on('storyRendered', MockDate.reset);
      }
      return (
        <div style={{ padding: '1em', backgroundColor: '#f9f9f9' }}>
          <Story />
        </div>
      );
    },
  ],
};

export const AtOriginNoPreviousSIT = () => <ShipmentSITExtensions sitStatus={SITStatusOrigin} shipment={SITShipment} />;

export const AtDestinationNoPreviousSIT = () => (
  <ShipmentSITExtensions sitStatus={SITStatusDestination} shipment={SITShipment} />
);

export const AtDestinationPreviousSITAtOrigin = () => (
  <ShipmentSITExtensions sitStatus={SITStatusWithPastSITOriginServiceItem} shipment={SITShipment} />
);

export const AtDestinationPreviousMulitpleSIT = () => (
  <ShipmentSITExtensions sitStatus={SITStatusWithPastSITServiceItems} shipment={SITShipment} />
);

export const AtDestinationPreviousSITAndExtension = () => (
  <ShipmentSITExtensions
    sitExtensions={SITExtensions}
    sitStatus={SITStatusWithPastSITServiceItems}
    shipment={SITShipment}
  />
);

export const AtDestinationPreviousSITAndExtensionWithComments = () => (
  <ShipmentSITExtensions
    sitExtensions={SITExtensionsWithComments}
    sitStatus={SITStatusWithPastSITServiceItems}
    shipment={SITShipment}
  />
);
