import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';
import { userEvent, within } from '@storybook/testing-library';

import ExpenseForm from 'components/Customer/PPM/Closeout/ExpenseForm/ExpenseForm';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { UnsupportedZipCodePPMErrorMsg } from 'utils/validation';

export default {
  title: 'Customer Components / PPM Closeout / Expenses PPM Form',
  component: ExpenseForm,
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
  argTypes: { onBack: { action: 'back button clicked' }, onSubmit: { action: 'submit button clicked' } },
};

const Template = (args) => <ExpenseForm {...args} />;

export const Blank = Template.bind({});
Blank.args = {
  mtoShipment: {
    id: 'f3c29ac7-823a-496a-90dd-b7ab0d4b0ece',
    moveTaskOrderId: 'e9864ee5-56e7-401d-9a7b-a5ea9a83bdea',
    shipmentType: SHIPMENT_OPTIONS.PPM,
    ppmShipment: {},
  },
  postalCodeValidator: () => {},
};

// export const BlankWithDefaultZIPs = Template.bind({});
// BlankWithDefaultZIPs.storyName = 'Blank With Default ZIPs';
// BlankWithDefaultZIPs.args = {
//   mtoShipment: {
//     id: 'f3c29ac7-823a-496a-90dd-b7ab0d4b0ece',
//     moveTaskOrderId: 'e9864ee5-56e7-401d-9a7b-a5ea9a83bdea',
//     shipmentType: SHIPMENT_OPTIONS.PPM,
//     ppmShipment: {
//       pickupPostalCode: '10001',
//       destinationPostalCode: '10002',
//     },
//   },
//   postalCodeValidator: () => {},
// };

// export const RequiredValues = Template.bind({});
// RequiredValues.args = {
//   mtoShipment: {
//     id: 'f3c29ac7-823a-496a-90dd-b7ab0d4b0ece',
//     moveTaskOrderId: 'e9864ee5-56e7-401d-9a7b-a5ea9a83bdea',
//     shipmentType: SHIPMENT_OPTIONS.PPM,
//     ppmShipment: {
//       actualMoveDate: '2022-05-19',
//       actualPickupPostalCode: '10001',
//       actualDestinationPostalCode: '60652',
//       hasReceivedAdvance: false,
//     },
//   },
//   postalCodeValidator: () => {},
// };

// export const OptionalValues = Template.bind({});
// OptionalValues.args = {
//   mtoShipment: {
//     id: 'f3c29ac7-823a-496a-90dd-b7ab0d4b0ece',
//     moveTaskOrderId: 'e9864ee5-56e7-401d-9a7b-a5ea9a83bdea',
//     shipmentType: SHIPMENT_OPTIONS.PPM,
//     ppmShipment: {
//       actualMoveDate: '2022-05-19',
//       actualPickupPostalCode: '10001',
//       actualDestinationPostalCode: '60652',
//       hasReceivedAdvance: true,
//       advanceAmountReceived: 456700,
//     },
//   },
//   postalCodeValidator: () => {},
// };

// export const RequiredErrors = Template.bind({});
// RequiredErrors.args = {
//   mtoShipment: {
//     id: 'f3c29ac7-823a-496a-90dd-b7ab0d4b0ece',
//     moveTaskOrderId: 'e9864ee5-56e7-401d-9a7b-a5ea9a83bdea',
//     shipmentType: SHIPMENT_OPTIONS.PPM,
//     ppmShipment: {
//       hasReceivedAdvance: true,
//     },
//   },
//   postalCodeValidator: () => {},
// };

// RequiredErrors.play = async ({ canvasElement }) => {
//   // Starts querying the component from its root element
//   const canvas = within(canvasElement);

//   await userEvent.click(canvas.getByText('Save & Continue'));
// };

// export const InvalidZIPs = Template.bind({});
// InvalidZIPs.storyName = 'Invalid ZIPs';
// InvalidZIPs.args = {
//   mtoShipment: {
//     id: 'f3c29ac7-823a-496a-90dd-b7ab0d4b0ece',
//     moveTaskOrderId: 'e9864ee5-56e7-401d-9a7b-a5ea9a83bdea',
//     shipmentType: SHIPMENT_OPTIONS.PPM,
//     ppmShipment: {
//       actualMoveDate: '2022-05-23',
//       actualPickupPostalCode: '10000',
//       actualDestinationPostalCode: '20000',
//     },
//   },
//   postalCodeValidator: () => UnsupportedZipCodePPMErrorMsg,
// };

// InvalidZIPs.play = async ({ canvasElement }) => {
//   // Starts querying the component from its root element
//   const canvas = within(canvasElement);

//   await userEvent.click(canvas.getByText('Save & Continue'));
// };
