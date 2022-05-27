import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';
import { Formik } from 'formik';
import * as Yup from 'yup';

import ShipmentIncentiveAdvance from './ShipmentIncentiveAdvance';

import styles from 'pages/Office/ServicesCounselingMoveInfo/ServicesCounselingTab.module.scss';
import shipmentFormStyles from 'components/Office/ShipmentForm/ShipmentForm.module.scss';
import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';
import { getFormattedMaxAdvancePercentage } from 'utils/incentives';

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
  <Formik initialValues={{ hasRequestedAdvance: false }}>
    {() => {
      return (
        <Form className={formStyles.form} style={{ maxWidth: 'none' }}>
          <ShipmentIncentiveAdvance />
        </Form>
      );
    }}
  </Formik>
);

export const advanceRequested = () => {
  const estimatedIncentive = 1000000;

  const validationSchema = Yup.object().shape({
    advance: Yup.number().max(
      (estimatedIncentive * 0.6) / 100,
      `Enter an amount that is less than or equal to the maximum advance (${getFormattedMaxAdvancePercentage()} of estimated incentive)`,
    ),
  });

  return (
    <Formik validationSchema={validationSchema} initialValues={{ advanceRequested: true, advance: '5000' }}>
      {() => {
        return (
          <Form className={formStyles.form} style={{ maxWidth: 'none' }}>
            <ShipmentIncentiveAdvance estimatedIncentive={estimatedIncentive} />
          </Form>
        );
      }}
    </Formik>
  );
};

export const advanceRequestedWithError = () => {
  const estimatedIncentive = 1111111;

  const validationSchema = Yup.object().shape({
    advance: Yup.number().max(
      (estimatedIncentive * 0.6) / 100,
      `Enter an amount that is less than or equal to the maximum advance (${getFormattedMaxAdvancePercentage()} of estimated incentive)`,
    ),
  });

  return (
    <Formik validationSchema={validationSchema} initialValues={{ advanceRequested: true, advance: '7000' }}>
      {() => {
        return (
          <Form className={formStyles.form} style={{ maxWidth: 'none' }}>
            <ShipmentIncentiveAdvance estimatedIncentive={estimatedIncentive} />
          </Form>
        );
      }}
    </Formik>
  );
};
