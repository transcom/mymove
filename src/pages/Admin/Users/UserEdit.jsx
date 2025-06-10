import { Alert } from '@trussworks/react-uswds';
import React, { useState } from 'react';
import {
  Confirm,
  DeleteButton,
  Edit,
  SaveButton,
  SelectInput,
  SimpleForm,
  TextInput,
  Toolbar,
  useRedirect,
} from 'react-admin';

import adminStyles from '../adminStyles.module.scss';

import styles from './UserEdit.module.scss';

import { deleteUser, updateUser } from 'services/adminApi';

const UserEdit = () => {
  const redirect = useRedirect();
  const [serverError, setServerError] = useState('');
  const [deleteOpen, setDeleteOpen] = useState(false);
  const [inactivateOpen, setInactivateOpen] = useState(false);
  const [userData, setUserData] = useState({});
  const handleDeleteClick = () => setDeleteOpen(true);
  const handleDeleteClose = () => setDeleteOpen(false);
  const handleInactivateClose = () => setInactivateOpen(false);

  const renderUserEditToolbar = () => {
    return (
      <Toolbar className={adminStyles.flexRight} sx={{ gap: '10px' }}>
        <DeleteButton
          mutationOptions={{
            onSuccess: async (data) => {
              // setting user data so we can use it in the delete function
              setUserData(data);
              handleDeleteClick();
            },
          }}
          sx={{
            backgroundColor: '#e1400a !important',
            width: 120,
            '&:hover': {
              opacity: '0.8',
            },
          }}
        />
        <SaveButton />
      </Toolbar>
    );
  };

  // hard deletes a user and associated roles/privileges/backupContacts
  const deleteUserHandler = async () => {
    try {
      await deleteUser(userData.id);
      redirect('./..');
    } catch (err) {
      if (err?.statusCode === 409) {
        setInactivateOpen(true);
      } else if (err?.statusCode === 403) {
        setServerError('This is an Admin User and cannot be deleted.');
      } else {
        setServerError(err?.message);
      }
      redirect(false);
    }
  };

  const inactivateUserHandler = async () => {
    const userUpdates = {
      active: false,
      oktaEmail: userData.oktaEmail,
    };
    try {
      await updateUser(userData.id, userUpdates);
      redirect('./show');
    } catch (err) {
      setServerError(err);
      redirect(false);
    }
  };

  const handleDeleteConfirm = () => {
    deleteUserHandler();
    setDeleteOpen(false);
  };

  const handleInactivateConfirm = () => {
    inactivateUserHandler();
    setInactivateOpen(false);
  };

  return (
    <Edit>
      <Confirm
        isOpen={deleteOpen}
        title={`Delete user ${userData.oktaEmail}?`}
        content="Are you sure you want to delete this user? It will delete all associated roles, privileges, and user data. This action cannot be undone."
        onConfirm={handleDeleteConfirm}
        onClose={handleDeleteClose}
      />
      <Confirm
        isOpen={inactivateOpen && userData.active}
        title={`Deletion failed for user ${userData.oktaEmail}.`}
        content="This deletion failed as this user is already tied to existing moves. Would you like to inactivate them instead?"
        onConfirm={handleInactivateConfirm}
        onClose={handleInactivateClose}
      />
      {inactivateOpen && !userData.active && (
        <Alert type="error" slim className={styles.error}>
          This deletion failed as this user is already tied to existing moves. The user is already inactive.
        </Alert>
      )}
      {serverError && (
        <Alert type="error" slim className={styles.error}>
          {serverError}
        </Alert>
      )}
      <SimpleForm
        toolbar={renderUserEditToolbar()}
        sx={{ '& .MuiInputBase-input': { width: 232 } }}
        mode="onBlur"
        reValidateMode="onBlur"
      >
        <TextInput source="id" disabled />
        <TextInput source="oktaEmail" />
        <SelectInput
          source="active"
          choices={[
            { id: true, name: 'Yes' },
            { id: false, name: 'No' },
          ]}
          sx={{ width: 256 }}
        />
        <SelectInput
          source="revokeAdminSession"
          choices={[
            { id: true, name: 'Yes' },
            { id: false, name: 'No' },
          ]}
          sx={{ width: 256 }}
        />
        <SelectInput
          source="revokeOfficeSession"
          choices={[
            { id: true, name: 'Yes' },
            { id: false, name: 'No' },
          ]}
          sx={{ width: 256 }}
        />
        <SelectInput
          source="revokeMilSession"
          choices={[
            { id: true, name: 'Yes' },
            { id: false, name: 'No' },
          ]}
          sx={{ width: 256 }}
        />
        <TextInput source="createdAt" disabled />
        <TextInput source="updatedAt" disabled />
      </SimpleForm>
    </Edit>
  );
};

export default UserEdit;
