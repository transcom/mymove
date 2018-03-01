import React from 'react';
import PropTypes from 'prop-types';
import SchemaField from './JsonSchemaField';

import { isEmpty } from 'lodash';
import { reduxForm } from 'redux-form';
import './index.css';

const renderGroupOrField = (fieldName, fields, uiSchema, nameSpace) => {
  /*TODO:
   dates look wonky in chrome
   styling in accordance with USWDS
   validate group names don't colide with field names
  */
  const group = uiSchema.groups && uiSchema.groups[fieldName];
  const isRef =
    fields[fieldName] &&
    fields[fieldName].$$ref &&
    fields[fieldName].properties;
  if (group) {
    const keys = group.fields;
    return (
      <fieldset key={fieldName}>
        <legend htmlFor={fieldName}>{group.title}</legend>
        {keys.map(f => renderGroupOrField(f, fields, uiSchema, nameSpace))}
      </fieldset>
    );
  } else if (isRef) {
    const refName = fields[fieldName].$$ref.split('/').pop();
    const refSchema = uiSchema.definitions[refName];
    return renderSchema(fields[fieldName], refSchema, fieldName);
  }
  return renderField(fieldName, fields, nameSpace);
};

const renderField = (fieldName, fields, nameSpace) => {
  const field = fields[fieldName];
  if (!field) {
    return;
  }
  return SchemaField.createSchemaField(fieldName, field, nameSpace);
};

// Because we have nested objects it's possible to have
// An object that is not-required that itself has required properties. This makes sense, in that
// If the entire object is omitted (say, an address) then the form is valid, but if a
// single property of the object is included, then all its required properties must be
// as well.
// Therefore, the rules for wether or not a field is required are:
// 1. If it is listed in the top level definition, it's required.
// 2. If it is required and it is an object, its required fields are required
// 3. If it is an object and some value in it has been set, then all it's required fields must be set too
// This is a recusive definition.
export const recursivelyValidateRequiredFields = (values, spec) => {
  let requiredErrors = {};
  // first, check that all required fields are present
  if (spec.required) {
    spec.required.forEach(requiredFieldName => {
      if (values[requiredFieldName] === undefined) {
        // check if the required thing is a object, in that case put it on its required fields. Otherwise recurse.
        let schemaForKey = spec.properties[requiredFieldName];
        if (schemaForKey) {
          if (schemaForKey.type === 'object') {
            let subErrors = recursivelyValidateRequiredFields({}, schemaForKey);
            if (!isEmpty(subErrors)) {
              requiredErrors[requiredFieldName] = subErrors;
            }
          } else {
            requiredErrors[requiredFieldName] = 'Required.';
          }
        } else {
          console.error('The schema should have all required fields in it.');
        }
      }
    });
  }

  // now go through every existing value, if its an object, we must recurse to see if its required properties are there.
  Object.keys(values).forEach(function(key) {
    let schemaForKey = spec.properties[key];
    if (schemaForKey) {
      if (schemaForKey.type === 'object') {
        let subErrors = recursivelyValidateRequiredFields(
          values[key],
          schemaForKey,
        );
        if (!isEmpty(subErrors)) {
          requiredErrors[key] = subErrors;
        }
      }
    } else {
      console.error('The schema should have fields for all present values..');
    }
  });

  return requiredErrors;
};

// To validate that fields are required, we look at the list of top level required
// fields and then validate them and their children.
const validateRequiredFields = (values, form, somethingelse, andhow) => {
  const swaggerSpec = form.schema;
  let requiredErrors;
  if (swaggerSpec && !isEmpty(swaggerSpec)) {
    requiredErrors = recursivelyValidateRequiredFields(values, swaggerSpec);
  }
  return requiredErrors;
};

const renderSchema = (schema, uiSchema, nameSpace = '') => {
  if (schema && !isEmpty(schema)) {
    const fields = schema.properties || [];
    return uiSchema.order.map(i =>
      renderGroupOrField(i, fields, uiSchema, nameSpace),
    );
  }
};
const JsonSchemaForm = props => {
  const { pristine, submitting, invalid } = props;
  const { handleSubmit, schema, showSubmit } = props;
  const title = schema ? schema.title : '';
  const uiSchema = props.subsetOfUiSchema
    ? Object.assign({}, props.uiSchema, {
        order: props.subsetOfUiSchema,
      })
    : props.uiSchema;
  return (
    <form className="default" onSubmit={handleSubmit}>
      <h1>{title}</h1>
      {renderSchema(schema, uiSchema)}
      {showSubmit && (
        <button type="submit" disabled={pristine || submitting || invalid}>
          Submit
        </button>
      )}
    </form>
  );
};

JsonSchemaForm.propTypes = {
  schema: PropTypes.object.isRequired,
  uiSchema: PropTypes.object.isRequired,
  handleSubmit: PropTypes.func.isRequired,
  showSubmit: PropTypes.bool,
  subsetOfUiSchema: PropTypes.arrayOf(PropTypes.string),
};

JsonSchemaForm.defaultProps = {
  showSubmit: true,
};

export const reduxifyForm = name =>
  reduxForm({ form: name, validate: validateRequiredFields })(JsonSchemaForm);
