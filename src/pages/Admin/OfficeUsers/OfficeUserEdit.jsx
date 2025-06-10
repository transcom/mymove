import { Alert } from '@trussworks/react-uswds';
import React, { useState } from 'react';
import {
  Edit,
  SimpleForm,
  TextInput,
  SelectInput,
  AutocompleteInput,
  ReferenceInput,
  ArrayInput,
  SimpleFormIterator,
  BooleanInput,
  useDataProvider,
  useRedirect,
  Button,
  DeleteButton,
  Confirm,
  SaveButton,
  Toolbar,
} from 'react-admin';
import { connect } from 'react-redux';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import adminStyles from '../adminStyles.module.scss';

import styles from './OfficeUserEdit.module.scss';

import { RolesPrivilegesCheckboxInput } from 'scenes/SystemAdmin/shared/RolesPrivilegesCheckboxes';
import { roleTypes } from 'constants/userRoles';
import { selectAdminUser } from 'store/entities/selectors';
import { deleteOfficeUser, updateOfficeUser } from 'services/adminApi';

const OfficeUserEdit = ({ adminUser }) => {
  const dataProvider = useDataProvider();
  const redirect = useRedirect();
  const [serverError, setServerError] = useState('');
  const [deleteOpen, setDeleteOpen] = useState(false);
  const [inactivateOpen, setInactivateOpen] = useState(false);
  const [userData, setUserData] = useState({});
  const handleDeleteClick = () => setDeleteOpen(true);
  const handleDeleteClose = () => setDeleteOpen(false);
  const handleInactivateClose = () => setInactivateOpen(false);

  const validateForm = async (values) => {
    const errors = {};
    if (!values.firstName) {
      errors.firstName = 'You must enter a first name.';
    }
    if (!values.lastName) {
      errors.lastName = 'You must enter a last name.';
    }
    if (!values.email) {
      errors.email = 'You must enter an email.';
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
        'You cannot select both Task Ordering Officer and Task Invoicing Officer. This is a policy managed by USTRANSCOM.';
    }

    if (!values?.transportationOfficeAssignments?.length) {
      errors.transportationOfficeAssignments = 'You must add at least one transportation office assignment. ';
    }

    if (values?.transportationOfficeAssignments?.length > 2) {
      errors.transportationOfficeAssignments = 'You cannot add more than two transportation office assignments.';
    }

    if (values?.transportationOfficeAssignments?.filter((toa) => toa.primaryOffice)?.length > 1) {
      errors.transportationOfficeAssignments = 'You cannot designate more than one primary transportation office.';
    }

    if (values?.transportationOfficeAssignments?.filter((toa) => toa.primaryOffice)?.length < 1) {
      errors.transportationOfficeAssignments = 'You must designate a primary transportation office.';
    }

    const gblocs = new Set();
    if (values?.transportationOfficeAssignments[0]?.transportationOfficeId) {
      const { data } = await dataProvider.getOne('offices', {
        id: values.transportationOfficeAssignments[0].transportationOfficeId,
      });
      await gblocs.add(data.gbloc);
    }

    if (values?.transportationOfficeAssignments[1]?.transportationOfficeId) {
      const { data } = await dataProvider.getOne('offices', {
        id: values.transportationOfficeAssignments[1].transportationOfficeId,
      });
      await gblocs.add(data.gbloc);
    }

    if (values?.transportationOfficeAssignments.length !== gblocs.size) {
      errors.transportationOfficeAssignments = 'Each assigned transportation office must be in a different GBLOC.';
    }

    if (values?.transportationOfficeAssignments?.filter((toa) => toa.transportationOfficeId == null)?.length > 0) {
      errors.transportationOfficeAssignments =
        'At least one of your transportation office assignments is blank. Please select a transportation office or remove that assignment.';
    }

    return errors;
  };

  // hard deletes an office user and associated roles/privileges
  const deleteUser = async () => {
    try {
      await deleteOfficeUser(userData.id);
      redirect('./..');
    } catch (err) {
      if (err?.statusCode === 409) {
        setInactivateOpen(true);
      } else {
        setServerError(err?.message);
      }
      redirect(false);
    }
  };

  const inactivateUser = async () => {
    const userUpdates = {
      active: false,
    };
    try {
      await updateOfficeUser(userData.id, userUpdates);
      redirect('./show');
    } catch (err) {
      setServerError(err);
      redirect(false);
    }
  };

  const handleDeleteConfirm = () => {
    deleteUser();
    setDeleteOpen(false);
  };

  const handleInactivateConfirm = () => {
    inactivateUser();
    setInactivateOpen(false);
  };

  // rendering tool bar
  const renderToolBar = () => {
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

  return (
    <Edit>
      <Confirm
        isOpen={deleteOpen}
        title={`Delete office user ${userData.firstName} ${userData.lastName}?`}
        content="Are you sure you want to delete this user? It will delete all associated roles, privileges, and user data. This action cannot be undone."
        onConfirm={handleDeleteConfirm}
        onClose={handleDeleteClose}
      />
      <Confirm
        isOpen={inactivateOpen && userData.active}
        title={`Deletion failed for user ${userData.firstName} ${userData.lastName}.`}
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
        toolbar={renderToolBar()}
        sx={{ '& .MuiInputBase-input': { width: 232 } }}
        mode="onSubmit"
        reValidateMode="onSubmit"
        validate={validateForm}
      >
        <TextInput source="id" disabled />
        <TextInput source="userId" label="User Id" disabled />
        <TextInput source="email" />
        <TextInput source="firstName" />
        <TextInput source="middleInitials" />
        <TextInput source="lastName" />
        <TextInput source="telephone" />
        <SelectInput
          source="active"
          choices={[
            { id: true, name: 'Yes' },
            { id: false, name: 'No' },
          ]}
          sx={{ width: 256 }}
        />
        <RolesPrivilegesCheckboxInput source="roles" adminUser={adminUser} />
        <ArrayInput source="transportationOfficeAssignments" label="Transportation Office Assignments (Maximum: 2)">
          <SimpleFormIterator
            inline
            disableClear
            disableReordering
            addButton={
              <Button
                type="button"
                size="extrasmall"
                data-testid="addTransportationOfficeButton"
                sx={{
                  backgroundColor: '#005ea2',
                  '&:hover': {
                    backgroundColor: '#1565c0',
                  },
                  color: 'white',
                }}
                label="Add transportation office assignment"
              />
            }
            removeButton={
              <Button
                type="button"
                size="extrasmall"
                data-testid="removeTransportationOfficeButton"
                sx={{
                  backgroundColor: '#e1400a',
                  '&:hover': {
                    backgroundColor: '#d23c0f',
                  },
                  width: '100px',
                  color: 'white',
                  visibility: 'visible',
                  opacity: 1,
                }}
                label="remove"
              >
                <FontAwesomeIcon icon="trash" />
              </Button>
            }
          >
            <ReferenceInput
              label="Transportation Office"
              reference="offices"
              source="transportationOfficeId"
              perPage={500}
            >
              <TransportationOfficePicker />
            </ReferenceInput>
            <BooleanInput source="primaryOffice" label="Primary Office" defaultValue={false} />
          </SimpleFormIterator>
        </ArrayInput>
        <TextInput source="createdAt" disabled />
        <TextInput source="updatedAt" disabled />
      </SimpleForm>
    </Edit>
  );
};

const TransportationOfficePicker = (props) => {
  return (
    <>
      <AutocompleteInput optionText="name" sx={{ width: 256 }} {...props} />
      <SelectInput source="offices.gbloc" label="GBLOC" optionText="gbloc" {...props} disabled sx={{ width: 128 }} />
    </>
  );
};

function mapStateToProps(state) {
  return {
    adminUser: selectAdminUser(state),
  };
}

export default connect(mapStateToProps)(OfficeUserEdit);
