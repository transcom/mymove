import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import ExpenseForm from 'components/Customer/PPM/Closeout/ExpenseForm/ExpenseForm';
import { SHIPMENT_OPTIONS } from 'shared/constants';

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
  argTypes: {
    onBack: { action: 'back button clicked' },
    onSubmit: { action: 'submit button clicked' },
    onCreateUpload: { action: 'upload created' },
    onUploadComplete: { action: 'upload completed' },
    onUploadDelete: { action: 'upload deleted' },
  },
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
};

export const ExistingExpenses = Template.bind({});
ExistingExpenses.args = {
  mtoShipment: {
    id: 'f3c29ac7-823a-496a-90dd-b7ab0d4b0ece',
    moveTaskOrderId: 'e9864ee5-56e7-401d-9a7b-a5ea9a83bdea',
    shipmentType: SHIPMENT_OPTIONS.PPM,
    ppmShipment: {
      id: '343bb456-63af-4f76-89bd-7403094a5c4d',
    },
  },
  expense: {
    id: '32ecb311-edbe-4fd4-96ee-bd693113f3f3',
    ppmShipmentId: '343bb456-63af-4f76-89bd-7403094a5c4d',
    expenseType: 'packing_materials',
    description: 'bubble wrap',
    missingReceipt: false,
    paidWithGTCC: false,
    amount: 600,
    receiptDocument: {
      uploads: [
        {
          id: 'db4713ae-6087-4330-8b0d-926b3d65c454',
          created_at: '2022-06-10T12:59:30.000Z',
          bytes: 204800,
          url: 'some/path/to/',
          filename: 'expenseReceipt.pdf',
          content_type: 'application/pdf',
        },
      ],
    },
  },
  receiptNumber: '1',
};

export const SITExpenses = Template.bind({});
SITExpenses.args = {
  mtoShipment: {
    id: 'f3c29ac7-823a-496a-90dd-b7ab0d4b0ece',
    moveTaskOrderId: 'e9864ee5-56e7-401d-9a7b-a5ea9a83bdea',
    shipmentType: SHIPMENT_OPTIONS.PPM,
    ppmShipment: {
      id: '343bb456-63af-4f76-89bd-7403094a5c4d',
    },
  },
  expense: {
    id: '32ecb311-edbe-4fd4-96ee-bd693113f3f3',
    ppmShipmentId: '343bb456-63af-4f76-89bd-7403094a5c4d',
    expenseType: 'storage',
    description: '10x10 storage pod',
    missingReceipt: false,
    paidWithGTCC: false,
    amount: 1600,
    sitStartDate: '2022-09-23',
    sitEndDate: '2022-12-25',
    receiptDocument: {
      uploads: [
        {
          id: 'db4713ae-6087-4330-8b0d-926b3d65c454',
          created_at: '2022-06-10T12:59:30.000Z',
          bytes: 204800,
          url: 'some/path/to/',
          filename: 'uhaulReceipt.pdf',
          content_type: 'application/pdf',
        },
      ],
    },
  },
  receiptNumber: '1',
};
