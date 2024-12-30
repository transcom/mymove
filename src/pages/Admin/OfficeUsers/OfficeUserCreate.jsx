import React from 'react';
import {
  Create,
  SimpleForm,
  TextInput,
  ReferenceInput,
  AutocompleteInput,
  ArrayInput,
  SimpleFormIterator,
  BooleanInput,
  SelectInput,
  useDataProvider,
  Button,
} from 'react-admin';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import { RolesPrivilegesCheckboxInput } from 'scenes/SystemAdmin/shared/RolesPrivilegesCheckboxes';
import { roleTypes } from 'constants/userRoles';

const OfficeUserCreate = () => {
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
    <Create>
      <SimpleForm
        sx={{ '& .MuiInputBase-input': { width: 232 } }}
        mode="onSubmit" // Array input validation doesn't update properly until the form is submitted due to a bug in a dependancy of the version of react-admin in use
        reValidateMode="onSubmit"
        validate={validateForm}
      >
        <TextInput source="firstName" mode="onBlur" />
        <TextInput source="middleInitials" />
        <TextInput source="lastName" mode="onBlur" />
        <TextInput source="email" type="email" mode="onBlur" />
        <TextInput source="telephone" mask="999-999-9999" mode="onBlur" />
        <RolesPrivilegesCheckboxInput source="roles" mode="onBlur" />
        <ArrayInput
          source="transportationOfficeAssignments"
          label="Transportation Office Assignments (Maximum: 2)"
          defaultValue={[{ primaryOffice: true }]}
        >
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
      </SimpleForm>
    </Create>
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

export default OfficeUserCreate;
