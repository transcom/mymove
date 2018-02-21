import React from 'react';
import PropTypes from 'prop-types';
import SchemaField from './JsonSchemaField';

import { isEmpty } from 'lodash';
import { reduxForm } from 'redux-form';
import './index.css';

const renderGroupOrField = (fieldName, fields, uiSchema, nameSpace) => {
  /*TODO:
   SSN/ pattern validation
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

const recursivelyValidateRequiredFields = (values, spec) => {
  let requiredErrors = {};
  // first, check that all required fields are present
  if (spec.required) {
    spec.required.forEach(requiredFieldName => {
      if (values[requiredFieldName] === undefined) {
        // check if the required thing is a ref, in that case put it on its required fields. Otherwise recurse.
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

  // now go through every existing value, if its an object, we gotta recurse to see if its required properties are there.
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
  // gotta start testing.
  // should be easy, tests are: a schema, a uischema, a data, and wether it's valid?(or the errors hash)
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
  const { handleSubmit, schema, uiSchema } = props;
  const title = schema ? schema.title : '';
  return (
    <form className="default" onSubmit={handleSubmit}>
      <h1>{title}</h1>
      {renderSchema(schema, uiSchema)}
      <button type="submit" disabled={pristine || submitting || invalid}>
        Submit
      </button>
    </form>
  );
};

JsonSchemaForm.propTypes = {
  schema: PropTypes.object.isRequired,
  uiSchema: PropTypes.object.isRequired,
  handleSubmit: PropTypes.func.isRequired,
};

export const reduxifyForm = name =>
  reduxForm({ form: name, validate: validateRequiredFields })(JsonSchemaForm);
