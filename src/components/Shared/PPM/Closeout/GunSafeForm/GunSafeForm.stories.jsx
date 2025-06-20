import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import GunSafeForm from 'components/Shared/PPM/Closeout/GunSafeForm/GunSafeForm';
import { MockProviders } from 'testUtils';
import { SHIPMENT_OPTIONS } from 'shared/constants';

export default {
  title: 'Shared Components / PPM Closeout / Gun Safe Form',
  component: GunSafeForm,
  decorators: [
    (Story) => (
      <MockProviders>
        <GridContainer>
          <Grid row>
            <Grid col desktop={{ col: 8, offset: 2 }}>
              <Story />
            </Grid>
          </Grid>
        </GridContainer>
      </MockProviders>
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

export const Template = (args) => <GunSafeForm {...args} />;

export const Blank = Template.bind({});
Blank.args = {
  mtoShipment: {
    id: 'f3c29ac7-823a-496a-90dd-b7ab0d4b0ece',
    moveTaskOrderId: 'e9864ee5-56e7-401d-9a7b-a5ea9a83bdea',
    shipmentType: SHIPMENT_OPTIONS.PPM,
    ppmShipment: {},
  },
  entitlements: {
    gunSafe: 1235,
    gunSafeWeight: 8500,
  },
  setNumber: 1,
};

export const ExistingGunSafeWeightTickets = Template.bind({});
ExistingGunSafeWeightTickets.args = {
  mtoShipment: {
    id: 'f3c29ac7-823a-496a-90dd-b7ab0d4b0ece',
    moveTaskOrderId: 'e9864ee5-56e7-401d-9a7b-a5ea9a83bdea',
    shipmentType: SHIPMENT_OPTIONS.PPM,
    ppmShipment: {
      id: '343bb456-63af-4f76-89bd-7403094a5c4d',
    },
  },
  gunSafe: {
    id: '32ecb311-edbe-4fd4-96ee-bd693113f3f3',
    ppmShipmentId: '343bb456-63af-4f76-89bd-7403094a5c4d',
    weight: 145,
    hasWeightTickets: true,
    missingWeightTicket: false,
    description: 'Gun safe test description',
    document: {
      uploads: [
        {
          id: '445d2896-571e-4d2e-8bd1-a9d5878ce21f',
          createdAt: '2022-06-08T07:15:01.000Z',
          bytes: 10240000,
          url: 'some/path/to/',
          filename: 'A very long file name with spaces included.jpg',
          contentType: 'image/jpeg',
          updatedAt: '2023-06-08T07:15:01.000Z',
        },
        {
          id: '445d2896-571e-4d2e-8bd1-a9d5878ce21f',
          createdAt: '2022-06-08T07:15:01.000Z',
          bytes: 10240000,
          url: 'some/path/to/',
          filename: 'flower.jpg',
          contentType: 'image/jpeg',
          updatedAt: '2023-06-08T07:15:01.000Z',
        },
      ],
    },
  },
  entitlements: {
    gunSafe: 1235,
    gunSafeWeight: 8500,
  },
  setNumber: 1,
};
