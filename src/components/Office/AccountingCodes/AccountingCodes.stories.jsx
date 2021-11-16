import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';
import { Formik } from 'formik';

import AccountingCodes from './AccountingCodes';

import styles from 'pages/Office/ServicesCounselingMoveInfo/ServicesCounselingTab.module.scss';
import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';

export default {
  title: 'Office Components / Forms / ServicesCounselingShipmentForm / AccountingCodes',
  component: AccountingCodes,
  decorators: [
    (Story) => (
      <GridContainer className={styles.gridContainer}>
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <Formik initialValues={{}}>
              {() => {
                return (
                  <Form className={formStyles.form} style={{ maxWidth: 'none' }}>
                    <Story />
                  </Form>
                );
              }}
            </Formik>
          </Grid>
        </Grid>
      </GridContainer>
    ),
  ],
};

export const AsRequired = () => <AccountingCodes optional={false} />;

export const WithNoTACsOrSACs = () => <AccountingCodes />;
WithNoTACsOrSACs.storyName = 'With No TACs or SACs';

export const WithSingleCode = () => <AccountingCodes TACs={{ hhg: '1234 ' }} />;

export const WithMultipleCodes = () => (
  <AccountingCodes TACs={{ hhg: '1234', nts: '5678' }} SACs={{ hhg: '98765', nts: '000012345' }} />
);
