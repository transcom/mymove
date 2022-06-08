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
