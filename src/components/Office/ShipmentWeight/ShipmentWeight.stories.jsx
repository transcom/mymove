import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';
import { Formik } from 'formik';

import ShipmentWeight from './ShipmentWeight';

import styles from 'pages/Office/ServicesCounselingMoveInfo/ServicesCounselingTab.module.scss';
import shipmentFormStyles from 'components/Office/ShipmentForm/ShipmentForm.module.scss';
import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';

export default {
  title: 'Office Components / Forms / ShipmentForm / ShipmentWeight',
  component: ShipmentWeight,
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

export const withoutProGear = () => (
  <Formik initialValues={{ hasProGear: false }}>
    {() => {
      return (
        <Form className={formStyles.form} style={{ maxWidth: 'none' }}>
          <ShipmentWeight />
        </Form>
      );
    }}
  </Formik>
);

export const withProGear = () => (
  <Formik
    initialValues={{
      hasProGear: true,
    }}
  >
    {() => {
      return (
        <Form className={formStyles.form} style={{ maxWidth: 'none' }}>
          <ShipmentWeight />
        </Form>
      );
    }}
  </Formik>
);

export const withProGearAndData = () => (
  <Formik
    initialValues={{
      hasProGear: true,
      estimatedWeight: '4000',
      proGearWeight: '3000',
      spouseProGearWeight: '2000',
    }}
  >
    {() => {
      return (
        <Form className={formStyles.form} style={{ maxWidth: 'none' }}>
          <ShipmentWeight />
        </Form>
      );
    }}
  </Formik>
);
