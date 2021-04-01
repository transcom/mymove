import React, { useState } from 'react';
import { Grid, GridContainer, Alert } from '@trussworks/react-uswds';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';

import ScrollToTop from 'components/ScrollToTop';
import BackupContactForm from 'components/Customer/BackupContactForm';
import { ServiceMemberShape, BackupContactShape } from 'types/customerShapes';
import {
  getResponseError,
  patchBackupContact,
  getServiceMember,
  createBackupContactForServiceMember,
} from 'services/internalApi';
import {
  updateBackupContact as updateBackupContactAction,
  updateServiceMember as updateServiceMemberAction,
} from 'store/entities/actions';
import { selectServiceMemberFromLoggedInUser, selectBackupContacts } from 'store/entities/selectors';
import requireCustomerState from 'containers/requireCustomerState/requireCustomerState';
import { profileStates } from 'constants/customerStates';
import { customerRoutes, generalRoutes } from 'constants/routes';
import scrollToTop from 'shared/scrollToTop';

export const BackupContact = ({
  serviceMember,
  currentBackupContacts,
  updateServiceMember,
  updateBackupContact,
  push,
}) => {
  const initialValues = {
    name: currentBackupContacts[0]?.name || '',
    telephone: currentBackupContacts[0]?.telephone || '',
    email: currentBackupContacts[0]?.email || '',
  };

  const NonePermission = 'NONE';

  const [serverError, setServerError] = useState(null);

  const handleBack = () => {
    return push(customerRoutes.BACKUP_ADDRESS_PATH);
  };

  const handleSubmit = (values) => {
    const payload = {
      name: values?.name || '',
      email: values?.email || '',
      telephone: values?.telephone || '',
      permission: values.permission === undefined ? NonePermission : values.permission,
    };

    const serviceMemberId = serviceMember.id;

    if (currentBackupContacts.length > 0) {
      const [firstBackupContact] = currentBackupContacts;
      payload.id = firstBackupContact.id;
      return patchBackupContact(payload)
        .then((response) => {
          updateBackupContact(response);
        })
        .then(() => getServiceMember(serviceMemberId))
        .then((response) => {
          updateServiceMember(response);
        })
        .then(() => push(generalRoutes.HOME_PATH))
        .catch((e) => {
          // TODO - error handling - below is rudimentary error handling to approximate existing UX
          // Error shape: https://github.com/swagger-api/swagger-js/blob/master/docs/usage/http-client.md#errors
          const { response } = e;
          const errorMessage = getResponseError(response, 'failed to update backup contact due to server error');
          setServerError(errorMessage);

          scrollToTop();
        });
    }
    return createBackupContactForServiceMember(serviceMemberId, payload)
      .then((response) => {
        updateBackupContact(response);
      })
      .then(() => getServiceMember(serviceMemberId))
      .then((response) => {
        updateServiceMember(response);
      })
      .then(() => push(generalRoutes.HOME_PATH))
      .catch((e) => {
        // TODO - error handling - below is rudimentary error handling to approximate existing UX
        // Error shape: https://github.com/swagger-api/swagger-js/blob/master/docs/usage/http-client.md#errors
        const { response } = e;
        const errorMessage = getResponseError(response, 'failed to create backup contact due to server error');
        setServerError(errorMessage);

        scrollToTop();
      });
  };

  return (
    <GridContainer>
      <ScrollToTop otherDep={serverError} />

      {serverError && (
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <Alert type="error" heading="An error occurred">
              {serverError}
            </Alert>
          </Grid>
        </Grid>
      )}

      <Grid row>
        <Grid col desktop={{ col: 8, offset: 2 }}>
          <BackupContactForm initialValues={initialValues} onSubmit={handleSubmit} onBack={handleBack} />
        </Grid>
      </Grid>
    </GridContainer>
  );
};

BackupContact.propTypes = {
  serviceMember: ServiceMemberShape.isRequired,
  updateServiceMember: PropTypes.func.isRequired,
  updateBackupContact: PropTypes.func.isRequired,
  currentBackupContacts: PropTypes.arrayOf(BackupContactShape).isRequired,
  push: PropTypes.func.isRequired,
};

const mapDispatchToProps = {
  updateBackupContact: updateBackupContactAction,
  updateServiceMember: updateServiceMemberAction,
};

const mapStateToProps = (state) => {
  return {
    serviceMember: selectServiceMemberFromLoggedInUser(state),
    currentBackupContacts: selectBackupContacts(state),
  };
};

export default connect(
  mapStateToProps,
  mapDispatchToProps,
)(requireCustomerState(BackupContact, profileStates.BACKUP_ADDRESS_COMPLETE));
