import PropTypes from 'prop-types';
import React, { useState } from 'react';
import { GridContainer, Grid, Alert } from '@trussworks/react-uswds';
import { useNavigate } from 'react-router-dom';

import styles from './CustomerName.module.scss';

import NotificationScrollToTop from 'components/NotificationScrollToTop';
import NameForm from 'components/Customer/NameForm/NameForm';
import { patchServiceMember, getResponseError } from 'services/internalApi';
import { ServiceMemberShape } from 'types/customerShapes';
import { generalRoutes } from 'constants/routes';

export const CustomerName = ({ serviceMember, updateServiceMember }) => {
  const [serverError, setServerError] = useState(null);
  const navigate = useNavigate();
  const initialValues = {
    first_name: serviceMember?.first_name || '',
    middle_name: serviceMember?.middle_name || '',
    last_name: serviceMember?.last_name || '',
    suffix: serviceMember?.suffix || '',
  };

  const handleNext = () => {
    // add next route
  };

  const handleBack = () => {
    navigate(generalRoutes.BASE_QUEUE_SEARCH_PATH);
  };

  const handleSubmit = (values) => {
    const payload = {
      id: serviceMember.id,
      first_name: values.first_name,
      middle_name: values.middle_name,
      last_name: values.last_name,
      suffix: values.suffix,
    };

    return patchServiceMember(payload)
      .then(updateServiceMember)
      .then(handleNext)
      .catch((e) => {
        // TODO - error handling - below is rudimentary error handling to approximate existing UX
        // Error shape: https://github.com/swagger-api/swagger-js/blob/master/docs/usage/http-client.md#errors
        const { response } = e;
        const errorMessage = getResponseError(response, 'failed to update service member due to server error');
        setServerError(errorMessage);
      });
  };

  return (
    <GridContainer>
      <NotificationScrollToTop dependency={serverError} />

      {serverError && (
        <Grid>
          <Grid col desktop={{ col: 8, offset: 2 }} style={{ width: '750px' }}>
            <Alert type="error" headingLevel="h4" heading="An error occurred">
              {serverError}
            </Alert>
          </Grid>
        </Grid>
      )}

      <Grid className={styles.nameFormContainer}>
        <Grid col desktop={{ col: 8 }} className={styles.nameForm}>
          <NameForm onSubmit={handleSubmit} onBack={handleBack} initialValues={initialValues} />
        </Grid>
      </Grid>
    </GridContainer>
  );
};

CustomerName.propTypes = {
  updateServiceMember: PropTypes.func.isRequired,
  serviceMember: ServiceMemberShape.isRequired,
};

export default CustomerName;
