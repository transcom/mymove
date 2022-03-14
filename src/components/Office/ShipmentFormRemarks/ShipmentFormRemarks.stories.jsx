import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';
import { Formik } from 'formik';

import ShipmentFormRemarks from './ShipmentFormRemarks';

import styles from 'pages/Office/ServicesCounselingMoveInfo/ServicesCounselingTab.module.scss';
import shipmentFormStyles from 'components/Office/ShipmentForm/ShipmentForm.module.scss';
import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';
import { roleTypes } from 'constants/userRoles';

export default {
  title: 'Office Components / Forms / ShipmentForm / ShipmentFormRemarks',
  component: ShipmentFormRemarks,
  decorators: [
    (Story) => (
      <GridContainer className={styles.gridContainer}>
        <Grid row>
          <Grid col={12}>
            <Formik initialValues={{ counselorRemarks: 'mock counselor remarks' }}>
              {() => {
                return (
                  <Form className={formStyles.form} style={{ maxWidth: 'none' }}>
                    <div className={shipmentFormStyles.ShipmentForm}>
                      <Story />
                    </div>
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

export const AsServiceCounselor = () => (
  <ShipmentFormRemarks userRole={roleTypes.SERVICES_COUNSELOR} customerRemarks="mock customer remarks" />
);
export const AsTOO = () => <ShipmentFormRemarks userRole={roleTypes.TOO} />;
