import React from 'react';
import {
  Edit,
  SimpleForm,
  TextInput,
  SelectInput,
  Toolbar,
  SaveButton,
  AutocompleteInput,
  ReferenceInput,
  ArrayInput,
  SimpleFormIterator,
  BooleanInput,
  useDataProvider,
  Button,
} from 'react-admin';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { connect } from 'react-redux';

import { RolesPrivilegesCheckboxInput } from 'scenes/SystemAdmin/shared/RolesPrivilegesCheckboxes';
import { roleTypes } from 'constants/userRoles';
import { selectAdminUser } from 'store/entities/selectors';

const OfficeUserEditToolbar = (props) => {
  return (
    <Toolbar {...props}>
      <SaveButton />
    </Toolbar>
  );
};

const OfficeUserEdit = ({ adminUser }) => {
  const dataProvider = useDataProvider();

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
        'At least one of your transportation office assigmnets is blank. Please select a transportation office or remove that assignment.';
    }

    return errors;
  };

  return (
    <Edit mutationMode="pessimistic">
      <SimpleForm
        toolbar={<OfficeUserEditToolbar />}
        sx={{ '& .MuiInputBase-input': { width: 232 } }}
        mode="onSubmit"
        reValidateMode="onSubmit"
        validate={validateForm}
      >
        <TextInput source="id" disabled />
        <TextInput source="userId" label="User Id" disabled />
        <TextInput source="email" disabled />
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
                  backgroundColor: '#1976d2',
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
                  color: 'white',
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
