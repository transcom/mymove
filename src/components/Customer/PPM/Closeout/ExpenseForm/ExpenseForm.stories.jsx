import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';
import { v4 as uuidv4 } from 'uuid';

import ExpenseForm from 'components/Customer/PPM/Closeout/ExpenseForm/ExpenseForm';
import { SHIPMENT_OPTIONS } from 'shared/constants';

const mockCreateUploadSuccess = (file) => {
  return Promise.resolve({
    id: uuidv4(),
    created_at: '2022-06-22T23:25:50.490Z',
    bytes: file.size,
    url: 'a/fake/path',
    filename: file.name,
    content_type: file.type,
  });
};

const mockUploadComplete = (upload, err, fieldName, values, setFieldValue) => {
  const newValue = {
    id: uuidv4(),
    created_at: '2022-06-22T23:25:50.490Z',
    bytes: upload.file.size,
    url: 'a/fake/path',
    filename: upload.file.name,
    content_type: upload.file.type,
  };
  setFieldValue(fieldName, [...values[`${fieldName}`], newValue]);
};

const mockUploadDelete = (uploadId, fieldName, values, setFieldTouched, setFieldValue) => {
  const remainingUploads = values[`${fieldName}`]?.filter((upload) => upload.id !== uploadId);
  setFieldTouched(fieldName, true, true);
  setFieldValue(fieldName, remainingUploads, true);
};

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
    receiptType: 'packing_materials',
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
  onCreateUpload: mockCreateUploadSuccess,
  onUploadComplete: mockUploadComplete,
  onUploadDelete: mockUploadDelete,
  tripNumber: '1',
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
    receiptType: 'storage',
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
          filename: 'expenseReceipt.pdf',
          content_type: 'application/pdf',
        },
      ],
    },
  },
  onCreateUpload: mockCreateUploadSuccess,
  onUploadComplete: mockUploadComplete,
  onUploadDelete: mockUploadDelete,
  tripNumber: '1',
};
