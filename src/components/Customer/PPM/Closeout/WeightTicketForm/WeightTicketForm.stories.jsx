import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import WeightTicketForm from 'components/Customer/PPM/Closeout/WeightTicketForm/WeightTicketForm';
import { SHIPMENT_OPTIONS } from 'shared/constants';

export default {
  title: 'Customer Components / PPM Closeout / Weight Ticket Form',
  component: WeightTicketForm,
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

const Template = (args) => <WeightTicketForm {...args} />;

export const Blank = Template.bind({});
Blank.args = {
  mtoShipment: {
    id: 'f3c29ac7-823a-496a-90dd-b7ab0d4b0ece',
    moveTaskOrderId: 'e9864ee5-56e7-401d-9a7b-a5ea9a83bdea',
    shipmentType: SHIPMENT_OPTIONS.PPM,
    ppmShipment: {},
  },
  tripNumber: 1,
};

export const ExistingWeightTickets = Template.bind({});
ExistingWeightTickets.args = {
  mtoShipment: {
    id: 'f3c29ac7-823a-496a-90dd-b7ab0d4b0ece',
    moveTaskOrderId: 'e9864ee5-56e7-401d-9a7b-a5ea9a83bdea',
    shipmentType: SHIPMENT_OPTIONS.PPM,
    ppmShipment: {},
  },
  weightTicket: {
    id: '32ecb311-edbe-4fd4-96ee-bd693113f3f3',
    vehicleDescription: 'DMC Delorean',
    missingEmptyWeightTicket: false,
    emptyWeight: 3456,
    emptyWeightTickets: [
      {
        id: 'db4713ae-6087-4330-8b0d-926b3d65c454',
        created_at: '2022-06-10T12:59:30.000Z',
        bytes: 204800,
        url: 'some/path/to/',
        filename: 'emptyWeight.pdf',
        content_type: 'application/pdf',
      },
    ],
    fullWeight: 6789,
    missingFullWeightTicket: false,
    fullWeightTickets: [
      {
        id: '28e6e387-7b2d-441b-b96f-f9ba7ed6e794',
        created_at: '2022-06-09T06:30:59.000Z',
        bytes: 4096000,
        url: 'some/path/to/',
        filename: 'Alongerfilenamewithoutspacestotestlinebreakdisplay.png',
        content_type: 'image/png',
      },
      {
        id: '28e6e387-7b2d-441b-b96f-f9ba7ed6e794',
        created_at: '2022-06-08T07:15:01.000Z',
        bytes: 10240000,
        url: 'some/path/to/',
        filename: 'A very long file name with spaces included.jpg',
        content_type: 'image/jpeg',
      },
    ],
    hasOwnTrailer: false,
    hasClaimedTrailer: false,
  },
  tripNumber: 1,
};

export const MissingWeightTickets = Template.bind({});
MissingWeightTickets.args = {
  mtoShipment: {
    id: 'f3c29ac7-823a-496a-90dd-b7ab0d4b0ece',
    moveTaskOrderId: 'e9864ee5-56e7-401d-9a7b-a5ea9a83bdea',
    shipmentType: SHIPMENT_OPTIONS.PPM,
    ppmShipment: {},
  },
  weightTicket: {
    id: '32ecb311-edbe-4fd4-96ee-bd693113f3f3',
    vehicleDescription: 'DMC Delorean',
    missingEmptyWeightTicket: true,
    emptyWeight: 3456,
    emptyWeightTickets: [
      {
        id: 'db4713ae-6087-4330-8b0d-926b3d65c454',
        created_at: '2022-06-10T12:59:30.000Z',
        bytes: 2048,
        url: 'some/path/to/',
        filename: 'emptyWeight.xls',
        content_type: 'application/vnd.ms-excel',
      },
    ],
    fullWeight: 6789,
    missingFullWeightTicket: true,
    fullWeightTickets: [
      {
        id: '28e6e387-7b2d-441b-b96f-f9ba7ed6e794',
        created_at: '2022-06-09T06:30:59.000Z',
        bytes: 4096,
        url: 'some/path/to/',
        filename: 'fullWeight.xlsx',
        content_type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
      },
    ],
    hasOwnTrailer: false,
    hasClaimedTrailer: false,
  },
  tripNumber: 1,
};
