import React from 'react';
import { action } from '@storybook/addon-actions';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import FinalCloseoutForm from 'components/Shared/PPM/Closeout/FinalCloseoutForm/FinalCloseoutForm';
import { createPPMShipmentWithFinalIncentive } from 'utils/test/factories/ppmShipment';
import { PPM_TYPES } from 'shared/constants';
import { APP_NAME } from 'constants/apps';

export default {
  title: 'Shared Components / PPM Closeout / Final Closeout Form',
  component: FinalCloseoutForm,
};

const exampleMove = {
  closeout_office: {
    name: 'Altus AFB',
  },
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
            affiliation="ARMY"
            selectedMove={exampleMove}
            appName={APP_NAME.MYMOVE}
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
            affiliation="ARMY"
            selectedMove={exampleMove}
            appName={APP_NAME.MYMOVE}
          />
        </Grid>
      </Grid>
    </GridContainer>
  );
};

export const NoCloseoutHelperText = () => {
  return (
    <GridContainer>
      <Grid row>
        <Grid desktop={{ col: 8, offset: 2 }}>
          <FinalCloseoutForm
            initialValues={{ date: '2022-11-01', signature: 'Grace Griffin' }}
            onBack={action('back button clicked')}
            onSubmit={action('submit button clicked')}
            mtoShipment={createPPMShipmentWithFinalIncentive()}
            affiliation="COAST_GUARD"
            selectedMove={exampleMove}
            appName={APP_NAME.MYMOVE}
          />
        </Grid>
      </Grid>
    </GridContainer>
  );
};

export const SmallPackagePPM = () => {
  return (
    <GridContainer>
      <Grid row>
        <Grid desktop={{ col: 8, offset: 2 }}>
          <FinalCloseoutForm
            initialValues={{ date: '2022-11-01', signature: '' }}
            onBack={action('back button clicked')}
            onSubmit={action('submit button clicked')}
            mtoShipment={createPPMShipmentWithFinalIncentive({
              ppmShipment: {
                ppmType: PPM_TYPES.SMALL_PACKAGE,
                movingExpenses: [
                  { isProGear: false, weightShipped: 1000, amount: 30000 },
                  { isProGear: true, weightShipped: 500, amount: 20000 },
                ],
              },
            })}
            affiliation="ARMY"
            selectedMove={exampleMove}
            appName={APP_NAME.MYMOVE}
          />
        </Grid>
      </Grid>
    </GridContainer>
  );
};

export const OfficeUserNoAgreement = () => {
  return (
    <GridContainer>
      <Grid row>
        <Grid desktop={{ col: 8, offset: 2 }}>
          <FinalCloseoutForm
            initialValues={{ date: '2022-11-01' }}
            onBack={action('back button clicked')}
            onSubmit={action('submit button clicked')}
            mtoShipment={createPPMShipmentWithFinalIncentive()}
            affiliation="ARMY"
            selectedMove={exampleMove}
            appName={APP_NAME.OFFICE}
          />
        </Grid>
      </Grid>
    </GridContainer>
  );
};
