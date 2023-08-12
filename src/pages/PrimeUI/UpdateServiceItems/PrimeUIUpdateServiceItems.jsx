import React, { useState } from 'react';
import { Alert, Grid, GridContainer } from '@trussworks/react-uswds';
import * as Yup from 'yup';

import PrimeUIRequestSITDestAddressChangeForm from './PrimeUIRequestSITDestAddressChangeForm';

import styles from 'components/Office/CustomerContactInfoForm/CustomerContactInfoForm.module.scss';
import primeStyles from 'pages/PrimeUI/Prime.module.scss';
import { addressSchema } from 'utils/validation';

const UpdateServiceItems = () => {
  const [errorMessage, setErrorMessage] = useState();
  setErrorMessage('just a placeholder so i can commit');

  const destAddressChangeRequestSchema = Yup.object().shape({
    addressID: Yup.string(),
    destinationAddress: Yup.object().shape({
      address: addressSchema,
    }),
    eTag: Yup.string(),
  });

  return (
    <div className={styles.tabContent}>
      <div className={styles.container}>
        <GridContainer className={styles.gridContainer}>
          <Grid row>
            <Grid col desktop={{ col: 9, offset: 2 }}>
              {errorMessage?.detail && (
                <div className={primeStyles.errorContainer}>
                  <Alert headingLevel="h4" type="error">
                    <span className={primeStyles.errorTitle}>{errorMessage.title}</span>
                    <span className={primeStyles.errorDetail}>{errorMessage.detail}</span>
                  </Alert>
                </div>
              )}
              <h1 className={styles.sectionHeader}>Update Service Items</h1>
              <PrimeUIRequestSITDestAddressChangeForm
                name="destinationAddress.address"
                destAddressChangeRequestSchema={destAddressChangeRequestSchema}
              />
            </Grid>
          </Grid>
        </GridContainer>
      </div>
    </div>
  );
};

export default UpdateServiceItems;
