import React, { Component } from 'react';
import Form from 'react-jsonschema-form';
import { GetSpec } from 'shared/api';

const log = type => console.log.bind(console, type);

export default class Alt extends Component {
  constructor(props) {
    super(props);
    this.state = {
      schema: {},
    };
  }
  componentDidMount() {
    const that = this;
    GetSpec().then(spec => {
      that.setState({ schema: spec.definitions.CreateForm1299Payload });
    });
  }
  render() {
    return (
      <Form
        schema={this.state.schema}
        onChange={log('changed')}
        onSubmit={log('submitted')}
        onError={log('errors')}
      />
    );
  }
}
