import React from 'react';
import { action } from '@storybook/addon-actions';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import FinalCloseoutForm from 'components/Customer/PPM/Closeout/FinalCloseoutForm/FinalCloseoutForm';
import { createPPMShipmentWithFinalIncentive } from 'utils/test/factories/ppmShipment';

export default {
  title: 'Customer Components / PPM Closeout / Final Closeout Form',
  component: FinalCloseoutForm,
};

export const Blank = () => {
  return (
    <GridContainer>
      <Grid row>
        <Grid desktop={{ col: 8, offset: 2 }}>
          <FinalCloseoutForm
            initialValues={{ date: '2022-11-01', signature: '' }}
            onBack={action('back button clicked')}
            onSubmit={action('submit button clicked')}
            mtoShipment={createPPMShipmentWithFinalIncentive()}
          />
        </Grid>
      </Grid>
    </GridContainer>
  );
};

export const WithSignature = () => {
  return (
    <GridContainer>
      <Grid row>
        <Grid desktop={{ col: 8, offset: 2 }}>
          <FinalCloseoutForm
            initialValues={{ date: '2022-11-01', signature: 'Grace Griffin' }}
            onBack={action('back button clicked')}
            onSubmit={action('submit button clicked')}
            mtoShipment={createPPMShipmentWithFinalIncentive()}
          />
        </Grid>
      </Grid>
    </GridContainer>
  );
};
