import React from 'react';
import {
  Edit,
  SimpleForm,
  TextInput,
  SelectInput,
  required,
  Toolbar,
  SaveButton,
  AutocompleteInput,
  ReferenceInput,
  ArrayInput,
  SimpleFormIterator,
  BooleanInput,
} from 'react-admin';
import { connect } from 'react-redux';

import { RolesPrivilegesCheckboxInput } from 'scenes/SystemAdmin/shared/RolesPrivilegesCheckboxes';
import { phoneValidators } from 'scenes/SystemAdmin/shared/form_validators';
import { roleTypes } from 'constants/userRoles';
import { selectAdminUser } from 'store/entities/selectors';

const OfficeUserEditToolbar = (props) => {
  return (
    <Toolbar {...props}>
      <SaveButton />
    </Toolbar>
  );
};
const validateForm = (values) => {
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

  if (!values.transportationOfficeAssignments?.length) {
    errors.transportationOfficeAssignments = 'You must select at least one transportation office.';
  } else if (values.transportationOfficeAssignments?.length > 2) {
    errors.transportationOfficeAssignments = 'You cannot select more than two transportation offices.';
  }

  if (values.transportationOfficeAssignments?.filter((toa) => toa.primaryOffice)?.length > 1) {
    errors.transportationOfficeAssignments = values.transportationOfficeAssignments.map((office) => {
      const officeErrors = {};
      if (office.primaryOffice) {
        officeErrors.primaryOffice = `You cannot designate more than one primary transportation office.`;
      }
      return officeErrors;
    });
  }

  if (values.transportationOfficeAssignments?.filter((toa) => toa.primaryOffice)?.length < 1) {
    errors.transportationOfficeAssignments = values.transportationOfficeAssignments.map((office) => {
      const officeErrors = {};
      if (!office.primaryOffice) {
        officeErrors.primaryOffice = `You must designate a primary transportation office.`;
      }
      return officeErrors;
    });
  }

  return errors;
};

const OfficeUserEdit = ({ adminUser }) => (
  <Edit>
    <SimpleForm
      toolbar={<OfficeUserEditToolbar />}
      sx={{ '& .MuiInputBase-input': { width: 232 } }}
      mode="onBlur"
      reValidateMode="onBlur"
      validate={validateForm}
    >
      <TextInput source="id" disabled />
      <TextInput source="userId" label="User Id" disabled />
      <TextInput source="email" disabled />
      <TextInput source="firstName" validate={required()} />
      <TextInput source="middleInitials" />
      <TextInput source="lastName" validate={required()} />
      <TextInput source="telephone" validate={phoneValidators} />
      <SelectInput
        source="active"
        choices={[
          { id: true, name: 'Yes' },
          { id: false, name: 'No' },
        ]}
        sx={{ width: 256 }}
      />
      <RolesPrivilegesCheckboxInput source="roles" validate={required()} adminUser={adminUser} />
      <ArrayInput source="transportationOfficeAssignments" label="Transportation Offices (Maximum: 2)">
        <SimpleFormIterator inline key={1} disableReordering>
          <ReferenceInput
            label="Transportation Office"
            reference="offices"
            source="transportationOfficeId"
            perPage={500}
            validate={required()}
          >
            <AutocompleteInput optionText="name" validate={required()} sx={{ width: 325 }} />
          </ReferenceInput>
          <BooleanInput source="primaryOffice" label="Primary Office" defaultValue={false} />
        </SimpleFormIterator>
      </ArrayInput>
      <TextInput source="createdAt" disabled />
      <TextInput source="updatedAt" disabled />
    </SimpleForm>
  </Edit>
);

function mapStateToProps(state) {
  return {
    adminUser: selectAdminUser(state),
  };
}

export default connect(mapStateToProps)(OfficeUserEdit);
