import React from 'react';
import { Grid, GridContainer, Alert } from '@trussworks/react-uswds';

import styles from './RequestAccount.module.scss';

import NotificationScrollToTop from 'components/NotificationScrollToTop';
import RequestAccountForm from 'components/Office/RequestAccountForm/RequestAccountForm';

const RequestAccountView = ({ serverError, onCancel, onSubmit, initialValues, rolesWithPrivs, privileges }) => (
  <GridContainer>
    <NotificationScrollToTop dependency={serverError} />
    {serverError && (
      <Grid row>
        <Alert data-testid="alert2" type="error" headingLevel="h4" heading="An error occurred" className={styles.error}>
          {serverError}
        </Alert>
      </Grid>
    )}
    <Grid row>
      <Grid col desktop={{ col: 8 }} className={styles.formContainer}>
        <RequestAccountForm
          onCancel={onCancel}
          onSubmit={onSubmit}
          initialValues={initialValues}
          rolesWithPrivs={rolesWithPrivs}
          privileges={privileges}
        />
      </Grid>
    </Grid>
  </GridContainer>
);

export default RequestAccountView;
