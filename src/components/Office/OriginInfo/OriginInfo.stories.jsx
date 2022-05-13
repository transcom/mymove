import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';
import { Formik } from 'formik';

import { ordersInfo } from '../RequestedShipments/RequestedShipmentsTestData';

import OriginInfo from './OriginInfo';

import styles from 'pages/Office/ServicesCounselingMoveInfo/ServicesCounselingTab.module.scss';
import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';

export const OriginInfoExample = () => (
  <Formik
    initialValues={{
      plannedDepartureDate: '',
      originPostalCode: '',
      useResidentialAddressZIP: false,
      secondOriginPostalCode: '',
    }}
  >
    {({ setFieldValue, values }) => {
      return (
        <Form className={formStyles.form} style={{ maxWidth: 'none' }}>
          <OriginInfo
            setFieldValue={setFieldValue}
            currentZip="90210"
            isUseResidentialAddressZIPChecked={values.useResidentialAddressZIP}
            postalCodeValidator={() => {}}
          />
        </Form>
      );
    }}
  </Formik>
);

export const OriginInfoExampleWithZipValidator = () => (
  <Formik
    initialValues={{
      plannedDepartureDate: '',
      originPostalCode: '',
      useResidentialAddressZIP: false,
      secondOriginPostalCode: '',
    }}
  >
    {({ setFieldValue, values }) => {
      return (
        <Form className={formStyles.form} style={{ maxWidth: 'none' }}>
          <OriginInfo
            setFieldValue={setFieldValue}
            currentZip="90210"
            isUseResidentialAddressZIPChecked={values.useResidentialAddressZIP}
            postalCodeValidator={() => 'We do not support that ZIP code.'}
          />
        </Form>
      );
    }}
  </Formik>
);

export default {
  title: 'Office Components / Origin Info',
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
