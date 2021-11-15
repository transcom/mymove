import React, { useRef } from 'react';
import PropTypes from 'prop-types';
import { v4 as uuidv4 } from 'uuid';
import { Fieldset } from '@trussworks/react-uswds';

import TextField from 'components/form/fields/TextField';
import { DropdownInput } from 'components/form/fields/DropdownInput';

const statesList = [
  { value: 'AL', key: 'AL' },
  { value: 'AK', key: 'AK' },
  { value: 'AR', key: 'AR' },
  { value: 'AZ', key: 'AZ' },
  { value: 'CA', key: 'CA' },
  { value: 'CO', key: 'CO' },
  { value: 'CT', key: 'CT' },
  { value: 'DC', key: 'DC' },
  { value: 'DE', key: 'DE' },
  { value: 'FL', key: 'FL' },
  { value: 'GA', key: 'GA' },
  { value: 'HI', key: 'HI' },
  { value: 'IA', key: 'IA' },
  { value: 'ID', key: 'ID' },
  { value: 'IL', key: 'IL' },
  { value: 'IN', key: 'IN' },
  { value: 'KS', key: 'KS' },
  { value: 'KY', key: 'KY' },
  { value: 'LA', key: 'LA' },
  { value: 'MA', key: 'MA' },
  { value: 'MD', key: 'MD' },
  { value: 'ME', key: 'ME' },
  { value: 'MI', key: 'MI' },
  { value: 'MN', key: 'MN' },
  { value: 'MO', key: 'MO' },
  { value: 'MS', key: 'MS' },
  { value: 'MT', key: 'MT' },
  { value: 'NC', key: 'NC' },
  { value: 'ND', key: 'ND' },
  { value: 'NE', key: 'NE' },
  { value: 'NH', key: 'NH' },
  { value: 'NJ', key: 'NJ' },
  { value: 'NM', key: 'NM' },
  { value: 'NV', key: 'NV' },
  { value: 'NY', key: 'NY' },
  { value: 'OH', key: 'OH' },
  { value: 'OK', key: 'OK' },
  { value: 'OR', key: 'OR' },
  { value: 'PA', key: 'PA' },
  { value: 'RI', key: 'RI' },
  { value: 'SC', key: 'SC' },
  { value: 'SD', key: 'SD' },
  { value: 'TN', key: 'TN' },
  { value: 'TX', key: 'TX' },
  { value: 'UT', key: 'UT' },
  { value: 'VA', key: 'VA' },
  { value: 'VT', key: 'VT' },
  { value: 'WA', key: 'WA' },
  { value: 'WI', key: 'WI' },
  { value: 'WV', key: 'WV' },
  { value: 'WY', key: 'WY' },
];

export const AddressFields = ({ legend, className, name, render, validators }) => {
  const addressFieldsUUID = useRef(uuidv4());

  return (
    <Fieldset legend={legend} className={className}>
      {render(
        <>
          <TextField
            label="Address 1"
            id={`mailingAddress1_${addressFieldsUUID.current}`}
            name={`${name}.street_address_1`}
            validate={validators?.streetAddress1}
          />
          <TextField
            label="Address 2"
            labelHint="Optional"
            id={`mailingAddress2_${addressFieldsUUID.current}`}
            name={`${name}.street_address_2`}
            validate={validators?.streetAddress2}
          />
          <TextField
            label="City"
            id={`city_${addressFieldsUUID.current}`}
            name={`${name}.city`}
            validate={validators?.city}
          />

          <div className="grid-row grid-gap">
            <div className="mobile-lg:grid-col-6">
              <DropdownInput
                name={`${name}.state`}
                id={`state_${addressFieldsUUID.current}`}
                label="State"
                options={statesList}
                validate={validators?.state}
              />
            </div>
            <div className="mobile-lg:grid-col-6">
              <TextField
                label="ZIP"
                id={`zip_${addressFieldsUUID.current}`}
                name={`${name}.postal_code`}
                maxLength={10}
                validate={validators?.postalCode}
              />
            </div>
          </div>
        </>,
      )}
    </Fieldset>
  );
};

AddressFields.propTypes = {
  legend: PropTypes.node,
  className: PropTypes.string,
  name: PropTypes.string.isRequired,
  render: PropTypes.func,
  validators: PropTypes.shape({
    streetAddress1: PropTypes.func,
    streetAddress2: PropTypes.func,
    city: PropTypes.func,
    state: PropTypes.func,
    postalCode: PropTypes.func,
  }),
};

AddressFields.defaultProps = {
  legend: '',
  className: '',
  render: (fields) => fields,
  validators: {},
};

export default AddressFields;
