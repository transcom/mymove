import React, { useState } from 'react';
import { connect } from 'react-redux';
import { func } from 'prop-types';
import { useNavigate } from 'react-router-dom';
import { Grid, GridContainer, Alert } from '@trussworks/react-uswds';

import { setFlashMessage } from 'store/flash/actions';
import RequestAccountForm from 'components/Office/RequestAccountForm/RequestAccountForm';
import { createOfficeAccountRequest } from 'services/ghcApi';
import NotificationScrollToTop from 'components/NotificationScrollToTop';
import { generalRoutes } from 'constants/routes';

export const RequestAccount = (props) => {
  const navigate = useNavigate();
  const [serverError, setServerError] = useState(null);
  const { setMessage } = props;

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

    if (values.transportationInvoicingOfficerCheckBox) {
      requestedRoles.push({
        name: 'Transportation Ordering Officer',
        roleType: 'transportation_ordering_officer',
      });
    }
    if (values.transportationOrderingOfficerCheckBox) {
      requestedRoles.push({
        name: 'Transportation Invoicing Officer',
        roleType: 'transportation_invoicing_officer',
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
    if (values.qualityAssuranceAndCustomerSupportCheckBox) {
      requestedRoles.push({
        name: 'Quality Assurance and Customer Service',
        roleType: 'qae_csr',
      });
    }

    const body = {
      email: values.officeAccountRequestEmail,
      edipi: values.officeAccountRequestEdipi,
      otherUniqueId: values.officeAccountRequestOtherUniqueId,
      firstName: values.officeAccountRequestFirstName,
      middleInitials: values.officeAccountRequestMiddleInitial,
      lastName: values.officeAccountRequestLastName,
      telephone: values.officeAccountRequestTelephone,
      transportationOfficeId: values.officeAccountTransportationOffice.id,
      roles: requestedRoles,
    };

    return createOfficeAccountRequest({ body })
      .then(() => {
        setMessage(
          'OFFICE_ACCOUNT_REQUEST_SUCCESS',
          'success',
          'Request Office Account form successfully submitted.',
          '',
          true,
        );
        navigate(generalRoutes.SIGN_IN_PATH);
      })
      .catch(() => {
        const errorMessage = 'Failed to submit office account request due to server error';
        setServerError(errorMessage);
      });
  };

  return (
    <GridContainer>
      <NotificationScrollToTop dependency={serverError} />

      {serverError && (
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <Alert data-testid="alert2" type="error" headingLevel="h4" heading="An error occurred">
              {serverError}
            </Alert>
          </Grid>
        </Grid>
      )}

      <Grid row>
        <Grid col desktop={{ col: 8, offset: 2 }}>
          <RequestAccountForm onCancel={handleCancel} onSubmit={handleSubmit} initialValues={initialValues} />
        </Grid>
      </Grid>
    </GridContainer>
  );
};

RequestAccount.propTypes = {
  setMessage: func.isRequired,
};

const mapDispatchToProps = {
  setMessage: setFlashMessage,
};

// export default RequestAccount;
export default connect(() => ({}), mapDispatchToProps)(RequestAccount);
