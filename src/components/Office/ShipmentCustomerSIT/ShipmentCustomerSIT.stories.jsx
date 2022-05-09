import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';
import { Formik } from 'formik';

import ShipmentCustomerSIT from './ShipmentCustomerSIT';

import styles from 'pages/Office/ServicesCounselingMoveInfo/ServicesCounselingTab.module.scss';
import shipmentFormStyles from 'components/Office/ShipmentForm/ShipmentForm.module.scss';
import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';

export default {
  title: 'Office Components / Forms / ShipmentForm / StorageInTransit',
  component: ShipmentCustomerSIT,
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

export const withoutSIT = () => (
  <Formik initialValues={{ sitExpected: false }}>
    {() => {
      return (
        <Form className={formStyles.form} style={{ maxWidth: 'none' }}>
          <ShipmentCustomerSIT />
        </Form>
      );
    }}
  </Formik>
);

export const withSIT = () => (
  <Formik
    initialValues={{
      sitExpected: true,
    }}
  >
    {() => {
      return (
        <Form className={formStyles.form} style={{ maxWidth: 'none' }}>
          <ShipmentCustomerSIT />
        </Form>
      );
    }}
  </Formik>
);
withSIT.storyName = 'With SIT';

export const withSITAndData = () => (
  <Formik
    initialValues={{
      sitExpected: true,
      sitLocation: 'destination',
      estimatedSITWeight: '5725',
      estimatedSITStart: '08/05/2022',
      estimatedSITEnd: '09/07/2022',
    }}
  >
    {() => {
      return (
        <Form className={formStyles.form} style={{ maxWidth: 'none' }}>
          <ShipmentCustomerSIT />
        </Form>
      );
    }}
  </Formik>
);
