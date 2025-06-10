import { Alert } from '@trussworks/react-uswds';
import React from 'react';
import {
  ArrayField,
  Datagrid,
  DateField,
  ReferenceField,
  Show,
  SimpleShowLayout,
  TextField,
  useRecordContext,
  useRedirect,
  DeleteButton,
  Toolbar,
  Confirm,
} from 'react-admin';
import { useNavigate } from 'react-router';

import adminStyles from '../adminStyles.module.scss';

import styles from './RejectedOfficeUserShow.module.scss';

import { deleteOfficeUser } from 'services/adminApi';
import { adminRoutes } from 'constants/routes';

const RejectedOfficeUserShowTitle = () => {
  const record = useRecordContext();

  return <span>{`${record?.firstName} ${record?.lastName}`}</span>;
};

const RejectedOfficeUserShowRoles = () => {
  const record = useRecordContext();
  if (!record?.roles) return <p>This user has not requested any roles.</p>;

  return (
    <ArrayField source="roles">
      <span>Roles Requested:</span>
      <Datagrid bulkActionButtons={false}>
        <TextField source="roleName" />
      </Datagrid>
    </ArrayField>
  );
};

const RejectedOfficeUserShow = () => {
  const redirect = useRedirect();
  const navigate = useNavigate();
  const [serverError, setServerError] = React.useState('');
  const [open, setOpen] = React.useState(false);
  const [userData, setUserData] = React.useState({});

  const handleClick = () => setOpen(true);
  const handleDialogClose = () => setOpen(false);

  // hard deletes a user and associated roles/privileges
  // cannot be undone, but the user is shown a confirmation modal to avoid oopsies
  const deleteUser = async () => {
    await deleteOfficeUser(userData.id)
      .then(() => {
        navigate(adminRoutes.REJECTED_OFFICE_USERS);
      })
      .catch((error) => {
        setServerError(error);
        redirect(false);
      });
  };

  const handleConfirm = () => {
    deleteUser();
    setOpen(false);
  };

  return (
    <Show title={<RejectedOfficeUserShowTitle />}>
      <Confirm
        isOpen={open}
        title={`Delete rejected office user ${userData.firstName} ${userData.lastName}?`}
        content="Are you sure you want to delete this user? It will delete all associated roles, privileges, and user data. This action cannot be undone."
        onConfirm={handleConfirm}
        onClose={handleDialogClose}
      />
      <SimpleShowLayout>
        <TextField source="id" />
        <TextField source="userId" label="User Id" />
        <TextField source="status" />
        <TextField source="email" />
        <TextField source="firstName" />
        <TextField source="middleInitials" />
        <TextField source="lastName" />
        <TextField source="telephone" />
        <TextField source="edipi" label="DODID#" />
        <TextField source="otherUniqueId" label="Other unique Id" />
        <RejectedOfficeUserShowRoles />
        <ReferenceField label="Transportation Office" source="transportationOfficeId" reference="offices" sortBy="name">
          <TextField component="pre" source="name" />
        </ReferenceField>
        <DateField label="Account rejected at" source="createdAt" showTime />
      </SimpleShowLayout>
      {serverError && (
        <Alert type="error" slim className={styles.error}>
          {serverError}
        </Alert>
      )}
      <Toolbar className={adminStyles.flexRight}>
        <DeleteButton
          mutationOptions={{
            onSuccess: async (data) => {
              // setting user data so we can use it in the delete function
              setUserData(data);
              handleClick();
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
      </Toolbar>
    </Show>
  );
};

export default RejectedOfficeUserShow;
