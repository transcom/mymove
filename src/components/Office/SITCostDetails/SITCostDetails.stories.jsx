import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import SITCostDetails from './SITCostDetails';

import styles from 'pages/Office/ServicesCounselingMoveInfo/ServicesCounselingTab.module.scss';

export default {
  title: 'Office Components/SIT Cost Details',
  component: SITCostDetails,
  decorators: [
    (Story) => (
      <GridContainer className={styles.gridContainer}>
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <Story />
          </Grid>
        </Grid>
      </GridContainer>
    ),
  ],
};

export const Details = (args) => <SITCostDetails {...args} />;

Details.args = {
  cost: 123400,
  weight: 2345,
  location: '12345',
  sitLocation: 'DESTINATION',
  departureDate: '2022-10-29',
  entryDate: '2022-08-06',
};
