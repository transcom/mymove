import React from 'react';
import { node } from 'prop-types';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import EditPPMNetWeight from './EditPPMNetWeight';

import { MockProviders } from 'testUtils';

export default {
  title: 'Office Components/PPM/EditPPMNetWeights',
  component: EditPPMNetWeight,
  decorators: [
    (Story) => (
      <MockProviders>
        <GridContainer>
          <Grid row>
            <Grid col desktop={{ col: 2, offset: 8 }}>
              <Story />
            </Grid>
          </Grid>
        </GridContainer>
      </MockProviders>
    ),
  ],
};

const Container = ({ children }) => <div style={{ width: 336, margin: '0 auto' }}>{children}</div>;

Container.propTypes = {
  children: node.isRequired,
};
const Template = (args) => <EditPPMNetWeight {...args} />;

export const EditPPMNetWeightDefault = Template.bind({});

EditPPMNetWeightDefault.args = {
  netWeightRemarks: '',
  moveCode: 'CLOSE0',
  weightTicket: {
    vehicleDescription: 'Kia Forte',
    emptyWeight: 600,
    fullWeight: 1200,
    ownsTrailer: true,
    trailerMeetsCriteria: false,
  },
};

export const EditPPMNetWeightExcessWeight = Template.bind({});
EditPPMNetWeightExcessWeight.args = {
  netWeightRemarks: '',
  moveCode: 'CLOSE0',
  weightTicket: {
    vehicleDescription: 'Kia Forte',
    emptyWeight: 600,
    fullWeight: 1200,
    ownsTrailer: true,
    trailerMeetsCriteria: false,
  },
};

export const EditPPMNetWeightReduceWeight = Template.bind({});
EditPPMNetWeightExcessWeight.args = {
  netWeightRemarks: '',
  moveCode: 'CLOSE0',
  weightTicket: {
    vehicleDescription: 'Kia Forte',
    emptyWeight: 600,
    fullWeight: 1200,
    ownsTrailer: true,
    trailerMeetsCriteria: false,
  },
};
