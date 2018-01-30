import React from 'react';
import DD1299Form from './DD1299Form';
import TestForm from './jsonschema';
import { GetSpec } from 'shared/api';

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
    if (true)
      return <DD1299Form onSubmit={this.submit} schema={this.state.schema} />;
    return <TestForm schema={this.state.schema} />;
  }
}
