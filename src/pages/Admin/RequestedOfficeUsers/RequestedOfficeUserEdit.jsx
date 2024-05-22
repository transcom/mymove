import { Alert } from '@trussworks/react-uswds';
import React, { useState } from 'react';
import {
  Edit,
  SimpleForm,
  TextInput,
  required,
  Toolbar,
  SaveButton,
  AutocompleteInput,
  ReferenceInput,
  useRecordContext,
  useRedirect,
} from 'react-admin';

import styles from './RequestedOfficeUserShow.module.scss';

import { RolesPrivilegesCheckboxInput } from 'scenes/SystemAdmin/shared/RolesPrivilegesCheckboxes';
import { edipiValidator, phoneValidators } from 'scenes/SystemAdmin/shared/form_validators';
import { updateRequestedOfficeUser } from 'services/adminApi';
import { roleTypes } from 'constants/userRoles';

const RequestedOfficeUserShowTitle = () => {
  const record = useRecordContext();

  return <span>{`${record?.firstName} ${record?.lastName}`}</span>;
};

const validateForm = (values) => {
  const errors = {};
  if (!values.firstName) {
    errors.firstName = 'You must enter a first name.';
  }
  if (!values.lastName) {
    errors.lastName = 'You must enter a last name.';
  }

  if (!values.telephone) {
    errors.telephone = 'You must enter a telephone number.';
  } else if (!values.telephone.match(/^[2-9]\d{2}-\d{3}-\d{4}$/)) {
    errors.telephone = 'Invalid phone number, should be 000-000-0000.';
  }

  if (!values.roles?.length) {
    errors.roles = 'You must select at least one role.';
  } else if (
    values.roles.find((role) => role.roleType === roleTypes.TIO) &&
    values.roles.find((role) => role.roleType === roleTypes.TOO)
  ) {
    errors.roles =
      'You cannot select both Transportation Ordering Officer and Transportation Invoicing Officer. This is a policy managed by USTRANSCOM.';
  }

  if (!values.transportationOfficeId) {
    errors.transportationOfficeId = 'You must select a transportation office.';
  }

  return errors;
};

const RequestedOfficeUserEdit = () => {
  const redirect = useRedirect();
  const [serverError, setServerError] = useState('');
  const [validationCheck, setValidationCheck] = useState('');

  // rejects the user with all relevant updates made by admin
  // performs validation to ensure the rejection reason was provided
  const reject = async (user) => {
    if (!user.rejectionReason) {
      setValidationCheck('You must provide a rejection reason when rejecting a user');
    } else {
      setValidationCheck('');
      const body = {
        email: user.email,
        edipi: user.edipi,
        firstName: user.firstName,
        middleInitials: user.middleInitials,
        lastName: user.lastName,
        otherUniqueId: user.otherUniqueId,
        rejectionReason: user.rejectionReason,
        roles: user.roles,
        status: 'REJECTED',
        telephone: user.telephone,
        transportationOfficeId: user.transportationOfficeId,
      };
      updateRequestedOfficeUser(user.id, body)
        .then(() => {
          redirect('/');
        })
        .catch((error) => {
          setServerError(error);
          redirect(false);
        });
    }
  };

  // approves the user with all relevant updates made by admin
  // performs validation to ensure either edipi or otherUniqueId was provided
  const approve = async (user) => {
    if (!user.edipi && !user.otherUniqueId) {
      setValidationCheck('You must provide an DODID# or unique ID for the user');
    } else {
      setValidationCheck('');
      const body = {
        email: user.email,
        edipi: user.edipi,
        firstName: user.firstName,
        middleInitials: user.middleInitials,
        lastName: user.lastName,
        otherUniqueId: user.otherUniqueId,
        rejectionReason: user.rejectionReason,
        roles: user.roles,
        status: 'APPROVED',
        telephone: user.telephone,
        transportationOfficeId: user.transportationOfficeId,
      };
      updateRequestedOfficeUser(user.id, body)
        .then(() => {
          redirect('/');
        })
        .catch((error) => {
          setServerError(error);
          redirect(false);
        });
    }
  };

  // rendering tool bar with added error/validation alerts
  const renderToolBar = () => {
    return (
      <>
        {serverError && (
          <Alert type="error" slim className={styles.error}>
            {serverError}
          </Alert>
        )}
        {validationCheck && (
          <Alert type="error" slim className={styles.rejErrorEdit}>
            {validationCheck}
          </Alert>
        )}
        <Toolbar sx={{ display: 'flex', gap: '10px' }}>
          <SaveButton
            type="button"
            alwaysEnable
            label="Approve"
            mutationOptions={{
              onSuccess: async (data) => {
                await approve(data);
              },
            }}
          />
          <SaveButton
            type="button"
            color="error"
            alwaysEnable
            label="Reject"
            mutationOptions={{
              onSuccess: async (data) => {
                await reject(data);
              },
            }}
          />
        </Toolbar>
      </>
    );
  };

  return (
    <Edit title={<RequestedOfficeUserShowTitle />}>
      <SimpleForm
        toolbar={renderToolBar()}
        sx={{ '& .MuiInputBase-input': { width: 232 } }}
        reValidateMode="onBlur"
        mode="onBlur"
        validate={validateForm}
      >
        <TextInput source="id" disabled />
        <TextInput source="userId" label="User Id" disabled />
        <TextInput source="email" disabled />
        <TextInput source="firstName" validate={required()} />
        <TextInput source="middleInitials" />
        <TextInput source="lastName" validate={required()} />
        <TextInput source="edipi" label="DODID#" validate={edipiValidator} />
        <TextInput source="otherUniqueId" label="Other unique Id" />
        <TextInput source="telephone" validate={phoneValidators} />
        <RolesPrivilegesCheckboxInput source="roles" />
        <ReferenceInput
          label="Transportation Office"
          reference="offices"
          source="transportationOfficeId"
          perPage={500}
          validate={required()}
        >
          <AutocompleteInput optionText="name" sx={{ width: 256 }} />
        </ReferenceInput>
        <TextInput source="createdAt" label="Requested at" disabled />
        <TextInput source="rejectionReason" className={styles.rejReasonInput} />
      </SimpleForm>
    </Edit>
  );
};

export default RequestedOfficeUserEdit;
