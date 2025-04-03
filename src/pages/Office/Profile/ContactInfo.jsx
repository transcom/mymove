import React, { useState } from 'react';
import 'styles/office.scss';
import { Grid, GridContainer, Alert } from '@trussworks/react-uswds';
import { connect } from 'react-redux';
import { useNavigate } from 'react-router-dom';
import PropTypes from 'prop-types';

import styles from './Profile.module.scss';

import NotificationScrollToTop from 'components/NotificationScrollToTop';
import ContactInfoForm from 'components/Office/Profile/ContactInfoForm/ContactInfoForm';
import { OfficeUserInfoShape } from 'types/index';
import { selectLoggedInUser } from 'store/entities/selectors';
import { officeRoutes } from 'constants/routes';
import { setFlashMessage as setFlashMessageAction } from 'store/flash/actions';
import { loadUser as loadUserAction } from 'store/auth/actions';
import { patchOfficeUser } from 'services/ghcApi';

export const ContactInfo = ({ officeUser, setFlashMessage, loadUser }) => {
  const navigate = useNavigate();
  const initialValues = {
    firstName: officeUser?.first_name || '',
    middleName: officeUser?.middle_name || '',
    lastName: officeUser?.last_name || '',
    telephone: officeUser?.telephone || '',
    email: officeUser?.email || '',
  };

  const [serverError, setServerError] = useState(null);

  const handleBack = () => {
    return navigate(officeRoutes.PROFILE_PATH);
  };

  const handleSubmit = async (values) => {
    const payload = {
      telephone: values?.telephone,
    };

    return patchOfficeUser(officeUser?.id, payload)
      .then(() => {
        setFlashMessage('EDIT_CONTACT_INFO_SUCCESS', 'success', "You've updated your information.");
        loadUser();
        navigate(officeRoutes.PROFILE_PATH);
      })
      .catch((e) => {
        const errorMessage = e.response?.body?.detail || e;
        const error = `Failed to update contact info due to server error: ${errorMessage}`;
        setServerError(error);
      });
  };

  return (
    <div className={styles.Profile}>
      <GridContainer>
        <Grid row>
          <Grid col={4} desktop={{ col: 12 }} tablet={{ col: 8 }}>
            <NotificationScrollToTop dependency={serverError} />

            {serverError && (
              <Alert type="error" headingLevel="h4" heading="An error occurred">
                {serverError}
              </Alert>
            )}
            <div>
              <h1>Edit contact info</h1>
            </div>
            <ContactInfoForm initialValues={initialValues} onSubmit={handleSubmit} onCancel={handleBack} />
          </Grid>
        </Grid>
      </GridContainer>
    </div>
  );
};

ContactInfo.propTypes = {
  officeUser: OfficeUserInfoShape,
  setFlashMessage: PropTypes.func.isRequired,
};

ContactInfo.defaultProps = {
  officeUser: {},
};

const mapDispatchToProps = {
  setFlashMessage: setFlashMessageAction,
  loadUser: loadUserAction,
};

const mapStateToProps = (state) => {
  const user = selectLoggedInUser(state);

  return {
    officeUser: user?.office_user || {},
  };
};

export default connect(mapStateToProps, mapDispatchToProps)(ContactInfo);
