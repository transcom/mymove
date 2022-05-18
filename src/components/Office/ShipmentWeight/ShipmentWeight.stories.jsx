import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';
import * as Yup from 'yup';
import { Formik } from 'formik';

import ShipmentWeight from './ShipmentWeight';

import styles from 'pages/Office/ServicesCounselingMoveInfo/ServicesCounselingTab.module.scss';
import shipmentFormStyles from 'components/Office/ShipmentForm/ShipmentForm.module.scss';
import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';

const validationSchema = Yup.object().shape({
  estimatedWeight: Yup.number().min(1, 'Enter a weight greater than 0 lbs').required('Required'),
  hasProGear: Yup.boolean().required('Required'),
  proGearWeight: Yup.number()
    .min(0, 'Enter a weight 0 lbs or greater')
    .when(['hasProGear', 'spouseProGearWeight'], {
      is: (hasProGear, spouseProGearWeight) => hasProGear && !spouseProGearWeight,
      then: (schema) =>
        schema
          .required(
            `Enter weight in at least one pro-gear field. If the customer will not move pro-gear in this PPM, select No above.`,
          )
          .max(2000, 'Enter a weight 2,000 lbs or less'),
      otherwise: Yup.number().min(0, 'Enter a weight 0 lbs or greater').max(2000, 'Enter a weight 2,000 lbs or less'),
    }),
  spouseProGearWeight: Yup.number()
    .min(0, 'Enter a weight 0 lbs or greater')
    .max(2000, 'Enter a weight 2,000 lbs or less'),
});

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
  <Formik
    initialValues={{ hasProGear: false, estimatedWeight: '', proGearWeight: '', spouseProGearWeight: '' }}
    validationSchema={validationSchema}
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

export const withProGear = () => (
  <Formik
    initialValues={{
      hasProGear: true,
      estimatedWeight: '',
      proGearWeight: '',
      spouseProGearWeight: '',
    }}
    validationSchema={validationSchema}
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
      estimatedWeight: '1500',
      proGearWeight: '1000',
      spouseProGearWeight: '300',
    }}
    validationSchema={validationSchema}
  >
    {() => {
      return (
        <Form className={formStyles.form} style={{ maxWidth: 'none' }}>
          <ShipmentWeight authorizedWeight="2000" />
        </Form>
      );
    }}
  </Formik>
);

export const withAuthorizedWeightWarning = () => (
  <Formik
    initialValues={{
      hasProGear: true,
      estimatedWeight: '4000',
      proGearWeight: '',
      spouseProGearWeight: '',
    }}
    validationSchema={validationSchema}
  >
    {() => {
      return (
        <Form className={formStyles.form} style={{ maxWidth: 'none' }}>
          <ShipmentWeight authorizedWeight="2000" />
        </Form>
      );
    }}
  </Formik>
);

export const withProGearWeightOverError = () => (
  <Formik
    initialValues={{
      hasProGear: true,
      estimatedWeight: '',
      proGearWeight: '3000',
      spouseProGearWeight: '3000',
    }}
    validationSchema={validationSchema}
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
