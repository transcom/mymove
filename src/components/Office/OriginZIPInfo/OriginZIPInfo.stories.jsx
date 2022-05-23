import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';
import * as Yup from 'yup';
import { Formik } from 'formik';

import { ordersInfo } from '../RequestedShipments/RequestedShipmentsTestData';

import OriginZIPInfo from './OriginZIPInfo';

import { ZIP5_CODE_REGEX, InvalidZIPTypeError } from 'utils/validation';
import styles from 'pages/Office/ServicesCounselingMoveInfo/ServicesCounselingTab.module.scss';
import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';

const validationSchema = Yup.object().shape({
  pickupPostalCode: Yup.string().matches(ZIP5_CODE_REGEX, InvalidZIPTypeError).required('Required'),
  secondPickupPostalCode: Yup.string().matches(ZIP5_CODE_REGEX, InvalidZIPTypeError),
  expectedDepartureDate: Yup.date()
    .typeError('Enter a complete date in DD MMM YYYY format (day, month, year).')
    .required('Required'),
});

export const OriginZIPInfoExample = () => (
  <Formik
    initialValues={{
      expectedDepartureDate: '',
      pickupPostalCode: '',
      secondPickupPostalCode: '',
    }}
    validationSchema={validationSchema}
  >
    {({ setFieldValue }) => {
      return (
        <Form className={formStyles.form} style={{ maxWidth: 'none' }}>
          <OriginZIPInfo setFieldValue={setFieldValue} currentZip="90210" />
        </Form>
      );
    }}
  </Formik>
);

export const OriginZIPInfoWithDataExample = () => (
  <Formik
    initialValues={{
      expectedDepartureDate: '2022-09-23',
      pickupPostalCode: '07003',
      secondPickupPostalCode: '08540',
    }}
  >
    {({ setFieldValue }) => {
      return (
        <Form className={formStyles.form} style={{ maxWidth: 'none' }}>
          <OriginZIPInfo setFieldValue={setFieldValue} currentZip="90210" />
        </Form>
      );
    }}
  </Formik>
);

export default {
  title: 'Office Components  / Forms / ShipmentForm / Origin ZIP Info',
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
