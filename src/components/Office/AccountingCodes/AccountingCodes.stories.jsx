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

export const AsRequired = () => (
  <div className="officeApp">
    <AccountingCodes optional={false} />{' '}
  </div>
);

export const WithNoTACsOrSACs = () => (
  <div className="officeApp">
    <AccountingCodes />
  </div>
);
WithNoTACsOrSACs.storyName = 'With No TACs or SACs';

export const WithSingleCode = () => (
  <div className="officeApp">
    <AccountingCodes TACs={{ HHG: '1234 ' }} />
  </div>
);

export const WithMultipleCodes = () => (
  <div className="officeApp">
    <AccountingCodes TACs={{ HHG: '1234', NTS: '5678' }} SACs={{ HHG: '98765', NTS: '000012345' }} />
  </div>
);
