import React, { useRef, useEffect, useState } from 'react';
import { PropTypes, shape } from 'prop-types';
import { v4 as uuidv4 } from 'uuid';
import { Fieldset } from '@trussworks/react-uswds';
import { useFormikContext } from 'formik';

import Hint from 'components/Hint/index';
import styles from 'components/form/AddressFields/AddressFields.module.scss';
import { technicalHelpDeskURL } from 'shared/constants';
import TextField from 'components/form/fields/TextField/TextField';
import LocationInput from 'components/form/fields/LocationInput';
import CountryInput from 'components/form/fields/CountryInput';

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
}) => {
  const addressFieldsUUID = useRef(uuidv4());
  const infoStr = 'If you encounter any inaccurate lookup information please contact the ';
  const assistanceStr = ' for further assistance.';

  const { values } = useFormikContext();

  // some 'name' are passed in like pickupAddress.address and we cannot use values to get pickupAddress.address, so we
  // remove the .address and get pickupAddress then get the country code from that. If there is not a .address we just
  // use the 'name' passed in as is
  const separator = '.';
  const separatorIndex = name.indexOf(separator);

  let addressName;
  if (separatorIndex !== -1) {
    addressName = name.substring(0, separatorIndex);
  } else {
    addressName = name;
  }

  const nameValue = values[addressName];
  let nameCountryCode = '';

<<<<<<< HEAD
  if (name.includes('.address') && nameValue?.address !== undefined && nameValue?.address.city !== '') {
    nameCountryCode = nameValue.address.country ? nameValue.address.country.code : '';
  } else if (!name.includes('.address') && nameValue?.country !== undefined) {
=======
  if (name.includes('.address') && nameValue !== undefined && nameValue.address !== undefined) {
    nameCountryCode = nameValue.address.country ? nameValue.address.country.code : '';
  } else if (!name.includes('.address') && nameValue !== undefined && nameValue.country !== undefined) {
>>>>>>> 246744c45567b63e0309b2cd2c2320ec54f9b9d0
    nameCountryCode = nameValue.country.code;
  }

  const [currentCountryCode, setCurrentCountryCode] = useState(nameCountryCode);

  useEffect(() => {
    setCurrentCountryCode(nameCountryCode);
  }, [nameCountryCode]);

  const getAddress1LabelHintText = (labelHint, address1Label) => {
    if (address1Label === null) {
      return labelHint;
    }

    // Override default and use what is passed in.
    if (address1Label && address1Label.trim().length > 0) {
      return address1Label;
    }

    return null;
  };

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

  const handleOnCountryChange = (value) => {
    const countryID = value ? value.id : null;
    const countryName = value ? value.name : null;
    const countryCode = value ? value.code : null;

    if (countryID == null) {
      handleOnLocationChange(null);
    }

    setCurrentCountryCode(countryCode);
    setFieldValue(`${name}.country.id`, countryID).then(() => {
      setFieldTouched(`${name}.country.id`, false);
    });
    setFieldValue(`${name}.country.code`, countryCode).then(() => {
      setFieldTouched(`${name}.country.code`, false);
    });
    setFieldValue(`${name}.country.name`, countryName).then(() => {
      setFieldTouched(`${name}.country.name`, false);
    });
    setFieldValue(`${name}.countryID`, countryID).then(() => {
      setFieldTouched(`${name}.countryID`, true);
    });
  };

  return (
    <Fieldset legend={legend} className={className}>
      {render(
        <>
          <TextField
            label="Address 1"
            id={`mailingAddress1_${addressFieldsUUID.current}`}
            name={`${name}.streetAddress1`}
            labelHint={getAddress1LabelHintText(labelHintProp, address1LabelHint)}
            data-testid={`${name}.streetAddress1`}
            validate={validators?.streetAddress1}
          />
          <TextField
            label="Address 2"
            labelHint={labelHintProp ? null : 'Optional'}
            id={`mailingAddress2_${addressFieldsUUID.current}`}
            name={`${name}.streetAddress2`}
            data-testid={`${name}.streetAddress2`}
            validate={validators?.streetAddress2}
          />
          <TextField
            label="Address 3"
            labelHint={labelHintProp ? null : 'Optional'}
            id={`mailingAddress3_${addressFieldsUUID.current}`}
            name={`${name}.streetAddress3`}
            data-testid={`${name}.streetAddress3`}
            validate={validators?.streetAddress3}
          />

          <CountryInput
            name={`${name}`}
            placeholder="Start typing a country name, code"
            label="Country Lookup"
            handleCountryChange={handleOnCountryChange}
          />

          <LocationInput
            name={`${name}`}
            placeholder="Start typing a Zip or City, State Zip"
            label="Location Lookup"
            handleLocationChange={handleOnLocationChange}
            isDisabled={currentCountryCode === null || currentCountryCode === '' || currentCountryCode !== 'US'}
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
                labelHint={labelHintProp}
                data-testid={`${name}.city`}
                display="readonly"
                validate={validators?.city}
              />
              <TextField
                label="State"
                id={`state_${addressFieldsUUID.current}`}
                name={`${name}.state`}
                data-testid={`${name}.state`}
                labelHint={labelHintProp}
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
                labelHint={labelHintProp}
                display="readonly"
                validate={validators?.postalCode}
              />
              <TextField
                label="County"
                id={`county_${addressFieldsUUID.current}`}
                name={`${name}.county`}
                labelHint={labelHintProp}
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
    countryID: PropTypes.func,
  }),
  address1LabelHint: PropTypes.string,
  formikProps: shape({
    touched: shape({}),
    errors: shape({}),
    setFieldTouched: PropTypes.func,
    setFieldValue: PropTypes.func,
  }),
};

AddressFields.defaultProps = {
  legend: '',
  className: '',
  render: (fields) => fields,
  validators: {},
  address1LabelHint: null,
  formikProps: {},
};

export default AddressFields;
