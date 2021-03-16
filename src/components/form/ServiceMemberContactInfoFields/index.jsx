import React from 'react';
import { func, node, string } from 'prop-types';
import { v4 as uuidv4 } from 'uuid';
import { Checkbox, Label, Fieldset } from '@trussworks/react-uswds';

import styles from './ServiceMemberContactInfoFields.module.scss';

import TextField from 'components/form/fields/TextField';

export const ServiceMemberContactInfoFields = ({
  legend,
  className,
  onChangePreferEmail,
  onChangePreferPhone,
  name,
  render,
}) => {
  const contactInfoFieldsetUUID = uuidv4();

  return (
    <Fieldset legend={legend} className={className}>
      {render(
        <div className={styles.ServiceMemberContactInfoFields}>
          <TextField
            className={styles.contactPhoneFields}
            label="Best contact phone"
            id={`phone_${contactInfoFieldsetUUID}`}
            data-testid="phone"
            name={`${name}.phone`}
            type="tel"
            maxLength="10"
          />
          <TextField
            className={styles.contactPhoneFields}
            label="Alt. phone"
            labelHint="Optional"
            id={`alternatePhone_${contactInfoFieldsetUUID}`}
            data-testid="alternamte-phone"
            name={`${name}.alternatePhone`}
            type="tel"
            maxLength="10"
          />
          <TextField
            label="Personal email"
            id={`email_${contactInfoFieldsetUUID}`}
            data-testid="email"
            name={`${name}.email`}
          />
          <Label>Preferred contact method</Label>
          <Checkbox
            id={`preferPhone_${contactInfoFieldsetUUID}`}
            label="Phone"
            name={`${name}.prefer_phone`}
            onChange={onChangePreferPhone}
          />
          <Checkbox
            id={`prefer_email${contactInfoFieldsetUUID}`}
            label="Email"
            name={`${name}.prefer_email`}
            onChange={onChangePreferEmail}
          />
        </div>,
      )}
    </Fieldset>
  );
};

ServiceMemberContactInfoFields.propTypes = {
  legend: node,
  className: string,
  name: string.isRequired,
  render: func,
  onChangePreferPhone: func.isRequired,
  onChangePreferEmail: func.isRequired,
};

ServiceMemberContactInfoFields.defaultProps = {
  legend: '',
  className: '',
  render: (fields) => fields,
};

export default ServiceMemberContactInfoFields;
