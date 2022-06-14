import PropTypes from 'prop-types';
import React, { useState } from 'react';
import { connect } from 'react-redux';
import { useHistory } from 'react-router-dom';
import { Alert, Grid, GridContainer } from '@trussworks/react-uswds';

import EditContactInfoForm, {
  backupAddressName,
  backupContactName,
  residentialAddressName,
} from 'components/Customer/EditContactInfoForm/EditContactInfoForm';
import ScrollToTop from 'components/ScrollToTop';
import { customerRoutes } from 'constants/routes';
import { getResponseError, patchBackupContact, patchServiceMember } from 'services/internalApi';
import {
  updateBackupContact as updateBackupContactAction,
  updateServiceMember as updateServiceMemberAction,
} from 'store/entities/actions';
import { selectBackupContacts, selectServiceMemberFromLoggedInUser } from 'store/entities/selectors';
import { setFlashMessage as setFlashMessageAction } from 'store/flash/actions';
import { BackupContactShape, ServiceMemberShape } from 'types/customerShapes';

export const EditContactInfo = ({
  currentBackupContacts,
  serviceMember,
  setFlashMessage,
  updateBackupContact,
  updateServiceMember,
}) => {
  const history = useHistory();
  const [serverError, setServerError] = useState(null);

  const initialValues = {
    telephone: serviceMember?.telephone || '',
    secondary_telephone: serviceMember?.secondary_telephone || '',
    personal_email: serviceMember?.personal_email || '',
    phone_is_preferred: serviceMember?.phone_is_preferred,
    email_is_preferred: serviceMember?.email_is_preferred,
    [residentialAddressName]: {
      streetAddress1: serviceMember.residential_address?.streetAddress1 || '',
      streetAddress2: serviceMember.residential_address?.streetAddress2 || '',
      city: serviceMember.residential_address?.city || '',
      state: serviceMember.residential_address?.state || '',
      postalCode: serviceMember.residential_address?.postalCode || '',
    },
    [backupAddressName]: {
      streetAddress1: serviceMember.backup_mailing_address?.streetAddress1 || '',
      streetAddress2: serviceMember.backup_mailing_address?.streetAddress2 || '',
      city: serviceMember.backup_mailing_address?.city || '',
      state: serviceMember.backup_mailing_address?.state || '',
      postalCode: serviceMember.backup_mailing_address?.postalCode || '',
    },
    [backupContactName]: {
      name: currentBackupContacts[0]?.name || '',
      telephone: currentBackupContacts[0]?.telephone || '',
      email: currentBackupContacts[0]?.email || '',
    },
  };

  const handleCancel = () => {
    history.push(customerRoutes.PROFILE_PATH);
  };

  const handleSubmit = async (values) => {
    const serviceMemberPayload = {
      id: serviceMember.id,
      telephone: values?.telephone,
      personal_email: values?.personal_email,
      phone_is_preferred: values?.phone_is_preferred,
      email_is_preferred: values?.email_is_preferred,
      residential_address: values[residentialAddressName.toString()],
      backup_mailing_address: values[backupAddressName.toString()],
    };

    if (values?.secondary_telephone) {
      serviceMemberPayload.secondary_telephone = values?.secondary_telephone;
    }

    const backupContactPayload = {
      id: currentBackupContacts[0].id,
      name: values[backupContactName.toString()]?.name || '',
      email: values[backupContactName.toString()]?.email || '',
      telephone: values[backupContactName.toString()]?.telephone || '',
      permission: currentBackupContacts[0].permission,
    };

    const backupContactChanged =
      initialValues[backupContactName.toString()].name !== backupContactPayload.name ||
      initialValues[backupContactName.toString()].email !== backupContactPayload.email ||
      initialValues[backupContactName.toString()].telephone !== backupContactPayload.telephone;

    // If only backup contact info is updated, we could call patchBackupContact, then we'd need to do an api call to
    // getServiceMember in order to get the latest info, and then call updateServiceMember. Conversely, we can just call
    // patchServiceMember even if no other service member data changed, since it returns the same data as
    // getServiceMember request would, and it makes our logic here simpler.

    if (backupContactChanged) {
      let error;
      const patchPromise = await patchBackupContact(backupContactPayload)
        .then(updateBackupContact)
        .catch((e) => {
          //     // TODO - error handling - below is rudimentary error handling to approximate existing UX
          //     // Error shape: https://github.com/swagger-api/swagger-js/blob/master/docs/usage/http-client.md#errors
          const { response } = e;
          const errorMessage = getResponseError(response, 'Failed to update backup contact due to server error');

          setServerError(errorMessage);
          error = true;
        });

      if (error) {
        return patchPromise;
      }
    }

    return patchServiceMember(serviceMemberPayload)
      .then(updateServiceMember)
      .then(() => {
        setFlashMessage('EDIT_CONTACT_INFO_SUCCESS', 'success', "You've updated your information.");
        history.push(customerRoutes.PROFILE_PATH);
      })
      .catch((e) => {
        //     // TODO - error handling - below is rudimentary error handling to approximate existing UX
        //     // Error shape: https://github.com/swagger-api/swagger-js/blob/master/docs/usage/http-client.md#errors
        const { response } = e;
        const errorMessage = getResponseError(response, 'Failed to update service member due to server error');

        setServerError(errorMessage);
      });
  };

  return (
    <GridContainer>
      <ScrollToTop otherDep={serverError} />

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
          <EditContactInfoForm initialValues={initialValues} onCancel={handleCancel} onSubmit={handleSubmit} />
        </Grid>
      </Grid>
    </GridContainer>
  );
};

EditContactInfo.propTypes = {
  currentBackupContacts: PropTypes.arrayOf(BackupContactShape).isRequired,
  serviceMember: ServiceMemberShape.isRequired,
  setFlashMessage: PropTypes.func.isRequired,
  updateBackupContact: PropTypes.func.isRequired,
  updateServiceMember: PropTypes.func.isRequired,
};

const mapDispatchToProps = {
  setFlashMessage: setFlashMessageAction,
  updateBackupContact: updateBackupContactAction,
  updateServiceMember: updateServiceMemberAction,
};

const mapStateToProps = (state) => ({
  currentBackupContacts: selectBackupContacts(state),
  serviceMember: selectServiceMemberFromLoggedInUser(state),
});

export default connect(mapStateToProps, mapDispatchToProps)(EditContactInfo);
