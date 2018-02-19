import React from 'react';
import PropTypes from 'prop-types';
import SchemaField from './JsonSchemaField';

import { reduxForm } from 'redux-form';
import './index.css';

const isEmpty = obj =>
  Object.keys(obj).length === 0 && obj.constructor === Object;
const renderGroupOrField = (fieldName, fields, uiSchema, nameSpace) => {
  /*TODO:
   telephone numbers/ pattern validation
   textbox vs textarea (e.g for addresses)
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

const renderSchema = (schema, uiSchema, nameSpace = '') => {
  if (schema && !isEmpty(schema)) {
    console.log('IOIOIOI');
    console.log(schema);

    // Mark all the required fields as required.
    if (schema.required) {
      schema.required.forEach(requiredFieldName => {
        console.log(requiredFieldName);
        schema.properties[requiredFieldName][
          SchemaField.IS_REQUIRED_KEY
        ] = true;
      });
    }

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

export const reduxifyForm = name => reduxForm({ form: name })(JsonSchemaForm);
