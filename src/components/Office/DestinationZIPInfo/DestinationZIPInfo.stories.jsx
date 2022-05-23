import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';
import * as Yup from 'yup';
import { Formik } from 'formik';

import { ordersInfo } from '../RequestedShipments/RequestedShipmentsTestData';

import DestinationZIPInfo from './DestinationZIPInfo';

import { ZIP5_CODE_REGEX, InvalidZIPTypeError } from 'utils/validation';
import styles from 'pages/Office/ServicesCounselingMoveInfo/ServicesCounselingTab.module.scss';
import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';

const validationSchema = Yup.object().shape({
  destinationPostalCode: Yup.string().matches(ZIP5_CODE_REGEX, InvalidZIPTypeError).required('Required'),
  secondDestinationPostalCode: Yup.string().matches(ZIP5_CODE_REGEX, InvalidZIPTypeError),
});

export const DestinationZIPInfoExample = () => (
  <Formik
    initialValues={{
      destinationPostalCode: '',
      secondDestinationPostalCode: '',
    }}
    validationSchema={validationSchema}
  >
    {({ setFieldValue }) => {
      return (
        <Form className={formStyles.form} style={{ maxWidth: 'none' }}>
          <DestinationZIPInfo setFieldValue={setFieldValue} dutyZip="90210" />
        </Form>
      );
    }}
  </Formik>
);

export const DestinationZIPInfoWithDataExample = () => (
  <Formik
    initialValues={{
      destinationPostalCode: '08540',
      secondDestinationPostalCode: '07003',
    }}
    validationSchema={validationSchema}
  >
    {({ setFieldValue }) => {
      return (
        <Form className={formStyles.form} style={{ maxWidth: 'none' }}>
          <DestinationZIPInfo setFieldValue={setFieldValue} dutyZip="90210" />
        </Form>
      );
    }}
  </Formik>
);

export default {
  title: 'Office Components / Forms / ShipmentForm / Destination ZIP Info',
  components: ordersInfo,
  decorators: [
    (Story) => (
      <div className="officeApp">
        <GridContainer className={styles.gridContainer}>
          <Grid row>
            <Grid col desktop={{ col: 8, offset: 2 }}>
              <Story />
            </Grid>
          </Grid>
        </GridContainer>
      </div>
    ),
  ],
};
