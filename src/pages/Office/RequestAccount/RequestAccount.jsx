import React, { useState } from 'react';
import { connect } from 'react-redux';
import { func } from 'prop-types';
import { useNavigate } from 'react-router-dom';
import { Grid, GridContainer, Alert } from '@trussworks/react-uswds';

import styles from './RequestAccount.module.scss';

import { setFlashMessage as setFlashMessageAction } from 'store/flash/actions';
import RequestAccountForm from 'components/Office/RequestAccountForm/RequestAccountForm';
import { createOfficeAccountRequest } from 'services/ghcApi';
import NotificationScrollToTop from 'components/NotificationScrollToTop';
import { generalRoutes } from 'constants/routes';

export const RequestAccount = ({ setFlashMessage }) => {
  const navigate = useNavigate();
  const [serverError, setServerError] = useState(null);

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
    const requestedRoles = [];

    if (values.taskInvoicingOfficerCheckBox) {
      requestedRoles.push({
        name: 'Task Invoicing Officer',
        roleType: 'task_invoicing_officer',
      });
    }
    if (values.taskOrderingOfficerCheckBox) {
      requestedRoles.push({
        name: 'Task Ordering Officer',
        roleType: 'task_ordering_officer',
      });
    }
    if (values.headquartersCheckBox) {
      requestedRoles.push({
        name: 'Headquarters',
        roleType: 'headquarters',
      });
    }
    if (values.transportationContractingOfficerCheckBox) {
      requestedRoles.push({
        name: 'Contracting Officer',
        roleType: 'contracting_officer',
      });
    }
    if (values.servicesCounselorCheckBox) {
      requestedRoles.push({
        name: 'Services Counselor',
        roleType: 'services_counselor',
      });
    }
    if (values.qualityAssuranceEvaluatorCheckBox) {
      requestedRoles.push({
        name: 'Quality Assurance Evaluator',
        roleType: 'qae',
      });
    }
    if (values.customerSupportRepresentativeCheckBox) {
      requestedRoles.push({
        name: 'Customer Service Representative',
        roleType: 'customer_service_representative',
      });
    }
    if (values.governmentSurveillanceRepresentativeCheckbox) {
      requestedRoles.push({
        name: 'Government Surveillance Representative',
        roleType: 'gsr',
      });
    }

    let body = {
      email: values.officeAccountRequestEmail,
      firstName: values.officeAccountRequestFirstName,
      middleInitials: values.officeAccountRequestMiddleInitial,
      lastName: values.officeAccountRequestLastName,
      telephone: values.officeAccountRequestTelephone,
      transportationOfficeId: values.officeAccountTransportationOffice.id,
      roles: requestedRoles,
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
          'Request Office Account form successfully submitted.',
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
          <RequestAccountForm onCancel={handleCancel} onSubmit={handleSubmit} initialValues={initialValues} />
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
