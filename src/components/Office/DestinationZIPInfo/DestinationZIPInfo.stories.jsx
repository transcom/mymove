import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';
import { Formik } from 'formik';

import { ordersInfo } from '../RequestedShipments/RequestedShipmentsTestData';

import DestinationZIPInfo from './DestinationZIPInfo';

import styles from 'pages/Office/ServicesCounselingMoveInfo/ServicesCounselingTab.module.scss';
import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';

export const DestinationZIPInfoExample = () => (
  <Formik
    initialValues={{
      destinationPostalCode: '',
      useDutyZIP: false,
      secondDestinationPostalCode: '',
    }}
  >
    {({ setFieldValue, values }) => {
      return (
        <Form className={formStyles.form} style={{ maxWidth: 'none' }}>
          <DestinationZIPInfo
            setFieldValue={setFieldValue}
            dutyZip="90210"
            isUseDutyZIPChecked={values.useDutyZIP}
            postalCodeValidator={() => {}}
          />
        </Form>
      );
    }}
  </Formik>
);

export const DestinationZIPInfoExampleWithZipValidator = () => (
  <Formik
    initialValues={{
      destinationPostalCode: '',
      useDutyZIP: false,
      secondDestinationPostalCode: '',
    }}
  >
    {({ setFieldValue, values }) => {
      return (
        <Form className={formStyles.form} style={{ maxWidth: 'none' }}>
          <DestinationZIPInfo
            setFieldValue={setFieldValue}
            dutyZip="90210"
            isUseDutyZIPChecked={values.useDutyZIP}
            postalCodeValidator={() => 'We do not support that ZIP code.'}
          />
        </Form>
      );
    }}
  </Formik>
);

export default {
  title: 'Office Components / Forms / Shipment Form / Destination ZIP Info',
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
