import React, { useRef } from 'react';
import { PropTypes, shape } from 'prop-types';
import { v4 as uuidv4 } from 'uuid';
import { Fieldset } from '@trussworks/react-uswds';

import { statesList } from '../../../constants/states';

import Hint from 'components/Hint/index';
import styles from 'components/form/AddressFields/AddressFields.module.scss';
import { technicalHelpDeskURL } from 'shared/constants';
import TextField from 'components/form/fields/TextField/TextField';
import { DropdownInput } from 'components/form/fields/DropdownInput';
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
  locationLookup,
  formikProps: { setFieldTouched, setFieldValue },
  labelHint: labelHintProp,
  address1LabelHint,
}) => {
  const addressFieldsUUID = useRef(uuidv4());
  const infoStr = 'If you encounter any inaccurate lookup information please contact the ';
  const assistanceStr = ' for further assistance.';

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
    setFieldValue(`${name}.city`, value.city).then(() => {
      setFieldTouched(`${name}.city`, false);
    });
    setFieldValue(`${name}.state`, value.state).then(() => {
      setFieldTouched(`${name}.state`, false);
    });
    setFieldValue(`${name}.county`, value.county).then(() => {
      setFieldTouched(`${name}.county`, false);
    });
    setFieldValue(`${name}.postalCode`, value.postalCode).then(() => {
      setFieldTouched(`${name}.postalCode`, false);
    });
    setFieldValue(`${name}.usPostRegionCitiesID`, value.usPostRegionCitiesID).then(() => {
      setFieldTouched(`${name}.usPostRegionCitiesID`, true);
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
          {locationLookup && (
            <>
              <LocationInput
                name={`${name}-location`}
                placeholder="Start typing a Zip or City, State Zip"
                label="Location Lookup"
                handleLocationChange={handleOnLocationChange}
              />
              <Hint className={styles.hint} id="locationInfo" data-testid="locationInfo">
                {infoStr}
                <a href={technicalHelpDeskURL} target="_blank" rel="noreferrer">
                  Technical Help Desk
                </a>
                {assistanceStr}
              </Hint>
            </>
          )}
          <div className={styles.container}>
            <div className={styles.column}>
              {!locationLookup && (
                <TextField
                  label="City"
                  id={`city_${addressFieldsUUID.current}`}
                  name={`${name}.city`}
                  labelHint={labelHintProp}
                  data-testid={`${name}.city`}
                  validate={validators?.city}
                />
              )}
              {locationLookup && (
                <>
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
                </>
              )}
            </div>
            <div className={styles.column}>
              {!locationLookup && (
                <DropdownInput
                  name={`${name}.state`}
                  data-testid={`${name}.state`}
                  id={`state_${addressFieldsUUID.current}`}
                  label="State"
                  labelHint={labelHintProp}
                  options={statesList}
                  validate={validators?.state}
                />
              )}
              <TextField
                label="ZIP"
                id={`zip_${addressFieldsUUID.current}`}
                name={`${name}.postalCode`}
                data-testid={`${name}.postalCode`}
                maxLength={10}
                labelHint={labelHintProp}
                display={!locationLookup ? '' : 'readonly'}
                validate={validators?.postalCode}
              />
              {locationLookup && (
                <TextField
                  label="County"
                  id={`county_${addressFieldsUUID.current}`}
                  name={`${name}.county`}
                  labelHint={labelHintProp}
                  data-testid={`${name}.county`}
                  display="readonly"
                  validate={validators?.county}
                />
              )}
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
  locationLookup: PropTypes.bool,
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
};

AddressFields.defaultProps = {
  legend: '',
  className: '',
  render: (fields) => fields,
  locationLookup: false,
  validators: {},
  address1LabelHint: null,
  formikProps: {},
};

export default AddressFields;
