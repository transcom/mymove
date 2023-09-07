import PropTypes from 'prop-types';
import React, { useState } from 'react';
import { connect } from 'react-redux';
import { useNavigate } from 'react-router-dom';
import { Alert, Grid, GridContainer } from '@trussworks/react-uswds';

import { OktaUserInfoShape } from 'types/user';
import EditOktaInfoForm from 'components/Customer/EditOktaInfoForm/EditOktaInfoForm';
import NotificationScrollToTop from 'components/NotificationScrollToTop';
import { customerRoutes } from 'constants/routes';
import { getResponseError, updateOktaUser } from 'services/internalApi';
import { selectServiceMemberFromLoggedInUser, selectOktaUser } from 'store/entities/selectors';
import { setFlashMessage as setFlashMessageAction } from 'store/flash/actions';

export const EditOktaInfo = ({ serviceMember, setFlashMessage, oktaUser }) => {
  const navigate = useNavigate();
  const [serverError, setServerError] = useState(null);

  const initialValues = {
    oktaUsername: oktaUser?.username || 'Not Provided',
    oktaEmail: oktaUser?.email || 'Not Provided',
    oktaFirstName: oktaUser?.first_name || 'Not Provided',
    oktaLastName: oktaUser?.last_name || 'Not Provided',
    oktaEdipi: oktaUser?.edipi || '',
    oktaSub: oktaUser?.sub,
  };

  const handleCancel = () => {
    navigate(customerRoutes.PROFILE_PATH);
  };

  // sends Okta data in form to backend to call Okta API to update profile values
  // TODO need to also update the users table with okta_email if it is different
  const handleSubmit = async (values) => {
    // including serviceMember.id in case we need to udpate users table with new okta_email
    const oktaPayload = {
      profile: {
        id: serviceMember.id,
        username: values?.oktaUsername,
        email: values?.oktaEmail,
        first_name: values?.oktaFirstName,
        last_name: values?.oktaLastName,
        cac_edipi: values?.oktaEdipi,
        sub: values?.oktaSub,
      },
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

    return updateOktaUser(oktaPayload)
      .then(() => {
        setFlashMessage('EDIT_OKTA_PROFILE_SUCCESS', 'success', "You've updated your Okta profile.");
        // navigate(customerRoutes.PROFILE_PATH);
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
  oktaUser: OktaUserInfoShape.isRequired,
};

const mapDispatchToProps = {
  setFlashMessage: setFlashMessageAction,
};

const mapStateToProps = (state) => ({
  serviceMember: selectServiceMemberFromLoggedInUser(state),
  oktaUser: selectOktaUser(state),
});

export default connect(mapStateToProps, mapDispatchToProps)(EditOktaInfo);
