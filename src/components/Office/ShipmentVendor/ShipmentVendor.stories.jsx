import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';
import { Formik } from 'formik';

import ShipmentVendor from './ShipmentVendor';

import styles from 'pages/Office/ServicesCounselingMoveInfo/ServicesCounselingTab.module.scss';
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

export const standard = () => <ShipmentVendor />;

export const externalVendorChecked = () => <ShipmentVendor usesExternalVendor />;
