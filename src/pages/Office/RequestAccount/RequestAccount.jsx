import React, { useEffect, useState } from 'react';
import { connect, useDispatch } from 'react-redux';
import { func } from 'prop-types';
import { useNavigate } from 'react-router-dom';
import { Grid, GridContainer, Alert } from '@trussworks/react-uswds';

import styles from './RequestAccount.module.scss';

import { setFlashMessage as setFlashMessageAction } from 'store/flash/actions';
import RequestAccountForm from 'components/Office/RequestAccountForm/RequestAccountForm';
import { createOfficeAccountRequest } from 'services/ghcApi';
import NotificationScrollToTop from 'components/NotificationScrollToTop';
import { generalRoutes } from 'constants/routes';
import { useRolesPrivilegesQueriesOfficeApp } from 'hooks/queries';
import { setShowLoadingSpinner } from 'store/general/actions';

export const RequestAccount = ({ setFlashMessage }) => {
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const [serverError, setServerError] = useState(null);

  const { result, isLoading } = useRolesPrivilegesQueriesOfficeApp();

  useEffect(() => {
    if (isLoading) {
      dispatch(setShowLoadingSpinner(true, null));
    } else {
      dispatch(setShowLoadingSpinner(false, null));
    }
  }, [isLoading, dispatch]);

  const initialValues = {
    officeAccountRequestFirstName: '',
    officeAccountRequestMiddleInitial: '',
    officeAccountRequestLastName: '',
    officeAccountRequestEmail: '',
    officeAccountRequestTelephone: '',
    officeAccountRequestEdipi: '',
    officeAccountRequestOtherUniqueId: '',
    officeAccountTransportationOffice: undefined,
  };

  const handleCancel = () => {
    navigate(generalRoutes.SIGN_IN_PATH);
  };

  const handleSubmit = async (values) => {
    // Dynamically build requestedRoles and requestedPrivileges
    const requestedRoles = (result.rolesWithPrivs || [])
      .filter((role) => values[`${role.roleType}Checkbox`])
      .map((role) => ({
        name: role.roleName,
        roleType: role.roleType,
      }));

    const requestedPrivileges = (result.privileges || [])
      .filter((priv) => values[`${priv.privilegeType}PrivilegeCheckbox`])
      .map((priv) => ({
        name: priv.privilegeName,
        privilegeType: priv.privilegeType,
      }));

    let body = {
      email: values.officeAccountRequestEmail,
      firstName: values.officeAccountRequestFirstName,
      middleInitials: values.officeAccountRequestMiddleInitial,
      lastName: values.officeAccountRequestLastName,
      telephone: values.officeAccountRequestTelephone,
      transportationOfficeId: values.officeAccountTransportationOffice.id,
      roles: requestedRoles,
      privileges: requestedPrivileges,
    };

    if (values.officeAccountRequestEdipi !== '') {
      body = {
        ...body,
        edipi: values.officeAccountRequestEdipi,
      };
    }

    if (values.officeAccountRequestOtherUniqueId !== '') {
      body = {
        ...body,
        otherUniqueId: values.officeAccountRequestOtherUniqueId,
      };
    }

    return createOfficeAccountRequest({ body })
      .then(() => {
        setFlashMessage(
          'OFFICE_ACCOUNT_REQUEST_SUCCESS',
          'success',
          'You have successfully requested access to MilMove. This request must be processed by an administrator prior to login. Once this process is completed, an approval or rejection email will be sent notifying you of the status of your account request.',
          '',
          true,
        );
        navigate(generalRoutes.SIGN_IN_PATH);
      })
      .catch((e) => {
        const { response } = e;
        let errorMessage = `Failed to submit office account request.`;

        if (response.body) {
          const responseBody = response.body;
          let responseMsg = '';

          if (responseBody.detail) {
            responseMsg += `${responseBody.detail}:`;
          }

          if (responseBody.invalid_fields) {
            const invalidFields = responseBody.invalid_fields;
            Object.keys(invalidFields).forEach((key) => {
              responseMsg += `\n${invalidFields[key]}`;
            });
          }
          errorMessage += `\n${responseMsg}`;
        }

        setServerError(errorMessage);
      });
  };

  if (isLoading || !result.rolesWithPrivs || !result.privileges) {
    return null;
  }

  return (
    <GridContainer>
      <NotificationScrollToTop dependency={serverError} />

      {serverError && (
        <Grid row>
          <Alert
            data-testid="alert2"
            type="error"
            headingLevel="h4"
            heading="An error occurred"
            className={styles.error}
          >
            {serverError}
          </Alert>
        </Grid>
      )}
      <Grid row>
        <Grid col desktop={{ col: 8 }} className={styles.formContainer}>
          <RequestAccountForm
            onCancel={handleCancel}
            onSubmit={handleSubmit}
            Add
            commentMore
            actions
            initialValues={initialValues}
            rolesWithPrivs={result.rolesWithPrivs}
            privileges={result.privileges}
          />
        </Grid>
      </Grid>
    </GridContainer>
  );
};

RequestAccount.propTypes = {
  setFlashMessage: func.isRequired,
};

const mapDispatchToProps = {
  setFlashMessage: setFlashMessageAction,
};

export default connect(() => ({}), mapDispatchToProps)(RequestAccount);
