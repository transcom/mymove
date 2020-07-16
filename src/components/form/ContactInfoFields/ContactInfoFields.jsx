import React from 'react';
import PropTypes from 'prop-types';
import { Fieldset, Label } from '@trussworks/react-uswds';

import { TextInput } from 'components/form/fields/TextInput';

export const ContactInfoFields = ({ contactType, legend }) => {
  return (
    <Fieldset legend={legend}>
      <TextInput
        label="First name"
        hint=" (optional)"
        id={`${contactType}-first-name`}
        name={`${contactType}-first-name`}
        type="text"
      />

      <TextInput
        label="Last name"
        hint=" (optional)"
        id={`${contactType}-last-name`}
        name={`${contactType}-last-name`}
        type="text"
      />

      <TextInput
        label="Phone"
        hint=" (optional)"
        id={`${contactType}-phone`}
        name={`${contactType}-phone`}
        type="text"
      />
      <TextInput
        label="Email"
        hint=" (optional)"
        id={`${contactType}-email`}
        name={`${contactType}-email`}
        type="text"
      />
    </Fieldset>
  );
};

ContactInfoFields.propTypes = {
  contactType: PropTypes.string,
  legend: PropTypes.string,
};

ContactInfoFields.defaultProps = {
  contactType: '',
  legend: '',
};

export default ContactInfoFields;
