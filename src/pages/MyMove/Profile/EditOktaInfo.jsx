import PropTypes from 'prop-types';
import React, { useState } from 'react';
import { connect } from 'react-redux';
import { useLocation, useNavigate } from 'react-router-dom';
import { Alert, Grid, GridContainer } from '@trussworks/react-uswds';

import { OktaUserInfoShape } from 'types/user';
import EditOktaInfoForm from 'components/Customer/EditOktaInfoForm/EditOktaInfoForm';
import NotificationScrollToTop from 'components/NotificationScrollToTop';
import { customerRoutes } from 'constants/routes';
import { getResponseError, updateOktaUser } from 'services/internalApi';
import { selectServiceMemberFromLoggedInUser, selectOktaUser } from 'store/entities/selectors';
import { updateOktaUserState as updateOktaUserStateAction } from 'store/entities/actions';
import { setFlashMessage as setFlashMessageAction } from 'store/flash/actions';

export const EditOktaInfo = ({ serviceMember, setFlashMessage, oktaUser, updateOktaUserState }) => {
  const navigate = useNavigate();
  const { state } = useLocation();
  const [serverError, setServerError] = useState(null);
  const [noChangeError, setNoChangeError] = useState(null);

  const initialValues = {
    oktaUsername: oktaUser?.login || 'Not Provided',
    oktaEmail: oktaUser?.email || 'Not Provided',
    oktaFirstName: oktaUser?.firstName || 'Not Provided',
    oktaLastName: oktaUser?.lastName || 'Not Provided',
    oktaEdipi: oktaUser?.cac_edipi || '',
    oktaSub: oktaUser?.sub,
  };

  const handleCancel = () => {
    navigate(customerRoutes.PROFILE_PATH, { state });
  };

  // sends POST request to Okta API with form values
  // then updates the state with updated values
  // sends the user back to profile page with confirmation banner
  const handleSubmit = async (values) => {
    // wrapping values in profile due to Okta API requirements
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

    // checking to see if the values are the same to avoid unnecessary api call
    if (
      oktaPayload.profile.cac_edipi === oktaUser.cac_edipi &&
      oktaPayload.profile.email === oktaUser.email &&
      oktaPayload.profile.firstName === oktaUser.firstName &&
      oktaPayload.profile.lastName === oktaUser.lastName
    ) {
      setNoChangeError(true); // if true, we'll let the customer know
    } else {
      return updateOktaUser(oktaPayload)
        .then((response) => {
          updateOktaUserState(response);
          setFlashMessage('EDIT_OKTA_PROFILE_SUCCESS', 'success', "You've updated your Okta profile.");
          navigate(customerRoutes.PROFILE_PATH, { state });
        })
        .catch((e) => {
          const { response } = e;
          const errorMessage = getResponseError(response, 'Failed to update okta profile due to server error');

          setServerError(errorMessage);
        });
    }
    return Promise.resolve();
  };

  return (
    <GridContainer>
      <NotificationScrollToTop dependency={serverError || noChangeError} />

      {serverError && (
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <Alert type="error" headingLevel="h4" heading="An error occurred">
              {serverError}
            </Alert>
          </Grid>
        </Grid>
      )}
      {noChangeError && (
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <Alert type="warning" headingLevel="h4" heading="No changes were made">
              You must make some changes if you want to edit your Okta profile.
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
  updateOktaUserState: updateOktaUserStateAction,
};

const mapStateToProps = (state) => ({
  serviceMember: selectServiceMemberFromLoggedInUser(state),
  oktaUser: selectOktaUser(state),
});

export default connect(mapStateToProps, mapDispatchToProps)(EditOktaInfo);
