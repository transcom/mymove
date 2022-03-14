import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';
import { Formik } from 'formik';

import StorageFacilityInfo from './StorageFacilityInfo';

import styles from 'pages/Office/ServicesCounselingMoveInfo/ServicesCounselingTab.module.scss';
import shipmentFormStyles from 'components/Office/ShipmentForm/ShipmentForm.module.scss';
import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';
import { roleTypes } from 'constants/userRoles';

export default {
  title: 'Office Components / Forms / ShipmentForm / StorageFacilityInfo',
  component: StorageFacilityInfo,
  decorators: [
    (Story) => (
      <GridContainer className={styles.gridContainer}>
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <Formik initialValues={{}}>
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

export const AsServiceCounselor = () => <StorageFacilityInfo userRole={roleTypes.SERVICES_COUNSELOR} />;
export const AsTOO = () => <StorageFacilityInfo userRole={roleTypes.TOO} />;
