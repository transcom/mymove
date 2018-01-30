import React from 'react';
import Form from 'react-jsonschema-form';

const log = type => console.log.bind(console, type);
const uiSchema = {
  mobile_home_services_requested: {
    'ui:widget': 'checkboxes',
  },
};

const render = props => (
  <Form
    schema={props.schema}
    uiSchema={uiSchema}
    onChange={log('changed')}
    onSubmit={log('submitted')}
    onError={log('errors')}
  />
);

export default render;
