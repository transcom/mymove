import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';
import { Formik } from 'formik';
import * as Yup from 'yup';

import ShipmentCustomerSIT from './ShipmentCustomerSIT';

import styles from 'pages/Office/ServicesCounselingMoveInfo/ServicesCounselingTab.module.scss';
import shipmentFormStyles from 'components/Office/ShipmentForm/ShipmentForm.module.scss';
import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';

const validationSchema = Yup.object().shape({
  sitExpected: Yup.boolean().required('Required'),
  sitEstimatedWeight: Yup.number().when('sitExpected', {
    is: true,
    then: (schema) => schema.required('Required'),
  }),
  sitEstimatedEntryDate: Yup.date().when('sitExpected', {
    is: true,
    then: (schema) =>
      schema.typeError('Enter a complete date in DD MMM YYYY format (day, month, year).').required('Required'),
  }),
  sitEstimatedDepartureDate: Yup.date().when('sitExpected', {
    is: true,
    then: (schema) =>
      schema.typeError('Enter a complete date in DD MMM YYYY format (day, month, year).').required('Required'),
  }),
});

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
  <Formik validationSchema={validationSchema} initialValues={{ sitExpected: false }}>
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
    validationSchema={validationSchema}
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
    validationSchema={validationSchema}
    initialValues={{
      sitExpected: true,
      sitLocation: 'DESTINATION',
      sitEstimatedWeight: '5725',
      sitEstimatedEntryDate: '2022-08-05',
      sitEstimatedDepartureDate: '2022-09-07',
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
