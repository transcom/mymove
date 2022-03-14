import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';
import { Formik } from 'formik';

import ShipmentWeightInput from './ShipmentWeightInput';

import styles from 'pages/Office/ServicesCounselingMoveInfo/ServicesCounselingTab.module.scss';
import shipmentFormStyles from 'components/Office/ShipmentForm/ShipmentForm.module.scss';
import { Form } from 'components/form/Form';
import { roleTypes } from 'constants/userRoles';
import formStyles from 'styles/form.module.scss';

export default {
  title: 'Office Components / Forms / ShipmentForm / ShipmentWeightInput',
  component: ShipmentWeightInput,
  argTypes: {
    userRole: {
      options: [roleTypes.SERVICES_COUNSELOR, roleTypes.TOO],
      control: { type: 'radio' },
    },
  },
  decorators: [
    (Story) => (
      <GridContainer className={styles.gridContainer}>
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <Formik initialValues={{ ntsRecordedWeight: '' }}>
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

const Template = (args) => <ShipmentWeightInput {...args} />;

export const Standard = Template.bind({});
