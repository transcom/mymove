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

import styles from './UserEdit.module.scss';

import { deleteUser, updateUser } from 'services/adminApi';

const UserEdit = () => {
  const redirect = useRedirect();
  const [serverError, setServerError] = useState('');
  const [deleteOpen, setDeleteOpen] = useState(false);
  const [disableOpen, setDisableOpen] = useState(false);
  const [userData, setUserData] = useState({});
  const handleDeleteClick = () => setDeleteOpen(true);
  const handleDeleteClose = () => setDeleteOpen(false);
  const handleDisableClose = () => setDisableOpen(false);

  const renderUserEditToolbar = () => {
    return (
      <Toolbar>
        <SaveButton />
        <DeleteButton
          mutationOptions={{
            onSuccess: async (data) => {
              // setting user data so we can use it in the delete function
              setUserData(data);
              handleDeleteClick();
            },
          }}
        />
      </Toolbar>
    );
  };

  // hard deletes a user and associated roles/privileges
  // cannot be undone, but the user is shown a confirmation modal to avoid oopsies
  const deleteUserHandler = async () => {
    await deleteUser(userData.id)
      .then(() => {
        redirect('/');
      })
      .catch(() => {
        setDisableOpen(true);
        redirect(false);
      });
  };

  const disableUserHandler = async () => {
    userData.active = false;
    await updateUser(userData.id, userData)
      .then(() => {
        redirect('/');
      })
      .catch((error) => {
        setServerError(error);
        redirect(false);
      });
  };

  const handleDeleteConfirm = () => {
    deleteUserHandler();
    setDeleteOpen(false);
  };

  const handleDisableConfirm = () => {
    disableUserHandler();
    setDisableOpen(false);
  };

  return (
    <Edit>
      <Confirm
        isOpen={deleteOpen}
        title={`Delete user ${userData.oktaEmail} ?`}
        content="Are you sure you want to delete this user? It will delete all associated roles, privileges, and user data. This action cannot be undone."
        onConfirm={handleDeleteConfirm}
        onClose={handleDeleteClose}
      />
      <Confirm
        isOpen={disableOpen && userData.active}
        title={`Deletion failed for user ${userData.oktaEmail}`}
        content="This deletion failed as this user is already tied to existing moves. Would you like to disable them instead?"
        onConfirm={handleDisableConfirm}
        onClose={handleDisableClose}
      />
      {disableOpen && !userData.active && (
        <Alert type="error" slim className={styles.error}>
          This deletion failed as this user is already tied to existing moves. The user is already disabled.
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
