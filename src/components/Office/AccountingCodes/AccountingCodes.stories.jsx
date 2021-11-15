import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import AccountingCodes from './AccountingCodes';

import formStyles from 'styles/form.module.scss';
import styles from 'pages/Office/ServicesCounselingMoveInfo/ServicesCounselingTab.module.scss';

export default {
  title: 'Office Components / Forms / ServicesCounselingShipmentForm / AccountingCodes',
  component: AccountingCodes,
  decorators: [
    (Story) => (
      <GridContainer className={styles.gridContainer}>
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <form className={formStyles.form}>
              <Story />
            </form>
          </Grid>
        </Grid>
      </GridContainer>
    ),
  ],
};

// create shipment stories (form should not prefill customer data)
export const ServicesCounselorView = () => <AccountingCodes />;
