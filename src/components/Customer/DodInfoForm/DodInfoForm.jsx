import React from 'react';
import { Field } from 'formik';

import { DropdownInput, TextInput } from 'components/form/fields';
import { Form } from 'components/form/Form';
import { dropdownInputOptions } from 'shared/formatters';
import { SERVICE_MEMBER_AGENCY_LABELS, rankLabels } from 'content/serviceMemberAgencies';

const DodInfoForm = () => {
  const affiliationOptions = dropdownInputOptions(SERVICE_MEMBER_AGENCY_LABELS);
  const rankOptions = dropdownInputOptions(rankLabels);

  return (
    <Form>
      <Field as={DropdownInput} label="Branch of service" name="affiliation" options={affiliationOptions} />
      <Field as={TextInput} label="DoD ID number" name="edipi" />
      <Field as={DropdownInput} label="Rank" name="rank" options={rankOptions} />
    </Form>
  );
};

export default DodInfoForm;
