import React, { useRef } from 'react';
import { PropTypes, shape } from 'prop-types';
import { v4 as uuidv4 } from 'uuid';
import { Fieldset } from '@trussworks/react-uswds';

import { requiredAsteriskMessage } from '../RequiredAsterisk';

import Hint from 'components/Hint/index';
import styles from 'components/form/AddressFields/AddressFields.module.scss';
import { technicalHelpDeskURL } from 'shared/constants';
import TextField from 'components/form/fields/TextField/TextField';
import LocationInput from 'components/form/fields/LocationInput';

/**
 * @param legend
 * @param className
 * @param name
 * @param render
 * @param validators
 * @param zipCity
 * @param address1LabelHint string to override display labelHint if street 1 is Optional/Required per context.
 * This is specifically designed to handle unique display between customer and office/prime sim for address 1.
 * @return {JSX.Element}
 * @constructor
 */
export const AddressFields = ({
  legend,
  className,
  name,
  render,
  validators,
  formikProps: { setFieldTouched, setFieldValue },
  labelHint: labelHintProp,
  address1LabelHint,
  includePOBoxes,
}) => {
  const addressFieldsUUID = useRef(uuidv4());
  const infoStr = 'If you encounter any inaccurate lookup information please contact the ';
  const assistanceStr = ' for further assistance.';

  const handleOnLocationChange = (value) => {
    const city = value ? value.city : null;
    const state = value ? value.state : null;
    const county = value ? value.county : null;
    const postalCode = value ? value.postalCode : null;
    const usPostRegionCitiesID = value ? value.usPostRegionCitiesID : null;

    setFieldValue(`${name}.city`, city).then(() => {
      setFieldTouched(`${name}.city`, false);
    });
    setFieldValue(`${name}.state`, state).then(() => {
      setFieldTouched(`${name}.state`, false);
    });
    setFieldValue(`${name}.county`, county).then(() => {
      setFieldTouched(`${name}.county`, false);
    });
    setFieldValue(`${name}.postalCode`, postalCode).then(() => {
      setFieldTouched(`${name}.postalCode`, false);
    });
    setFieldValue(`${name}.usPostRegionCitiesID`, usPostRegionCitiesID).then(() => {
      setFieldTouched(`${name}.usPostRegionCitiesID`, true);
    });
  };

  // E-05732: for PPMs, the destination address street 1 is now optional except for closeout
  // this field is usually always required other than PPMs
  // a value for address1LabelHint is passed in when we want address 1 to be optional
  const showRequiredAsteriskForAddress1 = address1LabelHint === null || labelHintProp === 'Required';

  return (
    <Fieldset legend={legend} className={className}>
      {requiredAsteriskMessage}
      {render(
        <>
          <TextField
            label="Address 1"
            id={`mailingAddress1_${addressFieldsUUID.current}`}
            name={`${name}.streetAddress1`}
            required={showRequiredAsteriskForAddress1}
            showRequiredAsterisk={showRequiredAsteriskForAddress1}
            data-testid={`${name}.streetAddress1`}
            validate={validators?.streetAddress1}
          />
          <TextField
            label="Address 2"
            id={`mailingAddress2_${addressFieldsUUID.current}`}
            name={`${name}.streetAddress2`}
            data-testid={`${name}.streetAddress2`}
            validate={validators?.streetAddress2}
          />
          <TextField
            label="Address 3"
            id={`mailingAddress3_${addressFieldsUUID.current}`}
            name={`${name}.streetAddress3`}
            data-testid={`${name}.streetAddress3`}
            validate={validators?.streetAddress3}
          />
          <LocationInput
            name={`${name}`}
            placeholder="Start typing a Zip or City, State Zip"
            label="Location Lookup"
            handleLocationChange={handleOnLocationChange}
            includePOBoxes={includePOBoxes}
          />

          <Hint className={styles.hint} id="locationInfo" data-testid="locationInfo">
            {infoStr}
            <a href={technicalHelpDeskURL} target="_blank" rel="noreferrer">
              Technical Help Desk
            </a>
            {assistanceStr}
          </Hint>
          <div className={styles.container}>
            <div className={styles.column}>
              <TextField
                label="City"
                id={`city_${addressFieldsUUID.current}`}
                name={`${name}.city`}
                showRequiredAsterisk
                required
                data-testid={`${name}.city`}
                display="readonly"
                validate={validators?.city}
              />
              <TextField
                label="State"
                id={`state_${addressFieldsUUID.current}`}
                name={`${name}.state`}
                data-testid={`${name}.state`}
                showRequiredAsterisk
                required
                display="readonly"
                validate={validators?.state}
                styles="margin-top: 1.5em"
              />
            </div>
            <div className={styles.column}>
              <TextField
                label="ZIP"
                id={`zip_${addressFieldsUUID.current}`}
                name={`${name}.postalCode`}
                data-testid={`${name}.postalCode`}
                maxLength={10}
                showRequiredAsterisk
                required
                display="readonly"
                validate={validators?.postalCode}
              />
              <TextField
                label="County"
                id={`county_${addressFieldsUUID.current}`}
                name={`${name}.county`}
                showRequiredAsterisk
                required
                data-testid={`${name}.county`}
                display="readonly"
                validate={validators?.county}
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
    county: PropTypes.func,
    usPostRegionCitiesID: PropTypes.func,
  }),
  address1LabelHint: PropTypes.string,
  formikProps: shape({
    touched: shape({}),
    errors: shape({}),
    setFieldTouched: PropTypes.func,
    setFieldValue: PropTypes.func,
  }),
  includePOBoxes: PropTypes.bool,
};

AddressFields.defaultProps = {
  legend: '',
  className: '',
  render: (fields) => fields,
  validators: {},
  address1LabelHint: null,
  formikProps: {},
  includePOBoxes: false,
};

export default AddressFields;
