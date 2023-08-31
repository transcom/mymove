import PropTypes from 'prop-types';
import React, { useState } from 'react';
import { connect } from 'react-redux';
import { useNavigate } from 'react-router-dom';
import { Alert, Grid, GridContainer } from '@trussworks/react-uswds';

import EditOktaInfoForm from 'components/Customer/EditOktaInfoForm/EditOktaInfoForm';
import NotificationScrollToTop from 'components/NotificationScrollToTop';
import { customerRoutes } from 'constants/routes';
import { getResponseError, patchOktaProfile } from 'services/internalApi';
import { updateOktaProfile as updateOktaProfileAction } from 'store/entities/actions';
import { selectBackupContacts, selectServiceMemberFromLoggedInUser } from 'store/entities/selectors';
import { setFlashMessage as setFlashMessageAction } from 'store/flash/actions';

export const EditOktaInfo = ({ serviceMember, setFlashMessage }) => {
  const navigate = useNavigate();
  const [serverError, setServerError] = useState(null);

  const initialValues = {
    oktaUsername: serviceMember?.oktaUsername || 'Not Provided',
    oktaEmail: serviceMember?.oktaEmail || 'Not Provided',
    oktaFirstName: serviceMember?.oktaFirstName || 'Not Provided',
    oktaLastName: serviceMember?.oktaLastName || 'Not Provided',
    oktaEdipi: serviceMember?.oktaEdipi || 'Not Provided',
  };

  const handleCancel = () => {
    navigate(customerRoutes.PROFILE_PATH);
  };

  // sends Okta data in form to backend to call Okta API to update profile values
  // TODO need to redirect the user back to customerRoutes.PROFILE_PATH
  // TODO need to also update the users table with okta_email if it is different
  const handleSubmit = async (values) => {
    const oktaPayload = {
      id: serviceMember.id,
      username: values?.oktaUsername,
      email: values?.oktaEmail,
      firstName: values?.oktaFirstName,
      lastName: values?.oktaLastName,
      cac_edipi: values?.oktaEdipi,
    };

    //! leaving this here for reference when implementing API calls for Okta
    // return patchOktaProfile(oktaPayload)
    // .then(updateServiceMember)
    // .then(() => {
    //   setFlashMessage('EDIT_CONTACT_INFO_SUCCESS', 'success', "You've updated your information.");
    //   navigate(customerRoutes.PROFILE_PATH);
    // })
    // .catch((e) => {
    //   const { response } = e;
    //   const errorMessage = getResponseError(response, 'Failed to update service member due to server error');

    //   setServerError(errorMessage);
    // });

    return patchOktaProfile(oktaPayload)
      .then(() => {
        setFlashMessage('EDIT_OKTA_PROFILE_SUCCESS', 'success', "You've updated your Okta profile.");
      })
      .catch((e) => {
        const { response } = e;
        const errorMessage = getResponseError(response, 'Failed to update okta profile due to server error');

        setServerError(errorMessage);
      });
  };

  return (
    <GridContainer>
      <NotificationScrollToTop dependency={serverError} />

      {serverError && (
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <Alert type="error" headingLevel="h4" heading="An error occurred">
              {serverError}
            </Alert>
          </Grid>
        </Grid>
      )}

      <Grid row>
        <Grid col desktop={{ col: 8, offset: 2 }}>
          <EditOktaInfoForm initialValues={initialValues} onCancel={handleCancel} onSubmit={handleSubmit} />
        </Grid>
      </Grid>
    </GridContainer>
  );
};

EditOktaInfo.propTypes = {
  setFlashMessage: PropTypes.func.isRequired,
};

const mapDispatchToProps = {
  setFlashMessage: setFlashMessageAction,
  updateOktaProfile: updateOktaProfileAction,
};

const mapStateToProps = (state) => ({
  currentBackupContacts: selectBackupContacts(state),
  serviceMember: selectServiceMemberFromLoggedInUser(state),
});

export default connect(mapStateToProps, mapDispatchToProps)(EditOktaInfo);
