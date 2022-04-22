import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import HHGShipmentCard from 'components/Customer/Review/ShipmentCard/HHGShipmentCard/HHGShipmentCard';
import PPMShipmentCard from 'components/Customer/Review/ShipmentCard/PPMShipmentCard/PPMShipmentCard';
import NTSShipmentCard from 'components/Customer/Review/ShipmentCard/NTSShipmentCard/NTSShipmentCard';
import NTSRShipmentCard from 'components/Customer/Review/ShipmentCard/NTSRShipmentCard/NTSRShipmentCard';
import { SHIPMENT_OPTIONS } from 'shared/constants';

export default {
  title: 'Customer Components / ShipmentCard',
  decorators: [
    (Story) => (
      <GridContainer>
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <Story />
          </Grid>
        </Grid>
      </GridContainer>
    ),
  ],
};

const noop = () => {};
const hhgDefaultProps = {
  moveId: 'testMove123',
  shipmentNumber: 1,
  shipmentType: 'HHG',
  shipmentId: 'ABC123K',
  requestedPickupDate: new Date('01/01/2020').toISOString(),
  pickupLocation: {
    streetAddress1: '17 8th St',
    city: 'New York',
    state: 'NY',
    postalCode: '111111',
  },
  releasingAgent: {
    firstName: 'Jo',
    lastName: 'Xi',
    phone: '(555) 555-5555',
    email: 'jo.xi@email.com',
  },
  requestedDeliveryDate: new Date('03/01/2020').toISOString(),
  destinationZIP: '73523',
  receivingAgent: {
    firstName: 'Dorothy',
    lastName: 'Lagomarsino',
    phone: '(999) 999-9999',
    email: 'dorothy.lagomarsino@email.com',
  },
  remarks:
    'This is 500 characters of customer remarks right here. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.',
  showEditAndDeleteBtn: true,
  onEditClick: noop,
  onDeleteClick: noop,
};

const ntsDefaultProps = {
  moveId: 'testMove123',
  shipmentType: 'HHG_INTO_NTS_DOMESTIC',
  shipmentId: 'ABC123K',
  showEditAndDeleteBtn: true,
  onEditClick: noop,
  onDeleteClick: noop,
  requestedPickupDate: new Date('01/01/2020').toISOString(),
  pickupLocation: {
    streetAddress1: '17 8th St',
    city: 'New York',
    state: 'NY',
    postalCode: '111111',
  },
  releasingAgent: {
    firstName: 'Jo',
    lastName: 'Xi',
    phone: '(555) 555-5555',
    email: 'jo.xi@email.com',
  },
  remarks:
    'This is 500 characters of customer remarks right here. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.',
};

const ntsrDefaultProps = {
  moveId: 'testMove123',
  shipmentNumber: 1,
  shipmentType: 'HHG_OUTOF_NTS_DOMESTIC',
  shipmentId: 'ABC123K',
  showEditAndDeleteBtn: true,
  onEditClick: noop,
  onDeleteClick: noop,
  requestedDeliveryDate: new Date('03/01/2020').toISOString(),
  receivingAgent: {
    firstName: 'Dorothy',
    lastName: 'Lagomarsino',
    phone: '(999) 999-9999',
    email: 'dorothy.lagomarsino@email.com',
  },
  remarks:
    'This is 500 characters of customer remarks right here. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.',
};

const ppmShipmentDefaultProps = {
  shipment: {
    moveTaskOrderID: 'testMove123',
    id: '93964f85-ae24-44a5-9e5b-73d1bc521e69',
    shipmentType: SHIPMENT_OPTIONS.PPM,
    ppmShipment: {
      pickupPostalCode: '00000',
      destinationPostalCode: '11111',
      estimatedWeight: 5000,
      expectedDepartureDate: new Date('01/01/2020').toISOString(),
      sitExpected: true,
      estimatedIncentive: 1000000,
      advance: 600000,
    },
  },
};

const ppmDefaultProps = {
  ...ppmShipmentDefaultProps,
  shipmentNumber: 1,
  showEditAndDeleteBtn: true,
  onEditClick: noop,
  onDeleteClick: noop,
};

const ppmShipmentSecondaryZIPProps = {
  shipment: {
    moveTaskOrderID: 'testMove123',
    id: '93964f85-ae24-44a5-9e5b-73d1bc521e69',
    shipmentType: SHIPMENT_OPTIONS.PPM,
    ppmShipment: {
      pickupPostalCode: '00000',
      secondaryPickupPostalCode: '00001',
      destinationPostalCode: '11111',
      secondaryDestinationPostalCode: '11112',
      estimatedWeight: 5000,
      expectedDepartureDate: new Date('01/01/2020').toISOString(),
      sitExpected: true,
      estimatedIncentive: 1000000,
    },
  },
};

const ppmShipmentProGearProps = {
  shipment: {
    moveTaskOrderID: 'testMove123',
    id: '93964f85-ae24-44a5-9e5b-73d1bc521e69',
    shipmentType: SHIPMENT_OPTIONS.PPM,
    ppmShipment: {
      pickupPostalCode: '00000',
      destinationPostalCode: '11111',
      estimatedWeight: 5000,
      hasProGear: true,
      proGearWeight: 1299,
      spouseProGearWeight: 366,
      expectedDepartureDate: new Date('01/01/2020').toISOString(),
      sitExpected: true,
      estimatedIncentive: 1000000,
    },
  },
};

const ppmShipmentIncompleteProps = {
  shipment: {
    moveTaskOrderID: 'testMove123',
    id: '93964f85-ae24-44a5-9e5b-73d1bc521e69',
    shipmentType: SHIPMENT_OPTIONS.PPM,
    ppmShipment: {
      pickupPostalCode: '00000',
      destinationPostalCode: '11111',
      expectedDepartureDate: new Date('01/01/2020').toISOString(),
    },
  },
};

const secondaryDeliveryAddress = {
  secondaryDeliveryAddress: {
    streetAddress1: '17 8th St',
    city: 'New York',
    state: 'NY',
    postalCode: '111111',
  },
};

const secondaryPickupAddress = {
  secondaryPickupAddress: {
    streetAddress1: '812 S 129th Street',
    city: 'New York',
    state: 'NY',
    postalCode: '111111',
  },
};

const HHGTemplate = (args) => <HHGShipmentCard {...args} />;

export const HHGShipment = HHGTemplate.bind({});
HHGShipment.args = hhgDefaultProps;

export const HHGShipmentWithSecondaryDestinationAddress = HHGTemplate.bind({});
HHGShipmentWithSecondaryDestinationAddress.args = {
  ...hhgDefaultProps,
  ...secondaryDeliveryAddress,
};

export const HHGShipmentWithSecondaryPickupAddress = HHGTemplate.bind({});
HHGShipmentWithSecondaryPickupAddress.args = {
  ...hhgDefaultProps,
  ...secondaryPickupAddress,
};

export const HHGShipmentWithSecondaryAddresses = HHGTemplate.bind({});
HHGShipmentWithSecondaryAddresses.args = {
  ...hhgDefaultProps,
  ...secondaryPickupAddress,
  ...secondaryDeliveryAddress,
};

const NTSTemplate = (args) => <NTSShipmentCard {...args} />;
export const NTSShipment = NTSTemplate.bind({});
NTSShipment.args = {
  ...ntsDefaultProps,
};

export const NTSShipmentWithSecondaryPickupAddress = NTSTemplate.bind({});
NTSShipmentWithSecondaryPickupAddress.args = {
  ...ntsDefaultProps,
  ...secondaryPickupAddress,
};

const NTSRTemplate = (args) => <NTSRShipmentCard {...args} />;
export const NTSRShipment = NTSRTemplate.bind({});
NTSRShipment.args = ntsrDefaultProps;

export const NTSRShipmentWithSecondaryDestinationAddress = NTSRTemplate.bind({});
NTSRShipmentWithSecondaryDestinationAddress.args = {
  ...ntsrDefaultProps,
  ...secondaryDeliveryAddress,
};

const PPMTemplate = (args) => <PPMShipmentCard {...args} />;
export const PPMShipment = PPMTemplate.bind({});
PPMShipment.args = ppmDefaultProps;

export const PPMShipmentWithSecondaryPostalCodes = PPMTemplate.bind({});
PPMShipmentWithSecondaryPostalCodes.args = {
  ...ppmDefaultProps,
  ...ppmShipmentSecondaryZIPProps,
};

export const PPMShipmentWithProGear = PPMTemplate.bind({});
PPMShipmentWithProGear.args = {
  ...ppmDefaultProps,
  ...ppmShipmentProGearProps,
};

export const PPMShipmentIncomplete = PPMTemplate.bind({});
PPMShipmentIncomplete.args = {
  ...ppmDefaultProps,
  ...ppmShipmentIncompleteProps,
};

export const WithBlankRemarks = () => (
  <>
    <HHGShipmentCard {...hhgDefaultProps} remarks="" />
    <PPMShipmentCard {...ppmDefaultProps} />
    <NTSShipmentCard {...ntsDefaultProps} remarks="" />
    <NTSRShipmentCard {...ntsrDefaultProps} remarks="" />
  </>
);
