import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';
import { Formik } from 'formik';

import ShipmentVendor from './ShipmentVendor';

import styles from 'pages/Office/ServicesCounselingMoveInfo/ServicesCounselingTab.module.scss';
import shipmentFormStyles from 'components/Office/ShipmentForm/ShipmentForm.module.scss';
import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';

export default {
  title: 'Office Components / Forms / ShipmentForm / ShipmentVendor',
  component: ShipmentVendor,
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

export const standard = () => (
  <Formik initialValues={{ usesExternalVendor: false }}>
    {() => {
      return (
        <Form className={formStyles.form} style={{ maxWidth: 'none' }}>
          <ShipmentVendor />
        </Form>
      );
    }}
  </Formik>
);

export const externalVendorChecked = () => (
  <Formik initialValues={{ usesExternalVendor: true }}>
    {() => {
      return (
        <Form className={formStyles.form} style={{ maxWidth: 'none' }}>
          <ShipmentVendor />
        </Form>
      );
    }}
  </Formik>
);
