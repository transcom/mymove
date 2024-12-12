/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { action } from '@storybook/addon-actions';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import ShipmentForm from './ShipmentForm';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { roleTypes } from 'constants/userRoles';
import styles from 'pages/Office/ServicesCounselingMoveInfo/ServicesCounselingTab.module.scss';
import { MockProviders } from 'testUtils';
import { servicesCounselingRoutes } from 'constants/routes';

const defaultProps = {
  moveTaskOrderID: 'task123',
  originDutyLocationAddress: {
    city: 'Washington',
    state: 'DC',
    postalCode: '20001',
  },
  newDutyLocationAddress: {
    city: 'Fort Benning',
    state: 'GA',
    postalCode: '31905',
  },
  currentResidence: {
    city: 'Fort Benning',
    state: 'GA',
    postalCode: '31905',
    streetAddress1: '123 Main',
  },
  useCurrentResidence: false,
  mtoShipment: {
    destinationAddress: undefined,
  },
  serviceMember: {
    weightAllotment: {
      totalWeightSelf: 5000,
      ubAllowance: 400,
    },
  },
  isCreatePage: true,
  submitHandler: action('submit MTO Shipment for create or update'),
};

const mockMtoShipment = {
  id: 'mock id',
  moveTaskOrderId: 'mock move id',
  customerRemarks: 'mock customer remarks',
  counselorRemarks: 'mock counselor remarks',
  requestedPickupDate: '2020-03-01',
  requestedDeliveryDate: '2020-03-30',
  agents: [
    {
      firstName: 'mock receiving',
      lastName: 'agent',
      telephone: '2225551234',
      email: 'mock.delivery.agent@example.com',
      agentType: 'RECEIVING_AGENT',
    },
    {
      firstName: 'Mock Releasing',
      lastName: 'Agent Jr, PhD, MD, DDS',
      telephone: '3335551234',
      email: 'mock.pickup.agent@example.com',
      agentType: 'RELEASING_AGENT',
    },
  ],
  pickupAddress: {
    streetAddress1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
  },
  destinationAddress: {
    streetAddress1: '441 SW Rio de la Plata Drive',
    city: 'Tacoma',
    state: 'WA',
    postalCode: '98421',
  },
};

const mockDeliveryAddressUpdate = {
  deliveryAddressUpdate: {
    contractorRemarks: 'Test Contractor Remark',
    id: 'c49f7921-5a6e-46b4-bb39-022583574453',
    newAddress: {
      city: 'Beverly Hills',
      country: 'US',
      eTag: 'MjAyMy0wNy0xN1QxODowODowNi42NTU5MTVa',
      id: '6b57ce91-cabd-4e3b-9f48-ed4627d4878f',
      postalCode: '90210',
      state: 'CA',
      streetAddress1: '123 Any Street',
      streetAddress2: 'P.O. Box 12345',
      streetAddress3: 'c/o Some Person',
    },
    originalAddress: {
      city: 'Fairfield',
      country: 'US',
      eTag: 'MjAyMy0wNy0xN1QxODowODowNi42NDkyNTha',
      id: '92509013-aafc-4892-a476-2e3b97e6933d',
      postalCode: '94535',
      state: 'CA',
      streetAddress1: '987 Any Avenue',
      streetAddress2: 'P.O. Box 9876',
      streetAddress3: 'c/o Some Person',
    },
    shipmentID: '5c84bcf3-92f7-448f-b0e1-e5378b6806df',
    status: 'REQUESTED',
  },
};

const mockPPMShipment = {
  id: '4774f99f-bc94-467a-9469-b6f81657b9ef',
  mtoShipmentId: mockMtoShipment.id,
  expectedDepartureDate: '2022-12-31',
  estimatedWeight: 2000,
  sitExpected: false,
  hasProGear: false,
  estimatedIncentive: 1000000,
  pickupAddress: {
    streetAddress1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
  },
  destinationAddress: {
    streetAddress1: '441 SW Rio de la Plata Drive',
    city: 'Tacoma',
    state: 'WA',
    postalCode: '98421',
  },
};

const mockMtoShipmentTypePPM = {
  id: '2523b014-aac6-4443-8181-6df0a754329b',
  moveTaskOrderId: 'e3b4eb5b-a19d-46a0-b155-dac3c49cc3f6',
  shipmentType: SHIPMENT_OPTIONS.PPM,
  ppmShipment: mockPPMShipment,
};

const mockMtoShipmentNoCustomerRemarks = {
  ...mockMtoShipment,
  customerRemarks: '',
};

export default {
  title: 'Office Components / Forms / ShipmentForm',
  component: ShipmentForm,
  decorators: [
    (Story) => (
      <MockProviders
        path={servicesCounselingRoutes.BASE_SHIPMENT_EDIT_PATH}
        params={{ moveCode: 'move123', shipmentId: 'shipment123' }}
      >
        <div className="officeApp">
          <GridContainer className={styles.gridContainer}>
            <Grid row>
              <Grid col desktop={{ col: 8, offset: 2 }}>
                <Story />
              </Grid>
            </Grid>
          </GridContainer>
        </div>
      </MockProviders>
    ),
  ],
};

const Template = (args) => <ShipmentForm {...args} />;

// create shipment stories (form should not prefill customer data)
export const HHGShipment = Template.bind({});
HHGShipment.args = {
  ...defaultProps,
  shipmentType: SHIPMENT_OPTIONS.HHG,
};

// edit shipment stories (form should prefill)
export const EditHHGShipment = Template.bind({});
EditHHGShipment.args = {
  ...defaultProps,
  shipmentType: SHIPMENT_OPTIONS.HHG,
  isCreatePage: false,
  mtoShipment: mockMtoShipment,
  userRole: roleTypes.SERVICES_COUNSELOR,
};

// edit shipment stories, no customer remarks (form should prefill)
export const EditHHGShipmentNoCustRemarks = Template.bind({});
EditHHGShipmentNoCustRemarks.args = {
  ...defaultProps,
  shipmentType: SHIPMENT_OPTIONS.HHG,
  isCreatePage: false,
  mtoShipment: mockMtoShipmentNoCustomerRemarks,
  userRole: roleTypes.SERVICES_COUNSELOR,
};

export const EditHHGShipmentAsTOOWithDeliveryAddressUpdateRequested = Template.bind({});
EditHHGShipmentAsTOOWithDeliveryAddressUpdateRequested.args = {
  ...defaultProps,
  shipmentType: SHIPMENT_OPTIONS.HHG,
  isCreatePage: false,
  mtoShipment: { ...mockMtoShipment, ...mockDeliveryAddressUpdate },
  userRole: roleTypes.TOO,
};

export const HHGShipmentAsTOO = Template.bind({});
HHGShipmentAsTOO.args = {
  ...defaultProps,
  shipmentType: SHIPMENT_OPTIONS.HHG,
  userRole: roleTypes.TOO,
};

export const NTSShipmentWithoutCodes = Template.bind({});
NTSShipmentWithoutCodes.args = {
  ...defaultProps,
  shipmentType: SHIPMENT_OPTIONS.NTS,
  userRole: roleTypes.SERVICES_COUNSELOR,
};

export const NTSShipmentWithCodes = Template.bind({});
NTSShipmentWithCodes.args = {
  ...defaultProps,
  shipmentType: SHIPMENT_OPTIONS.NTS,
  TACs: { HHG: '1234', NTS: '5678' },
  SACs: { HHG: '000012345' },
  userRole: roleTypes.SERVICES_COUNSELOR,
};

export const NTSReleaseShipment = Template.bind({});
NTSReleaseShipment.args = {
  ...defaultProps,
  shipmentType: SHIPMENT_OPTIONS.NTSR,
  TACs: { HHG: '1234', NTS: '5678' },
  SACs: { HHG: '000012345', NTS: '6789ABC' },
  userRole: roleTypes.SERVICES_COUNSELOR,
};

export const NTSShipmentAsTOO = Template.bind({});
NTSShipmentAsTOO.args = {
  ...defaultProps,
  shipmentType: SHIPMENT_OPTIONS.NTS,
  TACs: { HHG: '1234', NTS: '5678' },
  SACs: { HHG: '000012345' },
  userRole: roleTypes.TOO,
};

export const ExternalVendorShipment = Template.bind({});
ExternalVendorShipment.args = {
  ...defaultProps,
  shipmentType: SHIPMENT_OPTIONS.NTSR,
  TACs: { HHG: '1234', NTS: '5678' },
  SACs: { HHG: '000012345', NTS: '6789ABC' },
  mtoShipment: { ...mockMtoShipment, usesExternalVendor: true },
  userRole: roleTypes.TOO,
};

export const PPMShipment = Template.bind({});
PPMShipment.args = {
  ...defaultProps,
  shipmentType: SHIPMENT_OPTIONS.PPM,
  userRole: roleTypes.SERVICES_COUNSELOR,
};

export const PPMShipmentAdvance = Template.bind({});
PPMShipmentAdvance.args = {
  ...defaultProps,
  isCreatePage: false,
  shipmentType: SHIPMENT_OPTIONS.PPM,
  userRole: roleTypes.SERVICES_COUNSELOR,
  isAdvancePage: true,
  mtoShipment: mockMtoShipmentTypePPM,
};
