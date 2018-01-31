import React from 'react';
import { reduxifyForm } from 'shared/JsonSchemaForm';
import { GetSpec } from 'shared/api';
import { getUiSchema } from './uiSchema';

const DD1299Form = reduxifyForm('DD1299');
export default class DD1299 extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      schema: {},
      fields: [],
    };
  }
  componentDidMount() {
    const that = this;
    GetSpec().then(spec => {
      that.setState({
        schema: spec.definitions.CreateForm1299Payload,
      });
    });
  }
  submit = values => {
    // print the form values to the console
    console.log(values);
  };
  render() {
    const uiSchema = getUiSchema();
    return (
      <DD1299Form
        onSubmit={this.submit}
        schema={this.state.schema}
        uiSchema={uiSchema}
      />
    );
  }
}
