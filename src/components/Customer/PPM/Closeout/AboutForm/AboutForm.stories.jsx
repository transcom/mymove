import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';
import { userEvent, within } from '@storybook/testing-library';

import AboutForm from 'components/Customer/PPM/Closeout/AboutForm/AboutForm';
import { SHIPMENT_OPTIONS } from 'shared/constants';

export default {
  title: 'Customer Components / PPM Closeout / About PPM Form',
  component: AboutForm,
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

const Template = (args) => <AboutForm {...args} />;

export const BlankAboutPPMForm = Template.bind({});
BlankAboutPPMForm.args = {
  mtoShipment: {
    id: 'f3c29ac7-823a-496a-90dd-b7ab0d4b0ece',
    moveTaskOrderId: 'e9864ee5-56e7-401d-9a7b-a5ea9a83bdea',
    shipmentType: SHIPMENT_OPTIONS.PPM,
    ppmShipment: {},
  },
  postalCodeValidator: () => {},
};

export const RequiredValuesAboutPPMForm = Template.bind({});
RequiredValuesAboutPPMForm.args = {
  mtoShipment: {
    id: 'f3c29ac7-823a-496a-90dd-b7ab0d4b0ece',
    moveTaskOrderId: 'e9864ee5-56e7-401d-9a7b-a5ea9a83bdea',
    shipmentType: SHIPMENT_OPTIONS.PPM,
    ppmShipment: {
      actualMoveDate: '2022-05-19',
      actualPickupPostalCode: '10001',
      actualDestinationPostalCode: '60652',
      hasReceivedAdvance: false,
    },
  },
  postalCodeValidator: () => {},
};

export const OptionalValuesAboutPPMForm = Template.bind({});
OptionalValuesAboutPPMForm.args = {
  mtoShipment: {
    id: 'f3c29ac7-823a-496a-90dd-b7ab0d4b0ece',
    moveTaskOrderId: 'e9864ee5-56e7-401d-9a7b-a5ea9a83bdea',
    shipmentType: SHIPMENT_OPTIONS.PPM,
    ppmShipment: {
      actualMoveDate: '2022-05-19',
      actualPickupPostalCode: '10001',
      actualDestinationPostalCode: '60652',
      hasReceivedAdvance: true,
      advanceAmountReceived: 456700,
    },
  },
  postalCodeValidator: () => {},
};

export const RequiredErrorsAboutPPMForm = Template.bind({});
RequiredErrorsAboutPPMForm.args = {
  mtoShipment: {
    id: 'f3c29ac7-823a-496a-90dd-b7ab0d4b0ece',
    moveTaskOrderId: 'e9864ee5-56e7-401d-9a7b-a5ea9a83bdea',
    shipmentType: SHIPMENT_OPTIONS.PPM,
    ppmShipment: {
      hasReceivedAdvance: true,
    },
  },
  postalCodeValidator: () => {},
};
RequiredErrorsAboutPPMForm.play = async ({ canvasElement }) => {
  // Starts querying the component from its root element
  const canvas = within(canvasElement);

  await userEvent.click(canvas.getByText('Save & Continue'));
};
