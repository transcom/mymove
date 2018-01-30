import React from 'react';
import { Field, reduxForm } from 'redux-form';
import { getUiSchema } from './fields';
import './DD1299.css';

const isEmpty = obj =>
  Object.keys(obj).length === 0 && obj.constructor === Object;
const renderGroupOrField = (name, fields, uiSchema) => {
  /*TODO:
   handle enums
   telephone numbers
   textbox vs textarea (e.g for addresses)
   dates look wonky
   styling in accordance with USWDS
   formatting of labels?
   validate group names don't colide with field names
   tests!!!
  */
  const group = uiSchema.groups[name];
  if (group) {
    const keys = group.fields;
    return (
      <fieldset key={name}>
        <legend htmlFor={name}>{group.title}</legend>
        {keys.map(f => renderGroupOrField(f, fields, uiSchema))}
      </fieldset>
    );
  }
  return renderField(name, fields);
};
const renderField = (name, fields) => {
  const field = fields[name];
  if (!field) {
    return;
  }
  return (
    <div key={name}>
      <label htmlFor={name}>
        {field.title || name}
        <Field name={name} component="input" type={field.format} />
      </label>
    </div>
  );
};

const renderSchema = (schema, uiSchema) => {
  if (!isEmpty(schema)) {
    const fields = schema.properties || [];
    return uiSchema.order.map(i => renderGroupOrField(i, fields, uiSchema));
  }
};
const DD1299 = props => {
  const { handleSubmit, schema } = props;
  const title = schema ? schema.title : '';
  const ui = getUiSchema();
  debugger;
  return (
    <form className="dd1299" onSubmit={handleSubmit}>
      <h1>{title}</h1>
      {renderSchema(schema, ui)}
      <button type="submit">Submit</button>
    </form>
  );
};

//todo: may want to push this to parent, since there is no 1299 specific code here
export default reduxForm({
  // a unique name for the form
  form: 'DD1299',
})(DD1299);
