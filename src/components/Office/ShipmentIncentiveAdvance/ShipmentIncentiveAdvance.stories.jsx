import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';
import { Formik } from 'formik';

import ShipmentIncentiveAdvance from './ShipmentIncentiveAdvance';

import styles from 'pages/Office/ServicesCounselingMoveInfo/ServicesCounselingTab.module.scss';
import shipmentFormStyles from 'components/Office/ShipmentForm/ShipmentForm.module.scss';
import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';

export default {
  title: 'Office Components / Forms / ShipmentForm / ShipmentIncentiveAdvance',
  component: ShipmentIncentiveAdvance,
  decorators: [
    (Story) => (
      <GridContainer className={styles.gridContainer}>
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <div className={shipmentFormStyles.ShipmentForm}>
              <Story />
            </div>
          </Grid>
        </Grid>
      </GridContainer>
    ),
  ],
};

export const advanceNotRequested = () => (
  <Formik initialValues={{ advanceRequested: false }}>
    {() => {
      return (
        <Form className={formStyles.form} style={{ maxWidth: 'none' }}>
          <ShipmentIncentiveAdvance />
        </Form>
      );
    }}
  </Formik>
);

export const advanceRequested = () => (
  <Formik initialValues={{ advanceRequested: true, amountRequested: '5000' }}>
    {() => {
      return (
        <Form className={formStyles.form} style={{ maxWidth: 'none' }}>
          <ShipmentIncentiveAdvance estimatedIncentive={1000000} />
        </Form>
      );
    }}
  </Formik>
);

export const advanceRequestedWithError = () => (
  <Formik initialValues={{ advanceRequested: true, amountRequested: '7000' }}>
    {() => {
      return (
        <Form className={formStyles.form} style={{ maxWidth: 'none' }}>
          <ShipmentIncentiveAdvance estimatedIncentive={1111111} />
        </Form>
      );
    }}
  </Formik>
);
