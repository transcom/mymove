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
    oktaUsername: oktaUser?.login || 'Not Provided',
    oktaEmail: oktaUser?.email || 'Not Provided',
    oktaFirstName: oktaUser?.firstName || 'Not Provided',
    oktaLastName: oktaUser?.lastName || 'Not Provided',
    oktaEdipi: oktaUser?.cac_edipi || '',
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
        login: values?.oktaUsername,
        email: values?.oktaEmail,
        firstName: values?.oktaFirstName,
        lastName: values?.oktaLastName,
        cac_edipi: values?.oktaEdipi,
      },
    };

    return updateOktaUser(oktaPayload)
      .then((response) => {
        console.log(response);
        setFlashMessage('EDIT_OKTA_PROFILE_SUCCESS', 'success', "You've updated your Okta profile.");
        navigate(customerRoutes.PROFILE_PATH);
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
